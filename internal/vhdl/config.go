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
	s := fmt.Sprintf(";\n   %s_o : out std_logic_vector(%d downto 0)", cfg.Name, cfg.Width-1)
	fmts.EntityFunctionalPorts += s

	switch cfg.Access.(type) {
	case fbdl.AccessSingleSingle:
		generateConfigSingleSingle(cfg, fmts)
	default:
		panic("unknown single access strategy")
	}
}

func generateConfigSingleSingle(cfg *fbdl.Config, fmts *EntityFormatters) {
	fbdlAccess := cfg.Access.(fbdl.AccessSingleSingle)
	addr := fbdlAccess.Addr
	mask := fbdlAccess.Mask

	access := `
         %[1]s : if internal_addr = %[2]d then
            if internal_master_out.we = '0' then
               internal_master_in.dat(%[3]d downto %[4]d) <= registers(internal_addr)(%[3]d downto %[4]d);
               internal_master_in.ack <= '1';
               internal_master_in.err <= '0';
            end if;
            if internal_master_out.we = '1' then
               registers(internal_addr)(%[3]d downto %[4]d) <= internal_master_out.dat(%[3]d downto %[4]d);
               internal_master_in.ack <= '1';
               internal_master_in.err <= '0';
            end if;
         end if;
`
	access = fmt.Sprintf(access, cfg.Name, addr, mask.Upper, mask.Lower)
	fmts.ConfigsAccess += access

	var routing string
	routing = fmt.Sprintf(
		"   %s_o <= registers(%d)(%d downto %d);\n", cfg.Name, addr, mask.Upper, mask.Lower,
	)

	fmts.ConfigsRouting += routing
}
