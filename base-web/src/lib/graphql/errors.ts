/** GraphQL 服务端通过 extensions.code 暴露的稳定错误结构。 */
interface GraphQLErrorShape {
  extensions?: Record<string, unknown>;
}

const terminalSessionCodes = new Set([
  'AUTH_REQUIRED',
  'SESSION_REVOKED',
  'ORGANIZATION_SUSPENDED',
  'CREDENTIALS_CHANGED',
]);

/** 只读取结构化错误码，禁止按 message 文本匹配业务行为。 */
export function getGraphQLErrorCode(error: unknown): string | undefined {
  if (!isRecord(error)) return undefined;
  const direct = codeFromGraphQLError(error);
  if (direct) return direct;
  const errors = Array.isArray(error.errors) ? error.errors : [];
  for (const item of errors) {
    const code = codeFromGraphQLError(item);
    if (code) return code;
  }
  return undefined;
}

/** 判断错误码是否意味着当前浏览器会话必须立即清理。 */
export function isTerminalSessionCode(code: string | undefined): boolean {
  return typeof code === 'string' && terminalSessionCodes.has(code);
}

function codeFromGraphQLError(error: unknown): string | undefined {
  if (!isRecord(error) || !isRecord(error.extensions)) return undefined;
  return typeof error.extensions.code === 'string' ? error.extensions.code : undefined;
}

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === 'object' && value !== null;
}

export type { GraphQLErrorShape };
