package python

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/access"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/elem"
)

func genConfig(cfg *elem.Config, blk *elem.Block) string {
	if cfg.IsArray {
		return genConfigArray(cfg, blk)
	} else {
		return genConfigSingle(cfg, blk)
	}
}

func genConfigSingle(cfg *elem.Config, blk *elem.Block) string {
	var code string

	switch cfg.Access.(type) {
	case access.SingleSingle:
		a := cfg.Access.(access.SingleSingle)
		code += indent + fmt.Sprintf(
			"self.%s = ConfigSingleSingle(iface, %d, (%d, %d))\n",
			cfg.Name, blk.AddrSpace.Start()+a.Addr, a.EndBit(), a.StartBit(),
		)
	case access.SingleContinuous:
		a := cfg.Access.(access.SingleContinuous)
		code += indent + fmt.Sprintf(
			"self.%s = ConfigSingleContinuous(iface, %d, %d, (%d, %d), (%d, %d))\n",
			cfg.Name,
			blk.AddrSpace.Start()+a.StartAddr(),
			a.RegCount(),
			busWidth-1, a.StartBit(),
			a.EndBit(), 0,
		)
	default:
		panic("not yet implemented")
	}

	return code
}

func genConfigArray(cfg *elem.Config, blk *elem.Block) string {
	panic("not yet implemented")
}
