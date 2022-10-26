package python

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/access"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/elem"
)

func genStatic(st elem.Static, blk elem.Block) string {
	if st.IsArray() {
		panic("not yet implemented")
	} else {
		return genStaticSingle(st, blk)
	}
}

func genStaticSingle(st elem.Static, blk elem.Block) string {
	var code string

	switch st.Access().(type) {
	case access.SingleSingle:
		a := st.Access().(access.SingleSingle)
		code += indent + fmt.Sprintf(
			"self.%s = StaticSingleSingle(iface, %d, (%d, %d), 0b0%s)\n",
			st.Name(), blk.AddrSpace().Start()+a.Addr, a.EndBit(), a.StartBit(),
			st.Default().ToBin().ValueLiteral(),
		)
	case access.SingleContinuous:
		a := st.Access().(access.SingleContinuous)

		code += indent + fmt.Sprintf(
			"self.%s = StaticSingleContinuous(iface, %d, %d, (%d, %d), (%d, %d), 0b%s)\n",
			st.Name(),
			blk.AddrSpace().Start()+a.StartAddr(),
			a.RegCount(),
			busWidth-1, a.StartBit(),
			a.EndBit(), 0,
			st.Default().ToBin().ValueLiteral(),
		)
	default:
		panic("not yet implemented")
	}

	return code
}
