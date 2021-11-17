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

func generateClass(block *fbdl.Block) string {
	code := fmt.Sprintf("class %s:\n", block.Name)
	increaseIndent(1)
	code += indent + "def __init__(self, interface):\n"
	increaseIndent(1)

	code += generateStatuses(block)

	return code
}

func generateStatuses(block *fbdl.Block) string {
	var code string

	for _, st := range block.Statuses {
		if st.IsArray {

		} else {
			switch st.Access.(type) {
			case fbdl.AccessSingleSingle:
				access := st.Access.(fbdl.AccessSingleSingle)
				code += indent + fmt.Sprintf(
					"self.%s = StatusSingleSingle(interface, %d, (%d, %d))\n",
					st.Name, block.AddrSpace.Start()+access.Addr, access.Mask.Upper, access.Mask.Lower,
				)
			default:
				panic("not yet implemented")
			}
		}
	}

	return code
}
