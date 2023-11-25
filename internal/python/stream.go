package python

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/fn"
)

func genStream(stream *fn.Stream, blk *fn.Block) string {
	if stream.IsDownstream() {
		panic("unimplemented")
	}

	streamType := "Downstream"
	if stream.IsUpstream() {
		streamType = "Upstream"
	}

	code := indent + fmt.Sprintf("self.%s = %s(iface, %d, ",
		stream.Name, streamType, blk.StartAddr()+stream.StartAddr(),
	)

	if stream.IsDownstream() {
		code += genParamList(stream.Params, blk)
	} else {
		code += genReturnList(stream.Returns, blk)
	}

	code += ")\n"

	return code
}
