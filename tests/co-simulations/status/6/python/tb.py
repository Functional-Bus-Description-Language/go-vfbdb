import sys
import traceback

import cosim
import vfbdb


WRITE_FIFO_PATH = sys.argv[1]
READ_FIFO_PATH = sys.argv[2]

iface = cosim.Iface(WRITE_FIFO_PATH, READ_FIFO_PATH)

try:
    Main = vfbdb.Main(iface)

    print("Testing whole array read")
    data = Main.Sts.read()
    assert data[0] == Main.S0
    assert data[1] == Main.S1
    assert data[2] == Main.S2

    print("\nTesting index read")
    assert Main.Sts.read(0) == Main.S0
    assert Main.Sts.read(1) == Main.S1
    assert Main.Sts.read(2) == Main.S2

    data = Main.Sts.read([2, 0, 1])
    assert data[0] == Main.S2
    assert data[1] == Main.S0
    assert data[2] == Main.S1

    iface.end(0)

except Exception as E:
    iface.end(1)
    print(traceback.format_exc())
