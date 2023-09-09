import random
import sys
import traceback

import cosim
import vfbdb


WRITE_FIFO_PATH = sys.argv[1]
READ_FIFO_PATH = sys.argv[2]

iface = cosim.Iface(WRITE_FIFO_PATH, READ_FIFO_PATH)

try:
    print("\nstarting cosimulation")

    Main = vfbdb.Main(iface)

    print("list test")
    data = []
    for _ in range(len(Main.Cfgs)):
        data.append(random.randint(0, 2**Main.Cfgs.width - 1))

    Main.Cfgs.write(data)
    rdata = Main.Cfgs.read()
    assert rdata == data, f"invalid data read, got {rdata}, want {data}"

    # Clear data
    data = [0 for _ in range(10)]
    Main.Cfgs.write(data)

    print("dictionary test")
    data = {0: 123, 3: 9876, 7: 111, 9: 23456}
    Main.Cfgs.write(data)
    rdata = Main.Cfgs.read()
    assert rdata[0] == data[0], f"got {rdata[0]}, want {data[0]}"
    assert rdata[1] == 0, f"got {rdata[1]}, want 0"
    assert rdata[2] == 0, f"got {rdata[2]}, want 0"
    assert rdata[3] == data[3], f"got {rdata[3]}, want {data[3]}"
    assert rdata[4] == 0, f"got {rdata[4]}, want 0"
    assert rdata[5] == 0, f"got {rdata[5]}, want 0"
    assert rdata[6] == 0, f"got {rdata[6]}, want 0"
    assert rdata[7] == data[7], f"got {rdata[7]}, want {data[7]}"
    assert rdata[8] == 0, f"got {rdata[8]}, want 0"
    assert rdata[9] == data[9], f"got {rdata[9]}, want {data[9]}"

    # Clear data
    data = [0 for _ in range(10)]
    Main.Cfgs.write(data)

    print("offset test")
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
            assert rdata[i] == data[i - offset], f"got {rdata[i]}, want data[i - offset]"

    print("\nending cosimulation")
    iface.end(0)

except Exception as E:
    print(traceback.format_exc())
    iface.end(1)
