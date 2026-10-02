import assert from 'node:assert/strict';
import { test } from 'node:test';
import { ESLint } from 'eslint';

// 使用真实 ESLint 配置验证占位文案规则，防止 AST 选择器变化导致静默失效。
const eslint = new ESLint();
const cases = [
  ['<SelectItem>{placeholder}</SelectItem>', true],
  ['<SelectItem disabled>{props.placeholder}</SelectItem>', true],
  ['<option>{placeholder}</option>', true],
  ['<button role="option">{placeholder}</button>', true],
  ['<button role="option" disabled>{props.placeholder}</button>', true],
  ['<SelectValue placeholder={placeholder} />', false],
  ['<SelectItem>{clearLabel}</SelectItem>', false],
  ['<SelectItem>{option.label}</SelectItem>', false],
];

for (const [jsx, blocked] of cases) {
  test(`占位文案选项契约：${jsx}`, async () => {
    const [result] = await eslint.lintText(`const view = ${jsx};`, { filePath: 'src/placeholder-contract.tsx' });
    const violations = result.messages.filter((message) => message.ruleId === 'no-restricted-syntax');
    assert.equal(violations.length, blocked ? 1 : 0);
  });
}
