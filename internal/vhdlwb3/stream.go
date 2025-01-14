package vhdlwb3

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/access"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/fn"
)

func genStream(stream *fn.Stream, fmts *BlockEntityFormatters) {
	genStreamType(stream, fmts)
	genStreamPorts(stream, fmts)

	if stream.IsUpstream() {
		genUpstreamAccess(stream, fmts)
	} else {
		genDownstreamAccess(stream, fmts)
	}

	genStreamStrobe(stream, fmts)
}

func genStreamType(stream *fn.Stream, fmts *BlockEntityFormatters) {
	s := fmt.Sprintf("\ntype %s_t is record\n", stream.Name)

	// NOTE: Params and returns are generated in the same function.
	// However, a stream must have only params or only returns, so length
	// of at least one iteration is 0.

	// Downstream
	for _, p := range stream.Params {
		if p.IsArray {
			s += fmt.Sprintf("   %s : slv_vector(%d downto 0)(%d downto 0);\n", p.Name, p.Count-1, p.Width-1)
		} else {
			s += fmt.Sprintf("   %s : std_logic_vector(%d downto 0);\n", p.Name, p.Width-1)
		}
	}

	// Upstream
	for _, r := range stream.Returns {
		if r.IsArray {
			s += fmt.Sprintf("   %s : slv_vector(%d downto 0)(%d downto 0);\n", r.Name, r.Count-1, r.Width-1)
		} else {
			s += fmt.Sprintf("   %s : std_logic_vector(%d downto 0);\n", r.Name, r.Width-1)
		}
	}

	s += "end record;\n"

	fmts.StreamTypes += s
}

func genStreamPorts(stream *fn.Stream, fmts *BlockEntityFormatters) {
	dir := "out"
	suffix := "o"

	if stream.IsUpstream() {
		dir = "in"
		suffix = "i"
	}

	s := fmt.Sprintf(";\n   %s_%s : %s %[1]s_t;\n", stream.Name, suffix, dir)

	s += fmt.Sprintf("   %s_stb_o : out std_logic", stream.Name)

	fmts.EntityFunctionalPorts += s
}

func genUpstreamAccess(stream *fn.Stream, fmts *BlockEntityFormatters) {
	for _, r := range stream.Returns {
		switch acs := r.Access.(type) {
		case access.SingleOneReg:
			addr := [2]int64{acs.StartAddr(), acs.StartAddr()}
			code := fmt.Sprintf(
				"      master_in.dat(%d downto %d) <= %s_i.%s;\n",
				acs.EndBit(), acs.StartBit(), stream.Name, r.Name,
			)

			fmts.RegistersAccess.add(addr, code)
		case access.SingleNRegs:
			chunks := makeAccessChunksContinuous(acs, Compact)

			for _, c := range chunks {
				code := fmt.Sprintf(
					"      master_in.dat(%[5]d downto %[6]d) <= %[1]s_i.%[2]s(%[3]s downto %[4]s);",
					stream.Name, r.Name, c.range_[0], c.range_[1], c.endBit, c.startBit,
				)

				fmts.RegistersAccess.add([2]int64{c.addr[0], c.addr[1]}, code)
			}
		default:
			panic("unimplemented")
		}
	}
}

func genDownstreamAccess(stream *fn.Stream, fmts *BlockEntityFormatters) {
	for _, p := range stream.Params {
		switch acs := p.Access.(type) {
		case access.SingleOneReg:
			addr := [2]int64{acs.StartAddr(), acs.StartAddr()}
			code := fmt.Sprintf(`
      if master_out.we = '1' then
         %s_o.%s <= master_out.dat(%d downto %d);
      end if;
`,
				stream.Name, p.Name, acs.EndBit(), acs.StartBit(),
			)

			fmts.RegistersAccess.add(addr, code)
		case access.SingleNRegs:
			chunks := makeAccessChunksContinuous(acs, Compact)

			for _, c := range chunks {
				code := fmt.Sprintf(`
      if master_out.we = '1' then
         %[1]s_o.%[2]s(%[3]s downto %[4]s) <= master_out.dat(%[5]d downto %[6]d);
      end if;
`,
					stream.Name, p.Name, c.range_[0], c.range_[1], c.endBit, c.startBit,
				)

				fmts.RegistersAccess.add([2]int64{c.addr[0], c.addr[1]}, code)
			}
		default:
			panic("unimplemented")
		}
	}
}

func genStreamStrobe(stream *fn.Stream, fmts *BlockEntityFormatters) {
	clear := fmt.Sprintf("\n%s_stb_o <= '0';", stream.Name)

	fmts.StreamsStrobesClear += clear

	weVal := "1"
	if stream.IsUpstream() {
		weVal = "0"
	}

	stbSet := `
   %s_stb : if addr = %d then
      if master_out.we = '%s' then
         %[1]s_stb_o <= '1';
      end if;
   end if;
`
	set := fmt.Sprintf(stbSet, stream.Name, stream.StbAddr, weVal)

	fmts.StreamsStrobesSet += set
}
