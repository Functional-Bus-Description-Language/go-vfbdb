package python

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/elem"
)

// Generate block. Main must be true only for the main block.
func genBlock(blk elem.Block, main bool) string {
	className := blk.Name()
	if !main {
		className = blk.Name() + "Class"
	}

	code := indent + fmt.Sprintf("class %s:\n", className)
	increaseIndent(1)

	code += genConsts(blk)

	code += indent + "def __init__(self, iface):\n"
	increaseIndent(1)
	code += indent + "self.iface = iface\n"

	for _, st := range blk.Statuses() {
		code += genStatus(st, blk)
	}

	for _, cfg := range blk.Configs() {
		code += genConfig(cfg, blk)
	}

	for _, mask := range blk.Masks() {
		code += genMask(mask, blk)
	}

	for _, sb := range blk.Subblocks() {
		code += genSubblock(sb, blk)
	}

	for _, fun := range blk.Funcs() {
		code += genFunc(fun, blk)
	}

	for _, stream := range blk.Streams() {
		code += genStream(stream, blk)
	}

	decreaseIndent(1)

	for _, sb := range blk.Subblocks() {
		code += genBlock(sb, false)
	}

	decreaseIndent(1)

	return code
}

func genSubblock(sb elem.Block, blk elem.Block) string {
	if sb.IsArray() {
		return genSublockArray(sb, blk)
	} else {
		return genSublockSingle(sb, blk)
	}
}

func genSublockArray(sb elem.Block, blk elem.Block) string {
	panic("not yet implemented")
}

func genSublockSingle(sb elem.Block, blk elem.Block) string {
	code := indent + fmt.Sprintf("self.%[1]s = self.%[1]sClass(self.iface)\n", sb.Name())

	return code
}
