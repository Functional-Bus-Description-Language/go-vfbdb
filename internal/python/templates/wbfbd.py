# This file has been automatically generated by the wbfbd tool.
# Do not edit it manually, unless you really know what you do.
# https://github.com/Functional-Bus-Description-Language/PyWbFBD

import math

BUS_WIDTH = {{.BusWidth}}

class StatusSingleSingle:
    def __init__(self, interface, addr, mask):
        self.interface = interface
        self.addr = addr
        self.mask = ((1 << (mask[0] + 1)) - 1) ^ ((1 << mask[1]) - 1)
        self.shift = mask[1]

    def read(self):
        return (self.interface.read(self.addr) & self.mask) >> self.shift


class StatusArraySingle:
    def __init__(self, interface, addr, mask, item_count):
        self.interface = interface
        self.addr = addr
        self.mask = ((1 << (mask[0] + 1)) - 1) ^ ((1 << mask[1]) - 1)
        self.shift = mask[1]
        self.item_count = item_count

    def read(self, idx=None):
        if idx is None:
            idx = tuple(range(0, self.item_count))
        elif type(idx) == int:
            assert 0 <= idx < self.item_count
            return (self.interface.read(self.addr + idx) & self.mask) >> self.shift
        else:
            for i in idx:
                assert 0 <= i < self.item_count

        return [(self.interface.read(self.addr + i) & self.mask) >> self.shift for i in idx]


class StatusArrayMultiple:
    def __init__(self, interface, addr, start_bit, width, item_count, items_per_access):
        self.interface = interface
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
            return (self.interface.read(self.addr + reg_idx) >> shift) & mask
        else:
            reg_idx = set()
            for i in idx:
                assert 0 <= i < self.item_count
                reg_idx.add(i // self.items_per_access)

        reg_values = {reg_i : self.interface.read(self.addr + reg_i) for reg_i in reg_idx}

        values = []
        for i in idx:
            shift = self.start_bit + self.width * (i % self.items_per_access)
            mask = (1 << self.width) - 1
            values.append((reg_values[i // self.items_per_access] >> shift) & mask)

        return values

class ConfigSingleSingle:
    def __init__(self, interface, addr, mask):
        self.interface = interface
        self.addr = addr
        self.mask = ((1 << (mask[0] + 1)) - 1) ^ ((1 << mask[1]) - 1)
        self.shift = mask[1]

    def read(self):
        return (self.interface.read(self.addr) & self.mask) >> self.shift

    def write(self, val):
        self.interface.write(self.addr, val << self.shift)

{{.Code}}
