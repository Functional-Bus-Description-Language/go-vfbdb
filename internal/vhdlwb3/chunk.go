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
	startBit := a.GetStartBit()
	endBit := a.GetEndBit()

	cs := []accessChunk{}

	if strategy == Compact && a.GetStartRegWidth() == busWidth && a.GetEndRegWidth() == busWidth {
		cs = append(cs, accessChunk{
			addr: [2]int64{a.GetStartAddr(), a.GetEndAddr()},
			range_: [2]string{
				fmt.Sprintf("%d * (addr - %d + 1) - 1", busWidth, a.GetStartAddr()),
				fmt.Sprintf("%d * (addr - %d)", busWidth, a.GetStartAddr()),
			},
			startBit: 0,
			endBit:   busWidth - 1,
		})
	} else if a.GetRegCount() == 2 {
		cs = append(cs, accessChunk{
			addr:     [2]int64{a.GetStartAddr(), a.GetStartAddr()},
			range_:   [2]string{fmt.Sprintf("%d", a.GetStartRegWidth()-1), "0"},
			startBit: startBit,
			endBit:   busWidth - 1,
		})
		cs = append(cs, accessChunk{
			addr: [2]int64{a.GetEndAddr(), a.GetEndAddr()},
			range_: [2]string{
				fmt.Sprintf("%d", a.GetWidth()-1),
				fmt.Sprintf("%d", a.GetWidth()-a.GetEndRegWidth()),
			},
			startBit: 0,
			endBit:   endBit,
		})
	} else if strategy == SeparateLast && a.GetStartRegWidth() == busWidth {
		cs = append(cs, accessChunk{
			addr: [2]int64{a.GetStartAddr(), a.GetEndAddr() - 1},
			range_: [2]string{
				fmt.Sprintf("%d * (addr - %d + 1) - 1", busWidth, a.GetStartAddr()),
				fmt.Sprintf("%d * (addr - %d)", busWidth, a.GetStartAddr()),
			},
			startBit: 0,
			endBit:   busWidth - 1,
		})
		cs = append(cs, accessChunk{
			addr: [2]int64{a.GetEndAddr(), a.GetEndAddr()},
			range_: [2]string{
				fmt.Sprintf("%d", a.GetWidth()-1),
				fmt.Sprintf("%d", a.GetWidth()-a.GetEndRegWidth()),
			},
			startBit: 0,
			endBit:   endBit,
		})
	} else if strategy == SeparateFirst && a.GetEndRegWidth() == busWidth {
		cs = append(cs, accessChunk{
			addr:     [2]int64{a.GetStartAddr(), a.GetStartAddr()},
			range_:   [2]string{fmt.Sprintf("%d", a.GetStartRegWidth()-1), "0"},
			startBit: startBit,
			endBit:   busWidth - 1,
		})
		cs = append(cs, accessChunk{
			addr: [2]int64{a.GetStartAddr() + 1, a.GetEndAddr()},
			range_: [2]string{
				fmt.Sprintf("%d * (addr - %d + 1) + %d", busWidth, a.GetStartAddr(), a.GetStartRegWidth()-1),
				fmt.Sprintf("%d * (addr - %d) + %d", busWidth, a.GetStartAddr(), a.GetStartRegWidth()),
			},
			startBit: 0,
			endBit:   busWidth - 1,
		})
	} else {
		cs = append(cs, accessChunk{
			addr:     [2]int64{a.GetStartAddr(), a.GetStartAddr()},
			range_:   [2]string{fmt.Sprintf("%d", a.GetStartRegWidth()-1), "0"},
			startBit: startBit,
			endBit:   busWidth - 1,
		})
		cs = append(cs, accessChunk{
			addr: [2]int64{a.GetStartAddr() + 1, a.GetEndAddr() - 1},
			range_: [2]string{
				fmt.Sprintf("%d * (addr - %d) + %d", busWidth, a.GetStartAddr(), a.GetStartRegWidth()-1),
				fmt.Sprintf("%d * (addr - %d) + %d", busWidth, a.GetStartAddr()+1, a.GetStartRegWidth()),
			},
			startBit: 0,
			endBit:   busWidth - 1,
		})
		cs = append(cs, accessChunk{
			addr: [2]int64{a.GetEndAddr(), a.GetEndAddr()},
			range_: [2]string{
				fmt.Sprintf("%d", a.GetWidth()-1),
				fmt.Sprintf("%d", a.GetWidth()-a.GetEndRegWidth()),
			},
			startBit: 0,
			endBit:   endBit,
		})
	}

	return cs
}
