import { TraceExporter } from '@google-cloud/opentelemetry-cloud-trace-exporter';
import { trace } from '@opentelemetry/api';
import { NodeTracerProvider, SimpleSpanProcessor } from '@opentelemetry/sdk-trace-node';

export function setupTracing() {
  // Create and configure NodeTracerProvider
  const provider = new NodeTracerProvider();

  const exporter = new TraceExporter();

  // Configure the span processor to send spans to the exporter
  provider.addSpanProcessor(new SimpleSpanProcessor(exporter));

  // Initialize the provider
  provider.register();
}

export function getTracer() {
  return trace.getTracer('default');
}

export function buildTraceName(projectID: string, traceID: string): string {
  if (projectID == '' || traceID == '') {
    return '';
  }
  return `projects/${projectID}/traces/${traceID}`;
}
