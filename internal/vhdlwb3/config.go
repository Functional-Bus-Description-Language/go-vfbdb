package vhdlwb3

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/access"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/fn"
)

func genConfig(cfg *fn.Config, fmts *BlockEntityFormatters) {
	if cfg.IsArray {
		genConfigArray(cfg, fmts)
	} else {
		genConfigSingle(cfg, fmts)
	}
}

func genConfigArray(cfg *fn.Config, fmts *BlockEntityFormatters) {
	port := fmt.Sprintf(
		";\n   %s_o : buffer slv_vector(%d downto 0)(%d downto 0)",
		cfg.Name, cfg.Count-1, cfg.Width-1,
	)
	if cfg.InitValue != "" {
		port += fmt.Sprintf(" := (others => %s)", cfg.InitValue.Extend(cfg.Width))
	}
	fmts.EntityFunctionalPorts += port

	switch cfg.Access.(type) {
	case access.ArrayOneReg:
		genConfigArrayOneReg(cfg, fmts)
	case access.ArrayOneInReg:
		genConfigArrayOneInReg(cfg, fmts)
	case access.ArrayNInReg:
		genConfigArrayNInReg(cfg, fmts)
	case access.ArrayNInRegMInEndReg:
		genConfigArrayNInRegMInEndReg(cfg, fmts)
	default:
		panic("unimplemented")
	}
}

func genConfigSingle(cfg *fn.Config, fmts *BlockEntityFormatters) {
	dflt := ""
	if cfg.InitValue != "" {
		dflt = fmt.Sprintf(" := %s", cfg.InitValue.Extend(cfg.Width))
	}

	s := fmt.Sprintf(";\n   %s_o : buffer std_logic_vector(%d downto 0)%s", cfg.Name, cfg.Width-1, dflt)
	fmts.EntityFunctionalPorts += s

	switch cfg.Access.(type) {
	case access.SingleOneReg:
		genConfigSingleOneReg(cfg, fmts)
	case access.SingleNRegs:
		genConfigSingleNRegs(cfg, fmts)
	default:
		panic("unimplemented")
	}
}

func genConfigSingleOneReg(cfg *fn.Config, fmts *BlockEntityFormatters) {
	acs := cfg.Access.(access.SingleOneReg)

	code := fmt.Sprintf(`
      if master_out.we = '1' then
         %[1]s_o <= master_out.dat(%[2]d downto %[3]d);
      end if;
      master_in.dat(%[2]d downto %[3]d) <= %[1]s_o;`,
		cfg.Name, acs.EndBit(), acs.StartBit(),
	)

	addr := acs.StartAddr()
	fmts.RegistersAccess.add([2]int64{addr, addr}, code)
}

func genConfigSingleNRegs(cfg *fn.Config, fmts *BlockEntityFormatters) {
	if cfg.Atomic {
		genConfigSingleNRegsAtomic(cfg, fmts)
	} else {
		genConfigSingleNRegsNonAtomic(cfg, fmts)
	}
}

func genConfigSingleNRegsAtomic(cfg *fn.Config, fmts *BlockEntityFormatters) {
	acs := cfg.Access.(access.SingleNRegs)
	strategy := SeparateLast
	atomicShadowRange := [2]int64{cfg.Width - 1 - acs.EndRegWidth(), 0}
	chunks := makeAccessChunksContinuous(acs, strategy)

	fmts.SignalDeclarations += fmt.Sprintf(
		"signal %s_atomic : std_logic_vector(%d downto %d);\n",
		cfg.Name, atomicShadowRange[0], atomicShadowRange[1],
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
				cfg.Name, c.range_[0], c.range_[1], c.endBit, c.startBit,
				atomicShadowRange[0], atomicShadowRange[1],
			)
		} else {
			code = fmt.Sprintf(`
      if master_out.we = '1' then
         %[1]s_atomic(%[2]s downto %[3]s) <= master_out.dat(%[4]d downto %[5]d);
      end if;
      master_in.dat(%[4]d downto %[5]d) <= %[1]s_o(%[2]s downto %[3]s);
`,
				cfg.Name, c.range_[0], c.range_[1], c.endBit, c.startBit,
			)
		}

		fmts.RegistersAccess.add([2]int64{c.addr[0], c.addr[1]}, code)
	}
}

func genConfigSingleNRegsNonAtomic(cfg *fn.Config, fmts *BlockEntityFormatters) {
	acs := cfg.Access.(access.SingleNRegs)
	chunks := makeAccessChunksContinuous(acs, Compact)

	for _, c := range chunks {
		code := fmt.Sprintf(`
      if master_out.we = '1' then
         %[1]s_o(%[2]s downto %[3]s) <= master_out.dat(%[4]d downto %[5]d);
      end if;
      master_in.dat(%[4]d downto %[5]d) <= %[1]s_o(%[2]s downto %[3]s);`,
			cfg.Name, c.range_[0], c.range_[1], c.endBit, c.startBit,
		)

		fmts.RegistersAccess.add([2]int64{c.addr[0], c.addr[1]}, code)
	}
}

func genConfigArrayOneInReg(cfg *fn.Config, fmts *BlockEntityFormatters) {
	acs := cfg.Access.(access.ArrayOneInReg)

	code := fmt.Sprintf(`
      if master_out.we = '1' then
         %[1]s_o(addr - %[2]d) <= master_out.dat(%[3]d downto %[4]d);
      end if;
      master_in.dat(%[3]d downto %[4]d) <= %[1]s_o(addr - %[2]d);`,
		cfg.Name, acs.StartAddr(), acs.EndBit(), acs.StartBit(),
	)

	fmts.RegistersAccess.add(
		[2]int64{acs.StartAddr(), acs.StartAddr() + acs.RegCount() - 1},
		code,
	)
}

func genConfigArrayOneReg(cfg *fn.Config, fmts *BlockEntityFormatters) {
	acs := cfg.Access.(access.ArrayOneReg)

	addr := [2]int64{acs.StartAddr(), acs.EndAddr()}
	code := fmt.Sprintf(`
      for i in 0 to %[1]d loop
         if master_out.we = '1' then
            %[2]s_o(i) <= master_out.dat(%[3]d*(i+1)+%[4]d-1 downto %[3]d*i+%[4]d);
         end if;
         master_in.dat(%[3]d*(i+1)+%[4]d-1 downto %[3]d*i+%[4]d) <= %[2]s_o(i);
      end loop;`,
		cfg.Count-1, cfg.Name, acs.ItemWidth(), acs.StartBit(),
	)

	fmts.RegistersAccess.add(addr, code)
}

func genConfigArrayNInReg(cfg *fn.Config, fmts *BlockEntityFormatters) {
	acs := cfg.Access.(access.ArrayNInReg)

	addr := [2]int64{acs.StartAddr(), acs.EndAddr()}
	code := fmt.Sprintf(`
      for i in 0 to %[1]d loop
         if master_out.we = '1' then
            %[4]s_o((addr-%[5]d)*%[6]d+i) <= master_out.dat(%[2]d*(i+1)+%[3]d-1 downto %[2]d*i+%[3]d);
         end if;
         master_in.dat(%[2]d*(i+1)+%[3]d-1 downto %[2]d*i+%[3]d) <= %[4]s_o((addr-%[5]d)*%[6]d+i);
      end loop;`,
		acs.ItemsInReg()-1, acs.ItemWidth(), acs.StartBit(), cfg.Name, acs.StartAddr(), acs.ItemsInReg(),
	)

	fmts.RegistersAccess.add(addr, code)
}

func genConfigArrayNInRegMInEndReg(cfg *fn.Config, fmts *BlockEntityFormatters) {
	acs := cfg.Access.(access.ArrayNInRegMInEndReg)

	addr := [2]int64{acs.StartAddr(), acs.EndAddr() - 1}
	code := fmt.Sprintf(`
      for i in 0 to %[1]d loop
         if master_out.we = '1' then
            %[4]s_o((addr-%[5]d)*%[6]d+i) <= master_out.dat(%[2]d*(i+1) + %[3]d-1 downto %[2]d*i + %[3]d);
         end if;
         master_in.dat(%[2]d*(i+1) + %[3]d-1 downto %[2]d*i + %[3]d) <= %[4]s_o((addr-%[5]d)*%[6]d+i);
      end loop;`,
		acs.ItemsInReg()-1, acs.ItemWidth(), acs.StartBit(), cfg.Name, acs.StartAddr(), acs.ItemsInReg(),
	)
	fmts.RegistersAccess.add(addr, code)

	addr = [2]int64{acs.EndAddr(), acs.EndAddr()}
	code = fmt.Sprintf(`
      for i in 0 to %[1]d loop
         if master_out.we = '1' then
            %[4]s_o(%[5]d+i) <= master_out.dat(%[2]d*(i+1) + %[3]d-1 downto %[2]d*i+%[3]d);
         end if;
         master_in.dat(%[2]d*(i+1) + %[3]d-1 downto %[2]d*i+%[3]d) <= %[4]s_o(%[5]d+i);
      end loop;`,
		acs.ItemsInEndReg()-1, acs.ItemWidth(), acs.StartBit(), cfg.Name, (acs.RegCount()-1)*acs.ItemsInReg(),
	)

	fmts.RegistersAccess.add(addr, code)
}
