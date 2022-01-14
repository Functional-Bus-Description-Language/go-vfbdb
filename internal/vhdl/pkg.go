package vhdl

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl"
)

//go:embed templates/wbfbd.vhd
var wbfbdPkgStr string
var wbfbdPkgTmpl = template.Must(template.New("VHDL wbfbd package").Parse(wbfbdPkgStr))

type wbfbdPackageFormatters struct {
	PkgsConsts string
}

func generateWbfbdPackage(pkgsConsts map[string]fbdl.Package) {
	filePath := outputPath + "wbfbd.vhd"

	f, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("generate VHDL: %v", err)
	}
	defer f.Close()

	fmts := wbfbdPackageFormatters{PkgsConsts: generatePkgsConsts(pkgsConsts)}

	err = wbfbdPkgTmpl.Execute(f, fmts)
	if err != nil {
		log.Fatalf("generate VHDL: %v", err)
	}

	addGeneratedFile(filePath)
}

func generatePkgsConsts(pkgsConsts map[string]fbdl.Package) string {
	s := ""

	for pkgName, pkg := range pkgsConsts {
		if !pkg.HasConsts() {
			continue
		}

		// Package type definition
		s += fmt.Sprintf("type t_%s_pkg is record\n", pkgName)
		for name, _ := range pkg.IntConsts {
			s += fmt.Sprintf("   %s : int64;\n", name)
		}
		for name, _ := range pkg.StrConsts {
			s += fmt.Sprintf("   %s : string;\n", name)
		}
		s += fmt.Sprintf("end record;\n")

		// Package constant definition
		s += fmt.Sprintf("constant %[1]s_pkg : t_%[1]s_pkg := (\n", pkgName)
		for name, i := range pkg.IntConsts {
			s += fmt.Sprintf("   %s => signed'(x\"%016x\"),\n", name, i)
		}
		for name, str := range pkg.StrConsts {
			s += fmt.Sprintf("   %s => %q,\n", name, str)
		}
		s = s[:len(s)-2]
		s += fmt.Sprintf("\n);\n")
	}

	return s
}
