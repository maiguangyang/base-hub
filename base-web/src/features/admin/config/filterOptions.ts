import type { AdminFilterOption } from '../components/AdminFilterSelect';
import type { StoreLifecycle } from '@/__generated__/graphql';

export const storeLifecycleOptions = [
  { value: 'DRAFT', label: '草稿' },
  { value: 'PENDING_APPROVAL', label: '待审核' },
  { value: 'ACTIVE', label: '已启用' },
  { value: 'REJECTED', label: '已退回' },
] as const satisfies readonly AdminFilterOption[];

export function storeLifecycleLabel(lifecycle: StoreLifecycle): string {
  return storeLifecycleOptions.find((option) => option.value === lifecycle)?.label ?? '未知状态';
}

export const storeBusinessStatusOptions = [
  { value: 'OPEN', label: '营业中' },
  { value: 'CLOSED', label: '已停业' },
] as const satisfies readonly AdminFilterOption[];

export const organizationStatusOptions = [
  { value: 'ACTIVE', label: '正常' },
  { value: 'SUSPENDED', label: '已暂停' },
] as const satisfies readonly AdminFilterOption[];

export const membershipStatusOptions = [
  { value: 'INVITED', label: '待接受' },
  { value: 'ACTIVE', label: '正常' },
  { value: 'SUSPENDED', label: '已停用' },
  { value: 'LEFT', label: '已离开' },
] as const satisfies readonly AdminFilterOption[];

export const customerMemberStatusOptions = [
  { value: 'ACTIVE', label: '正常' },
  { value: 'SUSPENDED', label: '已暂停' },
  { value: 'CANCEL_PENDING', label: '待注销' },
  { value: 'CANCELLED', label: '已注销' },
] as const satisfies readonly AdminFilterOption[];

export const hqRoleKindOptions = [
  { value: 'HQ_SUPER_ADMIN', label: '系统角色' },
  { value: 'CUSTOM', label: '自定义角色' },
] as const satisfies readonly AdminFilterOption[];

export const franchiseRoleKindOptions = [
  { value: 'FRANCHISE_OWNER', label: '所有者角色' },
  { value: 'CUSTOM', label: '自定义角色' },
] as const satisfies readonly AdminFilterOption[];

export const storeAccessModeOptions = [
  { value: 'ALL_STORES', label: '全部门店' },
  { value: 'SELECTED_STORES', label: '指定门店' },
] as const satisfies readonly AdminFilterOption[];

export const auditResultOptions = [
  { value: 'SUCCESS', label: '成功' },
  { value: 'FAILURE', label: '失败' },
] as const satisfies readonly AdminFilterOption[];

export const auditResourceTypeOptions = [
  { value: 'account', label: '账号' },
  { value: 'organization', label: '组织' },
  { value: 'store', label: '门店' },
  { value: 'operatorMembership', label: '成员' },
  { value: 'operatorRole', label: '角色' },
  { value: 'membershipInvitation', label: '成员邀请' },
] as const satisfies readonly AdminFilterOption[];
