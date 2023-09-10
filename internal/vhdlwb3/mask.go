package vhdlwb3

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/access"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/fn"
)

func genMask(mask *fn.Mask, fmts *BlockEntityFormatters) {
	if mask.IsArray {
		genMaskArray(mask, fmts)
	} else {
		genMaskSingle(mask, fmts)
	}
}

func genMaskArray(mask *fn.Mask, fmts *BlockEntityFormatters) {
	panic("not yet implemented")
}

func genMaskSingle(mask *fn.Mask, fmts *BlockEntityFormatters) {
	dflt := ""
	if mask.InitValue != "" {
		dflt = fmt.Sprintf(" := %s", mask.InitValue.Extend(mask.Width))
	}

	s := fmt.Sprintf(";\n   %s_o : buffer std_logic_vector(%d downto 0)%s", mask.Name, mask.Width-1, dflt)
	fmts.EntityFunctionalPorts += s

	switch mask.Access.(type) {
	case access.SingleSingle:
		genMaskSingleSingle(mask, fmts)
	default:
		panic("unknown single access strategy")
	}
}

func genMaskSingleSingle(mask *fn.Mask, fmts *BlockEntityFormatters) {
	a := mask.Access.(access.SingleSingle)

	code := fmt.Sprintf(`
      if master_out.we = '1' then
         %[1]s_o <= master_out.dat(%[2]d downto %[3]d);
      end if;
      master_in.dat(%[2]d downto %[3]d) <= %[1]s_o;`,
		mask.Name, a.GetEndBit(), a.GetStartBit(),
	)

	fmts.RegistersAccess.add([2]int64{a.Addr, a.Addr}, code)
}
