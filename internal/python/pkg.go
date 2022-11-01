package python

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/elem"
)

func genPkgConsts(pkgsConsts map[string]*elem.Package) string {
	s := ""

	for pkgName, pkg := range pkgsConsts {
		if pkg.ConstContainer.Empty() {
			continue
		}

		s += fmt.Sprintf("class %sPkg:\n", pkgName)
		increaseIndent(1)
		s += genConsts(&pkg.ConstContainer)
		decreaseIndent(1)
	}

	return s
}
