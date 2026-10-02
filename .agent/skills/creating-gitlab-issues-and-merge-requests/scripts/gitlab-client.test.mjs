import assert from 'node:assert/strict';
import { mkdtemp, writeFile } from 'node:fs/promises';
import { tmpdir } from 'node:os';
import path from 'node:path';
import test from 'node:test';

import {
  ALLOWED_REPOSITORIES,
  GitLabClient,
  assertAllowedRepository,
  buildOwnershipMarker,
  loadGitLabConfig,
  parseGitLabRemote,
  validateGitLabLocation,
} from './gitlab-client.mjs';

test('configuration comes from process env before the ignored root .env', async () => {
  const cwd = await mkdtemp(path.join(tmpdir(), 'gitlab-client-config-'));
  await writeFile(path.join(cwd, '.env'), [
    'GITLAB_BASE_URL=https://gitlab.example.test',
    'GITLAB_API_TOKEN=file-token',
    '',
  ].join('\n'));

  const config = await loadGitLabConfig({
    cwd,
    env: { GITLAB_API_TOKEN: 'process-token' },
  });

  assert.deepEqual(config, {
    baseUrl: 'https://gitlab.example.test',
    token: 'process-token',
  });
});

test('configuration rejects a missing token without echoing secret candidates', async () => {
  const cwd = await mkdtemp(path.join(tmpdir(), 'gitlab-client-missing-'));
  await writeFile(path.join(cwd, '.env'), 'GITLAB_BASE_URL=https://gitlab.example.test\n');

  await assert.rejects(
    loadGitLabConfig({ cwd, env: {} }),
    /GITLAB_API_TOKEN.*\.env/,
  );
});

test('remote parsing supports HTTPS and SSH origins', () => {
  assert.deepEqual(
    parseGitLabRemote('https://gitlab.sjfood.us/oxygen/base-app.git'),
    { host: 'gitlab.sjfood.us', projectPath: 'oxygen/base-app' },
  );
  assert.deepEqual(
    parseGitLabRemote('git@gitlab.sjfood.us:oxygen/base-engine.git'),
    { host: 'gitlab.sjfood.us', projectPath: 'oxygen/base-engine' },
  );
});

test('GitLab location must be HTTPS and match the repository remote host', () => {
  assert.equal(
    validateGitLabLocation('https://gitlab.sjfood.us/', 'gitlab.sjfood.us'),
    'https://gitlab.sjfood.us',
  );
  assert.throws(
    () => validateGitLabLocation('http://gitlab.sjfood.us', 'gitlab.sjfood.us'),
    /HTTPS/,
  );
  assert.throws(
    () => validateGitLabLocation('https://evil.example', 'gitlab.sjfood.us'),
    /does not match origin host/,
  );
});

test('only the three code repositories are allowed', () => {
  assert.deepEqual(ALLOWED_REPOSITORIES, ['base-app', 'base-engine', 'base-web']);
  for (const name of ALLOWED_REPOSITORIES) assert.doesNotThrow(() => assertAllowedRepository(name));
  assert.throws(() => assertAllowedRepository('easypos-hub'), /not allowed/);
  assert.throws(() => assertAllowedRepository('.'), /not allowed/);
});

test('ownership marker is stable and never contains raw branch delimiters', () => {
  const marker = buildOwnershipMarker({
    projectPath: 'oxygen/base-app',
    sourceBranch: 'feature/a-->b',
    targetBranch: 'release/1.1',
  });
  assert.match(marker, /^<!-- easypos-gitlab-change-request:v1:[A-Za-z0-9_-]+ -->$/);
  assert.equal(marker.includes('feature/a-->b'), false);
});

test('HTTP failures redact tokens even when an upstream response echoes one', async () => {
  const token = 'glpat-super-secret-token';
  const client = new GitLabClient({
    baseUrl: 'https://gitlab.sjfood.us',
    token,
    fetchImpl: async (_url, options) => {
      assert.equal(options.headers['PRIVATE-TOKEN'], token);
      return new Response(JSON.stringify({ message: `rejected ${token}` }), {
        status: 401,
        headers: { 'content-type': 'application/json' },
      });
    },
  });

  await assert.rejects(
    client.request('GET', '/user'),
    (error) => {
      assert.equal(error.message.includes(token), false);
      assert.match(error.message, /\[REDACTED\]/);
      return true;
    },
  );
});

test('GitLab requests reject redirects so private tokens cannot cross origins', async () => {
  const client = new GitLabClient({
    baseUrl: 'https://gitlab.sjfood.us',
    token: 'glpat-private-token',
    fetchImpl: async (_url, options) => {
      assert.equal(options.redirect, 'error');
      return new Response('{}', {
        status: 200,
        headers: { 'content-type': 'application/json' },
      });
    },
  });

  await client.request('GET', '/user');
});
