package python

import (
	"fmt"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl"
)

func generateConfig(cfg *fbdl.Config, blk *fbdl.Block) string {
	if cfg.IsArray {
		return generateConfigArray(cfg, blk)
	} else {
		return generateConfigSingle(cfg, blk)
	}
}

func generateConfigSingle(cfg *fbdl.Config, blk *fbdl.Block) string {
	var code string

	switch cfg.Access.(type) {
	case fbdl.AccessSingleSingle:
		access := cfg.Access.(fbdl.AccessSingleSingle)
		code += indent + fmt.Sprintf(
			"self.%s = ConfigSingleSingle(interface, %d, (%d, %d))\n",
			cfg.Name, blk.AddrSpace.Start()+access.Addr, access.Mask.Upper, access.Mask.Lower,
		)
	default:
		panic("not yet implemented")
	}

	return code
}

func generateConfigArray(cfg *fbdl.Config, blk *fbdl.Block) string {
	panic("not yet implemented")
}
