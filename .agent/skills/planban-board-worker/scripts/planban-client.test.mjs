import assert from 'node:assert/strict';
import { existsSync, mkdtempSync, readFileSync, writeFileSync } from 'node:fs';
import { tmpdir } from 'node:os';
import path from 'node:path';
import test from 'node:test';

import {
  boardIdFrom,
  ensureAccountFile,
  normalizeListName,
  resolveAccount,
  resolveRequiredLists,
} from './planban-client.mjs';

test('boardIdFrom accepts board url and raw id', () => {
  assert.equal(boardIdFrom('https://planban.sjfood.us/boards/1793465014108029954'), '1793465014108029954');
  assert.equal(boardIdFrom('1793465014108029954'), '1793465014108029954');
  assert.throws(() => boardIdFrom('https://planban.sjfood.us/projects/1'), /Could not parse board id/);
});

test('normalizeListName folds whitespace and case', () => {
  assert.equal(normalizeListName(' Triage '), 'triage');
  assert.equal(normalizeListName('Need To Discuss'), 'need to discuss');
});

test('ensureAccountFile creates a local template when credentials are missing', async () => {
  const cwd = mkdtempSync(path.join(tmpdir(), 'planban-board-worker-'));
  const accountPath = path.join(cwd, 'planban-account');
  await ensureAccountFile(accountPath);
  assert.equal(existsSync(accountPath), true);
  assert.match(readFileSync(accountPath, 'utf8'), /Planban username or email/);
});

test('resolveAccount reads PLANBAN credentials from root .env when shell env is empty', async () => {
  const cwd = mkdtempSync(path.join(tmpdir(), 'planban-board-worker-'));
  writeFileSync(path.join(cwd, '.env'), [
    '# skill secrets',
    'PLANBAN_USERNAME=tester@example.com',
    'PLANBAN_PASSWORD=secret-pass',
    '',
  ].join('\n'));

  const account = await resolveAccount({
    cwd,
    accountFile: '',
    env: { PLANBAN_USERNAME: '', PLANBAN_PASSWORD: '' },
  });

  assert.deepEqual(account, {
    username: 'tester@example.com',
    password: 'secret-pass',
  });
});

test('resolveRequiredLists reports every missing list name', () => {
  assert.throws(
    () => resolveRequiredLists(
      [
        { id: '1', name: 'Triage' },
        { id: '2', name: 'To Do' },
      ],
      ['Triage', 'To Do', 'In Progress', 'Need To Discuss'],
    ),
    /Missing required lists: In Progress, Need To Discuss/,
  );
});
