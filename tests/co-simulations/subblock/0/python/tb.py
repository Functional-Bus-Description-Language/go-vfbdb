import random
import sys
import traceback

import cosim
import vfbdb

WRITE_FIFO_PATH = sys.argv[1]
READ_FIFO_PATH = sys.argv[2]

try:
    iface = cosim.Iface(WRITE_FIFO_PATH, READ_FIFO_PATH)

    Main = vfbdb.Main(iface)

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

    iface.end(0)

except Exception as E:
    iface.end(1)
    print(traceback.format_exc())
