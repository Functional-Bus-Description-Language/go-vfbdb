package python

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/access"

	"strings"
)

func genAccess(acs access.Access, blkAddr int64, b *strings.Builder) {
	b.WriteString(
		fmt.Sprintf(
			"{'StartAddr': %d, 'StartBit': %d, 'EndBit': %d, 'RegCount': %d, 'Type': ",
			blkAddr+acs.StartAddr(), acs.StartBit(), acs.EndBit(), acs.RegCount(),
		),
	)

	switch a := acs.(type) {
	case access.SingleOneReg:
		b.WriteString("'SingleOneReg'")
	case access.SingleNRegs:
		b.WriteString("'SingleNRegs'")
	case access.ArrayNRegs:
		b.WriteString(fmt.Sprintf("'ArrayNRegs', 'ItemCount': %d", a.ItemCount()))
	case access.ArrayOneInReg:
		panic("unimplemented")
	case access.ArrayNInReg:
		panic("unimplemented")
	case access.ArrayNInRegMInEndReg:
		panic("unimplemented")
	default:
		panic("should never happen")
	}

	b.WriteString("},")
}
