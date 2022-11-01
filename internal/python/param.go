package python

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/elem"

	"strings"
)

func genParamList(params []*elem.Param) string {
	if len(params) == 0 {
		return "None"
	}

	b := strings.Builder{}

	b.WriteString("[\n")
	increaseIndent(2)

	for _, p := range params {
		b.WriteString(
			fmt.Sprintf("%s{'Name': '%s', 'Width': %d, 'Access': ", indent, p.Name, p.Width),
		)
		genAccess(p.Access, &b)
		b.WriteString("},\n")
	}

	decreaseIndent(1)
	b.WriteString(fmt.Sprintf("%s]", indent))
	decreaseIndent(1)

	return b.String()
}
