package traceidctx

import (
	"context"

	"github.com/WithSoull/platform_common/pkg/contextx"
)

const TraceIDKey contextx.CtxKey = "trace_id"

func InjectTraceId(ctx context.Context, traceID int64) context.Context {
	return context.WithValue(ctx, TraceIDKey, traceID)
}

func ExtractTraceId(ctx context.Context) (int64, bool) {
	if traceId, ok := ctx.Value(TraceIDKey).(int64); ok {
		return traceId, true
	} 
	return 0, false
}
