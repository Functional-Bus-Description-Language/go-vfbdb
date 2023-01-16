package python

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/elem"
)

func genProc(p *elem.Proc, blk *elem.Block) string {
	if len(p.Params) == 0 && len(p.Returns) == 0 {
		return genEmptyProc(p, blk)
	}

	code := indent + fmt.Sprintf("self.%s = Proc(iface, %d, ",
		p.Name, blk.StartAddr()+p.ParamsStartAddr(),
	)
	code += genParamList(p.Params)
	code += ", "
	code += genReturnList(p.Returns)
	code += ")\n"

	return code
}

func genEmptyProc(p *elem.Proc, blk *elem.Block) string {
	delay := "None"
	exitAddr := "None"
	if p.Delay != nil {
		delay = fmt.Sprintf("%d + %d * 1e-9", p.Delay.S, p.Delay.Ns)
		exitAddr = fmt.Sprintf("%d", blk.StartAddr()+*p.ExitAddr)
	}
	code := indent + fmt.Sprintf("self.%s = EmptyProc(iface, %d, %s, %s)\n",
		p.Name, blk.StartAddr()+*p.CallAddr, delay, exitAddr,
	)

	return code
}
