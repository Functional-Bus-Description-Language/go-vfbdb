package python

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/access"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/elem"
)

func genStatus(st elem.Status, blk elem.Block) string {
	if st.IsArray() {
		return genStatusArray(st, blk)
	} else {
		return genStatusSingle(st, blk)
	}
}

func genStatusSingle(st elem.Status, blk elem.Block) string {
	var code string

	switch st.Access().(type) {
	case access.SingleSingle:
		access := st.Access().(access.SingleSingle)
		code += indent + fmt.Sprintf(
			"self.%s = StatusSingleSingle(iface, %d, (%d, %d))\n",
			st.Name(), blk.AddrSpace().Start()+access.Addr, access.Mask.Upper, access.Mask.Lower,
		)
	case access.SingleContinuous:
		a := st.Access().(access.SingleContinuous)
		decreasigOrder := "False"
		if st.HasDecreasingAccessOrder() {
			decreasigOrder = "True"
		}
		code += indent + fmt.Sprintf(
			"self.%s = StatusSingleContinuous(iface, %d, %d, (%d, %d), (%d, %d), %s)\n",
			st.Name(),
			blk.AddrSpace().Start()+a.StartAddr(),
			a.RegCount(),
			a.StartMask.Upper, a.StartMask.Lower,
			a.EndMask.Upper, a.EndMask.Lower,
			decreasigOrder,
		)
	default:
		panic("not yet implemented")
	}

	return code
}

func genStatusArray(st elem.Status, blk elem.Block) string {
	var code string

	switch st.Access().(type) {
	case access.ArraySingle:
		access := st.Access().(access.ArraySingle)
		code += indent + fmt.Sprintf(
			"self.%s = StatusArraySingle(iface, %d, (%d, %d), %d)\n",
			st.Name(),
			blk.AddrSpace().Start()+access.StartAddr(),
			access.Mask.Upper,
			access.Mask.Lower,
			access.RegCount(),
		)
	case access.ArrayMultiple:
		access := st.Access().(access.ArrayMultiple)
		code += indent + fmt.Sprintf(
			"self.%s = StatusArrayMultiple(iface, %d, %d, %d, %d, %d)\n",
			st.Name(),
			blk.AddrSpace().Start()+access.StartAddr(),
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
