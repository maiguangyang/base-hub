import { renderToStaticMarkup } from 'react-dom/server';
import { describe, expect, it, vi } from 'vitest';
import { SystemInitializationError } from '../../api/systemInitialization';
import { InitializationPage, submitSystemInitialization } from './InitializationPage';

const validInput = {
  phone: '13800000000',
  password: 'Correct-Horse-42',
  passwordConfirmation: 'Correct-Horse-42',
};

describe('系统初始化页', () => {
  it('展示三项输入、16px 字段间距、密码策略和关闭自动完成', () => {
    const markup = renderToStaticMarkup(<InitializationPage />);

    expect(markup).toContain('初始化系统');
    expect(markup).toContain('placeholder="请输入手机号"');
    expect(markup).toContain('maxLength="11"');
    expect(markup).toContain('pattern="1[34578][0-9]{9}"');
    expect(markup).toContain('placeholder="请输入密码"');
    expect(markup).toContain('placeholder="请再次输入密码"');
    expect(markup.match(/class="flex flex-col gap-4 text-sm font-medium"/g)).toHaveLength(3);
    expect(markup.match(/autoComplete="off"/g)).toHaveLength(4);
    expect(markup).toContain('密码需为 8–20 个字符。');
    expect(markup).not.toContain('大写字母');
    expect(markup.match(/h-11/g)).toHaveLength(4);
  });

  it.each([
    ['PASSWORD_WEAK', { ...validInput, password: 'short', passwordConfirmation: 'short' }],
    ['PASSWORD_CONFIRM_MISMATCH', { ...validInput, passwordConfirmation: 'Different-Horse-42' }],
  ])('客户端校验 %s 时不提交请求', async (code, input) => {
    const initialize = vi.fn();
    const navigate = vi.fn();

    await expect(submitSystemInitialization(input, '/admin/hq?a=1#x', { initialize, navigate }))
      .resolves.toBe(code);
    expect(initialize).not.toHaveBeenCalled();
    expect(navigate).not.toHaveBeenCalled();
  });

  it('成功后带着安全目的地址进入登录页', async () => {
    const initialize = vi.fn().mockResolvedValue(undefined);
    const navigate = vi.fn();

    await expect(submitSystemInitialization(validInput, '/admin/hq?a=1#x', { initialize, navigate }))
      .resolves.toBeUndefined();
    expect(initialize).toHaveBeenCalledWith(validInput);
    expect(navigate).toHaveBeenCalledWith('/admin/login?dt=%2Fadmin%2Fhq%3Fa%3D1%23x');
  });

  it('并发初始化冲突时转入登录页，其他错误返回稳定错误码', async () => {
    const navigate = vi.fn();
    const conflict = vi.fn().mockRejectedValue(new SystemInitializationError('HQ_ALREADY_BOOTSTRAPPED', 409));
    const failed = vi.fn().mockRejectedValue(new Error('network details'));

    await expect(submitSystemInitialization(validInput, '/admin/franchise', { initialize: conflict, navigate }))
      .resolves.toBeUndefined();
    expect(navigate).toHaveBeenCalledWith('/admin/login?dt=%2Fadmin%2Ffranchise');
    await expect(submitSystemInitialization(validInput, undefined, { initialize: failed, navigate }))
      .resolves.toBe('INTERNAL_ERROR');
  });
});
