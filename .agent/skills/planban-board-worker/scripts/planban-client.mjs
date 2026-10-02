import { readFile, writeFile } from 'node:fs/promises';
import path from 'node:path';

export const DEFAULT_API_BASE = 'https://planban-api.sjfood.us/api';
export const DEFAULT_SITE_BASE = 'https://planban.sjfood.us';

function isPlaceholderAccount(username, password) {
  const placeholderUsernames = new Set(['username_or_email', 'your_planban_username_or_email']);
  const placeholderPasswords = new Set(['password', 'your_planban_password']);
  return placeholderUsernames.has(String(username).toLowerCase())
    || placeholderPasswords.has(String(password).toLowerCase());
}

async function readAccountFile(filePath) {
  const text = await readFile(filePath, 'utf8');
  const lines = text
    .split(/\r?\n/)
    .map((line) => line.trim())
    .filter((line) => line && !line.startsWith('#'));
  if (lines.length < 2) {
    throw new Error(`Fill in Planban account file: ${filePath}\nExpected two lines:\nusername_or_email\npassword`);
  }
  if (isPlaceholderAccount(lines[0], lines[1])) {
    throw new Error(`Fill in Planban account file: ${filePath}\nReplace the template values with your Planban account and password.`);
  }
  return { username: lines[0], password: lines[1] };
}

async function tryReadAccountFile(filePath) {
  try {
    return await readAccountFile(filePath);
  } catch (error) {
    if (error?.code === 'ENOENT') return null;
    throw error;
  }
}

async function readRootDotEnv(cwd) {
  const envPath = path.resolve(cwd, '.env');
  try {
    const text = await readFile(envPath, 'utf8');
    const values = {};
    for (const rawLine of text.split(/\r?\n/)) {
      const line = rawLine.trim();
      if (!line || line.startsWith('#')) continue;
      const match = line.match(/^(?:export\s+)?([A-Za-z_][A-Za-z0-9_]*)\s*=\s*(.*)$/);
      if (!match) continue;
      let value = match[2].trim();
      if ((value.startsWith('"') && value.endsWith('"')) || (value.startsWith("'") && value.endsWith("'"))) {
        value = value.slice(1, -1);
      }
      values[match[1]] = value;
    }
    return values;
  } catch (error) {
    if (error?.code === 'ENOENT') return {};
    throw error;
  }
}

export function boardIdFrom(input) {
  const text = String(input);
  const match = text.match(/boards\/([0-9]+)/) || text.match(/^([0-9]+)$/);
  if (!match) throw new Error(`Could not parse board id from: ${input}`);
  return match[1];
}

export function normalizeListName(value) {
  return String(value || '').trim().replace(/\s+/g, ' ').toLowerCase();
}

export async function ensureAccountFile(accountPath) {
  await writeFile(accountPath, [
    '# Planban account',
    '# Line 1: Planban username or email',
    '# Line 2: Planban password',
    '',
  ].join('\n'), { encoding: 'utf8', flag: 'wx', mode: 0o600 });
}

export async function resolveAccount({
  cwd,
  accountFile,
  env,
  fallbackAccountFiles = [],
}) {
  if (accountFile) return readAccountFile(accountFile);
  const rootEnv = await readRootDotEnv(cwd);
  const username = env.PLANBAN_USERNAME || rootEnv.PLANBAN_USERNAME || '';
  const password = env.PLANBAN_PASSWORD || rootEnv.PLANBAN_PASSWORD || '';
  if (username && password) {
    return { username, password };
  }

  const candidates = [
    path.resolve(cwd, 'planban-account'),
    path.resolve(cwd, '.planban-account'),
    path.resolve(cwd, 'planban-account.local.txt'),
    ...fallbackAccountFiles,
  ];
  for (const filePath of candidates) {
    const account = await tryReadAccountFile(filePath);
    if (account) return account;
  }

  const accountFilePath = path.resolve(cwd, 'planban-account');
  await ensureAccountFile(accountFilePath);
  throw new Error(`Created Planban account file: ${accountFilePath}\nFill in your Planban account and password, then run the command again.`);
}

export async function requestJson(url, options = {}, timeoutMs = 15000) {
  const controller = new AbortController();
  const timeout = setTimeout(() => controller.abort(), timeoutMs);
  let response;
  try {
    response = await fetch(url, { ...options, signal: controller.signal });
  } finally {
    clearTimeout(timeout);
  }
  const text = await response.text();
  let json;
  try {
    json = text ? JSON.parse(text) : null;
  } catch {
    throw new Error(`Non-JSON response from ${url}: HTTP ${response.status}`);
  }
  if (!response.ok) {
    const message = json?.message || json?.code || `HTTP ${response.status}`;
    throw new Error(`${message} (${url})`);
  }
  return json;
}

export async function login(apiBase, username, password) {
  const json = await requestJson(`${apiBase}/access-tokens?withHttpOnlyToken=false`, {
    method: 'POST',
    headers: { 'content-type': 'application/json' },
    body: JSON.stringify({ emailOrUsername: username, password }),
  });
  if (!json?.item) throw new Error('Login did not return an access token');
  return json.item;
}

export function authHeaders(token, extra = {}) {
  return { authorization: `Bearer ${token}`, ...extra };
}

export function resolveRequiredLists(lists, requiredNames) {
  const byName = new Map(lists.map((list) => [normalizeListName(list.name), list]));
  const missing = requiredNames.filter((name) => !byName.has(normalizeListName(name)));
  if (missing.length) throw new Error(`Missing required lists: ${missing.join(', ')}`);
  return Object.fromEntries(requiredNames.map((name) => [name, byName.get(normalizeListName(name))]));
}

export function indexBoardSnapshot(boardJson) {
  const included = boardJson.included || {};
  return {
    board: boardJson.item,
    lists: included.lists || [],
    cards: included.cards || [],
    labels: included.labels || [],
    cardLabels: included.cardLabels || [],
    attachments: included.attachments || [],
    users: included.users || [],
  };
}
