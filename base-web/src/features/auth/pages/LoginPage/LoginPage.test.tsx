import { renderToStaticMarkup } from 'react-dom/server';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { LoginPage } from './LoginPage';

const setViewer = vi.fn();

vi.mock('@apollo/client/react', () => ({
  useMutation: () => [vi.fn(), { loading: false }],
}));

vi.mock('react-router', () => ({
  useNavigate: () => vi.fn(),
}));

vi.mock('@/__generated__', async (importOriginal) => {
  const actual = await importOriginal<typeof import('@/__generated__')>();
  return { ...actual, useFragment: vi.fn() };
});

vi.mock('@/features/auth/store/authStore', () => ({
  useAuthStore: (selector: (state: { setViewer: typeof setViewer }) => unknown) => selector({ setViewer }),
}));

describe('登录页视觉层级', () => {
  beforeEach(() => vi.clearAllMocks());

  it('只保留登录标题，并使用更舒展的字段和按钮尺寸', () => {
    const markup = renderToStaticMarkup(<LoginPage />);

    expect(markup).not.toContain('Base Hub');
    expect(markup).not.toContain('总部、直营店与加盟商共用同一账号入口。');
    expect(markup.match(/class="flex flex-col gap-4 text-sm font-medium"/g)).toHaveLength(2);
    expect(markup.match(/h-11/g)).toHaveLength(3);
  });

  it('关闭表单及字段的历史输入自动完成', () => {
    const markup = renderToStaticMarkup(<LoginPage />);

    expect(markup.match(/autoComplete="off"/g)).toHaveLength(3);
    expect(markup).not.toContain('autoComplete="username"');
    expect(markup).not.toContain('autoComplete="current-password"');
  });

  it('为手机号和密码输入框提供输入提示', () => {
    const markup = renderToStaticMarkup(<LoginPage />);

    expect(markup).toContain('placeholder="请输入手机号"');
    expect(markup).toContain('placeholder="请输入密码"');
    expect(markup).toContain('maxLength="11"');
    expect(markup).toContain('pattern="1[34578][0-9]{9}"');
  });
});
