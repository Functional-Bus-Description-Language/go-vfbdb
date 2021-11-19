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

CLK_PERIOD = 25


def delay_function():
    return CLK_PERIOD * random.randrange(5, 10)


try:
    log.info("Starting cosimulation")

    cosim_interface = CosimInterface(WRITE_FIFO_PATH, READ_FIFO_PATH, delay_function, True)

    main = wbfbd.main(cosim_interface)

    a = random.randint(0, 2**16-1)
    b = random.randint(0, 2**16-1)

    log.info(f"Calling add function, a = {a}, b = {b}")
    main.add(a, b)

    log.info(f"Reading result")
    result = main.result.read()

    if a + b != result:
        log.error(f"Wrong result, got {result}, expecting {a+b}")
        cosim_interface.end(1)

    cosim_interface.wait(5 * CLK_PERIOD)
    log.info("Ending cosimulation")
    cosim_interface.end(0)

except Exception as E:
    cosim_interface.end(1)
    log.exception(E)
