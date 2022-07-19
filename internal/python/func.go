package python

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/elem"
)

func genFunc(fun elem.Func, blk elem.Block) string {
	code := indent + fmt.Sprintf("self.%s = Func(iface, %d, ",
		fun.Name(), blk.AddrSpace().Start()+fun.ParamsStartAddr(),
	)
	code += genParamAccessList(fun.Params())
	code += genReturnAccessList(fun.Returns())
	code += ")\n"

	return code
}
