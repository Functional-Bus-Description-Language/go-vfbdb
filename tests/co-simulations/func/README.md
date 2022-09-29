# Func tests

## 0
Test checking `func` without any parameter or return.
In such case only strobe must be generated on address write.

## 1
Test checking `func` with 2 parameters and no returns.
Both parameters are single with `Single` access strategy, and are placed in the same register.

## 2
Test checking `func` with 2 parameters and no returns.
First `param` is single with `Continuous` access strategy, and second `param` is single with `Single` access strategy.
Both parameters are placed within 2 registers.

## 3
Test checking `func` with 2 parameters and no returns.
First `param` is single with `Single` access strategy ,and second `param` is single with `Continuous` access strategy,
Both parameters are placed within 2 registers.
