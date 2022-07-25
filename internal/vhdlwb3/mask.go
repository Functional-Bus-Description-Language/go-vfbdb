package vhdlwb3

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/access"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/elem"
)

func genMask(mask elem.Mask, fmts *BlockEntityFormatters) {
	if mask.IsArray() {
		genMaskArray(mask, fmts)
	} else {
		genMaskSingle(mask, fmts)
	}
}

func genMaskArray(mask elem.Mask, fmts *BlockEntityFormatters) {
	panic("not yet implemented")
}

func genMaskSingle(mask elem.Mask, fmts *BlockEntityFormatters) {
	dflt := ""
	if mask.Default() != "" {
		dflt = fmt.Sprintf(" := %s", mask.Default().Extend(mask.Width()))
	}

	s := fmt.Sprintf(";\n   %s_o : buffer std_logic_vector(%d downto 0)%s", mask.Name(), mask.Width()-1, dflt)
	fmts.EntityFunctionalPorts += s

	switch mask.Access().(type) {
	case access.SingleSingle:
		genMaskSingleSingle(mask, fmts)
	default:
		panic("unknown single access strategy")
	}
}

func genMaskSingleSingle(mask elem.Mask, fmts *BlockEntityFormatters) {
	a := mask.Access().(access.SingleSingle)

	code := fmt.Sprintf(
		"      if master_out.we = '1' then\n"+
			"         %[1]s_o <= master_out.dat(%[2]d downto %[3]d);\n"+
			"      end if;\n"+
			"      master_in.dat(%[2]d downto %[3]d) <= %[1]s_o;",
		mask.Name(), a.EndBit(), a.StartBit(),
	)

	fmts.RegistersAccess.add([2]int64{a.Addr, a.Addr}, code)
}
