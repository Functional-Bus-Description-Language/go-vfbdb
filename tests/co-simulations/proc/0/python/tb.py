import sys
import traceback

import cosim
import vfbdb

WRITE_FIFO_PATH = sys.argv[1]
READ_FIFO_PATH = sys.argv[2]


try:
    iface = cosim.Iface(WRITE_FIFO_PATH, READ_FIFO_PATH)

    Main = vfbdb.Main(iface)

    for i in range(10):
        print(f"calling foo function")
        Main.Foo()

        print(f"Reading count")
        count = Main.Count.read()

        if count != i + 1:
            log.error(f"Wrong count, got {count}, expecting {i+1}")
            iface.end(1)

    iface.end(0)

except Exception as E:
    iface.end(1)
    print(traceback.format_exc())
