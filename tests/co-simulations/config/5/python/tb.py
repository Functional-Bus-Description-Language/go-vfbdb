import random
import sys
import traceback

import cosim
import vfbdb


WRITE_FIFO_PATH = sys.argv[1]
READ_FIFO_PATH = sys.argv[2]

iface = cosim.Iface(WRITE_FIFO_PATH, READ_FIFO_PATH)

try:
    Main = vfbdb.Main(iface)

    print("\n\nlist test")
    data = []
    for _ in range(len(Main.Cfgs)):
        data.append(random.randint(0, 2**Main.Cfgs.width - 1))

    Main.Cfgs.write(data)
    rdata = Main.Cfgs.read()
    assert rdata == data, f"invalid data read, got {rdata}, want {data}"

    # Clear data
    data = [0 for _ in range(10)]
    Main.Cfgs.write(data)

    print("\n\ndictionary test")
    data = {0: 123, 3: 9876, 7: 111, 9: 23456}
    Main.Cfgs.write(data)
    rdata = Main.Cfgs.read()
    for i in range(len(Main.Cfgs)):
        if i in data:
            assert rdata[i] == data[i], f"{i}: got {rdata[0]}, want {data[0]}"
        else:
            assert rdata[i] == 0, f"{i}: got {rdata[0]}, want 0"

    # Clear data
    data = [0 for _ in range(10)]
    Main.Cfgs.write(data)

    print("\n\noffset test")
    offset = 3
    data = []
    for _ in range(len(Main.Cfgs) - offset):
        data.append(random.randint(0, 2**Main.Cfgs.width - 1))

    Main.Cfgs.write(data, offset)
    rdata = Main.Cfgs.read()
    for i in range(len(Main.Cfgs)):
        if i < offset:
            assert rdata[i] == 0, f"got {rdata[i]}, want 0"
        else:
            assert rdata[i] == data[i - offset], f"got {rdata[i]}, want {data[i - offset]}"

    iface.end(0)

except Exception as E:
    print(traceback.format_exc())
    iface.end(1)
