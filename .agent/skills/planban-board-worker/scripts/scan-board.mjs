#!/usr/bin/env node

import { mkdir, writeFile } from 'node:fs/promises';
import path from 'node:path';
import { fileURLToPath, pathToFileURL } from 'node:url';

import {
  authHeaders,
  boardIdFrom,
  DEFAULT_API_BASE,
  DEFAULT_SITE_BASE,
  indexBoardSnapshot,
  login,
  requestJson,
  resolveAccount,
  resolveRequiredLists,
} from './planban-client.mjs';
import { loadBoardState, saveBoardState } from './board-state.mjs';

const SKILL_DIR = path.dirname(path.dirname(fileURLToPath(import.meta.url)));

function parseArgs(argv) {
  const args = {};
  for (let i = 0; i < argv.length; i += 1) {
    const key = argv[i];
    if (!key.startsWith('--')) throw new Error(`Unexpected argument: ${key}`);
    const name = key.slice(2);
    const value = argv[i + 1];
    if (!value || value.startsWith('--')) throw new Error(`Missing value for ${key}`);
    args[name] = value;
    i += 1;
  }
  return args;
}

export async function resolveActiveBoardUrl({ cliBoardUrl, rememberedState }) {
  if (cliBoardUrl) return cliBoardUrl;
  if (rememberedState?.currentBoardUrl) return rememberedState.currentBoardUrl;
  throw new Error('Provide board URL');
}

export function validateBoardLists(lists) {
  const resolved = resolveRequiredLists(lists, ['Triage', 'To Do', 'In Progress', 'Need To Discuss']);
  return {
    triageId: resolved.Triage.id,
    toDoId: resolved['To Do'].id,
    inProgressId: resolved['In Progress'].id,
    needToDiscussId: resolved['Need To Discuss'].id,
  };
}

function sanitizeName(name) {
  return String(name || 'untitled')
    .replace(/[\\/:*?"<>|#%{}$!'@+`=]/g, '_')
    .replace(/\s+/g, ' ')
    .trim()
    .slice(0, 120) || 'untitled';
}

async function fetchWithTimeout(url, options = {}, timeoutMs = 15000) {
  const controller = new AbortController();
  const timeout = setTimeout(() => controller.abort(), timeoutMs);
  try {
    return await fetch(url, { ...options, signal: controller.signal });
  } finally {
    clearTimeout(timeout);
  }
}

async function downloadAttachment(attachment, token, targetDir) {
  const safeName = sanitizeName(attachment.name || `${attachment.id}.bin`);
  const suffix = attachment.image ? '' : '.url';
  const target = path.join(targetDir, `${attachment.id}-${safeName}${suffix}`);

  if (!attachment.image) {
    await writeFile(target, `${attachment.url}\n`, 'utf8');
    return {
      id: attachment.id,
      name: attachment.name,
      path: target,
      downloaded: false,
      type: attachment.type,
    };
  }

  try {
    const response = await fetchWithTimeout(attachment.url, {
      headers: { authorization: `Bearer ${token}` },
    });
    if (!response.ok) throw new Error(`HTTP ${response.status}`);
    const bytes = Buffer.from(await response.arrayBuffer());
    await writeFile(target, bytes);
    return {
      id: attachment.id,
      name: attachment.name,
      path: target,
      downloaded: true,
      type: attachment.type,
    };
  } catch (error) {
    await writeFile(`${target}.url`, `${attachment.url}\n# download failed: ${error.message}\n`, 'utf8');
    return {
      id: attachment.id,
      name: attachment.name,
      path: `${target}.url`,
      downloaded: false,
      type: attachment.type,
    };
  }
}

export function buildIssueMarkdown({
  board,
  card,
  listName,
  labelNames,
  creatorName,
  comments,
  attachments,
}) {
  const attachmentLines = attachments.length
    ? attachments.map((attachment) => `- ${path.relative(path.dirname(attachment.issuePath), attachment.path)}${attachment.downloaded ? '' : ' (link only)'}`).join('\n')
    : 'None';
  const commentLines = comments.length
    ? comments.map((comment) => `- ${comment.createdAt || 'unknown'}: ${comment.data?.plainText || comment.data?.text || JSON.stringify(comment.data)}`).join('\n')
    : 'None';

  return [
    `# ${card.name}`,
    '',
    `Board: ${board.name}`,
    `Card: ${card.id}`,
    `List: ${listName}`,
    `Labels: ${labelNames.join(', ') || 'None'}`,
    `Creator: ${creatorName}`,
    `Created: ${card.createdAt}`,
    `Updated: ${card.updatedAt || 'None'}`,
    '',
    '## Description',
    card.description || 'None',
    '',
    '## Comments',
    commentLines,
    '',
    '## Attachments',
    attachmentLines,
    '',
  ].join('\n');
}

function buildSummary(board, cards) {
  return [
    '# Planban Board Scan',
    '',
    `Board: ${board.name}`,
    `Board ID: ${board.id}`,
    `Triage Cards: ${cards.length}`,
    '',
    '## Triage',
    ...cards.map((card) => `- [${card.name}](${card.issuePath}) - comments: ${card.commentCount}, attachments: ${card.attachmentCount}, labels: ${card.labels.join(', ') || 'None'}`),
    '',
  ].join('\n');
}

export function deriveStateAfterScan({
  boardUrl,
  mode,
  triageCards,
  rememberedState,
  now,
}) {
  const processedCardIds = Array.from(new Set(rememberedState?.processedCardIds || []));
  const triageCardIds = triageCards.map((card) => card.id);
  const rememberedCurrentCardId = rememberedState?.currentCardId || null;
  const currentCardId = rememberedCurrentCardId
    && triageCardIds.includes(rememberedCurrentCardId)
    && !processedCardIds.includes(rememberedCurrentCardId)
    ? rememberedCurrentCardId
    : triageCardIds.find((cardId) => !processedCardIds.includes(cardId)) || null;

  return {
    currentBoardUrl: boardUrl,
    currentMode: mode,
    currentCardId,
    processedCardIds,
    lastRunAt: now,
  };
}

async function main() {
  const args = parseArgs(process.argv.slice(2));
  const stateFile = path.resolve(args['state-file'] || 'work/planban-board-worker/state/current-board.json');
  const rememberedState = await loadBoardState(stateFile);
  const boardUrl = await resolveActiveBoardUrl({
    cliBoardUrl: args.board || '',
    rememberedState,
  });
  const boardId = boardIdFrom(boardUrl);
  const mode = args.mode || rememberedState?.currentMode || 'auto';
  const outRoot = path.resolve(args.out || 'work/planban-board-worker/runs');
  const account = await resolveAccount({
    cwd: process.cwd(),
    accountFile: args['account-file'],
    env: process.env,
    fallbackAccountFiles: [
      path.join(SKILL_DIR, 'account.local.txt'),
      path.join(SKILL_DIR, 'account.txt'),
    ],
  });
  const token = await login(DEFAULT_API_BASE, account.username, account.password);
  const boardJson = await requestJson(`${DEFAULT_API_BASE}/boards/${boardId}?subscribe=true`, {
    headers: authHeaders(token),
  });
  const snapshot = indexBoardSnapshot(boardJson);
  const listIds = validateBoardLists(snapshot.lists);
  const triageCards = snapshot.cards.filter((card) => card.listId === listIds.triageId);
  const listsById = Object.fromEntries(snapshot.lists.map((item) => [item.id, item]));
  const usersById = Object.fromEntries(snapshot.users.map((item) => [item.id, item]));
  const labelNamesById = Object.fromEntries(snapshot.labels.map((item) => [item.id, item.name]));
  const attachmentsByCardId = new Map();
  for (const attachment of snapshot.attachments) {
    if (!attachmentsByCardId.has(attachment.cardId)) attachmentsByCardId.set(attachment.cardId, []);
    attachmentsByCardId.get(attachment.cardId).push(attachment);
  }
  const labelNamesByCardId = new Map();
  for (const cardLabel of snapshot.cardLabels) {
    if (!labelNamesByCardId.has(cardLabel.cardId)) labelNamesByCardId.set(cardLabel.cardId, []);
    const labelName = labelNamesById[cardLabel.labelId];
    if (labelName) labelNamesByCardId.get(cardLabel.cardId).push(labelName);
  }

  const runDir = path.join(outRoot, `${new Date().toISOString().replace(/[:.]/g, '-')}-${boardId}`);
  const cardsRoot = path.join(runDir, 'cards');
  await mkdir(runDir, { recursive: true });
  const cardEntries = [];
  for (const card of triageCards) {
    const cardDir = path.join(cardsRoot, `${card.id}-${sanitizeName(card.name)}`);
    const attachmentDir = path.join(cardDir, 'attachments');
    await mkdir(attachmentDir, { recursive: true });

    let comments = [];
    try {
      const actions = await requestJson(`${DEFAULT_API_BASE}/cards/${card.id}/actions`, {
        headers: authHeaders(token),
      });
      comments = (actions.items || []).filter((item) => /comment/i.test(item.type || '') || item.data?.text || item.data?.plainText);
    } catch (error) {
      comments = [{
        type: 'fetchError',
        createdAt: null,
        data: { plainText: `Could not fetch comments: ${error.message}` },
      }];
    }

    const downloadedAttachments = [];
    for (const attachment of attachmentsByCardId.get(card.id) || []) {
      downloadedAttachments.push(await downloadAttachment(attachment, token, attachmentDir));
    }

    const issuePath = path.join(cardDir, 'issue.md');
    for (const attachment of downloadedAttachments) attachment.issuePath = issuePath;
    await writeFile(issuePath, buildIssueMarkdown({
      board: snapshot.board,
      card,
      listName: listsById[card.listId]?.name || card.listId,
      labelNames: labelNamesByCardId.get(card.id) || [],
      creatorName: usersById[card.creatorUserId]?.name || card.creatorUserId || 'Unknown',
      comments,
      attachments: downloadedAttachments,
    }), 'utf8');

    cardEntries.push({
      id: card.id,
      name: card.name,
      description: card.description || '',
      issuePath: path.relative(runDir, issuePath),
      commentCount: comments.filter((item) => item.type !== 'fetchError').length,
      attachmentCount: downloadedAttachments.length,
      labels: labelNamesByCardId.get(card.id) || [],
      attachments: downloadedAttachments.map((attachment) => ({
        id: attachment.id,
        name: attachment.name,
        type: attachment.type,
        downloaded: attachment.downloaded,
        path: path.relative(runDir, attachment.path),
      })),
    });
  }

  await writeFile(path.join(runDir, 'summary.md'), buildSummary(snapshot.board, cardEntries), 'utf8');
  await writeFile(path.join(runDir, 'manifest.json'), `${JSON.stringify({
    board: {
      id: snapshot.board.id,
      name: snapshot.board.name,
      url: `${DEFAULT_SITE_BASE}/boards/${snapshot.board.id}`,
    },
    mode,
    lists: listIds,
    triageCardIds: triageCards.map((card) => card.id),
    triageCards: cardEntries,
  }, null, 2)}\n`, 'utf8');

  await saveBoardState(stateFile, deriveStateAfterScan({
    boardUrl,
    mode,
    triageCards,
    rememberedState,
    now: new Date().toISOString(),
  }));

  console.log(`Scanned board ${snapshot.board.name} (${snapshot.board.id})`);
  console.log(`Run directory: ${runDir}`);
}

const isDirectExecution = process.argv[1]
  && import.meta.url === pathToFileURL(path.resolve(process.argv[1])).href;

if (isDirectExecution) {
  main().catch((error) => {
    console.error(`Error: ${error.message}`);
    process.exit(1);
  });
}
