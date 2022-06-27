package python

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/access"
)

func genParamAccessList(params []*fbdl.Param) string {
	if len(params) == 0 {
		return "None,"
	}

	code := "[\n"
	increaseIndent(2)

	for _, p := range params {
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
			panic("should never happen")
		}
	}

	decreaseIndent(1)
	code = code[:len(code)]
	code += indent + "],"
	decreaseIndent(1)

	return code
}
