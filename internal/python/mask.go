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

	switch acs := mask.Access.(type) {
	case access.SingleOneReg:
		code += indent + fmt.Sprintf(
			"self.%s = MaskSingleOneReg(iface, %d, %d, %d)\n",
			mask.Name,
			blk.StartAddr()+acs.Addr,
			acs.StartBit,
			acs.EndBit,
		)
	default:
		panic("unimplemented")
	}

	return code
}

func genMaskArray(mask *fn.Mask, blk *fn.Block) string {
	panic("unimplemented")
}
