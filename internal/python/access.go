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
				"{'Type': 'SingleSingle', 'Addr': %d, 'Mask': (%d, %d), 'RegCount': 1},",
				a.Addr, a.Mask.Upper, a.Mask.Lower,
			),
		)
	case access.SingleContinuous:
		a := acs.(access.SingleContinuous)
		b.WriteString(
			fmt.Sprintf(
				"{'Type': 'SingleContinuous', 'StartAddr': %d, 'StartMask': (%d, %d), 'RegCount': %d},",
				a.StartAddr(), a.StartMask.Upper, a.StartMask.Lower, a.RegCount(),
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
