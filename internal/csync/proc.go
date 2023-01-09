package csync

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/access"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/elem"
	"github.com/Functional-Bus-Description-Language/go-vfbdb/internal/c"
	_ "github.com/Functional-Bus-Description-Language/go-vfbdb/internal/utils"
	"strings"
)

func genProc(p *elem.Proc, blk *elem.Block, hFmts *BlockHFormatters, cFmts *BlockCFormatters) {
	sig := genProcSignature(p, blk, hFmts)

	hFmts.Code += "\n" + sig + ";\n"

	cFmts.Code += fmt.Sprintf("\n%s {\n", sig)
	if len(p.Params) == 0 && len(p.Returns) == 0 {
		cFmts.Code += fmt.Sprintf("\treturn iface->write(%d, 0);\n};\n", blk.StartAddr()+p.StbAddr)
		return
	}

	if len(p.Params) > 0 {
		genProcParamsAccess(p, blk, cFmts)
	}

	if len(p.Returns) > 0 {
		genProcReturnsAccess(p, blk, cFmts)
	}

	cFmts.Code += "};\n"
}

func genProcSignature(p *elem.Proc, blk *elem.Block, hFmts *BlockHFormatters) string {
	prefix := "int vfbdb_" + hFmts.BlockName + "_" + p.Name

	params := strings.Builder{}
	params.WriteString("const vfbdb_iface_t * const iface")

	for _, p := range p.Params {
		params.WriteString(
			", const " + c.WidthToWriteType(p.Width).String() + " " + p.Name,
		)
	}

	for _, r := range p.Returns {
		params.WriteString(
			", " + c.WidthToReadType(r.Width).String() + " const " + r.Name,
		)
	}

	return prefix + "(" + params.String() + ")"
}

func genProcParamsAccess(p *elem.Proc, blk *elem.Block, cFmts *BlockCFormatters) {
	if p.ParamsBufSize() == 1 {
		genProcParamsAccessSingleReg(p, blk, cFmts)
	}
}

func genProcParamsAccessSingleReg(p *elem.Proc, blk *elem.Block, cFmts *BlockCFormatters) {
	cFmts.Code += fmt.Sprintf("\treturn iface->write(%d, ", blk.StartAddr()+p.StbAddr)
	for i, p := range p.Params {
		if i != 0 {
			cFmts.Code += " | "
		}

		switch a := p.Access.(type) {
		case access.SingleSingle:
			cFmts.Code += fmt.Sprintf("%s << %d", p.Name, a.StartBit())
		default:
			panic("not yet implemented")
		}
	}
	cFmts.Code += ");\n"
}

func genProcReturnsAccess(p *elem.Proc, blk *elem.Block, cFmts *BlockCFormatters) {
	panic("not implemented")
}
