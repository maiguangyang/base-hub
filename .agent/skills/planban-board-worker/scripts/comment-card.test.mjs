import assert from 'node:assert/strict';
import test from 'node:test';

import { buildCommentRequest, parseCommentArgs } from './comment-card.mjs';

test('parseCommentArgs defaults to dry-run and requires card text', () => {
  assert.throws(() => parseCommentArgs([]), /Missing value for --card/);

  const args = parseCommentArgs(['--card', '123', '--text', 'hello']);

  assert.equal(args.apply, false);
  assert.equal(args.card, '123');
  assert.equal(args.text, 'hello');
});

test('buildCommentRequest matches the documented live contract', () => {
  const request = buildCommentRequest({
    apiBase: 'https://planban-api.sjfood.us/api',
    cardId: '123',
    text: 'Codex write probe',
  });

  assert.equal(request.method, 'POST');
  assert.equal(request.url, 'https://planban-api.sjfood.us/api/cards/123/comment-actions');
  assert.equal(request.body, JSON.stringify({
    text: 'Codex write probe',
    plainText: 'Codex write probe',
    atUserIds: [],
  }));
});
