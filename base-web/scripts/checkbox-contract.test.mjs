import assert from 'node:assert/strict';
import { test } from 'node:test';
import { ESLint } from 'eslint';

const eslint = new ESLint();
for (const [jsx, blocked] of [
  ['<input type="checkbox" />', true],
  ['<input type="checkbox" checked={false} />', true],
  ['<Checkbox checked={false} />', false],
  ['<Checkbox checked="indeterminate" />', false],
  ['<input type="file" />', false],
]) {
  test(`复选框统一组件契约：${jsx}`, async () => {
    const [result] = await eslint.lintText(`const view = ${jsx};`, { filePath: 'src/checkbox-contract.tsx' });
    const violations = result.messages.filter((message) => message.ruleId === 'no-restricted-syntax');
    assert.equal(violations.length, blocked ? 1 : 0);
  });
}

for (const moduleName of ['radix-ui', '@radix-ui/react-checkbox']) {
  for (const filePath of ['src/features/checkbox-contract.tsx', 'src/features/admin/checkbox-contract.tsx']) {
  test(`业务页面不得绕过统一复选框：${filePath} ${moduleName}`, async () => {
    const [result] = await eslint.lintText(`import { Checkbox } from '${moduleName}';`, {
      filePath,
    });
    assert.equal(result.messages.filter((message) => message.ruleId === 'no-restricted-imports').length, 1);
  });
}
}
