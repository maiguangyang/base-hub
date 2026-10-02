#!/usr/bin/env node

import { execFile as execFileCallback } from 'node:child_process';
import { readFile } from 'node:fs/promises';
import path from 'node:path';
import { promisify } from 'node:util';
import { fileURLToPath, pathToFileURL } from 'node:url';

import {
  ALLOWED_REPOSITORIES,
  GitLabClient,
  assertAllowedRepository,
  assertExpectedProjectPath,
  buildOwnershipMarker,
  loadGitLabConfig,
  parseGitLabRemote,
  validateGitLabLocation,
} from './gitlab-client.mjs';

const execFile = promisify(execFileCallback);
export const REQUIRED_CHANGELOG_SECTIONS = Object.freeze([
  '新增功能',
  '功能优化',
  '问题修复',
  '升级注意事项',
  '测试验证清单',
]);

function requireNonEmptyString(value, field) {
  const normalized = String(value || '').trim();
  if (!normalized) throw new Error(`Manifest field must be a non-empty string: ${field}`);
  return normalized;
}

export function normalizeChangelogEntry(rawEntry, repository) {
  if (!rawEntry || typeof rawEntry !== 'object' || Array.isArray(rawEntry)) {
    throw new Error(`Manifest repository entry must be an object: ${repository}`);
  }
  const title = requireNonEmptyString(rawEntry.title, `${repository}.title`);
  if (/[\r\n]/.test(title)) throw new Error(`${repository}.title must be a single line.`);
  const changelog = requireNonEmptyString(rawEntry.changelog, `${repository}.changelog`);
  for (const section of REQUIRED_CHANGELOG_SECTIONS) {
    if (!new RegExp(`^##\\s+${section}\\s*$`, 'm').test(changelog)) {
      throw new Error(`${repository} is missing changelog section: ${section}`);
    }
  }
  return { title, changelog };
}

function normalizePushedChangelogEntry(rawEntry, repository) {
  return {
    ...normalizeChangelogEntry(rawEntry, repository),
    sourceBranch: requireNonEmptyString(rawEntry.sourceBranch, `${repository}.sourceBranch`),
    sourceSha: requireNonEmptyString(rawEntry.sourceSha, `${repository}.sourceSha`),
    targetBranch: requireNonEmptyString(rawEntry.targetBranch, `${repository}.targetBranch`),
    targetSha: requireNonEmptyString(rawEntry.targetSha, `${repository}.targetSha`),
  };
}

export function normalizeManifest(input) {
  if (!input || typeof input !== 'object' || Array.isArray(input)) {
    throw new Error('Manifest must be a JSON object.');
  }
  if (input.version !== 1) throw new Error('Manifest version must be 1.');

  const source = requireNonEmptyString(input.source, 'source');
  const target = requireNonEmptyString(input.target, 'target');
  if (source === target) throw new Error('Source and target branches must differ.');
  if (!input.repositories || typeof input.repositories !== 'object' || Array.isArray(input.repositories)) {
    throw new Error('Manifest repositories must be an object.');
  }

  const repositories = {};
  for (const [repository, rawEntry] of Object.entries(input.repositories)) {
    assertAllowedRepository(repository);
    repositories[repository] = normalizePushedChangelogEntry(rawEntry, repository);
  }

  return { version: 1, source, target, repositories };
}

export function normalizeIssueRecordManifest(input) {
  if (!input || typeof input !== 'object' || Array.isArray(input)) {
    throw new Error('Issue-only manifest must be a JSON object.');
  }
  if (input.version !== 2 || input.mode !== 'issue-only') {
    throw new Error('Issue-only manifest must use version 2 and mode issue-only.');
  }
  const source = requireNonEmptyString(input.source, 'source');
  const target = requireNonEmptyString(input.target, 'target');
  if (source === target) throw new Error('Issue-only source and target must differ.');
  if (!input.repositories || typeof input.repositories !== 'object' || Array.isArray(input.repositories)) {
    throw new Error('Issue-only repositories must be an object.');
  }
  const repositoryEntries = Object.entries(input.repositories);
  if (!repositoryEntries.length) throw new Error('Issue-only repositories must not be empty.');
  const repositories = {};
  for (const [repository, rawEntry] of repositoryEntries) {
    assertAllowedRepository(repository);
    repositories[repository] = {
      ...normalizeChangelogEntry(rawEntry, repository),
      sourceSha: requireNonEmptyString(rawEntry.sourceSha, `${repository}.sourceSha`),
      targetSha: requireNonEmptyString(rawEntry.targetSha, `${repository}.targetSha`),
    };
  }
  return { version: 2, mode: 'issue-only', source, target, repositories };
}

function buildIssueRecordMarker({ projectPath, source, target, sourceSha, targetSha }) {
  const payload = Buffer.from(JSON.stringify({
    mode: 'issue-only',
    projectPath,
    source,
    target,
    sourceSha,
    targetSha,
  }), 'utf8').toString('base64url');
  return `<!-- easypos-gitlab-issue-record:v2:${payload} -->`;
}

export class GitInspector {
  constructor({ rootDir, execGit = execFile, timeoutMs = 120_000 }) {
    this.rootDir = rootDir;
    this.execGit = execGit;
    this.timeoutMs = timeoutMs;
  }

  async git(repository, args) {
    const cwd = path.join(this.rootDir, assertAllowedRepository(repository));
    try {
      const { stdout } = await this.execGit('git', args, {
        cwd,
        encoding: 'utf8',
        timeout: this.timeoutMs,
        env: { ...process.env, GIT_TERMINAL_PROMPT: '0' },
      });
      return stdout.trim();
    } catch (error) {
      const detail = String(error?.stderr || error?.message || error).trim();
      throw new Error(`${repository}: Git check failed (${args.join(' ')}): ${detail}`);
    }
  }

  async resolveTargetCommit(repository, targetBranch) {
    const candidates = [
      `refs/heads/${targetBranch}`,
      `refs/remotes/origin/${targetBranch}`,
    ];
    for (const ref of candidates) {
      try {
        const sha = await this.git(repository, ['rev-parse', '--verify', `${ref}^{commit}`]);
        return { localTargetRef: ref, localTargetSha: sha };
      } catch {
        // Try the remote-tracking ref when the local target branch is absent.
      }
    }
    throw new Error(`${repository}: target branch is not available as a local or origin-tracking ref: ${targetBranch}`);
  }

  async inspect(repository, requestedSource, requestedTarget) {
    assertAllowedRepository(repository);
    const branch = requestedSource === 'CURRENT'
      ? await this.git(repository, ['symbolic-ref', '--quiet', '--short', 'HEAD'])
      : requireNonEmptyString(requestedSource, 'source');
    const targetBranch = requireNonEmptyString(requestedTarget, 'target');

    await this.git(repository, ['check-ref-format', '--branch', branch]);
    await this.git(repository, ['check-ref-format', '--branch', targetBranch]);
    const [status, localSha, remote, targetCommit] = await Promise.all([
      this.git(repository, ['status', '--porcelain=v1', '--untracked-files=all']),
      this.git(repository, ['rev-parse', '--verify', `refs/heads/${branch}^{commit}`]),
      this.git(repository, ['remote', 'get-url', 'origin']),
      this.resolveTargetCommit(repository, targetBranch),
    ]);
    const { host: remoteHost, projectPath } = parseGitLabRemote(remote);
    assertExpectedProjectPath(repository, projectPath);
    return {
      repository,
      branch,
      dirty: Boolean(status),
      localSha,
      ...targetCommit,
      remoteHost,
      projectPath,
    };
  }

  async resolvePublicationRef(repository, ref) {
    const value = requireNonEmptyString(ref, 'Issue-only Git reference');
    const candidates = /^[a-f0-9]{40}$/.test(value)
      ? [value]
      : [`refs/heads/${value}`, `refs/remotes/origin/${value}`, value];
    for (const candidate of candidates) {
      try {
        return await this.git(repository, ['rev-parse', '--verify', `${candidate}^{commit}`]);
      } catch {
        // Continue through exact local, origin-tracking, and raw commit-ish candidates.
      }
    }
    throw new Error(`${repository}: Issue-only Git reference cannot be resolved locally: ${value}`);
  }

  async inspectIssueRecord(repository, source, target) {
    assertAllowedRepository(repository);
    const [sourceSha, targetSha, remote] = await Promise.all([
      this.resolvePublicationRef(repository, source),
      this.resolvePublicationRef(repository, target),
      this.git(repository, ['config', '--get', 'remote.origin.url']),
    ]);
    let sourceMerged = false;
    try {
      await this.git(repository, ['merge-base', '--is-ancestor', sourceSha, targetSha]);
      sourceMerged = true;
    } catch {
      // Exit 1 means the source is not an ancestor of the target.
    }
    const { host: remoteHost, projectPath } = parseGitLabRemote(remote);
    assertExpectedProjectPath(repository, projectPath);
    return { repository, sourceSha, targetSha, sourceMerged, remoteHost, projectPath };
  }
}

function projectAccessLevel(project) {
  return Math.max(
    Number(project?.permissions?.project_access?.access_level || 0),
    Number(project?.permissions?.group_access?.access_level || 0),
  );
}

export class GitLabApi {
  constructor(client) {
    this.client = client;
  }

  projectPath(projectPath) {
    return `/projects/${encodeURIComponent(projectPath)}`;
  }

  getProject(projectPath) {
    return this.client.request('GET', this.projectPath(projectPath));
  }

  getBranch(projectPath, branch) {
    return this.client.request(
      'GET',
      `${this.projectPath(projectPath)}/repository/branches/${encodeURIComponent(branch)}`,
    );
  }

  getCommit(projectPath, ref) {
    return this.client.request(
      'GET',
      `${this.projectPath(projectPath)}/repository/commits/${encodeURIComponent(ref)}`,
    );
  }

  compare(projectPath, targetBranch, sourceBranch) {
    return this.client.request('GET', `${this.projectPath(projectPath)}/repository/compare`, {
      query: { from: targetBranch, to: sourceBranch },
    });
  }

  listOpenMergeRequests(projectPath, sourceBranch, targetBranch) {
    return this.client.request('GET', `${this.projectPath(projectPath)}/merge_requests`, {
      query: {
        state: 'opened',
        scope: 'all',
        source_branch: sourceBranch,
        target_branch: targetBranch,
        per_page: 100,
      },
    });
  }

  listMergedMergeRequests(projectPath, sourceBranch, targetBranch) {
    return this.client.request('GET', `${this.projectPath(projectPath)}/merge_requests`, {
      query: {
        state: 'merged',
        scope: 'all',
        source_branch: sourceBranch,
        target_branch: targetBranch,
        per_page: 100,
      },
    });
  }

  async findOpenIssueByMarker(projectPath, marker) {
    const issues = await this.client.request('GET', `${this.projectPath(projectPath)}/issues`, {
      query: { state: 'opened', scope: 'all', search: marker, in: 'description', per_page: 100 },
    });
    return (issues || []).find((issue) => String(issue.description || '').includes(marker)) || null;
  }

  async findIssueByMarker(projectPath, marker) {
    const issues = await this.client.request('GET', `${this.projectPath(projectPath)}/issues`, {
      query: { state: 'all', scope: 'all', search: marker, in: 'description', per_page: 100 },
    });
    return (issues || []).find((issue) => String(issue.description || '').includes(marker)) || null;
  }

  createIssue(projectPath, payload) {
    return this.client.request('POST', `${this.projectPath(projectPath)}/issues`, { body: payload });
  }

  updateIssue(projectPath, iid, payload) {
    return this.client.request('PUT', `${this.projectPath(projectPath)}/issues/${iid}`, { body: payload });
  }

  createMergeRequest(projectPath, payload) {
    return this.client.request('POST', `${this.projectPath(projectPath)}/merge_requests`, { body: payload });
  }

  updateMergeRequest(projectPath, iid, payload) {
    return this.client.request('PUT', `${this.projectPath(projectPath)}/merge_requests/${iid}`, { body: payload });
  }
}

function hasComparisonChanges(comparison) {
  if (comparison?.compare_timeout === true) {
    throw new Error('GitLab compare timed out; final differences may be incomplete.');
  }
  if (!Array.isArray(comparison?.diffs)) {
    throw new Error('GitLab compare response is missing the diffs array.');
  }
  return comparison.diffs.length > 0;
}

export async function preflightBatch({ rootDir, source, target, git, api, baseUrl }) {
  const normalizedSource = requireNonEmptyString(source, 'source');
  const normalizedTarget = requireNonEmptyString(target, 'target');
  if (normalizedSource === normalizedTarget) throw new Error('Source and target branches must differ.');
  const inspector = git || new GitInspector({ rootDir });

  const inspected = [];
  const localErrors = [];
  for (const repository of ALLOWED_REPOSITORIES) {
    try {
      const state = await inspector.inspect(repository, normalizedSource, normalizedTarget);
      assertExpectedProjectPath(repository, state.projectPath);
      if (state.dirty) localErrors.push(`${repository}: repository has uncommitted changes`);
      inspected.push(state);
    } catch (error) {
      localErrors.push(error.message);
    }
  }
  if (localErrors.length) throw new Error(`GitLab publication preflight failed:\n- ${localErrors.join('\n- ')}`);

  const repositories = [];
  for (const state of inspected) {
    try {
      if (baseUrl) validateGitLabLocation(baseUrl, state.remoteHost);
      const project = await api.getProject(state.projectPath);
      if (projectAccessLevel(project) < 30) {
        throw new Error('GitLab token requires Developer access or higher.');
      }

      const [sourceBranch, targetBranch, comparison, mergedMergeRequests] = await Promise.all([
        api.getBranch(state.projectPath, state.branch),
        api.getBranch(state.projectPath, normalizedTarget),
        api.compare(state.projectPath, normalizedTarget, state.branch),
        api.listMergedMergeRequests(state.projectPath, state.branch, normalizedTarget),
      ]);
      if (sourceBranch?.commit?.id !== state.localSha) {
        throw new Error('source branch HEAD is not pushed to GitLab at the local SHA.');
      }
      if (targetBranch?.commit?.id !== state.localTargetSha) {
        throw new Error('target branch local SHA does not match remote SHA.');
      }

      const marker = buildOwnershipMarker({
        projectPath: state.projectPath,
        sourceBranch: state.branch,
        targetBranch: normalizedTarget,
      });
      const comparisonHasChanges = hasComparisonChanges(comparison);
      const exactMergedRevision = (mergedMergeRequests || []).find((mergeRequest) => (
        mergeRequest?.diff_refs?.head_sha === state.localSha
        || mergeRequest?.sha === state.localSha
      ));
      if (exactMergedRevision) {
        repositories.push({
          ...state,
          targetBranch: targetBranch.name,
          comparison,
          marker,
          status: 'skipped',
          reason: `source revision already merged by MR !${exactMergedRevision.iid}`,
        });
        continue;
      }
      if (!comparisonHasChanges) {
        repositories.push({
          ...state,
          targetBranch: targetBranch.name,
          comparison,
          marker,
          status: 'skipped',
          reason: 'no remote difference',
        });
        continue;
      }

      const existingIssue = await api.findIssueByMarker(state.projectPath, marker);
      if (existingIssue?.state && existingIssue.state !== 'opened') {
        throw new Error(`owned Issue #${existingIssue.iid} is ${existingIssue.state}; reopen it before publishing.`);
      }
      const openMergeRequests = await api.listOpenMergeRequests(
        state.projectPath,
        state.branch,
        normalizedTarget,
      );
      if ((openMergeRequests || []).length > 1) {
        throw new Error('multiple existing open MRs use the same source and target branches.');
      }
      const existingMergeRequest = openMergeRequests?.[0] || null;
      if (existingMergeRequest && !String(existingMergeRequest.description || '').includes(marker)) {
        throw new Error('existing open MR is not owned by this skill; refusing to overwrite it.');
      }

      repositories.push({
        ...state,
        targetBranch: targetBranch.name,
        comparison,
        marker,
        status: 'ready',
        existingIssue,
        existingMergeRequest,
      });
    } catch (error) {
      throw new Error(`${state.repository}: ${error.message}`);
    }
  }

  return { source: normalizedSource, target: normalizedTarget, repositories };
}

export async function preflightIssueRecordBatch({
  rootDir,
  manifest: rawManifest,
  git,
  api,
  baseUrl,
}) {
  const manifest = normalizeIssueRecordManifest(rawManifest);
  const inspector = git || new GitInspector({ rootDir });
  const inspected = [];
  const localErrors = [];
  for (const repository of Object.keys(manifest.repositories)) {
    try {
      const state = await inspector.inspectIssueRecord(repository, manifest.source, manifest.target);
      assertExpectedProjectPath(repository, state.projectPath);
      inspected.push(state);
    } catch (error) {
      localErrors.push(error.message);
    }
  }
  if (localErrors.length) {
    throw new Error(`Issue-only preflight failed:\n- ${localErrors.join('\n- ')}`);
  }

  const comparisons = {};
  for (const state of inspected) {
    try {
      if (baseUrl) validateGitLabLocation(baseUrl, state.remoteHost);
      const project = await api.getProject(state.projectPath);
      if (projectAccessLevel(project) < 30) {
        throw new Error('GitLab token requires Developer access or higher.');
      }
      const [sourceCommit, targetCommit, comparison] = await Promise.all([
        api.getCommit(state.projectPath, state.sourceSha),
        api.getCommit(state.projectPath, state.targetSha),
        api.compare(state.projectPath, state.targetSha, state.sourceSha),
      ]);
      if (sourceCommit?.id !== state.sourceSha || targetCommit?.id !== state.targetSha) {
        throw new Error('resolved source/target commits do not match GitLab.');
      }
      hasComparisonChanges(comparison);
      comparisons[state.repository] = { ...state, ...comparison };
    } catch (error) {
      throw new Error(`${state.repository}: ${error.message}`);
    }
  }
  return comparisons;
}

function validateManifestCoverage(manifest, preflight) {
  if (manifest.source !== preflight.source || manifest.target !== preflight.target) {
    throw new Error('Manifest source/target does not match the completed preflight.');
  }
  const missing = preflight.repositories
    .filter(({ status }) => status === 'ready')
    .map(({ repository }) => repository)
    .filter((repository) => !manifest.repositories[repository]);
  if (missing.length) throw new Error(`Manifest is missing changed repositories: ${missing.join(', ')}`);
  for (const state of preflight.repositories.filter(({ status }) => status === 'ready')) {
    const entry = manifest.repositories[state.repository];
    if (entry.sourceBranch !== state.branch || entry.sourceSha !== state.localSha) {
      throw new Error(`${state.repository}: source branch or SHA differs from the changelog manifest.`);
    }
    if (entry.targetBranch !== state.targetBranch || entry.targetSha !== state.localTargetSha) {
      throw new Error(`${state.repository}: target branch or SHA differs from the changelog manifest.`);
    }
  }
}

function completedPublicationSummary(results) {
  const completed = results.filter(({ issue, mergeRequest }) => issue && mergeRequest);
  if (!completed.length) return '';
  return ` Completed before failure: ${completed.map((item) => (
    `${item.repository} Issue ${item.issue.url}, MR ${item.mergeRequest.url}`
  )).join('; ')}.`;
}

export async function publishBatch({ manifest: rawManifest, preflight, api, apply = false }) {
  const manifest = normalizeManifest(rawManifest);
  validateManifestCoverage(manifest, preflight);
  const results = [];

  for (const state of preflight.repositories) {
    if (state.status === 'skipped') {
      results.push({ repository: state.repository, status: 'skipped', reason: state.reason });
      continue;
    }
    const entry = manifest.repositories[state.repository];
    if (!apply) {
      results.push({
        repository: state.repository,
        status: 'preview',
        action: state.existingIssue || state.existingMergeRequest ? 'update-or-complete' : 'create',
      });
      continue;
    }

    const issueDescription = `${entry.changelog.trim()}\n\n${state.marker}`;
    let issue;
    try {
      issue = state.existingIssue
        ? await api.updateIssue(state.projectPath, state.existingIssue.iid, {
          title: entry.title,
          description: issueDescription,
        })
        : await api.createIssue(state.projectPath, {
          title: entry.title,
          description: issueDescription,
        });
    } catch (error) {
      throw new Error(
        `${state.repository}: Issue publication failed: ${error.message}.${completedPublicationSummary(results)}`,
      );
    }

    const mergeRequestDescription = [
      `Closes #${issue.iid}`,
      '',
      entry.changelog.trim(),
      '',
      state.marker,
    ].join('\n');
    const mergeRequestPayload = {
      source_branch: state.branch,
      target_branch: state.targetBranch,
      title: entry.title,
      description: mergeRequestDescription,
    };
    let mergeRequest;
    try {
      mergeRequest = state.existingMergeRequest
        ? await api.updateMergeRequest(
          state.projectPath,
          state.existingMergeRequest.iid,
          mergeRequestPayload,
        )
        : await api.createMergeRequest(state.projectPath, mergeRequestPayload);
    } catch (error) {
      throw new Error(
        `${state.repository}: Issue #${issue.iid} (${issue.web_url}) exists, but MR publication failed: ${error.message}.${completedPublicationSummary(results)}`,
      );
    }

    results.push({
      repository: state.repository,
      status: state.existingIssue || state.existingMergeRequest ? 'updated' : 'created',
      issue: { iid: issue.iid, url: issue.web_url },
      mergeRequest: { iid: mergeRequest.iid, url: mergeRequest.web_url },
    });
  }

  return { applied: apply, source: manifest.source, target: manifest.target, results };
}

export async function publishIssueRecords({ manifest: rawManifest, comparison, api, apply = false }) {
  const manifest = normalizeIssueRecordManifest(rawManifest);
  const plans = [];
  for (const repository of Object.keys(manifest.repositories)) {
    const state = comparison?.[repository];
    if (!state) throw new Error(`${repository}: Issue-only comparison result is required.`);
    assertExpectedProjectPath(repository, state.projectPath);
    const entry = manifest.repositories[repository];
    if (entry.sourceSha !== state.sourceSha || entry.targetSha !== state.targetSha) {
      throw new Error(`${repository}: source or target SHA differs from the changelog manifest.`);
    }
    const hasFinalDiff = hasComparisonChanges(state);
    if (!state.sourceMerged && hasFinalDiff) {
      throw new Error(`${repository}: a current final diff exists; use the Merge Request workflow.`);
    }
    plans.push({
      repository,
      projectPath: state.projectPath,
      marker: buildIssueRecordMarker({
        projectPath: state.projectPath,
        source: manifest.source,
        target: manifest.target,
        sourceSha: entry.sourceSha,
        targetSha: entry.targetSha,
      }),
    });
  }
  if (!apply) {
    return {
      applied: false,
      source: manifest.source,
      target: manifest.target,
      results: plans.map(({ repository }) => ({
        repository,
        status: 'preview',
        action: 'create-or-update-release-record-issue',
      })),
    };
  }

  const results = [];
  for (const plan of plans) {
    const entry = manifest.repositories[plan.repository];
    const description = [
      '此 Issue 仅作为已合并或无最终差异的发布变更记录，不代表仍有待合并代码。',
      '',
      entry.changelog.trim(),
      '',
      plan.marker,
    ].join('\n');
    let existingIssue;
    let issue;
    try {
      existingIssue = await api.findIssueByMarker(plan.projectPath, plan.marker);
      issue = existingIssue
        ? await api.updateIssue(plan.projectPath, existingIssue.iid, {
          title: entry.title,
          description,
        })
        : await api.createIssue(plan.projectPath, {
          title: entry.title,
          description,
        });
    } catch (error) {
      const completed = results.length
        ? ` Completed before failure: ${results.map((item) => `${item.repository} Issue ${item.issue.url}`).join('; ')}.`
        : '';
      throw new Error(`${plan.repository}: Issue-only publication failed: ${error.message}.${completed}`);
    }
    results.push({
      repository: plan.repository,
      status: existingIssue ? 'updated' : 'created',
      issue: { iid: issue.iid, url: issue.web_url },
    });
  }
  return { applied: true, source: manifest.source, target: manifest.target, results };
}

function parseArguments(argv) {
  const [command, ...rest] = argv;
  const values = {};
  for (let index = 0; index < rest.length; index += 1) {
    const argument = rest[index];
    if (argument === '--apply') {
      values.apply = true;
      continue;
    }
    if (!argument.startsWith('--') || index + 1 >= rest.length) {
      throw new Error(`Invalid command argument: ${argument}`);
    }
    values[argument.slice(2)] = rest[index + 1];
    index += 1;
  }
  return { command, values };
}

function usage() {
  return [
    'Usage:',
    '  publish-change-requests.mjs preflight --source <branch|CURRENT> --target <branch>',
    '  publish-change-requests.mjs publish --manifest <json-file> [--apply]',
    '  publish-change-requests.mjs local-preflight --target <branch|CURRENT>',
    '  publish-change-requests.mjs local-publish --manifest <json-file> [--apply --confirm <digest>]',
    '  publish-change-requests.mjs issue-record --manifest <json-file> [--apply]',
    '',
    'The publish command is preview-only unless --apply is present.',
  ].join('\n');
}

function validateCommandArguments(command, values) {
  const allowed = {
    preflight: new Set(['source', 'target']),
    publish: new Set(['manifest', 'apply']),
    'local-preflight': new Set(['target']),
    'local-publish': new Set(['manifest', 'apply', 'confirm']),
    'issue-record': new Set(['manifest', 'apply']),
  };
  if (!allowed[command]) throw new Error(usage());
  for (const key of Object.keys(values)) {
    if (!allowed[command].has(key)) throw new Error(`Unknown argument for ${command}: --${key}`);
  }
  if (command === 'preflight' && (!values.source || !values.target)) throw new Error(usage());
  if (command === 'publish' && !values.manifest) throw new Error(usage());
  if (command === 'local-preflight' && !values.target) throw new Error(usage());
  if (command === 'local-publish' && !values.manifest) throw new Error(usage());
  if (command === 'issue-record' && !values.manifest) throw new Error(usage());
  if (command === 'local-publish' && values.apply && !values.confirm) {
    throw new Error('local-publish --apply requires the exact --confirm digest from preview.');
  }
}

export async function runCli(argv, {
  rootDir,
  config: suppliedConfig,
  api: suppliedApi,
  git,
  localModule: suppliedLocalModule,
} = {}) {
  const resolvedRoot = rootDir || path.resolve(path.dirname(fileURLToPath(import.meta.url)), '../../../..');
  const { command, values } = parseArguments(argv);
  validateCommandArguments(command, values);

  const config = suppliedConfig || await loadGitLabConfig({ cwd: resolvedRoot });
  const api = suppliedApi || new GitLabApi(new GitLabClient(config));
  if (command === 'local-preflight' || command === 'local-publish') {
    const localModule = suppliedLocalModule || await import('./local-worktree-publication.mjs');
    if (command === 'local-preflight') {
      const preflight = await localModule.preflightLocalBatch({
        rootDir: resolvedRoot,
        target: values.target,
        git,
        api,
        baseUrl: config.baseUrl,
      });
      return { applied: false, ...preflight };
    }
    const manifest = localModule.normalizeLocalManifest(
      JSON.parse(await readFile(path.resolve(values.manifest), 'utf8')),
    );
    const preflight = await localModule.preflightLocalBatch({
      rootDir: resolvedRoot,
      target: 'CURRENT',
      manifest,
      git,
      api,
      baseUrl: config.baseUrl,
    });
    return localModule.publishLocalBatch({
      manifest,
      preflight,
      confirmationDigest: values.confirm,
      git,
      api,
      apply: Boolean(values.apply),
    });
  }
  if (command === 'issue-record') {
    const manifest = normalizeIssueRecordManifest(
      JSON.parse(await readFile(path.resolve(values.manifest), 'utf8')),
    );
    const comparison = await preflightIssueRecordBatch({
      rootDir: resolvedRoot,
      manifest,
      git,
      api,
      baseUrl: config.baseUrl,
    });
    return publishIssueRecords({
      manifest,
      comparison,
      api,
      apply: Boolean(values.apply),
    });
  }
  let source;
  let target;
  let manifest;

  if (command === 'publish') {
    if (!values.manifest) throw new Error(usage());
    manifest = normalizeManifest(JSON.parse(await readFile(path.resolve(values.manifest), 'utf8')));
    ({ source, target } = manifest);
  } else {
    source = values.source;
    target = values.target;
  }

  const preflight = await preflightBatch({
    rootDir: resolvedRoot,
    source,
    target,
    git,
    api,
    baseUrl: config.baseUrl,
  });
  if (command === 'preflight') return { applied: false, ...preflight };
  return publishBatch({ manifest, preflight, api, apply: Boolean(values.apply) });
}

const isMain = process.argv[1]
  && pathToFileURL(path.resolve(process.argv[1])).href === import.meta.url;
if (isMain) {
  runCli(process.argv.slice(2))
    .then((result) => process.stdout.write(`${JSON.stringify(result, null, 2)}\n`))
    .catch((error) => {
      process.stderr.write(`${error.message}\n`);
      process.exitCode = 1;
    });
}
