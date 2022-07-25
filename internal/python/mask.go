package python

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/access"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/elem"
)

func genMask(mask elem.Mask, blk elem.Block) string {
	if mask.IsArray() {
		return genMaskArray(mask, blk)
	} else {
		return genMaskSingle(mask, blk)
	}
}

func genMaskSingle(mask elem.Mask, blk elem.Block) string {
	var code string

	switch mask.Access().(type) {
	case access.SingleSingle:
		a := mask.Access().(access.SingleSingle)
		code += indent + fmt.Sprintf(
			"self.%s = MaskSingleSingle(iface, %d, (%d, %d))\n",
			mask.Name(), blk.AddrSpace().Start()+a.Addr, a.EndBit(), a.StartBit(),
		)
	default:
		panic("not yet implemented")
	}

	return code
}

func genMaskArray(mask elem.Mask, blk elem.Block) string {
	panic("not yet implemented")
}
