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
	panic("unimplemented")
}

func genMaskSingle(mask *fn.Mask, fmts *BlockEntityFormatters) {
	dflt := ""
	if mask.InitValue != "" {
		dflt = fmt.Sprintf(" := %s", mask.InitValue.Extend(mask.Width))
	}

	s := fmt.Sprintf(";\n   %s_o : buffer std_logic_vector(%d downto 0)%s", mask.Name, mask.Width-1, dflt)
	fmts.EntityFunctionalPorts += s

	switch mask.Access.(type) {
	case access.SingleOneReg:
		genMaskSingleOneReg(mask, fmts)
	case access.SingleNRegs:
		genMaskSingleNRegs(mask, fmts)
	default:
		panic("unimplemented")
	}
}

func genMaskSingleOneReg(mask *fn.Mask, fmts *BlockEntityFormatters) {
	acs := mask.Access.(access.SingleOneReg)

	code := fmt.Sprintf(`
      if master_out.we = '1' then
         %[1]s_o <= master_out.dat(%[2]d downto %[3]d);
      end if;
      master_in.dat(%[2]d downto %[3]d) <= %[1]s_o;`,
		mask.Name, acs.EndBit(), acs.StartBit(),
	)

	addr := acs.StartAddr()
	fmts.RegistersAccess.add([2]int64{addr, addr}, code)
}

func genMaskSingleNRegs(mask *fn.Mask, fmts *BlockEntityFormatters) {
	if mask.Atomic {
		genMaskSingleNRegsAtomic(mask, fmts)
	} else {
		genMaskSingleNRegsNonAtomic(mask, fmts)
	}
}

func genMaskSingleNRegsAtomic(mask *fn.Mask, fmts *BlockEntityFormatters) {
	acs := mask.Access.(access.SingleNRegs)
	strategy := SeparateLast
	atomicShadowRange := [2]int64{mask.Width - 1 - acs.EndRegWidth(), 0}
	chunks := makeAccessChunksContinuous(acs, strategy)

	fmts.SignalDeclarations += fmt.Sprintf(
		"signal %s_atomic : std_logic_vector(%d downto %d);\n",
		mask.Name, atomicShadowRange[0], atomicShadowRange[1],
	)

	for i, c := range chunks {
		var code string
		if (strategy == SeparateFirst && i == 0) || (strategy == SeparateLast && i == len(chunks)-1) {
			code = fmt.Sprintf(`
      if master_out.we = '1' then
         %[1]s_o(%[2]s downto %[3]s) <= master_out.dat(%[4]d downto %[5]d);
         %[1]s_o(%[6]d downto %[7]d) <= %[1]s_atomic(%[6]d downto %[7]d);
      end if;
      master_in.dat(%[4]d downto %[5]d) <= %[1]s_o(%[2]s downto %[3]s);`,
				mask.Name, c.range_[0], c.range_[1], c.endBit, c.startBit,
				atomicShadowRange[0], atomicShadowRange[1],
			)
		} else {
			code = fmt.Sprintf(`
      if master_out.we = '1' then
         %[1]s_atomic(%[2]s downto %[3]s) <= master_out.dat(%[4]d downto %[5]d);
      end if;
      master_in.dat(%[4]d downto %[5]d) <= %[1]s_o(%[2]s downto %[3]s);
`,
				mask.Name, c.range_[0], c.range_[1], c.endBit, c.startBit,
			)
		}

		fmts.RegistersAccess.add([2]int64{c.addr[0], c.addr[1]}, code)
	}
}

func genMaskSingleNRegsNonAtomic(mask *fn.Mask, fmts *BlockEntityFormatters) {
	acs := mask.Access.(access.SingleNRegs)
	chunks := makeAccessChunksContinuous(acs, Compact)

	for _, c := range chunks {
		code := fmt.Sprintf(`
      if master_out.we = '1' then
         %[1]s_o(%[2]s downto %[3]s) <= master_out.dat(%[4]d downto %[5]d);
      end if;
      master_in.dat(%[4]d downto %[5]d) <= %[1]s_o(%[2]s downto %[3]s);`,
			mask.Name, c.range_[0], c.range_[1], c.endBit, c.startBit,
		)

		fmts.RegistersAccess.add([2]int64{c.addr[0], c.addr[1]}, code)
	}
}
