import sys

from cosim_interface import CosimInterface
import vfbdb


WRITE_FIFO_PATH = sys.argv[1]
READ_FIFO_PATH = sys.argv[2]

cosim_interface = CosimInterface(WRITE_FIFO_PATH, READ_FIFO_PATH)

try:
    print("\nstarting cosimulation")

    Main = vfbdb.Main(cosim_interface)

    max_val = 2 ** vfbdb.Main.WIDTH - 1

    print("\nTesting Mask setting")

    Main.Mask.set()

    read = Main.Mask.read()
    assert read == max_val, f"read {read}, expecting {max_val}"
    read = Main.St.read()
    assert read == max_val, f"read {read}, expecting {max_val}"

    print("\nTesting Mask clear")

    Main.Mask.set([])

    read = Main.Mask.read()
    assert read == 0, f"read {read}, expecting 0"
    read = Main.St.read()
    assert read == 0, f"read {read}, expecting 0"

    print("\nTesting Mask setting single bit")

    Main.Mask.set(4)

    read = Main.Mask.read()
    assert read == 1 << 4, f"read {read}, expecting 4"
    read = Main.St.read()
    assert read == 1 << 4, f"read {read}, expecting 4"

    # Clear before next test.
    Main.Mask.set([])

    print("\nTesting Mask setting multiple bits")

    Main.Mask.set([0, 3])

    read = Main.Mask.read()
    assert read == 9, f"read {read}, expecting 9"
    read = Main.St.read()
    assert read == 9, f"read {read}, expecting 9"

    print("\nTesting Mask update")

    Main.Mask.set([])

    Main.Mask.update([0, 2])
    read = Main.Mask.read()
    assert read == 5, f"read {read}, expecting 5"
    read = Main.St.read()
    assert read == 5, f"read {read}, expecting 5"

    Main.Mask.update([2], mode="clear")
    read = Main.Mask.read()
    assert read == 1, f"read {read}, expecting 1"
    read = Main.St.read()
    assert read == 1, f"read {read}, expecting 1"

    print("\nending cosimulation")
    cosim_interface.end(0)

except Exception as E:
    cosim_interface.end(1)
    log.exception(E)
