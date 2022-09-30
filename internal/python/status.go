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
		a := st.Access().(access.SingleSingle)
		code += indent + fmt.Sprintf(
			"self.%s = StatusSingleSingle(iface, %d, (%d, %d))\n",
			st.Name(), blk.AddrSpace().Start()+a.Addr, a.EndBit(), a.StartBit(),
		)
	case access.SingleContinuous:
		a := st.Access().(access.SingleContinuous)

		code += indent + fmt.Sprintf(
			"self.%s = StatusSingleContinuous(iface, %d, %d, (%d, %d), (%d, %d))\n",
			st.Name(),
			blk.AddrSpace().Start()+a.StartAddr(),
			a.RegCount(),
			busWidth-1, a.StartBit(),
			a.EndBit(), 0,
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
		a := st.Access().(access.ArraySingle)
		code += indent + fmt.Sprintf(
			"self.%s = StatusArraySingle(iface, %d, (%d, %d), %d)\n",
			st.Name(),
			blk.AddrSpace().Start()+a.StartAddr(),
			a.EndBit(),
			a.StartBit(),
			a.RegCount(),
		)
	case access.ArrayMultiple:
		access := st.Access().(access.ArrayMultiple)
		code += indent + fmt.Sprintf(
			"self.%s = StatusArrayMultiple(iface, %d, %d, %d, %d, %d)\n",
			st.Name(),
			blk.AddrSpace().Start()+access.StartAddr(),
			access.StartBit(),
			access.ItemWidth,
			access.ItemCount,
			access.ItemsPerAccess,
		)
	default:
		panic("not yet implemented")
	}

	return code
}
