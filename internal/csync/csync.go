package csync

import (
	_ "embed"
	"fmt"
	"log"
	"math"
	"os"
	"text/template"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl"
	"github.com/Functional-Bus-Description-Language/go-vfbdb/internal/c"
	"github.com/Functional-Bus-Description-Language/go-vfbdb/internal/utils"
	"strconv"
	"sync"
)

var busWidth int64
var outputPath string

var addrType c.Type
var readDataType c.Type
var writeDataType c.Type

//go:embed templates/vfbdb.h
var vfbdbHeaderTmplStr string
var vfbdbHeaderTmpl = template.Must(template.New("C-Sync vfbdb.h").Parse(vfbdbHeaderTmplStr))

//go:embed templates/vfbdb.c
var vfbdbSourceTmplStr string
var vfbdbSourceTmpl = template.Must(template.New("C-Sync vfbdb.c").Parse(vfbdbSourceTmplStr))

type vfbdbHeaderFormatters struct {
	AddrType      string
	ReadDataType  string
	WriteDataType string
}

type vfbdbSourceFormatters struct {
	ID        string
	TIMESTAMP string
}

func Generate(bus *fbdl.Block, pkgsConsts map[string]fbdl.Package, cmdLineArgs map[string]string) {
	busWidth = bus.Width
	outputPath = cmdLineArgs["-path"] + "/"

	err := os.MkdirAll(outputPath, os.FileMode(int(0775)))
	if err != nil {
		log.Fatalf("generate C-Sync: %v", err)
	}

	hFile, err := os.Create(outputPath + "vfbdb.h")
	if err != nil {
		log.Fatalf("generate C-Sync: %v", err)
	}
	defer hFile.Close()

	addrType = c.WidthToWriteType(
		int64(math.Log2(float64(bus.Sizes.BlockAligned))),
	)
	readDataType = c.WidthToReadType(bus.Width)
	writeDataType = c.WidthToWriteType(bus.Width)

	hFmts := vfbdbHeaderFormatters{
		AddrType:      addrType.String(),
		ReadDataType:  readDataType.String(),
		WriteDataType: writeDataType.String(),
	}

	err = vfbdbHeaderTmpl.Execute(hFile, hFmts)
	if err != nil {
		log.Fatalf("generate C-Sync: %v", err)
	}

	srcFile, err := os.Create(outputPath + "vfbdb.c")
	if err != nil {
		log.Fatalf("generate C-Sync: %v", err)
	}
	defer srcFile.Close()

	srcFmts := vfbdbSourceFormatters{
		ID:        fmt.Sprintf("0x%s", strconv.FormatUint(bus.Status("ID").Default.Uint64(), 16)),
		TIMESTAMP: fmt.Sprintf("0x%s", strconv.FormatUint(bus.Status("TIMESTAMP").Default.Uint64(), 16)),
	}

	err = vfbdbSourceTmpl.Execute(srcFile, srcFmts)
	if err != nil {
		log.Fatalf("generate C-Sync: %v", err)
	}

	blocks := utils.CollectBlocks(bus, nil, []string{})
	utils.ResolveBlockNameConflicts(blocks)

	var wg sync.WaitGroup
	defer wg.Wait()

	for _, b := range blocks {
		wg.Add(1)
		go genBlock(b, &wg)
	}
}
