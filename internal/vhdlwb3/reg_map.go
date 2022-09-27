package vhdlwb3

type RegisterMap map[[2]int64]string

func (rm RegisterMap) add(addr [2]int64, code string) {
	if addr[1] < addr[0] {
		panic("addr[1] < addr[0]")
	}

	overlaps := [][2]int64{}
	for a := range rm {
		if (a[0] <= addr[0] && addr[0] <= a[1]) ||
			a[0] <= addr[1] && addr[1] <= a[1] {
			overlaps = append(overlaps, a)
		}
	}

	if len(overlaps) == 0 {
		rm[addr] = code
		return
	}

	if len(overlaps) == 1 && overlaps[0][0] == addr[0] && overlaps[0][1] == addr[1] {
		rm[addr] += code
		return
	}

	for _, o := range overlaps {
		tmpCode := rm[o]
		delete(rm, o)

		// Middle overlap
		if o[0] < addr[0] && addr[1] < o[1] {
			rm[[2]int64{o[0], addr[0] - 1}] = tmpCode
			rm[addr] = tmpCode + code
			rm[[2]int64{addr[1] + 1, o[1]}] = tmpCode
		}
		// Start overlap
		if addr[0] <= o[0] && addr[1] < o[1] {
			rm[[2]int64{addr[1] + 1, o[1]}] = tmpCode
			if o[0] == addr[0] {
				rm[addr] = tmpCode + code
			} else {
				rm[[2]int64{addr[0], o[0] - 1}] = code
				rm[[2]int64{o[0], addr[1]}] = tmpCode + code
			}
		}
		// End overlap
		if o[0] < addr[0] && o[1] <= addr[1] {
			rm[[2]int64{o[0], addr[0] - 1}] = tmpCode
			if o[1] == addr[1] {
				rm[addr] = tmpCode + code
			} else {
				rm[[2]int64{addr[0], o[1]}] = tmpCode + code
				rm[[2]int64{o[1] + 1, addr[1]}] = code
			}
		}
	}
}
