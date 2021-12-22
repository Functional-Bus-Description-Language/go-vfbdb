import sys
import random

from cosim_interface import CosimInterface
import wbfbd

WRITE_FIFO_PATH = sys.argv[1]
READ_FIFO_PATH = sys.argv[2]


try:
    print("\nstarting cosimulation")

    cosim_interface = CosimInterface(WRITE_FIFO_PATH, READ_FIFO_PATH)

    main = wbfbd.main(cosim_interface)

    a = random.randint(0, 2 ** 16 - 1)
    b = random.randint(0, 2 ** 16 - 1)

    print(f"Calling add function, a = {a}, b = {b}")
    main.add(a, b)

    print(f"Reading result")
    result = main.result.read()

    if a + b != result:
        log.error(f"Wrong result, got {result}, expecting {a+b}")
        cosim_interface.end(1)

    print("\nending cosimulation")
    cosim_interface.end(0)

except Exception as E:
    cosim_interface.end(1)
    log.exception(E)
