import { describe, expect, it } from 'vitest';
import { authErrorMessage } from './authMessages';

describe('登录失败提示', () => {
  it('分别提示不存在的手机号和错误的密码', () => {
    expect(authErrorMessage('PHONE_NOT_FOUND')).toBe('手机号不存在。');
    expect(authErrorMessage('PASSWORD_INCORRECT')).toBe('密码错误。');
  });

  it('其他凭据失败仍使用通用提示', () => {
    expect(authErrorMessage('INVALID_CREDENTIALS')).toBe('手机号或密码不正确。');
  });

  it('说明账号不可用的原因', () => {
    expect(authErrorMessage('ACCOUNT_INACTIVE')).toBe('账号已停用，请联系管理员。');
    expect(authErrorMessage('ACCOUNT_UNAVAILABLE')).toBe('账号暂不可登录，请联系管理员。');
  });
});
