package python

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl"
)

func generateFunc(fun *fbdl.Func, blk *fbdl.Block) string {
	/*
		code := generateFuncFunctionSignature(fun)
	*/

	increaseIndent(1)
	code := indent + fmt.Sprintf("self.%s = Func(iface, %d, ",
		fun.Name, blk.AddrSpace.Start()+fun.ParamsStartAddr(),
	)
	code += generateFuncParamAccessList(fun)
	code += generateFuncReturnAccessList(fun)
	code += ")\n"

	/*
		if fun.AreAllParamsSingleSingle() {
			code = generateFuncSingleSingle(fun, blk)
		} else {
			panic("not yet implemented")
		}
	*/

	return code
}

func generateFuncFunctionSignature(fun *fbdl.Func) string {
	code := indent + fmt.Sprintf("def %s(self, ", fun.Name)
	for _, p := range fun.Params {
		code += p.Name + ", "
	}
	code = code[:len(code)-2]
	code += "):\n"

	increaseIndent(1)

	return code
}

func generateFuncParamAccessList(fun *fbdl.Func) string {
	if len(fun.Params) == 0 {
		return "None,"
	}

	code := "[\n"
	increaseIndent(2)

	for _, p := range fun.Params {
		switch p.Access.(type) {
		case fbdl.AccessSingleSingle:
			ass := p.Access.(fbdl.AccessSingleSingle)
			code += indent + fmt.Sprintf(
				"{'Type': 'SingleSingle', 'Width': %d, 'Addr': %d, 'Shift': %d},\n",
				p.Width, ass.Addr, ass.Mask.Lower,
			)
		case fbdl.AccessSingleContinuous:
			asc := p.Access.(fbdl.AccessSingleContinuous)
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

func generateFuncReturnAccessList(fun *fbdl.Func) string {
	return "None"
}

// generateFuncSingleSingle generates function body for func which all parameters are of type AccessSingleSingle.
// In such case there is no Python for loop as it is relatiely easy to unroll during code generation.
func generateFuncSingleSingle(fun *fbdl.Func, blk *fbdl.Block) string {
	code := generateFuncFunctionSignature(fun)

	val := ""
	for i, p := range fun.Params {
		access := p.Access.(fbdl.AccessSingleSingle)
		val += fmt.Sprintf("%s << %d | ", p.Name, access.Mask.Lower)
		if i == len(fun.Params)-1 || fun.Params[i+1].Access.StartAddr() != access.Addr {
			val = val[:len(val)-3]
			code += indent + fmt.Sprintf("self.iface.write(%d, %s)\n", blk.AddrSpace.Start()+access.Addr, val)
			val = ""
		}
	}

	if len(fun.Params) == 0 {
		code += indent + fmt.Sprintf("self.iface.write(%d, 0)\n", blk.AddrSpace.Start()+fun.StbAddr)
	}

	decreaseIndent(1)

	return code
}
