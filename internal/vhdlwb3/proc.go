package vhdlwb3

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/access"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/elem"
)

func genProc(p *elem.Proc, fmts *BlockEntityFormatters) {
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

func genProcOutType(proc *elem.Proc, fmts *BlockEntityFormatters) {
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

func genProcInType(proc *elem.Proc, fmts *BlockEntityFormatters) {
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

func genProcPorts(proc *elem.Proc, fmts *BlockEntityFormatters) {
	s := fmt.Sprintf(";\n   %s_o : out %[1]s_out_t", proc.Name)
	if len(proc.Returns) != 0 {
		s += fmt.Sprintf(";\n   %s_i : in %[1]s_in_t", proc.Name)
	}
	fmts.EntityFunctionalPorts += s
}

func genProcAccess(proc *elem.Proc, fmts *BlockEntityFormatters) {
	genProcParamsAccess(proc, fmts)
	genProcReturnsAccess(proc, fmts)
}

func genProcParamsAccess(proc *elem.Proc, fmts *BlockEntityFormatters) {
	for _, param := range proc.Params {
		switch param.Access.(type) {
		case access.SingleSingle:
			genProcParamAccessSingleSingle(proc, fmts, param)
		case access.SingleContinuous:
			genProcParamAccessSingleContinuous(proc, fmts, param)
		case access.ArrayContinuous:
			genProcParamAccessArrayContinuous(proc, fmts, param)
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

func genProcParamAccessSingleSingle(proc *elem.Proc, fmts *BlockEntityFormatters, param *elem.Param) {
	a := param.Access.(access.SingleSingle)

	code := fmt.Sprintf(
		"      if master_out.we = '1' then\n"+
			"         %[1]s_o.%[2]s <= master_out.dat(%[3]d downto %[4]d);\n"+
			"      end if;\n"+
			"      master_in.dat(%[3]d downto %[4]d) <= %[1]s_o.%[2]s;\n",
		proc.Name, param.Name, a.EndBit(), a.StartBit(),
	)
	fmts.RegistersAccess.add([2]int64{a.StartAddr(), a.StartAddr()}, code)
}

func genProcParamAccessSingleContinuous(proc *elem.Proc, fmts *BlockEntityFormatters, param *elem.Param) {
	a := param.Access.(access.SingleContinuous)

	chunks := makeAccessChunksContinuous(a, Compact)
	for _, c := range chunks {
		code := fmt.Sprintf(
			"      if master_out.we = '1' then\n"+
				"         %[1]s_o.%[2]s(%[3]s downto %[4]s) <= master_out.dat(%[5]d downto %[6]d);\n"+
				"      end if;\n"+
				"      master_in.dat(%[5]d downto %[6]d) <= %[1]s_o.%[2]s(%[3]s downto %[4]s);\n",
			proc.Name, param.Name, c.range_[0], c.range_[1], c.endBit, c.startBit,
		)
		fmts.RegistersAccess.add([2]int64{c.addr[0], c.addr[1]}, code)
	}
}

func genProcParamAccessArrayContinuous(proc *elem.Proc, fmts *BlockEntityFormatters, param *elem.Param) {
	a := param.Access.(access.ArrayContinuous)

	fmts.SignalDeclarations += fmt.Sprintf(
		"signal %s_%s : slv_vector(%d downto 0)(%d downto 0);\n",
		proc.Name, param.Name, a.RegCount(), busWidth-1,
	)

	code := fmt.Sprintf(
		"      if master_out.we = '1' then\n"+
			"         %s_%s(addr - %d) <= master_out.dat;\n"+
			"      end if;\n",
		proc.Name, param.Name, a.StartAddr(),
	)
	fmts.RegistersAccess.add([2]int64{a.StartAddr(), a.EndAddr()}, code)

	code = fmt.Sprintf(
		"\n%s_%s_driver : process(%[1]s_%[2]s) is\n"+
			"   constant bus_width : natural := %d;\n"+
			"   constant item_width : natural := %d;\n"+
			"   constant item_count : natural := %d;\n"+
			"   variable start_bit : natural;\n"+
			"   variable width : natural;\n"+
			"   variable chunk_width : natural;\n"+
			"   variable item_idx : natural;\n"+
			"   variable next_addr : boolean;\n"+
			"begin\n"+
			"   start_bit := %d;\n"+
			"   width := 0;\n"+
			"   item_idx := 0;\n"+
			"   for addr in 0 to %d loop\n"+
			"      next_addr := false;\n"+
			"      while not next_addr loop\n"+
			"         if item_width - width < bus_width - start_bit then\n"+
			"            chunk_width := item_width - width;\n"+
			"         else\n"+
			"            chunk_width := bus_width - start_bit;\n"+
			"         end if;\n"+
			"         %[1]s_o.%[2]s(item_idx)(chunk_width + width - 1 downto width) <= %[1]s_%[2]s(addr)(chunk_width + start_bit - 1 downto start_bit);\n"+
			"         width := width + chunk_width;\n"+
			"\n         if width = item_width then\n"+
			"            item_idx := item_idx + 1;\n"+
			"            if item_idx = item_count then\n"+
			"               exit;\n"+
			"            end if;\n"+
			"            width := 0;\n"+
			"         end if;\n"+
			"\n         start_bit := start_bit + chunk_width;\n"+
			"         if start_bit = bus_Width then\n"+
			"            next_addr := true;\n"+
			"            start_bit := 0;\n"+
			"         end if;\n"+
			"      end loop;\n" +
			"   end loop;\n" +
			"end process;\n",
		proc.Name, param.Name, busWidth, a.ItemWidth, a.ItemCount, a.StartBit(), a.RegCount()-1,
	)
	fmts.CombinationalProcesses += code

}

func genProcReturnsAccess(proc *elem.Proc, fmts *BlockEntityFormatters) {
	for _, r := range proc.Returns {
		switch r.Access.(type) {
		case access.SingleSingle:
			access := r.Access.(access.SingleSingle)

			addr := [2]int64{access.StartAddr(), access.StartAddr()}
			code := fmt.Sprintf(
				"      master_in.dat(%[1]d downto %[2]d) <= %[3]s_i.%[4]s;\n",
				access.EndBit(), access.StartBit(), proc.Name, r.Name,
			)

			fmts.RegistersAccess.add(addr, code)
		default:
			panic("not yet implemented")
		}
	}
	if (proc.IsEmpty() || proc.IsParam()) && proc.Delay != nil {
		if proc.ExitAddr != nil {
			fmts.RegistersAccess.add([2]int64{*proc.ExitAddr, *proc.ExitAddr}, "")
		}
	}
}

func genProcCall(proc *elem.Proc, fmts *BlockEntityFormatters) {
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

func genProcExit(proc *elem.Proc, fmts *BlockEntityFormatters) {
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
