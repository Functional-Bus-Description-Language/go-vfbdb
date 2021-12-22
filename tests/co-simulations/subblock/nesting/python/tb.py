import sys
import random

from cosim_interface import CosimInterface
import wbfbd

WRITE_FIFO_PATH = sys.argv[1]
READ_FIFO_PATH = sys.argv[2]

try:
    print("\nstarting cosimulation")

    cosim_interface = CosimInterface(WRITE_FIFO_PATH, READ_FIFO_PATH)

    main = wbfbd.main(cosim_interface)

    subblocks = [main.blk0, main.blk1, main.blk1.blk2]

    for i, sb in enumerate(subblocks):
        print(f"Testing access to blk{i}")
        r = random.randrange(0, 2 ** 32 - 1)
        print(f"Writing value {r} to cfg register")
        sb.cfg.write(r)

        print(f"Reading cfg register")
        read = sb.cfg.read()
        assert read == r, f"Read wrong value from cfg register {read}"

        print(f"Reading st register")
        read = sb.st.read()
        assert read == r, f"Read wrong value from st register {read}"

    print("\nending cosimulation")
    cosim_interface.end(0)

except Exception as E:
    cosim_interface.end(1)
    log.exception(E)
