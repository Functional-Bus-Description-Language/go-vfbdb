import sys
import traceback

from cosim_interface import CosimInterface
import wbfbd


WRITE_FIFO_PATH = sys.argv[1]
READ_FIFO_PATH = sys.argv[2]

cosim_interface = CosimInterface(WRITE_FIFO_PATH, READ_FIFO_PATH)

try:
    print("\nstarting cosimulation")

    Main = wbfbd.Main(cosim_interface)

    print("Reading ST")
    read_val = Main.St.read()
    # The expected value is magic. It has been manually checked,
    # and it may change in case of BFM changes. It depends on how much time
    # passes in the simulation from the start to the first read.
    if read_val != 0xFDFFFFFFFF:
        raise Exception(f"Read wrong value form St {read_val}")

    print("\nending cosimulation")
    cosim_interface.end(0)

except Exception as E:
    cosim_interface.end(1)
    print(traceback.format_exc())
