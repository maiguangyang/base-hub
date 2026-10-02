import { renderToStaticMarkup } from 'react-dom/server';
import { describe, expect, it, vi } from 'vitest';
import { PasswordChangePage } from './PasswordChangePage';

vi.mock('@apollo/client/react', () => ({
  useMutation: () => [vi.fn(), { loading: false }],
}));

vi.mock('@/features/auth/store/authStore', () => ({
  useAuthStore: (selector: (state: unknown) => unknown) => selector({
    viewer: { account: { mustChangePassword: true } },
    setViewer: vi.fn(),
  }),
}));

describe('首次登录改密页', () => {
  it('只呈现临时密码强制修改文案', () => {
    const markup = renderToStaticMarkup(<PasswordChangePage />);

    expect(markup).toContain('首次登录');
    expect(markup).toContain('当前临时密码');
  });
});
