package vhdl

import (
	_ "embed"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl"
	"log"
	"os"
	"sync"
	"text/template"
)

var busWidth int64
var outputPath string

func Generate(bus *fbdl.Block, cmdLineArgs map[string]string) {
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

	generateWbfbdPackage()

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

//go:embed templates/wbfbd.vhd
var wbfbdPkgStr string
var wbfbdPkgTmpl = template.Must(template.New("VHDL entity").Parse(wbfbdPkgStr))

func generateWbfbdPackage() {
	filePath := outputPath + "wbfbd.vhd"

	f, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("generate VHDL: %v", err)
	}
	defer f.Close()

	err = wbfbdPkgTmpl.Execute(f, nil)
	if err != nil {
		log.Fatalf("generate VHDL: %v", err)
	}

	addGeneratedFile(filePath)
}
