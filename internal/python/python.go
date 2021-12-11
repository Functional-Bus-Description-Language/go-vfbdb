package python

import (
	_ "embed"
	"fmt"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl"
	"log"
	"os"
	_ "strings"
	"text/template"
)

var busWidth int64
var outputPath string

//go:embed templates/wbfbd.py
var pythonTmplStr string

var pythonTmpl = template.Must(template.New("Python module").Parse(pythonTmplStr))

type pythonFormatters struct {
	BusWidth int64
	Code     string
}

func Generate(bus *fbdl.Block, cmdLineArgs map[string]string) {
	busWidth = bus.Width
	outputPath = cmdLineArgs["--path"] + "/"

	err := os.MkdirAll(outputPath, os.FileMode(int(0775)))
	if err != nil {
		log.Fatalf("generate vhdl: %v", err)
	}

	code := generateClass(bus)

	f, err := os.Create(outputPath + "wbfbd.py")
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
	for i := 0; i < val; i++ {
		indent += "    "
	}
}

func decreaseIndent(val int) {
	indent = indent[:len(indent)-val*4]
}

func generateClass(blk *fbdl.Block) string {
	className := "main"
	if blk.Name != "main" {
		className = blk.Name + "Class"
	}

	code := indent + fmt.Sprintf("class %s:\n", className)
	increaseIndent(1)
	code += indent + "def __init__(self, interface):\n"
	increaseIndent(1)
	code += indent + "self.interface = interface\n"

	for _, st := range blk.Statuses {
		code += generateStatus(st, blk)
	}

	for _, cfg := range blk.Configs {
		code += generateConfig(cfg, blk)
	}

	for _, sb := range blk.Subblocks {
		code += generateSubblock(sb, blk)
	}

	decreaseIndent(1)

	for _, fun := range blk.Funcs {
		code += generateFunc(fun, blk)
	}

	for _, sb := range blk.Subblocks {
		code += generateClass(sb)
		decreaseIndent(1)
	}

	return code
}
