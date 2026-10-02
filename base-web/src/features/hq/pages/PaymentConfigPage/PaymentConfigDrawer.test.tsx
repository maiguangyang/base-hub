// @vitest-environment jsdom

import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import type { ScopeRef, ScopeView } from '@/features/hq/api/paymentConfig';
import { PaymentConfigDrawer } from './PaymentConfigDrawer';

const mocks = vi.hoisted(() => ({
  permissions: ['paymentConfig:read', 'paymentConfig:manage'],
  read: vi.fn<(ref: ScopeRef) => Promise<ScopeView>>(),
  close: vi.fn(),
}));

vi.mock('@/features/auth/store/authStore', () => ({ useAuthStore: (select: (state: unknown) => unknown) => select({ viewer: {
  currentWorkspace: { workspaceType: 'HEADQUARTERS' }, permissions: mocks.permissions,
} }) }));
vi.mock('@/features/hq/api/paymentConfig', async (importOriginal) => ({
  ...await importOriginal<typeof import('@/features/hq/api/paymentConfig')>(),
  paymentConfigAPI: { read: mocks.read },
}));

beforeEach(() => {
  mocks.permissions = ['paymentConfig:read', 'paymentConfig:manage'];
  mocks.read.mockReset().mockResolvedValue({ channels: [] });
  mocks.close.mockReset();
});
afterEach(cleanup);

it('加盟商操作打开右侧抽屉并读取该加盟商配置', async () => {
  render(<PaymentConfigDrawer target={{ scope: 'FRANCHISE', organizationId: 'org-1', name: '甲加盟商' }} onClose={mocks.close} />);
  expect(screen.getByRole('dialog').textContent).toContain('甲加盟商');
  expect(screen.getByRole('dialog').getAttribute('data-side')).toBe('right');
  await waitFor(() => expect(mocks.read).toHaveBeenCalledWith({ scope: 'FRANCHISE', organizationId: 'org-1' }));
  fireEvent.click(screen.getByRole('button', { name: 'Close' }));
  expect(mocks.close).not.toHaveBeenCalled();
  await waitFor(() => expect(mocks.close).toHaveBeenCalledTimes(1));
});

it('门店操作只读取所选门店，缺少支付读取权限时隐藏抽屉', async () => {
  const drawer = render(<PaymentConfigDrawer target={{ scope: 'STORE', storeId: 'store-2', name: '乙门店' }} onClose={mocks.close} />);
  await waitFor(() => expect(mocks.read).toHaveBeenCalledWith({ scope: 'STORE', storeId: 'store-2' }));
  fireEvent.click(screen.getByRole('button', { name: 'Close' }));
  expect(mocks.close).not.toHaveBeenCalled();
  await waitFor(() => expect(mocks.close).toHaveBeenCalledTimes(1));
  drawer.unmount();
  mocks.permissions = [];
  render(<PaymentConfigDrawer target={{ scope: 'STORE', storeId: 'store-2', name: '乙门店' }} onClose={mocks.close} />);
  expect(screen.queryByRole('dialog')).toBeNull();
});
