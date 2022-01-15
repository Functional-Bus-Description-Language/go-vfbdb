package vhdl

import (
	_ "embed"
	"fmt"
	"log"
	"math"
	"os"
	"sync"
	"text/template"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl"
)

type BlockEntity struct {
	Name      string
	NameLevel int
	Path      []string
	Block     *fbdl.Block
}

// Rename renames BlockEntity based on the NameLevel and Path.
// NameLevel indicates how many Path elements should be included in the Name.
// NameLevel field is incremented internally. Path elements are taken from the end.
func (be *BlockEntity) Rename() {
	if be.NameLevel < len(be.Path) {
		be.NameLevel += 1
	}

	name := be.Path[len(be.Path)-1]
	for i := 1; i < be.NameLevel; i++ {
		name = be.Path[len(be.Path)-i-1] + "_" + name
	}

	be.Name = name
}

//go:embed templates/blockEntity.vhd
var blockEntityTmplStr string

var blockEntityTmpl = template.Must(template.New("VHDL entity").Parse(blockEntityTmplStr))

type BlockEntityFormatters struct {
	BusWidth   int64
	EntityName string

	MastersCount          int64
	RegistersCount        int64
	InternalAddrBitsCount int64
	SubblocksCount        int64

	// Things going to package.
	Constants string
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

func generateBlock(be BlockEntity, wg *sync.WaitGroup) {
	defer wg.Done()

	fmts := BlockEntityFormatters{
		BusWidth:              busWidth,
		EntityName:            be.Name,
		MastersCount:          be.Block.Masters,
		RegistersCount:        be.Block.Sizes.Own,
		InternalAddrBitsCount: int64(math.Ceil(math.Log2(float64(be.Block.Sizes.Own)))),
		AddressValues:         fmt.Sprintf("0 => \"%032b\"", 0),
		RegistersAccess:       make(RegisterMap),
	}

	addrBitsCount := int(math.Log2(float64(be.Block.Sizes.BlockAligned)))

	mask := 0
	if len(be.Block.Subblocks) > 0 {
		mask = ((1 << addrBitsCount) - 1) ^ ((1 << fmts.InternalAddrBitsCount) - 1)
	}
	fmts.MaskValues = fmt.Sprintf("0 => \"%032b\"", mask)

	generateConsts(be.Block, &fmts)

	for _, sb := range be.Block.Subblocks {
		generateSubblock(sb, be.Block.AddrSpace.Start(), addrBitsCount, &fmts)
	}

	for _, fun := range be.Block.Funcs {
		generateFunc(fun, &fmts)
	}

	for _, st := range be.Block.Statuses {
		generateStatus(st, &fmts)
	}

	for _, cfg := range be.Block.Configs {
		generateConfig(cfg, &fmts)
	}

	for _, mask := range be.Block.Masks {
		generateMask(mask, &fmts)
	}

	filePath := outputPath + be.Name + ".vhd"
	f, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("generate VHDL: %v", err)
	}
	defer f.Close()

	err = blockEntityTmpl.Execute(f, fmts)
	if err != nil {
		log.Fatalf("generate VHDL: %v", err)
	}

	addGeneratedFile(filePath)
}

func generateSubblock(
	sb *fbdl.Block,
	superBlockAddrStart int64,
	superBlockAddrBitsCount int,
	fmts *BlockEntityFormatters,
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
		lowerBound := initSubblocksCount + 1
		upperBound := lowerBound + sb.Count - 1

		s := fmt.Sprintf("\n   master_i(%d downto %d) => %s_master_i,", lowerBound, upperBound, sb.Name)
		fmts.CrossbarSubblockPortsIn += s

		s = fmt.Sprintf(",\n   master_o(%d downto %d) => %s_master_o", lowerBound, upperBound, sb.Name)
		fmts.CrossbarSubblockPortsOut += s
	}

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

func generateConsts(blk *fbdl.Block, fmts *BlockEntityFormatters) {
	s := ""

	for name, i := range blk.IntConsts {
		s += fmt.Sprintf("constant %s : int64 := signed'(x\"%016x\");\n", name, i)
	}
	for name, list := range blk.IntListConsts {
		s += fmt.Sprintf("constant %s : int64_vector(0 to %d) := (", name, len(list)-1)
		for i, v := range list {
			s += fmt.Sprintf("%d => signed'(x\"%016x\"), ", i, v)
		}
		s = s[:len(s)-2]
		s += ");\n"
	}
	for name, str := range blk.StrConsts {
		s += fmt.Sprintf("constant %s : string := %q;\n", name, str)
	}

	fmts.Constants += s
}
