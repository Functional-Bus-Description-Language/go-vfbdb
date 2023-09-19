# Status tests

## 1
Test checking whether the `status` functionality, spanning two registers, with the `SingleNRegs` access type and nonatomic access is read correctly.
The HDL must not contain any extra logic related with the atomic access.

## 2
Test checking access to the `status` array, with `ArrayOneInReg` access type for the array element.

## 3
Test checking access to the `status` array, with `ArrayNInRegs`, `ArrayOneReg`, `ArrayNInRegsMInEndReg` access types for the array element.

## 4
Test checking whether the `status` functionality, spanning 2 registers, with the `SingleNRegs` access type and atomic access is read correctly.
The read order must be increasing.

## 5
Test checking whether the `status` functionality, spanning 3 registers, with the `SingleNRegs` access type and atomic access is read correctly.
The read order must be increasing.

## 6
Test checking access to the `status` array with `ArrayOneInNRegs` nonatomic access.
