package python

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/access"

	"strings"
)

func genAccess(acs access.Access, b *strings.Builder) {
	switch acs.(type) {
	case access.SingleSingle:
		a := acs.(access.SingleSingle)
		b.WriteString(
			fmt.Sprintf(
				"{'Type': 'SingleSingle', 'Addr': %d, 'StartBit': %d, 'EndBit': %d, 'RegCount': 1},",
				a.Addr, a.StartBit(), a.EndBit(),
			),
		)
	case access.SingleContinuous:
		a := acs.(access.SingleContinuous)
		b.WriteString(
			fmt.Sprintf(
				"{'Type': 'SingleContinuous', 'StartAddr': %d, 'StartBit': %d,: 'EndBit': %d, 'RegCount': %d},",
				a.StartAddr(), a.StartBit(), a.EndBit(), a.RegCount(),
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
