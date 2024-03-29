import { describe, expect, test } from 'vitest';
import { withTracing, setupTracing, runInSpan } from './index';
import { Hono } from 'hono';
import { getConfigWithEnv } from '@app/domain/config';
import { trace } from '@opentelemetry/api';
import { X_CLOUD_TRACE_HEADER } from '@google-cloud/opentelemetry-cloud-trace-propagator';

describe('tracing', () => {
  test('trace context from XCTC header', async () => {
    const config = await getConfigWithEnv();
    const { shutdown } = setupTracing(config);
    const app = new Hono();
    app.use('*', withTracing());
    app.get('/', (c) => {
      const span = trace.getActiveSpan();
      if (span == null) {
        return c.json({ error: 'No active span' }, 400);
      }
      return c.json(span.spanContext());
    });

    const res = await app.request('/', {
      method: 'GET',
      headers: {
        'X-Cloud-Trace-Context': '105445aa7843bc8bf206b12000100000/1;o=1',
      },
    });
    await shutdown();

    if (!res.ok) {
      throw new Error(await res.text());
    }
    const span = await res.json();
    expect(span.traceId).toBe('105445aa7843bc8bf206b12000100000');
    expect(span.spanId).not.toBe('1');
  });

  test('runInSpan should create a new child span', async (t) => {
    const config = await getConfigWithEnv();
    const { shutdown } = setupTracing(config);
    const app = new Hono();
    app.use('*', withTracing());
    app.get('/', async (c) => {
      return runInSpan('test', async (span) => {
        if (!span.isRecording()) {
          return c.json({ error: 'No active span' }, 400);
        }
        return c.json(span.spanContext());
      });
    });

    const res = await app.request('/', {
      method: 'GET',
      headers: {
        [X_CLOUD_TRACE_HEADER]: '105445aa7843bc8bf206b12000100000/1;o=1',
      },
    });
    await shutdown();
    if (!res.ok) {
      throw new Error(await res.text());
    }
  });
});
