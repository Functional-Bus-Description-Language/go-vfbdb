package python

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl"
)

func genConfig(cfg *fbdl.Config, blk *fbdl.Block) string {
	if cfg.IsArray {
		return genConfigArray(cfg, blk)
	} else {
		return genConfigSingle(cfg, blk)
	}
}

func genConfigSingle(cfg *fbdl.Config, blk *fbdl.Block) string {
	var code string

	switch cfg.Access.(type) {
	case fbdl.AccessSingleSingle:
		a := cfg.Access.(fbdl.AccessSingleSingle)
		code += indent + fmt.Sprintf(
			"self.%s = ConfigSingleSingle(iface, %d, (%d, %d))\n",
			cfg.Name, blk.AddrSpace.Start()+a.Addr, a.Mask.Upper, a.Mask.Lower,
		)
	case fbdl.AccessSingleContinuous:
		a := cfg.Access.(fbdl.AccessSingleContinuous)
		decreasigOrder := "False"
		if cfg.HasDecreasingAccessOrder() {
			decreasigOrder = "True"
		}
		code += indent + fmt.Sprintf(
			"self.%s = ConfigSingleContinuous(iface, %d, %d, (%d, %d), (%d, %d), %s)\n",
			cfg.Name,
			blk.AddrSpace.Start()+a.StartAddr(),
			a.RegCount(),
			a.StartMask.Upper, a.StartMask.Lower,
			a.EndMask.Upper, a.EndMask.Lower,
			decreasigOrder,
		)
	default:
		panic("not yet implemented")
	}

	return code
}

func genConfigArray(cfg *fbdl.Config, blk *fbdl.Block) string {
	panic("not yet implemented")
}
