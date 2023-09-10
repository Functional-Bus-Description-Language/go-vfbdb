package csync

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/access"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/fn"
	"github.com/Functional-Bus-Description-Language/go-vfbdb/internal/c"
	"github.com/Functional-Bus-Description-Language/go-vfbdb/internal/utils"
)

func genStatus(st *fn.Status, blk *fn.Block, hFmts *BlockHFormatters, cFmts *BlockCFormatters) {
	if st.IsArray {
		panic("unimplemented")
	} else {
		genStatusSingle(st, blk, hFmts, cFmts)
	}
}

func genStatusSingle(st *fn.Status, blk *fn.Block, hFmts *BlockHFormatters, cFmts *BlockCFormatters) {
	switch st.Access.(type) {
	case access.SingleOneReg:
		genStatusSingleOneReg(st, blk, hFmts, cFmts)
	default:
		panic("unimplemented")
	}
}

func genStatusSingleOneReg(st *fn.Status, blk *fn.Block, hFmts *BlockHFormatters, cFmts *BlockCFormatters) {
	typ := c.WidthToReadType(st.Width)
	signature := fmt.Sprintf(
		"int vfbdb_%s_%s_read(const vfbdb_iface_t * const iface, %s const data)",
		hFmts.BlockName, st.Name, typ.String(),
	)

	hFmts.Code += fmt.Sprintf("\n%s;\n", signature)

	acs := st.Access.(access.SingleOneReg)
	cFmts.Code += fmt.Sprintf("\n%s {\n", signature)
	if readType.Typ() != "ByteArray" && typ.Typ() != "ByteArray" {
		if busWidth == st.Width {
			cFmts.Code += fmt.Sprintf(
				"\treturn iface->read(%d, data);\n};\n", blk.StartAddr()+acs.Addr,
			)
		} else {
			cFmts.Code += fmt.Sprintf(`	%s aux;
	const int err = iface->read(%d, &aux);
	if (err)
		return err;
	*data = (aux >> %d) & 0x%x;
	return 0;
};
`, readType.Depointer().String(), blk.StartAddr()+acs.Addr, acs.StartBit, utils.Uint64Mask(acs.StartBit, acs.EndBit),
			)
		}
	} else {
		panic("unimplemented")
	}
}
