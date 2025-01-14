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

	switch acs := cfg.Access.(type) {
	case access.SingleOneReg:
		code += indent + fmt.Sprintf(
			"self.%s = ConfigSingleOneReg(iface, %d, %d, %d)\n",
			cfg.Name,
			blk.StartAddr()+acs.StartAddr(),
			acs.StartBit(),
			acs.EndBit(),
		)
	case access.SingleNRegs:
		code += indent + fmt.Sprintf(
			"self.%s = ConfigSingleNRegs(iface, %d, %d, (%d, %d), (%d, %d))\n",
			cfg.Name,
			blk.StartAddr()+acs.StartAddr(),
			acs.RegCount(),
			busWidth-1, acs.StartBit(),
			acs.EndBit(), 0,
		)
	default:
		panic("unimplemented")
	}

	return code
}

func genConfigArray(cfg *fn.Config, blk *fn.Block) string {
	var code string

	switch acs := cfg.Access.(type) {
	case access.ArrayOneReg:
		code += indent + fmt.Sprintf(
			"self.%s = ConfigArrayOneReg(iface, %d, %d, %d, %d)\n",
			cfg.Name,
			blk.StartAddr()+acs.StartAddr(),
			acs.StartBit(),
			acs.ItemWidth(),
			acs.ItemCount(),
		)
	case access.ArrayOneInReg:
		code += indent + fmt.Sprintf(
			"self.%s = ConfigArrayOneInReg(iface, %d, (%d, %d), %d)\n",
			cfg.Name,
			blk.StartAddr()+acs.StartAddr(),
			acs.EndBit(),
			acs.StartBit(),
			acs.RegCount(),
		)
	case access.ArrayNInReg:
		code += indent + fmt.Sprintf(
			"self.%s = ConfigArrayNInReg(iface, %d, %d, %d, %d, %d)\n",
			cfg.Name,
			blk.StartAddr()+acs.StartAddr(),
			acs.StartBit(),
			acs.ItemWidth(),
			acs.ItemCount(),
			acs.ItemsInReg(),
		)
	case access.ArrayNInRegMInEndReg:
		code += indent + fmt.Sprintf(
			"self.%s = ConfigArrayNInReg(iface, %d, %d, %d, %d, %d)\n",
			cfg.Name,
			blk.StartAddr()+acs.StartAddr(),
			acs.StartBit(),
			acs.ItemWidth(),
			acs.ItemCount(),
			acs.ItemsInReg(),
		)
	default:
		panic("unimplemented")
	}

	return code
}
