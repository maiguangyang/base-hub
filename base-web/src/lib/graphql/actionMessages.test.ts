import { Kind, type DocumentNode } from 'graphql';
import { expect, it } from 'vitest';
import * as generated from '@/__generated__/graphql';
import { actionLabel } from './actionMessages';

it('所有现有后台写操作均有中文反馈名称，新增操作需补齐文案', () => {
  let mutations = 0;
  for (const value of Object.values(generated)) {
    if (!value || typeof value !== 'object' || !('kind' in value) || value.kind !== Kind.DOCUMENT) continue;
    const document = value as DocumentNode;
    mutations += checkMutationLabels(document);
  }
  expect(mutations).toBeGreaterThan(30);
});

function checkMutationLabels(document: DocumentNode): number {
  let count = 0;
  for (const definition of document.definitions) {
    if (definition.kind !== Kind.OPERATION_DEFINITION || definition.operation !== 'mutation') continue;
    count += 1;
    for (const selection of definition.selectionSet.selections) {
      if (selection.kind === Kind.FIELD) expect(actionLabel(selection.name.value), selection.name.value).not.toBe('操作');
    }
  }
  return count;
}
