package python

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/access"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/fn"
)

func genMask(mask *fn.Mask, blk *fn.Block) string {
	if mask.IsArray {
		return genMaskArray(mask, blk)
	} else {
		return genMaskSingle(mask, blk)
	}
}

func genMaskSingle(mask *fn.Mask, blk *fn.Block) string {
	var code string

	switch a := mask.Access.(type) {
	case access.SingleSingle:
		code += indent + fmt.Sprintf(
			"self.%s = MaskSingleSingle(iface, %d, (%d, %d))\n",
			mask.Name, blk.StartAddr()+a.Addr, a.EndBit(), a.StartBit(),
		)
	default:
		panic("not yet implemented")
	}

	return code
}

func genMaskArray(mask *fn.Mask, blk *fn.Block) string {
	panic("not yet implemented")
}
