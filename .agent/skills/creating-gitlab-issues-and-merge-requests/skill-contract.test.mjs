import assert from 'node:assert/strict';
import { readFile } from 'node:fs/promises';
import test from 'node:test';

const skillUrl = new URL('./SKILL.md', import.meta.url);
const constitutionCoreUrl = new URL('../../memory/constitution-core.md', import.meta.url);
const constitutionUrl = new URL('../../memory/constitution.md', import.meta.url);

test('skill preserves repository, output, authorization, and single-agent contracts', async () => {
  const text = await readFile(skillUrl, 'utf8');

  assert.match(text, /^---\nname: creating-gitlab-issues-and-merge-requests\n/m);
  assert.match(text, /Use when/);
  for (const repository of ['base-app', 'base-engine', 'base-web']) assert.match(text, new RegExp(`\\b${repository}\\b`));
  assert.match(text, /根仓库.*(?:禁止|不得|不创建)/);
  assert.match(text, /不得使用 sub agent|禁止使用 sub agent/i);
  const changelogDependency = text.indexOf('**REQUIRED SUB-SKILL:** Use $generating-git-changelogs');
  const manifestCreation = text.indexOf('创建 UTF-8 临时 manifest');
  assert.notEqual(changelogDependency, -1, 'generating-git-changelogs must be a required sub-skill');
  assert.equal(
    changelogDependency < manifestCreation,
    true,
    'the required changelog sub-skill must run before manifest creation',
  );
  assert.match(text, /localSha/);
  assert.match(text, /localTargetSha/);
  assert.match(text, /sourceBranch.*sourceSha.*targetBranch.*targetSha/is);
  assert.match(text, /Issue-only.*sourceSha.*targetSha/is);
  for (const heading of ['新增功能', '功能优化', '问题修复', '升级注意事项', '测试验证清单']) assert.match(text, new RegExp(heading));
  assert.match(text, /Closes #<issue_iid>/);
  assert.match(text, /--apply/);
  assert.match(text, /本地未提交.*local-preflight.*preview.*单独.*确认.*local-publish/is);
  assert.match(text, /已推送.*preflight.*publish/is);
  assert.match(text, /已合并|无最终差异/);
  assert.match(text, /不创建.*(?:空\s*)?MR|禁止.*(?:空\s*)?MR/);
  assert.match(text, /显式.*Issue-only|明确.*Issue-only/i);
  assert.match(text, /feature\/<issue_iid>-<branchSlug>/);
  assert.match(text, /失败.*保留.*manifest|manifest.*失败.*保留/is);
  assert.match(text, /成功|取消/);
  assert.match(text, /宪章例外.*仅.*local-worktree/is);
  assert.match(text, /agent.*自动串联.*generating-git-changelogs/is);
  assert.match(text, /当前分支.*1\.1\.0-rc/);
  assert.match(text, /merged MR.*sourceSha|sourceSha.*merged MR/is);
  assert.match(text, /120.*(?:秒|second).*GIT_TERMINAL_PROMPT/is);
});

test('constitution permits only confirmed local-worktree publication mutations', async () => {
  const [core, full] = await Promise.all([
    readFile(constitutionCoreUrl, 'utf8'),
    readFile(constitutionUrl, 'utf8'),
  ]);

  for (const text of [core, full]) {
    assert.match(text, /creating-gitlab-issues-and-merge-requests/);
    assert.match(text, /preview.*separate.*confirm/is);
    assert.match(text, /base-app.*base-engine.*base-web/s);
    assert.match(text, /force-push.*reset.*clean.*stash/is);
  }
});
