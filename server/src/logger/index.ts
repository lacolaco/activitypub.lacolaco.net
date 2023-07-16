import { Config } from '@app/domain/config';
import pino from 'pino';
import { getPinoOptions } from '@relaycorp/pino-cloud';

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
    return pino({ level: 'debug' });
  }
}
