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
    return CLOCK_PERIOD * random.randrange(5, 10)


try:
    log.info("Starting cosimulation")

    cosim_interface = CosimInterface(WRITE_FIFO_PATH, READ_FIFO_PATH, delay_function, True)

    main = wbfbd.main(cosim_interface)

    subblocks = [main.blk0, main.blk1, main.blk1.blk2]

    for i, sb in enumerate(subblocks):
        log.info(f"Testing access to blk{i}")
        r = random.randrange(0, 2**32-1)
        log.info(f"Writing value {r} to cfg register")
        sb.cfg.write(r)

        log.info(f"Reading cfg register")
        read = sb.cfg.read()
        assert read == r, f"Read wrong value from cfg register {read}"

        log.info(f"Reading st register")
        read = sb.st.read()
        assert read == r, f"Read wrong value from st register {read}"

    cosim_interface.wait(5 * CLOCK_PERIOD)
    log.info("Ending cosimulation")
    cosim_interface.end(0)

except Exception as E:
    cosim_interface.end(1)
    log.exception(E)
