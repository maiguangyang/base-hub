import { expect, it } from 'vitest';
import { storeBusinessStatusInput } from '@/features/admin/pages/StoreEditorPage/storeInputs';

it('changes only the direct store business status', () => {
  expect(storeBusinessStatusInput(true)).toEqual({ businessStatus: 'OPEN' });
  expect(storeBusinessStatusInput(false)).toEqual({ businessStatus: 'CLOSED' });
});
