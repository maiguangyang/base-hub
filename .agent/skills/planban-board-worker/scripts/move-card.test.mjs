import assert from 'node:assert/strict';
import test from 'node:test';

import {
  buildMoveRequest,
  deriveStateAfterMove,
  resolveTargetListName,
  validateTransition,
} from './move-card.mjs';

test('resolveTargetListName only allows approved destinations', () => {
  assert.equal(resolveTargetListName('To Do'), 'To Do');
  assert.equal(resolveTargetListName('In Progress'), 'In Progress');
  assert.equal(resolveTargetListName('Need To Discuss'), 'Need To Discuss');
  assert.throws(() => resolveTargetListName('Done'), /Unsupported target list/);
});

test('buildMoveRequest matches the documented live contract', () => {
  const request = buildMoveRequest({
    apiBase: 'https://planban-api.sjfood.us/api',
    cardId: '123',
    listId: '456',
    position: 131070,
  });

  assert.equal(request.method, 'PATCH');
  assert.equal(request.url, 'https://planban-api.sjfood.us/api/cards/123');
  assert.equal(request.body, JSON.stringify({
    listId: '456',
    position: 131070,
  }));
});

test('validateTransition only allows the documented workflow transitions', () => {
  assert.equal(validateTransition('Triage', 'To Do'), 'To Do');
  assert.equal(validateTransition('Triage', 'Need To Discuss'), 'Need To Discuss');
  assert.equal(validateTransition('To Do', 'In Progress'), 'In Progress');
  assert.throws(() => validateTransition('In Progress', 'To Do'), /Unsupported transition/);
  assert.throws(() => validateTransition('Done', 'In Progress'), /Unsupported transition/);
});

test('deriveStateAfterMove marks classified triage cards as processed and advances current card', () => {
  const state = deriveStateAfterMove({
    rememberedState: {
      currentBoardUrl: 'https://planban.sjfood.us/boards/1793465014108029954',
      currentMode: 'semi-auto',
      currentCardId: 'card-1',
      processedCardIds: [],
    },
    boardUrl: 'https://planban.sjfood.us/boards/1793465014108029954',
    movedCardId: 'card-1',
    fromListName: 'Triage',
    targetListName: 'To Do',
    triageCardsAfterMove: [
      { id: 'card-2' },
      { id: 'card-3' },
    ],
    now: '2026-06-25T03:00:00.000Z',
  });

  assert.equal(state.currentCardId, 'card-2');
  assert.deepEqual(state.processedCardIds, ['card-1']);
});
