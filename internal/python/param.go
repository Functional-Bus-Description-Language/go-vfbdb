package python

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/access"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/elem"
)

func genAccessList(accesses []access.Access) string {
	if len(accesses) == 0 {
		return "None,"
	}

	code := "[\n"
	increaseIndent(2)

	for _, a := range accesses {
		switch a.(type) {
		case access.SingleSingle:
			ass := a.(access.SingleSingle)
			code += indent + fmt.Sprintf(
				"{'Type': 'SingleSingle', 'Width': %d, 'Addr': %d, 'Shift': %d},\n",
				a.Width(), ass.Addr, ass.Mask.Lower,
			)
		case access.SingleContinuous:
			asc := a.(access.SingleContinuous)
			code += indent + fmt.Sprintf(
				"{'Type': 'SingleContinuous', 'Width': %d, 'StartAddr': %d, 'RegCount': %d, 'StartShift': %d},\n",
				a.Width(), asc.RegCount(), asc.StartAddr(), asc.StartMask.Lower,
			)
		case access.ArrayContinuous:
			panic("not yet implemented")
		case access.ArrayMultiple:
			panic("not yet implemented")
		case access.ArraySingle:
			panic("not yet implemented")
		default:
			panic("should never happen")
		}
	}

	decreaseIndent(1)
	code = code[:len(code)]
	code += indent + "],"
	decreaseIndent(1)

	return code
}

func genParamAccessList(params []elem.Param) string {
	accesses := []access.Access{}
	for _, p := range params {
		accesses = append(accesses, p.Access())
	}

	return genAccessList(accesses)
}
