package csync

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/access"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/elem"
	"github.com/Functional-Bus-Description-Language/go-vfbdb/internal/c"
	_ "github.com/Functional-Bus-Description-Language/go-vfbdb/internal/utils"
	"strings"
)

func genFunc(fun elem.Func, hFmts *BlockHFormatters, cFmts *BlockCFormatters) {
	sig := genFuncSignature(fun, hFmts)

	hFmts.Code += "\n" + sig + ";\n"

	cFmts.Code += fmt.Sprintf("\n%s {\n", sig)
	if len(fun.Params()) == 0 && len(fun.Returns()) == 0 {
		cFmts.Code += fmt.Sprintf("\treturn iface->write(%d, 0);\n};\n", fun.StbAddr())
		return
	}

	if len(fun.Params()) > 0 {
		genFuncParamsAccess(fun, cFmts)
	}

	if len(fun.Returns()) > 0 {
		genFuncReturnsAccess(fun, cFmts)
	}

	cFmts.Code += "};\n"
}

func genFuncSignature(fun elem.Func, hFmts *BlockHFormatters) string {
	prefix := "int vfbdb_" + hFmts.BlockName + "_" + fun.Name()

	params := strings.Builder{}
	params.WriteString("const vfbdb_iface_t * const iface")

	for _, p := range fun.Params() {
		params.WriteString(
			", const " + c.WidthToWriteType(p.Width()).String() + " " + p.Name(),
		)
	}

	for _, r := range fun.Returns() {
		params.WriteString(
			", " + c.WidthToReadType(r.Width()).String() + " const " + r.Name(),
		)
	}

	return prefix + "(" + params.String() + ")"
}

func genFuncParamsAccess(fun elem.Func, cFmts *BlockCFormatters) {
	if fun.ParamsBufSize() == 1 {
		genFuncParamsAccessSingleReg(fun, cFmts)
	}
}

func genFuncParamsAccessSingleReg(fun elem.Func, cFmts *BlockCFormatters) {
	cFmts.Code += fmt.Sprintf("\treturn iface->write(%d, ", fun.StbAddr())
	for i, p := range fun.Params() {
		if i != 0 {
			cFmts.Code += " | "
		}

		switch a := p.Access().(type) {
		case access.SingleSingle:
			cFmts.Code += fmt.Sprintf("%s << %d", p.Name(), a.StartBit())
		default:
			panic("not yet implemented")
		}
	}
	cFmts.Code += ");\n"
}

func genFuncReturnsAccess(fun elem.Func, cFmts *BlockCFormatters) {

}
