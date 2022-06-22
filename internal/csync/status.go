package csync

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/access"
	"github.com/Functional-Bus-Description-Language/go-vfbdb/internal/c"
)

func genStatus(st *fbdl.Status, hFmts *BlockHeaderFormatters, srcFmts *BlockSourceFormatters) {
	if st.IsArray {
		panic("not yet implemented")
	} else {
		genStatusSingle(st, hFmts, srcFmts)
	}
}

func genStatusSingle(st *fbdl.Status, hFmts *BlockHeaderFormatters, srcFmts *BlockSourceFormatters) {
	switch st.Access.(type) {
	case access.SingleSingle:
		genStatusSingleSingle(st, hFmts, srcFmts)
	case access.SingleContinuous:
		panic("not yet implemented")
	default:
		panic("unknown single access strategy")
	}
}

func genStatusSingleSingle(st *fbdl.Status, hFmts *BlockHeaderFormatters, srcFmts *BlockSourceFormatters) {
	typ := c.WidthToReadType(st.Width)
	signature := fmt.Sprintf(
		"\n\nint vfbdb_%s_%s_read(const vfbdb_iface_t * const iface, %s const data)",
		hFmts.BlockName, st.Name, typ.String(),
	)

	hFmts.Code += fmt.Sprintf("%s;", signature)

	access := st.Access.(access.SingleSingle)
	srcFmts.Code += fmt.Sprintf("%s {\n", signature)
	if readDataType.Typ() != "ByteArray" && typ.Typ() != "ByteArray" {
		if busWidth == st.Width {
			srcFmts.Code += fmt.Sprintf(
				"\treturn iface->read(%d, data);\n};", access.Addr,
			)
		} else {
			srcFmts.Code += fmt.Sprintf(`	int err;
	%s aux;
	err = iface->read(%d, &aux);
	if (err) {
		return err;
	}
	*data = (aux >> %d) & %x;
	return 0;
};`, readDataType.Depointer().String(), access.Addr, access.Mask.Lower, access.Mask.Uint64(),
			)
		}
	} else {
		panic("not yet implemented")
	}
}
