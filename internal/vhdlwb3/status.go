package vhdlwb3

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/access"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/elem"
)

func genStatus(st elem.Status, fmts *BlockEntityFormatters, main bool) {
	if st.IsArray() {
		genStatusArray(st, fmts)
	} else {
		genStatusSingle(st, fmts, main)
	}
}

func genStatusArray(st elem.Status, fmts *BlockEntityFormatters) {
	switch st.Access().(type) {
	case access.ArraySingle:
		genStatusArraySingle(st, fmts)
	case access.ArrayMultiple:
		genStatusArrayMultiple(st, fmts)
	default:
		panic("not yet implemented")
	}
}

func genStatusSingle(st elem.Status, fmts *BlockEntityFormatters, main bool) {
	if !main || (st.Name() != "ID" && st.Name() != "TIMESTAMP") {
		s := fmt.Sprintf(";\n   %s_i : in std_logic_vector(%d downto 0)", st.Name(), st.Width()-1)
		fmts.EntityFunctionalPorts += s
	}

	switch st.Access().(type) {
	case access.SingleSingle:
		genStatusSingleSingle(st, fmts, main)
	case access.SingleContinuous:
		genStatusSingleContinuous(st, fmts)
	default:
		panic("unknown single access strategy")
	}
}

func genStatusSingleSingle(st elem.Status, fmts *BlockEntityFormatters, main bool) {
	a := st.Access().(access.SingleSingle)

	var code string
	if main && (st.Name() == "ID" || st.Name() == "TIMESTAMP") {
		code = fmt.Sprintf(
			"      master_in.dat(%d downto %d) <= %s; -- %s",
			a.EndBit(), a.StartBit(), string(st.Default()), st.Name(),
		)
	} else {
		code = fmt.Sprintf(
			"      master_in.dat(%d downto %d) <= %s_i;\n",
			a.EndBit(), a.StartBit(), st.Name(),
		)
	}

	fmts.RegistersAccess.add([2]int64{a.Addr, a.Addr}, code)
}

func genStatusSingleContinuous(st elem.Status, fmts *BlockEntityFormatters) {
	if st.Atomic() == true {
		genStatusSingleContinuousAtomic(st, fmts)
	} else {
		genStatusSingleContinuousNonAtomic(st, fmts)
	}
}

func genStatusSingleContinuousAtomic(st elem.Status, fmts *BlockEntityFormatters) {
	a := st.Access().(access.SingleContinuous)
	strategy := SeparateFirst
	atomicShadowRange := [2]int64{st.Width() - 1, a.FirstRegWidth()}
	if st.HasDecreasingAccessOrder() {
		strategy = SeparateLast
		atomicShadowRange[0] = st.Width() - 1 - a.LastRegWidth()
		atomicShadowRange[1] = 0
	}
	chunks := makeAccessChunksContinuous(a, strategy)

	fmts.SignalDeclarations += fmt.Sprintf(
		"signal %s_atomic : std_logic_vector(%d downto %d);\n",
		st.Name(), atomicShadowRange[0], atomicShadowRange[1],
	)

	for i, c := range chunks {
		var code string
		if (strategy == SeparateFirst && i == 0) || (strategy == SeparateLast && i == len(chunks)-1) {
			code = fmt.Sprintf(
				"      %[1]s_atomic(%[2]d downto %[3]d) <= %[1]s_i(%[2]d downto %[3]d);\n"+
					"      master_in.dat(%[4]d downto %[5]d) <= %[1]s_i(%[6]s downto %[7]s);",
				st.Name(), atomicShadowRange[0], atomicShadowRange[1],
				c.endBit, c.startBit, c.range_[0], c.range_[1],
			)
		} else {
			code = fmt.Sprintf(
				"      master_in.dat(%d downto %d) <= %s_atomic(%s downto %s);",
				c.endBit, c.startBit, st.Name(), c.range_[0], c.range_[1],
			)
		}

		fmts.RegistersAccess.add([2]int64{c.addr[0], c.addr[1]}, code)
	}
}

func genStatusSingleContinuousNonAtomic(st elem.Status, fmts *BlockEntityFormatters) {
	chunks := makeAccessChunksContinuous(st.Access().(access.SingleContinuous), Compact)

	for _, c := range chunks {
		code := fmt.Sprintf(
			"      master_in.dat(%d downto %d) <= %s_i(%s downto %s);",
			c.endBit, c.startBit, st.Name(), c.range_[0], c.range_[1],
		)

		fmts.RegistersAccess.add([2]int64{c.addr[0], c.addr[1]}, code)
	}
}

func genStatusArraySingle(st elem.Status, fmts *BlockEntityFormatters) {
	a := st.Access().(access.ArraySingle)

	port := fmt.Sprintf(";\n   %s_i : in slv_vector(%d downto 0)(%d downto 0)", st.Name(), st.Count()-1, st.Width()-1)
	fmts.EntityFunctionalPorts += port

	code := fmt.Sprintf(
		"      master_in.dat(%d downto %d) <= %s_i(addr - %d);",
		a.EndBit(), a.StartBit(), st.Name(), a.StartAddr(),
	)

	fmts.RegistersAccess.add(
		[2]int64{a.StartAddr(), a.StartAddr() + a.RegCount() - 1},
		code,
	)
}

func genStatusArrayMultiple(st elem.Status, fmts *BlockEntityFormatters) {
	access := st.Access().(access.ArrayMultiple)

	port := fmt.Sprintf(";\n   %s_i : in slv_vector(%d downto 0)(%d downto 0)", st.Name(), st.Count()-1, st.Width()-1)
	fmts.EntityFunctionalPorts += port

	itemsPerAccess := access.ItemsPerAccess

	var addr [2]int64
	var code string

	if access.ItemCount <= itemsPerAccess {
		addr = [2]int64{access.StartAddr(), access.EndAddr()}
		code = fmt.Sprintf(`      for i in 0 to %[1]d loop
         master_in.dat(%[2]d*(i+1)+%[3]d-1 downto %[2]d*i+%[3]d) <= %[4]s_i(i);
      end loop;`,
			st.Count()-1, access.ItemWidth, access.StartBit(), st.Name(),
		)
	} else if access.ItemsInLastReg() == 0 {
		addr = [2]int64{access.StartAddr(), access.EndAddr()}
		code = fmt.Sprintf(`      for i in 0 to %[1]d loop
         master_in.dat(%[2]d*(i+1)+%[3]d-1 downto %[2]d*i+%[3]d) <= %[4]s_i((addr-%[5]d)*%[6]d+i);
      end loop;`,
			itemsPerAccess-1, access.ItemWidth, access.StartBit(), st.Name(), access.StartAddr(), access.ItemsPerAccess,
		)
	} else {
		addr = [2]int64{access.StartAddr(), access.EndAddr() - 1}
		code = fmt.Sprintf(`      for i in 0 to %[1]d loop
         master_in.dat(%[2]d*(i+1) + %[3]d-1 downto %[2]d*i + %[3]d) <= %[4]s_i((addr-%[5]d)*%[6]d+i);
      end loop;`,
			itemsPerAccess-1, access.ItemWidth, access.StartBit(), st.Name(), access.StartAddr(), access.ItemsPerAccess,
		)
		fmts.RegistersAccess.add(addr, code)

		addr = [2]int64{access.EndAddr(), access.EndAddr()}
		code = fmt.Sprintf(`      for i in 0 to %[1]d loop
         master_in.dat(%[2]d*(i+1) + %[3]d-1 downto %[2]d*i+%[3]d) <= %[4]s_i(%[5]d+i);
      end loop;`,
			access.ItemsInLastReg()-1, access.ItemWidth, access.StartBit(), st.Name(), (access.RegCount()-1)*access.ItemsPerAccess,
		)
	}

	fmts.RegistersAccess.add(addr, code)
}
