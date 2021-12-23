import sys

from cosim_interface import CosimInterface
import wbfbd

WRITE_FIFO_PATH = sys.argv[1]
READ_FIFO_PATH = sys.argv[2]


try:
    print("\nstarting cosimulation")

    cosim_interface = CosimInterface(WRITE_FIFO_PATH, READ_FIFO_PATH)

    main = wbfbd.main(cosim_interface)

    for i in range(10):
        print(f"calling foo function")
        main.foo()

        print(f"Reading count")
        count = main.count.read()

        if count != i + 1:
            log.error(f"Wrong count, got {count}, expecting {i+1}")
            cosim_interface.end(1)

    print("\nending cosimulation")
    cosim_interface.end(0)

except Exception as E:
    cosim_interface.end(1)
    log.exception(E)
