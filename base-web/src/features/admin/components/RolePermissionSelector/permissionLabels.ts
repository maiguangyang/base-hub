export interface PermissionOption {
  id: string;
  name: string;
  action: string;
  module: string;
  scope: 'SYSTEM' | 'TENANT';
}

// CRUD 名称对应 base-engine/model/model.graphql 的 @entity(title)。
const entityTitles: Record<string, string> = {
  account: '账号', organization: '加盟商', operatorMembership: '加盟商成员', permission: '权限',
  operatorRole: '角色', store: '门店', session: '登录会话', membershipInvitation: '成员邀请',
  auditLog: '审计日志', franchiseOpeningRecord: '加盟商开通记录',
  globalPaymentConfig: '全局支付配置', franchisePaymentConfig: '加盟商支付配置',
  storePaymentConfig: '门店支付配置',
};

const domainTitles: Record<string, string> = {
  franchise: '加盟商开通', hqRole: '总部角色', hqMembership: '管理员',
  hqStore: '直营门店', aiModelConfig: 'AI 模型配置', paymentConfig: '支付配置',
  tenantAudit: '加盟商审计日志',
};

const actionLabels: Record<string, string> = {
  read: '查看', create: '创建', update: '编辑', delete: '删除',
  restore: '恢复', suspend: '停用', provision: '开通',
  read_all: '查看全部', approve: '批准', reject: '退回',
  read_sensitive: '查看敏感信息', export: '导出', manage: '管理',
  cancel: '取消', grant: '发放', reverse: '冲正', correct: '更正',
  revoke: '撤销', record: '登记', post: '过账', submit: '提交',
};

export function permissionGroupLabel(resource: string): string {
  return entityTitles[resource] ?? domainTitles[resource] ?? resource;
}

export function permissionDisplayName(permission: PermissionOption): string {
  if (permission.action === 'franchise:provision') return '开通加盟商';
  const [resource, action] = permission.action.split(':');
  const title = permissionGroupLabel(resource);
  const verb = actionLabels[action];
  if (title !== resource && verb) return `${verb}${title}`;
  return permission.name && !permission.name.startsWith('permission.') ? permission.name : permission.action;
}
