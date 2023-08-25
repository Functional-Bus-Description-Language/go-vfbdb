package python

import (
	"fmt"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/cnst"
)

func genConsts(cc *cnst.Container) string {
	s := ""

	for name, b := range cc.Bools {
		v := "False"
		if b {
			v = "True"
		}
		s += indent + fmt.Sprintf("%s = %s\n", name, v)
	}
	for name, list := range cc.BoolLists {
		s += indent + fmt.Sprintf("%s = [", name)
		for _, b := range list {
			v := "False"
			if b {
				v = "True"
			}
			s += fmt.Sprintf("%s, ", v)
		}
		s = s[:len(s)-2]
		s += "]\n"
	}
	for name, i := range cc.Ints {
		s += indent + fmt.Sprintf("%s = %d\n", name, i)
	}
	for name, list := range cc.IntLists {
		s += indent + fmt.Sprintf("%s = [", name)
		for _, i := range list {
			s += fmt.Sprintf("%d, ", i)
		}
		s = s[:len(s)-2]
		s += "]\n"
	}
	for name, str := range cc.Strings {
		s += indent + fmt.Sprintf("%s = %q\n", name, str)
	}

	return s
}
