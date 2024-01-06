import os


class Iface:
    def __init__(
        self, write_fifo_path, read_fifo_path, delay_function=None, delay=False
    ):
        """Create co-simulation interface.
        write_fifo_path - path to software -> firmware named pipe
        read_fifo_path  - path to firmware -> software named pipe
        delay_function  - reference to function returning random value when delay is set to 'True'
        delay - if set to 'True' there is a random delay between any write or read operation.
                Useful for modelling real access times.
        """
        self.write_fifo_path = write_fifo_path
        self.read_fifo_path = read_fifo_path

        self._make_fifos()
        self.write_fifo = open(write_fifo_path, "w")
        self.read_fifo = open(read_fifo_path, "r")

        if delay and delay_function is None:
            raise Exception("delay set to 'True', but delay_function not provided")

        self.delay_function = delay_function
        self.delay = delay

        # Attributes related with statistics collection.
        self.write_count = 0
        self.read_count = 0
        self.cwrite_count = 0
        self.cread_count = 0
        self.writeb_count = 0
        self.cwriteb_count = 0
        self.creadb_count = 0
        self.readb_count = 0
        self.rmw_count = 0

    def _make_fifos(self):
        """Create named pipes needed for inter-process communication."""
        self._remove_fifos()
        os.mkfifo(self.write_fifo_path)
        os.mkfifo(self.read_fifo_path)

    def _remove_fifos(self):
        """Remove named pipes."""
        try:
            os.remove(self.write_fifo_path)
            os.remove(self.read_fifo_path)
        except:
            pass

    def write(self, addr, data):
        """Single Write
        addr  - register address
        data  - data to be written
        """
        if self.delay:
            self.wait(self.delay_function())

        print(
            "write: addr 0x{:08x}, data {} (0x{:08x}) (0b{:032b})".format(addr, data, data, data)
        )

        cmd = "W" + ("%.8x" % addr) + "," + ("%.8x" % data) + "\n"
        self.write_fifo.write(cmd)
        self.write_fifo.flush()

        s = self.read_fifo.readline()
        if s.strip() == "ACK":
            self.write_count += 1
            return
        else:
            raise Exception("Wrong status returned:" + s.strip())

    def read(self, addr):
        """Single Read
        addr - register address
        """
        if self.delay:
            self.wait(self.delay_function())

        print("read: addr 0x{:08x}".format(addr))

        cmd = "R" + ("%.8x" % addr) + "\n"
        self.write_fifo.write(cmd)
        self.write_fifo.flush()

        s = self.read_fifo.readline()
        if s.strip() == "ERR":
            raise Exception("Error status returned")

        self.read_count += 1
        data = int(s, 2)
        print("read: data {} (0x{:08x}) (0b{:032b})".format(data, data, data))

        return data

    def cwrite(self, addr, data):
        """Cyclic Read
        addr - register address
        data - data
        """
        for d in data:
            self.write(addr, d)
        self.cwrite_count += 1

    def cread(self, addr, n):
        """Cyclic Read
        addr - register address
        n    - number of reads
        """
        data = [self.read(addr) for _ in range(n)]
        self.cread_count += 1
        return data

    def writeb(self, addr, data):
        """Block Write
        addr - start address
        data - buffer with data to be written
        """
        if self.delay:
            self.wait(self.delay_function())

        print(
            "writeb: addr 0x{:08x}, count {}".format(addr, len(data))
        )

        for i, d in enumerate(data):
            self.write(addr + i, d)

        self.writeb_count += 1

    def readb(self, addr, block_size):
        """Block Read
        addr - start address
        block_size - block size to be read (in words, not bytes)
        """
        if self.delay:
            self.wait(self.delay_function())

        print("readb: addr 0x{:08x}, block size {}".format(addr, block_size))

        buf = []
        for i in range(block_size):
            buf.append(self.read(addr + i))

        self.readb_count += 1

        return buf

    def cwriteb(self, addr, data):
        """Cyclic Block Write
        addr - start address
        data - buffer with buffers with data to be written
        """
        print(
            "cwriteb: addr 0x{:08x}, count {}".format(addr, len(data))
        )

        for dataset in data:
            self.writeb(addr, dataset)

        self.cwriteb_count += 1

    def creadb(self, addr, block_size, count):
        """Cyclic Block Read
        addr - start address
        block_size - block size to read
        count - number of block reads
        """
        print("creadb: addr 0x{:08x}, block size {}, count {}".format(addr, block_size, count))

        buf = []

        for i in range(count):
            buf.append(self.readb(addr, block_size))

        self.creadb_count += 1

        return buf

    def rmw(self, addr, data, mask):
        """Perform read-modify-write operation.
        New data is determined by following formula: X := (X & ~mask) | (data & mask).

        addr - register address
        data - data
        mask - mask
        """
        print(
            "rmw: addr 0x%.8x, data %d (0x%.8x) (%s), mask %d (%s)"
            % (addr, data, data, bin(data), mask, bin(mask))
        )
        X = self.read(addr)
        self.write(addr, (X & abs(mask - 0xFFFFFFFF)) | (data & mask))

        self.rmw_count += 1

    def wait(self, time_ns):
        """Wait in the simulator for a given amount of time.
        time_ns - time to wait in nanoseconds
        """
        assert time_ns > 0, "Wait time must be greater than 0"

        print("wait for %d ns" % time_ns)

        cmd = "T" + ("%.8x" % time_ns) + "\n"
        self.write_fifo.write(cmd)
        self.write_fifo.flush()

        s = self.read_fifo.readline()
        if s.strip() == "ACK":
            return
        else:
            raise Exception("Wrong status returned:" + s.strip())

    def end(self, status):
        """End a co-simulation with a given status.
        status - status to be returned by the simulation process
        """
        print("\n\nCosimIface: ending with status %d" % status)

        cmd = "E" + ("%.8x" % status) + "\n"
        self.write_fifo.write(cmd)
        self.write_fifo.flush()

        s = self.read_fifo.readline()
        if s.strip() == "ACK":
            self._remove_fifos()
        else:
            raise Exception("Wrong status returned:" + s.strip())
        self.print_stats()

    def print_stats(self):
        print(
            f"CosimIface: transactions statistics:\n"
            + f"  read:    {self.read_count}\n"
            + f"  write:   {self.write_count}\n"
            + f"  cread:   {self.cread_count}\n"
            + f"  cwrite:  {self.cwrite_count}\n"
            + f"  readb:   {self.readb_count}\n"
            + f"  writeb:  {self.writeb_count}\n"
            + f"  creadb:  {self.creadb_count}\n"
            + f"  cwriteb: {self.cwriteb_count}\n"
            + f"  rmw:     {self.rmw_count}"
        )
