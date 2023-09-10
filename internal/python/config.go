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
	case access.SingleOneReg:
		code += indent + fmt.Sprintf(
			"self.%s = ConfigSingleOneReg(iface, %d, (%d, %d))\n",
			cfg.Name, blk.StartAddr()+a.Addr, a.GetEndBit(), a.GetStartBit(),
		)
	case access.SingleNRegs:
		code += indent + fmt.Sprintf(
			"self.%s = ConfigSingleNRegs(iface, %d, %d, (%d, %d), (%d, %d))\n",
			cfg.Name,
			blk.StartAddr()+a.GetStartAddr(),
			a.GetRegCount(),
			busWidth-1, a.GetStartBit(),
			a.GetEndBit(), 0,
		)
	default:
		panic("unimplemented")
	}

	return code
}

func genConfigArray(cfg *fn.Config, blk *fn.Block) string {
	var code string

	switch acs := cfg.Access.(type) {
	case access.ArrayOneInReg:
		code += indent + fmt.Sprintf(
			"self.%s = ConfigArrayOneInReg(iface, %d, (%d, %d), %d)\n",
			cfg.Name,
			blk.StartAddr()+acs.StartAddr,
			acs.EndBit,
			acs.StartBit,
			acs.RegCount,
		)
	default:
		panic("unimplemented")
	}

	return code
}
