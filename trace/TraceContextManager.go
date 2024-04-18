package trace

import (
	"github.com/whatap/golib/util/hmap"
	agentconfig "github.com/zaelani23/go-api/agent/agent/config"
)

const (
	STAT_METHOD = 0
	STAT_SQL    = 1
	STAT_HTTPC  = 2
	STAT_DBC    = 3
	STAT_SOCKET = 4
)

var conf *agentconfig.Config = agentconfig.GetConfig()

var ctxTable *hmap.LongKeyLinkedMap = hmap.NewLongKeyLinkedMapDefault().SetMax(int(conf.TxMaxCount))

func AddGIDTraceCtx(GID int64, traceCtx *TraceCtx) {
	if !conf.GoUseGoroutineIDEnabled {
		return
	}
	ctxTable.Put(GID, traceCtx)
}
func GetGIDTraceCtx(GID int64) *TraceCtx {
	if !conf.GoUseGoroutineIDEnabled {
		return nil
	}

	if obj := ctxTable.Get(GID); obj != nil {
		if v, ok := obj.(*TraceCtx); ok {
			return v
		}
	}
	return nil
}

func RemoveGIDTraceCtx(GID int64) {
	if !conf.GoUseGoroutineIDEnabled {
		return
	}
	ctxTable.Remove(GID)
}
