import assert from 'node:assert/strict';
import test from 'node:test';

import {
  buildIssueMarkdown,
  deriveStateAfterScan,
  resolveActiveBoardUrl,
  validateBoardLists,
} from './scan-board.mjs';

test('resolveActiveBoardUrl prefers explicit board url', async () => {
  const url = await resolveActiveBoardUrl({
    cliBoardUrl: 'https://planban.sjfood.us/boards/1793465014108029954',
    rememberedState: { currentBoardUrl: 'https://planban.sjfood.us/boards/old' },
  });
  assert.equal(url, 'https://planban.sjfood.us/boards/1793465014108029954');
});

test('resolveActiveBoardUrl throws when no board is remembered', async () => {
  await assert.rejects(
    () => resolveActiveBoardUrl({ cliBoardUrl: '', rememberedState: null }),
    /Provide board URL/,
  );
});

test('validateBoardLists rejects boards missing Need To Discuss', () => {
  assert.throws(
    () => validateBoardLists([
      { id: '1', name: 'Triage' },
      { id: '2', name: 'To Do' },
      { id: '3', name: 'In Progress' },
    ]),
    /Missing required lists: Need To Discuss/,
  );
});

test('buildIssueMarkdown exports description comments and attachments', () => {
  const issue = buildIssueMarkdown({
    board: { name: 'AI POS Mobile v1.0' },
    card: {
      id: 'card-1',
      name: '无法批量上架',
      description: '批量上架按钮点击后没有反应',
      createdAt: '2026-06-25T00:00:00.000Z',
      updatedAt: '2026-06-25T01:00:00.000Z',
    },
    listName: 'Triage',
    labelNames: ['engine'],
    creatorName: 'Tester',
    comments: [
      {
        createdAt: '2026-06-25T01:30:00.000Z',
        data: { plainText: '点了没有任何提示' },
      },
    ],
    attachments: [
      {
        path: '/tmp/run/cards/card-1/attachments/shot.png',
        issuePath: '/tmp/run/cards/card-1/issue.md',
        downloaded: true,
      },
    ],
  });

  assert.match(issue, /批量上架按钮点击后没有反应/);
  assert.match(issue, /点了没有任何提示/);
  assert.match(issue, /attachments\/shot\.png/);
});

test('deriveStateAfterScan keeps progress and selects the next unprocessed triage card', () => {
  const state = deriveStateAfterScan({
    boardUrl: 'https://planban.sjfood.us/boards/1793465014108029954',
    mode: 'semi-auto',
    triageCards: [
      { id: 'card-1' },
      { id: 'card-2' },
      { id: 'card-3' },
    ],
    rememberedState: {
      currentCardId: 'missing-card',
      processedCardIds: ['card-1', 'done-card'],
    },
    now: '2026-06-25T02:50:00.000Z',
  });

  assert.equal(state.currentCardId, 'card-2');
  assert.deepEqual(state.processedCardIds, ['card-1', 'done-card']);
});
