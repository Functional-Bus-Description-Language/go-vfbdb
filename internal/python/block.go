package python

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl"
)

func genBlock(blk *fbdl.Block) string {
	className := "Main"
	if blk.Name != "Main" {
		className = blk.Name + "Class"
	}

	code := indent + fmt.Sprintf("class %s:\n", className)
	increaseIndent(1)

	code += genConsts(blk.ConstContainer)

	code += indent + "def __init__(self, iface):\n"
	increaseIndent(1)
	code += indent + "self.iface = iface\n"

	for _, st := range blk.Statuses {
		code += genStatus(st, blk)
	}

	for _, cfg := range blk.Configs {
		code += genConfig(cfg, blk)
	}

	for _, mask := range blk.Masks {
		code += genMask(mask, blk)
	}

	for _, sb := range blk.Subblocks {
		code += genSubblock(sb, blk)
	}

	decreaseIndent(1)

	for _, fun := range blk.Funcs {
		code += genFunc(fun, blk)
	}

	for _, sb := range blk.Subblocks {
		code += genBlock(sb)
	}

	decreaseIndent(1)

	return code
}

func genSubblock(sb *fbdl.Block, blk *fbdl.Block) string {
	if sb.IsArray {
		return genSublockArray(sb, blk)
	} else {
		return genSublockSingle(sb, blk)
	}
}

func genSublockArray(sb *fbdl.Block, blk *fbdl.Block) string {
	panic("not yet implemented")
}

func genSublockSingle(sb *fbdl.Block, blk *fbdl.Block) string {
	code := indent + fmt.Sprintf("self.%[1]s = self.%[1]sClass(self.iface)\n", sb.Name)

	return code
}
