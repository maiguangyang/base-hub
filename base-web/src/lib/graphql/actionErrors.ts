import { getGraphQLErrorCode } from './errors';

const reasons: Record<string, string> = {
  PERMISSION_DENIED: '当前账号没有执行此操作的权限。',
  STORE_SCOPE_DENIED: '当前账号无权操作该门店。',
  AUTH_REQUIRED: '登录已过期，请重新登录。',
  SESSION_REVOKED: '当前登录已失效，请重新登录。',
  CREDENTIALS_CHANGED: '账号凭证已变更，请重新登录。',
  ORGANIZATION_SUSPENDED: '所属组织已停用，请联系总部。',
  VALIDATION_FAILED: '资料不符合要求，请检查必填项和输入格式。',
  CONFLICT: '数据或状态已发生变化，请刷新并核对后重试。',
  NOT_FOUND: '操作对象不存在或已被删除，请刷新列表。',
  ROLE_IN_USE: '角色仍被有效成员使用，请先解除关联。',
  PERMISSION_DELEGATION_DENIED: '所选权限超出当前账号的可授权范围。',
  SELF_MEMBERSHIP_CHANGE_DENIED: '不能停用、删除或降低当前登录账号的权限。',
  LAST_HQ_SUPER_ADMIN_REQUIRED: '必须保留至少一位有效总部超级管理员。',
  LAST_OWNER_REQUIRED: '必须保留至少一位有效加盟商负责人。',
  INVALID_CREDENTIALS: '账号或密码不正确，请核对后重试。',
  PHONE_NOT_FOUND: '该手机号尚未注册，请核对账号。',
  PASSWORD_INCORRECT: '密码不正确，请重新输入。',
  ACCOUNT_INACTIVE: '账号已停用，请联系管理员。',
  ACCOUNT_UNAVAILABLE: '账号不可用，请联系管理员。',
  ACCOUNT_LOCKED: '账号已锁定，请稍后重试或联系管理员。',
  PASSWORD_WEAK: '密码强度不足，请按页面要求设置。',
  PASSWORD_CONFIRM_MISMATCH: '两次输入的密码不一致。',
  TEMP_PASSWORD_CHANGE_REQUIRED: '请先修改临时密码。',
  WORKSPACE_FORBIDDEN: '当前账号无权进入该工作台。',
  MEMBERSHIP_INACTIVE: '当前成员关系已停用，请联系管理员。',
  STORE_NOT_ACTIVE: '门店尚未启用，请先完成审核。',
  INVITATION_PENDING: '成员邀请尚未接受，请先完成邀请确认。',
  INVITATION_EXPIRED: '成员邀请已过期，请重新邀请。',
  STOCKTAKE_RECONCILIATION_REQUIRED: '盘点数据需要核对，请补齐差异原因后重试。',
  CUSTOMER_POINTS_MISMATCH: '积分余额与流水不一致，请先核对并校正。',
  HQ_ALREADY_BOOTSTRAPPED: '总部已完成初始化，请直接登录。',
  OPENING_RECORD_NUMBER_CONFLICT: '开通记录编号已被使用，请核对原始记录。',
  RATE_LIMITED: '操作过于频繁，请稍后重试。',
  NETWORK_ERROR: '网络连接失败，请检查网络后重试。',
  INTERNAL_ERROR: '服务端处理异常，请稍后重试。',
  CONFIG_UNAVAILABLE: '服务端配置尚未就绪，请联系管理员。',
  RESULT_UNKNOWN: '请求结果尚未确认，请刷新并核对数据，避免重复操作。',
};

/** Display approved reasons, never raw server exceptions or credentials. */
export function actionErrorReason(cause: unknown): string {
  const code = getGraphQLErrorCode(cause) ?? directCode(cause);
  if (code && reasons[code]) return reasons[code];
  if (cause instanceof TypeError) return reasons.NETWORK_ERROR;
  return '操作未完成，请稍后重试；如涉及新增或发放，请先核对结果。';
}

function directCode(cause: unknown): string | undefined {
  if (!cause || typeof cause !== 'object') return undefined;
  if ('code' in cause && typeof cause.code === 'string') return cause.code;
  if (cause instanceof Error && reasons[cause.message]) return cause.message;
  return undefined;
}
