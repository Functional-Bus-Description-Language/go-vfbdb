package csync

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/Functional-Bus-Description-Language/go-vfbdb/internal/utils"

	"text/template"
)

//go:embed templates/block.h
var blockHeaderTmplStr string
var blockHeaderTmpl = template.Must(template.New("C-Sync header").Parse(blockHeaderTmplStr))

//go:embed templates/block.c
var blockSourceTmplStr string
var blockSourceTmpl = template.Must(template.New("C-Sync source").Parse(blockSourceTmplStr))

type BlockHFormatters struct {
	BlockName string
	Code      string
}

type BlockCFormatters struct {
	Code string
}

func genBlock(b utils.Block, wg *sync.WaitGroup) {
	defer wg.Done()

	hFmts := BlockHFormatters{
		BlockName: b.Name,
		Code:      "",
	}
	cFmts := BlockCFormatters{Code: ""}

	for _, st := range b.Block.Statics {
		genStatic(st, b.Block, &hFmts, &cFmts)
	}

	for _, st := range b.Block.Statuses {
		genStatus(st, b.Block, &hFmts, &cFmts)
	}

	for _, cfg := range b.Block.Configs {
		genConfig(cfg, b.Block, &hFmts, &cFmts)
	}

	for _, proc := range b.Block.Procs {
		genProc(proc, b.Block, &hFmts, &cFmts)
	}

	genBlockH(b, hFmts)
	genBlockC(b, cFmts)
}

func genBlockH(b utils.Block, hFmts BlockHFormatters) {
	f, err := os.Create(outputPath + fmt.Sprintf("%s.h", b.Name))
	if err != nil {
		log.Fatalf("generate C-Sync: %v", err)
	}
	defer f.Close()

	err = blockHeaderTmpl.Execute(f, hFmts)
	if err != nil {
		log.Fatalf("generate C-Sync: %v", err)
	}
}

func genBlockC(b utils.Block, cFmts BlockCFormatters) {
	f, err := os.Create(outputPath + fmt.Sprintf("%s.c", b.Name))
	if err != nil {
		log.Fatalf("generate C-Sync: %v", err)
	}
	defer f.Close()

	err = blockSourceTmpl.Execute(f, cFmts)
	if err != nil {
		log.Fatalf("generate C-Sync: %v", err)
	}
}
