package vhdlwb3

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/access"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/fn"
)

func genStatic(st *fn.Static, fmts *BlockEntityFormatters) {
	if st.IsArray {
		panic("unimplemented")
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
	case access.SingleOneReg:
		genStaticSingleOneReg(st, fmts)
	default:
		panic("unimplemented")
	}
}

func genStaticSingleOneReg(st *fn.Static, fmts *BlockEntityFormatters) {
	acs := st.Access.(access.SingleOneReg)

	code := fmt.Sprintf(
		"      master_in.dat(%d downto %d) <= %s; -- %s\n",
		acs.EndBit(), acs.StartBit(), string(st.InitValue), st.Name,
	)
	addr := acs.StartAddr()
	fmts.RegistersAccess.add([2]int64{addr, addr}, code)
}
