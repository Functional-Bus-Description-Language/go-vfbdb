import sys

from cosim_interface import CosimInterface
import wbfbd


WRITE_FIFO_PATH = sys.argv[1]
READ_FIFO_PATH = sys.argv[2]

cosim_interface = CosimInterface(WRITE_FIFO_PATH, READ_FIFO_PATH)

try:
    print("\nstarting cosimulation")

    Main = wbfbd.Main(cosim_interface)

    print(f"Writing VALID_VALUE ({wbfbd.mainPkg.VALID_VALUE}) to Cfg register")
    Main.Cfg.write(wbfbd.mainPkg.VALID_VALUE)

    print("Reading Cfg")
    read_val = Main.Cfg.read()
    if read_val != wbfbd.mainPkg.VALID_VALUE:
        raise Exception(f"Read wrong value form Cfg {read_val}")

    print("\nending cosimulation")
    cosim_interface.end(0)

except Exception as E:
    cosim_interface.end(1)
    print(E)
