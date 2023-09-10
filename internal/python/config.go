package python

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/access"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/fn"
)

func genConfig(cfg *fn.Config, blk *fn.Block) string {
	if cfg.IsArray {
		return genConfigArray(cfg, blk)
	} else {
		return genConfigSingle(cfg, blk)
	}
}

func genConfigSingle(cfg *fn.Config, blk *fn.Block) string {
	var code string

	switch a := cfg.Access.(type) {
	case access.SingleSingle:
		code += indent + fmt.Sprintf(
			"self.%s = ConfigSingleSingle(iface, %d, (%d, %d))\n",
			cfg.Name, blk.StartAddr()+a.Addr, a.GetEndBit(), a.GetStartBit(),
		)
	case access.SingleContinuous:
		code += indent + fmt.Sprintf(
			"self.%s = ConfigSingleContinuous(iface, %d, %d, (%d, %d), (%d, %d))\n",
			cfg.Name,
			blk.StartAddr()+a.GetStartAddr(),
			a.GetRegCount(),
			busWidth-1, a.GetStartBit(),
			a.GetEndBit(), 0,
		)
	default:
		panic("not yet implemented")
	}

	return code
}

func genConfigArray(cfg *fn.Config, blk *fn.Block) string {
	var code string

	switch a := cfg.Access.(type) {
	case access.ArraySingle:
		code += indent + fmt.Sprintf(
			"self.%s = ConfigArraySingle(iface, %d, (%d, %d), %d)\n",
			cfg.Name,
			blk.StartAddr()+a.GetStartAddr(),
			a.GetEndBit(),
			a.GetStartBit(),
			a.GetRegCount(),
		)
	default:
		panic("unimplemented")
	}

	return code
}
