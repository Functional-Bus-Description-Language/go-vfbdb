package python

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl"
)

func generatePkgConsts(pkgsConsts map[string]fbdl.Package) string {
	s := ""

	for pkgName, pkg := range pkgsConsts {
		if !pkg.HasConsts() {
			continue
		}

		s += fmt.Sprintf("class %sPkg:\n", pkgName)
		increaseIndent(1)
		s += generateConsts(pkg.ConstContainer)
		decreaseIndent(1)
	}

	return s
}
