import sys

from cosim_interface import CosimInterface
import wbfbd


WRITE_FIFO_PATH = sys.argv[1]
READ_FIFO_PATH  = sys.argv[2]

cosim_interface = CosimInterface(WRITE_FIFO_PATH, READ_FIFO_PATH)

try:
    print("\nStarting cosimulation")

    main = wbfbd.main(cosim_interface)

    print("Reading st register")
    read = main.st.read()
    assert read == main.C, f"read value {read} differs from constant value {main.C}"

    print("ending cosimulation")
    cosim_interface.end(0)

except Exception as E:
    cosim_interface.end(1)
    log.exception(E)
