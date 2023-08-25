package csync

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/access"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/fn"
	"github.com/Functional-Bus-Description-Language/go-vfbdb/internal/c"
	"github.com/Functional-Bus-Description-Language/go-vfbdb/internal/utils"
)

func genConfig(cfg *fn.Config, blk *fn.Block, hFmts *BlockHFormatters, cFmts *BlockCFormatters) {
	if cfg.IsArray {
		panic("not yet implemented")
	} else {
		genConfigSingle(cfg, blk, hFmts, cFmts)
	}
}

func genConfigSingle(cfg *fn.Config, blk *fn.Block, hFmts *BlockHFormatters, cFmts *BlockCFormatters) {
	switch cfg.Access.(type) {
	case access.SingleSingle:
		genConfigSingleSingle(cfg, blk, hFmts, cFmts)
	case access.SingleContinuous:
		panic("not yet implemented")
	default:
		panic("unknown single access strategy")
	}
}

func genConfigSingleSingle(cfg *fn.Config, blk *fn.Block, hFmts *BlockHFormatters, cFmts *BlockCFormatters) {
	rType := c.WidthToReadType(cfg.Width)
	wType := c.WidthToWriteType(cfg.Width)

	readSignature := fmt.Sprintf(
		"int vfbdb_%s_%s_read(const vfbdb_iface_t * const iface, %s const data)",
		hFmts.BlockName, cfg.Name, rType.String(),
	)
	writeSignature := fmt.Sprintf(
		"int vfbdb_%s_%s_write(const vfbdb_iface_t * const iface, %s const data)",
		hFmts.BlockName, cfg.Name, wType.String(),
	)

	hFmts.Code += fmt.Sprintf("\n%s;\n%s;\n", readSignature, writeSignature)

	a := cfg.Access.(access.SingleSingle)
	cFmts.Code += fmt.Sprintf("\n%s {\n", readSignature)
	if readType.Typ() != "ByteArray" && rType.Typ() != "ByteArray" {
		if busWidth == cfg.Width {
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

	cFmts.Code += fmt.Sprintf("\n%s {\n", writeSignature)
	if readType.Typ() != "ByteArray" && rType.Typ() != "ByteArray" {
		if busWidth == cfg.Width {
			cFmts.Code += fmt.Sprintf(
				"\treturn iface->write(%d, data);\n};\n", blk.StartAddr()+a.Addr,
			)
		} else {
			cFmts.Code += fmt.Sprintf(
				"	return iface->write(%d, (data << %d));\n };", blk.StartAddr()+a.Addr, a.StartBit(),
			)
		}
	} else {
		panic("not yet implemented")
	}
}
