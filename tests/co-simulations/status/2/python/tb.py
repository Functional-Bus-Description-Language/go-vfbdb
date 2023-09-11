import sys

import cosim
import vfbdb

WRITE_FIFO_PATH = sys.argv[1]
READ_FIFO_PATH = sys.argv[2]

try:
    iface = cosim.Iface(WRITE_FIFO_PATH, READ_FIFO_PATH)

    Main = vfbdb.Main(iface)

    values = Main.Status_array.read()
    assert len(values) == 9
    for i, v in enumerate(values):
        assert v == i

    idx = [2, 7]
    values = Main.Status_array.read(idx)
    assert values[0] == 2
    assert values[1] == 7

    value = Main.Status_array.read(5)
    assert value == 5

    iface.end(0)

except Exception as E:
    iface.end(1)
    print(traceback.format_exc())
