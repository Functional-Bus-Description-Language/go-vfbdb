import sys

from cosim_interface import CosimInterface
import wbfbd


WRITE_FIFO_PATH = sys.argv[1]
READ_FIFO_PATH = sys.argv[2]

cosim_interface = CosimInterface(WRITE_FIFO_PATH, READ_FIFO_PATH)

try:
    print("\nstarting cosimulation")

    main = wbfbd.main(cosim_interface)

    expected0 = 0b010101
    expected1 = 0b11

    print(f"Expecting cfg0 value: {expected0}")

    print("Reading cfg0")
    read_val = main.cfg0.read()
    if read_val != expected0:
        raise Exception(f"Read wrong value form cfg0 {read_val}")

    print("Reading st0")
    read_val = main.st0.read()
    if read_val != expected0:
        raise Exception(f"Read wrong value form st0 {read_val}")

    print(f"Expecting cfg1 value: {expected1}")

    print("Reading cfg1")
    read_val = main.cfg1.read()
    if read_val != expected1:
        raise Exception(f"Read wrong value form cfg1 {read_val}")

    print("Reading st1")
    read_val = main.st1.read()
    if read_val != expected1:
        raise Exception(f"Read wrong value form st1 {read_val}")

    print("\nending cosimulation")
    cosim_interface.end(0)

except Exception as E:
    cosim_interface.end(1)
    print(E)
