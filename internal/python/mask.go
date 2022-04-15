package python

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl"
)

func genMask(mask *fbdl.Mask, blk *fbdl.Block) string {
	if mask.IsArray {
		return genMaskArray(mask, blk)
	} else {
		return genMaskSingle(mask, blk)
	}
}

func genMaskSingle(mask *fbdl.Mask, blk *fbdl.Block) string {
	var code string

	switch mask.Access.(type) {
	case fbdl.AccessSingleSingle:
		access := mask.Access.(fbdl.AccessSingleSingle)
		code += indent + fmt.Sprintf(
			"self.%s = MaskSingleSingle(iface, %d, (%d, %d))\n",
			mask.Name, blk.AddrSpace.Start()+access.Addr, access.Mask.Upper, access.Mask.Lower,
		)
	default:
		panic("not yet implemented")
	}

	return code
}

func genMaskArray(mask *fbdl.Mask, blk *fbdl.Block) string {
	panic("not yet implemented")
}
