package csync

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl"
	"github.com/Functional-Bus-Description-Language/go-wbfbd/internal/utils"
	"strconv"
)

var busWidth int64
var outputPath string

// addrSize and dataSize are used only when address or data are too wide
// to be represented as basic type, and they must be represented as array.
// NOTE: Handling such cases is not yet implemented.
var addrSize int // Size of address array
var dataSize int // Size of data array

//go:embed templates/wbfbd.h
var wbfbdHeaderTmplStr string
var wbfbdHeaderTmpl = template.Must(template.New("C-Sync wbfbd.h").Parse(wbfbdHeaderTmplStr))

type wbfbdHeaderFormatters struct {
	AddrType      string
	ReadDataType  string
	WriteDataType string
	ID            string
	TIMESTAMP     string
}

func Generate(bus *fbdl.Block, pkgsConsts map[string]fbdl.Package, cmdLineArgs map[string]string) {
	busWidth = bus.Width
	outputPath = cmdLineArgs["-path"] + "/"

	err := os.MkdirAll(outputPath, os.FileMode(int(0775)))
	if err != nil {
		log.Fatalf("generate C-Sync: %v", err)
	}

	f, err := os.Create(outputPath + "wbfbd.h")
	if err != nil {
		log.Fatalf("generate C-Sync: %v", err)
	}
	defer f.Close()

	fmts := wbfbdHeaderFormatters{
		AddrType:      utils.WidthToCTypeWrite(busWidth),
		ReadDataType:  utils.WidthToCTypeRead(bus.Width),
		WriteDataType: utils.WidthToCTypeWrite(bus.Width),
		ID:            fmt.Sprintf("0x%s", strconv.FormatUint(bus.Status("ID").Default.Uint64(), 16)),
		TIMESTAMP:     fmt.Sprintf("0x%s", strconv.FormatUint(bus.Status("TIMESTAMP").Default.Uint64(), 16)),
	}

	err = wbfbdHeaderTmpl.Execute(f, fmts)
	if err != nil {
		log.Fatalf("generate C-Sync: %v", err)
	}
}
