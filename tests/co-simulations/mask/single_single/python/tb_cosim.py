import sys

import logging as log
log.basicConfig(
    level=log.DEBUG,
    format="%(module)s: %(levelname)s: %(message)s",
    datefmt="%H:%M:%S",
)
import random

from cosim_interface import CosimInterface
import wbfbd


WRITE_FIFO_PATH = sys.argv[1]
READ_FIFO_PATH  = sys.argv[2]

CLK_PERIOD = 10

def delay_function():
    return CLK_PERIOD * random.randrange(5, 10)


cosim_interface = CosimInterface(WRITE_FIFO_PATH, READ_FIFO_PATH, delay_function, True)

try:
    log.info("Starting cosimulation")

    main = wbfbd.main(cosim_interface)

    max_val = 2**wbfbd.main.WIDTH - 1

    log.info("Testing mask setting")

    main.mask.set()

    read = main.mask.read()
    assert read == max_val, f"read {read}, expecting {max_val}"
    read = main.st.read()
    assert read == max_val, f"read {read}, expecting {max_val}"

    log.info("Testing mask clear")

    main.mask.set([])

    read = main.mask.read()
    assert read == 0, f"read {read}, expecting 0"
    read = main.st.read()
    assert read == 0, f"read {read}, expecting 0"

    log.info("Testing mask setting single bit")

    main.mask.set(4)

    read = main.mask.read()
    assert read == 1<<4, f"read {read}, expecting 4"
    read = main.st.read()
    assert read == 1<<4, f"read {read}, expecting 4"

    # Clear before next test.
    main.mask.set([])

    log.info("Testing mask setting multiple bits")

    main.mask.set([0,3])

    read = main.mask.read()
    assert read == 9, f"read {read}, expecting 9"
    read = main.st.read()
    assert read == 9, f"read {read}, expecting 9"

    log.info("Testing mask update")

    main.mask.set([])

    main.mask.update([0, 2])
    read = main.mask.read()
    assert read == 5, f"read {read}, expecting 5"
    read = main.st.read()
    assert read == 5, f"read {read}, expecting 5"

    main.mask.update([2], mode="clear")
    read = main.mask.read()
    assert read == 1, f"read {read}, expecting 1"
    read = main.st.read()
    assert read == 1, f"read {read}, expecting 1"

    cosim_interface.wait(5 * CLK_PERIOD)
    log.info("Ending cosimulation")
    cosim_interface.end(0)

except Exception as E:
    cosim_interface.end(1)
    log.exception(E)
