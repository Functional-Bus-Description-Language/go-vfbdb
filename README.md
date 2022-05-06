[![Tests](https://github.com/Functional-Bus-Description-Language/go-vfbdb/actions/workflows/tests.yml/badge.svg?branch=master)](https://github.com/Functional-Bus-Description-Language/go-vfbdb/actions?query=master)

# go-vfbdb

Versatile Functional Bus Description Language compiler backend written in Go.

Supported targets:
- c-sync - C target with synchronous (blocking) interface functions,
- python - Python target,
- vhdl-wb3 - VHDL target for Wishbone compilant with revision B.3.

## Installation

### go
```
go install github.com/Functional-Bus-Description-Language/go-vfbdb/cmd/vfbdb@latest
```

Go installation installs to go configured path.

### Manual

```
git clone https://github.com/Functional-Bus-Description-Language/go-vfbdb.git
make
make install
```

Manual installation installs to `/usr/bin`.
