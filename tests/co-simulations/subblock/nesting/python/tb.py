import sys
import random

from cosim_interface import CosimInterface
import wbfbd

WRITE_FIFO_PATH = sys.argv[1]
READ_FIFO_PATH = sys.argv[2]

try:
    print("\nstarting cosimulation")

    cosim_interface = CosimInterface(WRITE_FIFO_PATH, READ_FIFO_PATH)

    Main = wbfbd.Main(cosim_interface)

    subblocks = [Main.Blk0, Main.Blk1, Main.Blk1.Blk2]

    for i, sb in enumerate(subblocks):
        print(f"Testing access to blk{i}")
        r = random.randrange(0, 2 ** 32 - 1)
        print(f"Writing value {r} to cfg register")
        sb.Cfg.write(r)

        print(f"Reading Cfg register")
        read = sb.Cfg.read()
        assert read == r, f"Read wrong value from Cfg register {read}"

        print(f"Reading St register")
        read = sb.St.read()
        assert read == r, f"Read wrong value from St register {read}"

    print("\nending cosimulation")
    cosim_interface.end(0)

except Exception as E:
    cosim_interface.end(1)
    log.exception(E)
