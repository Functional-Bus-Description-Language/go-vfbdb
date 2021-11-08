package vhdl

import (
	_ "embed"
	"fmt"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl"
	"log"
	"math"
	"os"
	"sync"
	"text/template"
)

type Entity struct {
	Name  string
	Path  []string
	Block *fbdl.Block
}

//go:embed templates/block.vhd
var entityTmplStr string

var entityTmpl = template.Must(template.New("VHDL entity").Parse(entityTmplStr))

type EntityFormatters struct {
	BusWidth   int64
	EntityName string

	MastersCount          int64
	RegistersCount        int64
	InternalAddrBitsCount int64
	SubblocksCount        int64

	// Constants string // TODO: Decide how to implement this.

	EntitySubblockPorts   string
	EntityFunctionalPorts string
	CrossbarSubblockPorts string
	SignalDeclarations    string
	AddressValues         string
	MaskValues            string
	StatusesAccess        string
	StatusesRouting       string
	DefaultValues         string
}

func generateEntity(entity Entity, wg *sync.WaitGroup) {
	defer wg.Done()

	fmts := EntityFormatters{
		BusWidth:              busWidth,
		EntityName:            entity.Name,
		MastersCount:          entity.Block.Masters,
		RegistersCount:        entity.Block.Sizes.Own,
		InternalAddrBitsCount: int64(math.Ceil(math.Log2(float64(entity.Block.Sizes.Own)))),
		AddressValues:         fmt.Sprintf("0 => \"%032b\"", 0),
	}

	addrBitsCount := int(math.Log2(float64(entity.Block.Sizes.BlockAligned)))

	mask := 0
	if len(entity.Block.Subblocks) > 0 {
		mask = ((1 << addrBitsCount) - 1) ^ ((1 << fmts.InternalAddrBitsCount) - 1)
	}
	fmts.MaskValues = fmt.Sprintf("0 => \"%032b\"", mask)

	currentSubblockAddr := entity.Block.Sizes.BlockAligned
	for _, sb := range entity.Block.Subblocks {
		currentSubblockAddr = generateSubblock(&sb, addrBitsCount, currentSubblockAddr, fmts)
	}

	f, err := os.Create(outputPath + entity.Name + ".vhd")
	if err != nil {
		log.Fatalf("generate VHDL: %v", err)
	}
	defer f.Close()

	err = entityTmpl.Execute(f, fmts)
	if err != nil {
		log.Fatalf("generate VHDL: %v", err)
	}
}

func generateSubblock(
	sb *fbdl.Block,
	superBlockAddrBitsCount int,
	currentSubblockAddr int64,
	fmts EntityFormatters,
) int64 {
	initSubblocksCount := fmts.SubblocksCount

	s := fmt.Sprintf(
		";\n      %s_master_o : out t_wishbone_master_out_array(%d - 1 downto 0);\n"+
			"      %[1]s_master_i : in  t_wishbone_master_in_array(%[2]d - 1 downto 0)",
		sb.Name, sb.Count,
	)
	fmts.EntitySubblockPorts += s

	if sb.Count == 1 {
		s := fmt.Sprintf(
			",\n      master_i(%d + 1) => %s_master_i,\n"+
				"      master_o(%[1]d + 1) => %[2]s_master_o",
			initSubblocksCount, sb.Name,
		)
		fmts.CrossbarSubblockPorts += s
	} else {
		lower_bound := initSubblocksCount + 1
		upper_bound := lower_bound + sb.Count - 1

		s := fmt.Sprintf(
			",\n      master_i(%d downto %d) => %s_master_i,\n"+
				"      master_o(%[1]d downto %[2]d) => %[3]s_master_o",
			lower_bound, upper_bound, sb.Name,
		)
		fmts.CrossbarSubblockPorts += s
	}

	for i := int64(0); i < sb.Count; i++ {
		fmts.SubblocksCount += 1

		currentSubblockAddr -= sb.Sizes.BlockAligned

		s := fmt.Sprintf(", %d => \"%032b\"", fmts.SubblocksCount, currentSubblockAddr)
		fmts.AddressValues += s

		mask := ((1 << superBlockAddrBitsCount) - 1) ^ ((1 << int(math.Ceil(math.Log2(float64(sb.Sizes.Own))))) - 1)
		s = fmt.Sprintf(", %d => \"%032b\"", fmts.SubblocksCount, mask)
		fmts.MaskValues += s
	}

	return currentSubblockAddr
}
