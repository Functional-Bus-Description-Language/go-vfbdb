package vhdl

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl"
)

func generateMask(mask *fbdl.Mask, fmts *BlockEntityFormatters) {
	if mask.IsArray {
		generateMaskArray(mask, fmts)
	} else {
		generateMaskSingle(mask, fmts)
	}
}

func generateMaskArray(mask *fbdl.Mask, fmts *BlockEntityFormatters) {
	panic("not yet implemented")
}

func generateMaskSingle(mask *fbdl.Mask, fmts *BlockEntityFormatters) {
	dflt := ""
	if mask.Default != "" {
		dflt = fmt.Sprintf(" := %s", mask.Default.Extend(mask.Width))
	}

	s := fmt.Sprintf(";\n   %s_o : buffer std_logic_vector(%d downto 0)%s", mask.Name, mask.Width-1, dflt)
	fmts.EntityFunctionalPorts += s

	switch mask.Access.(type) {
	case fbdl.AccessSingleSingle:
		generateMaskSingleSingle(mask, fmts)
	default:
		panic("unknown single access strategy")
	}
}

func generateMaskSingleSingle(mask *fbdl.Mask, fmts *BlockEntityFormatters) {
	access := mask.Access.(fbdl.AccessSingleSingle)
	accessMask := access.Mask

	code := fmt.Sprintf(
		"      if master_out.we = '1' then\n"+
			"         %[1]s_o <= master_out.dat(%[2]d downto %[3]d);\n"+
			"      end if;\n"+
			"      master_in.dat(%[2]d downto %[3]d) <= %[1]s_o;",
		mask.Name, accessMask.Upper, accessMask.Lower,
	)

	fmts.RegistersAccess.add([2]int64{access.Addr, access.Addr}, code)
}
