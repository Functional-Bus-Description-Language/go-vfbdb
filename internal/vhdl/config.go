package vhdl

import (
	"fmt"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl"
)

func generateConfig(cfg *fbdl.Config, fmts *EntityFormatters) {
	if cfg.IsArray {
		generateConfigArray(cfg, fmts)
	} else {
		generateConfigSingle(cfg, fmts)
	}
}

func generateConfigArray(cfg *fbdl.Config, fmts *EntityFormatters) {
	panic("not yet implemented")
}

func generateConfigSingle(cfg *fbdl.Config, fmts *EntityFormatters) {
	dflt := ""
	if cfg.Default != "" {
		dflt = fmt.Sprintf(" := %s", cfg.Default.Extend(cfg.Width))
	}

	s := fmt.Sprintf(";\n   %s_o : buffer std_logic_vector(%d downto 0)%s", cfg.Name, cfg.Width-1, dflt)
	fmts.EntityFunctionalPorts += s

	switch cfg.Access.(type) {
	case fbdl.AccessSingleSingle:
		generateConfigSingleSingle(cfg, fmts)
	default:
		panic("unknown single access strategy")
	}
}

func generateConfigSingleSingle(cfg *fbdl.Config, fmts *EntityFormatters) {
	access := cfg.Access.(fbdl.AccessSingleSingle)
	mask := access.Mask

	code := fmt.Sprintf(
		"      if internal_master_out.we = '1' then\n"+
			"         %[1]s_o <= internal_master_out.dat(%[2]d downto %[3]d);\n"+
			"      end if;\n"+
			"      internal_master_in.dat(%[2]d downto %[3]d) <= %[1]s_o;",
		cfg.Name, mask.Upper, mask.Lower,
	)

	fmts.RegistersAccess.add([2]int64{access.Addr, access.Addr}, code)
}
