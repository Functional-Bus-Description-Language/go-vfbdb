package csync

import (
	"fmt"
	"strconv"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/access"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/elem"
	"github.com/Functional-Bus-Description-Language/go-vfbdb/internal/c"
	"github.com/Functional-Bus-Description-Language/go-vfbdb/internal/utils"
)

func genStatic(st *elem.Static, blk *elem.Block, hFmts *BlockHFormatters, cFmts *BlockCFormatters) {
	if st.IsArray {
		panic("not yet implemented")
	} else {
		genStaticSingle(st, blk, hFmts, cFmts)
	}
}

func genStaticSingle(st *elem.Static, blk *elem.Block, hFmts *BlockHFormatters, cFmts *BlockCFormatters) {
	switch st.Access.(type) {
	case access.SingleSingle:
		genStaticSingleSingle(st, blk, hFmts, cFmts)
	case access.SingleContinuous:
		panic("not yet implemented")
	default:
		panic("unknown single access strategy")
	}
}

func genStaticSingleSingle(st *elem.Static, blk *elem.Block, hFmts *BlockHFormatters, cFmts *BlockCFormatters) {
	wTyp := c.WidthToWriteType(st.Width)
	rTyp := c.WidthToReadType(st.Width)

	hFmts.Code += fmt.Sprintf(
		"\nextern const %s vfbdb_%s_%s;\n",
		wTyp.String(), hFmts.BlockName, st.Name,
	)

	signature := fmt.Sprintf(
		"int vfbdb_%s_%s_read(const vfbdb_iface_t * const iface, %s const data)",
		hFmts.BlockName, st.Name, rTyp.String(),
	)

	hFmts.Code += fmt.Sprintf("%s;\n", signature)

	cFmts.Code += fmt.Sprintf(
		"\nconst %s vfbdb_%s_%s = %s;\n",
		wTyp.String(), hFmts.BlockName, st.Name,
		// XXX: Uint64 is currently used. Below code needs fix if static is longer than 64 bits.
		fmt.Sprintf("0x%s", strconv.FormatUint(st.Default.Uint64(), 16)),
	)

	a := st.Access.(access.SingleSingle)
	cFmts.Code += fmt.Sprintf("%s {\n", signature)
	if readType.Typ() != "ByteArray" && rTyp.Typ() != "ByteArray" {
		if busWidth == st.Width {
			cFmts.Code += fmt.Sprintf(
				"\treturn iface->read(%d, data);\n};\n", blk.StartAddr()+a.Addr,
			)
		} else {
			cFmts.Code += fmt.Sprintf(`	%s aux;
	const int err = iface->read(%d, &aux);
	if (err)
		return err;
	*data = (aux >> %d) & 0x%x;
	return 0;
};
`, readType.Depointer().String(), blk.StartAddr()+a.Addr, a.StartBit(), utils.Uint64Mask(a.StartBit(), a.EndBit()),
			)
		}
	} else {
		panic("not yet implemented")
	}
}
