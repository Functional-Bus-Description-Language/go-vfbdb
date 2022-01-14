# This file has been automatically generated by the wbfbd tool.
# Do not edit it manually, unless you really know what you do.
# https://github.com/Functional-Bus-Description-Language/PyWbFBD

import math

BUS_WIDTH = {{.BusWidth}}

def calc_mask(m):
    return (((1 << (m[0] + 1)) - 1) ^ ((1 << m[1]) - 1)) >> m[1]

class SingleSingle:
    def __init__(self, iface, addr, mask):
        self.iface = iface
        self.addr = addr
        self.mask = calc_mask(mask)
        self.width = mask[0] - mask[1] + 1
        self.shift = mask[1]

    def read(self):
        return (self.iface.read(self.addr) >> self.shift) & self.mask

class ConfigSingleSingle(SingleSingle):
    def __init__(self, iface, addr, mask):
        super().__init__(iface, addr, mask)

    def write(self, val):
        assert 0 <= val < 2 ** self.width, "error: value overrange ({})".format(val)
        self.iface.write(self.addr, val << self.shift)

class ConfigSingleContinuous:
    def __init__(self, iface, start_addr, reg_count, start_mask, end_mask, decreasing_order):
        self.iface = iface
        self.addrs = list(range(start_addr, start_addr + reg_count))
        self.width = 0
        self.masks = []
        self.reg_shifts = []
        self.val_shifts = []

        for i in range(reg_count):
            if i == 0:
                self.masks.append(calc_mask(start_mask))
                self.reg_shifts.append(start_mask[1])
                self.val_shifts.append(0)
                self.width += start_mask[0] - start_mask[1] + 1
            else:
                self.reg_shifts.append(0)
                self.val_shifts.append(self.width)
                if i == reg_count - 1:
                    self.masks.append(calc_mask(end_mask))
                    self.width += end_mask[0] - end_mask[1] + 1
                else:
                    self.masks.append(calc_mask((BUS_WIDTH-1, 0)))
                    self.width += BUS_WIDTH

        if decreasing_order:
            self.addrs.reverse()
            self.masks.reverse()
            self.reg_shifts.reverse()
            self.val_shifts.reverse()

    def read(self):
        val = 0
        for i, a in enumerate(self.addrs):
            val |= ((self.iface.read(a) >> self.reg_shifts[i]) & self.masks[i]) << self.val_shifts[i]
        return val

    def write(self, val):
        assert 0 <= val < 2 ** self.width, "error: value overrange ({})".format(val)
        for i, a in enumerate(self.addrs):
            self.iface.write(a, ((val >> self.val_shifts[i]) & self.masks[i]) << self.reg_shifts[i])

class MaskSingleSingle(SingleSingle):
    def __init__(self, iface, addr, mask):
        super().__init__(iface, addr, mask)

    def set(self, bits=None):
        if bits == None:
            bits = range(self.width)
        elif type(bits) == int:
            bits = [bits]

        mask = 0
        for b in bits:
            assert 0 <= b < self.width, "mask overrange"
            mask |= 1 << b

        self.iface.write(self.addr, mask << self.shift)

    def update(self, bits, mode="set"):
        if mode not in ["set", "clear"]:
            raise Exception("invalid mode '" + mode + "'")
        if bits == None:
            raise Exception("bits to update cannot have None value")
        if type(bits).__name__ in ["list", "tuple", "range", "set"] and len(bits) == 0:
            raise Exception("empty " + type(bits) + " of bits to update")

        mask = 0
        reg_mask = 0
        for b in bits:
            assert 0 <= b < self.width, "mask overrange"
            if mode == "set":
                mask |= 1 << b
            reg_mask |= 1 << b

        self.iface.rmw(self.addr, mask << self.shift, reg_mask << self.shift)

class StatusSingleSingle(SingleSingle):
    def __init__(self, iface, addr, mask):
        super().__init__(iface, addr, mask)

class StatusArraySingle:
    def __init__(self, iface, addr, mask, item_count):
        self.iface = iface
        self.addr = addr
        self.mask = calc_mask(mask)
        self.shift = mask[1]
        self.item_count = item_count

    def read(self, idx=None):
        if idx is None:
            idx = tuple(range(0, self.item_count))
        elif type(idx) == int:
            assert 0 <= idx < self.item_count
            return (self.iface.read(self.addr + idx) >> self.shift) & self.mask
        else:
            for i in idx:
                assert 0 <= i < self.item_count

        return [(self.iface.read(self.addr + i) >> self.shift) & self.mask for i in idx]


class StatusArrayMultiple:
    def __init__(self, iface, addr, start_bit, width, item_count, items_per_access):
        self.iface = iface
        self.addr = addr
        self.start_bit = start_bit
        self.width = width
        self.item_count = item_count
        self.items_per_access = items_per_access
        self.reg_count = math.ceil(item_count / self.items_per_access)

    def read(self, idx=None):
        if idx is None:
            idx = tuple(range(0, self.item_count))
            reg_idx = tuple(range(self.reg_count))
        elif type(idx) == int:
            assert 0 <= idx < self.item_count
            reg_idx = idx // self.items_per_access
            shift = self.start_bit + self.width * (idx % self.items_per_access)
            mask = (1 << self.width) - 1
            return (self.iface.read(self.addr + reg_idx) >> shift) & mask
        else:
            reg_idx = set()
            for i in idx:
                assert 0 <= i < self.item_count
                reg_idx.add(i // self.items_per_access)

        reg_values = {reg_i : self.iface.read(self.addr + reg_i) for reg_i in reg_idx}

        values = []
        for i in idx:
            shift = self.start_bit + self.width * (i % self.items_per_access)
            mask = (1 << self.width) - 1
            values.append((reg_values[i // self.items_per_access] >> shift) & mask)

        return values

{{.Code}}
