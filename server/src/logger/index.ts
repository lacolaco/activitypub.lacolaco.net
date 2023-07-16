import { Config } from '@app/domain/config';
import { SpanContext, trace } from '@opentelemetry/api';
import { getPinoOptions } from '@relaycorp/pino-cloud';
import pino from 'pino';

export type Logger = pino.Logger;

const gcpPinoOptions = getPinoOptions('gcp', {
  name: 'lacolaco-activitypub',
});

export function createLogger(config: Config) {
  if (config.isRunningOnCloud) {
    return pino({
      ...gcpPinoOptions,
      level: 'info',
    });
  } else {
    return pino({ level: 'trace' });
  }
}

const traceKey = 'logging.googleapis.com/trace';
const spanKey = 'logging.googleapis.com/spanId';
const traceSampledKey = 'logging.googleapis.com/trace_sampled';

function buildTraceName(projectID: string, traceID: string): string {
  if (projectID == '' || traceID == '') {
    return '';
  }
  return `projects/${projectID}/traces/${traceID}`;
}

export function createLoggerWithTrace(parent: Logger, config: Config, spanContext: SpanContext) {
  if (!trace.isSpanContextValid(spanContext) || config.gcpProjectID == null) {
    return parent.child({});
  }
  const traceName = buildTraceName(config.gcpProjectID, spanContext.traceId);
  return parent.child({
    [traceKey]: traceName,
    [spanKey]: spanContext.spanId,
    [traceSampledKey]: spanContext.traceFlags === 1,
  });
}
