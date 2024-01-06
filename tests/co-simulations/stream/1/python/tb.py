import random
import sys
import traceback

import cosim
import vfbdb

WRITE_FIFO_PATH = sys.argv[1]
READ_FIFO_PATH = sys.argv[2]


try:
    iface = cosim.Iface(WRITE_FIFO_PATH, READ_FIFO_PATH)

    Main = vfbdb.Main(iface)

    data = []
    for i in range(vfbdb.mainPkg.DEPTH):
        dataset = []
        dataset.append(random.randint(0, 2 ** Main.Add.params[0]['Width'] - 1))
        dataset.append(random.randint(0, 2 ** Main.Add.params[1]['Width'] - 1))
        dataset.append(random.randint(0, 2 ** Main.Add.params[2]['Width'] - 1))
        data.append(dataset)

    print(f"Writing downstream {vfbdb.mainPkg.DEPTH} times")
    Main.Add.write(data)

    results = Main.Result.read(vfbdb.mainPkg.DEPTH)

    for i, dataset in enumerate(data):
        got = results[i][0]
        want = sum(dataset)
        assert got == want, f"{i}: got {got}, want {want}"

    iface.end(0)

except Exception as E:
    iface.end(1)
    print(traceback.format_exc())
