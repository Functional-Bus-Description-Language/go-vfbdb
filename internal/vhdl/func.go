package vhdl

import (
	"fmt"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl"
)

func generateFunc(fun *fbdl.Func, fmts *EntityFormatters) {
	generateFuncType(fun, fmts)
	generateFuncPort(fun, fmts)
	generateFuncRouting(fun, fmts)
	generateFuncAccess(fun, fmts)
	generateFuncStrobe(fun, fmts)
}

func generateFuncType(fun *fbdl.Func, fmts *EntityFormatters) {
	s := fmt.Sprintf("\ntype t_%s is record\n", fun.Name)

	for _, p := range fun.Params {
		if p.IsArray {
			s += fmt.Sprintf("   %s : t_slv_vector(%d downto 0)(%d downto 0);\n", p.Name, p.Count-1, p.Width-1)
		} else {
			s += fmt.Sprintf("   %s : std_logic_vector(%d downto 0);\n", p.Name, p.Width-1)
		}
	}

	s += "   stb : std_logic;\nend record;\n"

	fmts.FuncTypes += s
}

func generateFuncPort(fun *fbdl.Func, fmts *EntityFormatters) {
	s := fmt.Sprintf(";\n   %s_o : out t_%[1]s", fun.Name)
	fmts.EntityFunctionalPorts += s
}

func generateFuncRouting(fun *fbdl.Func, fmts *EntityFormatters) {
	s := ""

	for _, p := range fun.Params {
		switch p.Access.(type) {
		case fbdl.AccessSingleSingle:
			a := p.Access.(fbdl.AccessSingleSingle)
			s += fmt.Sprintf(
				"   %s_o.%s <= registers(%d)(%d downto %d);\n",
				fun.Name, p.Name, a.Addr, a.Mask.Upper, a.Mask.Lower,
			)
		default:
			panic("not yet implemented")
		}
	}

	fmts.FuncsRouting += s
}

func generateFuncAccess(fun *fbdl.Func, fmts *EntityFormatters) {
	param_access := `
         %s_%s : if internal_addr = %d then
            if internal_master_out.we = '1' then
               registers(internal_addr)(%d downto %d) <= internal_master_out.dat(%[4]d downto %[5]d);
               internal_master_in.ack <= '1';
               internal_master_in.err <= '0';
            end if;
         end if;
`

	s := ""
	for _, p := range fun.Params {
		switch p.Access.(type) {
		case fbdl.AccessSingleSingle:
			a := p.Access.(fbdl.AccessSingleSingle)
			s += fmt.Sprintf(param_access, fun.Name, p.Name, a.Addr, a.Mask.Upper, a.Mask.Lower)
		default:
			panic("not yet implemented")
		}
	}

	fmts.FuncsAccess += s
}

func generateFuncStrobe(fun *fbdl.Func, fmts *EntityFormatters) {
	clear := fmt.Sprintf("\n      %s_o.stb <= '0';", fun.Name)

	fmts.FuncsStrobesClear += clear

	stb_set := `
         %s_stb : if internal_addr = %d then
            if internal_master_out.we = '1' then
               %[1]s_o.stb <= '1';
            end if;
         end if;
`
	set := fmt.Sprintf(stb_set, fun.Name, fun.EndAddr())

	fmts.FuncsStrobesSet += set
}
