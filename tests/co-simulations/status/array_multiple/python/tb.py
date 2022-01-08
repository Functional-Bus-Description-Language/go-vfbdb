import sys

from cosim_interface import CosimInterface
import wbfbd

WRITE_FIFO_PATH = sys.argv[1]
READ_FIFO_PATH = sys.argv[2]


try:
    print("\nstarting cosimulation")

    cosim_interface = CosimInterface(WRITE_FIFO_PATH, READ_FIFO_PATH)

    Main = wbfbd.Main(cosim_interface)

    print("Testing count % iterm per access = 0 scenerio.")

    values = Main.Status_array0.read()
    assert len(values) == 8
    for i, v in enumerate(values):
        assert v == i

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

    print("\nending cosimulation")
    cosim_interface.end(0)

except Exception as E:
    cosim_interface.end(1)
    log.exception(E)