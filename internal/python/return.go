package python

import (
	"fmt"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/fn"
	"strings"
)

func genReturnList(returns []*fn.Return, blk *fn.Block) string {
	if len(returns) == 0 {
		return "None"
	}

	b := strings.Builder{}

	b.WriteString("[\n")
	increaseIndent(2)

	for _, r := range returns {
		b.WriteString(fmt.Sprintf("%s{'Name': '%s', 'Access': ", indent, r.Name))
		genAccess(r.Access, blk.StartAddr(), &b)
		b.WriteString("},\n")
	}

	decreaseIndent(1)
	b.WriteString(fmt.Sprintf("%s]", indent))
	decreaseIndent(1)

	return b.String()
}
