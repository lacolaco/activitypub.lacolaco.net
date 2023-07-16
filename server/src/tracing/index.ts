import { Config } from '@app/domain/config';
import { TraceExporter } from '@google-cloud/opentelemetry-cloud-trace-exporter';
import { Span, trace } from '@opentelemetry/api';
import { ConsoleSpanExporter, NodeTracerProvider, SimpleSpanProcessor } from '@opentelemetry/sdk-trace-node';

export function setupTracing(config: Config) {
  // Create and configure NodeTracerProvider
  const provider = new NodeTracerProvider();

  const exporter = config.isRunningOnCloud ? new TraceExporter() : new ConsoleSpanExporter();

  // Configure the span processor to send spans to the exporter
  provider.addSpanProcessor(new SimpleSpanProcessor(exporter));

  // Initialize the provider
  provider.register();

  trace.setGlobalTracerProvider(provider);
}

export function getTracer() {
  return trace.getTracer('default');
}

export function runInSpan<T>(name: string, fn: (span: Span) => T): Promise<T> {
  return getTracer().startActiveSpan(name, async (span) => {
    try {
      return await fn(span);
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
