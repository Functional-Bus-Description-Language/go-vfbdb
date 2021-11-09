package vhdl

import (
	_ "embed"
	"fmt"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl"
	"log"
	"os"
	_ "strings"
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

	entities := collectEntities(bus, nil, []string{})

	//resolveEntityNameConflicts(entities)

	var wg sync.WaitGroup
	defer wg.Wait()

	generateWbfbdPackage()

	for _, ent := range entities {
		wg.Add(1)
		go generateEntity(ent, &wg)
	}
}

func collectEntities(block *fbdl.Block, entities []Entity, path []string) []Entity {
	if entities == nil {
		entities = []Entity{Entity{Name: "main", Path: []string{"main"}, Block: block}}
		path = append(path, "main")
	} else {
		p := make([]string, len(path))
		n := copy(p, path)
		if n != len(path) {
			log.Fatalf("generate vhdl: copying entity path failed, copied %d, expected %d", n, len(path))
		}

		ent := Entity{Name: block.Name, Path: p, Block: block}
		entities = append(entities, ent)
	}

	fmt.Println(path)
	for _, b := range block.Subblocks {
		path = append(path, b.Name)
		entities = collectEntities(b, entities, path)
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
