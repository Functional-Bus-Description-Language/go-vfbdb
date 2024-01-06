# Stream tests

## 0
Test checking upstream with single `return` narrower than single register width.

## 1
Test checking both downstream and upstream.
The downstream occupies 2 registers with `SingleNRegs` and `SingleOneReg` params access types.
The upstream occupies 2 registers with one `SingleNRegs` return.
