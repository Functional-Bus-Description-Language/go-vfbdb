package python

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/elem"
)

func genProc(p *elem.Proc, blk *elem.Block) string {
	if len(p.Params) == 0 && len(p.Returns) == 0 {
		return genEmptyProc(p, blk)
	} else if len(p.Params) > 0 && len(p.Returns) == 0 {
		return genParamsProc(p, blk)
	} else if len(p.Params) == 0 && len(p.Returns) > 0 {
		return genReturnsProc(p, blk)
	}
	panic("not yet implemented")
}

func genEmptyProc(p *elem.Proc, blk *elem.Block) string {
	delay, exitAddr := genDelayAndExitAddr(p, blk)
	code := indent + fmt.Sprintf("self.%s = EmptyProc(iface, %d, %s, %s)\n",
		p.Name, blk.StartAddr()+*p.CallAddr, delay, exitAddr,
	)

	return code
}

func genParamsProc(p *elem.Proc, blk *elem.Block) string {
	code := indent + fmt.Sprintf("self.%s = ParamsProc(iface, %d, ",
		p.Name, blk.StartAddr()+p.ParamsStartAddr(),
	)
	code += genParamList(p.Params)
	delay, exitAddr := genDelayAndExitAddr(p, blk)
	code += fmt.Sprintf(", %s, %s)\n", delay, exitAddr)

	return code
}

func genDelayAndCallAddr(p *elem.Proc, blk *elem.Block) (string, string) {
	delay := "None"
	callAddr := "None"
	if p.Delay != nil {
		delay = fmt.Sprintf("%d + %d * 1e-9", p.Delay.S, p.Delay.Ns)
		callAddr = fmt.Sprintf("%d", blk.StartAddr()+*p.CallAddr)
	}
	return delay, callAddr
}

func genDelayAndExitAddr(p *elem.Proc, blk *elem.Block) (string, string) {
	delay := "None"
	exitAddr := "None"
	if p.Delay != nil {
		delay = fmt.Sprintf("%d + %d * 1e-9", p.Delay.S, p.Delay.Ns)
		exitAddr = fmt.Sprintf("%d", blk.StartAddr()+*p.ExitAddr)
	}
	return delay, exitAddr
}

func genReturnsProc(p *elem.Proc, blk *elem.Block) string {
	code := indent + fmt.Sprintf("self.%s = ReturnsProc(iface, %d, ",
		p.Name, blk.StartAddr()+p.ReturnsStartAddr(),
	)
	code += genReturnList(p.Returns)
	delay, callAddr := genDelayAndCallAddr(p, blk)
	code += fmt.Sprintf(", %s, %s)\n", delay, callAddr)

	return code
}
