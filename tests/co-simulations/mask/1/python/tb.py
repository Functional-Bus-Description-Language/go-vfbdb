import sys
import traceback

import cosim
import vfbdb


WRITE_FIFO_PATH = sys.argv[1]
READ_FIFO_PATH = sys.argv[2]

iface = cosim.Iface(WRITE_FIFO_PATH, READ_FIFO_PATH)

try:
    Main = vfbdb.Main(iface)

    max_val = 2 ** vfbdb.Main.WIDTH - 1

    print("\nTesting Mask setting")
    Main.Mask.set()
    read = Main.Mask.read()
    assert read == max_val, f"read {read}"
    read = Main.St.read()
    assert read == max_val, f"read {read}"

    print("\nTesting Mask clear")
    Main.Mask.set([])
    read = Main.Mask.read()
    assert read == 0, f"read {read}"
    read = Main.St.read()
    assert read == 0, f"read {read}"

    print("\nTesting Mask setting single bit")
    Main.Mask.set(32)
    read = Main.Mask.read()
    assert read == 1 << 32, f"read {read}"
    read = Main.St.read()
    assert read == 1 << 32, f"read {read}"

    # Clear before next test.
    Main.Mask.set([])
    print("\nTesting Mask setting multiple bits")
    bits = [0, 3, 40]
    Main.Mask.set(bits)
    want = sum([1 << b for b in bits])
    read = Main.Mask.read()
    assert read == want, f"read {read}, want {want}"
    read = Main.St.read()
    assert read == want, f"read {read}, want {want}"

    print("\nTesting Mask update_set")
    Main.Mask.set([])
    bits = [0, 2, 39]
    Main.Mask.update_set(bits)
    want = sum([1 << b for b in bits])
    read = Main.Mask.read()
    assert read == want, f"read {read}, want {want}"
    read = Main.St.read()
    assert read == want, f"read {read}, want {want}"

    print("\nTesting Mask update_clear")
    Main.Mask.update_clear([2])
    read = Main.Mask.read()
    assert read == 1 | (1 << 39), f"read {read}"
    read = Main.St.read()
    assert read == 1 | (1 << 39), f"read {read}"

    # Clear before next test.
    Main.Mask.set([])
    print("\nTesting Mask toggle")
    bits = [0, 2, 43]
    Main.Mask.toggle(bits)
    want = sum([1 << b for b in bits])
    read = Main.Mask.read()
    assert read == want, f"read {read}, want {want}"
    read = Main.St.read()
    assert read == want, f"read {read}, want {want}"

    iface.end(0)

except Exception as E:
    iface.end(1)
    print(traceback.format_exc())
