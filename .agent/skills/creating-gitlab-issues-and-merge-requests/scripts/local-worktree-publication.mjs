import { execFile as execFileCallback } from 'node:child_process';
import { createHash, timingSafeEqual } from 'node:crypto';
import { access, mkdir, mkdtemp, readFile, rm } from 'node:fs/promises';
import { tmpdir } from 'node:os';
import path from 'node:path';
import { promisify } from 'node:util';

import {
  ALLOWED_REPOSITORIES,
  assertAllowedRepository,
  assertExpectedProjectPath,
  parseGitLabRemote,
  validateGitLabLocation,
} from './gitlab-client.mjs';
import { normalizeChangelogEntry } from './publish-change-requests.mjs';

const execFile = promisify(execFileCallback);
const FULL_SHA = /^[a-f0-9]{40}$/;
const BRANCH_SLUG = /^[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?$/;

function requireString(value, field) {
  const normalized = String(value || '').trim();
  if (!normalized) throw new Error(`Local manifest field must be a non-empty string: ${field}`);
  return normalized;
}

function requireSha(value, field) {
  const normalized = requireString(value, field);
  if (!FULL_SHA.test(normalized)) throw new Error(`${field} must be a full lowercase Git SHA.`);
  return normalized;
}

function normalizeManifestPath(value, repository) {
  const candidate = String(value || '');
  if (
    !candidate
    || candidate.includes('\0')
    || candidate.includes('\\')
    || path.posix.isAbsolute(candidate)
    || candidate.split('/').some((part) => !part || part === '.' || part === '..')
  ) {
    throw new Error(`${repository}.paths contains an unsafe relative path: ${candidate}`);
  }
  return candidate;
}

export function validateBranchSlug(value) {
  const slug = requireString(value, 'branchSlug');
  if (!BRANCH_SLUG.test(slug)) {
    throw new Error('branchSlug must contain 1-63 lowercase ASCII letters, digits, or hyphens.');
  }
  return slug;
}

export function normalizeLocalManifest(input) {
  if (!input || typeof input !== 'object' || Array.isArray(input)) {
    throw new Error('Local manifest must be a JSON object.');
  }
  if (input.version !== 2) throw new Error('Local manifest version must be 2.');
  if (input.mode !== 'local-worktree') throw new Error('Local manifest mode must be local-worktree.');
  if (!input.repositories || typeof input.repositories !== 'object' || Array.isArray(input.repositories)) {
    throw new Error('Local manifest repositories must be an object.');
  }
  const entries = Object.entries(input.repositories);
  if (!entries.length) throw new Error('Local manifest repositories must not be empty.');
  for (const [repository] of entries) assertAllowedRepository(repository);

  const repositories = {};
  for (const repository of ALLOWED_REPOSITORIES) {
    const rawEntry = input.repositories[repository];
    if (!rawEntry) continue;
    const { title, changelog } = normalizeChangelogEntry(rawEntry, repository);
    if (!Array.isArray(rawEntry.paths) || !rawEntry.paths.length) {
      throw new Error(`${repository}.paths must be a non-empty array.`);
    }
    const paths = rawEntry.paths.map((value) => normalizeManifestPath(value, repository)).sort();
    if (new Set(paths).size !== paths.length) throw new Error(`${repository}.paths contains duplicate paths.`);
    repositories[repository] = {
      title,
      changelog,
      baseBranch: requireString(rawEntry.baseBranch, `${repository}.baseBranch`),
      baseSha: requireSha(rawEntry.baseSha, `${repository}.baseSha`),
      snapshotTreeSha: requireSha(rawEntry.snapshotTreeSha, `${repository}.snapshotTreeSha`),
      paths,
      branchSlug: validateBranchSlug(rawEntry.branchSlug),
    };
  }

  return { version: 2, mode: 'local-worktree', repositories };
}

export class LocalGitWorkspace {
  constructor({ rootDir, execGit = execFile, timeoutMs = 120_000 }) {
    this.rootDir = rootDir;
    this.execGit = execGit;
    this.timeoutMs = timeoutMs;
  }

  repositoryDir(repository) {
    return path.join(this.rootDir, assertAllowedRepository(repository));
  }

  async git(repository, args, { env, trimOutput = true } = {}) {
    const cwd = this.repositoryDir(repository);
    try {
      const { stdout } = await this.execGit('git', args, {
        cwd,
        encoding: 'utf8',
        timeout: this.timeoutMs,
        env: {
          ...process.env,
          ...(env || {}),
          GIT_TERMINAL_PROMPT: '0',
        },
      });
      return trimOutput ? stdout.trim() : stdout;
    } catch (error) {
      const detail = String(error?.stderr || error?.message || error).trim();
      throw new Error(`${repository}: Git check failed (${args.join(' ')}): ${detail}`);
    }
  }

  async gitOptional(repository, args) {
    try {
      return await this.git(repository, args);
    } catch {
      return '';
    }
  }

  async operationInProgress(repository) {
    for (const [name, ref] of [
      ['merge', 'MERGE_HEAD'],
      ['cherry-pick', 'CHERRY_PICK_HEAD'],
      ['revert', 'REVERT_HEAD'],
    ]) {
      if (await this.gitOptional(repository, ['rev-parse', '--quiet', '--verify', ref])) return name;
    }
    for (const [name, gitPath] of [['rebase', 'rebase-merge'], ['rebase', 'rebase-apply']]) {
      const value = await this.git(repository, ['rev-parse', '--git-path', gitPath]);
      try {
        await access(path.resolve(this.repositoryDir(repository), value));
        return name;
      } catch {
        // The operation directory is absent.
      }
    }
    return null;
  }

  async prospectiveState(repository) {
    const objectPathValue = await this.git(repository, ['rev-parse', '--git-path', 'objects']);
    const objectPath = path.resolve(this.repositoryDir(repository), objectPathValue);
    const snapshotDir = await mkdtemp(path.join(tmpdir(), 'easypos-git-snapshot-'));
    const temporaryObjectPath = path.join(snapshotDir, 'objects');
    await mkdir(temporaryObjectPath);
    const snapshotEnv = {
      ...process.env,
      GIT_INDEX_FILE: path.join(snapshotDir, 'index'),
      GIT_OBJECT_DIRECTORY: temporaryObjectPath,
      GIT_ALTERNATE_OBJECT_DIRECTORIES: [
        objectPath,
        process.env.GIT_ALTERNATE_OBJECT_DIRECTORIES,
      ].filter(Boolean).join(path.delimiter),
    };

    try {
      await this.git(repository, ['read-tree', 'HEAD'], { env: snapshotEnv });
      await this.git(repository, ['add', '-A'], { env: snapshotEnv });
      const snapshotTreeSha = await this.git(repository, ['write-tree'], { env: snapshotEnv });
      const changedPaths = await this.git(
        repository,
        ['diff-index', '--cached', '--name-only', '-z', 'HEAD'],
        { env: snapshotEnv, trimOutput: false },
      );
      return {
        snapshotTreeSha,
        paths: changedPaths.split('\0').filter(Boolean).sort(),
      };
    } finally {
      await rm(snapshotDir, { recursive: true, force: true });
    }
  }

  async snapshot(repository) {
    assertAllowedRepository(repository);
    const [
      baseBranch,
      baseSha,
      status,
      remote,
      indexPathValue,
      conflictsValue,
      operationInProgress,
      gitUserName,
      gitUserEmail,
    ] = await Promise.all([
      this.git(repository, ['symbolic-ref', '--quiet', '--short', 'HEAD']),
      this.git(repository, ['rev-parse', '--verify', 'HEAD^{commit}']),
      this.git(repository, ['status', '--porcelain=v1', '--untracked-files=all'], { trimOutput: false }),
      this.git(repository, ['config', '--get', 'remote.origin.url']),
      this.git(repository, ['rev-parse', '--git-path', 'index']),
      this.git(repository, ['diff', '--name-only', '--diff-filter=U', '-z'], { trimOutput: false }),
      this.operationInProgress(repository),
      this.gitOptional(repository, ['config', '--get', 'user.name']),
      this.gitOptional(repository, ['config', '--get', 'user.email']),
    ]);
    const cwd = this.repositoryDir(repository);
    const indexPath = path.resolve(cwd, indexPathValue);
    const beforeIndex = await readFile(indexPath);
    const { host: remoteHost, projectPath } = parseGitLabRemote(remote);
    assertExpectedProjectPath(repository, projectPath);

    const { snapshotTreeSha, paths } = await this.prospectiveState(repository);

    const [afterIndex, afterHead, afterBranch, afterStatus] = await Promise.all([
      readFile(indexPath),
      this.git(repository, ['rev-parse', '--verify', 'HEAD^{commit}']),
      this.git(repository, ['symbolic-ref', '--quiet', '--short', 'HEAD']),
      this.git(repository, ['status', '--porcelain=v1', '--untracked-files=all'], { trimOutput: false }),
    ]);
    if (
      !beforeIndex.equals(afterIndex)
      || afterHead !== baseSha
      || afterBranch !== baseBranch
      || afterStatus !== status
    ) {
      throw new Error(`${repository}: repository changed while creating the local snapshot.`);
    }

    return {
      repository,
      baseBranch,
      baseSha,
      snapshotTreeSha,
      paths,
      dirty: paths.length > 0,
      remoteHost,
      projectPath,
      conflicts: conflictsValue.split('\0').filter(Boolean),
      operationInProgress,
      gitUserName,
      gitUserEmail,
    };
  }
}

function accessLevel(project) {
  return Math.max(
    Number(project?.permissions?.project_access?.access_level || 0),
    Number(project?.permissions?.group_access?.access_level || 0),
  );
}

export async function preflightLocalBatch({ rootDir, target, manifest: rawManifest, git, api, baseUrl }) {
  const normalizedTarget = requireString(target, 'target');
  const manifest = rawManifest ? normalizeLocalManifest(rawManifest) : null;
  const workspace = git || new LocalGitWorkspace({ rootDir });
  const inspected = [];
  const localErrors = [];

  for (const repository of ALLOWED_REPOSITORIES) {
    try {
      const currentState = await workspace.snapshot(repository);
      assertExpectedProjectPath(repository, currentState.projectPath);
      const entry = manifest?.repositories[repository];
      if (!entry && !currentState.dirty) {
        inspected.push({ ...currentState, status: 'skipped', reason: 'no local changes' });
        continue;
      }
      if (manifest && !entry) {
        localErrors.push(`${repository}: local changes are not covered by the publication manifest.`);
        continue;
      }
      let state = currentState;
      if (entry) {
        const currentBranch = currentState.baseBranch;
        const currentSha = currentState.baseSha;
        if (currentBranch === entry.baseBranch) {
          if (
            !currentState.dirty
            || currentSha !== entry.baseSha
            || currentState.snapshotTreeSha !== entry.snapshotTreeSha
            || !samePaths(currentState.paths, entry.paths)
          ) {
            localErrors.push(`${repository}: local base worktree drifted from the publication manifest.`);
          }
        } else {
          const recoveryPattern = new RegExp(`^feature/([1-9][0-9]*)-${entry.branchSlug}$`);
          const match = currentBranch.match(recoveryPattern);
          if (!match) {
            localErrors.push(`${repository}: current branch is not the manifest base or Issue-derived recovery branch.`);
          } else if (currentState.dirty) {
            if (
              currentSha !== entry.baseSha
              || currentState.snapshotTreeSha !== entry.snapshotTreeSha
              || !samePaths(currentState.paths, entry.paths)
            ) {
              localErrors.push(`${repository}: branch-ready worktree drifted from the publication manifest.`);
            }
          } else if (typeof workspace.git === 'function') {
            try {
              const [parent, tree, message] = await Promise.all([
                workspace.git(repository, ['rev-parse', 'HEAD^']),
                workspace.git(repository, ['rev-parse', 'HEAD^{tree}']),
                workspace.git(repository, ['log', '-1', '--format=%B']),
              ]);
              if (parent !== entry.baseSha || tree !== entry.snapshotTreeSha) {
                throw new Error('recovery commit parent or tree does not match the manifest.');
              }
              assertCommitMessage(message, Number(match[1]), entry.title, entry.snapshotTreeSha);
            } catch (error) {
              localErrors.push(`${repository}: ${error.message}`);
            }
          }
        }
        state = {
          ...currentState,
          currentBranch,
          currentSha,
          recoveryIssueIid: currentBranch === entry.baseBranch
            ? null
            : Number(currentBranch.match(/^feature\/([1-9][0-9]*)-/)?.[1] || 0),
          baseBranch: entry.baseBranch,
          baseSha: entry.baseSha,
          snapshotTreeSha: entry.snapshotTreeSha,
          paths: entry.paths,
          dirty: true,
        };
      }
      if (normalizedTarget !== 'CURRENT' && normalizedTarget !== state.baseBranch) {
        localErrors.push(`${repository}: target must match current branch/publication base ${state.baseBranch}.`);
      }
      if (currentState.conflicts?.length) {
        localErrors.push(`${repository}: unresolved conflicts: ${currentState.conflicts.join(', ')}`);
      }
      if (currentState.operationInProgress) {
        localErrors.push(`${repository}: ${currentState.operationInProgress} operation is in progress.`);
      }
      if (!String(currentState.gitUserName || '').trim()) {
        localErrors.push(`${repository}: Git user.name is not configured.`);
      }
      if (!String(currentState.gitUserEmail || '').trim()) {
        localErrors.push(`${repository}: Git user.email is not configured.`);
      }
      inspected.push({ ...state, status: 'local-checked' });
    } catch (error) {
      localErrors.push(error.message);
    }
  }
  if (localErrors.length) {
    throw new Error(`Local publication preflight failed:\n- ${localErrors.join('\n- ')}`);
  }

  const repositories = [];
  for (const state of inspected) {
    if (state.status === 'skipped') {
      repositories.push(state);
      continue;
    }
    try {
      if (baseUrl) validateGitLabLocation(baseUrl, state.remoteHost);
      const project = await api.getProject(state.projectPath);
      if (accessLevel(project) < 30) throw new Error('GitLab token requires Developer access or higher.');
      const targetBranch = await api.getBranch(state.projectPath, state.baseBranch);
      if (targetBranch?.commit?.id !== state.baseSha) {
        throw new Error('base branch HEAD is not pushed to GitLab at the local SHA.');
      }
      repositories.push({ ...state, status: 'ready', targetBranch: state.baseBranch });
    } catch (error) {
      throw new Error(`${state.repository}: ${error.message}`);
    }
  }
  return { target: normalizedTarget, repositories };
}

function samePaths(left, right) {
  return JSON.stringify([...(left || [])].sort()) === JSON.stringify([...(right || [])].sort());
}

function validateManifestCoverage(manifest, preflight) {
  const ready = preflight.repositories.filter(({ status }) => status === 'ready');
  const expected = ready.map(({ repository }) => repository);
  const actual = Object.keys(manifest.repositories);
  if (JSON.stringify(expected) !== JSON.stringify(actual)) {
    throw new Error(`Local publication drift: manifest repositories must be ${expected.join(', ')}.`);
  }
  for (const state of ready) {
    const entry = manifest.repositories[state.repository];
    if (
      entry.baseBranch !== state.baseBranch
      || entry.baseSha !== state.baseSha
      || entry.snapshotTreeSha !== state.snapshotTreeSha
      || !samePaths(entry.paths, state.paths)
    ) {
      throw new Error(`${state.repository}: local publication drift detected after preview.`);
    }
  }
}

function digestManifest(manifest) {
  return createHash('sha256').update(JSON.stringify(manifest), 'utf8').digest('hex');
}

export function previewLocalBatch({ manifest: rawManifest, preflight }) {
  const manifest = normalizeLocalManifest(rawManifest);
  validateManifestCoverage(manifest, preflight);
  const results = preflight.repositories.map((state) => {
    if (state.status === 'skipped') {
      return { repository: state.repository, status: 'skipped', reason: state.reason };
    }
    const entry = manifest.repositories[state.repository];
    return {
      repository: state.repository,
      status: 'preview',
      paths: entry.paths,
      title: entry.title,
      targetBranch: state.baseBranch,
      sourceBranchPattern: `feature/<issue_iid>-${entry.branchSlug}`,
      commitMessage: `Refs #<issue_iid>: ${entry.title}`,
    };
  });
  return { confirmationDigest: digestManifest(manifest), results };
}

export function validateLocalApplyState({ manifest: rawManifest, preflight, confirmationDigest }) {
  const manifest = normalizeLocalManifest(rawManifest);
  validateManifestCoverage(manifest, preflight);
  const digestText = String(confirmationDigest || '');
  if (!/^[a-f0-9]{64}$/.test(digestText)) {
    throw new Error('Local publication confirmation digest must be exactly 64 lowercase hexadecimal characters.');
  }
  const expected = Buffer.from(digestManifest(manifest), 'hex');
  const candidate = Buffer.from(digestText, 'hex');
  if (candidate.length !== expected.length || !timingSafeEqual(candidate, expected)) {
    throw new Error('Local publication confirmation digest does not match the preview.');
  }
}

export function deriveFeatureBranch(issueIid, branchSlug) {
  const iid = Number(issueIid);
  if (!Number.isSafeInteger(iid) || iid <= 0) {
    throw new Error('GitLab Issue IID must be a positive integer.');
  }
  return `feature/${iid}-${validateBranchSlug(branchSlug)}`;
}

function parseRemoteSha(value) {
  const line = String(value || '').trim();
  if (!line) return '';
  const [sha] = line.split(/\s+/);
  if (!FULL_SHA.test(sha)) throw new Error('Git remote returned an invalid branch SHA.');
  return sha;
}

function requireMatchingPublicationState(repository, state, entry) {
  if (
    entry.baseBranch !== state.baseBranch
    || entry.baseSha !== state.baseSha
    || entry.snapshotTreeSha !== state.snapshotTreeSha
    || !samePaths(entry.paths, state.paths)
  ) {
    throw new Error(`${repository}: local publication state does not match the confirmed snapshot.`);
  }
}

async function requireGitIdentity(repository, git) {
  const [name, email] = await Promise.all([
    git.gitOptional(repository, ['config', '--get', 'user.name']),
    git.gitOptional(repository, ['config', '--get', 'user.email']),
  ]);
  if (!String(name || '').trim()) throw new Error(`${repository}: Git user.name is not configured.`);
  if (!String(email || '').trim()) throw new Error(`${repository}: Git user.email is not configured.`);
}

function assertCommitMessage(message, issueIid, title, snapshotTreeSha) {
  const lines = String(message || '').trim().split('\n');
  if (lines[0] !== `Refs #${issueIid}: ${title}`) {
    throw new Error('Created commit subject does not match the owned Issue.');
  }
  const trailer = `EasyPos-Change-Request: ${snapshotTreeSha}`;
  if (!lines.includes(trailer)) {
    throw new Error('Created commit is missing the exact EasyPos ownership trailer.');
  }
}

function literalPathspec(changedPath) {
  return `:(top,literal)${changedPath}`;
}

export async function publishLocalGitState({ repository, state, entry, issue, git }) {
  assertAllowedRepository(repository);
  requireMatchingPublicationState(repository, state, entry);
  const currentRemote = parseGitLabRemote(
    await git.git(repository, ['config', '--get', 'remote.origin.url']),
  );
  assertExpectedProjectPath(repository, currentRemote.projectPath);
  if (currentRemote.host !== state.remoteHost) {
    throw new Error(`${repository}: origin host drifted from ${state.remoteHost} to ${currentRemote.host}.`);
  }
  const branch = deriveFeatureBranch(issue?.iid, entry.branchSlug);
  const branchRef = `refs/heads/${branch}`;
  await requireGitIdentity(repository, git);

  const [currentBranch, currentHead, localBranch] = await Promise.all([
    git.gitOptional(repository, ['symbolic-ref', '--quiet', '--short', 'HEAD']),
    git.gitOptional(repository, ['rev-parse', '--verify', 'HEAD^{commit}']),
    git.gitOptional(repository, ['rev-parse', '--verify', `${branchRef}^{commit}`]),
  ]);
  const remoteBefore = parseRemoteSha(
    await git.git(repository, ['ls-remote', '--heads', 'origin', branchRef]),
  );

  let recovered = false;
  let shouldCommit = false;
  let shouldStage = false;
  if (!localBranch) {
    if (remoteBefore) throw new Error(`${repository}: remote branch collision at ${branch}.`);
    if (currentBranch !== state.baseBranch || currentHead !== state.baseSha) {
      throw new Error(`${repository}: current branch or HEAD drifted before feature branch creation.`);
    }
    await git.git(repository, ['switch', '-c', branch, state.baseSha]);
    shouldCommit = true;
    shouldStage = true;
  } else {
    if (currentBranch !== branch || currentHead !== localBranch) {
      throw new Error(`${repository}: local branch collision at ${branch}.`);
    }
    recovered = true;
    if (localBranch === state.baseSha) {
      if (remoteBefore) throw new Error(`${repository}: remote branch collision at ${branch}.`);
      if (typeof git.prospectiveState === 'function') {
        const recoverySnapshot = await git.prospectiveState(repository);
        if (
          recoverySnapshot.snapshotTreeSha !== entry.snapshotTreeSha
          || !samePaths(recoverySnapshot.paths, entry.paths)
        ) {
          throw new Error(`${repository}: branch-ready worktree does not match the confirmed snapshot.`);
        }
      }
      shouldCommit = true;
      shouldStage = await git.git(repository, ['write-tree']) !== entry.snapshotTreeSha;
    }
  }

  if (shouldCommit) {
    if (shouldStage) {
      if (recovered) {
        for (const changedPath of entry.paths) {
          await git.gitOptional(repository, ['add', '-A', '--', literalPathspec(changedPath)]);
        }
      } else {
        await git.git(repository, ['add', '-A', '--', ...entry.paths.map(literalPathspec)]);
      }
    }
    if (await git.git(repository, ['write-tree']) !== entry.snapshotTreeSha) {
      throw new Error(`${repository}: staged tree does not match the confirmed snapshot.`);
    }
    await git.git(repository, [
      'commit',
      '-m',
      `Refs #${issue.iid}: ${entry.title}`,
      '-m',
      `EasyPos-Change-Request: ${entry.snapshotTreeSha}`,
    ]);
  }

  const [parentSha, treeSha, status, message, commitSha] = await Promise.all([
    git.git(repository, ['rev-parse', 'HEAD^']),
    git.git(repository, ['rev-parse', 'HEAD^{tree}']),
    git.git(repository, ['status', '--porcelain=v1', '--untracked-files=all']),
    git.git(repository, ['log', '-1', '--format=%B']),
    git.git(repository, ['rev-parse', 'HEAD']),
  ]);
  if (parentSha !== state.baseSha) {
    throw new Error(`${repository}: created commit parent differs from the confirmed base SHA.`);
  }
  if (treeSha !== entry.snapshotTreeSha) {
    throw new Error(`${repository}: created commit tree differs from the confirmed snapshot tree.`);
  }
  if (status) {
    throw new Error(`${repository}: worktree changed while creating the publication commit.`);
  }
  assertCommitMessage(message, issue.iid, entry.title, entry.snapshotTreeSha);

  if (remoteBefore) {
    if (remoteBefore !== commitSha) {
      throw new Error(`${repository}: remote branch collision at ${branch}.`);
    }
    return { branch, commitSha, pushed: false, recovered: true };
  }
  await git.git(repository, ['push', '--set-upstream', 'origin', branch]);
  const remoteAfter = parseRemoteSha(
    await git.git(repository, ['ls-remote', '--heads', 'origin', branchRef]),
  );
  if (remoteAfter !== commitSha) {
    throw new Error(`${repository}: pushed remote SHA does not match the verified local commit.`);
  }
  return { branch, commitSha, pushed: true, recovered };
}

export function buildLocalOwnershipMarker({ projectPath, baseBranch, baseSha, snapshotTreeSha }) {
  const payload = Buffer.from(JSON.stringify({
    projectPath: requireString(projectPath, 'projectPath'),
    baseBranch: requireString(baseBranch, 'baseBranch'),
    baseSha: requireSha(baseSha, 'baseSha'),
    snapshotTreeSha: requireSha(snapshotTreeSha, 'snapshotTreeSha'),
  }), 'utf8').toString('base64url');
  return `<!-- easypos-gitlab-local-change-request:v2:${payload} -->`;
}

function localProgressSummary(results) {
  const progress = results.filter(({ issue, gitState, mergeRequest }) => issue || gitState || mergeRequest);
  if (!progress.length) return '';
  return ` Completed before failure: ${progress.map((item) => {
    const parts = [item.repository];
    if (item.issue) parts.push(`Issue ${item.issue.url}`);
    if (item.gitState) parts.push(`commit ${item.gitState.commitSha}`);
    if (item.mergeRequest) parts.push(`MR ${item.mergeRequest.url}`);
    return parts.join(' ');
  }).join('; ')}.`;
}

export async function publishLocalBatch({
  manifest: rawManifest,
  preflight,
  confirmationDigest,
  git,
  api,
  apply = false,
  publishGitState = publishLocalGitState,
}) {
  const manifest = normalizeLocalManifest(rawManifest);
  if (!apply) {
    const preview = previewLocalBatch({ manifest, preflight });
    return { applied: false, ...preview };
  }
  validateLocalApplyState({ manifest, preflight, confirmationDigest });
  const results = [];

  for (const state of preflight.repositories) {
    if (state.status === 'skipped') {
      results.push({ repository: state.repository, status: 'skipped', reason: state.reason });
      continue;
    }
    const entry = manifest.repositories[state.repository];
    const marker = buildLocalOwnershipMarker({
      projectPath: state.projectPath,
      baseBranch: state.baseBranch,
      baseSha: state.baseSha,
      snapshotTreeSha: state.snapshotTreeSha,
    });
    const issueDescription = `${entry.changelog.trim()}\n\n${marker}`;
    let existingIssue;
    let issue;
    try {
      existingIssue = await api.findIssueByMarker(state.projectPath, marker);
      if (existingIssue?.state && existingIssue.state !== 'opened') {
        throw new Error(
          `owned Issue #${existingIssue.iid} is ${existingIssue.state}; reopen it before local recovery.`,
        );
      }
      issue = existingIssue
        ? await api.updateIssue(state.projectPath, existingIssue.iid, {
          title: entry.title,
          description: issueDescription,
        })
        : await api.createIssue(state.projectPath, {
          title: entry.title,
          description: issueDescription,
        });
    } catch (error) {
      throw new Error(
        `${state.repository}: Issue publication failed: ${error.message}.${localProgressSummary(results)}`,
      );
    }

    const partial = {
      repository: state.repository,
      status: 'issue-ready',
      issue: { iid: issue.iid, url: issue.web_url },
    };
    results.push(partial);
    let gitState;
    try {
      gitState = await publishGitState({ repository: state.repository, state, entry, issue, git });
      partial.status = gitState.pushed ? 'push-ready' : 'commit-ready';
      partial.gitState = gitState;
    } catch (error) {
      throw new Error(
        `${state.repository}: Issue #${issue.iid} (${issue.web_url}) exists, but Git publication failed: ${error.message}.${localProgressSummary(results)}`,
      );
    }

    const mergeRequestDescription = [
      `Closes #${issue.iid}`,
      '',
      entry.changelog.trim(),
      '',
      marker,
    ].join('\n');
    const mergeRequestPayload = {
      source_branch: gitState.branch,
      target_branch: state.targetBranch,
      title: entry.title,
      description: mergeRequestDescription,
    };
    try {
      const openMergeRequests = await api.listOpenMergeRequests(
        state.projectPath,
        gitState.branch,
        state.targetBranch,
      );
      if ((openMergeRequests || []).length > 1) {
        throw new Error('multiple existing open MRs use the generated source and target branches.');
      }
      const existingMergeRequest = openMergeRequests?.[0] || null;
      if (existingMergeRequest && !String(existingMergeRequest.description || '').includes(marker)) {
        throw new Error('existing open MR is not owned by this local publication; refusing to overwrite it.');
      }
      const mergeRequest = existingMergeRequest
        ? await api.updateMergeRequest(
          state.projectPath,
          existingMergeRequest.iid,
          mergeRequestPayload,
        )
        : await api.createMergeRequest(state.projectPath, mergeRequestPayload);
      partial.status = 'complete';
      partial.mergeRequest = { iid: mergeRequest.iid, url: mergeRequest.web_url };
      partial.action = existingIssue || existingMergeRequest ? 'updated-or-reused' : 'created';
    } catch (error) {
      throw new Error(
        `${state.repository}: Issue #${issue.iid} (${issue.web_url}) exists, but MR publication failed: ${error.message}.${localProgressSummary(results)}`,
      );
    }
  }

  return { applied: true, results };
}
