interface RoleLabelSource { name: string; kind?: string }

/** 系统角色持久化名称是稳定键；界面在显示边界转换为本地文案。 */
export function roleDisplayName(role: RoleLabelSource): string {
  if (role.kind === 'FRANCHISE_OWNER' || role.name === 'role.franchiseOwner') return '加盟商负责人';
  if (role.kind === 'HQ_SUPER_ADMIN' || role.name === 'role.hqSuperAdministrator') return '总部超级管理员';
  return role.name;
}
