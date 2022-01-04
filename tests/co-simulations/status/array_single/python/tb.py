import sys

from cosim_interface import CosimInterface
import wbfbd

WRITE_FIFO_PATH = sys.argv[1]
READ_FIFO_PATH = sys.argv[2]

try:
    print("\nstarting cosimulation")

    cosim_interface = CosimInterface(WRITE_FIFO_PATH, READ_FIFO_PATH)

    Main = wbfbd.Main(cosim_interface)

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

    print("\nending cosimulation")
    cosim_interface.end(0)

except Exception as E:
    cosim_interface.end(1)
    log.exception(E)
