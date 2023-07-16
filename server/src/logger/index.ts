import { Config } from '@app/domain/config';
import { LoggingWinston } from '@google-cloud/logging-winston';
import * as winston from 'winston';

export function createLogger(config: Config) {
  const logger = winston.createLogger({
    level: 'verbose',
    transports: [new winston.transports.Console({ format: winston.format.simple() })],
  });

  if (config.isRunningOnCloud) {
    logger.configure({
      level: 'info',
      transports: [logger.add(new LoggingWinston({}))],
    });
  }
  return logger;
}
