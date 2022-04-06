package python

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl"
	"strconv"
)

var busWidth int64
var outputPath string

//go:embed templates/wbfbd.py
var pythonTmplStr string

var pythonTmpl = template.Must(template.New("Python module").Parse(pythonTmplStr))

type pythonFormatters struct {
	BusWidth  int64
	ID        string
	TIMESTAMP string
	Code      string
}

func Generate(bus *fbdl.Block, pkgsConsts map[string]fbdl.Package, cmdLineArgs map[string]string) {
	busWidth = bus.Width
	outputPath = cmdLineArgs["--path"] + "/"

	err := os.MkdirAll(outputPath, os.FileMode(int(0775)))
	if err != nil {
		log.Fatalf("generate Python: %v", err)
	}

	code := generateBlock(bus)

	code += generatePkgConsts(pkgsConsts)

	f, err := os.Create(outputPath + "wbfbd.py")
	if err != nil {
		log.Fatalf("generate Python: %v", err)
	}
	defer f.Close()

	fmts := pythonFormatters{
		BusWidth:  busWidth,
		ID:        fmt.Sprintf("0x%s", strconv.FormatUint(bus.Status("ID").Default.Uint64(), 16)),
		TIMESTAMP: fmt.Sprintf("0x%s", strconv.FormatUint(bus.Status("TIMESTAMP").Default.Uint64(), 16)),
		Code:      code,
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
