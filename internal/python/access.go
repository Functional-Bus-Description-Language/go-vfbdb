package python

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/access"

	"strings"
)

func genAccess(acs access.Access, b *strings.Builder) {
	b.WriteString(
		fmt.Sprintf(
			"{'StartAddr': %d, 'StartBit': %d, 'EndBit': %d, 'RegCount': %d, ",
			acs.StartAddr(), acs.StartBit(), acs.EndBit(), acs.RegCount(),
		),
	)

	switch acs.(type) {
	case access.SingleSingle:
		b.WriteString("'Type': 'SingleSingle'")
	case access.SingleContinuous:
		b.WriteString("'Type': 'SingleContinuous'")
	case access.ArrayContinuous:
		panic("not yet implemented")
	case access.ArrayMultiple:
		panic("not yet implemented")
	case access.ArraySingle:
		panic("not yet implemented")
	default:
		panic("should never happen")
	}

	b.WriteString("},")
}
