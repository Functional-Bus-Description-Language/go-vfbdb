package vhdlwb3

import (
	"log"
	"os"
	"sync"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/elem"
	"github.com/Functional-Bus-Description-Language/go-vfbdb/internal/utils"
)

var busWidth int64
var outputPath string

func Generate(bus *elem.Block, pkgsConsts map[string]*elem.Package, cmdLineArgs map[string]string) {
	busWidth = bus.Width
	outputPath = cmdLineArgs["-path"] + "/"

	err := os.MkdirAll(outputPath, os.FileMode(int(0775)))
	if err != nil {
		log.Fatalf("generate vhdl-wb3: %v", err)
	}

	blocks := utils.CollectBlocks(bus, nil, []string{})
	utils.ResolveBlockNameConflicts(blocks)

	var wg sync.WaitGroup
	defer wg.Wait()

	genWb3Package(pkgsConsts)

	for _, b := range blocks {
		wg.Add(1)
		go genBlock(b, &wg)
	}
}
