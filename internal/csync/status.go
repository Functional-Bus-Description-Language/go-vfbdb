package csync

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/access"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/elem"
	"github.com/Functional-Bus-Description-Language/go-vfbdb/internal/c"
	"github.com/Functional-Bus-Description-Language/go-vfbdb/internal/utils"
)

func genStatus(st elem.Status, hFmts *BlockHeaderFormatters, srcFmts *BlockSourceFormatters) {
	if st.IsArray() {
		panic("not yet implemented")
	} else {
		genStatusSingle(st, hFmts, srcFmts)
	}
}

func genStatusSingle(st elem.Status, hFmts *BlockHeaderFormatters, srcFmts *BlockSourceFormatters) {
	switch st.Access().(type) {
	case access.SingleSingle:
		genStatusSingleSingle(st, hFmts, srcFmts)
	case access.SingleContinuous:
		panic("not yet implemented")
	default:
		panic("unknown single access strategy")
	}
}

func genStatusSingleSingle(st elem.Status, hFmts *BlockHeaderFormatters, srcFmts *BlockSourceFormatters) {
	typ := c.WidthToReadType(st.Width())
	signature := fmt.Sprintf(
		"\n\nint vfbdb_%s_%s_read(const vfbdb_iface_t * const iface, %s const data)",
		hFmts.BlockName, st.Name(), typ.String(),
	)

	hFmts.Code += fmt.Sprintf("%s;", signature)

	a := st.Access().(access.SingleSingle)
	srcFmts.Code += fmt.Sprintf("%s {\n", signature)
	if readDataType.Typ() != "ByteArray" && typ.Typ() != "ByteArray" {
		if busWidth == st.Width() {
			srcFmts.Code += fmt.Sprintf(
				"\treturn iface->read(%d, data);\n};", a.Addr,
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
};`, readDataType.Depointer().String(), a.Addr, a.StartBit(), utils.Uint64Mask(a.StartBit(), a.EndBit()),
			)
		}
	} else {
		panic("not yet implemented")
	}
}
