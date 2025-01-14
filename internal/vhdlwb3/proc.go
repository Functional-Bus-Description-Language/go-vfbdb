package vhdlwb3

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/access"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/fn"
)

func genProc(p *fn.Proc, fmts *BlockEntityFormatters) {
	genProcOutType(p, fmts)
	genProcInType(p, fmts)
	genProcPorts(p, fmts)
	genProcAccess(p, fmts)
	if p.CallAddr != nil {
		genProcCall(p, fmts)
	}
	if p.ExitAddr != nil {
		genProcExit(p, fmts)
	}
}

func genProcOutType(proc *fn.Proc, fmts *BlockEntityFormatters) {
	s := fmt.Sprintf("\ntype %s_out_t is record\n", proc.Name)

	for _, p := range proc.Params {
		if p.IsArray {
			s += fmt.Sprintf("   %s : slv_vector(%d downto 0)(%d downto 0);\n", p.Name, p.Count-1, p.Width-1)
		} else {
			s += fmt.Sprintf("   %s : std_logic_vector(%d downto 0);\n", p.Name, p.Width-1)
		}
	}

	s += "   call : std_logic;\n"

	if len(proc.Returns) != 0 {
		s += "   exitt : std_logic;\n"
	}
	s += "end record;\n"

	fmts.ProcTypes += s
}

func genProcInType(proc *fn.Proc, fmts *BlockEntityFormatters) {
	if len(proc.Returns) == 0 {
		return
	}

	s := fmt.Sprintf("\ntype %s_in_t is record\n", proc.Name)

	for _, r := range proc.Returns {
		if r.IsArray {
			s += fmt.Sprintf("   %s : slv_vector(%d downto 0)(%d downto 0);\n", r.Name, r.Count-1, r.Width-1)
		} else {
			s += fmt.Sprintf("   %s : std_logic_vector(%d downto 0);\n", r.Name, r.Width-1)
		}
	}

	s += "end record;\n"

	fmts.ProcTypes += s
}

func genProcPorts(proc *fn.Proc, fmts *BlockEntityFormatters) {
	s := fmt.Sprintf(";\n   %s_o : out %[1]s_out_t", proc.Name)
	if len(proc.Returns) != 0 {
		s += fmt.Sprintf(";\n   %s_i : in %[1]s_in_t", proc.Name)
	}
	fmts.EntityFunctionalPorts += s
}

func genProcAccess(proc *fn.Proc, fmts *BlockEntityFormatters) {
	genProcParamsAccess(proc, fmts)
	genProcReturnsAccess(proc, fmts)
}

func genProcParamsAccess(proc *fn.Proc, fmts *BlockEntityFormatters) {
	for _, param := range proc.Params {
		switch param.Access.(type) {
		case access.SingleOneReg:
			genProcParamAccessSingleOneReg(proc, fmts, param)
		case access.SingleNRegs:
			genProcParamAccessSingleNRegs(proc, fmts, param)
		case access.ArrayNRegs:
			genProcParamAccessArrayNRegs(proc, fmts, param)
		default:
			panic("should never happen")
		}
	}

	if proc.IsEmpty() || (proc.IsReturn() && proc.Delay != nil) {
		if proc.CallAddr != nil {
			fmts.RegistersAccess.add([2]int64{*proc.CallAddr, *proc.CallAddr}, "")
		}
	}
}

func genProcParamAccessSingleOneReg(proc *fn.Proc, fmts *BlockEntityFormatters, param *fn.Param) {
	acs := param.Access.(access.SingleOneReg)

	code := fmt.Sprintf(`
      if master_out.we = '1' then
         %[1]s_o.%[2]s <= master_out.dat(%[3]d downto %[4]d);
      end if;
      master_in.dat(%[3]d downto %[4]d) <= %[1]s_o.%[2]s;`,
		proc.Name, param.Name, acs.EndBit(), acs.StartBit(),
	)
	addr := acs.StartAddr()
	fmts.RegistersAccess.add([2]int64{addr, addr}, code)
}

func genProcParamAccessSingleNRegs(proc *fn.Proc, fmts *BlockEntityFormatters, param *fn.Param) {
	acs := param.Access.(access.SingleNRegs)

	chunks := makeAccessChunksContinuous(acs, Compact)
	for _, c := range chunks {
		code := fmt.Sprintf(`
      if master_out.we = '1' then
         %[1]s_o.%[2]s(%[3]s downto %[4]s) <= master_out.dat(%[5]d downto %[6]d);
      end if;
      master_in.dat(%[5]d downto %[6]d) <= %[1]s_o.%[2]s(%[3]s downto %[4]s);`,
			proc.Name, param.Name, c.range_[0], c.range_[1], c.endBit, c.startBit,
		)
		fmts.RegistersAccess.add([2]int64{c.addr[0], c.addr[1]}, code)
	}
}

func genProcParamAccessArrayNRegs(proc *fn.Proc, fmts *BlockEntityFormatters, param *fn.Param) {
	acs := param.Access.(access.ArrayNRegs)

	fmts.SignalDeclarations += fmt.Sprintf(
		"signal %s_%s : slv_vector(%d downto 0)(%d downto 0);\n",
		proc.Name, param.Name, acs.RegCount(), busWidth-1,
	)

	code := fmt.Sprintf(`
      if master_out.we = '1' then
         %s_%s(addr - %d) <= master_out.dat;
      end if;`,
		proc.Name, param.Name, acs.StartAddr(),
	)
	fmts.RegistersAccess.add([2]int64{acs.StartAddr(), acs.EndAddr()}, code)

	code = fmt.Sprintf(
		`
%s_%s_driver : process(%[1]s_%[2]s) is
   variable start_bit : natural;
   variable width : natural;
   variable chunk_width : natural;
   variable item_idx : natural;
begin
   start_bit := %[6]d;
   width := 0;
   item_idx := 0;
   for addr in 0 to %[7]d loop
      loop
         chunk_width := %[4]d - width when %[4]d - width < %[3]d - start_bit else %[3]d - start_bit;
         %[1]s_o.%[2]s(item_idx)(chunk_width + width - 1 downto width) <= %[1]s_%[2]s(addr)(chunk_width + start_bit - 1 downto start_bit);

         width := width + chunk_width;
         if width = %[4]d then
            item_idx := item_idx + 1;
            exit when item_idx = %[5]d;
            width := 0;
         end if;

         start_bit := start_bit + chunk_width;
         if start_bit = %[3]d then
            start_bit := 0;
            exit;
         end if;
      end loop;
   end loop;
end process;
`,
		proc.Name, param.Name, busWidth, acs.ItemWidth(), acs.ItemCount(), acs.StartBit(), acs.RegCount()-1,
	)
	fmts.CombinationalProcesses += code
}

func genProcReturnsAccess(proc *fn.Proc, fmts *BlockEntityFormatters) {
	for _, r := range proc.Returns {
		switch acs := r.Access.(type) {
		case access.SingleOneReg:
			addr := [2]int64{acs.StartAddr(), acs.StartAddr()}
			code := fmt.Sprintf(
				"      master_in.dat(%[1]d downto %[2]d) <= %[3]s_i.%[4]s;\n",
				acs.EndBit(), acs.StartBit(), proc.Name, r.Name,
			)

			fmts.RegistersAccess.add(addr, code)
		default:
			panic("unimplemented")
		}
	}
	if (proc.IsEmpty() || proc.IsParam()) && proc.Delay != nil {
		if proc.ExitAddr != nil {
			fmts.RegistersAccess.add([2]int64{*proc.ExitAddr, *proc.ExitAddr}, "")
		}
	}
}

func genProcCall(proc *fn.Proc, fmts *BlockEntityFormatters) {
	clear := fmt.Sprintf("\n%s_o.call <= '0';", proc.Name)

	fmts.ProcsCallsClear += clear

	callSet := `
   %s_call : if addr = %d then
      if master_out.we = '1' then
         %[1]s_o.call <= '1';
      end if;
   end if;
`
	set := fmt.Sprintf(callSet, proc.Name, *proc.CallAddr)

	fmts.ProcsCallsSet += set
}

func genProcExit(proc *fn.Proc, fmts *BlockEntityFormatters) {
	if len(proc.Returns) == 0 {
		return
	}

	clear := fmt.Sprintf("\n%s_o.exitt <= '0';", proc.Name)

	fmts.ProcsExitsClear += clear

	exitSet := `
   %s_exit : if addr = %d then
      if master_out.we = '0' then
         %[1]s_o.exitt <= '1';
      end if;
   end if;
`
	set := fmt.Sprintf(exitSet, proc.Name, *proc.ExitAddr)

	fmts.ProcsExitsSet += set
}
