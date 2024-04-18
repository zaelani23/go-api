package method

import (
	"github.com/whatap/golib/lang/pack/udp"
	"github.com/zaelani/go-api/trace"
)

type MethodCtx struct {
	ctx  *trace.TraceCtx
	step *udp.UdpTxMethodPack
}

func NewMethodCtx() *MethodCtx {
	p := new(MethodCtx)
	return p
}
