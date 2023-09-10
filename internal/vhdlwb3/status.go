package vhdlwb3

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/access"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/fn"
)

func genStatus(st *fn.Status, fmts *BlockEntityFormatters) {
	if st.IsArray {
		genStatusArray(st, fmts)
	} else {
		genStatusSingle(st, fmts)
	}
}

func genStatusArray(st *fn.Status, fmts *BlockEntityFormatters) {
	switch st.Access.(type) {
	case access.ArrayOneReg:
		genStatusArrayOneReg(st, fmts)
	case access.ArrayOneInReg:
		genStatusArrayOneInReg(st, fmts)
	case access.ArrayNInReg:
		genStatusArrayNInReg(st, fmts)
	default:
		panic("unimplemented")
	}
}

func genStatusSingle(st *fn.Status, fmts *BlockEntityFormatters) {
	s := fmt.Sprintf(";\n   %s_i : in std_logic_vector(%d downto 0)", st.Name, st.Width-1)
	fmts.EntityFunctionalPorts += s

	switch st.Access.(type) {
	case access.SingleOneReg:
		genStatusSingleOneReg(st, fmts)
	case access.SingleNRegs:
		genStatusSingleNRegs(st, fmts)
	default:
		panic("unknown single access strategy")
	}
}

func genStatusSingleOneReg(st *fn.Status, fmts *BlockEntityFormatters) {
	acs := st.Access.(access.SingleOneReg)

	code := fmt.Sprintf(
		"      master_in.dat(%d downto %d) <= %s_i;\n",
		acs.EndBit, acs.StartBit, st.Name,
	)

	fmts.RegistersAccess.add([2]int64{acs.Addr, acs.Addr}, code)
}

func genStatusSingleNRegs(st *fn.Status, fmts *BlockEntityFormatters) {
	if st.Atomic {
		genStatusSingleNRegsAtomic(st, fmts)
	} else {
		genStatusSingleNRegsNonAtomic(st, fmts)
	}
}

func genStatusSingleNRegsAtomic(st *fn.Status, fmts *BlockEntityFormatters) {
	acs := st.Access.(access.SingleNRegs)
	strategy := SeparateFirst
	atomicShadowRange := [2]int64{st.Width - 1, acs.GetStartRegWidth()}
	chunks := makeAccessChunksContinuous(acs, strategy)

	fmts.SignalDeclarations += fmt.Sprintf(
		"signal %s_atomic : std_logic_vector(%d downto %d);\n",
		st.Name, atomicShadowRange[0], atomicShadowRange[1],
	)

	for i, c := range chunks {
		var code string
		if (strategy == SeparateFirst && i == 0) || (strategy == SeparateLast && i == len(chunks)-1) {
			code = fmt.Sprintf(`
      %[1]s_atomic(%[2]d downto %[3]d) <= %[1]s_i(%[2]d downto %[3]d);
      master_in.dat(%[4]d downto %[5]d) <= %[1]s_i(%[6]s downto %[7]s);`,
				st.Name, atomicShadowRange[0], atomicShadowRange[1],
				c.endBit, c.startBit, c.range_[0], c.range_[1],
			)
		} else {
			code = fmt.Sprintf(
				"      master_in.dat(%d downto %d) <= %s_atomic(%s downto %s);",
				c.endBit, c.startBit, st.Name, c.range_[0], c.range_[1],
			)
		}

		fmts.RegistersAccess.add([2]int64{c.addr[0], c.addr[1]}, code)
	}
}

func genStatusSingleNRegsNonAtomic(st *fn.Status, fmts *BlockEntityFormatters) {
	chunks := makeAccessChunksContinuous(st.Access.(access.SingleNRegs), Compact)

	for _, c := range chunks {
		code := fmt.Sprintf(
			"      master_in.dat(%d downto %d) <= %s_i(%s downto %s);",
			c.endBit, c.startBit, st.Name, c.range_[0], c.range_[1],
		)

		fmts.RegistersAccess.add([2]int64{c.addr[0], c.addr[1]}, code)
	}
}

func genStatusArrayOneInReg(st *fn.Status, fmts *BlockEntityFormatters) {
	acs := st.Access.(access.ArrayOneInReg)

	port := fmt.Sprintf(";\n   %s_i : in slv_vector(%d downto 0)(%d downto 0)", st.Name, st.Count-1, st.Width-1)
	fmts.EntityFunctionalPorts += port

	code := fmt.Sprintf(
		"      master_in.dat(%d downto %d) <= %s_i(addr - %d);",
		acs.EndBit, acs.StartBit, st.Name, acs.StartAddr,
	)

	fmts.RegistersAccess.add(
		[2]int64{acs.StartAddr, acs.StartAddr + acs.RegCount - 1},
		code,
	)
}

func genStatusArrayOneReg(st *fn.Status, fmts *BlockEntityFormatters) {
	acs := st.Access.(access.ArrayOneReg)

	port := fmt.Sprintf(
		";\n   %s_i : in slv_vector(%d downto 0)(%d downto 0)",
		st.Name, st.Count-1, st.Width-1,
	)
	fmts.EntityFunctionalPorts += port

	addr := [2]int64{acs.GetStartAddr(), acs.GetEndAddr()}
	code := fmt.Sprintf(`
      for i in 0 to %[1]d loop
         master_in.dat(%[2]d*(i+1)+%[3]d-1 downto %[2]d*i+%[3]d) <= %[4]s_i(i);
      end loop;`,
		st.Count-1, acs.ItemWidth, acs.StartBit, st.Name,
	)

	fmts.RegistersAccess.add(addr, code)
}

func genStatusArrayNInReg(st *fn.Status, fmts *BlockEntityFormatters) {
	acs := st.Access.(access.ArrayNInReg)

	port := fmt.Sprintf(
		";\n   %s_i : in slv_vector(%d downto 0)(%d downto 0)",
		st.Name, st.Count-1, st.Width-1,
	)
	fmts.EntityFunctionalPorts += port

	addr := [2]int64{acs.StartAddr, acs.GetEndAddr()}
	code := fmt.Sprintf(`
      for i in 0 to %[1]d loop
         master_in.dat(%[2]d*(i+1)+%[3]d-1 downto %[2]d*i+%[3]d) <= %[4]s_i((addr-%[5]d)*%[6]d+i);
      end loop;`,
		acs.ItemsInReg-1, acs.ItemWidth, acs.StartBit, st.Name, acs.StartAddr, acs.ItemsInReg,
	)

	fmts.RegistersAccess.add(addr, code)
}

/*
func genStatusArrayMultiple(st *fn.Status, fmts *BlockEntityFormatters) {
	a := st.Access.(access.ArrayMultiple)

	port := fmt.Sprintf(
		";\n   %s_i : in slv_vector(%d downto 0)(%d downto 0)",
		st.Name, st.Count-1, st.Width-1,
	)
	fmts.EntityFunctionalPorts += port

	var addr [2]int64
	var code string

	addr = [2]int64{a.GetStartAddr(), a.GetEndAddr() - 1}
	code = fmt.Sprintf(`
      for i in 0 to %[1]d loop
         master_in.dat(%[2]d*(i+1) + %[3]d-1 downto %[2]d*i + %[3]d) <= %[4]s_i((addr-%[5]d)*%[6]d+i);
      end loop;`,
		a.ItemsPerReg-1, a.ItemWidth, a.GetStartBit(), st.Name, a.GetStartAddr(), a.ItemsPerReg,
	)
	fmts.RegistersAccess.add(addr, code)

	addr = [2]int64{a.GetEndAddr(), a.GetEndAddr()}
	code = fmt.Sprintf(`
      for i in 0 to %[1]d loop
         master_in.dat(%[2]d*(i+1) + %[3]d-1 downto %[2]d*i+%[3]d) <= %[4]s_i(%[5]d+i);
      end loop;`,
		a.ItemsInLastReg()-1, a.ItemWidth, a.GetStartBit(), st.Name, (a.GetRegCount()-1)*a.ItemsPerReg,
	)

	fmts.RegistersAccess.add(addr, code)
}
*/
