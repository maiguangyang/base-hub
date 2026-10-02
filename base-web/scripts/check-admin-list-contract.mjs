import { readFileSync, readdirSync, statSync } from 'node:fs';
import { join } from 'node:path';
import { fileURLToPath } from 'node:url';
import { analyzeListPage } from './admin-list-contract.mjs';
import { analyzeAdminDialog, dialogCoverageFailures } from './admin-dialog-contract.mjs';

// WEB-ADMIN-001 / rules §4.6：后台列表结构契约自动发现列表页与其他列表工作台，
// 禁止维护容易漏页的手工页面白名单。
const pagesRoot = fileURLToPath(new URL('../src/features/', import.meta.url));

/** 递归收集页面目录的 TSX；实际列表工作台由 AST 检查结果识别。 */
function collectListPages(dir) {
  const found = [];
  for (const entry of readdirSync(dir)) {
    const full = join(dir, entry);
    if (statSync(full).isDirectory()) found.push(...collectListPages(full));
    else if (entry.endsWith('.tsx') && !entry.endsWith('.test.tsx')) found.push(full);
  }
  return found;
}

/** 截出 src/ 起的相对路径用于报错。路径不含 src/ 属于目录结构异常，显式报错而非静默降级。 */
function toRelative(file) {
  const index = file.indexOf('src/');
  if (index === -1) throw new Error(`发现 src/ 之外的列表页，目录结构异常: ${file}`);
  return file.slice(index);
}

const failures = [];
const pages = collectListPages(pagesRoot);
let checkedPages = 0;
let formSurfaceFileCount = 0;

if (pages.length === 0) {
  failures.push('未发现任何后台页面；契约检查失去意义，请确认页面目录');
}

for (const file of pages) {
  const relative = toRelative(file);
  const result = analyzeListPage(readFileSync(file, 'utf8'));
  if (result.hasFormSurface) formSurfaceFileCount += 1;
  if (result.hasMissingFormSurfaceHook) {
    failures.push(`${relative} 的自定义表单触发器或浮层未接入 data-admin-form-surface；WEB-UI-009 要求共享语义覆盖新增控件`);
  }
  if (result.hasInvalidFormSurfaceException) {
    failures.push(`${relative} 使用了未登记的 data-admin-form-surface-exception；例外必须是契约允许的静态值`);
  }
  if (result.hasPageLocalFormSurfaceColor) {
    failures.push(`${relative} 在表单控件或复合字段外壳上局部覆盖白底黑字；WEB-UI-009 要求由共享表单表面语义统一持有`);
  }
  const strictListPage = file.endsWith('ListPage.tsx');
  if (!strictListPage && !result.hasTableShell) continue;
  checkedPages += 1;
  if (strictListPage) {
    for (const name of result.missingImports) {
      failures.push(`${relative} 未引入 §4.6 要求的 ${name}`);
    }
    for (const name of result.notRendered) {
      failures.push(`${relative} 引入了 ${name} 但未渲染；§4.6 要求的是实际装配，不是声明`);
    }
  }
  if (result.hasForbiddenColor) {
    failures.push(`${relative} 手写了状态颜色，§4.4 要求 tone 必须来自 config/statusTone.ts`);
  }
  if (result.hasRawListInput) {
    failures.push(`${relative} 在 AdminToolbar 内直接使用原生 input 或基础 Input，WEB-UI-009 要求搜索使用 AdminListSearch`);
  }
  if (result.missingFilterContainer) {
    failures.push(`${relative} 未在 AdminToolbar 内渲染 AdminListFilters，WEB-UI-010 要求使用共享筛选容器`);
  }
  if (result.hasInvalidFilterControl) {
    failures.push(`${relative} 在 AdminListFilters 外渲染筛选控件或直接使用基础 Select，WEB-UI-010 要求使用共享筛选边界`);
  }
}

let formDialogCount = 0;
for (const file of pages) {
  const relative = toRelative(file);
  if (relative === 'src/features/admin/components/AdminFormDialogShell/AdminFormDialogShell.tsx') continue;
  const result = analyzeAdminDialog(readFileSync(file, 'utf8'), file);
  if (result.hasInvalidFormDialogException) {
    failures.push(`${relative} 使用了未登记的 data-admin-form-dialog-exception；例外必须是契约允许的静态值`);
  }
  if (!result.isFormDialog) continue;
  formDialogCount += 1;
  if (result.missingShellImport) failures.push(`${relative} 是表单弹窗但未使用 AdminFormDialogShell 或完整的 shadcn AlertDialog 确认结构`);
  if (result.shellNotRendered) failures.push(`${relative} 引入了 AdminFormDialogShell 但未实际渲染`);
  if (result.hasRawDialogContent) failures.push(`${relative} 直接渲染 Dialog.Content，WEB-UI-009 要求使用 AdminFormDialogShell`);
}
failures.push(...dialogCoverageFailures(pages.length, formDialogCount));

if (failures.length > 0) {
  for (const message of failures) process.stderr.write(`${message}\n`);
  process.exitCode = 1;
} else {
  process.stdout.write(`后台列表结构契约通过，覆盖 ${checkedPages} 个列表页\n`);
  process.stdout.write(`后台表单表面契约通过，按控件语义覆盖 ${formSurfaceFileCount} 个文件\n`);
  process.stdout.write(`后台表单弹窗契约通过，按渲染语义覆盖 ${formDialogCount} 个表单弹窗（扫描 ${pages.length} 个 TSX 文件）\n`);
}
