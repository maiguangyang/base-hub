import { readFileSync, readdirSync, statSync } from 'node:fs';
import { join } from 'node:path';
import { fileURLToPath } from 'node:url';
import { describe, expect, it } from 'vitest';
import ts from 'typescript';

const featuresRoot = fileURLToPath(new URL('../src/features/', import.meta.url));
const operationColumnClass = 'admin-operation-column';
const globalStyles = readFileSync(fileURLToPath(new URL('../src/styles/global.css', import.meta.url)), 'utf8');

function collectFeatureFiles(dir) {
  return readdirSync(dir).flatMap((entry) => {
    const full = join(dir, entry);
    if (statSync(full).isDirectory()) return collectFeatureFiles(full);
    return entry.endsWith('.tsx') && !entry.endsWith('.test.tsx') ? [full] : [];
  });
}

function operationColumnCoverage(source) {
  const file = ts.createSourceFile('page.tsx', source, ts.ScriptTarget.Latest, true, ts.ScriptKind.TSX);
  let headers = 0;
  let styledHeaders = 0;

  const visit = (node) => {
    if (ts.isJsxElement(node)) {
      const tag = node.openingElement.tagName.getText(file);
      const label = node.children
        .filter(ts.isJsxText)
        .map((child) => child.text)
        .join('')
        .trim();
      if ((tag === 'th' || tag === 'TableHead') && label === '操作') headers += 1;
    }
    if (ts.isStringLiteral(node) && node.text === '操作' && ts.isArrayLiteralExpression(node.parent)) headers += 1;
    if (ts.isStringLiteral(node) && node.text.split(/\s+/).includes(operationColumnClass)) styledHeaders += 1;
    ts.forEachChild(node, visit);
  };
  visit(file);
  return { headers, styledHeaders };
}

describe('后台列表操作列契约', () => {
  it('每个操作列表头都标记统一的右侧固定样式', () => {
    const failures = collectFeatureFiles(featuresRoot).flatMap((file) => {
      const coverage = operationColumnCoverage(readFileSync(file, 'utf8'));
      if (coverage.headers === 0 || coverage.styledHeaders >= coverage.headers) return [];
      return [`${file}: 发现 ${coverage.headers} 个操作列，仅 ${coverage.styledHeaders} 个表头应用统一样式`];
    });

    expect(failures).toEqual([]);
  });

  it('短操作组保持横排，仅在超过 300px 时换行', () => {
    const cellRule = globalStyles.slice(
      globalStyles.indexOf('.admin-operation-column,'),
      globalStyles.indexOf('table:has(> thead .admin-operation-column) > thead'),
    );
    const actionGroupRule = globalStyles.slice(
      globalStyles.indexOf('table:has(> thead .admin-operation-column) > tbody > tr > :last-child > div'),
    );
    expect(cellRule).toContain('max-width: 300px');
    expect(cellRule).toContain('white-space: nowrap');
    expect(actionGroupRule).toContain('width: max-content');
    expect(actionGroupRule).toContain('max-width: 300px');
  });
});
