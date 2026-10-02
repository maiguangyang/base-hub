// @vitest-environment jsdom

import { cleanup, fireEvent, render, screen } from '@testing-library/react';
import { afterEach, expect, it } from 'vitest';
import { AdminToastProvider } from '@/features/admin/components/AdminToast';
import { AiSecretCard } from './AiSecretCard';

afterEach(() => cleanup());

it('临时密码默认遮盖，用户主动显示后才呈现明文', () => {
  render(<AdminToastProvider><AiSecretCard secret={{ type: 'secret', toolId: 'FranchiseInitialAccount', targetAccountId: 'account-1', value: 'private-secret-123' }} /></AdminToastProvider>);
  expect(screen.queryByText('private-secret-123')).toBeNull();
  expect(screen.getByRole('button', { name: '复制临时密码' })).toBeTruthy();
  fireEvent.click(screen.getByRole('button', { name: '显示临时密码' }));
  expect(screen.getByText('private-secret-123')).toBeTruthy();
  fireEvent.click(screen.getByRole('button', { name: '隐藏临时密码' }));
  expect(screen.queryByText('private-secret-123')).toBeNull();
});
