package csync

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/access"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/elem"
	"github.com/Functional-Bus-Description-Language/go-vfbdb/internal/c"
	"github.com/Functional-Bus-Description-Language/go-vfbdb/internal/utils"
)

func genStatus(st elem.Status, hFmts *BlockHFormatters, cFmts *BlockCFormatters) {
	if st.IsArray() {
		panic("not yet implemented")
	} else {
		genStatusSingle(st, hFmts, cFmts)
	}
}

func genStatusSingle(st elem.Status, hFmts *BlockHFormatters, cFmts *BlockCFormatters) {
	switch st.Access().(type) {
	case access.SingleSingle:
		genStatusSingleSingle(st, hFmts, cFmts)
	case access.SingleContinuous:
		panic("not yet implemented")
	default:
		panic("unknown single access strategy")
	}
}

func genStatusSingleSingle(st elem.Status, hFmts *BlockHFormatters, cFmts *BlockCFormatters) {
	typ := c.WidthToReadType(st.Width())
	signature := fmt.Sprintf(
		"\n\nint vfbdb_%s_%s_read(const vfbdb_iface_t * const iface, %s const data)",
		hFmts.BlockName, st.Name(), typ.String(),
	)

	hFmts.Code += fmt.Sprintf("%s;", signature)

	a := st.Access().(access.SingleSingle)
	cFmts.Code += fmt.Sprintf("%s {\n", signature)
	if readType.Typ() != "ByteArray" && typ.Typ() != "ByteArray" {
		if busWidth == st.Width() {
			cFmts.Code += fmt.Sprintf(
				"\treturn iface->read(%d, data);\n};", a.Addr,
			)
		} else {
			cFmts.Code += fmt.Sprintf(`	int err;
	%s aux;
	err = iface->read(%d, &aux);
	if (err) {
		return err;
	}
	*data = (aux >> %d) & %x;
	return 0;
};`, readType.Depointer().String(), a.Addr, a.StartBit(), utils.Uint64Mask(a.StartBit(), a.EndBit()),
			)
		}
	} else {
		panic("not yet implemented")
	}
}
