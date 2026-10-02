#!/usr/bin/env node

import { pathToFileURL } from 'node:url';
import path from 'node:path';

import {
  authHeaders,
  boardIdFrom,
  DEFAULT_API_BASE,
  indexBoardSnapshot,
  login,
  normalizeListName,
  requestJson,
  resolveAccount,
  resolveRequiredLists,
} from './planban-client.mjs';
import { loadBoardState, saveBoardState } from './board-state.mjs';
import { moveContract } from './write-contract.mjs';

const POSITION_GAP = 65535;
const DEFAULT_STATE_FILE = 'work/planban-board-worker/state/current-board.json';
const CANONICAL_TARGETS = new Map([
  ['triage', 'Triage'],
  ['to do', 'To Do'],
  ['in progress', 'In Progress'],
  ['need to discuss', 'Need To Discuss'],
]);

export function parseMoveArgs(argv) {
  const args = {
    apply: false,
    board: '',
    card: '',
    toList: '',
  };

  for (let i = 0; i < argv.length; i += 1) {
    const key = argv[i];
    if (key === '--apply') {
      args.apply = true;
      continue;
    }

    const value = argv[i + 1];
    if (!value || value.startsWith('--')) throw new Error(`Missing value for ${key}`);

    if (key === '--board') args.board = value;
    else if (key === '--card') args.card = value;
    else if (key === '--to-list') args.toList = value;
    else if (key === '--state-file') args.stateFile = value;
    else throw new Error(`Unexpected argument: ${key}`);

    i += 1;
  }

  if (!args.board) throw new Error('Missing value for --board');
  if (!args.card) throw new Error('Missing value for --card');
  if (!args.toList) throw new Error('Missing value for --to-list');

  return args;
}

export function resolveTargetListName(value) {
  const canonical = CANONICAL_TARGETS.get(normalizeListName(value));
  if (!canonical) throw new Error(`Unsupported target list: ${value}`);
  return canonical;
}

function resolveCurrentListName(value) {
  const canonical = CANONICAL_TARGETS.get(normalizeListName(value));
  if (!canonical) throw new Error(`Unsupported current list: ${value}`);
  return canonical;
}

export function validateTransition(currentListName, targetListName) {
  let current;
  let target;
  try {
    current = resolveCurrentListName(currentListName);
    target = resolveTargetListName(targetListName);
  } catch {
    throw new Error(`Unsupported transition: ${currentListName} -> ${targetListName}`);
  }
  const allowed = (
    (current === 'Triage' && (target === 'To Do' || target === 'Need To Discuss'))
    || (current === 'To Do' && target === 'In Progress')
  );
  if (!allowed) throw new Error(`Unsupported transition: ${current} -> ${target}`);
  return target;
}

export function buildMoveRequest({ apiBase, cardId, listId, position }) {
  return {
    method: moveContract.method,
    url: `${apiBase}${moveContract.pathForCard(cardId)}`,
    body: JSON.stringify(moveContract.buildBody({ listId, position })),
  };
}

function nextCardPosition(cards, targetListId, movingCardId) {
  const siblings = cards
    .filter((card) => card.listId === targetListId && card.id !== movingCardId)
    .sort((left, right) => Number(left.position || 0) - Number(right.position || 0));
  const last = siblings[siblings.length - 1];
  return (last ? Number(last.position || 0) : 0) + POSITION_GAP;
}

export function deriveStateAfterMove({
  rememberedState,
  boardUrl,
  movedCardId,
  fromListName,
  targetListName,
  triageCardsAfterMove,
  now,
}) {
  const processedCardIds = new Set(rememberedState?.processedCardIds || []);
  const normalizedFrom = resolveCurrentListName(fromListName);
  const normalizedTarget = resolveTargetListName(targetListName);
  if (normalizedFrom === 'Triage' && (normalizedTarget === 'To Do' || normalizedTarget === 'Need To Discuss')) {
    processedCardIds.add(movedCardId);
  }

  const triageCardIds = triageCardsAfterMove.map((card) => card.id);
  const rememberedCurrentCardId = rememberedState?.currentCardId || null;
  const currentCardId = rememberedCurrentCardId
    && triageCardIds.includes(rememberedCurrentCardId)
    && !processedCardIds.has(rememberedCurrentCardId)
    ? rememberedCurrentCardId
    : triageCardIds.find((cardId) => !processedCardIds.has(cardId)) || null;

  return {
    currentBoardUrl: boardUrl,
    currentMode: rememberedState?.currentMode || 'auto',
    currentCardId,
    processedCardIds: [...processedCardIds],
    lastRunAt: now,
  };
}

async function main() {
  const args = parseMoveArgs(process.argv.slice(2));
  const stateFile = path.resolve(args.stateFile || DEFAULT_STATE_FILE);
  const rememberedState = await loadBoardState(stateFile);
  const account = await resolveAccount({
    cwd: process.cwd(),
    accountFile: '',
    env: process.env,
  });
  const token = await login(DEFAULT_API_BASE, account.username, account.password);
  const boardJson = await requestJson(`${DEFAULT_API_BASE}/boards/${boardIdFrom(args.board)}?subscribe=true`, {
    headers: authHeaders(token),
  });
  const snapshot = indexBoardSnapshot(boardJson);
  const card = snapshot.cards.find((item) => item.id === args.card);
  if (!card) throw new Error(`Card not found on board: ${args.card}`);
  const currentList = snapshot.lists.find((item) => item.id === card.listId);
  if (!currentList) throw new Error(`Current list not found for card: ${args.card}`);

  const targetListName = validateTransition(currentList.name, args.toList);
  const lists = resolveRequiredLists(snapshot.lists, ['To Do', 'In Progress', 'Need To Discuss']);
  const targetList = lists[targetListName];
  const request = buildMoveRequest({
    apiBase: DEFAULT_API_BASE,
    cardId: args.card,
    listId: targetList.id,
    position: nextCardPosition(snapshot.cards, targetList.id, args.card),
  });

  if (!args.apply) {
    console.log(`DRY RUN: would move card ${args.card} to ${targetListName}`);
    return;
  }

  await requestJson(request.url, {
    method: request.method,
    headers: authHeaders(token, { 'content-type': 'application/json' }),
    body: request.body,
  });

  const refreshedBoardJson = await requestJson(`${DEFAULT_API_BASE}/boards/${boardIdFrom(args.board)}?subscribe=true`, {
    headers: authHeaders(token),
  });
  const refreshedSnapshot = indexBoardSnapshot(refreshedBoardJson);
  const refreshedCard = refreshedSnapshot.cards.find((item) => item.id === args.card);
  if (!refreshedCard) throw new Error(`Moved card missing after refresh: ${args.card}`);
  if (refreshedCard.listId !== targetList.id) {
    throw new Error(`Move verification failed for card ${args.card}`);
  }
  const refreshedLists = resolveRequiredLists(refreshedSnapshot.lists, ['Triage', 'To Do', 'In Progress', 'Need To Discuss']);

  await saveBoardState(stateFile, deriveStateAfterMove({
    rememberedState,
    boardUrl: args.board,
    movedCardId: args.card,
    fromListName: currentList.name,
    targetListName,
    triageCardsAfterMove: refreshedSnapshot.cards.filter((item) => item.listId === refreshedLists.Triage.id),
    now: new Date().toISOString(),
  }));

  console.log(`Moved card ${args.card} to ${targetListName}`);
}

const isDirectExecution = process.argv[1]
  && import.meta.url === pathToFileURL(path.resolve(process.argv[1])).href;

if (isDirectExecution) {
  main().catch((error) => {
    console.error(`Error: ${error.message}`);
    process.exit(1);
  });
}
