package python

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/elem"
)

func genStream(stream *elem.Stream, blk *elem.Block) string {
	if stream.IsDownstream() {
		panic("downstream not yet supported")
	}

	streamType := "Downstream"
	if stream.IsUpstream() {
		streamType = "Upstream"
	}

	code := indent + fmt.Sprintf("self.%s = %s(iface, %d, ",
		stream.Name, streamType, blk.StartAddr()+stream.StartAddr(),
	)

	if stream.IsDownstream() {
		code += genParamList(stream.Params)
	} else {
		code += genReturnList(stream.Returns)
	}

	code += ")\n"

	return code
}
