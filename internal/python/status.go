package python

import (
	"fmt"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl"
)

func generateStatus(st *fbdl.Status, blk *fbdl.Block) string {
	if st.IsArray {
		return generateStatusArray(st, blk)
	} else {
		return generateStatusSingle(st, blk)
	}
}

func generateStatusSingle(st *fbdl.Status, blk *fbdl.Block) string {
	var code string

	switch st.Access.(type) {
	case fbdl.AccessSingleSingle:
		access := st.Access.(fbdl.AccessSingleSingle)
		code += indent + fmt.Sprintf(
			"self.%s = StatusSingleSingle(interface, %d, (%d, %d))\n",
			st.Name, blk.AddrSpace.Start()+access.Addr, access.Mask.Upper, access.Mask.Lower,
		)
	default:
		panic("not yet implemented")
	}

	return code
}

func generateStatusArray(st *fbdl.Status, blk *fbdl.Block) string {
	var code string

	switch st.Access.(type) {
	case fbdl.AccessArraySingle:
		access := st.Access.(fbdl.AccessArraySingle)
		code += indent + fmt.Sprintf(
			"self.%s = StatusArraySingle(interface, %d, (%d, %d), %d)\n",
			st.Name,
			blk.AddrSpace.Start()+access.StartAddr(),
			access.Mask.Upper,
			access.Mask.Lower,
			access.RegCount(),
		)
	case fbdl.AccessArrayMultiple:
		access := st.Access.(fbdl.AccessArrayMultiple)
		code += indent + fmt.Sprintf(
			"self.%s = StatusArrayMultiple(interface, %d, %d, %d, %d, %d)\n",
			st.Name,
			blk.AddrSpace.Start()+access.StartAddr(),
			access.StartBit,
			access.ItemWidth,
			access.ItemCount,
			access.ItemsPerAccess,
		)
	default:
		panic("not yet implemented")
	}

	return code
}
