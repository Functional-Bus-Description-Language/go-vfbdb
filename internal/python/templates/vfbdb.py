# This file has been automatically generated by the vfbdb tool.
# Do not edit it manually, unless you really know what you do.
# https://github.com/Functional-Bus-Description-Language/go-vfbdb

import math
import time

BUS_WIDTH = {{.BusWidth}}

def calc_mask(m):
    """
    calc_mask calculates mask based on tuple (End Bit, Start Bit).
    The returned mask is shifted to the right.
    """
    return (((1 << (m[0] + 1)) - 1) ^ ((1 << m[1]) - 1)) >> m[1]

class _BufferIface:
    """
    _BufferIface is the internal interface used for reading/writing internal buffer
    (after reading)/(before writing) the target buffer. It is very useful
    as it allows treating proc or stream params/returns as configs/statuses.
    """
    def set_buf(self, buf):
        self.buf = buf

    def write(self, addr, data):
        self.buf[addr] = data

    def read(self, addr):
        return self.buf[addr]

def check_arg_values(params, *args):
    """
    check_arg_values checks that all arguments are in valid range and raises
    an exception if any argument is out of range.
    """
    for arg_idx, arg in enumerate(args):
        param = params[arg_idx]

        type = param['Access']['Type']

        if type.startswith("Single"):
            assert 0 <= arg < 2 ** param['Width'], \
                "{} value overrange ({})".format(param['Name'], arg)
        elif type.startswith("Array"):
            assert len(arg) == param['Access']['ItemCount'], \
                "invalid number of items ({}) for {} param, expecting {} items".format(len(arg), param['Name'], param['ItemCount'])

            for val_idx, v in enumerate(arg):
                assert 0 <= v < 2 ** param['Width'], "{}[{}] value overrange ({})".format(param['Name'], val_idx, v)
        else:
            raise Exception("invalid param access type {}".format(type))

def pack_params(params, *args):
    check_arg_values(params, *args)

    buf = []
    addr = None # Current argument address
    data = 0

    for arg_idx, arg in enumerate(args):
        param = params[arg_idx]
        a = param['Access']

        if addr is None:
            addr = a['StartAddr']
        elif a['StartAddr'] > addr:
            buf.append(data)
            data = 0
            addr = a['StartAddr']

        if a['Type'] == 'SingleOneReg':
            data |= arg << a['StartBit']
        elif a['Type'] == 'SingleNRegs':
            for r in range(a['RegCount']):
                if r == 0:
                    data |= (arg & calc_mask((BUS_WIDTH - 1, a['StartBit']))) << a['StartBit']
                    buf.append(data)
                    arg = arg >> (BUS_WIDTH - a['StartBit'])
                else:
                    addr += 1
                    data = arg & calc_mask((BUS_WIDTH, 0))
                    arg = arg >> BUS_WIDTH
                    if r < a['RegCount'] - 1:
                        buf.append(data)
                        data = 0
        elif a['Type'] == 'ArrayNRegs':
            start_bit = a['StartBit']
            for i, v in  enumerate(arg):
                width = param['Width']
                # Number of registers ith argument from vector occupies.
                reg_count = int(math.ceil((width - (BUS_WIDTH - start_bit)) / BUS_WIDTH)) + 1
                for _ in range(reg_count):
                    reg_width = width
                    if reg_width > BUS_WIDTH - start_bit:
                        reg_width = BUS_WIDTH - start_bit
                    data |= (v & ((1 << reg_width) - 1)) << start_bit
                    v >>= reg_width
                    start_bit = (start_bit + reg_width)
                    if start_bit >= BUS_WIDTH:
                        buf.append(data)
                        data = 0
                        start_bit %= BUS_WIDTH
                    width -= reg_width
        else:
            raise Exception("unhandled access type {}".format(a['Type']))

    buf.append(data)

    return buf

def create_mock_returns(buf_iface, start_addr, returns):
    """
    Create_mock_returns creates mock returns that can be used with internal software buffer.
    It is useful to be used with proc with returns and with upstram.
    """
    buf_size = 0
    rets = []
    for ret in returns:
        a = ret['Access']
        buf_size += a['RegCount']
        r = {}
        r['Name'] = ret['Name']
        # TODO: Add support for groups.

        if a['Type'] == 'SingleOneReg':
            r['Status'] = StatusSingleOneReg(
                buf_iface, a['StartAddr'] - start_addr, a['StartBit'], a['EndBit']
            )
        elif a['Type'] == 'SingleNRegs':
            r['Status'] = StatusSingleNRegs(
                buf_iface,
                a['StartAddr'] - start_addr,
                a['RegCount'],
                (BUS_WIDTH - 1, a['StartBit']),
                (a['EndBit'], 0),
            )
        else:
            raise Exception("unimplemented")

        rets.append(r)

    return buf_size, rets

class EmptyProc:
    def __init__(self, iface, call_addr, delay, exit_addr):
        self.iface = iface
        self.call_addr = call_addr
        self.delay = delay
        self.exit_addr = exit_addr
    def __call__(self):
        self.iface.write(self.call_addr, 0)
        if self.delay is not None:
            if self.delay != 0:
                time.sleep(self.delay)
            self.iface.read(self.exit_addr)

class ParamsProc:
    def __init__(self, iface, params_start_addr, params, delay, exit_addr):
        self.iface = iface
        self.params_start_addr = params_start_addr
        self.params = params
        self.delay = delay
        self.exit_addr = exit_addr

    def __call__(self, *args):
        assert len(args) == len(self.params), \
            "{}() takes {} arguments but {} were given".format(self.__name__, len(self.params), len(args))

        buf = pack_params(self.params, *args)

        if len(buf) == 1:
            self.iface.write(self.params_start_addr, buf[0])
        else:
            self.iface.writeb(self.params_start_addr, buf)

        if self.delay is not None:
            if self.delay != 0:
                time.sleep(self.delay)
            self.iface.read(self.exit_addr)

class ReturnsProc:
    def __init__(self, iface, returns_start_addr, returns, delay, call_addr):
        self.iface = iface
        self.returns_start_addr = returns_start_addr
        self.delay = delay
        self.call_addr = call_addr

        self.buf_iface = _BufferIface()
        self.buf_size, self.returns = create_mock_returns(self.buf_iface, returns_start_addr, returns)

    def __call__(self):
        if self.delay is not None:
            self.iface.write(self.call_addr, 0)
            if self.delay != 0:
                time.sleep(self.delay)

        if self.buf_size == 1:
            buf = [self.iface.read(self.returns_start_addr)]
        else:
            buf = self.iface.readb(self.returns_start_addr, self.buf_size)

        self.buf_iface.set_buf(buf)
        tup = [] # List to allow append but must be cast to tuple.

        for ret in self.returns:
            # NOTE: Groups are not yet supported so it is safe to immediately append.
            tup.append(ret['Status'].read())

        return tuple(tup)

class ParamsAndReturnsProc:
    def __init__(self, iface, params_start_addr, params, returns_start_addr, returns, delay):
        self.iface = iface

        self.params_start_addr = params_start_addr
        self.params = params

        self.returns_start_addr = returns_start_addr
        self.returns_buf_iface = _BufferIface()
        self.returns_buf_size, self.returns = create_mock_returns(self.returns_buf_iface, returns_start_addr, returns)

        self.delay = delay

    def __call__(self, *args):
        assert len(args) == len(self.params), \
            "{}() takes {} arguments but {} were given".format(self.__name__, len(self.params), len(args))

        params_buf = pack_params(self.params, *args)
        if len(params_buf) == 1:
            self.iface.write(self.params_start_addr, params_buf[0])
        else:
            self.iface.writeb(self.params_start_addr, params_buf)

        if self.delay is not None:
            if self.delay != 0:
                time.sleep(self.delay)

        if self.returns_buf_size == 1:
                returns_buf = [self.iface.read(self.returns_start_addr)]
        else:
            returns_buf = self.iface.readb(self.returns_start_addr, self.returns_buf_size)
        self.returns_buf_iface.set_buf(returns_buf)
        tup = [] # List to allow append but must be cast to tuple.
        for ret in self.returns:
            # NOTE: Groups are not yet supported so it is safe to immediately append.
            tup.append(ret['Status'].read())

        return tuple(tup)


class Static:
    def __init__(self, value):
        self._value = value
    @property
    def value(self):
        return self._value
    @value.setter
    def value(self, v):
        raise Exception(f"cannot set value of static element")


class StatusSingleOneReg:
    def __init__(self, iface, addr, start_bit, end_bit):
        self.iface = iface

        self.addr = addr
        self.start_bit = start_bit

        self.mask = calc_mask((end_bit, start_bit))
        self.width = end_bit - start_bit + 1

    def read(self):
        return (self.iface.read(self.addr) >> self.start_bit) & self.mask

class StaticSingleOneReg(Static, StatusSingleOneReg):
    def __init__(self, iface, addr, start_bit, end_bit, value):
        Static.__init__(self, value)
        StatusSingleOneReg.__init__(self, iface, addr, start_bit, end_bit)

class ConfigSingleOneReg(StatusSingleOneReg):
    def __init__(self, iface, addr, start_bit, end_bit):
        super().__init__(iface, addr, start_bit, end_bit)

    def write(self, data):
        assert 0 <= data < 2 ** self.width, "value overrange ({})".format(data)
        self.iface.write(self.addr, data << self.start_bit)

class MaskSingleOneReg(StatusSingleOneReg):
    def __init__(self, iface, addr, start_bit, end_bit):
        super().__init__(iface, addr, start_bit, end_bit)

    def _bits_to_iterable(self, bits):
        if bits == None:
            return range(self.width)
        elif type(bits) == int:
            return (bits,)
        return bits

    def _assert_bits_in_range(self, bits):
        for b in bits:
            assert 0 <= b < self.width, "mask overrange"

    def _assert_bits_to_update(self, bits):
        if bits == None:
            raise Exception("bits to update cannot have None value")
        if type(bits).__name__ in ["list", "tuple", "range", "set"] and len(bits) == 0:
            raise Exception("empty " + type(bits) + " of bits to update")

    def set(self, bits=None):
        bits = self._bits_to_iterable(bits)
        self._assert_bits_in_range(bits)

        mask = 0
        for b in bits:
            mask |= 1 << b

        self.iface.write(self.addr, mask << self.start_bit)

    def clear(self, bits=None):
        bits = self._bits_to_iterable(bits)
        self._assert_bits_in_range(bits)

        mask = self.mask
        for b in bits:
            mask ^= 1 << b

        self.iface.write(self.addr, mask << self.start_bit)

    def toggle(self, bits=None):
        bits = self._bits_to_iterable(bits)
        self._assert_bits_in_range(bits)

        xor_mask = 0
        for b in bits:
            xor_mask |= 1 << b
        xor_mask <<= self.start_bit

        mask = self.iface.read(self.addr) ^ xor_mask
        self.iface.write(self.addr, mask)

    def update_set(self, bits):
        self._assert_bits_to_update(bits)

        bits = self._bits_to_iterable(bits)
        self._assert_bits_in_range(bits)

        mask = 0
        for b in bits:
            mask |= 1 << b

        mask = self.iface.read(self.addr) | (mask << self.start_bit)
        self.iface.write(self.addr, mask)

    def update_clear(self, bits):
        self._assert_bits_to_update(bits)

        bits = self._bits_to_iterable(bits)
        self._assert_bits_in_range(bits)

        mask = 2**BUS_WIDTH - 1
        for b in bits:
            mask ^= 1 << b

        mask = self.iface.read(self.addr) & (mask << self.start_bit)
        self.iface.write(self.addr, mask)


class StatusSingleNRegs:
    def __init__(self, iface, start_addr, reg_count, start_mask, end_mask):
        self.iface = iface
        self.addrs = list(range(start_addr, start_addr + reg_count))
        self.width = 0
        self.masks = []
        self.reg_shifts = []
        self.data_shifts = []

        for i in range(reg_count):
            if i == 0:
                self.masks.append(calc_mask(start_mask))
                self.reg_shifts.append(start_mask[1])
                self.data_shifts.append(0)
                self.width += start_mask[0] - start_mask[1] + 1
            else:
                self.reg_shifts.append(0)
                self.data_shifts.append(self.width)
                if i == reg_count - 1:
                    self.masks.append(calc_mask(end_mask))
                    self.width += end_mask[0] - end_mask[1] + 1
                else:
                    self.masks.append(calc_mask((BUS_WIDTH - 1, 0)))
                    self.width += BUS_WIDTH

    def read(self):
        data = 0
        for i, a in enumerate(self.addrs):
            data |= ((self.iface.read(a) >> self.reg_shifts[i]) & self.masks[i]) << self.data_shifts[i]
        return data

class ConfigSingleNRegs(StatusSingleNRegs):
    def __init__(self, iface, start_addr, reg_count, start_mask, end_mask):
        super().__init__(iface, start_addr, reg_count, start_mask, end_mask)

    def write(self, data):
        assert 0 <= data < 2 ** self.width, "value overrange ({})".format(data)
        for i, a in enumerate(self.addrs):
            self.iface.write(a, ((data >> self.data_shifts[i]) & self.masks[i]) << self.reg_shifts[i])

class StaticSingleNRegs(Static, StatusSingleNRegs):
    def __init__(self, iface, start_addr, reg_count, start_mask, end_mask, value):
        Static.__init__(self, value)
        StatusSingleNRegs.__init__(self, iface, start_addr, reg_count, start_mask, end_mask)


class StatusArrayOneReg:
    def __init__(self, iface, addr, start_bit, width, item_count):
        self.iface = iface
        self.addr = addr
        self.start_bit = start_bit
        self.width = width
        self.item_count = item_count

    def __len__(self):
        return self.item_count

    def read(self, idx=None):
        reg = self.iface.read(self.addr)
        mask = (1 << self.width) - 1

        if type(idx) == int:
            assert 0 <= idx < self.item_count
            shift = self.start_bit + self.width * idx
            return (reg >> shift) & mask
        elif idx is None:
            idx = tuple(range(0, self.item_count))

        for i in idx:
            assert 0 <= i < self.item_count

        data = []
        for i in idx:
            shift = self.start_bit + self.width * i
            data.append((reg >> shift) & mask)

        return data

class ConfigArrayOneReg(StatusArrayOneReg):
    def __init__(self, iface, addr, start_bit, width, item_count):
        super().__init__(iface, addr, start_bit, width, item_count)

    def write(self, data, offset=0):
        """ offset - elements index offset, applied also when data is dictionary """
        assert 0 <= len(data) <= self.item_count, f"invalid data len {len(data)}"

        val = 0
        mask = 0

        if type(data) == dict:
            for i, v in data.items():
                assert type(i) == int, f'invalid index type {type(i)}'
                assert i >= 0, f"negative index {i}"
                assert i + offset < self.item_count, f"index overrange {i}"
                assert 0 <= v < 2 ** self.width, f"data out of range, index {i}, value {v}"
                shift = self.start_bit + (i + offset) * self.width
                val |= v << shift
                mask |= (2 ** self.width - 1) << shift
        else:
            assert len(data) + offset <= self.item_count

            for i, v in enumerate(data):
                assert 0 <= v < 2 ** self.width, f"data out of range, index {i}, value {v}"
                shift = (self.start_bit + (i + offset) * self.width)
                val |= v << shift
                mask |= 2 ** self.width - 1  << shift

        if len(data) == self.item_count:
            self.iface.write(self.addr, val)
        else:
            self.iface.rmw(self.addr, val, mask)


class StatusArrayOneInReg:
    def __init__(self, iface, addr, mask, item_count):
        self.iface = iface
        self.addr = addr
        self.mask = calc_mask(mask)
        self.shift = mask[1]
        self.width = mask[0] - mask[1] + 1
        self.item_count = item_count

    def __len__(self):
        return self.item_count

    def read(self, idx=None):
        if idx is None:
            idx = tuple(range(0, self.item_count))
            if self.item_count == 1:
                return (self.iface.read(self.addr) >> self.shift) & self.mask
            else:
                buf = self.iface.readb(self.addr, self.item_count)
                return [(data >> self.shift) & self.mask for data in buf]
        elif type(idx) == int:
            assert 0 <= idx < self.item_count
            return (self.iface.read(self.addr + idx) >> self.shift) & self.mask
        else:
            for i in idx:
                assert 0 <= i < self.item_count
            return [(self.iface.read(self.addr + i) >> self.shift) & self.mask for i in idx]

class ConfigArrayOneInReg(StatusArrayOneInReg):
    def __init__(self, iface, addr, mask, item_count):
        super().__init__(iface, addr, mask, item_count)

    def write(self, data, offset=0):
        """ offset - elements index offset, applied also when data is dictionary """
        assert 0 <= len(data) <= self.item_count, f"invalid data len {len(data)}"

        if type(data) == dict:
            idxs = sorted(data.keys())
            for idx in idxs:
                self.iface.write(self.addr + offset + idx, data[idx] << self.shift)
        else:
            assert len(data) + offset <= self.item_count

            if len(data) == 1:
                self.iface.write(self.addr + offset, data[0] << self.shift)
            else:
                buf = []
                for d in data:
                    buf.append(d << self.shift)
                self.iface.writeb(self.addr + offset, buf)


class StatusArrayNInReg:
    def __init__(self, iface, addr, start_bit, width, item_count, items_in_reg):
        self.iface = iface
        self.addr = addr
        self.start_bit = start_bit
        self.width = width
        self.item_count = item_count
        self.items_in_reg = items_in_reg
        self.reg_count = math.ceil(item_count / self.items_in_reg)

    def __len__(self):
        return self.item_count

    def read(self, idx=None):
        mask = (1 << self.width) - 1

        if idx is None:
            idx = tuple(range(0, self.item_count))
            reg_idx = tuple(range(self.reg_count))
        elif type(idx) == int:
            assert 0 <= idx < self.item_count
            reg_idx = idx // self.items_in_reg
            shift = self.start_bit + self.width * (idx % self.items_in_reg)
            return (self.iface.read(self.addr + reg_idx) >> shift) & mask
        else:
            reg_idx = set()
            for i in idx:
                assert 0 <= i < self.item_count
                reg_idx.add(i // self.items_in_reg)

        reg_data = {reg_i : self.iface.read(self.addr + reg_i) for reg_i in reg_idx}

        data = []
        for i in idx:
            shift = self.start_bit + self.width * (i % self.items_in_reg)
            data.append((reg_data[i // self.items_in_reg] >> shift) & mask)

        return data

class ConfigArrayNInReg(StatusArrayNInReg):
    def __init__(self, iface, addr, start_bit, width, item_count, items_in_reg):
        super().__init__(iface, addr, start_bit, width, item_count, items_in_reg)

    def write(self, data, offset=0):
        """ offset - elements index offset, applied also when data is dictionary """
        assert 0 <= len(data) <= self.item_count, f"invalid data len {len(data)}"

        regs = dict()
        def add_to_regs(idx, val):
            idx = idx + offset
            assert idx <= self.item_count, f"index overrange {idx + offset}"
            reg_idx = idx // self.items_in_reg
            if reg_idx not in regs:
                regs[reg_idx] = [0, 0] # [value, mask]
            shift = self.start_bit + (idx % self.items_in_reg) * self.width
            regs[reg_idx][0] |= val << shift
            regs[reg_idx][1] |= (2 ** self.width - 1)  << shift

        if type(data) == dict:
            for idx, val in data.items():
                add_to_regs(idx, val)
        else:
            for idx, val in enumerate(data):
                add_to_regs(idx, val)

        reg_idxs = sorted(regs.keys())
        for idx in reg_idxs:
            self.iface.rmw(self.addr + idx, regs[idx][0], regs[idx][1])


class StatusArrayNInRegMInEndReg(StatusArrayNInReg):
    def __init__(self, iface, addr, start_bit, width, item_count, items_in_reg):
        super().__init__(iface, addr, start_bit, width, item_count, items_in_reg)

class ConfigArrayNInRegMInEndReg(ConfigArrayNInReg):
    def __init__(self, iface, addr, start_bit, width, item_count, items_in_reg):
        super().__init__(iface, addr, start_bit, width, item_count, items_in_reg)


class StatusArrayOneInNRegs:
    def __init__(self, iface, addr, width, item_count, regs_per_item, reg_count, end_bit):
        self.iface = iface

        self.addr = addr
        self.width = width
        self.item_count = item_count

        self.regs_per_item = regs_per_item
        self.reg_count = reg_count
        self.last_reg_mask = calc_mask((end_bit, 0))

    def __len__(self):
        return self.item_count

    def _regs_to_data(self, buf):
        assert len(buf) == self.regs_per_item
        data = 0
        for i, bite in enumerate(buf):
            if i == len(buf) - 1:
                data |= (bite & self.last_reg_mask) << (i * BUS_WIDTH)
            else:
                data |= bite << (i * BUS_WIDTH)
        return data

    def read(self, idx=None):
        if idx is None:
            buf = self.iface.readb(self.addr, self.reg_count)
            data = []
            for i in range(self.item_count):
                data.append(self._regs_to_data(buf[i*self.regs_per_item:(i+1)*self.regs_per_item]))
            return data
        elif type(idx) == int:
            assert 0 <= idx < self.item_count
            buf = self.iface.readb(self.addr + idx * self.regs_per_item, self.regs_per_item)
            return self._regs_to_data(buf)
        else:
            data = []
            for i in idx:
                assert 0 <= i < self.item_count
                buf = self.iface.readb(self.addr + i * self.regs_per_item, self.regs_per_item)
                data.append(self._regs_to_data(buf))
            return data


class Upstream:
    def __init__(self, iface, addr, delay, returns):
        self.iface = iface
        self.addr = addr
        self.delay = delay
        self.buf_iface = _BufferIface()
        self.buf_size, self.returns = create_mock_returns(self.buf_iface, addr, returns)

    def read(self, n):
        """
        Read the stream n times.
        Read returns a tuple of tuples. Grouped returns are returned as dictionary (not yet supported).
        Non grouped returns are returned as values within tuple.
        """
        if self.buf_size == 1:
            read_data = [[x] for x in self.iface.cread(self.addr, n)]
        else:
            read_data = self.iface.creadb(self.addr, self.buf_size, n)

        data = []
        for buf in read_data:
            self.buf_iface.set_buf(buf)
            tup = [] # List to allow append but must be cast to tuple.

            for ret in self.returns:
                # NOTE: Groups are not yet supported so it is safe to immediately append.
                tup.append(ret['Status'].read())

            data.append(tuple(tup))

        return tuple(data)

class Downstream:
    def __init__(self, iface, addr, delay, params):
        self.iface = iface
        self.addr = addr
        self.params = params
        self.delay = delay

    def write(self, data):
        wbuf = [] # Write buffer
        args_in_one_reg = False # All arguments occupy one register

        for args in data:
            assert len(args) == len(self.params), f"invalid number of arguments {len(args)}, want {len(self.params)}"

            buf = pack_params(self.params, *args)
            if len(buf) == 1:
                args_in_one_reg = True
                wbuf.append(buf[0])
            else:
                wbuf.append(buf)

        if self.delay is None:
            if args_in_one_reg:
                self.iface.cwrite(self.addr, wbuf)
            else:
                self.iface.cwriteb(self.addr, wbuf)
        else:
            for i, val in enumerate(wbuf):
                if args_in_one_reg:
                    self.iface.write(self.addr, val)
                else:
                    self.iface.writeb(self.addr, buf)

                if i < len(wbuf) - 1:
                    time.sleep(self.delay)

{{.Code}}
