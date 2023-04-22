# Func tests

## One parameter tests table

```
------------------------------------------------------
|                    Param                            |
|-----------------------------------------------------|
| Single Single | Single Continuous | Array Continuous|
|-----------------------------------------------------|
|       -       |         -         |       4         |
------------------------------------------------------|
```

## Two parameters tests matrix

```
-----------------------------------------------------------------------------------
|       |                                        Second param                     |
|       |-------------------------------------------------------------------------|
|       | ----------------- | Single Single | Single Continuous | Array Continuous|
|-------|-------------------|-----------------------------------------------------|
|       | Single Single     |       1       |         3         |       -         |
| First |-------------------|---------------------------------------------------- |
| param | Single Continuous |       2       |         -         |       -         |
|       |-------------------|---------------------------------------------------- |
|       | Array Continuous  |       -       |         -         |       -         |
----------------------------------------------------------------------------------
```

## 0
Test checking `proc` without any parameter or return.
In such case only call signal must be generated on address write.

## 1
Test checking `proc` with 2 parameters and no returns.
Both parameters are single with `Single` access strategy, and are placed in the same register.

## 2
Test checking `proc` with 2 parameters and no returns.
First `param` is single with `Continuous` access strategy, and second `param` is single with `Single` access strategy.
Both parameters are placed within 2 registers.

## 3
Test checking `proc` with 2 parameters and no returns.
First `param` is single with `Single` access strategy ,and second `param` is single with `Continuous` access strategy,
Both parameters are placed within 2 registers.

## 4
Test checking `proc` with 1 parameter and no returns.
The `param` is array with `Continuous` access strategy.
Single argument is narrower than signle register, but no more than one argument fits single register.
