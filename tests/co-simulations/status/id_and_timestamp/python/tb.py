import sys
import traceback

from cosim_interface import CosimInterface
import wbfbd


WRITE_FIFO_PATH = sys.argv[1]
READ_FIFO_PATH = sys.argv[2]

CLK_PERIOD = 10

cosim_interface = CosimInterface(WRITE_FIFO_PATH, READ_FIFO_PATH)

try:
    print("\nstarting cosimulation\n")

    Main = wbfbd.Main(cosim_interface)

    id = Main.ID.read()

    assert id == wbfbd.ID, f"Read wrong ID {id}, expecting {wbfbd.ID}"

    print(f"ID: {id}\n")
    print(f"Timestamp: {Main.TIMESTAMP.read()}\n")

    print("\nending cosimulation")
    cosim_interface.end(0)

except Exception as E:
    cosim_interface.end(1)
    print(traceback.format_exc())
