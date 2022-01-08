package vhdl

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl"
)

type accessChunk struct {
	addr   [2]int64
	range_ [2]string
	mask   fbdl.AccessMask
}

func makeAccessChunks(access fbdl.Access) []accessChunk {
	if asc, ok := access.(fbdl.AccessSingleContinuous); ok {
		return makeAccessChunksContinuous(asc)
	} else {
		panic("not yet implemented")
	}
}

func makeAccessChunksContinuous(a fbdl.AccessSingleContinuous) []accessChunk {
	startMask := a.StartMask
	endMask := a.EndMask

	cs := []accessChunk{}

	if a.RegCount() == 2 {
		cs = append(cs, accessChunk{
			addr:   [2]int64{a.StartAddr(), a.StartAddr()},
			range_: [2]string{fmt.Sprintf("%d", startMask.Width()-1), "0"},
			mask:   startMask,
		})
		cs = append(cs, accessChunk{
			addr: [2]int64{a.EndAddr(), a.EndAddr()},
			range_: [2]string{
				fmt.Sprintf("%d", a.Width()-1),
				fmt.Sprintf("%d", a.Width()-endMask.Width()),
			},
			mask: endMask,
		})
	} else if startMask.Width() == busWidth && endMask.Width() == busWidth {
		cs = append(cs, accessChunk{
			addr: [2]int64{a.StartAddr(), a.EndAddr()},
			range_: [2]string{
				fmt.Sprintf("%d * (addr - %d + 1) - 1", busWidth, a.StartAddr()),
				fmt.Sprintf("%d * (addr - %d)", busWidth, a.StartAddr()),
			},
			mask: startMask,
		})
	} else if startMask.Width() == busWidth {
		cs = append(cs, accessChunk{
			addr: [2]int64{a.StartAddr(), a.EndAddr() - 1},
			range_: [2]string{
				fmt.Sprintf("%d * (addr - %d + 1) - 1", busWidth, a.StartAddr()),
				fmt.Sprintf("%d * (addr - %d)", busWidth, a.StartAddr()),
			},
			mask: startMask,
		})
		cs = append(cs, accessChunk{
			addr: [2]int64{a.EndAddr(), a.EndAddr()},
			range_: [2]string{
				fmt.Sprintf("%d", a.Width()-1),
				fmt.Sprintf("%d", a.Width()-endMask.Width()),
			},
			mask: endMask,
		})
	} else if endMask.Width() == busWidth {
		cs = append(cs, accessChunk{
			addr:   [2]int64{a.StartAddr(), a.StartAddr()},
			range_: [2]string{fmt.Sprintf("%d", startMask.Width()-1), "0"},
			mask:   startMask,
		})
		cs = append(cs, accessChunk{
			addr: [2]int64{a.StartAddr() + 1, a.EndAddr()},
			range_: [2]string{
				fmt.Sprintf("%d * (addr - %d + 1) + %d", busWidth, a.StartAddr(), startMask.Width()-1),
				fmt.Sprintf("%d * (addr - %d) + %d", busWidth, a.StartAddr(), startMask.Width()),
			},
			mask: startMask,
		})
	} else {
		cs = append(cs, accessChunk{
			addr:   [2]int64{a.StartAddr(), a.StartAddr()},
			range_: [2]string{fmt.Sprintf("%d", startMask.Width()-1), "0"},
			mask:   startMask,
		})
		cs = append(cs, accessChunk{
			addr: [2]int64{a.StartAddr() + 1, a.EndAddr() - 1},
			range_: [2]string{
				fmt.Sprintf("%d * (addr - %d + 1) + %d", busWidth, a.StartAddr(), startMask.Width()-1),
				fmt.Sprintf("%d * (addr - %d) + %d", busWidth, a.StartAddr(), startMask.Width()),
			},
			mask: startMask,
		})
		cs = append(cs, accessChunk{
			addr: [2]int64{a.EndAddr(), a.EndAddr()},
			range_: [2]string{
				fmt.Sprintf("%d", a.Width()-1),
				fmt.Sprintf("%d", a.Width()-endMask.Width()),
			},
			mask: endMask,
		})
	}

	return cs
}