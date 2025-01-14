package python

import (
	_ "embed"
	"log"
	"os"
	"text/template"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/fn"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/pkg"
)

var busWidth int64
var outputPath string

//go:embed templates/vfbdb.py
var pythonTmplStr string
var pythonTmpl = template.Must(template.New("Python module").Parse(pythonTmplStr))

type pythonFormatters struct {
	BusWidth int64
	Code     string
}

func Generate(bus *fn.Block, pkgsConsts map[string]*pkg.Package, cmdLineArgs map[string]string) {
	busWidth = bus.Width
	outputPath = cmdLineArgs["-path"] + "/"

	err := os.MkdirAll(outputPath, os.FileMode(int(0775)))
	if err != nil {
		log.Fatalf("generate Python: %v", err)
	}

	code := genBlock(bus, true)

	code += genPkgConsts(pkgsConsts)

	f, err := os.Create(outputPath + "vfbdb.py")
	if err != nil {
		log.Fatalf("generate Python: %v", err)
	}
	defer f.Close()

	fmts := pythonFormatters{
		BusWidth: busWidth,
		Code:     code,
	}

	err = pythonTmpl.Execute(f, fmts)
	if err != nil {
		log.Fatalf("generate Python: %v", err)
	}
}

var indent string

func increaseIndent(val int) {
	// NOTE: Inefficient implementaion.
	for range val {
		indent += "    "
	}
}

func decreaseIndent(val int) {
	indent = indent[:len(indent)-val*4]
}
