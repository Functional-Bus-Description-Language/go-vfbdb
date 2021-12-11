package python

import (
	"fmt"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl"
)

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
