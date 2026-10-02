import assert from 'node:assert/strict';
import { chmod, lstat, mkdir, mkdtemp, readFile, symlink, unlink, writeFile } from 'node:fs/promises';
import { tmpdir } from 'node:os';
import path from 'node:path';
import test from 'node:test';
import { execFile as execFileCallback } from 'node:child_process';
import { promisify } from 'node:util';

import {
  LocalGitWorkspace,
  buildLocalOwnershipMarker,
  deriveFeatureBranch,
  normalizeLocalManifest,
  preflightLocalBatch,
  previewLocalBatch,
  publishLocalBatch,
  publishLocalGitState,
  validateLocalApplyState,
  validateBranchSlug,
} from './local-worktree-publication.mjs';

const execFile = promisify(execFileCallback);
const BASE_SHA = '1'.repeat(40);
const TREE_SHA = '2'.repeat(40);

function changelog(item = '支持本地工作区发布') {
  return [
    '# 本地发布',
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
    '- [ ] [P0] 验证本地发布结果',
  ].join('\n');
}

function localManifest(overrides = {}) {
  return {
    version: 2,
    mode: 'local-worktree',
    repositories: {
      'base-app': {
        title: '新增订单折扣',
        changelog: changelog(),
        baseBranch: '1.1.0-rc',
        baseSha: BASE_SHA,
        snapshotTreeSha: TREE_SHA,
        paths: ['lib/order.dart'],
        branchSlug: 'order-discount',
        ...overrides,
      },
    },
  };
}

test('local manifest normalizes an exact version 2 worktree snapshot', () => {
  assert.deepEqual(normalizeLocalManifest(localManifest()), localManifest());
  assert.equal(validateBranchSlug('order-discount'), 'order-discount');
});

test('local manifest rejects unsafe repositories, paths, SHAs, and branch slugs', () => {
  assert.throws(
    () => normalizeLocalManifest({ ...localManifest(), repositories: { '.': localManifest().repositories['base-app'] } }),
    /not allowed/i,
  );
  for (const unsafePath of ['/absolute', '../escape', 'nested/../../escape', '']) {
    assert.throws(
      () => normalizeLocalManifest(localManifest({ paths: [unsafePath] })),
      /path/i,
    );
  }
  assert.throws(() => normalizeLocalManifest(localManifest({ paths: ['a', 'a'] })), /duplicate/i);
  assert.throws(() => normalizeLocalManifest(localManifest({ baseSha: 'short' })), /baseSha/i);
  assert.throws(() => normalizeLocalManifest(localManifest({ snapshotTreeSha: 'short' })), /snapshotTreeSha/i);
  for (const slug of ['Uppercase', '-leading', 'contains/slash', 'a'.repeat(64)]) {
    assert.throws(() => validateBranchSlug(slug), /branchSlug/i);
  }
  assert.throws(() => normalizeLocalManifest(localManifest({ title: 'line one\nline two' })), /title.*single line/i);
});

test('local manifest requires repositories and all five changelog sections', () => {
  assert.throws(
    () => normalizeLocalManifest({ version: 2, mode: 'local-worktree', repositories: {} }),
    /repositories/i,
  );
  assert.throws(
    () => normalizeLocalManifest(localManifest({ changelog: '## 新增功能\n- A' })),
    /功能优化/i,
  );
});

test('local Git commands have a timeout and disable interactive credential prompts', async () => {
  const calls = [];
  const workspace = new LocalGitWorkspace({
    rootDir: '/unused',
    timeoutMs: 1_234,
    execGit: async (command, args, options) => {
      calls.push({ command, args, options });
      return { stdout: 'ok\n' };
    },
  });

  assert.equal(await workspace.git('base-app', ['status']), 'ok');
  assert.equal(calls[0].options.timeout, 1_234);
  assert.equal(calls[0].options.env.GIT_TERMINAL_PROMPT, '0');
});

async function git(cwd, args, options = {}) {
  const { stdout } = await execFile('git', args, { cwd, encoding: 'utf8', ...options });
  return stdout.trim();
}

async function createSnapshotFixture() {
  const rootDir = await mkdtemp(path.join(tmpdir(), 'easypos-local-snapshot-'));
  const repositoryDir = path.join(rootDir, 'base-app');
  await mkdir(repositoryDir);
  await git(repositoryDir, ['init', '-b', '1.1.0-rc']);
  await git(repositoryDir, ['config', 'user.name', 'Snapshot Test']);
  await git(repositoryDir, ['config', 'user.email', 'snapshot@example.test']);
  await writeFile(path.join(repositoryDir, 'staged.txt'), 'base staged\n');
  await writeFile(path.join(repositoryDir, 'unstaged.txt'), 'base unstaged\n');
  await writeFile(path.join(repositoryDir, 'delete.txt'), 'delete me\n');
  await writeFile(path.join(repositoryDir, 'rename-old.txt'), 'rename me\n');
  await writeFile(path.join(repositoryDir, 'executable.sh'), '#!/bin/sh\nexit 0\n');
  await git(repositoryDir, ['add', '-A']);
  await git(repositoryDir, ['commit', '-m', 'base']);
  await git(repositoryDir, ['remote', 'add', 'origin', 'git@gitlab.sjfood.us:oxygen/base-app.git']);

  await writeFile(path.join(repositoryDir, 'staged.txt'), 'staged change\n');
  await git(repositoryDir, ['add', 'staged.txt']);
  await writeFile(path.join(repositoryDir, 'unstaged.txt'), 'unstaged change\n');
  await writeFile(path.join(repositoryDir, 'untracked.txt'), 'untracked\n');
  await unlink(path.join(repositoryDir, 'delete.txt'));
  await git(repositoryDir, ['mv', 'rename-old.txt', 'rename-new.txt']);
  await chmod(path.join(repositoryDir, 'executable.sh'), 0o755);
  await symlink('staged.txt', path.join(repositoryDir, 'linked.txt'));
  return { rootDir, repositoryDir };
}

test('snapshot captures all worktree forms without changing real Git state', async () => {
  const { rootDir, repositoryDir } = await createSnapshotFixture();
  const indexPath = await git(repositoryDir, ['rev-parse', '--git-path', 'index']);
  const resolvedIndexPath = path.resolve(repositoryDir, indexPath);
  const before = {
    branch: await git(repositoryDir, ['branch', '--show-current']),
    head: await git(repositoryDir, ['rev-parse', 'HEAD']),
    status: await git(repositoryDir, ['status', '--porcelain=v1', '--untracked-files=all']),
    index: await readFile(resolvedIndexPath),
  };

  const snapshot = await new LocalGitWorkspace({ rootDir }).snapshot('base-app');

  assert.equal(snapshot.baseBranch, '1.1.0-rc');
  assert.equal(snapshot.baseSha, before.head);
  assert.match(snapshot.snapshotTreeSha, /^[a-f0-9]{40}$/);
  assert.deepEqual(snapshot.paths, [
    'delete.txt',
    'executable.sh',
    'linked.txt',
    'rename-new.txt',
    'rename-old.txt',
    'staged.txt',
    'unstaged.txt',
    'untracked.txt',
  ]);
  assert.equal(snapshot.projectPath, 'oxygen/base-app');
  assert.equal(snapshot.remoteHost, 'gitlab.sjfood.us');
  assert.equal(await git(repositoryDir, ['branch', '--show-current']), before.branch);
  assert.equal(await git(repositoryDir, ['rev-parse', 'HEAD']), before.head);
  assert.equal(await git(repositoryDir, ['status', '--porcelain=v1', '--untracked-files=all']), before.status);
  assert.deepEqual(await readFile(resolvedIndexPath), before.index);
  assert.equal((await lstat(path.join(repositoryDir, 'linked.txt'))).isSymbolicLink(), true);
});

const localStates = Object.freeze({
  'base-app': {
    repository: 'base-app',
    baseBranch: '1.1.0-rc',
    baseSha: BASE_SHA,
    snapshotTreeSha: TREE_SHA,
    paths: ['lib/order.dart'],
    dirty: true,
    remoteHost: 'gitlab.sjfood.us',
    projectPath: 'oxygen/base-app',
    conflicts: [],
    operationInProgress: null,
    gitUserName: 'Developer',
    gitUserEmail: 'developer@example.test',
  },
  'base-engine': {
    repository: 'base-engine',
    baseBranch: '1.1.0-rc',
    baseSha: '3'.repeat(40),
    snapshotTreeSha: '4'.repeat(40),
    paths: [],
    dirty: false,
    remoteHost: 'gitlab.sjfood.us',
    projectPath: 'oxygen/base-engine',
    conflicts: [],
    operationInProgress: null,
    gitUserName: 'Developer',
    gitUserEmail: 'developer@example.test',
  },
  'base-web': {
    repository: 'base-web',
    baseBranch: '1.1.0-rc',
    baseSha: '5'.repeat(40),
    snapshotTreeSha: '6'.repeat(40),
    paths: [],
    dirty: false,
    remoteHost: 'gitlab.sjfood.us',
    projectPath: 'oxygen/base-web',
    conflicts: [],
    operationInProgress: null,
    gitUserName: 'Developer',
    gitUserEmail: 'developer@example.test',
  },
});

function localInspector(overrides = {}) {
  return {
    snapshot: async (repository) => ({ ...localStates[repository], ...overrides[repository] }),
  };
}

function localApi(overrides = {}) {
  const writes = [];
  return {
    writes,
    getProject: async (projectPath) => ({
      id: projectPath,
      permissions: { project_access: { access_level: 30 } },
    }),
    getBranch: async (projectPath, branch) => ({
      name: branch,
      commit: { id: localStates[projectPath.split('/').at(-1)].baseSha },
    }),
    createIssue: async (...args) => writes.push(['issue', ...args]),
    createMergeRequest: async (...args) => writes.push(['mr', ...args]),
    ...overrides,
  };
}

test('local preflight marks dirty repositories ready and clean repositories skipped', async () => {
  const api = localApi();
  const result = await preflightLocalBatch({
    target: 'CURRENT',
    git: localInspector(),
    api,
    baseUrl: 'https://gitlab.sjfood.us',
  });

  assert.deepEqual(result.repositories.map(({ repository, status }) => [repository, status]), [
    ['base-app', 'ready'],
    ['base-engine', 'skipped'],
    ['base-web', 'skipped'],
  ]);
  assert.equal(result.repositories[0].baseBranch, '1.1.0-rc');
  assert.equal(api.writes.length, 0);
});

test('local preflight aggregates unsafe Git states before GitLab access', async () => {
  let reads = 0;
  const api = localApi({
    getProject: async () => {
      reads += 1;
      return {};
    },
  });

  await assert.rejects(
    preflightLocalBatch({
      target: 'release/other',
      git: localInspector({
        'base-app': {
          conflicts: ['lib/conflict.dart'],
          operationInProgress: 'merge',
          gitUserEmail: '',
        },
      }),
      api,
    }),
    (error) => {
      assert.match(error.message, /target.*current branch/i);
      assert.match(error.message, /conflict/i);
      assert.match(error.message, /merge/i);
      assert.match(error.message, /user\.email/i);
      return true;
    },
  );
  assert.equal(reads, 0);
  assert.equal(api.writes.length, 0);
});

test('local preflight rejects a base commit that is not pushed', async () => {
  const api = localApi({
    getBranch: async (_projectPath, branch) => ({ name: branch, commit: { id: '9'.repeat(40) } }),
  });
  await assert.rejects(
    preflightLocalBatch({ target: 'CURRENT', git: localInspector(), api }),
    /base-app.*base branch.*not pushed/i,
  );
  assert.equal(api.writes.length, 0);
});

test('local preflight recognizes an owned recovery branch from manifest state', async () => {
  const api = localApi();
  const result = await preflightLocalBatch({
    target: 'CURRENT',
    manifest: localManifest(),
    git: localInspector({
      'base-app': {
        baseBranch: 'feature/123-order-discount',
        baseSha: '3'.repeat(40),
        snapshotTreeSha: TREE_SHA,
        paths: [],
        dirty: false,
      },
    }),
    api,
  });

  assert.equal(result.repositories[0].status, 'ready');
  assert.equal(result.repositories[0].baseBranch, '1.1.0-rc');
  assert.equal(result.repositories[0].baseSha, BASE_SHA);
  assert.equal(result.repositories[0].currentBranch, 'feature/123-order-discount');
  assert.equal(api.writes.length, 0);
});

function readyPreflight(overrides = {}) {
  return {
    target: 'CURRENT',
    repositories: [
      { ...localStates['base-app'], status: 'ready', ...overrides },
      { ...localStates['base-engine'], status: 'skipped', reason: 'no local changes' },
      { ...localStates['base-web'], status: 'skipped', reason: 'no local changes' },
    ],
  };
}

test('preview binds the exact manifest to a confirmation digest without writes', () => {
  const result = previewLocalBatch({
    manifest: localManifest(),
    preflight: readyPreflight(),
  });

  assert.match(result.confirmationDigest, /^[a-f0-9]{64}$/);
  assert.deepEqual(result.results[0], {
    repository: 'base-app',
    status: 'preview',
    paths: ['lib/order.dart'],
    title: '新增订单折扣',
    targetBranch: '1.1.0-rc',
    sourceBranchPattern: 'feature/<issue_iid>-order-discount',
    commitMessage: 'Refs #<issue_iid>: 新增订单折扣',
  });
});

test('apply validation rejects manifest, repository state, and confirmation drift', () => {
  const manifest = localManifest();
  const preflight = readyPreflight();
  const preview = previewLocalBatch({ manifest, preflight });
  assert.doesNotThrow(() => validateLocalApplyState({
    manifest,
    preflight,
    confirmationDigest: preview.confirmationDigest,
  }));

  for (const drift of [
    readyPreflight({ baseBranch: 'other' }),
    readyPreflight({ baseSha: '7'.repeat(40) }),
    readyPreflight({ snapshotTreeSha: '8'.repeat(40) }),
    readyPreflight({ paths: ['different.txt'] }),
  ]) {
    assert.throws(
      () => validateLocalApplyState({ manifest, preflight: drift, confirmationDigest: preview.confirmationDigest }),
      /drift/i,
    );
  }
  assert.throws(
    () => validateLocalApplyState({ manifest, preflight, confirmationDigest: '0'.repeat(64) }),
    /confirmation/i,
  );
  assert.throws(
    () => validateLocalApplyState({
      manifest,
      preflight,
      confirmationDigest: `${preview.confirmationDigest}zz`,
    }),
    /confirmation/i,
  );
});

test('feature branch derivation and command policy use only guarded Git operations', async () => {
  assert.equal(deriveFeatureBranch(123, 'order-discount'), 'feature/123-order-discount');
  const commands = [];
  let remoteReads = 0;
  const responses = new Map([
    ['config --get user.name', 'Developer'],
    ['config --get user.email', 'developer@example.test'],
    ['config --get remote.origin.url', 'git@gitlab.sjfood.us:oxygen/base-app.git'],
    ['symbolic-ref --quiet --short HEAD', '1.1.0-rc'],
    ['rev-parse --verify HEAD^{commit}', BASE_SHA],
    ['show-ref --verify --quiet refs/heads/feature/123-order-discount', ''],
    ['ls-remote --heads origin refs/heads/feature/123-order-discount', ''],
    ['rev-parse HEAD^', BASE_SHA],
    ['rev-parse HEAD^{tree}', TREE_SHA],
    ['write-tree', TREE_SHA],
    ['status --porcelain=v1 --untracked-files=all', ''],
    ['log -1 --format=%B', `Refs #123: 新增订单折扣\n\nEasyPos-Change-Request: ${TREE_SHA}`],
    ['rev-parse HEAD', '3'.repeat(40)],
  ]);
  const git = {
    git: async (_repository, args) => {
      commands.push(args);
      if (args.join(' ') === 'ls-remote --heads origin refs/heads/feature/123-order-discount') {
        remoteReads += 1;
        return remoteReads === 1 ? '' : `${'3'.repeat(40)}\trefs/heads/feature/123-order-discount`;
      }
      return responses.get(args.join(' ')) || '';
    },
    gitOptional: async (_repository, args) => {
      commands.push(args);
      if (args.join(' ') === 'ls-remote --heads origin refs/heads/feature/123-order-discount') {
        remoteReads += 1;
        return remoteReads === 1 ? '' : `${'3'.repeat(40)}\trefs/heads/feature/123-order-discount`;
      }
      return responses.get(args.join(' ')) || '';
    },
  };

  await publishLocalGitState({
    repository: 'base-app',
    state: localStates['base-app'],
    entry: localManifest().repositories['base-app'],
    issue: { iid: 123, web_url: 'https://gitlab.test/issues/123' },
    git,
  });

  const rendered = commands.map((args) => args.join(' '));
  assert.ok(rendered.includes('switch -c feature/123-order-discount 1111111111111111111111111111111111111111'));
  assert.ok(rendered.includes('add -A -- :(top,literal)lib/order.dart'));
  assert.ok(rendered.includes(`commit -m Refs #123: 新增订单折扣 -m EasyPos-Change-Request: ${TREE_SHA}`));
  assert.ok(rendered.includes('push --set-upstream origin feature/123-order-discount'));
  assert.equal(rendered.some((command) => /--force|--amend|\brebase\b|\breset\b|\brestore\b|\bstash\b|\bclean\b|\bmerge\b|branch -[dD]|\btag\b/.test(command)), false);
});

test('Git publication rejects origin project drift before branch creation', async () => {
  const calls = [];
  const git = {
    git: async (_repository, args) => {
      calls.push(args.join(' '));
      if (args.join(' ') === 'config --get remote.origin.url') {
        return 'git@gitlab.sjfood.us:unrelated/base-app.git';
      }
      return '';
    },
    gitOptional: async (_repository, args) => {
      calls.push(args.join(' '));
      if (args.join(' ') === 'config --get user.name') return 'Developer';
      if (args.join(' ') === 'config --get user.email') return 'developer@example.test';
      if (args.join(' ') === 'symbolic-ref --quiet --short HEAD') return '1.1.0-rc';
      if (args.join(' ') === 'rev-parse --verify HEAD^{commit}') return BASE_SHA;
      return '';
    },
  };

  await assert.rejects(
    publishLocalGitState({
      repository: 'base-app',
      state: localStates['base-app'],
      entry: localManifest().repositories['base-app'],
      issue: { iid: 123 },
      git,
    }),
    /origin project path.*oxygen\/base-app|oxygen\/base-app.*received/i,
  );
  assert.equal(calls.some((command) => command.startsWith('switch ')), false);
});

test('Git publication treats ls-remote failure as blocking and does not create a branch', async () => {
  const calls = [];
  const responses = new Map([
    ['config --get remote.origin.url', 'git@gitlab.sjfood.us:oxygen/base-app.git'],
    ['config --get user.name', 'Developer'],
    ['config --get user.email', 'developer@example.test'],
    ['symbolic-ref --quiet --short HEAD', '1.1.0-rc'],
    ['rev-parse --verify HEAD^{commit}', BASE_SHA],
  ]);
  const workspace = {
    git: async (_repository, args) => {
      const command = args.join(' ');
      calls.push(command);
      if (command.startsWith('ls-remote ')) throw new Error('network unavailable');
      return responses.get(command) || '';
    },
    gitOptional: async (_repository, args) => {
      const command = args.join(' ');
      calls.push(command);
      try {
        if (command.startsWith('ls-remote ')) throw new Error('network unavailable');
        return responses.get(command) || '';
      } catch {
        return '';
      }
    },
  };

  await assert.rejects(
    publishLocalGitState({
      repository: 'base-app',
      state: localStates['base-app'],
      entry: localManifest().repositories['base-app'],
      issue: { iid: 123 },
      git: workspace,
    }),
    /network unavailable/i,
  );
  assert.equal(calls.some((command) => command.startsWith('switch ')), false);
});

async function createPublicationFixture() {
  const rootDir = await mkdtemp(path.join(tmpdir(), 'easypos-local-publication-'));
  const repositoryDir = path.join(rootDir, 'base-app');
  const remoteDir = path.join(rootDir, 'origin.git');
  await mkdir(repositoryDir);
  await git(rootDir, ['init', '--bare', remoteDir]);
  await git(repositoryDir, ['init', '-b', '1.1.0-rc']);
  await git(repositoryDir, ['config', 'user.name', 'Publication Test']);
  await git(repositoryDir, ['config', 'user.email', 'publication@example.test']);
  await writeFile(path.join(repositoryDir, 'tracked.txt'), 'base\n');
  await writeFile(path.join(repositoryDir, 'delete.txt'), 'delete\n');
  await git(repositoryDir, ['add', '-A']);
  await git(repositoryDir, ['commit', '-m', 'base']);
  await git(repositoryDir, ['remote', 'add', 'origin', remoteDir]);
  await git(repositoryDir, ['push', '--set-upstream', 'origin', '1.1.0-rc']);
  await git(repositoryDir, ['remote', 'set-url', 'origin', 'git@gitlab.sjfood.us:oxygen/base-app.git']);
  await git(repositoryDir, [
    'config',
    `url.${remoteDir}.insteadOf`,
    'git@gitlab.sjfood.us:oxygen/base-app.git',
  ]);

  await writeFile(path.join(repositoryDir, 'tracked.txt'), 'changed\n');
  await git(repositoryDir, ['add', 'tracked.txt']);
  await unlink(path.join(repositoryDir, 'delete.txt'));
  await writeFile(path.join(repositoryDir, 'untracked.txt'), 'new\n');
  const workspace = new LocalGitWorkspace({ rootDir });
  const state = await workspace.snapshot('base-app');
  return { rootDir, repositoryDir, remoteDir, workspace, state };
}

test('temporary repository lifecycle creates one verified commit and normal push', async () => {
  const { repositoryDir, remoteDir, workspace, state } = await createPublicationFixture();
  const entry = {
    ...localManifest().repositories['base-app'],
    baseSha: state.baseSha,
    snapshotTreeSha: state.snapshotTreeSha,
    paths: state.paths,
  };

  const result = await publishLocalGitState({
    repository: 'base-app',
    state,
    entry,
    issue: { iid: 123, web_url: 'https://gitlab.test/issues/123' },
    git: workspace,
  });

  assert.equal(await git(repositoryDir, ['branch', '--show-current']), 'feature/123-order-discount');
  assert.equal(await git(repositoryDir, ['rev-parse', 'HEAD^']), state.baseSha);
  assert.equal(await git(repositoryDir, ['rev-parse', 'HEAD^{tree}']), state.snapshotTreeSha);
  assert.match(await git(repositoryDir, ['log', '-1', '--format=%B']), new RegExp(`EasyPos-Change-Request: ${state.snapshotTreeSha}`));
  assert.equal(await git(repositoryDir, ['status', '--porcelain=v1', '--untracked-files=all']), '');
  assert.equal(await git(remoteDir, ['rev-parse', 'refs/heads/feature/123-order-discount']), result.commitSha);
  assert.equal(await git(remoteDir, ['rev-parse', 'refs/heads/1.1.0-rc']), state.baseSha);
  assert.equal(result.pushed, true);
});

test('foreign branch collisions and missing Git identity stop before publication', async () => {
  const { repositoryDir, workspace, state } = await createPublicationFixture();
  const entry = {
    ...localManifest().repositories['base-app'],
    baseSha: state.baseSha,
    snapshotTreeSha: state.snapshotTreeSha,
    paths: state.paths,
  };
  await git(repositoryDir, ['branch', 'feature/123-order-discount', state.baseSha]);
  await assert.rejects(
    publishLocalGitState({ repository: 'base-app', state, entry, issue: { iid: 123 }, git: workspace }),
    /local branch.*collision/i,
  );
  assert.equal(await git(repositoryDir, ['branch', '--show-current']), '1.1.0-rc');

  await git(repositoryDir, ['branch', '-D', 'feature/123-order-discount']);
  await git(repositoryDir, ['config', 'user.email', '']);
  await assert.rejects(
    publishLocalGitState({ repository: 'base-app', state, entry, issue: { iid: 123 }, git: workspace }),
    /user\.email/i,
  );
  assert.equal(await git(repositoryDir, ['branch', '--show-current']), '1.1.0-rc');
});

test('a failing pre-commit hook leaves the generated branch unpushed without cleanup', async () => {
  const { repositoryDir, remoteDir, workspace, state } = await createPublicationFixture();
  const entry = {
    ...localManifest().repositories['base-app'],
    baseSha: state.baseSha,
    snapshotTreeSha: state.snapshotTreeSha,
    paths: state.paths,
  };
  const hookPath = path.join(repositoryDir, '.git', 'hooks', 'pre-commit');
  await writeFile(hookPath, '#!/bin/sh\nexit 1\n');
  await chmod(hookPath, 0o755);

  await assert.rejects(
    publishLocalGitState({ repository: 'base-app', state, entry, issue: { iid: 123 }, git: workspace }),
    /commit/i,
  );
  assert.equal(await git(repositoryDir, ['branch', '--show-current']), 'feature/123-order-discount');
  await assert.rejects(git(remoteDir, ['rev-parse', 'refs/heads/feature/123-order-discount']));
});

test('foreign remote branch collision stops before creating a local feature branch', async () => {
  const { repositoryDir, remoteDir, workspace, state } = await createPublicationFixture();
  const entry = {
    ...localManifest().repositories['base-app'],
    baseSha: state.baseSha,
    snapshotTreeSha: state.snapshotTreeSha,
    paths: state.paths,
  };
  await git(remoteDir, ['update-ref', 'refs/heads/feature/123-order-discount', state.baseSha]);

  await assert.rejects(
    publishLocalGitState({ repository: 'base-app', state, entry, issue: { iid: 123 }, git: workspace }),
    /remote branch.*collision/i,
  );
  assert.equal(await git(repositoryDir, ['branch', '--show-current']), '1.1.0-rc');
  await assert.rejects(git(repositoryDir, ['rev-parse', '--verify', 'refs/heads/feature/123-order-discount']));
});

test('post-commit hook content drift is detected before push without automatic cleanup', async () => {
  const { repositoryDir, remoteDir, workspace, state } = await createPublicationFixture();
  const entry = {
    ...localManifest().repositories['base-app'],
    baseSha: state.baseSha,
    snapshotTreeSha: state.snapshotTreeSha,
    paths: state.paths,
  };
  const hookPath = path.join(repositoryDir, '.git', 'hooks', 'post-commit');
  await writeFile(hookPath, '#!/bin/sh\nprintf drift > hook-drift.txt\n');
  await chmod(hookPath, 0o755);

  await assert.rejects(
    publishLocalGitState({ repository: 'base-app', state, entry, issue: { iid: 123 }, git: workspace }),
    /worktree changed/i,
  );
  assert.equal(await git(repositoryDir, ['branch', '--show-current']), 'feature/123-order-discount');
  assert.equal(await readFile(path.join(repositoryDir, 'hook-drift.txt'), 'utf8'), 'drift');
  await assert.rejects(git(remoteDir, ['rev-parse', 'refs/heads/feature/123-order-discount']));
});

function statefulPublicationApi(events, { failIssue = false, failMr = false, foreignMr = false } = {}) {
  let issue = null;
  let mergeRequest = foreignMr
    ? { iid: 77, web_url: 'https://gitlab.test/mr/77', description: 'foreign' }
    : null;
  return {
    get issue() { return issue; },
    get mergeRequest() { return mergeRequest; },
    findOpenIssueByMarker: async (_projectPath, marker) => {
      events.push('find-issue');
      return issue && issue.description.includes(marker) ? issue : null;
    },
    findIssueByMarker: async (_projectPath, marker) => {
      events.push('find-issue');
      return issue && issue.description.includes(marker) ? issue : null;
    },
    createIssue: async (_projectPath, payload) => {
      events.push('create-issue');
      if (failIssue) throw new Error('issue failed');
      issue = { iid: 123, web_url: 'https://gitlab.test/issues/123', ...payload };
      return issue;
    },
    updateIssue: async (_projectPath, _iid, payload) => {
      events.push('update-issue');
      issue = { ...issue, ...payload };
      return issue;
    },
    listOpenMergeRequests: async () => {
      events.push('find-mr');
      return mergeRequest ? [mergeRequest] : [];
    },
    createMergeRequest: async (_projectPath, payload) => {
      events.push('create-mr');
      if (failMr) throw new Error('mr failed');
      mergeRequest = { iid: 456, web_url: 'https://gitlab.test/mr/456', ...payload };
      return mergeRequest;
    },
    updateMergeRequest: async (_projectPath, _iid, payload) => {
      events.push('update-mr');
      mergeRequest = { ...mergeRequest, ...payload };
      return mergeRequest;
    },
    allowIssue() { failIssue = false; },
    allowMr() { failMr = false; },
  };
}

test('local ownership marker is stable and Issue-first publication links identical changelog content', async () => {
  const marker = buildLocalOwnershipMarker({
    projectPath: 'oxygen/base-app',
    baseBranch: 'release/1.1',
    baseSha: BASE_SHA,
    snapshotTreeSha: TREE_SHA,
  });
  assert.equal(marker, buildLocalOwnershipMarker({
    projectPath: 'oxygen/base-app',
    baseBranch: 'release/1.1',
    baseSha: BASE_SHA,
    snapshotTreeSha: TREE_SHA,
  }));
  assert.match(marker, /^<!-- easypos-gitlab-local-change-request:v2:[A-Za-z0-9_-]+ -->$/);
  assert.doesNotMatch(marker, /release\/1\.1|oxygen\/base-app/);

  const events = [];
  const api = statefulPublicationApi(events);
  const preflight = readyPreflight();
  const manifest = localManifest();
  const { confirmationDigest } = previewLocalBatch({ manifest, preflight });
  const gitPublisher = async ({ issue }) => {
    events.push(`git-${issue.iid}`);
    return { branch: 'feature/123-order-discount', commitSha: '3'.repeat(40), pushed: true };
  };

  const preview = await publishLocalBatch({
    manifest,
    preflight,
    confirmationDigest,
    git: {},
    api,
    apply: false,
    publishGitState: gitPublisher,
  });
  assert.equal(preview.applied, false);
  assert.deepEqual(events, []);

  const result = await publishLocalBatch({
    manifest,
    preflight,
    confirmationDigest,
    git: {},
    api,
    apply: true,
    publishGitState: gitPublisher,
  });
  assert.deepEqual(events, ['find-issue', 'create-issue', 'git-123', 'find-mr', 'create-mr']);
  assert.equal(result.results[0].status, 'complete');
  assert.equal(result.results[0].action, 'created');
  assert.ok(api.issue.description.includes(manifest.repositories['base-app'].changelog));
  assert.ok(api.mergeRequest.description.startsWith('Closes #123\n\n'));
  assert.ok(api.mergeRequest.description.includes(manifest.repositories['base-app'].changelog));
});

test('local partial failure recovery reuses owned Issue and completed Git state', async () => {
  const events = [];
  const api = statefulPublicationApi(events, { failMr: true });
  const preflight = readyPreflight();
  const manifest = localManifest();
  const { confirmationDigest } = previewLocalBatch({ manifest, preflight });
  let gitCalls = 0;
  const gitPublisher = async () => {
    gitCalls += 1;
    events.push(gitCalls === 1 ? 'git-created' : 'git-reused');
    return {
      branch: 'feature/123-order-discount',
      commitSha: '3'.repeat(40),
      pushed: gitCalls === 1,
      recovered: gitCalls > 1,
    };
  };
  await assert.rejects(
    publishLocalBatch({ manifest, preflight, confirmationDigest, git: {}, api, apply: true, publishGitState: gitPublisher }),
    /Issue #123.*MR publication failed/i,
  );
  api.allowMr();
  const result = await publishLocalBatch({
    manifest,
    preflight,
    confirmationDigest,
    git: {},
    api,
    apply: true,
    publishGitState: gitPublisher,
  });
  assert.equal(events.filter((event) => event === 'create-issue').length, 1);
  assert.ok(events.includes('update-issue'));
  assert.ok(events.includes('git-reused'));
  assert.equal(result.results[0].mergeRequest.iid, 456);
});

test('foreign local merge request state is rejected without overwrite', async () => {
  const events = [];
  const api = statefulPublicationApi(events, { foreignMr: true });
  const preflight = readyPreflight();
  const manifest = localManifest();
  const { confirmationDigest } = previewLocalBatch({ manifest, preflight });
  await assert.rejects(
    publishLocalBatch({
      manifest,
      preflight,
      confirmationDigest,
      git: {},
      api,
      apply: true,
      publishGitState: async () => ({
        branch: 'feature/123-order-discount',
        commitSha: '3'.repeat(40),
        pushed: true,
      }),
    }),
    /not owned/i,
  );
  assert.equal(events.includes('update-mr'), false);
});

test('closed owned Issue blocks local recovery without creating a duplicate', async () => {
  const events = [];
  const api = statefulPublicationApi(events);
  api.findIssueByMarker = async () => ({
    iid: 123,
    state: 'closed',
    web_url: 'https://gitlab.test/issues/123',
    description: 'owned marker is supplied by fake lookup',
  });
  let gitCalled = false;
  const manifest = localManifest();
  const preflight = readyPreflight();
  const { confirmationDigest } = previewLocalBatch({ manifest, preflight });
  await assert.rejects(
    publishLocalBatch({
      manifest,
      preflight,
      confirmationDigest,
      git: {},
      api,
      apply: true,
      publishGitState: async () => { gitCalled = true; },
    }),
    /owned Issue #123.*closed.*reopen/i,
  );
  assert.equal(gitCalled, false);
  assert.equal(events.includes('create-issue'), false);
});

test('temporary repository recovery resumes branch-ready and reuses pushed commit state', async () => {
  const { repositoryDir, workspace, state } = await createPublicationFixture();
  const entry = {
    ...localManifest().repositories['base-app'],
    baseSha: state.baseSha,
    snapshotTreeSha: state.snapshotTreeSha,
    paths: state.paths,
  };
  const hookPath = path.join(repositoryDir, '.git', 'hooks', 'pre-commit');
  await writeFile(hookPath, '#!/bin/sh\nexit 1\n');
  await chmod(hookPath, 0o755);
  await assert.rejects(
    publishLocalGitState({ repository: 'base-app', state, entry, issue: { iid: 123 }, git: workspace }),
    /commit/i,
  );
  await writeFile(hookPath, '#!/bin/sh\nexit 0\n');

  const resumed = await publishLocalGitState({
    repository: 'base-app',
    state,
    entry,
    issue: { iid: 123 },
    git: workspace,
  });
  assert.equal(resumed.pushed, true);
  assert.equal(resumed.recovered, true);

  const reused = await publishLocalGitState({
    repository: 'base-app',
    state,
    entry,
    issue: { iid: 123 },
    git: workspace,
  });
  assert.equal(reused.commitSha, resumed.commitSha);
  assert.equal(reused.pushed, false);
  assert.equal(reused.recovered, true);
});

test('three-repository temporary end-to-end publishes dirty repositories and skips clean one', async () => {
  const rootDir = await mkdtemp(path.join(tmpdir(), 'easypos-local-e2e-'));
  const projectPaths = {
    'base-app': 'oxygen/base-app',
    'base-engine': 'oxygen/base-engine',
    'base-web': 'oxygen/base-web',
  };
  const baseShas = {};
  const remoteDirs = {};
  for (const repository of Object.keys(projectPaths)) {
    const repositoryDir = path.join(rootDir, repository);
    const remoteDir = path.join(rootDir, `${repository}.git`);
    remoteDirs[repository] = remoteDir;
    await mkdir(repositoryDir);
    await git(rootDir, ['init', '--bare', remoteDir]);
    await git(repositoryDir, ['init', '-b', '1.1.0-rc']);
    await git(repositoryDir, ['config', 'user.name', 'E2E Test']);
    await git(repositoryDir, ['config', 'user.email', 'e2e@example.test']);
    await writeFile(path.join(repositoryDir, 'feature.txt'), 'base\n');
    await git(repositoryDir, ['add', '-A']);
    await git(repositoryDir, ['commit', '-m', 'base']);
    baseShas[repository] = await git(repositoryDir, ['rev-parse', 'HEAD']);
    await git(repositoryDir, ['remote', 'add', 'origin', remoteDir]);
    await git(repositoryDir, ['push', '--set-upstream', 'origin', '1.1.0-rc']);
    const gitLabRemote = `git@gitlab.sjfood.us:${projectPaths[repository]}.git`;
    await git(repositoryDir, ['remote', 'set-url', 'origin', gitLabRemote]);
    await git(repositoryDir, ['config', `url.${remoteDir}.insteadOf`, gitLabRemote]);
    if (repository !== 'base-web') {
      await writeFile(path.join(repositoryDir, 'feature.txt'), `${repository} change\n`);
      await writeFile(path.join(repositoryDir, 'new.txt'), 'new\n');
    }
  }

  const events = [];
  const issues = new Map();
  const mergeRequests = new Map();
  let nextIssue = 100;
  const api = {
    getProject: async (projectPath) => ({
      id: projectPath,
      permissions: { project_access: { access_level: 30 } },
    }),
    getBranch: async (projectPath, branch) => ({
      name: branch,
      commit: { id: baseShas[projectPath.split('/').at(-1)] },
    }),
    findOpenIssueByMarker: async (projectPath, marker) => {
      const issue = issues.get(projectPath);
      return issue?.description.includes(marker) ? issue : null;
    },
    findIssueByMarker: async (projectPath, marker) => {
      const issue = issues.get(projectPath);
      return issue?.description.includes(marker) ? issue : null;
    },
    createIssue: async (projectPath, payload) => {
      nextIssue += 1;
      events.push(['issue', projectPath]);
      const issue = { iid: nextIssue, web_url: `https://gitlab.test/${projectPath}/issues/${nextIssue}`, ...payload };
      issues.set(projectPath, issue);
      return issue;
    },
    updateIssue: async (projectPath, _iid, payload) => {
      const issue = { ...issues.get(projectPath), ...payload };
      issues.set(projectPath, issue);
      return issue;
    },
    listOpenMergeRequests: async (projectPath) => mergeRequests.has(projectPath)
      ? [mergeRequests.get(projectPath)]
      : [],
    createMergeRequest: async (projectPath, payload) => {
      events.push(['mr', projectPath]);
      const mr = { iid: nextIssue, web_url: `https://gitlab.test/${projectPath}/mr/${nextIssue}`, ...payload };
      mergeRequests.set(projectPath, mr);
      return mr;
    },
    updateMergeRequest: async (projectPath, _iid, payload) => ({
      ...mergeRequests.get(projectPath),
      ...payload,
    }),
  };
  const workspace = new LocalGitWorkspace({ rootDir });
  const preflight = await preflightLocalBatch({ target: 'CURRENT', git: workspace, api });
  const repositories = {};
  for (const state of preflight.repositories.filter(({ status }) => status === 'ready')) {
    repositories[state.repository] = {
      title: `${state.repository} 本地发布`,
      changelog: changelog(`${state.repository} 新功能`),
      baseBranch: state.baseBranch,
      baseSha: state.baseSha,
      snapshotTreeSha: state.snapshotTreeSha,
      paths: state.paths,
      branchSlug: `${state.repository.replace('base-', '')}-change`,
    };
  }
  const manifest = { version: 2, mode: 'local-worktree', repositories };
  const preview = previewLocalBatch({ manifest, preflight });
  const result = await publishLocalBatch({
    manifest,
    preflight,
    confirmationDigest: preview.confirmationDigest,
    git: workspace,
    api,
    apply: true,
  });

  assert.deepEqual(result.results.map(({ status }) => status), ['complete', 'complete', 'skipped']);
  assert.deepEqual(events.map(([kind]) => kind), ['issue', 'mr', 'issue', 'mr']);
  for (const repository of ['base-app', 'base-engine']) {
    const repositoryDir = path.join(rootDir, repository);
    const issue = issues.get(projectPaths[repository]);
    const branch = `feature/${issue.iid}-${repository.replace('base-', '')}-change`;
    assert.equal(await git(repositoryDir, ['branch', '--show-current']), branch);
    assert.equal(
      await git(remoteDirs[repository], ['rev-parse', `refs/heads/${branch}`]),
      await git(repositoryDir, ['rev-parse', 'HEAD']),
    );
  }
  assert.equal(await git(path.join(rootDir, 'base-web'), ['branch', '--show-current']), '1.1.0-rc');
});
