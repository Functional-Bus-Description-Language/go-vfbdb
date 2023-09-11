import sys
import traceback

import cosim
import vfbdb


WRITE_FIFO_PATH = sys.argv[1]
READ_FIFO_PATH = sys.argv[2]

iface = cosim.Iface(WRITE_FIFO_PATH, READ_FIFO_PATH)

try:
    Main = vfbdb.Main(iface)

    print("Reading ST")
    read_val = Main.St.read()
    # The expected value is magic. It has been manually checked,
    # and it may change in case of BFM changes. It depends on how much time
    # passes in the simulation from the start to the first read.
    if read_val != 0xFDFFFFFFFF:
        raise Exception(f"Read wrong value form St {read_val}")

    iface.end(0)

except Exception as E:
    iface.end(1)
    print(traceback.format_exc())
