package vhdlwb3

import (
	"testing"
)

type testTuple struct {
	addr [2]int64
	code string
}

func TestEndOverlap(t *testing.T) {
	var tests = []struct {
		first  testTuple
		second testTuple
		want   RegisterMap
	}{
		{
			testTuple{addr: [2]int64{0, 1}, code: "f"},
			testTuple{addr: [2]int64{1, 1}, code: "s"},
			RegisterMap{
				[2]int64{0, 0}: "f",
				[2]int64{1, 1}: "fs",
			},
		},
		{
			testTuple{addr: [2]int64{0, 2}, code: "f"},
			testTuple{addr: [2]int64{1, 2}, code: "s"},
			RegisterMap{
				[2]int64{0, 0}: "f",
				[2]int64{1, 2}: "fs",
			},
		},
		{
			testTuple{addr: [2]int64{0, 3}, code: "f"},
			testTuple{addr: [2]int64{2, 5}, code: "s"},
			RegisterMap{
				[2]int64{0, 1}: "f",
				[2]int64{2, 3}: "fs",
				[2]int64{4, 5}: "s",
			},
		},
	}

	for i, test := range tests {
		rm := RegisterMap{}

		rm.add(test.first.addr, test.first.code)
		rm.add(test.second.addr, test.second.code)

		for addr := range test.want {
			if rm[addr] != test.want[addr] {
				t.Errorf("[%d]: got %v, want %v", i, rm, test.want)
			}
		}
	}
}

func TestMiddleOverlap(t *testing.T) {
	var tests = []struct {
		first  testTuple
		second testTuple
		want   RegisterMap
	}{
		{
			testTuple{addr: [2]int64{0, 3}, code: "f"},
			testTuple{addr: [2]int64{1, 2}, code: "s"},
			RegisterMap{
				[2]int64{0, 0}: "f",
				[2]int64{1, 2}: "fs",
				[2]int64{3, 3}: "f",
			},
		},
		{
			testTuple{addr: [2]int64{0, 6}, code: "f"},
			testTuple{addr: [2]int64{5, 6}, code: "s"},
			RegisterMap{
				[2]int64{0, 4}: "f",
				[2]int64{5, 6}: "fs",
			},
		},
	}

	for i, test := range tests {
		rm := RegisterMap{}

		rm.add(test.first.addr, test.first.code)
		rm.add(test.second.addr, test.second.code)

		for addr := range test.want {
			if rm[addr] != test.want[addr] {
				t.Errorf("[%d]: got %v, want %v", i, rm, test.want)
			}
		}
	}
}

func TestStartOverlap(t *testing.T) {
	var tests = []struct {
		first  testTuple
		second testTuple
		want   RegisterMap
	}{
		{
			testTuple{addr: [2]int64{0, 1}, code: "f"},
			testTuple{addr: [2]int64{0, 0}, code: "s"},
			RegisterMap{
				[2]int64{0, 0}: "fs",
				[2]int64{1, 1}: "f",
			},
		},
		{
			testTuple{addr: [2]int64{0, 2}, code: "f"},
			testTuple{addr: [2]int64{0, 1}, code: "s"},
			RegisterMap{
				[2]int64{0, 1}: "fs",
				[2]int64{2, 2}: "f",
			},
		},
		{
			testTuple{addr: [2]int64{3, 5}, code: "f"},
			testTuple{addr: [2]int64{2, 3}, code: "s"},
			RegisterMap{
				[2]int64{2, 2}: "s",
				[2]int64{3, 3}: "fs",
				[2]int64{4, 5}: "f",
			},
		},
	}

	for i, test := range tests {
		rm := RegisterMap{}

		rm.add(test.first.addr, test.first.code)
		rm.add(test.second.addr, test.second.code)

		for addr := range test.want {
			if rm[addr] != test.want[addr] {
				t.Errorf("[%d]: got %v, want %v", i, rm, test.want)
			}
		}
	}
}
