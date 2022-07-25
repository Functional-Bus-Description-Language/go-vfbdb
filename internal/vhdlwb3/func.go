package vhdlwb3

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/access"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/elem"
)

func genFunc(fun elem.Func, fmts *BlockEntityFormatters) {
	genFuncType(fun, fmts)
	genFuncPort(fun, fmts)
	genFuncAccess(fun, fmts)
	genFuncStrobe(fun, fmts)
}

func genFuncType(fun elem.Func, fmts *BlockEntityFormatters) {
	s := fmt.Sprintf("\ntype t_%s is record\n", fun.Name())

	for _, p := range fun.Params() {
		if p.IsArray() {
			s += fmt.Sprintf("   %s : slv_vector(%d downto 0)(%d downto 0);\n", p.Name(), p.Count()-1, p.Width()-1)
		} else {
			s += fmt.Sprintf("   %s : std_logic_vector(%d downto 0);\n", p.Name(), p.Width()-1)
		}
	}

	s += "   stb : std_logic;\nend record;\n"

	fmts.FuncTypes += s
}

func genFuncPort(fun elem.Func, fmts *BlockEntityFormatters) {
	s := fmt.Sprintf(";\n   %s_o : out t_%[1]s", fun.Name())
	fmts.EntityFunctionalPorts += s
}

func genFuncAccess(fun elem.Func, fmts *BlockEntityFormatters) {
	for _, p := range fun.Params() {
		switch p.Access().(type) {
		case access.SingleSingle:
			access := p.Access().(access.SingleSingle)

			addr := [2]int64{access.StartAddr(), access.StartAddr()}
			code := fmt.Sprintf(
				"      if master_out.we = '1' then\n"+
					"         %[1]s_o.%[2]s <= master_out.dat(%[3]d downto %[4]d);\n"+
					"      end if;\n"+
					"      master_in.dat(%[3]d downto %[4]d) <= %[1]s_o.%[2]s;\n",
				fun.Name(), p.Name(), access.EndBit(), access.StartBit(),
			)

			fmts.RegistersAccess.add(addr, code)
		case access.SingleContinuous:
			chunks := makeAccessChunksContinuous(p.Access().(access.SingleContinuous), Compact)

			for _, c := range chunks {
				code := fmt.Sprintf(
					"      if master_out.we = '1' then\n"+
						"         %[1]s_o.%[2]s(%[3]s downto %[4]s) <= master_out.dat(%[5]d downto %[6]d);\n"+
						"      end if;\n"+
						"      master_in.dat(%[5]d downto %[6]d) <= %[1]s_o.%[2]s(%[3]s downto %[4]s);\n",
					fun.Name(), p.Name(), c.range_[0], c.range_[1], c.endBit, c.startBit,
				)

				fmts.RegistersAccess.add([2]int64{c.addr[0], c.addr[1]}, code)
			}
		default:
			panic("not yet implemented")
		}
	}
	if len(fun.Params()) == 0 {
		fmts.RegistersAccess.add([2]int64{fun.StbAddr(), fun.StbAddr()}, "")
	}
}

func genFuncStrobe(fun elem.Func, fmts *BlockEntityFormatters) {
	clear := fmt.Sprintf("\n%s_o.stb <= '0';", fun.Name())

	fmts.FuncsStrobesClear += clear

	stbSet := `
   %s_stb : if addr = %d then
      if master_out.we = '1' then
         %[1]s_o.stb <= '1';
      end if;
   end if;
`
	set := fmt.Sprintf(stbSet, fun.Name(), fun.StbAddr())

	fmts.FuncsStrobesSet += set
}
