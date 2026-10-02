import { readFile } from 'node:fs/promises';
import path from 'node:path';

export const ALLOWED_REPOSITORIES = Object.freeze([
  'base-app',
  'base-engine',
  'base-web',
]);

export const EXPECTED_PROJECT_PATHS = Object.freeze({
  'base-app': 'oxygen/base-app',
  'base-engine': 'oxygen/base-engine',
  'base-web': 'oxygen/base-web',
});

async function readDotEnv(cwd) {
  try {
    const text = await readFile(path.resolve(cwd, '.env'), 'utf8');
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

export async function loadGitLabConfig({ cwd, env = process.env }) {
  const fileValues = await readDotEnv(cwd);
  const baseUrl = String(env.GITLAB_BASE_URL || fileValues.GITLAB_BASE_URL || '').trim();
  const token = String(env.GITLAB_API_TOKEN || fileValues.GITLAB_API_TOKEN || '').trim();
  if (!baseUrl) throw new Error('Set GITLAB_BASE_URL in the root .env or process environment.');
  if (!token) throw new Error('Set GITLAB_API_TOKEN in the ignored root .env or process environment.');
  return { baseUrl: baseUrl.replace(/\/+$/, ''), token };
}

export function parseGitLabRemote(remote) {
  const value = String(remote || '').trim();
  let host;
  let projectPath;

  if (/^https?:\/\//i.test(value) || /^ssh:\/\//i.test(value)) {
    const url = new URL(value);
    host = url.hostname;
    projectPath = url.pathname.replace(/^\/+|\/+$/g, '');
  } else {
    const match = value.match(/^(?:[^@/]+@)?([^:/]+):(.+)$/);
    if (!match) throw new Error(`Unsupported Git remote format: ${value}`);
    [, host, projectPath] = match;
  }

  projectPath = projectPath.replace(/\.git$/, '').replace(/^\/+|\/+$/g, '');
  if (!host || !projectPath || !projectPath.includes('/')) {
    throw new Error(`Git remote does not identify a GitLab project: ${value}`);
  }
  return { host: host.toLowerCase(), projectPath };
}

export function validateGitLabLocation(baseUrl, remoteHost) {
  const url = new URL(baseUrl);
  if (url.protocol !== 'https:') throw new Error('GITLAB_BASE_URL must use HTTPS.');
  if (url.username || url.password || url.search || url.hash) {
    throw new Error('GITLAB_BASE_URL must not contain credentials, query parameters, or fragments.');
  }
  if (url.hostname.toLowerCase() !== String(remoteHost).toLowerCase()) {
    throw new Error(`GITLAB_BASE_URL host ${url.hostname} does not match origin host ${remoteHost}.`);
  }
  return url.toString().replace(/\/+$/, '');
}

export function assertAllowedRepository(repository) {
  if (!ALLOWED_REPOSITORIES.includes(repository)) {
    throw new Error(`Repository is not allowed for GitLab Issue/MR publication: ${repository}`);
  }
  return repository;
}

export function assertExpectedProjectPath(repository, projectPath) {
  assertAllowedRepository(repository);
  const expected = EXPECTED_PROJECT_PATHS[repository];
  if (projectPath !== expected) {
    throw new Error(`${repository}: origin project path must be ${expected}, received ${projectPath}.`);
  }
  return projectPath;
}

export function buildOwnershipMarker({ projectPath, sourceBranch, targetBranch }) {
  const payload = Buffer.from(JSON.stringify({ projectPath, sourceBranch, targetBranch }), 'utf8')
    .toString('base64url');
  return `<!-- easypos-gitlab-change-request:v1:${payload} -->`;
}

function redact(text, token) {
  return String(text).split(token).join('[REDACTED]');
}

export class GitLabClient {
  constructor({ baseUrl, token, fetchImpl = globalThis.fetch, timeoutMs = 20_000 }) {
    this.baseUrl = String(baseUrl).replace(/\/+$/, '');
    this.token = token;
    this.fetchImpl = fetchImpl;
    this.timeoutMs = timeoutMs;
  }

  async request(method, apiPath, { query, body } = {}) {
    const url = new URL(`${this.baseUrl}/api/v4${apiPath}`);
    for (const [key, value] of Object.entries(query || {})) {
      if (value !== undefined && value !== null) url.searchParams.set(key, String(value));
    }

    const controller = new AbortController();
    const timeout = setTimeout(() => controller.abort(), this.timeoutMs);
    let response;
    try {
      response = await this.fetchImpl(url, {
        method,
        redirect: 'error',
        signal: controller.signal,
        headers: {
          'PRIVATE-TOKEN': this.token,
          ...(body === undefined ? {} : { 'content-type': 'application/json' }),
        },
        ...(body === undefined ? {} : { body: JSON.stringify(body) }),
      });
    } catch (error) {
      throw new Error(redact(`GitLab ${method} request failed: ${error?.message || error}`, this.token));
    } finally {
      clearTimeout(timeout);
    }

    const text = await response.text();
    let data = null;
    if (text) {
      try {
        data = JSON.parse(text);
      } catch {
        if (response.ok) throw new Error(`GitLab ${method} returned a non-JSON response (HTTP ${response.status}).`);
      }
    }
    if (!response.ok) {
      const upstream = typeof data?.message === 'string' ? data.message : `HTTP ${response.status}`;
      throw new Error(redact(`GitLab ${method} failed: ${upstream} (HTTP ${response.status})`, this.token));
    }
    return data;
  }
}
