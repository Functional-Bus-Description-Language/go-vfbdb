import sys
import traceback

import cosim
import vfbdb


WRITE_FIFO_PATH = sys.argv[1]
READ_FIFO_PATH = sys.argv[2]

iface = cosim.Iface(WRITE_FIFO_PATH, READ_FIFO_PATH)

try:
    Main = vfbdb.Main(iface)

    value = 2 ** vfbdb.mainPkg.WIDTH - 1

    print(f"Writing VALID_VALUE ({value}) to Cfg register")
    Main.Cfg.write(value)

    print("Reading Cfg")
    read_val = Main.Cfg.read()
    if read_val != value:
        raise Exception(f"Read wrong value form Cfg {read_val}")

    iface.end(0)

except Exception as E:
    print(traceback.format_exc())
    iface.end(1)
