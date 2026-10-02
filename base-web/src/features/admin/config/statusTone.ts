/** 后台状态语气。§4.4 规定状态颜色必须来自本映射，页面禁止手写颜色 class。 */
export type AdminStatusTone = 'success' | 'warning' | 'danger' | 'info' | 'neutral';

/** 每个语气的三层 class。取自 admin-theme.css 的作用域 token，
 *  这些工具类只在 [data-surface="admin"] 下有值，不得在后台之外使用。 */
export const statusToneClasses: Record<AdminStatusTone, string> = {
  success: 'bg-success-bg border-success-border text-success-fg',
  warning: 'bg-warning-bg border-warning-border text-warning-fg',
  danger: 'bg-danger-bg border-danger-border text-danger-fg',
  info: 'bg-info-bg border-info-border text-info-fg',
  neutral: 'bg-neutral-bg border-neutral-border text-neutral-fg',
};

/** 后端 state 字段：1 为启用，其余按禁用处理。 */
export function toneForEnabledState(state: number | null | undefined): AdminStatusTone {
  return state === 1 ? 'success' : 'danger';
}

/** 后端 isDelete 字段：1 为已删除，其余为正常。 */
export function toneForDeleteState(isDelete: number | null | undefined): AdminStatusTone {
  return isDelete === 1 ? 'danger' : 'neutral';
}

/** 任务完成态：已完成为 success，未完成按待处理计为 warning。 */
export function toneForCompletion(completed: boolean | null | undefined): AdminStatusTone {
  return completed === true ? 'success' : 'warning';
}

/** 组织状态语气。 */
export function toneForOrganizationStatus(status: string): AdminStatusTone {
  return status === 'ACTIVE' ? 'success' : 'danger';
}

/** 门店准入生命周期语气。 */
export function toneForStoreLifecycle(lifecycle: string): AdminStatusTone {
  if (lifecycle === 'ACTIVE') return 'success';
  if (lifecycle === 'PENDING_APPROVAL') return 'warning';
  if (lifecycle === 'REJECTED') return 'danger';
  return 'neutral';
}

/** 审计结果语气。 */
export function toneForAuditResult(result: string): AdminStatusTone {
  return result === 'SUCCESS' ? 'success' : 'danger';
}

/** 成员状态语气。 */
export function toneForMembershipStatus(status: string): AdminStatusTone {
  if (status === 'ACTIVE') return 'success';
  if (status === 'INVITED') return 'warning';
  return 'danger';
}

/** 门店营业状态语气。 */
export function toneForStoreBusinessStatus(status: string | null | undefined): AdminStatusTone {
  return status === 'OPEN' ? 'success' : 'neutral';
}
