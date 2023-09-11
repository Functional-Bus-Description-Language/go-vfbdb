import random
import traceback
import sys

import cosim
import vfbdb

WRITE_FIFO_PATH = sys.argv[1]
READ_FIFO_PATH = sys.argv[2]


try:
    iface = cosim.Iface(WRITE_FIFO_PATH, READ_FIFO_PATH)

    Main = vfbdb.Main(iface)

    a = random.randint(0, 2 ** 16 - 1)
    b = random.randint(0, 2 ** 16 - 1)

    print(f"Calling add function, a = {a}, b = {b}")
    Main.Add(a, b)

    print(f"Reading result")
    result = Main.Result.read()

    if a + b != result:
        print(f"Wrong result, got {result}, expecting {a+b}")
        iface.end(1)

    iface.end(0)

except Exception as E:
    iface.end(1)
    print(traceback.format_exc())
