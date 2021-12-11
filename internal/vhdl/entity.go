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

	// Things going to package.
	// Constants string // TODO: Decide how to implement this.
	FuncTypes string

	EntitySubblockPorts   string
	EntityFunctionalPorts string

	CrossbarSubblockPortsIn  string
	CrossbarSubblockPortsOut string

	SignalDeclarations string
	AddressValues      string
	MaskValues         string

	RegistersAccess RegisterMap

	FuncsStrobesClear string
	FuncsStrobesSet   string

	DefaultValues string
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
		RegistersAccess:       make(RegisterMap),
	}

	addrBitsCount := int(math.Log2(float64(entity.Block.Sizes.BlockAligned)))

	mask := 0
	if len(entity.Block.Subblocks) > 0 {
		mask = ((1 << addrBitsCount) - 1) ^ ((1 << fmts.InternalAddrBitsCount) - 1)
	}
	fmts.MaskValues = fmt.Sprintf("0 => \"%032b\"", mask)

	for _, sb := range entity.Block.Subblocks {
		generateSubblock(sb, entity.Block.AddrSpace.Start(), addrBitsCount, &fmts)
	}

	for _, fun := range entity.Block.Funcs {
		generateFunc(fun, &fmts)
	}

	for _, st := range entity.Block.Statuses {
		generateStatus(st, &fmts)
	}

	for _, cfg := range entity.Block.Configs {
		generateConfig(cfg, &fmts)
	}

	filePath := outputPath + entity.Name + ".vhd"
	f, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("generate VHDL: %v", err)
	}
	defer f.Close()

	err = entityTmpl.Execute(f, fmts)
	if err != nil {
		log.Fatalf("generate VHDL: %v", err)
	}

	addGeneratedFile(filePath)
}

func generateSubblock(
	sb *fbdl.Block,
	superBlockAddrStart int64,
	superBlockAddrBitsCount int,
	fmts *EntityFormatters,
) {
	initSubblocksCount := fmts.SubblocksCount

	s := fmt.Sprintf(
		";\n   %s_master_o : out t_wishbone_master_out_array(%d downto 0);\n"+
			"   %[1]s_master_i : in  t_wishbone_master_in_array(%[2]d downto 0)",
		sb.Name, sb.Count-1,
	)
	fmts.EntitySubblockPorts += s

	if sb.Count == 1 {
		s := fmt.Sprintf("\n   master_i(%d) => %s_master_i(0),", initSubblocksCount+1, sb.Name)
		fmts.CrossbarSubblockPortsIn += s

		s = fmt.Sprintf(",\n   master_o(%d) => %s_master_o(0)", initSubblocksCount+1, sb.Name)
		fmts.CrossbarSubblockPortsOut += s
	} else {
		lower_bound := initSubblocksCount + 1
		upper_bound := lower_bound + sb.Count - 1

		s := fmt.Sprintf("\n   master_i(%d downto %d) => %s_master_i,", lower_bound, upper_bound, sb.Name)
		fmts.CrossbarSubblockPortsIn += s

		s = fmt.Sprintf(",\n   master_o(%d downto %d) => %s_master_o", lower_bound, upper_bound, sb.Name)
		fmts.CrossbarSubblockPortsOut += s
	}

	fmt.Printf("%s: %s\n", fmts.EntityName, sb.Name)

	subblockAddr := sb.AddrSpace.Start() - superBlockAddrStart
	for i := int64(0); i < sb.Count; i++ {
		fmts.SubblocksCount += 1

		s := fmt.Sprintf(", %d => \"%032b\"", fmts.SubblocksCount, subblockAddr)
		fmts.AddressValues += s

		mask := ((1 << superBlockAddrBitsCount) - 1) ^ ((1 << int(math.Log2(float64(sb.Sizes.BlockAligned)))) - 1)
		s = fmt.Sprintf(", %d => \"%032b\"", fmts.SubblocksCount, mask)
		fmts.MaskValues += s

		subblockAddr += sb.Sizes.BlockAligned
	}
}
