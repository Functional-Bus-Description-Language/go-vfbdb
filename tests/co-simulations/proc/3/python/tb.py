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

    s = random.randint(0, 2 ** 16 - 1)
    c = random.randint(2**33, 2 ** 40 - 1)

    print(f"Calling add function, s = {s}, c = {c}")
    Main.Add(s, c)

    print(f"Reading result")
    result = Main.Result.read()

    if s + c != result:
        print(f"Wrong result, got {result}, expecting {s+c}")
        iface.end(1)

    iface.end(0)

except Exception as E:
    iface.end(1)
    print(traceback.format_exc())
