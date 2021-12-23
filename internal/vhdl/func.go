package vhdl

import (
	"fmt"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl"
)

func generateFunc(fun *fbdl.Func, fmts *BlockEntityFormatters) {
	generateFuncType(fun, fmts)
	generateFuncPort(fun, fmts)
	generateFuncAccess(fun, fmts)
	generateFuncStrobe(fun, fmts)
}

func generateFuncType(fun *fbdl.Func, fmts *BlockEntityFormatters) {
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

func generateFuncPort(fun *fbdl.Func, fmts *BlockEntityFormatters) {
	s := fmt.Sprintf(";\n   %s_o : out t_%[1]s", fun.Name)
	fmts.EntityFunctionalPorts += s
}

func generateFuncAccess(fun *fbdl.Func, fmts *BlockEntityFormatters) {
	for _, p := range fun.Params {
		switch p.Access.(type) {
		case fbdl.AccessSingleSingle:
			access := p.Access.(fbdl.AccessSingleSingle)

			addr := [2]int64{access.StartAddr(), access.StartAddr()}
			code := fmt.Sprintf(
				"      %s_o.%s <= master_out.dat(%d downto %d);\n",
				fun.Name, p.Name, access.Mask.Upper, access.Mask.Lower,
			)

			fmts.RegistersAccess.add(addr, code)
		default:
			panic("not yet implemented")
		}
	}
	if len(fun.Params) == 0 {
		fmts.RegistersAccess.add([2]int64{fun.EndAddr(), fun.EndAddr()}, "")
	}
}

func generateFuncStrobe(fun *fbdl.Func, fmts *BlockEntityFormatters) {
	clear := fmt.Sprintf("\n%s_o.stb <= '0';", fun.Name)

	fmts.FuncsStrobesClear += clear

	stb_set := `
   %s_stb : if addr = %d then
      if master_out.we = '1' then
         %[1]s_o.stb <= '1';
      end if;
   end if;
`
	set := fmt.Sprintf(stb_set, fun.Name, fun.EndAddr())

	fmts.FuncsStrobesSet += set
}
