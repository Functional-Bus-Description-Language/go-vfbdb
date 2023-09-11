import sys
import random

import cosim
import vfbdb


WRITE_FIFO_PATH = sys.argv[1]
READ_FIFO_PATH = sys.argv[2]

iface = cosim.Iface(WRITE_FIFO_PATH, READ_FIFO_PATH)

try:
    Main = vfbdb.Main(iface)

    val = random.randint(2 ** 33, 2 ** 48  - 1)

    print(f"Generated random value: {val}")

    print("Writing Cfg")
    Main.Cfg.write(val)

    print("Reading Cfg")
    read_val = Main.Cfg.read()
    if read_val != val:
        raise Exception(f"Read wrong value form Cfg {read_val}")

    iface.end(0)

except Exception as E:
    iface.end(1)
    print(E)
