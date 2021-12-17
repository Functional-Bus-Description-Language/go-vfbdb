package python

import (
	"fmt"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl"
)

func generateMask(mask *fbdl.Mask, blk *fbdl.Block) string {
	if mask.IsArray {
		return generateMaskArray(mask, blk)
	} else {
		return generateMaskSingle(mask, blk)
	}
}

func generateMaskSingle(mask *fbdl.Mask, blk *fbdl.Block) string {
	var code string

	switch mask.Access.(type) {
	case fbdl.AccessSingleSingle:
		access := mask.Access.(fbdl.AccessSingleSingle)
		code += indent + fmt.Sprintf(
			"self.%s = MaskSingleSingle(interface, %d, (%d, %d))\n",
			mask.Name, blk.AddrSpace.Start()+access.Addr, access.Mask.Upper, access.Mask.Lower,
		)
	default:
		panic("not yet implemented")
	}

	return code
}

func generateMaskArray(mask *fbdl.Mask, blk *fbdl.Block) string {
	panic("not yet implemented")
}
