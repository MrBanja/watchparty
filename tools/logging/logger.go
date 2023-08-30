package logging

import (
	"context"

	"go.uber.org/zap"
)

const (
	traceIDKey  = "__trace_xid"
	clientIDKey = "__client_xid"
)

func WithTrace(logger *zap.Logger, requestCtx context.Context) *zap.Logger {
	traceID := requestCtx.Value(traceIDKey).(string)
	clientID := requestCtx.Value(clientIDKey).(string)
	return logger.With(zap.String("RequestID", traceID), zap.String("ClientID", clientID))
}
