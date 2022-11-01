package python

import (
	"fmt"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/elem"
	"strings"
)

func genReturnList(returns []*elem.Return) string {
	if len(returns) == 0 {
		return "None"
	}

	b := strings.Builder{}

	b.WriteString("[\n")
	increaseIndent(2)

	for _, r := range returns {
		b.WriteString(fmt.Sprintf("%s{'Name': '%s', 'Access': ", indent, r.Name))
		genAccess(r.Access, &b)
		b.WriteString("},\n")
	}

	decreaseIndent(1)
	b.WriteString(fmt.Sprintf("%s]", indent))
	decreaseIndent(1)

	return b.String()
}
