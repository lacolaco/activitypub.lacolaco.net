package logging

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/lacolaco/activitypub.lacolaco.net/config"
	"github.com/lacolaco/activitypub.lacolaco.net/tracing"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type contextKey int

const (
	loggerContextKey contextKey = iota
)

// Middleware returns a gin middleware that sets the logger in the context.
func Middleware(cfg *config.Config) gin.HandlerFunc {
	logger := NewLogger(cfg)
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		ctx = ContextWithLogger(ctx, logger)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// ContextWithLogger returns a new context with the given logger.
func ContextWithLogger(c context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(c, loggerContextKey, logger)
}

// FromContext returns the logger for the given context.
// The logger has been set a trace context if the request is traced.
func FromContext(c context.Context) *zap.Logger {
	logger, ok := c.Value(loggerContextKey).(*zap.Logger)
	if !ok {
		return zap.NewNop()
	}
	return logger.WithOptions(traceContext(c))
}

const (
	traceKey        = "logging.googleapis.com/trace"
	spanKey         = "logging.googleapis.com/spanId"
	traceSampledKey = "logging.googleapis.com/trace_sampled"
)

func traceContext(c context.Context) zap.Option {
	traceName := tracing.TraceNameFromContext(c)
	spanContext := trace.SpanContextFromContext(c)
	if traceName == "" || !spanContext.IsValid() {
		return zap.Fields()
	}
	return zap.Fields(
		zap.String(traceKey, traceName),
		zap.String(spanKey, spanContext.SpanID().String()),
		zap.Bool(traceSampledKey, spanContext.IsSampled()),
	)
}
