package vhdlwb3

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/access"
)

type chunkStrategy uint8

const (
	Compact chunkStrategy = iota // Use only for non atomic elements.
	SeparateFirst
	SeparateLast
)

type accessChunk struct {
	addr     [2]int64
	range_   [2]string
	startBit int64
	endBit   int64
}

func makeAccessChunksContinuous(a access.SingleNRegs, strategy chunkStrategy) []accessChunk {
	startBit := a.StartBit()
	endBit := a.EndBit()

	cs := []accessChunk{}

	if strategy == Compact && a.StartRegWidth() == busWidth && a.EndRegWidth() == busWidth {
		cs = append(cs, accessChunk{
			addr: [2]int64{a.StartAddr(), a.EndAddr()},
			range_: [2]string{
				fmt.Sprintf("%d * (addr - %d + 1) - 1", busWidth, a.StartAddr()),
				fmt.Sprintf("%d * (addr - %d)", busWidth, a.StartAddr()),
			},
			startBit: 0,
			endBit:   busWidth - 1,
		})
	} else if a.RegCount() == 2 {
		cs = append(cs, accessChunk{
			addr:     [2]int64{a.StartAddr(), a.StartAddr()},
			range_:   [2]string{fmt.Sprintf("%d", a.StartRegWidth()-1), "0"},
			startBit: startBit,
			endBit:   busWidth - 1,
		})
		cs = append(cs, accessChunk{
			addr: [2]int64{a.EndAddr(), a.EndAddr()},
			range_: [2]string{
				fmt.Sprintf("%d", a.Width()-1),
				fmt.Sprintf("%d", a.Width()-a.EndRegWidth()),
			},
			startBit: 0,
			endBit:   endBit,
		})
	} else if strategy == SeparateLast && a.StartRegWidth() == busWidth {
		cs = append(cs, accessChunk{
			addr: [2]int64{a.StartAddr(), a.EndAddr() - 1},
			range_: [2]string{
				fmt.Sprintf("%d * (addr - %d + 1) - 1", busWidth, a.StartAddr()),
				fmt.Sprintf("%d * (addr - %d)", busWidth, a.StartAddr()),
			},
			startBit: 0,
			endBit:   busWidth - 1,
		})
		cs = append(cs, accessChunk{
			addr: [2]int64{a.EndAddr(), a.EndAddr()},
			range_: [2]string{
				fmt.Sprintf("%d", a.Width()-1),
				fmt.Sprintf("%d", a.Width()-a.EndRegWidth()),
			},
			startBit: 0,
			endBit:   endBit,
		})
	} else if strategy == SeparateFirst && a.EndRegWidth() == busWidth {
		cs = append(cs, accessChunk{
			addr:     [2]int64{a.StartAddr(), a.StartAddr()},
			range_:   [2]string{fmt.Sprintf("%d", a.StartRegWidth()-1), "0"},
			startBit: startBit,
			endBit:   busWidth - 1,
		})
		cs = append(cs, accessChunk{
			addr: [2]int64{a.StartAddr() + 1, a.EndAddr()},
			range_: [2]string{
				fmt.Sprintf("%d * (addr - %d + 1) + %d", busWidth, a.StartAddr(), a.StartRegWidth()-1),
				fmt.Sprintf("%d * (addr - %d) + %d", busWidth, a.StartAddr(), a.StartRegWidth()),
			},
			startBit: 0,
			endBit:   busWidth - 1,
		})
	} else {
		cs = append(cs, accessChunk{
			addr:     [2]int64{a.StartAddr(), a.StartAddr()},
			range_:   [2]string{fmt.Sprintf("%d", a.StartRegWidth()-1), "0"},
			startBit: startBit,
			endBit:   busWidth - 1,
		})
		cs = append(cs, accessChunk{
			addr: [2]int64{a.StartAddr() + 1, a.EndAddr() - 1},
			range_: [2]string{
				fmt.Sprintf("%d * (addr - %d) + %d", busWidth, a.StartAddr(), a.StartRegWidth()-1),
				fmt.Sprintf("%d * (addr - %d) + %d", busWidth, a.StartAddr()+1, a.StartRegWidth()),
			},
			startBit: 0,
			endBit:   busWidth - 1,
		})
		cs = append(cs, accessChunk{
			addr: [2]int64{a.EndAddr(), a.EndAddr()},
			range_: [2]string{
				fmt.Sprintf("%d", a.Width()-1),
				fmt.Sprintf("%d", a.Width()-a.EndRegWidth()),
			},
			startBit: 0,
			endBit:   endBit,
		})
	}

	return cs
}
