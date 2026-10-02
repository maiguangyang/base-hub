import { readFileSync } from 'node:fs';
import ts from 'typescript';

const sourcePath = new URL('./translate.ts', import.meta.url);
const baselinePath = new URL('./translate-baseline.json', import.meta.url);
const source = readFileSync(sourcePath, 'utf8');
const file = ts.createSourceFile('translate.ts', source, ts.ScriptTarget.Latest, true);

function activeLines() {
  const lines = new Set();
  const scanner = ts.createScanner(ts.ScriptTarget.Latest, false, ts.LanguageVariant.Standard, source);
  for (let kind = scanner.scan(); kind !== ts.SyntaxKind.EndOfFileToken; kind = scanner.scan()) {
    if (kind >= ts.SyntaxKind.FirstTriviaToken && kind <= ts.SyntaxKind.LastTriviaToken) continue;
    const start = file.getLineAndCharacterOfPosition(scanner.getTokenPos()).line;
    const end = file.getLineAndCharacterOfPosition(scanner.getTextPos() - 1).line;
    for (let line = start; line <= end; line++) lines.add(line);
  }
  return lines;
}

const lines = activeLines();

function countLines(node) {
  const start = file.getLineAndCharacterOfPosition(node.getStart(file)).line;
  const end = file.getLineAndCharacterOfPosition(node.end - 1).line;
  return [...lines].filter((line) => line >= start && line <= end).length;
}

function flowMetrics(body) {
  const result = { complexity: 1, depth: 0 };
  function visit(node, depth) {
    const control = ts.isIfStatement(node) || ts.isForStatement(node) || ts.isForOfStatement(node) ||
      ts.isForInStatement(node) || ts.isWhileStatement(node) || ts.isDoStatement(node) || ts.isSwitchStatement(node);
    const decision = control || ts.isCaseClause(node) || ts.isCatchClause(node) || ts.isConditionalExpression(node) ||
      (ts.isBinaryExpression(node) && [ts.SyntaxKind.AmpersandAmpersandToken, ts.SyntaxKind.BarBarToken].includes(node.operatorToken.kind));
    if (decision) result.complexity++;
    const nextDepth = depth + Number(control);
    result.depth = Math.max(result.depth, nextDepth);
    ts.forEachChild(node, (child) => visit(child, nextDepth));
  }
  visit(body, 0);
  return result;
}

function measurements() {
  const result = { 'translate.ts': { lines: lines.size, complexity: 0, depth: 0 } };
  function visit(node) {
    if (ts.isFunctionDeclaration(node) && node.name && node.body) {
      result[`translate.ts#${node.name.text}`] = { lines: countLines(node), ...flowMetrics(node.body) };
    }
    ts.forEachChild(node, visit);
  }
  visit(file);
  return result;
}

const actual = measurements();
if (process.argv.includes('--report')) {
  process.stdout.write(`${JSON.stringify(actual, null, 2)}\n`);
} else {
  const baseline = JSON.parse(readFileSync(baselinePath, 'utf8'));
  const failures = Object.entries(actual).filter(([key, value]) => {
    const allowed = baseline[key];
    return !allowed || Object.keys(value).some((metric) => value[metric] > allowed[metric]);
  });
  if (failures.length) {
    for (const [key, value] of failures) process.stderr.write(`${key} grew beyond baseline: ${JSON.stringify(value)}\n`);
    process.exitCode = 1;
  }
}
