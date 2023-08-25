package python

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/pkg"
)

func genPkgConsts(pkgsConsts map[string]*pkg.Package) string {
	s := ""

	for pkgName, pkg := range pkgsConsts {
		if pkg.Consts.Empty() {
			continue
		}

		s += fmt.Sprintf("class %sPkg:\n", pkgName)
		increaseIndent(1)
		s += genConsts(&pkg.Consts)
		decreaseIndent(1)
	}

	return s
}
