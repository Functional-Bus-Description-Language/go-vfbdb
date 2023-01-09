package vhdlwb3

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/access"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/elem"
)

func genProc(p *elem.Proc, fmts *BlockEntityFormatters) {
	genProcType(p, fmts)
	genProcPort(p, fmts)
	genProcAccess(p, fmts)
	genProcStrobe(p, fmts)
}

func genProcType(proc *elem.Proc, fmts *BlockEntityFormatters) {
	s := fmt.Sprintf("\ntype %s_t is record\n", proc.Name)

	for _, p := range proc.Params {
		if p.IsArray {
			s += fmt.Sprintf("   %s : slv_vector(%d downto 0)(%d downto 0);\n", p.Name, p.Count-1, p.Width-1)
		} else {
			s += fmt.Sprintf("   %s : std_logic_vector(%d downto 0);\n", p.Name, p.Width-1)
		}
	}

	s += "   stb : std_logic;\nend record;\n"

	fmts.ProcTypes += s
}

func genProcPort(proc *elem.Proc, fmts *BlockEntityFormatters) {
	s := fmt.Sprintf(";\n   %s_o : out %[1]s_t", proc.Name)
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
	if len(proc.Params) == 0 {
		fmts.RegistersAccess.add([2]int64{proc.StbAddr, proc.StbAddr}, "")
	}
}

func genProcStrobe(proc *elem.Proc, fmts *BlockEntityFormatters) {
	clear := fmt.Sprintf("\n%s_o.stb <= '0';", proc.Name)

	fmts.ProcsStrobesClear += clear

	stbSet := `
   %s_stb : if addr = %d then
      if master_out.we = '1' then
         %[1]s_o.stb <= '1';
      end if;
   end if;
`
	set := fmt.Sprintf(stbSet, proc.Name, proc.StbAddr)

	fmts.ProcsStrobesSet += set
}
