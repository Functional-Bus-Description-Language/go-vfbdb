import sys
import traceback

from cosim_interface import CosimInterface
import vfbdb


WRITE_FIFO_PATH = sys.argv[1]
READ_FIFO_PATH = sys.argv[2]

cosim_interface = CosimInterface(WRITE_FIFO_PATH, READ_FIFO_PATH)

try:
    print("\nstarting cosimulation")

    Main = vfbdb.Main(cosim_interface)

    value = 2 ** vfbdb.mainPkg.WIDTH - 1

    print(f"Writing VALID_VALUE ({value}) to Cfg register")
    Main.Cfg.write(value)

    print("Reading Cfg")
    read_val = Main.Cfg.read()
    if read_val != value:
        raise Exception(f"Read wrong value form Cfg {read_val}")

    print("\nending cosimulation")
    cosim_interface.end(0)

except Exception as E:
    print(traceback.format_exc())
    cosim_interface.end(1)
