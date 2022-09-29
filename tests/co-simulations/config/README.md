# Config tests

## 0
Test checking whether the single `config` with `Single` access strategy is correctly written and read.
It also checks whether the value change is visible in the HDL.

## 1
Test checking whether the single `config`, spanning 2 registers, with `Continuous` access strategy and nonatomic access is correctly written and read.
It also checks whether the value change is visible in the HDL.
The HDL must not contain any extra logic related with the atomic access.
