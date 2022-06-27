package python

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/access"
)

func genFunc(fun *fbdl.Func, blk *fbdl.Block) string {
	code := indent + fmt.Sprintf("self.%s = Func(iface, %d, ",
		fun.Name, blk.AddrSpace.Start()+fun.ParamsStartAddr(),
	)
	code += genFuncParamAccessList(fun)
	code += genFuncReturnAccessList(fun)
	code += ")\n"

	return code
}

func genFuncParamAccessList(fun *fbdl.Func) string {
	if len(fun.Params) == 0 {
		return "None,"
	}

	code := "[\n"
	increaseIndent(2)

	for _, p := range fun.Params {
		switch p.Access.(type) {
		case access.SingleSingle:
			ass := p.Access.(access.SingleSingle)
			code += indent + fmt.Sprintf(
				"{'Type': 'SingleSingle', 'Width': %d, 'Addr': %d, 'Shift': %d},\n",
				p.Width, ass.Addr, ass.Mask.Lower,
			)
		case access.SingleContinuous:
			asc := p.Access.(access.SingleContinuous)
			code += indent + fmt.Sprintf(
				"{'Type': 'SingleContinuous', 'Width': %d, 'StartAddr': %d, 'RegCount': %d, 'StartShift': %d},\n",
				p.Width, asc.RegCount(), asc.StartAddr(), asc.StartMask.Lower,
			)
		default:
			//panic("should never happen")
		}
	}

	decreaseIndent(1)
	code = code[:len(code)]
	code += indent + "],"
	decreaseIndent(1)

	return code
}

func genFuncReturnAccessList(fun *fbdl.Func) string {
	return "None"
}
