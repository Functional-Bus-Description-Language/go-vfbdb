import sys

from cosim_interface import CosimInterface
import vfbdb


WRITE_FIFO_PATH = sys.argv[1]
READ_FIFO_PATH = sys.argv[2]

cosim_interface = CosimInterface(WRITE_FIFO_PATH, READ_FIFO_PATH)

try:
    print("\nstarting cosimulation")

    Main = vfbdb.Main(cosim_interface)

    print(f"Writing VALID_VALUE ({vfbdb.mainPkg.VALID_VALUE}) to Cfg register")
    Main.Cfg.write(vfbdb.mainPkg.VALID_VALUE)

    print("Reading Cfg")
    read_val = Main.Cfg.read()
    if read_val != vfbdb.mainPkg.VALID_VALUE:
        raise Exception(f"Read wrong value form Cfg {read_val}")

    print("\nending cosimulation")
    cosim_interface.end(0)

except Exception as E:
    cosim_interface.end(1)
    print(E)
