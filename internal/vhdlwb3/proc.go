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
	for _, p := range proc.Params {
		switch p.Access.(type) {
		case access.SingleSingle:
			access := p.Access.(access.SingleSingle)

			addr := [2]int64{access.StartAddr(), access.StartAddr()}
			code := fmt.Sprintf(
				"      if master_out.we = '1' then\n"+
					"         %[1]s_o.%[2]s <= master_out.dat(%[3]d downto %[4]d);\n"+
					"      end if;\n"+
					"      master_in.dat(%[3]d downto %[4]d) <= %[1]s_o.%[2]s;\n",
				proc.Name, p.Name, access.EndBit(), access.StartBit(),
			)

			fmts.RegistersAccess.add(addr, code)
		case access.SingleContinuous:
			chunks := makeAccessChunksContinuous(p.Access.(access.SingleContinuous), Compact)

			for _, c := range chunks {
				code := fmt.Sprintf(
					"      if master_out.we = '1' then\n"+
						"         %[1]s_o.%[2]s(%[3]s downto %[4]s) <= master_out.dat(%[5]d downto %[6]d);\n"+
						"      end if;\n"+
						"      master_in.dat(%[5]d downto %[6]d) <= %[1]s_o.%[2]s(%[3]s downto %[4]s);\n",
					proc.Name, p.Name, c.range_[0], c.range_[1], c.endBit, c.startBit,
				)

				fmts.RegistersAccess.add([2]int64{c.addr[0], c.addr[1]}, code)
			}
		default:
			panic("not yet implemented")
		}
	}
	if proc.CallAddr != nil {
		fmts.RegistersAccess.add([2]int64{*proc.CallAddr, *proc.CallAddr}, "")
	}
	if proc.ExitAddr != nil {
		fmts.RegistersAccess.add([2]int64{*proc.ExitAddr, *proc.ExitAddr}, "")
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
