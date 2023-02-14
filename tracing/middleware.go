package tracing

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/lacolaco/activitypub.lacolaco.net/config"
	"go.opentelemetry.io/otel"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type contextKey int

const (
	traceNameContextKey contextKey = iota
)

func WithTracing(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		originalCtx := c.Request.Context()
		defer func() {
			c.Request = c.Request.WithContext(originalCtx)
		}()
		ctx := otel.GetTextMapPropagator().Extract(originalCtx, propagation.HeaderCarrier(c.Request.Header))
		spanContext := trace.SpanContextFromContext(ctx)
		ctx = context.WithValue(ctx, traceNameContextKey, buildTraceName(cfg.ProjectID(), spanContext.TraceID().String()))
		if !spanContext.HasSpanID() {
			var rootSpan trace.Span
			ctx, rootSpan = Tracer().Start(ctx, c.Request.URL.Path)
			rootSpan.SetAttributes(attribute.String("/http/method", c.Request.Method))
			rootSpan.SetAttributes(attribute.String("/http/host", c.Request.Host))
			rootSpan.SetAttributes(attribute.String("/http/path", c.Request.URL.String()))
			rootSpan.SetAttributes(attribute.String("/http/user_agent", c.Request.UserAgent()))
			defer func() {
				rootSpan.SetAttributes(attribute.Int("/http/status_code", c.Writer.Status()))
				rootSpan.SetAttributes(attribute.Int("/http/response/size", c.Writer.Size()))
				rootSpan.End()
			}()
		}
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func TraceNameFromContext(c context.Context) string {
	if traceName, ok := c.Value(traceNameContextKey).(string); ok {
		return traceName
	}
	return ""
}

func SpanContextFromContext(ctx context.Context) trace.SpanContext {
	return trace.SpanContextFromContext(ctx)
}

func SpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

func StartSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	return Tracer().Start(ctx, name)
}
