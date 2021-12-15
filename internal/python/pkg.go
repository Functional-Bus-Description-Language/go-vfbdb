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

		for name, i := range pkg.IntConsts {
			s += indent + fmt.Sprintf("%s = %d\n", name, i)
		}
		for name, str := range pkg.StrConsts {
			s += indent + fmt.Sprintf("%s = %q\n", name, str)
		}

		decreaseIndent(1)
	}

	return s
}
