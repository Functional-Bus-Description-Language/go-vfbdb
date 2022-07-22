package python

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/access"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/elem"

	"strings"
)

func genAccess(acs access.Access, b *strings.Builder) {
	switch acs.(type) {
	case access.SingleSingle:
		a := acs.(access.SingleSingle)
		b.WriteString(
			fmt.Sprintf(
				"{'Type': 'SingleSingle', 'Width': %d, 'Addr': %d, 'Shift': %d},",
				a.Width(), a.Addr, a.Mask.Lower,
			),
		)
	case access.SingleContinuous:
		a := acs.(access.SingleContinuous)
		b.WriteString(
			fmt.Sprintf(
				"{'Type': 'SingleContinuous', 'Width': %d, 'StartAddr': %d, 'RegCount': %d, 'StartShift': %d},",
				a.Width(), a.RegCount(), a.StartAddr(), a.StartMask.Lower,
			),
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

func genParamList(params []elem.Param) string {
	if len(params) == 0 {
		return "None"
	}

	b := strings.Builder{}

	b.WriteString("[\n")
	increaseIndent(2)

	for _, p := range params {
		b.WriteString(fmt.Sprintf("%s{'Name': '%s', 'Access': ", indent, p.Name()))
		genAccess(p.Access(), &b)
		b.WriteString("},\n")
	}

	decreaseIndent(1)
	b.WriteString(fmt.Sprintf("%s]", indent))
	decreaseIndent(1)

	return b.String()
}
