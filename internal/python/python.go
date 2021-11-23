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
	code := fmt.Sprintf("class %s:\n", blk.Name)
	increaseIndent(1)
	code += indent + "def __init__(self, interface):\n"
	increaseIndent(1)
	code += indent + "self.interface = interface\n"

	for _, st := range blk.Statuses {
		code += generateStatus(blk, st)
	}
	decreaseIndent(1)

	for _, fun := range blk.Funcs {
		code += generateFunc(blk, fun)
	}

	return code
}

func generateStatus(blk *fbdl.Block, st *fbdl.Status) string {
	var code string

	if st.IsArray {
		switch st.Access.(type) {
		case fbdl.AccessArrayMultiple:
			access := st.Access.(fbdl.AccessArrayMultiple)
			code += indent + fmt.Sprintf(
				"self.%s = StatusArrayMultiple(interface, %d, %d, %d)\n",
				st.Name, blk.AddrSpace.Start()+access.StartAddr(), access.ItemWidth, access.ItemCount,
			)
		default:
			panic("not yet implemented")
		}
	} else {
		switch st.Access.(type) {
		case fbdl.AccessSingleSingle:
			access := st.Access.(fbdl.AccessSingleSingle)
			code += indent + fmt.Sprintf(
				"self.%s = StatusSingleSingle(interface, %d, (%d, %d))\n",
				st.Name, blk.AddrSpace.Start()+access.Addr, access.Mask.Upper, access.Mask.Lower,
			)
		default:
			panic("not yet implemented")
		}
	}

	return code
}
