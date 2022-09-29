import sys
import traceback

import cosim
import vfbdb


WRITE_FIFO_PATH = sys.argv[1]
READ_FIFO_PATH = sys.argv[2]

CLK_PERIOD = 10

iface = cosim.Iface(WRITE_FIFO_PATH, READ_FIFO_PATH)

try:
    print("\nstarting cosimulation\n")

    Main = vfbdb.Main(iface)

    id = Main.ID.read()

    assert id == vfbdb.ID, f"Read wrong ID {id}, expecting {vfbdb.ID}"

    print(f"ID: {id}\n")
    print(f"Timestamp: {Main.TIMESTAMP.read()}\n")

    print("\nending cosimulation")
    iface.end(0)

except Exception as E:
    iface.end(1)
    print(traceback.format_exc())
