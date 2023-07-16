import { Config } from '@app/domain/config';
import { TraceExporter } from '@google-cloud/opentelemetry-cloud-trace-exporter';
import { CloudPropagator } from '@google-cloud/opentelemetry-cloud-trace-propagator';
import { Span, trace } from '@opentelemetry/api';
import { CompositePropagator, W3CBaggagePropagator, W3CTraceContextPropagator } from '@opentelemetry/core';
import { registerInstrumentations } from '@opentelemetry/instrumentation';
import { BasicTracerProvider, ConsoleSpanExporter, SimpleSpanProcessor } from '@opentelemetry/sdk-trace-node';
import { MiddlewareHandler } from 'hono';

export function setupTracing(config: Config) {
  const provider = new BasicTracerProvider({});

  registerInstrumentations({ tracerProvider: provider });

  const exporter = config.isRunningOnCloud ? new TraceExporter() : new ConsoleSpanExporter();
  provider.addSpanProcessor(new SimpleSpanProcessor(exporter));

  provider.register({
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
    console.log(JSON.stringify(Object.fromEntries(c.req.headers.entries())));
    const spanContext = trace.getActiveSpan()?.spanContext();
    console.log(`currentSpan: ${spanContext?.spanId} traceID: ${spanContext?.traceId}`);
    const span = getTracer().startSpan('hono.request');
    try {
      await next();
    } finally {
      span.end();
    }
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
