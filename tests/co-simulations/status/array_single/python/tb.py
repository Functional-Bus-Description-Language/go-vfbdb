import sys

import logging as log
log.basicConfig(
    level=log.DEBUG,
    format="%(module)s:%(levelname)s:%(message)s",
    datefmt="%H:%M:%S",
)
import random

from cosim_interface import CosimInterface
import wbfbd

WRITE_FIFO_PATH = sys.argv[1]
READ_FIFO_PATH  = sys.argv[2]

CLOCK_PERIOD = 10


def delay_function():
    return CLOCK_PERIOD * random.randrange(10, 40)


try:
    log.info("Starting cosimulation")

    cosim_interface = CosimInterface(WRITE_FIFO_PATH, READ_FIFO_PATH, delay_function, True)

    main = wbfbd.main(cosim_interface)

    values = main.status_array.read()
    assert len(values) == 9
    for i, v in enumerate(values):
        assert v == i

    idx = [2, 7]
    values = main.status_array.read(idx)
    assert values[0] == 2
    assert values[1] == 7

    value = main.status_array.read(5)
    assert value == 5

    cosim_interface.wait(5 * CLOCK_PERIOD)
    log.info("Ending cosimulation")
    cosim_interface.end(0)

except Exception as E:
    cosim_interface.end(1)
    log.exception(E)
