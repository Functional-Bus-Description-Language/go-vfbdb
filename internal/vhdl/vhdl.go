package vhdl

import (
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl"
	"log"
	"os"
	"sync"
)

var busWidth int64
var outputPath string

func Generate(bus *fbdl.Block, pkgsConsts map[string]fbdl.Package, cmdLineArgs map[string]string) {
	busWidth = bus.Width
	outputPath = cmdLineArgs["--path"] + "/"

	err := os.MkdirAll(outputPath, os.FileMode(int(0775)))
	if err != nil {
		log.Fatalf("generate vhdl: %v", err)
	}

	blockEntities := collectBlockEntities(bus, nil, []string{})

	//resolveEntityNameConflicts(entities)

	var wg sync.WaitGroup
	defer wg.Wait()

	generateWbfbdPackage(pkgsConsts)

	for _, be := range blockEntities {
		wg.Add(1)
		go generateBlock(be, &wg)
	}
}

func collectBlockEntities(blk *fbdl.Block, entities []BlockEntity, path []string) []BlockEntity {
	if entities == nil {
		entities = []BlockEntity{BlockEntity{Name: "main", Path: []string{"main"}, Block: blk}}
		path = append(path, "main")
	} else {
		p := make([]string, len(path))
		n := copy(p, path)
		if n != len(path) {
			log.Fatalf("generate vhdl: copying entity path failed, copied %d, expected %d", n, len(path))
		}

		ent := BlockEntity{Name: blk.Name, Path: p, Block: blk}
		entities = append(entities, ent)
	}

	for _, b := range blk.Subblocks {
		path = append(path, b.Name)
		entities = collectBlockEntities(b, entities, path)
		path = path[:len(path)-1]
	}

	return entities
}
