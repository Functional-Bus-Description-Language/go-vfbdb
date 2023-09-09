# Config tests

## 0
Test checking whether the single `config` with `Single` access strategy is correctly written and read.
It also checks whether the value change is visible in the HDL.

## 1
Test checking whether the single `config`, spanning 2 registers, with `Continuous` access strategy and nonatomic access is correctly written and read.
It also checks whether the value change is visible in the HDL.
The HDL must not contain any extra logic related with the atomic access.

## 2
Test checking whether the single `config`, spanning 2 registers, with `Continuous` access strategy and atomic access is correctly written and read.
It also checks whether the value change is visible in the HDL.

## 3
Test checking whether the single `config`, spanning 3 registers, with `Continuous` access strategy and atomic access is correctly written and read.
It also checks whether the value change is visible in the HDL.

## 4
Test checking whether the single `config`, spanning 4 registers, with `Continuous` access strategy and atomic access is correctly written and read.
It also checks whether the value change is visible in the HDL.

## 5
Test checking whether the array `config` with `Single access` strategy and atomic access is correctly written and read.
