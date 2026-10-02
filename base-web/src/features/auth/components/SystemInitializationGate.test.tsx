// @vitest-environment jsdom

import { renderToStaticMarkup } from 'react-dom/server';
import { act, cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter, useLocation } from 'react-router';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { getSystemInitializationStatus } from '@/features/auth/api/systemInitialization';
import {
  initializationRedirect, InitializationStatusError, SystemInitializationGate,
} from './SystemInitializationGate';

vi.mock('@/features/auth/api/systemInitialization', () => ({
  getSystemInitializationStatus: vi.fn(),
}));

const getStatus = vi.mocked(getSystemInitializationStatus);

afterEach(() => {
  cleanup();
  getStatus.mockReset();
});

function LocationProbe() {
  const location = useLocation();
  return <output data-testid="location">{`${location.pathname}${location.search}${location.hash}`}</output>;
}

function renderGate(initialEntry: string) {
  return render(
    <MemoryRouter initialEntries={[initialEntry]}>
      <LocationProbe />
      <SystemInitializationGate><div data-testid="protected-content">受保护内容</div></SystemInitializationGate>
    </MemoryRouter>,
  );
}

describe('系统初始化门禁跳转规则', () => {
  it('未初始化时把完整业务地址放进初始化页 dt', () => {
    expect(initializationRedirect(false, {
      pathname: '/admin/hq/franchises', search: '?page=2&q=%E9%9F%A9%E6%96%87', hash: '#row-1',
    })).toBe('/admin/initialize?dt=%2Fadmin%2Fhq%2Ffranchises%3Fpage%3D2%26q%3D%25E9%259F%25A9%25E6%2596%2587%23row-1');
  });

  it('从带 dt 的登录页进入初始化时保留原业务地址', () => {
    expect(initializationRedirect(false, {
      pathname: '/admin/login', search: '?dt=%2Fadmin%2Ffranchise%2Fstores%3Fpage%3D2', hash: '',
    })).toBe('/admin/initialize?dt=%2Fadmin%2Ffranchise%2Fstores%3Fpage%3D2');
  });

  it('未初始化时留在初始化页，已初始化时携带 dt 转入登录页', () => {
    const location = {
      pathname: '/admin/initialize', search: '?dt=%2Fadmin%2Fhq%2Faudit%3Fpage%3D3', hash: '',
    };
    expect(initializationRedirect(false, location)).toBeUndefined();
    expect(initializationRedirect(true, location))
      .toBe('/admin/login?dt=%2Fadmin%2Fhq%2Faudit%3Fpage%3D3');
  });

  it('状态读取失败只提供重试，不渲染受保护内容', () => {
    const retry = vi.fn();
    const markup = renderToStaticMarkup(<InitializationStatusError onRetry={retry} />);

    expect(markup).toContain('无法检查系统初始化状态。');
    expect(markup).toContain('重试');
    expect(markup).not.toContain('data-protected-content');
  });
});

describe('系统初始化门禁生命周期', () => {
  it('状态确认前不挂载受保护内容，确认后才放行', async () => {
    let resolveStatus: ((initialized: boolean) => void) | undefined;
    getStatus.mockReturnValue(new Promise((resolve) => { resolveStatus = resolve; }));
    renderGate('/admin/hq/audit?page=2');

    expect(screen.queryByTestId('protected-content')).toBeNull();
    await act(async () => resolveStatus?.(true));
    expect(await screen.findByTestId('protected-content')).toBeTruthy();
  });

  it('未初始化时把完整地址带到初始化页', async () => {
    getStatus.mockResolvedValue(false);
    renderGate('/admin/hq/franchises?page=2#row-1');

    await waitFor(() => expect(screen.getByTestId('location').textContent)
      .toBe('/admin/initialize?dt=%2Fadmin%2Fhq%2Ffranchises%3Fpage%3D2%23row-1'));
    expect(screen.queryByTestId('protected-content')).toBeNull();
  });

  it('已初始化时把初始化页携带的 dt 转入登录页', async () => {
    getStatus.mockResolvedValue(true);
    renderGate('/admin/initialize?dt=%2Fadmin%2Ffranchise%2Fstores%3Fpage%3D3');

    await waitFor(() => expect(screen.getByTestId('location').textContent)
      .toBe('/admin/login?dt=%2Fadmin%2Ffranchise%2Fstores%3Fpage%3D3'));
  });

  it('状态失败后保持关闭，重试成功才放行', async () => {
    getStatus.mockRejectedValueOnce(new Error('offline')).mockResolvedValueOnce(true);
    renderGate('/admin/login');

    fireEvent.click(await screen.findByRole('button', { name: '重试' }));
    expect(await screen.findByTestId('protected-content')).toBeTruthy();
    expect(getStatus).toHaveBeenCalledTimes(2);
  });

  it('卸载时中止仍在进行的状态请求', () => {
    let requestSignal: AbortSignal | undefined;
    getStatus.mockImplementation((signal) => {
      requestSignal = signal;
      return new Promise(() => undefined);
    });
    const view = renderGate('/admin/login');

    view.unmount();
    expect(requestSignal?.aborted).toBe(true);
  });
});
