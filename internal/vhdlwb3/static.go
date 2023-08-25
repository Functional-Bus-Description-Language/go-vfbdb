package vhdlwb3

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/access"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/fn"
)

func genStatic(st *fn.Static, fmts *BlockEntityFormatters) {
	if st.IsArray {
		panic("not implemented")
	} else {
		genStaticSingle(st, fmts)
	}
}

func genStaticSingle(st *fn.Static, fmts *BlockEntityFormatters) {
	s := fmt.Sprintf(
		";\n   %s_o : out std_logic_vector(%d downto 0) := %s",
		st.Name, st.Width-1, string(st.InitValue),
	)
	fmts.EntityFunctionalPorts += s

	switch st.Access.(type) {
	case access.SingleSingle:
		genStaticSingleSingle(st, fmts)
	case access.SingleContinuous:
		panic("not implemented")
	default:
		panic("unknown single access strategy")
	}
}

func genStaticSingleSingle(st *fn.Static, fmts *BlockEntityFormatters) {
	a := st.Access.(access.SingleSingle)

	code := fmt.Sprintf(
		"      master_in.dat(%d downto %d) <= %s; -- %s",
		a.EndBit(), a.StartBit(), string(st.InitValue), st.Name,
	)

	fmts.RegistersAccess.add([2]int64{a.Addr, a.Addr}, code)
}
