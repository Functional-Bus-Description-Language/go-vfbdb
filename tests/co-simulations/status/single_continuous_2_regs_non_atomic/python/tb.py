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

    lower = random.randint(0, 2 ** 30  - 1)
    upper = random.randint(0, 2 ** 20  - 1)

    print(f"Generated random values: lower = {lower}, upper = {upper}")

    print("Writing Lower")
    Main.Lower.write(lower)

    print("Writing Upper")
    Main.Upper.write(upper)

    print("Reading St")
    st = Main.St.read()
    if st != (upper << 30) | lower:
        raise Exception(f"Read wrong value form St {st}, expects {(upper << 30) | lower}")

    print("\nending cosimulation")
    cosim_interface.end(0)

except Exception as E:
    cosim_interface.end(1)
    print(E)
