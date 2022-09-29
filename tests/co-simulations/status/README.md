# Status tests

## 0
Test checking whether the ID and TIMESTAMP `status` elements are generated correctly for the `bus` element.
It also checks whether they are accessible and coherent.

## 1
Test checking whether the `status` element, spanning two registers, with the `Continuous` access strategy and nonatomic access is read correctly.
The HDL must not contain any extra logic related with the atomic access.

## 2
Test checking access to the `status` array, with `Single` access strategy for the array element.

## 3
Test checking access to the `status` array, with `Multiple` access strategy for the array element.

## 4
Test checking whether the `status` element, spanning 2 registers, with the `Continuous` access strategy and atomic access is read correctly.
The read order must be increasing.

## 5
Test checking whether the `status` element, spanning 3 registers, with the `Continuous` access strategy and atomic access is read correctly.
The read order must be increasing.
