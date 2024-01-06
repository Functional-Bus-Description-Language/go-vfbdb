package python

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/fn"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/val"
)

func genProc(p *fn.Proc, blk *fn.Block) string {
	if len(p.Params) == 0 && len(p.Returns) == 0 {
		return genEmptyProc(p, blk)
	} else if len(p.Params) > 0 && len(p.Returns) == 0 {
		return genParamsProc(p, blk)
	} else if len(p.Params) == 0 && len(p.Returns) > 0 {
		return genReturnsProc(p, blk)
	} else {
		return genParamsAndReturnsProc(p, blk)
	}
}

func genEmptyProc(p *fn.Proc, blk *fn.Block) string {
	delay, exitAddr := genDelayAndExitAddr(p, blk)
	code := indent + fmt.Sprintf("self.%s = EmptyProc(iface, %d, %s, %s)\n",
		p.Name, blk.StartAddr()+*p.CallAddr, delay, exitAddr,
	)

	return code
}

func genParamsProc(p *fn.Proc, blk *fn.Block) string {
	code := indent + fmt.Sprintf("self.%s = ParamsProc(iface, %d, ",
		p.Name, blk.StartAddr()+p.ParamsStartAddr(),
	)
	code += genParamList(p.Params, blk)
	delay, exitAddr := genDelayAndExitAddr(p, blk)
	code += fmt.Sprintf(", %s, %s)\n", delay, exitAddr)

	return code
}

func genDelay(d *val.Time) string {
	delay := "None"
	if d != nil {
		delay = fmt.Sprintf("%d + %d * 1e-9", d.S, d.Ns)
	}
	return delay
}

func genDelayAndCallAddr(p *fn.Proc, blk *fn.Block) (string, string) {
	delay := "None"
	callAddr := "None"
	if p.Delay != nil {
		delay = fmt.Sprintf("%d + %d * 1e-9", p.Delay.S, p.Delay.Ns)
		callAddr = fmt.Sprintf("%d", blk.StartAddr()+*p.CallAddr)
	}
	return delay, callAddr
}

func genDelayAndExitAddr(p *fn.Proc, blk *fn.Block) (string, string) {
	delay := "None"
	exitAddr := "None"
	if p.Delay != nil {
		delay = fmt.Sprintf("%d + %d * 1e-9", p.Delay.S, p.Delay.Ns)
		exitAddr = fmt.Sprintf("%d", blk.StartAddr()+*p.ExitAddr)
	}
	return delay, exitAddr
}

func genReturnsProc(p *fn.Proc, blk *fn.Block) string {
	code := indent + fmt.Sprintf("self.%s = ReturnsProc(iface, %d, ",
		p.Name, blk.StartAddr()+p.ReturnsStartAddr(),
	)
	code += genReturnList(p.Returns, blk)
	delay, callAddr := genDelayAndCallAddr(p, blk)
	code += fmt.Sprintf(", %s, %s)\n", delay, callAddr)

	return code
}

func genParamsAndReturnsProc(p *fn.Proc, blk *fn.Block) string {
	code := indent + fmt.Sprintf(
		"self.%s = ParamsAndReturnsProc(iface, %d, %s, %d, %s, %s)\n",
		p.Name,
		blk.StartAddr()+p.ParamsStartAddr(), genParamList(p.Params, blk),
		blk.StartAddr()+p.ReturnsStartAddr(), genReturnList(p.Returns, blk),
		genDelay(p.Delay),
	)

	return code
}
