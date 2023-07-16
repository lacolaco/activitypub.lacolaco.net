import { Config } from '@app/domain/config';
import { TraceExporter } from '@google-cloud/opentelemetry-cloud-trace-exporter';
import { CloudPropagator, X_CLOUD_TRACE_HEADER } from '@google-cloud/opentelemetry-cloud-trace-propagator';
import { DiagConsoleLogger, DiagLogLevel, Span, SpanKind, context, diag, propagation, trace } from '@opentelemetry/api';
import { CompositePropagator, W3CBaggagePropagator, W3CTraceContextPropagator } from '@opentelemetry/core';
import { ConsoleSpanExporter, NodeTracerProvider, SimpleSpanProcessor } from '@opentelemetry/sdk-trace-node';
import { MiddlewareHandler } from 'hono';

let initialized = false;

export function setupTracing(config: Config): { shutdown: () => Promise<void> } {
  if (initialized) {
    console.log('TraceProvider already setup');
    return {
      shutdown: async () => {},
    };
  }
  initialized = true;

  const provider = new NodeTracerProvider({});
  const exporter = config.isRunningOnCloud ? new TraceExporter() : new ConsoleSpanExporter();
  provider.addSpanProcessor(new SimpleSpanProcessor(exporter));
  provider.register({
    propagator: new CompositePropagator({
      propagators: [new CloudPropagator(), new W3CTraceContextPropagator(), new W3CBaggagePropagator()],
    }),
  });
  diag.setLogger(new DiagConsoleLogger(), DiagLogLevel.INFO);

  return {
    shutdown: async () => {
      await provider.shutdown();
    },
  };
}

export function getTracer() {
  return trace.getTracer('default');
}

export function withTracing(): MiddlewareHandler {
  return async (c, next) => {
    const traceHeaders = {
      [X_CLOUD_TRACE_HEADER]: c.req.headers.get(X_CLOUD_TRACE_HEADER) ?? '',
    };
    let traceContext = propagation.extract(context.active(), traceHeaders);
    if (trace.getSpanContext(traceContext) == null) {
      console.log('No trace context found');
    }
    await context.with(traceContext, async () => {
      await getTracer().startActiveSpan(
        'request',
        {
          attributes: {
            '/http/method': c.req.method,
            '/http/url': c.req.url,
          },
          kind: SpanKind.SERVER,
        },
        async (span) => {
          await next();
          span.setAttributes({
            '/http/status_code': c.res.status,
          });
          span.end();
        },
      );
    });
  };
}

export function runInSpan<T>(name: string, fn: (span: Span) => T): Promise<T> {
  return getTracer().startActiveSpan(name, (span) => {
    const ret = Promise.resolve(fn(span));
    return ret.finally(() => {
      span.end();
    });
  });
}

export function buildTraceName(projectID: string, traceID: string): string {
  if (projectID == '' || traceID == '') {
    return '';
  }
  return `projects/${projectID}/traces/${traceID}`;
}
