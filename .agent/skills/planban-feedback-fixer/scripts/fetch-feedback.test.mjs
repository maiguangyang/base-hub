import assert from 'node:assert/strict';
import { existsSync, mkdtempSync, readFileSync } from 'node:fs';
import { tmpdir } from 'node:os';
import path from 'node:path';
import { spawnSync } from 'node:child_process';
import { fileURLToPath } from 'node:url';
import test from 'node:test';

const scriptPath = fileURLToPath(new URL('./fetch-feedback.mjs', import.meta.url));

test('creates a local account file with setup instructions when credentials are missing', () => {
  const cwd = mkdtempSync(path.join(tmpdir(), 'planban-feedback-'));
  const result = spawnSync(process.execPath, [
    scriptPath,
    '--board',
    '1793465014108029954',
    '--label',
    '後台',
    '--out',
    'work/planban-feedback',
  ], {
    cwd,
    env: { ...process.env, PLANBAN_USERNAME: '', PLANBAN_PASSWORD: '' },
    encoding: 'utf8',
  });

  const accountPath = path.join(cwd, 'planban-account');

  assert.equal(result.status, 1);
  assert.equal(existsSync(accountPath), true);
  assert.match(result.stderr, /Created Planban account file:/);
  assert.match(result.stderr, /Fill in your Planban account/);
  assert.match(readFileSync(accountPath, 'utf8'), /Line 1: Planban username or email/);
});
