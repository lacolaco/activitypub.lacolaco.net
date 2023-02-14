package tracing

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/lacolaco/activitypub.lacolaco.net/config"
	"go.opentelemetry.io/otel"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
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
		ctx, span := Tracer().Start(ctx, c.Request.URL.String())
		defer span.End()
		span.SetAttributes(attribute.String("/http/method", c.Request.Method))
		span.SetAttributes(attribute.String("/http/host", c.Request.Host))
		span.SetAttributes(attribute.String("/http/path", c.Request.URL.String()))
		span.SetAttributes(attribute.String("/http/user_agent", c.Request.UserAgent()))

		ctx = context.WithValue(ctx, traceNameContextKey, buildTraceName(cfg.ProjectID(), span.SpanContext().TraceID().String()))
		c.Request = c.Request.WithContext(ctx)
		c.Next()

		span.SetAttributes(attribute.Int("/http/status_code", c.Writer.Status()))
		span.SetAttributes(attribute.Int("/http/response/size", c.Writer.Size()))
	}
}

func TraceNameFromContext(c context.Context) string {
	if traceName, ok := c.Value(traceNameContextKey).(string); ok {
		return traceName
	}
	return ""
}
