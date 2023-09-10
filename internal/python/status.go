package python

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/access"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/fn"
)

func genStatus(st *fn.Status, blk *fn.Block) string {
	if st.IsArray {
		return genStatusArray(st, blk)
	} else {
		return genStatusSingle(st, blk)
	}
}

func genStatusSingle(st *fn.Status, blk *fn.Block) string {
	var code string

	switch acs := st.Access.(type) {
	case access.SingleOneReg:
		code += indent + fmt.Sprintf(
			"self.%s = StatusSingleOneReg(iface, %d, (%d, %d))\n",
			st.Name, blk.StartAddr()+acs.Addr, acs.EndBit, acs.StartBit,
		)
	case access.SingleNRegs:
		code += indent + fmt.Sprintf(
			"self.%s = StatusSingleNRegs(iface, %d, %d, (%d, %d), (%d, %d))\n",
			st.Name,
			blk.StartAddr()+acs.GetStartAddr(),
			acs.GetRegCount(),
			busWidth-1, acs.GetStartBit(),
			acs.GetEndBit(), 0,
		)
	default:
		panic("unimplemented")
	}

	return code
}

func genStatusArray(st *fn.Status, blk *fn.Block) string {
	var code string

	switch a := st.Access.(type) {
	case access.ArraySingle:
		code += indent + fmt.Sprintf(
			"self.%s = StatusArraySingle(iface, %d, (%d, %d), %d)\n",
			st.Name,
			blk.StartAddr()+a.GetStartAddr(),
			a.GetEndBit(),
			a.GetStartBit(),
			a.GetRegCount(),
		)
	case access.ArrayOneReg:
		code += indent + fmt.Sprintf(
			"self.%s = StatusArrayOneReg(iface, %d, %d, %d, %d)\n",
			st.Name,
			blk.StartAddr()+a.GetStartAddr(),
			a.GetStartBit(),
			a.ItemWidth,
			a.ItemCount,
		)
	case access.ArrayMultiple:
		code += indent + fmt.Sprintf(
			"self.%s = StatusArrayMultiple(iface, %d, %d, %d, %d, %d)\n",
			st.Name,
			blk.StartAddr()+a.GetStartAddr(),
			a.GetStartBit(),
			a.ItemWidth,
			a.ItemCount,
			a.ItemsPerReg,
		)
	default:
		panic("unimplemented")
	}

	return code
}
