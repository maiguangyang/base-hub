import { appConfig } from '@/config/app';

export interface SystemInitializationInput {
  phone: string;
  password: string;
  passwordConfirmation: string;
}

/** 初始化 REST 边界的稳定客户端错误。 */
export class SystemInitializationError extends Error {
  constructor(public readonly code: string, public readonly status: number) {
    super(code);
    this.name = 'SystemInitializationError';
  }
}

/** 查询一次性总部初始化标记。 */
export async function getSystemInitializationStatus(signal?: AbortSignal): Promise<boolean> {
  const response = await fetch(appConfig.engine.systemInitializationUrl, {
    credentials: 'include', headers: { Accept: 'application/json' }, signal,
  });
  if (!response.ok) throw await initializationError(response);
  const body: unknown = await response.json();
  if (!isInitializedBody(body)) throw new SystemInitializationError('INTERNAL_ERROR', response.status);
  return body.initialized;
}

/** 提交首个总部管理员的手机号与正式密码。 */
export async function initializeSystem(input: SystemInitializationInput, signal?: AbortSignal): Promise<void> {
  const response = await fetch(appConfig.engine.systemInitializationUrl, {
    method: 'POST', credentials: 'include',
    headers: { Accept: 'application/json', 'Content-Type': 'application/json' },
    body: JSON.stringify(input), signal,
  });
  if (!response.ok) throw await initializationError(response);
  const body: unknown = await response.json();
  if (!isInitializedBody(body) || !body.initialized) throw new SystemInitializationError('INTERNAL_ERROR', response.status);
}

async function initializationError(response: Response): Promise<SystemInitializationError> {
  try {
    const body: unknown = await response.json();
    if (isErrorBody(body)) return new SystemInitializationError(body.code, response.status);
  } catch {
    // 损坏或非 JSON 的失败响应统一收敛，不向界面暴露正文。
  }
  return new SystemInitializationError('INTERNAL_ERROR', response.status);
}

function isInitializedBody(value: unknown): value is { initialized: boolean } {
  return typeof value === 'object' && value !== null && 'initialized' in value
    && typeof value.initialized === 'boolean';
}

function isErrorBody(value: unknown): value is { code: string } {
  return typeof value === 'object' && value !== null && 'code' in value
    && typeof value.code === 'string';
}
