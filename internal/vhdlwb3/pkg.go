package vhdlwb3

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/pkg"
)

//go:embed templates/wb3.vhd
var wb3PkgStr string
var wb3PkgTmpl = template.Must(template.New("VHDL wb3 package").Parse(wb3PkgStr))

type wb3PackageFormatters struct {
	PkgsConsts string
}

func genWb3Package(pkgsConsts map[string]*pkg.Package) {
	filePath := outputPath + "wb3.vhd"

	f, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("generate vhdl-wb3: %v", err)
	}
	defer f.Close()

	fmts := wb3PackageFormatters{PkgsConsts: genPkgsConsts(pkgsConsts)}

	err = wb3PkgTmpl.Execute(f, fmts)
	if err != nil {
		log.Fatalf("generate vhdl-wb3: %v", err)
	}

	addGeneratedFile(filePath)
}

func genPkgsConsts(pkgsConsts map[string]*pkg.Package) string {
	s := ""

	for pkgName, pkg := range pkgsConsts {
		if pkg.Consts.Empty() {
			continue
		}

		// Package type definition
		s += fmt.Sprintf("type %s_pkg_t is record\n", pkgName)
		for name := range pkg.Consts.Bools {
			s += fmt.Sprintf("   %s : boolean;\n", name)
		}
		for name, list := range pkg.Consts.BoolLists {
			s += fmt.Sprintf("   %s : boolean_vector(0 to %d);\n", name, len(list)-1)
		}
		for name := range pkg.Consts.Ints {
			s += fmt.Sprintf("   %s : int64;\n", name)
		}
		for name, list := range pkg.Consts.IntLists {
			s += fmt.Sprintf("   %s : int64_vector(0 to %d);\n", name, len(list)-1)
		}
		for name := range pkg.Consts.Strings {
			s += fmt.Sprintf("   %s : string;\n", name)
		}
		s += "end record;\n"

		// Package constant definition
		s += fmt.Sprintf("constant %[1]s_pkg : %[1]s_pkg_t := (\n", pkgName)
		for name, b := range pkg.Consts.Bools {
			s += fmt.Sprintf("   %s => %t,\n", name, b)
		}
		for name, list := range pkg.Consts.BoolLists {
			s += fmt.Sprintf("   %s => (", name)
			for i, b := range list {
				s += fmt.Sprintf("%d => %t, ", i, b)
			}
			s = s[:len(s)-2]
			s += "),\n"
		}
		for name, i := range pkg.Consts.Ints {
			s += fmt.Sprintf("   %s => signed'(x\"%016x\"),\n", name, i)
		}
		for name, list := range pkg.Consts.IntLists {
			s += fmt.Sprintf("   %s => (", name)
			for i, v := range list {
				s += fmt.Sprintf("%d => signed'(x\"%016x\"), ", i, v)
			}
			s = s[:len(s)-2]
			s += "),\n"
		}
		for name, str := range pkg.Consts.Strings {
			s += fmt.Sprintf("   %s => %q,\n", name, str)
		}
		s = s[:len(s)-2]
		s += "\n);\n"
	}

	return s
}
