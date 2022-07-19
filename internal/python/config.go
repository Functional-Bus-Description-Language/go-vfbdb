package python

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/access"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/elem"
)

func genConfig(cfg elem.Config, blk elem.Block) string {
	if cfg.IsArray() {
		return genConfigArray(cfg, blk)
	} else {
		return genConfigSingle(cfg, blk)
	}
}

func genConfigSingle(cfg elem.Config, blk elem.Block) string {
	var code string

	switch cfg.Access().(type) {
	case access.SingleSingle:
		a := cfg.Access().(access.SingleSingle)
		code += indent + fmt.Sprintf(
			"self.%s = ConfigSingleSingle(iface, %d, (%d, %d))\n",
			cfg.Name(), blk.AddrSpace().Start()+a.Addr, a.Mask.Upper, a.Mask.Lower,
		)
	case access.SingleContinuous:
		a := cfg.Access().(access.SingleContinuous)
		decreasigOrder := "False"
		if cfg.HasDecreasingAccessOrder() {
			decreasigOrder = "True"
		}
		code += indent + fmt.Sprintf(
			"self.%s = ConfigSingleContinuous(iface, %d, %d, (%d, %d), (%d, %d), %s)\n",
			cfg.Name(),
			blk.AddrSpace().Start()+a.StartAddr(),
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

func genConfigArray(cfg elem.Config, blk elem.Block) string {
	panic("not yet implemented")
}
