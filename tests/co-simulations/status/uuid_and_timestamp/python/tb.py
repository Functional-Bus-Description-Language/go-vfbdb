import sys

from cosim_interface import CosimInterface
import wbfbd


WRITE_FIFO_PATH = sys.argv[1]
READ_FIFO_PATH = sys.argv[2]

CLK_PERIOD = 10

cosim_interface = CosimInterface(WRITE_FIFO_PATH, READ_FIFO_PATH)

try:
    print("\nstarting cosimulation\n")

    main = wbfbd.main(cosim_interface)

    print(f"UUID: {main.x_uuid_x.read()}\n")
    print(f"Timestamp: {main.x_timestamp_x.read()}\n")

    print("\nending cosimulation")
    cosim_interface.end(0)

except Exception as E:
    cosim_interface.end(1)
    log.exception(E)
