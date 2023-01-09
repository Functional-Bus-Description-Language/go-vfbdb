package python

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/elem"
)

func genProc(p *elem.Proc, blk *elem.Block) string {
	code := indent + fmt.Sprintf("self.%s = Proc(iface, %d, ",
		p.Name, blk.StartAddr()+p.ParamsStartAddr(),
	)
	code += genParamList(p.Params)
	code += ", "
	code += genReturnList(p.Returns)
	code += ")\n"

	return code
}
