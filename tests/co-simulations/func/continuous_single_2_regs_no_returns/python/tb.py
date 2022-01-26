import random
import sys
import traceback

from cosim_interface import CosimInterface
import wbfbd

WRITE_FIFO_PATH = sys.argv[1]
READ_FIFO_PATH = sys.argv[2]


try:
    print("\nstarting cosimulation")

    cosim_interface = CosimInterface(WRITE_FIFO_PATH, READ_FIFO_PATH)

    Main = wbfbd.Main(cosim_interface)

    c = random.randint(2**33, 2 ** 40 - 1)
    s = random.randint(0, 2 ** 16 - 1)

    print(f"Calling add function, c = {c}, s = {s}")
    Main.Add(c, s)

    print(f"Reading result")
    result = Main.Result.read()

    if c + s != result:
        print(f"Wrong result, got {result}, expecting {c+s}")
        cosim_interface.end(1)

    print("\nending cosimulation")
    cosim_interface.end(0)

except Exception as E:
    cosim_interface.end(1)
    print(traceback.format_exc())
