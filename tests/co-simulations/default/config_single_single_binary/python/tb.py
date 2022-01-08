import sys

from cosim_interface import CosimInterface
import wbfbd


WRITE_FIFO_PATH = sys.argv[1]
READ_FIFO_PATH = sys.argv[2]

cosim_interface = CosimInterface(WRITE_FIFO_PATH, READ_FIFO_PATH)

try:
    print("\nstarting cosimulation")

    Main = wbfbd.Main(cosim_interface)

    expected0 = 0b010101
    expected1 = 0b11

    print(f"Expecting Cfg0 value: {expected0}")

    print("Reading Cfg0")
    read_val = Main.Cfg0.read()
    if read_val != expected0:
        raise Exception(f"Read wrong value form Cfg0 {read_val}")

    print("Reading St0")
    read_val = Main.St0.read()
    if read_val != expected0:
        raise Exception(f"Read wrong value form St0 {read_val}")

    print(f"Expecting Cfg1 value: {expected1}")

    print("Reading Cfg1")
    read_val = Main.Cfg1.read()
    if read_val != expected1:
        raise Exception(f"Read wrong value form Cfg1 {read_val}")

    print("Reading St1")
    read_val = Main.St1.read()
    if read_val != expected1:
        raise Exception(f"Read wrong value form St1 {read_val}")

    print("\nending cosimulation")
    cosim_interface.end(0)

except Exception as E:
    cosim_interface.end(1)
    print(E)