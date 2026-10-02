import { messagesFor } from '@/i18n';

const messages = messagesFor('zh-CN').auth;

/** 把稳定错误码映射为界面文案；未知码使用统一安全兜底。 */
export function authErrorMessage(code: string | undefined): string {
  if (!code) return messages.errors.UNKNOWN;
  return messages.errors[code as keyof typeof messages.errors] ?? messages.errors.UNKNOWN;
}

export { messages as authMessages };
