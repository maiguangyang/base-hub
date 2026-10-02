#!/usr/bin/env node

import { mkdir, writeFile } from 'node:fs/promises';
import path from 'node:path';
import { fileURLToPath } from 'node:url';
import {
  authHeaders,
  boardIdFrom,
  DEFAULT_API_BASE,
  DEFAULT_SITE_BASE,
  login,
  requestJson,
  resolveAccount,
} from '../../planban-board-worker/scripts/planban-client.mjs';

const SKILL_DIR = path.dirname(path.dirname(fileURLToPath(import.meta.url)));

function usage() {
  return `Usage:
  fetch-feedback.mjs --board <board-url-or-id> --label <label> --out <dir> [--account-file <path>] [--api-base <url>]

Credentials:
  PLANBAN_USERNAME and PLANBAN_PASSWORD, --account-file, or an auto-loaded local file:
  planban-account
  .planban-account
  planban-account.local.txt
  .agent/skills/planban-feedback-fixer/account.local.txt

If no credentials are found, this command creates ./planban-account with setup instructions.

Account files are plain text with two non-empty lines:
  username
  password
`;
}

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

function normalizeLabel(value) {
  return String(value || '')
    .trim()
    .toLowerCase()
    .replaceAll('後', '后');
}

function sanitizeName(name) {
  return String(name || 'untitled')
    .replace(/[\\/:*?"<>|#%{}$!'@+`=]/g, '_')
    .replace(/\s+/g, ' ')
    .trim()
    .slice(0, 120) || 'untitled';
}

function section(text, heading) {
  const source = String(text || '').replace(/\r\n/g, '\n');
  const escaped = heading.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
  const re = new RegExp(`\\*\\*${escaped}[^\\n]*\\*\\*\\s*\\n([\\s\\S]*?)(?=\\n\\*\\*|$)`, 'i');
  const match = source.match(re);
  return match ? match[1].trim() : '';
}

function stripTemplateLines(text) {
  return String(text || '')
    .split('\n')
    .map((line) => line.trim())
    .filter((line) => line && !/^步骤[123]\s*$/.test(line) && !/^前提[123]\s*$/.test(line))
    .filter((line) => !/^\(?Summarize the bug encountered concisely\)?$/i.test(line))
    .filter((line) => !/^从上面的选择一个或多个$/.test(line))
    .filter((line) => !/^机型$/.test(line))
    .filter((line) => !/^期望結果緣由$/.test(line))
    .join('\n')
    .trim();
}

async function fetchWithTimeout(url, options = {}, timeoutMs = 8000) {
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
    return { id: attachment.id, name: attachment.name, path: target, downloaded: false, type: attachment.type };
  }

  try {
    const response = await fetchWithTimeout(attachment.url, {
      headers: { authorization: `Bearer ${token}` },
    }, 15000);
    if (!response.ok) throw new Error(`HTTP ${response.status}`);
    const bytes = Buffer.from(await response.arrayBuffer());
    await writeFile(target, bytes);
    return { id: attachment.id, name: attachment.name, path: target, downloaded: true, type: attachment.type };
  } catch (error) {
    await writeFile(`${target}.url`, `${attachment.url}\n# download failed: ${error.message}\n`, 'utf8');
    return { id: attachment.id, name: attachment.name, path: `${target}.url`, downloaded: false, type: attachment.type };
  }
}

function makeIssueMarkdown({ board, card, listName, labelNames, creatorName, comments, attachments }) {
  const problem = stripTemplateLines(section(card.description, '問題表現')) || stripTemplateLines(section(card.description, 'Problem')) || '';
  const expected = stripTemplateLines(section(card.description, '期待的正确结果或表达')) || stripTemplateLines(section(card.description, '期待的正確結果')) || '';
  const summary = stripTemplateLines(section(card.description, 'Summary')) || '';

  const attachmentLines = attachments.length
    ? attachments.map((a) => `- ${path.relative(path.dirname(a.issuePath), a.path)}${a.downloaded ? '' : ' (link only)'}`).join('\n')
    : 'None';
  const commentLines = comments.length
    ? comments.map((c) => `- ${c.createdAt || 'unknown'}: ${c.data?.plainText || c.data?.text || JSON.stringify(c.data)}`).join('\n')
    : 'None';

  return `# ${card.name}

Board: ${board.name}
Card: ${card.id}
List: ${listName}
Labels: ${labelNames.join(', ') || 'None'}
Creator: ${creatorName}
Created: ${card.createdAt}
Updated: ${card.updatedAt || 'None'}

## Summary
${summary || card.name}

## Problem
${problem || 'Not provided. Infer narrowly from title, description, and attachments.'}

## Expected
${expected || 'Not provided. Clarify if the correct behavior is not obvious from the screenshot/title.'}

## Full Description
${card.description || 'None'}

## Comments
${commentLines}

## Attachments
${attachmentLines}

## Repair Notes
Pending.

## Verification
Pending.
`;
}

async function main() {
  const args = parseArgs(process.argv.slice(2));
  if (!args.board || !args.label || !args.out) {
    console.error(usage());
    process.exit(2);
  }

  const apiBase = (args['api-base'] || DEFAULT_API_BASE).replace(/\/$/, '');
  const boardId = boardIdFrom(args.board);
  const account = await resolveAccount({
    cwd: process.cwd(),
    accountFile: args['account-file'],
    env: process.env,
    fallbackAccountFiles: [
      path.join(SKILL_DIR, 'account.local.txt'),
      path.join(SKILL_DIR, 'account.txt'),
    ],
  });

  if (!account.username || !account.password) {
    throw new Error('Missing credentials. Set PLANBAN_USERNAME/PLANBAN_PASSWORD or fill in planban-account.');
  }

  const token = await login(apiBase, account.username, account.password);
  const boardJson = await requestJson(`${apiBase}/boards/${boardId}?subscribe=true`, {
    headers: authHeaders(token),
  });

  const board = boardJson.item;
  const inc = boardJson.included || {};
  const labels = inc.labels || [];
  const listsById = Object.fromEntries((inc.lists || []).map((item) => [item.id, item]));
  const usersById = Object.fromEntries((inc.users || []).map((item) => [item.id, item]));
  const label = labels.find((item) => item.name === args.label)
    || labels.find((item) => normalizeLabel(item.name) === normalizeLabel(args.label));
  if (!label) {
    throw new Error(`Label not found: ${args.label}. Available labels: ${labels.map((item) => item.name).join(', ')}`);
  }

  const labelNamesById = Object.fromEntries(labels.map((item) => [item.id, item.name]));
  const cardLabels = inc.cardLabels || [];
  const selectedCardIds = new Set(cardLabels.filter((item) => item.labelId === label.id).map((item) => item.cardId));
  const cards = (inc.cards || []).filter((card) => selectedCardIds.has(card.id));
  const attachmentsByCardId = new Map();
  for (const attachment of inc.attachments || []) {
    if (!selectedCardIds.has(attachment.cardId)) continue;
    if (!attachmentsByCardId.has(attachment.cardId)) attachmentsByCardId.set(attachment.cardId, []);
    attachmentsByCardId.get(attachment.cardId).push(attachment);
  }

  const outDir = path.resolve(args.out);
  await mkdir(outDir, { recursive: true });

  const manifestCards = [];
  const summaryLines = [`# Planban Feedback Export`, '', `Board: ${board.name}`, `Board ID: ${board.id}`, `Label: ${label.name}`, `Cards: ${cards.length}`, '', '## Cards'];

  for (const card of cards) {
    const listName = listsById[card.listId]?.name || card.listId;
    const creatorName = usersById[card.creatorUserId]?.name || card.creatorUserId || 'Unknown';
    const labelNames = cardLabels.filter((item) => item.cardId === card.id).map((item) => labelNamesById[item.labelId]).filter(Boolean);
    const cardDir = path.join(outDir, `${card.id}-${sanitizeName(card.name)}`);
    const attachmentDir = path.join(cardDir, 'attachments');
    await mkdir(attachmentDir, { recursive: true });

    let comments = [];
    try {
      const actions = await requestJson(`${apiBase}/cards/${card.id}/actions`, {
        headers: authHeaders(token),
      });
      comments = (actions.items || []).filter((item) => /comment/i.test(item.type || '') || item.data?.text || item.data?.plainText);
    } catch (error) {
      comments = [{ type: 'fetchError', createdAt: null, data: { plainText: `Could not fetch comments: ${error.message}` } }];
    }

    const downloaded = [];
    for (const attachment of attachmentsByCardId.get(card.id) || []) {
      downloaded.push(await downloadAttachment(attachment, token, attachmentDir));
    }

    const issuePath = path.join(cardDir, 'issue.md');
    for (const attachment of downloaded) attachment.issuePath = issuePath;
    await writeFile(issuePath, makeIssueMarkdown({
      board,
      card,
      listName,
      labelNames,
      creatorName,
      comments,
      attachments: downloaded,
    }), 'utf8');

    const cardEntry = {
      id: card.id,
      url: `${DEFAULT_SITE_BASE}/cards/${card.id}`,
      name: card.name,
      list: listName,
      labels: labelNames,
      creator: creatorName,
      createdAt: card.createdAt,
      updatedAt: card.updatedAt,
      issuePath,
      attachmentCount: downloaded.length,
      commentCount: comments.filter((item) => item.type !== 'fetchError').length,
    };
    manifestCards.push(cardEntry);
    summaryLines.push(`- [${card.name}](${path.relative(outDir, issuePath)}) - ${listName} - ${card.id}`);
  }

  const manifest = {
    exportedAt: new Date().toISOString(),
    board: { id: board.id, name: board.name, url: `${DEFAULT_SITE_BASE}/boards/${board.id}` },
    label: { id: label.id, name: label.name },
    counts: { cards: manifestCards.length },
    cards: manifestCards,
  };
  await writeFile(path.join(outDir, 'manifest.json'), JSON.stringify(manifest, null, 2), 'utf8');
  await writeFile(path.join(outDir, 'summary.md'), `${summaryLines.join('\n')}\n`, 'utf8');

  console.log(`Exported ${manifestCards.length} card(s) to ${outDir}`);
  console.log(`Summary: ${path.join(outDir, 'summary.md')}`);
}

main().catch((error) => {
  console.error(`Error: ${error.message}`);
  process.exit(1);
});
