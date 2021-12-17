package python

import (
	"fmt"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl"
)

func generateBlock(blk *fbdl.Block) string {
	className := "main"
	if blk.Name != "main" {
		className = blk.Name + "Class"
	}

	code := indent + fmt.Sprintf("class %s:\n", className)
	increaseIndent(1)

	code += generateConsts(blk)

	code += indent + "def __init__(self, interface):\n"
	increaseIndent(1)
	code += indent + "self.interface = interface\n"

	for _, st := range blk.Statuses {
		code += generateStatus(st, blk)
	}

	for _, cfg := range blk.Configs {
		code += generateConfig(cfg, blk)
	}

	for _, mask := range blk.Masks {
		code += generateMask(mask, blk)
	}

	for _, sb := range blk.Subblocks {
		code += generateSubblock(sb, blk)
	}

	decreaseIndent(1)

	for _, fun := range blk.Funcs {
		code += generateFunc(fun, blk)
	}

	for _, sb := range blk.Subblocks {
		code += generateBlock(sb)
	}

	decreaseIndent(1)

	return code
}

func generateSubblock(sb *fbdl.Block, blk *fbdl.Block) string {
	if sb.IsArray {
		return generateSublockArray(sb, blk)
	} else {
		return generateSublockSingle(sb, blk)
	}
}

func generateSublockArray(sb *fbdl.Block, blk *fbdl.Block) string {
	panic("not yet implemented")
}

func generateSublockSingle(sb *fbdl.Block, blk *fbdl.Block) string {
	code := indent + fmt.Sprintf("self.%[1]s = self.%[1]sClass(self.interface)\n", sb.Name)

	return code
}

func generateConsts(blk *fbdl.Block) string {
	code := ""

	for name, i := range blk.IntConsts {
		code += indent + fmt.Sprintf("%s = %d\n", name, i)
	}
	for name, str := range blk.StrConsts {
		code += indent + fmt.Sprintf("%s = %q\n", name, str)
	}

	return code
}
