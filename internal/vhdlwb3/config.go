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
	switch cfg.Access.(type) {
	case access.ArraySingle:
		genConfigArraySingle(cfg, fmts)
	case access.ArrayMultiple:
		genConfigArrayMultiple(cfg, fmts)
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
	case access.SingleSingle:
		genConfigSingleSingle(cfg, fmts)
	case access.SingleContinuous:
		genConfigSingleContinuous(cfg, fmts)
	default:
		panic("unknown single access strategy")
	}
}

func genConfigSingleSingle(cfg *fn.Config, fmts *BlockEntityFormatters) {
	a := cfg.Access.(access.SingleSingle)

	code := fmt.Sprintf(`
      if master_out.we = '1' then
         %[1]s_o <= master_out.dat(%[2]d downto %[3]d);
      end if;
      master_in.dat(%[2]d downto %[3]d) <= %[1]s_o;`,
		cfg.Name, a.EndBit(), a.StartBit(),
	)

	fmts.RegistersAccess.add([2]int64{a.Addr, a.Addr}, code)
}

func genConfigSingleContinuous(cfg *fn.Config, fmts *BlockEntityFormatters) {
	if cfg.Atomic {
		genConfigSingleContinuousAtomic(cfg, fmts)
	} else {
		genConfigSingleContinuousNonAtomic(cfg, fmts)
	}
}

func genConfigSingleContinuousAtomic(cfg *fn.Config, fmts *BlockEntityFormatters) {
	a := cfg.Access.(access.SingleContinuous)
	strategy := SeparateLast
	atomicShadowRange := [2]int64{cfg.Width - 1 - a.EndRegWidth(), 0}
	chunks := makeAccessChunksContinuous(a, strategy)

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

func genConfigSingleContinuousNonAtomic(cfg *fn.Config, fmts *BlockEntityFormatters) {
	a := cfg.Access.(access.SingleContinuous)
	chunks := makeAccessChunksContinuous(a, Compact)

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

func genConfigArraySingle(cfg *fn.Config, fmts *BlockEntityFormatters) {
	a := cfg.Access.(access.ArraySingle)

	port := fmt.Sprintf(";\n   %s_o : buffer slv_vector(%d downto 0)(%d downto 0)", cfg.Name, cfg.Count-1, cfg.Width-1)
	fmts.EntityFunctionalPorts += port

	code := fmt.Sprintf(`
      if master_out.we = '1' then
         %[1]s_o(addr - %[2]d) <= master_out.dat(%[3]d downto %[4]d);
      end if;
      master_in.dat(%[3]d downto %[4]d) <= %[1]s_o(addr - %[2]d);`,
		cfg.Name, a.StartAddr(), a.EndBit(), a.StartBit(),
	)

	fmts.RegistersAccess.add(
		[2]int64{a.StartAddr(), a.StartAddr() + a.RegCount() - 1},
		code,
	)
}

func genConfigArrayMultiple(cfg *fn.Config, fmts *BlockEntityFormatters) {
	a := cfg.Access.(access.ArrayMultiple)

	port := fmt.Sprintf(
		";\n   %s_o : buffer slv_vector(%d downto 0)(%d downto 0)",
		cfg.Name, cfg.Count-1, cfg.Width-1,
	)
	fmts.EntityFunctionalPorts += port

	var addr [2]int64
	var code string

	if a.ItemCount <= a.ItemsPerReg {
		addr = [2]int64{a.StartAddr(), a.EndAddr()}
		code = fmt.Sprintf(`
      for i in 0 to %[1]d loop
         if master_out.we = '1' then
            %[2]s_o(i) <= master_out.dat(%[3]d*(i+1)+%[4]d-1 downto %[3]d*i+%[4]d);
         end if;
         master_in.dat(%[3]d*(i+1)+%[4]d-1 downto %[3]d*i+%[4]d) <= %[2]s_o(i);
      end loop;`,
			cfg.Count-1, cfg.Name, a.ItemWidth, a.StartBit(),
		)
	} else if a.ItemsInLastReg() == a.ItemsPerReg {
		addr = [2]int64{a.StartAddr(), a.EndAddr()}
		code = fmt.Sprintf(`
      for i in 0 to %[1]d loop
         if master_out.we = '1' then
            %[4]s_o((addr-%[5]d)*%[6]d+i) <= master_out.dat(%[2]d*(i+1)+%[3]d-1 downto %[2]d*i+%[3]d);
         end if;
         master_in.dat(%[2]d*(i+1)+%[3]d-1 downto %[2]d*i+%[3]d) <= %[4]s_o((addr-%[5]d)*%[6]d+i);
      end loop;`,
			a.ItemsPerReg-1, a.ItemWidth, a.StartBit(), cfg.Name, a.StartAddr(), a.ItemsPerReg,
		)
	} else {
		addr = [2]int64{a.StartAddr(), a.EndAddr() - 1}
		code = fmt.Sprintf(`
      for i in 0 to %[1]d loop
         if master_out.we = '1' then
            %[4]s_o((addr-%[5]d)*%[6]d+i) <= master_out.dat(%[2]d*(i+1) + %[3]d-1 downto %[2]d*i + %[3]d);
         end if;
         master_in.dat(%[2]d*(i+1) + %[3]d-1 downto %[2]d*i + %[3]d) <= %[4]s_o((addr-%[5]d)*%[6]d+i);
      end loop;`,
			a.ItemsPerReg-1, a.ItemWidth, a.StartBit(), cfg.Name, a.StartAddr(), a.ItemsPerReg,
		)
		fmts.RegistersAccess.add(addr, code)

		addr = [2]int64{a.EndAddr(), a.EndAddr()}
		code = fmt.Sprintf(`
      for i in 0 to %[1]d loop
         if master_out.we = '1' then
            %[4]s_o(%[5]d+i) <= master_out.dat(%[2]d*(i+1) + %[3]d-1 downto %[2]d*i+%[3]d);
         end if;
         master_in.dat(%[2]d*(i+1) + %[3]d-1 downto %[2]d*i+%[3]d) <= %[4]s_o(%[5]d+i);
      end loop;`,
			a.ItemsInLastReg()-1, a.ItemWidth, a.StartBit(), cfg.Name, (a.RegCount()-1)*a.ItemsPerReg,
		)
	}

	fmts.RegistersAccess.add(addr, code)
}
