#!/usr/bin/env node

import { pathToFileURL } from 'node:url';
import path from 'node:path';

import {
  authHeaders,
  DEFAULT_API_BASE,
  login,
  requestJson,
  resolveAccount,
} from './planban-client.mjs';
import { commentContract } from './write-contract.mjs';

export function parseCommentArgs(argv) {
  const args = {
    apply: false,
    card: '',
    text: '',
    board: '',
  };

  for (let i = 0; i < argv.length; i += 1) {
    const key = argv[i];
    if (key === '--apply') {
      args.apply = true;
      continue;
    }

    const value = argv[i + 1];
    if (!value || value.startsWith('--')) throw new Error(`Missing value for ${key}`);

    if (key === '--card') args.card = value;
    else if (key === '--text') args.text = value;
    else if (key === '--board') args.board = value;
    else throw new Error(`Unexpected argument: ${key}`);

    i += 1;
  }

  if (!args.card) throw new Error('Missing value for --card');
  if (!args.text.trim()) throw new Error('Missing value for --text');

  return args;
}

export function buildCommentRequest({ apiBase, cardId, text }) {
  const trimmedText = text.trim();
  return {
    method: commentContract.method,
    url: `${apiBase}${commentContract.pathForCard(cardId)}`,
    body: JSON.stringify(commentContract.buildBody({
      text: trimmedText,
      plainText: trimmedText,
      atUserIds: [],
    })),
  };
}

async function main() {
  const args = parseCommentArgs(process.argv.slice(2));

  if (!args.apply) {
    console.log(`DRY RUN: would comment on card ${args.card}`);
    return;
  }

  const account = await resolveAccount({
    cwd: process.cwd(),
    accountFile: '',
    env: process.env,
  });
  const token = await login(DEFAULT_API_BASE, account.username, account.password);
  const request = buildCommentRequest({
    apiBase: DEFAULT_API_BASE,
    cardId: args.card,
    text: args.text,
  });

  await requestJson(request.url, {
    method: request.method,
    headers: authHeaders(token, { 'content-type': 'application/json' }),
    body: request.body,
  });

  console.log(`Commented on card ${args.card}`);
}

const isDirectExecution = process.argv[1]
  && import.meta.url === pathToFileURL(path.resolve(process.argv[1])).href;

if (isDirectExecution) {
  main().catch((error) => {
    console.error(`Error: ${error.message}`);
    process.exit(1);
  });
}
