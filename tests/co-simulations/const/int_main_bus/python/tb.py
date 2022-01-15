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


    print("\n\nTesting int constant")
    print("Reading St register")
    read = Main.St.read()
    assert read == Main.C, f"read value {read} differs from constant value {Main.C}"


    print("\n\nTesting int list constants")
    print("Reading Stl register")
    read = Main.Stl.read()
    for i, v in enumerate(Main.CL):
        assert (
            read[i] == v
        ), f"read value {read[i]} differs from constant value {v}"


    print("\nending cosimulation")
    cosim_interface.end(0)

except Exception as E:
    cosim_interface.end(1)
    print(traceback.format_exc())
