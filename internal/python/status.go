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
			"self.%s = StatusSingleOneReg(iface, %d, %d, %d)\n",
			st.Name,
			blk.StartAddr()+acs.StartAddr(),
			acs.StartBit(),
			acs.EndBit(),
		)
	case access.SingleNRegs:
		code += indent + fmt.Sprintf(
			"self.%s = StatusSingleNRegs(iface, %d, %d, (%d, %d), (%d, %d))\n",
			st.Name,
			blk.StartAddr()+acs.StartAddr(),
			acs.RegCount(),
			busWidth-1, acs.StartBit(),
			acs.EndBit(), 0,
		)
	default:
		panic("unimplemented")
	}

	return code
}

func genStatusArray(st *fn.Status, blk *fn.Block) string {
	var code string

	switch acs := st.Access.(type) {
	case access.ArrayOneInReg:
		code += indent + fmt.Sprintf(
			"self.%s = StatusArrayOneInReg(iface, %d, (%d, %d), %d)\n",
			st.Name,
			blk.StartAddr()+acs.StartAddr(),
			acs.EndBit(),
			acs.StartBit(),
			acs.RegCount(),
		)
	case access.ArrayOneReg:
		code += indent + fmt.Sprintf(
			"self.%s = StatusArrayOneReg(iface, %d, %d, %d, %d)\n",
			st.Name,
			blk.StartAddr()+acs.StartAddr(),
			acs.StartBit(),
			acs.ItemWidth(),
			acs.ItemCount(),
		)
	case access.ArrayNInReg:
		code += indent + fmt.Sprintf(
			"self.%s = StatusArrayNInReg(iface, %d, %d, %d, %d, %d)\n",
			st.Name,
			blk.StartAddr()+acs.StartAddr(),
			acs.StartBit(),
			acs.ItemWidth(),
			acs.ItemCount(),
			acs.ItemsInReg(),
		)
	case access.ArrayNInRegMInEndReg:
		code += indent + fmt.Sprintf(
			"self.%s = StatusArrayNInRegMInEndReg(iface, %d, %d, %d, %d, %d)\n",
			st.Name,
			blk.StartAddr()+acs.StartAddr(),
			acs.StartBit(),
			acs.ItemWidth(),
			acs.ItemCount(),
			acs.ItemsInReg(),
		)
	case access.ArrayOneInNRegs:
		code += indent + fmt.Sprintf(
			"self.%s = StatusArrayOneInNRegs(iface, %d, %d, %d, %d, %d, %d)\n",
			st.Name,
			blk.StartAddr()+acs.StartAddr(),
			acs.ItemWidth(),
			acs.ItemCount(),
			acs.RegsPerItem(),
			acs.RegCount(),
			acs.EndBit(),
		)
	default:
		panic("unimplemented")
	}

	return code
}
