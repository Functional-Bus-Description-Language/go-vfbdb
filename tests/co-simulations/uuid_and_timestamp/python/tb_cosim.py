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


try:
    log.info("Starting cosimulation")

    cosim_interface = CosimInterface(WRITE_FIFO_PATH, READ_FIFO_PATH, delay_function, True)

    main = wbfbd.main(cosim_interface)

    log.info(f"UUID: {main.x_uuid_x.read()}")
    log.info(f"Timestamp: {main.x_timestamp_x.read()}")

    cosim_interface.wait(5 * CLK_PERIOD)
    log.info("Ending cosimulation")
    cosim_interface.end(0)

except Exception as E:
    cosim_interface.end(1)
    log.exception(E)
