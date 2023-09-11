# Func tests

## One parameter tests table

```
-------------------------------------------
|                  Param                  |
|-----------------------------------------|
| SingleOneReg | SingleNRegs | ArrayNRegs |
|-----------------------------------------|
|      -       |      -       |   4, 5    |
------------------------------------------|
```

## Two parameters tests matrix

```
------------------------------------------------------------------
|       |                              Second param              |
|       |--------------------------------------------------------|
|       |--------------| SingleOneReg | SingleNRegs | ArrayNRegs |
|-------|--------------------------------------------------------|
|       | SingleOneReg |      1       |      3      |     -      |
| First |--------------|-----------------------------------------|
| param | SingleNRegs  |      2       |      -      |     -      |
|       |--------------|---------------------------------------- |
|       | ArrayNRegs   |      -       |      -      |     -      |
------------------------------------------------------------------
```

## 0
Test checking `proc` without any parameter or return.
In such case only call signal must be generated on address write.

## 1
Test checking `proc` with 2 parameters and no returns.
Both parameters are single with `SingleOneReg` access type, and are placed in the same register.

## 2
Test checking `proc` with 2 parameters and no returns.
First `param` is single with `SingleNRegs` access type, and second `param` is single with `SingleOneReg` access type.
Both parameters are placed within 2 registers.

## 3
Test checking `proc` with 2 parameters and no returns.
First `param` is single with `SingleOneReg` access type ,and second `param` is single with `SingleNRegs` access type,
Both parameters are placed within 2 registers.

## 4
Test checking `proc` with 1 parameter and no returns.
The `param` is array with `ArrayNRegs` access type.
Single argument is narrower than signle register, but no more than one argument fits single register.

## 5
Test checking `proc` with 1 parameter and no returns.
The `param` is array with `ArrayNRegs` access type.
Single register stores 2 or 3 whole items.
