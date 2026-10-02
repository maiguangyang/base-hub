import assert from 'node:assert/strict';
import { mkdtempSync } from 'node:fs';
import { tmpdir } from 'node:os';
import path from 'node:path';
import test from 'node:test';

import { loadBoardState, saveBoardState } from './board-state.mjs';

test('loadBoardState returns null when the state file is missing', async () => {
  const dir = mkdtempSync(path.join(tmpdir(), 'planban-state-'));
  const state = await loadBoardState(path.join(dir, 'current-board.json'));
  assert.equal(state, null);
});

test('saveBoardState round-trips current board context', async () => {
  const dir = mkdtempSync(path.join(tmpdir(), 'planban-state-'));
  const filePath = path.join(dir, 'current-board.json');
  await saveBoardState(filePath, {
    currentBoardUrl: 'https://planban.sjfood.us/boards/1793465014108029954',
    currentMode: 'auto',
    currentCardId: null,
    processedCardIds: ['card-1'],
    lastRunAt: '2026-06-24T00:00:00.000Z',
  });
  const loaded = await loadBoardState(filePath);
  assert.equal(loaded.currentMode, 'auto');
  assert.deepEqual(loaded.processedCardIds, ['card-1']);
});
