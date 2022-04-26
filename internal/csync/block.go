package csync

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/Functional-Bus-Description-Language/go-wbfbd/internal/utils"

	"text/template"
)

//go:embed templates/block.h
var blockHeaderTmplStr string
var blockHeaderTmpl = template.Must(template.New("C-Sync header").Parse(blockHeaderTmplStr))

//go:embed templates/block.c
var blockSourceTmplStr string
var blockSourceTmpl = template.Must(template.New("C-Sync source").Parse(blockSourceTmplStr))

type BlockHeaderFormatters struct {
	BlockName string
	Code      string
}

type BlockSourceFormatters struct {
	Code string
}

func genBlock(b utils.Block, wg *sync.WaitGroup) {
	defer wg.Done()

	hFmts := BlockHeaderFormatters{
		BlockName: b.Name,
		Code:      "",
	}
	srcFmts := BlockSourceFormatters{Code: ""}

	for _, st := range b.Block.Statuses {
		genStatus(st, &hFmts, &srcFmts)
	}

	genBlockHeader(b, hFmts)
	genBlockSource(b, srcFmts)
}

func genBlockHeader(b utils.Block, hFmts BlockHeaderFormatters) {
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

func genBlockSource(b utils.Block, srcFmts BlockSourceFormatters) {
	f, err := os.Create(outputPath + fmt.Sprintf("%s.c", b.Name))
	if err != nil {
		log.Fatalf("generate C-Sync: %v", err)
	}
	defer f.Close()

	err = blockSourceTmpl.Execute(f, srcFmts)
	if err != nil {
		log.Fatalf("generate C-Sync: %v", err)
	}
}
