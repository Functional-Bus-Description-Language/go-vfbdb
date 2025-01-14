package python

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/access"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/fn"
)

func genStatic(st *fn.Static, blk *fn.Block) string {
	if st.IsArray {
		panic("unimplemented")
	} else {
		return genStaticSingle(st, blk)
	}
}

func genStaticSingle(st *fn.Static, blk *fn.Block) string {
	var code string

	switch acs := st.Access.(type) {
	case access.SingleOneReg:
		code += indent + fmt.Sprintf(
			"self.%s = StaticSingleOneReg(iface, %d, %d, %d, 0b0%s)\n",
			st.Name,
			blk.StartAddr()+acs.StartAddr(),
			acs.StartBit(),
			acs.EndBit(),
			st.InitValue.ToBin().ValueLiteral(),
		)
	case access.SingleNRegs:
		code += indent + fmt.Sprintf(
			"self.%s = StaticSingleNRegs(iface, %d, %d, (%d, %d), (%d, %d), 0b%s)\n",
			st.Name,
			blk.StartAddr()+acs.StartAddr(),
			acs.RegCount(),
			busWidth-1, acs.StartBit(),
			acs.EndBit(), 0,
			st.InitValue.ToBin().ValueLiteral(),
		)
	default:
		panic("unimplemented")
	}

	return code
}
