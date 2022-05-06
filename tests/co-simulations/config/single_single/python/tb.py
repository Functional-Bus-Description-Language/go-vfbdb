import sys
import random

from cosim_interface import CosimInterface
import vfbdb


WRITE_FIFO_PATH = sys.argv[1]
READ_FIFO_PATH = sys.argv[2]

cosim_interface = CosimInterface(WRITE_FIFO_PATH, READ_FIFO_PATH)

try:
    print("\nstarting cosimulation")

    Main = vfbdb.Main(cosim_interface)

    val = random.randint(0, 2 ** 7 - 1)

    print(f"Generated random value: {val}")

    print("Writing Cfg")
    Main.Cfg.write(val)

    print("Reading Cfg")
    read_val = Main.Cfg.read()
    if read_val != val:
        raise Exception(f"Read wrong value form Cfg {read_val}")

    print("Reading St")
    read_val = Main.St.read()
    if read_val != val:
        raise Exception(f"Read wrong value form St {read_val}")

    print("\nending cosimulation")
    cosim_interface.end(0)

except Exception as E:
    cosim_interface.end(1)
    print(E)
