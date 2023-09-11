import sys
import traceback

import cosim
import vfbdb

WRITE_FIFO_PATH = sys.argv[1]
READ_FIFO_PATH = sys.argv[2]


try:
    iface = cosim.Iface(WRITE_FIFO_PATH, READ_FIFO_PATH)

    Main = vfbdb.Main(iface)

    print("Testing count % items per access = 0 scenerio.")

    values = Main.Status_array0.read()
    assert len(values) == 8
    for i, v in enumerate(values):
        assert v == i, f"got {v}, expecting {i}"

    idx = [2, 7]
    values = Main.Status_array0.read(idx)
    assert values[0] == 2
    assert values[1] == 7

    value = Main.Status_array0.read(5)
    assert value == 5

    print("Testing count < items per access scenario.")

    values = Main.Status_array1.read()
    assert len(values) == 4
    for i, v in enumerate(values):
        assert v == i

    idx = [0, 3]
    values = Main.Status_array1.read(idx)
    assert values[0] == 0
    assert values[1] == 3

    value = Main.Status_array1.read(2)
    assert value == 2

    print(
        "Testing scenerio when the number of items in the last register is different."
    )

    values = Main.Status_array2.read()
    assert len(values) == 6
    for i, v in enumerate(values):
        assert v == i

    idx = [1, 5]
    values = Main.Status_array2.read(idx)
    assert values[0] == 1
    assert values[1] == 5

    value = Main.Status_array2.read(0)
    assert value == 0

    iface.end(0)

except Exception as E:
    iface.end(1)
    print(traceback.format_exc())
