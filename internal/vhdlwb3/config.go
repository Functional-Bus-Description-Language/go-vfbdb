package vhdlwb3

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/access"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/elem"
)

func genConfig(cfg elem.Config, fmts *BlockEntityFormatters) {
	if cfg.IsArray() {
		genConfigArray(cfg, fmts)
	} else {
		genConfigSingle(cfg, fmts)
	}
}

func genConfigArray(cfg elem.Config, fmts *BlockEntityFormatters) {
	panic("not yet implemented")
}

func genConfigSingle(cfg elem.Config, fmts *BlockEntityFormatters) {
	dflt := ""
	if cfg.Default() != "" {
		dflt = fmt.Sprintf(" := %s", cfg.Default().Extend(cfg.Width()))
	}

	s := fmt.Sprintf(";\n   %s_o : buffer std_logic_vector(%d downto 0)%s", cfg.Name(), cfg.Width()-1, dflt)
	fmts.EntityFunctionalPorts += s

	switch cfg.Access().(type) {
	case access.SingleSingle:
		genConfigSingleSingle(cfg, fmts)
	case access.SingleContinuous:
		genConfigSingleContinuous(cfg, fmts)
	default:
		panic("unknown single access strategy")
	}
}

func genConfigSingleSingle(cfg elem.Config, fmts *BlockEntityFormatters) {
	access := cfg.Access().(access.SingleSingle)
	mask := access.Mask

	code := fmt.Sprintf(
		"      if master_out.we = '1' then\n"+
			"         %[1]s_o <= master_out.dat(%[2]d downto %[3]d);\n"+
			"      end if;\n"+
			"      master_in.dat(%[2]d downto %[3]d) <= %[1]s_o;",
		cfg.Name(), mask.Upper, mask.Lower,
	)

	fmts.RegistersAccess.add([2]int64{access.Addr, access.Addr}, code)
}

func genConfigSingleContinuous(cfg elem.Config, fmts *BlockEntityFormatters) {
	if cfg.Atomic() == true {
		genConfigSingleContinuousAtomic(cfg, fmts)
	} else {
		genConfigSingleContinuousNonAtomic(cfg, fmts)
	}
}

func genConfigSingleContinuousAtomic(cfg elem.Config, fmts *BlockEntityFormatters) {
	a := cfg.Access().(access.SingleContinuous)
	strategy := SeparateLast
	atomicShadowRange := [2]int64{cfg.Width() - 1 - a.EndMask.Width(), 0}
	if cfg.HasDecreasingAccessOrder() {
		strategy = SeparateFirst
		atomicShadowRange[0] = cfg.Width() - 1
		atomicShadowRange[1] = a.StartMask.Width()
	}
	chunks := makeAccessChunksContinuous(a, strategy)

	fmts.SignalDeclarations += fmt.Sprintf(
		"signal %s_atomic : std_logic_vector(%d downto %d);\n",
		cfg.Name(), atomicShadowRange[0], atomicShadowRange[1],
	)

	for i, c := range chunks {
		var code string
		if (strategy == SeparateFirst && i == 0) || (strategy == SeparateLast && i == len(chunks)-1) {
			code = fmt.Sprintf(
				"      if master_out.we = '1' then\n"+
					"         %[1]s_o(%[2]s downto %[3]s) <= master_out.dat(%[4]d downto %[5]d);\n"+
					"         %[1]s_o(%[6]d downto %[7]d) <= %[1]s_atomic(%[6]d downto %[7]d);\n"+
					"      end if;\n"+
					"      master_in.dat(%[4]d downto %[5]d) <= %[1]s_o(%[2]s downto %[3]s);",
				cfg.Name(), c.range_[0], c.range_[1], c.mask.Upper, c.mask.Lower,
				atomicShadowRange[0], atomicShadowRange[1],
			)
		} else {
			code = fmt.Sprintf(
				"      if master_out.we = '1' then\n"+
					"         %[1]s_atomic(%[2]s downto %[3]s) <= master_out.dat(%[4]d downto %[5]d);\n"+
					"      end if;\n"+
					"      master_in.dat(%[4]d downto %[5]d) <= %[1]s_o(%[2]s downto %[3]s);",
				cfg.Name(), c.range_[0], c.range_[1], c.mask.Upper, c.mask.Lower,
			)
		}

		fmts.RegistersAccess.add([2]int64{c.addr[0], c.addr[1]}, code)
	}
}

func genConfigSingleContinuousNonAtomic(cfg elem.Config, fmts *BlockEntityFormatters) {
	a := cfg.Access().(access.SingleContinuous)
	chunks := makeAccessChunksContinuous(a, Compact)

	for _, c := range chunks {
		code := fmt.Sprintf(
			"      if master_out.we = '1' then\n"+
				"         %[1]s_o(%[2]s downto %[3]s) <= master_out.dat(%[4]d downto %[5]d);\n"+
				"      end if;\n"+
				"      master_in.dat(%[4]d downto %[5]d) <= %[1]s_o(%[2]s downto %[3]s);",
			cfg.Name(), c.range_[0], c.range_[1], c.mask.Upper, c.mask.Lower,
		)

		fmts.RegistersAccess.add([2]int64{c.addr[0], c.addr[1]}, code)
	}
}
