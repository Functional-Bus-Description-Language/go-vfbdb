import sys
import traceback

import cosim
import vfbdb


WRITE_FIFO_PATH = sys.argv[1]
READ_FIFO_PATH = sys.argv[2]

CLK_PERIOD = 10

iface = cosim.Iface(WRITE_FIFO_PATH, READ_FIFO_PATH)

try:
    Main = vfbdb.Main(iface)

    id = Main.ID.read()
    assert id == Main.ID.value, f"Read wrong ID {id}, expecting {Main.ID.value}"
    print(f"ID: {id}\n")

    ts = Main.TIMESTAMP.read()
    assert ts == Main.TIMESTAMP.value, f"Read wrong TIMESTAMP {ts}, expecting {Main.TIMESTAMP.value}"
    print(f"Timestamp: {ts}\n")

    iface.end(0)

except Exception as E:
    iface.end(1)
    print(traceback.format_exc())
