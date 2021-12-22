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

    val = random.randint(0, 2**7 - 1)

    log.info(f"Generated random value: {val}")

    log.info("Writing cfg")
    main.cfg.write(val)

    log.info("Reading cfg")
    read_val = main.cfg.read()
    if read_val != val:
        raise Exception(f"Read wrong value form cfg {read_val}")

    log.info("Reading st")
    read_val = main.st.read()
    if read_val != val:
        raise Exception(f"Read wrong value form st {read_val}")

    cosim_interface.wait(5 * CLK_PERIOD)
    log.info("Ending cosimulation")
    cosim_interface.end(0)

except Exception as E:
    cosim_interface.end(1)
    log.exception(E)
