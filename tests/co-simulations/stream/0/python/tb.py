import random
import traceback
import sys

import cosim
import vfbdb

WRITE_FIFO_PATH = sys.argv[1]
READ_FIFO_PATH = sys.argv[2]


try:
    iface = cosim.Iface(WRITE_FIFO_PATH, READ_FIFO_PATH)

    Main = vfbdb.Main(iface)

    N = random.randint(1, 32)
    print(f"Reading FIFO stream {N} times")
    vals = Main.FIFO.read(N)

    for i, v in enumerate(vals):
        assert i == v[0], f"read {v(0)}, expecting {i}"

    iface.end(0)

except Exception as E:
    iface.end(1)
    print(traceback.format_exc())
