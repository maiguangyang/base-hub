import assert from 'node:assert/strict';
import { mkdtemp, writeFile } from 'node:fs/promises';
import { tmpdir } from 'node:os';
import path from 'node:path';
import test from 'node:test';

import {
  GitInspector,
  GitLabApi,
  normalizeIssueRecordManifest,
  normalizeManifest,
  preflightBatch,
  preflightIssueRecordBatch,
  publishBatch,
  publishIssueRecords,
  runCli,
} from './publish-change-requests.mjs';

const repositoryStates = Object.freeze({
  'base-app': { projectPath: 'oxygen/base-app', branch: 'feature/release', localSha: 'app-source', localTargetSha: 'target' },
  'base-engine': { projectPath: 'oxygen/base-engine', branch: 'feature/release', localSha: 'engine-source', localTargetSha: 'target' },
  'base-web': { projectPath: 'oxygen/base-web', branch: 'feature/release', localSha: 'web-source', localTargetSha: 'target' },
});

function createGitInspector(overrides = {}) {
  return {
    inspect: async (repository, requestedSource) => ({
      repository,
      dirty: false,
      ...repositoryStates[repository],
      branch: requestedSource === 'CURRENT' ? repositoryStates[repository].branch : requestedSource,
      ...overrides[repository],
    }),
    inspectIssueRecord: async (repository, source, target) => ({
      repository,
      projectPath: repositoryStates[repository].projectPath,
      remoteHost: 'gitlab.sjfood.us',
      sourceSha: repositoryStates[repository].localSha,
      targetSha: repositoryStates[repository].localTargetSha,
      sourceMerged: false,
      source,
      target,
      ...overrides[repository],
    }),
  };
}

function createApi(overrides = {}) {
  const writes = [];
  const api = {
    writes,
    getProject: async (projectPath) => ({ id: projectPath, permissions: { project_access: { access_level: 30 } } }),
    getCommit: async (_projectPath, sha) => ({ id: sha }),
    getBranch: async (_projectPath, branch) => ({ name: branch, commit: { id: `${branch === 'release/1.1' ? 'target' : 'source'}` } }),
    compare: async (projectPath) => ({
      commits: projectPath.endsWith('web') ? [] : [{ id: 'change' }],
      diffs: projectPath.endsWith('web') ? [] : [{ new_path: 'changed-file' }],
    }),
    listOpenMergeRequests: async () => [],
    listMergedMergeRequests: async () => [],
    findOpenIssueByMarker: async () => null,
    findIssueByMarker: async () => null,
    createIssue: async (projectPath, payload) => {
      writes.push({ kind: 'issue', projectPath, payload });
      return { iid: writes.length, web_url: `https://gitlab.test/${projectPath}/issues/${writes.length}` };
    },
    updateIssue: async (projectPath, iid, payload) => {
      writes.push({ kind: 'issue-update', projectPath, iid, payload });
      return { iid, web_url: `https://gitlab.test/${projectPath}/issues/${iid}` };
    },
    createMergeRequest: async (projectPath, payload) => {
      writes.push({ kind: 'mr', projectPath, payload });
      return { iid: writes.length, web_url: `https://gitlab.test/${projectPath}/merge_requests/${writes.length}` };
    },
    updateMergeRequest: async (projectPath, iid, payload) => {
      writes.push({ kind: 'mr-update', projectPath, iid, payload });
      return { iid, web_url: `https://gitlab.test/${projectPath}/merge_requests/${iid}` };
    },
    ...overrides,
  };
  return api;
}

function manifest() {
  const changelog = (name, item) => [
    `# ${name}`,
    '',
    '## 新增功能',
    `- ${item}`,
    '',
    '## 功能优化',
    '- 无',
    '',
    '## 问题修复',
    '- 无',
    '',
    '## 升级注意事项',
    '- 无',
    '',
    '## 测试验证清单',
    '- [ ] [P0] 验证发布结果',
  ].join('\n');
  return {
    version: 1,
    source: 'CURRENT',
    target: 'release/1.1',
    repositories: {
      'base-app': {
        title: 'App release changes',
        changelog: changelog('App', 'A'),
        sourceBranch: repositoryStates['base-app'].branch,
        sourceSha: repositoryStates['base-app'].localSha,
        targetBranch: 'release/1.1',
        targetSha: repositoryStates['base-app'].localTargetSha,
      },
      'base-engine': {
        title: 'Engine release changes',
        changelog: changelog('Engine', 'B'),
        sourceBranch: repositoryStates['base-engine'].branch,
        sourceSha: repositoryStates['base-engine'].localSha,
        targetBranch: 'release/1.1',
        targetSha: repositoryStates['base-engine'].localTargetSha,
      },
    },
  };
}

test('manifest rejects root and unknown repositories', () => {
  assert.throws(
    () => normalizeManifest({ ...manifest(), repositories: { '.': { title: 'x', changelog: 'y' } } }),
    /not allowed/,
  );
  assert.throws(
    () => normalizeManifest({ ...manifest(), repositories: { 'easypos-hub': { title: 'x', changelog: 'y' } } }),
    /not allowed/,
  );
});

test('manifest requires the five fixed changelog sections', () => {
  assert.throws(
    () => normalizeManifest({
      ...manifest(),
      repositories: {
        'base-app': { title: 'Incomplete', changelog: '## 新增功能\n- A' },
      },
    }),
    /missing changelog section.*功能优化/i,
  );
});

test('GitLab searches all actors when detecting existing issues and MRs', async () => {
  const calls = [];
  const api = new GitLabApi({
    request: async (method, apiPath, options) => {
      calls.push({ method, apiPath, options });
      return [];
    },
  });

  await api.listOpenMergeRequests('oxygen/base-app', 'feature/a', 'release/1.1');
  await api.listMergedMergeRequests('oxygen/base-app', 'feature/a', 'release/1.1');
  await api.findOpenIssueByMarker('oxygen/base-app', '<!-- marker -->');
  await api.findIssueByMarker('oxygen/base-app', '<!-- marker -->');

  assert.equal(calls[0].options.query.scope, 'all');
  assert.equal(calls[1].options.query.state, 'merged');
  assert.equal(calls[1].options.query.scope, 'all');
  assert.equal(calls[2].options.query.scope, 'all');
  assert.equal(calls[3].options.query.scope, 'all');
  assert.equal(calls[3].options.query.state, 'all');
});

test('preflight inspects all repositories and skips a repository with no remote difference', async () => {
  const api = createApi({
    getBranch: async (projectPath, branch) => ({
      name: branch,
      commit: { id: branch === 'release/1.1' ? 'target' : repositoryStates[projectPath.split('/').at(-1)].localSha },
    }),
  });
  const result = await preflightBatch({
    source: 'CURRENT',
    target: 'release/1.1',
    git: createGitInspector(),
    api,
  });

  assert.deepEqual(result.repositories.map(({ repository, status }) => [repository, status]), [
    ['base-app', 'ready'],
    ['base-engine', 'ready'],
    ['base-web', 'skipped'],
  ]);
});

test('pushed manifest rejects source or target revision drift after changelog generation', async () => {
  const api = createApi({
    getBranch: async (projectPath, branch) => ({
      name: branch,
      commit: { id: branch === 'release/1.1' ? 'target' : repositoryStates[projectPath.split('/').at(-1)].localSha },
    }),
  });
  const preflight = await preflightBatch({
    source: 'CURRENT',
    target: 'release/1.1',
    git: createGitInspector(),
    api,
  });
  const drifted = {
    ...preflight,
    repositories: preflight.repositories.map((state) => state.repository === 'base-app'
      ? { ...state, branch: 'feature/other', localSha: 'different-source-sha' }
      : state),
  };

  await assert.rejects(
    publishBatch({ manifest: manifest(), preflight: drifted, api, apply: false }),
    /base-app.*source branch or SHA.*changelog manifest/i,
  );
  const targetDrifted = {
    ...preflight,
    repositories: preflight.repositories.map((state) => state.repository === 'base-app'
      ? { ...state, targetBranch: 'release/other', localTargetSha: 'different-target-sha' }
      : state),
  };
  await assert.rejects(
    publishBatch({ manifest: manifest(), preflight: targetDrifted, api, apply: false }),
    /base-app.*target branch or SHA.*changelog manifest/i,
  );
  assert.equal(api.writes.length, 0);
});

test('preflight rejects dirty and unpushed repositories before publication', async () => {
  const dirtyGit = createGitInspector({ 'base-engine': { dirty: true } });
  await assert.rejects(
    preflightBatch({ source: 'CURRENT', target: 'release/1.1', git: dirtyGit, api: createApi() }),
    /base-engine.*uncommitted/i,
  );

  const api = createApi({
    getBranch: async (projectPath, branch) => ({
      name: branch,
      commit: { id: projectPath.endsWith('app') && branch !== 'release/1.1' ? 'different' : repositoryStates[projectPath.split('/').at(-1)]?.localSha || 'target' },
    }),
  });
  await assert.rejects(
    preflightBatch({ source: 'CURRENT', target: 'release/1.1', git: createGitInspector(), api }),
    /base-app.*not pushed/i,
  );
});

test('Git inspector resolves the target branch to an exact local commit', async () => {
  const inspector = new GitInspector({ rootDir: '/unused' });
  const calls = [];
  inspector.git = async (_repository, args) => {
    calls.push(args);
    if (args[0] === 'symbolic-ref') return 'feature/release';
    if (args[0] === 'check-ref-format') return '';
    if (args[0] === 'status') return '';
    if (args[0] === 'remote') return 'git@gitlab.sjfood.us:oxygen/base-app.git';
    if (args.at(-1) === 'refs/heads/feature/release^{commit}') return 'source-sha';
    if (args.at(-1) === 'refs/heads/release/1.1^{commit}') return 'target-sha';
    throw new Error(`Unexpected git arguments: ${args.join(' ')}`);
  };

  const state = await inspector.inspect('base-app', 'CURRENT', 'release/1.1');

  assert.equal(state.localTargetSha, 'target-sha');
  assert.equal(state.localTargetRef, 'refs/heads/release/1.1');
  assert.equal(calls.some((args) => args.at(-1) === 'refs/heads/release/1.1^{commit}'), true);
});

test('pushed-branch Git checks have a timeout and disable interactive credential prompts', async () => {
  const calls = [];
  const inspector = new GitInspector({
    rootDir: '/unused',
    timeoutMs: 2_345,
    execGit: async (command, args, options) => {
      calls.push({ command, args, options });
      return { stdout: 'ok\n' };
    },
  });

  assert.equal(await inspector.git('base-app', ['status']), 'ok');
  assert.equal(calls[0].options.timeout, 2_345);
  assert.equal(calls[0].options.env.GIT_TERMINAL_PROMPT, '0');
});

test('preflight rejects a stale local target branch before publication', async () => {
  const api = createApi({
    getBranch: async (projectPath, branch) => ({
      name: branch,
      commit: {
        id: branch === 'release/1.1'
          ? 'remote-target'
          : repositoryStates[projectPath.split('/').at(-1)].localSha,
      },
    }),
  });

  await assert.rejects(
    preflightBatch({ source: 'CURRENT', target: 'release/1.1', git: createGitInspector(), api }),
    /base-app.*target branch.*local.*remote/i,
  );
  assert.equal(api.writes.length, 0);
});

test('preflight rejects a repository whose origin points at the wrong GitLab project path', async () => {
  const api = createApi({
    getBranch: async (projectPath, branch) => ({
      name: branch,
      commit: {
        id: branch === 'release/1.1'
          ? 'target'
          : repositoryStates[projectPath.split('/').at(-1)].localSha,
      },
    }),
  });

  await assert.rejects(
    preflightBatch({
      source: 'CURRENT',
      target: 'release/1.1',
      git: createGitInspector({
        'base-app': { projectPath: 'unrelated/base-app' },
      }),
      api,
    }),
    /base-app.*oxygen\/base-app/i,
  );
  assert.equal(api.writes.length, 0);
});

test('preflight skips reverted commits when the final comparison has no file differences', async () => {
  const api = createApi({
    getBranch: async (projectPath, branch) => ({
      name: branch,
      commit: { id: branch === 'release/1.1' ? 'target' : repositoryStates[projectPath.split('/').at(-1)].localSha },
    }),
    compare: async () => ({ commits: [{ id: 'reverted-change' }], diffs: [] }),
  });

  const result = await preflightBatch({
    source: 'CURRENT',
    target: 'release/1.1',
    git: createGitInspector(),
    api,
  });

  assert.deepEqual(result.repositories.map(({ status }) => status), ['skipped', 'skipped', 'skipped']);
});

test('preflight skips an exact source revision already represented by a merged MR', async () => {
  const api = createApi({
    getBranch: async (projectPath, branch) => ({
      name: branch,
      commit: { id: branch === 'release/1.1' ? 'target' : repositoryStates[projectPath.split('/').at(-1)].localSha },
    }),
    listMergedMergeRequests: async (projectPath) => projectPath.endsWith('app')
      ? [{ iid: 71, state: 'merged', diff_refs: { head_sha: repositoryStates['base-app'].localSha } }]
      : [],
  });

  const result = await preflightBatch({
    source: 'CURRENT',
    target: 'release/1.1',
    git: createGitInspector(),
    api,
  });

  assert.equal(result.repositories[0].status, 'skipped');
  assert.match(result.repositories[0].reason, /already merged/i);
  assert.equal(api.writes.length, 0);
});

test('preflight permits new commits added after an older MR from the same branch was merged', async () => {
  const api = createApi({
    getBranch: async (projectPath, branch) => ({
      name: branch,
      commit: { id: branch === 'release/1.1' ? 'target' : repositoryStates[projectPath.split('/').at(-1)].localSha },
    }),
    listMergedMergeRequests: async (projectPath) => projectPath.endsWith('app')
      ? [{ iid: 70, state: 'merged', diff_refs: { head_sha: 'older-source-sha' } }]
      : [],
  });

  const result = await preflightBatch({
    source: 'CURRENT',
    target: 'release/1.1',
    git: createGitInspector(),
    api,
  });

  assert.equal(result.repositories[0].status, 'ready');
});

test('preflight rejects an incomplete comparison even when the exact source revision has a merged MR', async () => {
  const api = createApi({
    getBranch: async (projectPath, branch) => ({
      name: branch,
      commit: { id: branch === 'release/1.1' ? 'target' : repositoryStates[projectPath.split('/').at(-1)].localSha },
    }),
    compare: async (projectPath) => projectPath.endsWith('app')
      ? { compare_timeout: true, diffs: [] }
      : {
        compare_timeout: false,
        diffs: projectPath.endsWith('web') ? [] : [{ new_path: 'changed-file' }],
      },
    listMergedMergeRequests: async (projectPath) => projectPath.endsWith('app')
      ? [{ iid: 71, state: 'merged', diff_refs: { head_sha: repositoryStates['base-app'].localSha } }]
      : [],
  });

  await assert.rejects(
    preflightBatch({ source: 'CURRENT', target: 'release/1.1', git: createGitInspector(), api }),
    /base-app.*compare.*timed out/i,
  );
  assert.equal(api.writes.length, 0);
});

test('preflight fails safely when GitLab reports a compare timeout', async () => {
  const api = createApi({
    getBranch: async (projectPath, branch) => ({
      name: branch,
      commit: { id: branch === 'release/1.1' ? 'target' : repositoryStates[projectPath.split('/').at(-1)].localSha },
    }),
    compare: async () => ({
      compare_timeout: true,
      commits: [{ id: 'change' }],
      diffs: [{ new_path: 'possibly-incomplete' }],
    }),
  });

  await assert.rejects(
    preflightBatch({ source: 'CURRENT', target: 'release/1.1', git: createGitInspector(), api }),
    /base-app.*compare.*timed out/i,
  );
  assert.equal(api.writes.length, 0);
});

test('preflight rejects a malformed GitLab comparison without a diffs array', async () => {
  const api = createApi({
    getBranch: async (projectPath, branch) => ({
      name: branch,
      commit: { id: branch === 'release/1.1' ? 'target' : repositoryStates[projectPath.split('/').at(-1)].localSha },
    }),
    compare: async () => ({ commits: [] }),
  });

  await assert.rejects(
    preflightBatch({ source: 'CURRENT', target: 'release/1.1', git: createGitInspector(), api }),
    /base-app.*diffs array/i,
  );
  assert.equal(api.writes.length, 0);
});

test('foreign open MR collision fails before any write', async () => {
  const api = createApi({
    getBranch: async (projectPath, branch) => ({
      name: branch,
      commit: { id: branch === 'release/1.1' ? 'target' : repositoryStates[projectPath.split('/').at(-1)].localSha },
    }),
    listOpenMergeRequests: async (projectPath) => projectPath.endsWith('engine')
      ? [{ iid: 9, description: 'written manually' }]
      : [],
  });

  await assert.rejects(
    preflightBatch({ source: 'CURRENT', target: 'release/1.1', git: createGitInspector(), api }),
    /existing open MR.*not owned/i,
  );
  assert.equal(api.writes.length, 0);
});

test('remote branch failures identify the affected repository before any write', async () => {
  const api = createApi({
    getBranch: async (projectPath, branch) => {
      if (projectPath.endsWith('engine') && branch !== 'release/1.1') {
        throw new Error('source branch missing');
      }
      return {
        name: branch,
        commit: { id: branch === 'release/1.1' ? 'target' : repositoryStates[projectPath.split('/').at(-1)].localSha },
      };
    },
  });

  await assert.rejects(
    preflightBatch({ source: 'CURRENT', target: 'release/1.1', git: createGitInspector(), api }),
    /base-engine.*source branch missing/i,
  );
  assert.equal(api.writes.length, 0);
});

test('preview performs no writes and apply creates issue before MR with automatic closure link', async () => {
  const api = createApi({
    getBranch: async (projectPath, branch) => ({
      name: branch,
      commit: { id: branch === 'release/1.1' ? 'target' : repositoryStates[projectPath.split('/').at(-1)].localSha },
    }),
  });
  const preflight = await preflightBatch({ source: 'CURRENT', target: 'release/1.1', git: createGitInspector(), api });

  const preview = await publishBatch({ manifest: manifest(), preflight, api, apply: false });
  assert.equal(api.writes.length, 0);
  assert.deepEqual(preview.results.map((item) => item.status), ['preview', 'preview', 'skipped']);

  const applied = await publishBatch({ manifest: manifest(), preflight, api, apply: true });
  assert.deepEqual(api.writes.map(({ kind }) => kind), ['issue', 'mr', 'issue', 'mr']);
  assert.match(api.writes[1].payload.description, /^Closes #1\n/);
  assert.match(api.writes[1].payload.description, /easypos-gitlab-change-request/);
  assert.equal(applied.results[2].status, 'skipped');
});

test('owned issue and MR are reused and updated on rerun', async () => {
  let marker = '';
  const api = createApi({
    getBranch: async (projectPath, branch) => ({
      name: branch,
      commit: { id: branch === 'release/1.1' ? 'target' : repositoryStates[projectPath.split('/').at(-1)].localSha },
    }),
    findIssueByMarker: async (_projectPath, candidate) => {
      marker = candidate;
      return { iid: 41, state: 'opened', description: `old\n\n${candidate}` };
    },
    listOpenMergeRequests: async (projectPath) => projectPath.endsWith('app')
      ? [{ iid: 42, description: `old\n\n${marker}` }]
      : [],
  });
  const preflight = await preflightBatch({ source: 'CURRENT', target: 'release/1.1', git: createGitInspector(), api });
  await publishBatch({ manifest: manifest(), preflight, api, apply: true });

  assert.equal(api.writes.some(({ kind }) => kind === 'issue-update'), true);
  assert.equal(api.writes.some(({ kind }) => kind === 'mr-update'), true);
});

test('closed owned Issue blocks pushed-branch publication without creating a duplicate', async () => {
  const api = createApi({
    getBranch: async (projectPath, branch) => ({
      name: branch,
      commit: { id: branch === 'release/1.1' ? 'target' : repositoryStates[projectPath.split('/').at(-1)].localSha },
    }),
    findIssueByMarker: async (_projectPath, marker) => ({
      iid: 91,
      state: 'closed',
      web_url: 'https://gitlab.test/issues/91',
      description: `old\n\n${marker}`,
    }),
  });

  await assert.rejects(
    preflightBatch({ source: 'CURRENT', target: 'release/1.1', git: createGitInspector(), api }),
    /base-app.*owned Issue #91.*closed.*reopen/is,
  );
  assert.equal(api.writes.length, 0);
});

test('MR failure reports the already-created issue so rerun can recover safely', async () => {
  const api = createApi({
    getBranch: async (projectPath, branch) => ({
      name: branch,
      commit: { id: branch === 'release/1.1' ? 'target' : repositoryStates[projectPath.split('/').at(-1)].localSha },
    }),
    createMergeRequest: async () => {
      throw new Error('temporary GitLab outage');
    },
  });
  const preflight = await preflightBatch({ source: 'CURRENT', target: 'release/1.1', git: createGitInspector(), api });

  await assert.rejects(
    publishBatch({ manifest: manifest(), preflight, api, apply: true }),
    /base-app.*Issue #1.*temporary GitLab outage/i,
  );
  assert.deepEqual(api.writes.map(({ kind }) => kind), ['issue']);
});

test('later repository failure reports objects already completed in earlier repositories', async () => {
  const api = createApi({
    getBranch: async (projectPath, branch) => ({
      name: branch,
      commit: { id: branch === 'release/1.1' ? 'target' : repositoryStates[projectPath.split('/').at(-1)].localSha },
    }),
    createMergeRequest: async (projectPath, payload) => {
      if (projectPath.endsWith('engine')) throw new Error('engine MR failed');
      api.writes.push({ kind: 'mr', projectPath, payload });
      return { iid: 2, web_url: 'https://gitlab.test/oxygen/base-app/merge_requests/2' };
    },
  });
  const preflight = await preflightBatch({ source: 'CURRENT', target: 'release/1.1', git: createGitInspector(), api });

  await assert.rejects(
    publishBatch({ manifest: manifest(), preflight, api, apply: true }),
    (error) => {
      assert.match(error.message, /base-engine.*engine MR failed/i);
      assert.match(error.message, /base-app.*merge_requests\/2/i);
      return true;
    },
  );
});

test('local CLI rejects missing required values, confirmation, and unknown arguments before writes', async () => {
  await assert.rejects(runCli(['local-preflight']), /--target|Usage/i);
  await assert.rejects(runCli(['local-publish']), /--manifest|Usage/i);
  await assert.rejects(
    runCli(['local-publish', '--manifest', 'unused.json', '--apply']),
    /--confirm|confirmation/i,
  );
  await assert.rejects(
    runCli(['local-publish', '--manifest', 'unused.json', '--bogus', 'value']),
    /unknown.*--bogus/i,
  );
});

test('local CLI preview returns digest and apply forwards the exact confirmation', async () => {
  const rootDir = await mkdtemp(path.join(tmpdir(), 'easypos-local-cli-'));
  const manifestPath = path.join(rootDir, 'manifest.json');
  const localManifest = { version: 2, mode: 'local-worktree', repositories: {} };
  await writeFile(manifestPath, JSON.stringify(localManifest));
  const calls = [];
  const localModule = {
    normalizeLocalManifest: (value) => value,
    preflightLocalBatch: async (options) => {
      calls.push(['preflight', options.target]);
      return { target: options.target, repositories: [] };
    },
    publishLocalBatch: async (options) => {
      calls.push(['publish', options.apply, options.confirmationDigest]);
      return options.apply
        ? { applied: true, results: [] }
        : { applied: false, confirmationDigest: 'a'.repeat(64), results: [] };
    },
  };
  const dependencies = {
    rootDir,
    config: { baseUrl: 'https://gitlab.test', token: 'test-token' },
    api: {},
    git: {},
    localModule,
  };

  const preflight = await runCli(['local-preflight', '--target', 'CURRENT'], dependencies);
  assert.equal(preflight.applied, false);
  const preview = await runCli(['local-publish', '--manifest', manifestPath], dependencies);
  assert.equal(preview.confirmationDigest, 'a'.repeat(64));
  const applied = await runCli([
    'local-publish',
    '--manifest', manifestPath,
    '--apply',
    '--confirm', 'a'.repeat(64),
  ], dependencies);
  assert.equal(applied.applied, true);
  assert.deepEqual(calls, [
    ['preflight', 'CURRENT'],
    ['preflight', 'CURRENT'],
    ['publish', false, undefined],
    ['preflight', 'CURRENT'],
    ['publish', true, 'a'.repeat(64)],
  ]);
});

test('pushed-branch CLI uses the injected Git inspector', async () => {
  const api = createApi({
    getBranch: async (projectPath, branch) => ({
      name: branch,
      commit: { id: branch === 'release/1.1' ? 'target' : repositoryStates[projectPath.split('/').at(-1)].localSha },
    }),
  });

  const result = await runCli([
    'preflight',
    '--source', 'CURRENT',
    '--target', 'release/1.1',
  ], {
    rootDir: '/unused-and-must-not-be-read',
    config: {},
    api,
    git: createGitInspector(),
  });

  assert.deepEqual(result.repositories.map(({ status }) => status), ['ready', 'ready', 'skipped']);
});

function issueRecordManifest() {
  const original = manifest();
  return {
    version: 2,
    mode: 'issue-only',
    source: 'feature/already-merged',
    target: 'release/1.1',
    repositories: {
      'base-app': {
        title: original.repositories['base-app'].title,
        changelog: original.repositories['base-app'].changelog,
        sourceSha: repositoryStates['base-app'].localSha,
        targetSha: repositoryStates['base-app'].localTargetSha,
      },
    },
  };
}

test('Issue-only manifest is distinct and rejects root repositories', () => {
  assert.deepEqual(normalizeIssueRecordManifest(issueRecordManifest()), issueRecordManifest());
  assert.throws(
    () => normalizeIssueRecordManifest({
      ...issueRecordManifest(),
      repositories: { '.': issueRecordManifest().repositories['base-app'] },
    }),
    /not allowed/i,
  );
});

test('pushed-branch and Issue-only manifests reject a multi-line title', () => {
  assert.throws(
    () => normalizeManifest({
      ...manifest(),
      repositories: {
        'base-app': { ...manifest().repositories['base-app'], title: 'line one\nline two' },
      },
    }),
    /base-app\.title.*single line/i,
  );
  assert.throws(
    () => normalizeIssueRecordManifest({
      ...issueRecordManifest(),
      repositories: {
        'base-app': { ...issueRecordManifest().repositories['base-app'], title: 'line one\nline two' },
      },
    }),
    /base-app\.title.*single line/i,
  );
});

test('Issue-only release record previews by default and never invokes Merge Request methods', async () => {
  const calls = [];
  const api = {
    findOpenIssueByMarker: async () => {
      calls.push('find-issue');
      return null;
    },
    findIssueByMarker: async () => {
      calls.push('find-issue');
      return null;
    },
    createIssue: async (_projectPath, payload) => {
      calls.push(['create-issue', payload]);
      return { iid: 81, web_url: 'https://gitlab.test/issues/81' };
    },
    updateIssue: async () => calls.push('update-issue'),
    listOpenMergeRequests: async () => calls.push('list-mr'),
    createMergeRequest: async () => calls.push('create-mr'),
    updateMergeRequest: async () => calls.push('update-mr'),
  };
  const comparison = {
    'base-app': {
      projectPath: 'oxygen/base-app',
      sourceSha: repositoryStates['base-app'].localSha,
      targetSha: repositoryStates['base-app'].localTargetSha,
      compare_timeout: false,
      diffs: [],
    },
  };

  const preview = await publishIssueRecords({
    manifest: issueRecordManifest(),
    comparison,
    api,
    apply: false,
  });
  assert.equal(preview.applied, false);
  assert.deepEqual(calls, []);

  const applied = await publishIssueRecords({
    manifest: issueRecordManifest(),
    comparison,
    api,
    apply: true,
  });
  assert.equal(applied.results[0].status, 'created');
  assert.deepEqual(calls.map((call) => Array.isArray(call) ? call[0] : call), ['find-issue', 'create-issue']);
  const description = calls[1][1].description;
  assert.match(description, /发布记录|变更记录/);
  assert.match(description, /easypos-gitlab-issue-record:v2/);
  assert.doesNotMatch(description, /Closes\s+#/i);
});

test('Issue-only release record rejects repositories with a current final diff', async () => {
  const api = createApi();
  await assert.rejects(
    publishIssueRecords({
      manifest: issueRecordManifest(),
      comparison: {
        'base-app': {
          projectPath: 'oxygen/base-app',
          sourceSha: repositoryStates['base-app'].localSha,
          targetSha: repositoryStates['base-app'].localTargetSha,
          diffs: [{ new_path: 'still-unmerged.txt' }],
        },
      },
      api,
      apply: true,
    }),
    /final diff.*Merge Request workflow/i,
  );
  assert.equal(api.writes.length, 0);
});

test('Issue-only manifest rejects resolved revision drift after changelog generation', async () => {
  const api = createApi();
  await assert.rejects(
    publishIssueRecords({
      manifest: issueRecordManifest(),
      comparison: {
        'base-app': {
          projectPath: 'oxygen/base-app',
          sourceSha: 'different-source-sha',
          targetSha: repositoryStates['base-app'].localTargetSha,
          diffs: [],
          sourceMerged: true,
        },
      },
      api,
      apply: false,
    }),
    /base-app.*source or target SHA.*changelog manifest/i,
  );
  await assert.rejects(
    publishIssueRecords({
      manifest: issueRecordManifest(),
      comparison: {
        'base-app': {
          projectPath: 'oxygen/base-app',
          sourceSha: repositoryStates['base-app'].localSha,
          targetSha: 'different-target-sha',
          diffs: [],
          sourceMerged: true,
        },
      },
      api,
      apply: false,
    }),
    /base-app.*source or target SHA.*changelog manifest/i,
  );
  assert.equal(api.writes.length, 0);
});

test('Issue-only preflight accepts an already-merged source by commit ancestry despite later target changes', async () => {
  const api = createApi({
    getCommit: async (_projectPath, sha) => ({ id: sha }),
    compare: async () => ({ diffs: [{ new_path: 'later-target-change.txt' }] }),
  });
  const inspector = {
    inspectIssueRecord: async (repository) => ({
      repository,
      projectPath: `oxygen/${repository}`,
      remoteHost: 'gitlab.sjfood.us',
      sourceSha: repositoryStates[repository].localSha,
      targetSha: repositoryStates[repository].localTargetSha,
      sourceMerged: true,
    }),
  };

  const comparison = await preflightIssueRecordBatch({
    manifest: issueRecordManifest(),
    git: inspector,
    api,
  });
  const result = await publishIssueRecords({
    manifest: issueRecordManifest(),
    comparison,
    api,
    apply: false,
  });
  assert.equal(comparison['base-app'].sourceMerged, true);
  assert.equal(result.results[0].status, 'preview');
  assert.equal(api.writes.length, 0);
});

test('Issue-only CLI routes an explicit release record without creating an MR', async () => {
  const rootDir = await mkdtemp(path.join(tmpdir(), 'easypos-issue-record-cli-'));
  const manifestPath = path.join(rootDir, 'issue-record.json');
  await writeFile(manifestPath, JSON.stringify(issueRecordManifest()));
  const api = createApi({
    getBranch: async (projectPath, branch) => ({
      name: branch,
      commit: {
        id: branch === 'release/1.1'
          ? 'target'
          : repositoryStates[projectPath.split('/').at(-1)].localSha,
      },
    }),
    compare: async () => ({ diffs: [] }),
  });
  const dependencies = {
    rootDir,
    config: {},
    api,
    git: createGitInspector(),
  };

  const preview = await runCli(['issue-record', '--manifest', manifestPath], dependencies);
  assert.equal(preview.applied, false);
  assert.equal(api.writes.length, 0);
  const applied = await runCli(['issue-record', '--manifest', manifestPath, '--apply'], dependencies);
  assert.equal(applied.results[0].status, 'created');
  assert.deepEqual(api.writes.map(({ kind }) => kind), ['issue']);
});

test('Issue-only release record updates a closed owned Issue instead of duplicating it', async () => {
  const writes = [];
  const api = {
    findIssueByMarker: async () => ({ iid: 91, state: 'closed', web_url: 'https://gitlab.test/issues/91' }),
    updateIssue: async (_projectPath, iid, payload) => {
      writes.push(['update', iid]);
      return { iid, web_url: 'https://gitlab.test/issues/91', ...payload };
    },
    createIssue: async () => writes.push(['create']),
  };
  const result = await publishIssueRecords({
    manifest: issueRecordManifest(),
    comparison: {
      'base-app': {
        projectPath: 'oxygen/base-app',
        sourceSha: repositoryStates['base-app'].localSha,
        targetSha: repositoryStates['base-app'].localTargetSha,
        diffs: [],
        sourceMerged: true,
      },
    },
    api,
    apply: true,
  });
  assert.equal(result.results[0].status, 'updated');
  assert.deepEqual(writes, [['update', 91]]);
});

test('Issue-only failure reports issues completed in earlier repositories', async () => {
  const original = manifest();
  const releaseRecord = issueRecordManifest();
  releaseRecord.repositories['base-engine'] = {
    title: original.repositories['base-engine'].title,
    changelog: original.repositories['base-engine'].changelog,
    sourceSha: repositoryStates['base-engine'].localSha,
    targetSha: repositoryStates['base-engine'].localTargetSha,
  };
  const api = createApi({
    findIssueByMarker: async () => null,
    createIssue: async (projectPath, payload) => {
      if (projectPath.endsWith('engine')) throw new Error('engine Issue failed');
      api.writes.push({ kind: 'issue', projectPath, payload });
      return { iid: 101, web_url: 'https://gitlab.test/oxygen/base-app/issues/101' };
    },
  });

  await assert.rejects(
    publishIssueRecords({
      manifest: releaseRecord,
      comparison: {
        'base-app': {
          projectPath: 'oxygen/base-app',
          sourceSha: repositoryStates['base-app'].localSha,
          targetSha: repositoryStates['base-app'].localTargetSha,
          diffs: [],
          sourceMerged: true,
        },
        'base-engine': {
          projectPath: 'oxygen/base-engine',
          sourceSha: repositoryStates['base-engine'].localSha,
          targetSha: repositoryStates['base-engine'].localTargetSha,
          diffs: [],
          sourceMerged: true,
        },
      },
      api,
      apply: true,
    }),
    (error) => {
      assert.match(error.message, /base-engine.*engine Issue failed/i);
      assert.match(error.message, /base-app.*issues\/101/i);
      return true;
    },
  );
});
