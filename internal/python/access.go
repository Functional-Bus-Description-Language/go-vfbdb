package python

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/access"

	"strings"
)

func genAccess(acs access.Access, b *strings.Builder) {
	b.WriteString(
		fmt.Sprintf(
			"{'StartAddr': %d, 'StartBit': %d, 'EndBit': %d, 'RegCount': %d, 'Type': ",
			acs.GetStartAddr(), acs.GetStartBit(), acs.GetEndBit(), acs.GetRegCount(),
		),
	)

	switch a := acs.(type) {
	case access.SingleOneReg:
		b.WriteString("'SingleSingle'")
	case access.SingleContinuous:
		b.WriteString("'SingleContinuous'")
	case access.ArrayContinuous:
		b.WriteString(fmt.Sprintf("'ArrayContinuous', 'ItemCount': %d", a.ItemCount))
	case access.ArrayMultiple:
		panic("not yet implemented")
	case access.ArraySingle:
		panic("not yet implemented")
	default:
		panic("should never happen")
	}

	b.WriteString("},")
}
