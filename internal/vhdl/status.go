package vhdl

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl"
)

func generateStatus(st *fbdl.Status, fmts *BlockEntityFormatters) {
	if st.IsArray {
		generateStatusArray(st, fmts)
	} else {
		generateStatusSingle(st, fmts)
	}
}

func generateStatusArray(st *fbdl.Status, fmts *BlockEntityFormatters) {
	switch st.Access.(type) {
	case fbdl.AccessArraySingle:
		generateStatusArraySingle(st, fmts)
	case fbdl.AccessArrayMultiple:
		generateStatusArrayMultiple(st, fmts)
	default:
		panic("not yet implemented")
	}
}

func generateStatusSingle(st *fbdl.Status, fmts *BlockEntityFormatters) {
	if st.Name != "X_ID_X" && st.Name != "X_TIMESTAMP_X" {
		s := fmt.Sprintf(";\n   %s_i : in std_logic_vector(%d downto 0)", st.Name, st.Width-1)
		fmts.EntityFunctionalPorts += s
	}

	switch st.Access.(type) {
	case fbdl.AccessSingleSingle:
		generateStatusSingleSingle(st, fmts)
	case fbdl.AccessSingleContinuous:
		generateStatusSingleContinuous(st, fmts)
	default:
		panic("unknown single access strategy")
	}
}

func generateStatusSingleSingle(st *fbdl.Status, fmts *BlockEntityFormatters) {
	fbdlAccess := st.Access.(fbdl.AccessSingleSingle)
	addr := fbdlAccess.Addr
	mask := fbdlAccess.Mask

	var code string
	if st.Name == "X_ID_X" || st.Name == "X_TIMESTAMP_X" {
		code = fmt.Sprintf(
			"      master_in.dat(%d downto %d) <= %s; -- %s",
			mask.Upper, mask.Lower, string(st.Default), st.Name,
		)
	} else {
		code = fmt.Sprintf(
			"      master_in.dat(%d downto %d) <= %s_i;",
			mask.Upper, mask.Lower, st.Name,
		)
	}

	fmts.RegistersAccess.add([2]int64{addr, addr}, code)
}

func generateStatusSingleContinuous(st *fbdl.Status, fmts *BlockEntityFormatters) {
	if st.Atomic == true {
		panic("not yet implemented")
	} else {
		generateStatusSingleContinuousNonAtomic(st, fmts)
	}
}

func generateStatusSingleContinuousNonAtomic(st *fbdl.Status, fmts *BlockEntityFormatters) {
	chunks := makeAccessChunks(st.Access)

	for _, c := range chunks {
		code := fmt.Sprintf(
			"      master_in.dat(%d downto %d) <= %s_i(%s downto %s);",
			c.mask.Upper, c.mask.Lower, st.Name, c.range_[0], c.range_[1],
		)

		fmts.RegistersAccess.add([2]int64{c.addr[0], c.addr[1]}, code)
	}
}

func generateStatusArraySingle(st *fbdl.Status, fmts *BlockEntityFormatters) {
	access := st.Access.(fbdl.AccessArraySingle)

	port := fmt.Sprintf(";\n   %s_i : in t_slv_vector(%d downto 0)(%d downto 0)", st.Name, st.Count-1, st.Width-1)
	fmts.EntityFunctionalPorts += port

	code := fmt.Sprintf(
		"      master_in.dat(%d downto %d) <= %s_i(addr - %d);",
		access.Mask.Upper, access.Mask.Lower, st.Name, access.StartAddr(),
	)

	fmts.RegistersAccess.add(
		[2]int64{access.StartAddr(), access.StartAddr() + access.RegCount() - 1},
		code,
	)
}

func generateStatusArrayMultiple(st *fbdl.Status, fmts *BlockEntityFormatters) {
	access := st.Access.(fbdl.AccessArrayMultiple)

	port := fmt.Sprintf(";\n   %s_i : in t_slv_vector(%d downto 0)(%d downto 0)", st.Name, st.Count-1, st.Width-1)
	fmts.EntityFunctionalPorts += port

	itemsPerAccess := access.ItemsPerAccess

	var addr [2]int64
	var code string

	if access.ItemCount <= itemsPerAccess {
		addr = [2]int64{access.StartAddr(), access.EndAddr()}
		code = fmt.Sprintf(`      for i in 0 to %[1]d loop
         master_in.dat(%[2]d*(i+1)+%[3]d-1 downto %[2]d*i+%[3]d) <= %[4]s_i(i);
      end loop;`,
			st.Count-1, access.ItemWidth, access.StartBit, st.Name,
		)
	} else if access.ItemsInLastReg() == 0 {
		addr = [2]int64{access.StartAddr(), access.EndAddr()}
		code = fmt.Sprintf(`      for i in 0 to %[1]d loop
         master_in.dat(%[2]d*(i+1)+%[3]d-1 downto %[2]d*i+%[3]d) <= %[4]s_i((addr-%[5]d)*%[6]d+i);
      end loop;`,
			itemsPerAccess-1, access.ItemWidth, access.StartBit, st.Name, access.StartAddr(), access.ItemsPerAccess,
		)
	} else {
		addr = [2]int64{access.StartAddr(), access.EndAddr() - 1}
		code = fmt.Sprintf(`      for i in 0 to %[1]d loop
         master_in.dat(%[2]d*(i+1) + %[3]d-1 downto %[2]d*i + %[3]d) <= %[4]s_i((addr-%[5]d)*%[6]d+i);
      end loop;`,
			itemsPerAccess-1, access.ItemWidth, access.StartBit, st.Name, access.StartAddr(), access.ItemsPerAccess,
		)
		fmts.RegistersAccess.add(addr, code)

		addr = [2]int64{access.EndAddr(), access.EndAddr()}
		code = fmt.Sprintf(`      for i in 0 to %[1]d loop
         master_in.dat(%[2]d*(i+1) + %[3]d-1 downto %[2]d*i+%[3]d) <= %[4]s_i(%[5]d+i);
      end loop;`,
			access.ItemsInLastReg()-1, access.ItemWidth, access.StartBit, st.Name, (access.RegCount()-1)*access.ItemsPerAccess,
		)
	}

	fmts.RegistersAccess.add(addr, code)
}
