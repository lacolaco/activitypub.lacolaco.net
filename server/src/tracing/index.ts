import { Config } from '@app/domain/config';
import { TraceExporter } from '@google-cloud/opentelemetry-cloud-trace-exporter';
import { CloudPropagator } from '@google-cloud/opentelemetry-cloud-trace-propagator';
import { Span, context, propagation, trace } from '@opentelemetry/api';
import { AsyncLocalStorageContextManager } from '@opentelemetry/context-async-hooks';
import { CompositePropagator, W3CBaggagePropagator, W3CTraceContextPropagator } from '@opentelemetry/core';
import { ConsoleSpanExporter, NodeTracerProvider, SimpleSpanProcessor } from '@opentelemetry/sdk-trace-node';
import { MiddlewareHandler } from 'hono';

export function setupTracing(config: Config) {
  const provider = new NodeTracerProvider({});

  const exporter = config.isRunningOnCloud ? new TraceExporter() : new ConsoleSpanExporter();
  provider.addSpanProcessor(new SimpleSpanProcessor(exporter));

  provider.register({
    contextManager: new AsyncLocalStorageContextManager(),
    propagator: new CompositePropagator({
      propagators: [new CloudPropagator(), new W3CTraceContextPropagator(), new W3CBaggagePropagator()],
    }),
  });

  trace.setGlobalTracerProvider(provider);
}

export function getTracer() {
  return trace.getTracer('default');
}

export function withTracing(): MiddlewareHandler {
  return async (c, next) => {
    const headers = Object.fromEntries(c.req.headers.entries());
    const traceContext = propagation.extract(context.active(), headers);
    await context.with(traceContext, async () => {
      await next();
    });
  };
}

export function runInSpan<T>(name: string, fn: (span: Span) => T): Promise<T> {
  return getTracer().startActiveSpan(name, (span) => {
    try {
      const ret = Promise.resolve(fn(span));
      return ret.finally(() => {
        span.end();
      });
    } finally {
      span.end();
    }
  });
}

export function buildTraceName(projectID: string, traceID: string): string {
  if (projectID == '' || traceID == '') {
    return '';
  }
  return `projects/${projectID}/traces/${traceID}`;
}
