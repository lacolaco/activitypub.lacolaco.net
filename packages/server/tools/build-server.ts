import TsconfigPathsPlugin from '@esbuild-plugins/tsconfig-paths';
import { build } from 'esbuild';

build({
  bundle: true,
  packages: 'external',
  target: 'node18',
  entryPoints: ['src/main.ts'],
  outfile: 'dist/main.js',
  tsconfig: 'tsconfig.json',
  plugins: [TsconfigPathsPlugin({ tsconfig: 'tsconfig.json' })],
  logLevel: 'info',
  supported: {
    'top-level-await': false,
  },
}).catch(() => process.exit(1));
