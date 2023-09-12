# Config tests

## Test table

```
------------------------------------------------
|----------------------|       Atomicity       |
|----------------------------------------------|
|     Access Type      |  Atomic  | Non-atomic |
|----------------------------------------------|
|     SingleOneReg     |    0     |     NA     |
|----------------------------------------------|
|     SingleNRegs      | 2, 3, 4  |     1      |
|----------------------------------------------|
|     ArrayOneReg      |    6     |     NA     |
|----------------------------------------------|
|     ArrayOneInReg    |    5     |     NA     |
|----------------------------------------------|
|     ArrayNInReg      |    7     |     NA     |
|----------------------------------------------|
| ArrayNInRegMInEndReg |    8     |     NA     |
|----------------------------------------------|
|     ArrayNRegs       |    -     |     -      |
------------------------------------------------
```

## 0
Test checking whether the single `config` with `SingleOneReg` access type is correctly written and read.
It also checks whether the value change is visible in the HDL.

## 1
Test checking whether the single `config`, spanning 2 registers, with `SingleNRegs` access type and nonatomic access is correctly written and read.
It also checks whether the value change is visible in the HDL.
The HDL must not contain any extra logic related with the atomic access.

## 2
Test checking whether the single `config`, spanning 2 registers, with `SingleNRegs` access type and atomic access is correctly written and read.
It also checks whether the value change is visible in the HDL.

## 3
Test checking whether the single `config`, spanning 3 registers, with `SingleNRegs` access type and atomic access is correctly written and read.
It also checks whether the value change is visible in the HDL.

## 4
Test checking whether the single `config`, spanning 4 registers, with `SingleNRegs` access type and atomic access is correctly written and read.
It also checks whether the value change is visible in the HDL.

## 5
Test checking whether the array `config` with `ArrayOneInReg` access type and atomic access is correctly written and read.

## 6
Test checking whether the array `config` with `ArrayOneReg` access type and atomic access is correctly written and read.

## 7
Test checking whether the array `config` with `ArrayNInReg` access type and atomic access is correctly written and read.

## 8
Test checking whether the array `config` with `ArrayNInRegMInEnd` access type and atomic access is correctly written and read.
