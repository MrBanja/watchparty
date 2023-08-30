package logging

import (
	"context"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/rs/xid"
	"go.uber.org/zap"
)

func Middleware(next http.Handler, logger *zap.Logger) http.Handler {
	next = handlers.CustomLoggingHandler(os.Stdout, next, formatter(logger))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := xid.New().String()
		r = r.WithContext(context.WithValue(r.Context(), traceIDKey, id))

		clientID := r.Header.Get("X-Client-Id")
		if clientID == "" {
			clientID = "empty"
		}
		r = r.WithContext(context.WithValue(r.Context(), clientIDKey, clientID))

		l := WithTrace(logger, r.Context())
		l.Info(
			">>>",
			zap.String("Method", r.Method),
			zap.String("URL", r.RequestURI),
		)
		next.ServeHTTP(w, r)
	})
}

func formatter(logger *zap.Logger) handlers.LogFormatter {
	return func(writer io.Writer, params handlers.LogFormatterParams) {
		l := WithTrace(logger, params.Request.Context())
		l.Info(
			"<<<",
			zap.String("URL", params.URL.String()),
			zap.Int("Status", params.StatusCode),
		)
	}
}
