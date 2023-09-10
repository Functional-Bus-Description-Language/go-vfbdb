package csync

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/access"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/fn"
	"github.com/Functional-Bus-Description-Language/go-vfbdb/internal/c"
	_ "github.com/Functional-Bus-Description-Language/go-vfbdb/internal/utils"
	"strings"
)

func genProc(p *fn.Proc, blk *fn.Block, hFmts *BlockHFormatters, cFmts *BlockCFormatters) {
	sig := genProcSignature(p, blk, hFmts)

	hFmts.Code += "\n" + sig + ";\n"

	cFmts.Code += fmt.Sprintf("\n%s {\n", sig)
	if len(p.Params) == 0 && len(p.Returns) == 0 {
		cFmts.Code += fmt.Sprintf("\treturn iface->write(%d, 0);\n", blk.StartAddr()+*p.CallAddr)
	}

	if len(p.Params) > 0 {
		genProcParamsAccess(p, blk, cFmts)
	}

	if len(p.Returns) > 0 {
		genProcReturnsAccess(p, blk, cFmts)
	}

	cFmts.Code += "};\n"
}

func genProcSignature(p *fn.Proc, blk *fn.Block, hFmts *BlockHFormatters) string {
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

func genProcParamsAccess(p *fn.Proc, blk *fn.Block, cFmts *BlockCFormatters) {
	if p.ParamsBufSize() == 1 {
		genProcParamsAccessSingleWrite(p, blk, cFmts)
	} else {
		genProcParamsAccessBlockWrite(p, blk, cFmts)
	}
}

func genProcParamsAccessSingleWrite(p *fn.Proc, blk *fn.Block, cFmts *BlockCFormatters) {
	if p.Delay == nil && len(p.Returns) == 0 {
		genProcParamsAccessSingleWriteNoDelayNoReturns(p, blk, cFmts)
	} else {
		panic("not yet implemented")
	}
}

func genProcParamsAccessSingleWriteNoDelayNoReturns(p *fn.Proc, blk *fn.Block, cFmts *BlockCFormatters) {
	cFmts.Code += fmt.Sprintf("\treturn iface->write(%d, ", blk.StartAddr()+*p.CallAddr)
	for i, p := range p.Params {
		if i != 0 {
			cFmts.Code += " | "
		}

		switch acs := p.Access.(type) {
		case access.SingleOneReg:
			cFmts.Code += fmt.Sprintf("%s << %d", p.Name, acs.StartBit)
		default:
			panic("not yet implemented")
		}
	}
	cFmts.Code += ");\n"
}

func genProcParamsAccessBlockWrite(p *fn.Proc, blk *fn.Block, cFmts *BlockCFormatters) {
	if p.Delay == nil && len(p.Returns) == 0 {
		genProcParamsAccessBlockWriteNoDelayNoReturns(p, blk, cFmts)
	} else {
		panic("not yet implemented")
	}
}

func genProcParamsAccessBlockWriteNoDelayNoReturns(proc *fn.Proc, blk *fn.Block, cFmts *BlockCFormatters) {
	cFmts.Code += fmt.Sprintf("\t%s buf[%d] = {0};\n\n", c.WidthToWriteType(blk.Width), proc.ParamsBufSize())

	for _, p := range proc.Params {
		switch acs := p.Access.(type) {
		case access.SingleOneReg:
			cFmts.Code += fmt.Sprintf(
				"\tbuf[%d] |= %s << %d;\n",
				acs.Addr-proc.ParamsStartAddr(), p.Name, acs.StartBit,
			)
		default:
			panic("not yet implemented")
		}
	}

	cFmts.Code += fmt.Sprintf(
		"\n\treturn iface->writeb(%d, buf, %d);\n",
		blk.StartAddr()+proc.ParamsStartAddr(), proc.ParamsBufSize(),
	)
}

func genProcReturnsAccess(p *fn.Proc, blk *fn.Block, cFmts *BlockCFormatters) {
	if p.ReturnsBufSize() == 1 {
		genProcReturnsAccessSingleRead(p, blk, cFmts)
	} else {
		genProcReturnsAccessBlockRead(p, blk, cFmts)
	}
}

func genProcReturnsAccessSingleRead(p *fn.Proc, blk *fn.Block, cFmts *BlockCFormatters) {
	cFmts.Code += fmt.Sprintf("\t%s _rdata;\n", c.WidthToWriteType(blk.Width))

	cFmts.Code += fmt.Sprintf("\tconst int err = iface->read(%d, &_rdata);\n", *p.ExitAddr)
	cFmts.Code += "\tif (err)\n\t\t return err;\n"

	for _, r := range p.Returns {
		switch acs := r.Access.(type) {
		case access.SingleOneReg:
			cFmts.Code += fmt.Sprintf(
				"\t*%s = (_rdata >> %d) & 0x%X;\n",
				r.Name, acs.StartBit, c.MaskToValue(acs.StartBit, acs.EndBit),
			)
		default:
			panic("not yet implemented")
		}
	}
	cFmts.Code += "\treturn 0;\n"
}

func genProcReturnsAccessBlockRead(p *fn.Proc, blk *fn.Block, cFmts *BlockCFormatters) {
	panic("not yet implemented")
	/*
		cFmts.Code += fmt.Sprintf(
			"\t%s _rbuff[%d];\n", c.WidthToWriteType(blk.Width), p.ReturnsBufSize(),
		)
		cFmts.Code += fmt.Sprintf(
			"\tconst int err = iface.readb(%d, _rbuff, %d);\n", p.ReturnsStartAddr(), p.ReturnsBufSize(),
		)
		cFmts.Code += "\tif (err)\n\t\t return err;\n"
	*/
}
