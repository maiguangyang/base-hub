import { describe, expect, it } from 'vitest';
import { AdminFormDialogShell } from '@/features/admin/components/AdminFormDialogShell';
import { storeCreateInput } from '@/features/admin/pages/StoreEditorPage/storeInputs';
import { canDeleteDirectStore } from './DirectStoreListPage/DirectStoreTable';
import { ProvisionFranchiseDialog, provisionFranchiseInput, provisionOutcome } from './FranchiseListPage/ProvisionFranchiseDialog';
import { ReviewStoreDialog, availableReviewActions, reviewStoreInput } from './StoreApprovalListPage/ReviewStoreDialog';

describe('总部表单与输入', () => {
  it('加盟开通表单使用 AdminFormDialogShell 且不暴露编码字段', () => {
    const provision = ProvisionFranchiseDialog({ open: true, onOpenChange: () => undefined, onSubmit: async () => undefined });
    expect(provision.type).toBe(AdminFormDialogShell);
    expect(JSON.stringify(provision)).not.toContain('编码');
    expect(JSON.stringify(provision)).not.toContain('编号');
  });

  it('直营店输入由工作台注入组织和 ACTIVE，支持内部编码自动生成', () => {
    const input = storeCreateInput({ code: 'S1', name: 'One', contactPhone: '13800000000' }, 'hq', 'HEADQUARTERS');
    expect(input.code).toBe('S1');
    expect(input.name).toBe('One');
    expect(input.organizationId).toBe('hq');
    expect(input.lifecycle).toBe('ACTIVE');
    expect(input.contactPhone).toBe('13800000000');
    expect(input.businessStatus).toBe('OPEN');
    expect(input.supportDineIn).toBe(true);
    const auto = storeCreateInput({ name: 'AutoStore' }, 'hq', 'HEADQUARTERS');
    expect(auto.name).toBe('AutoStore');
    expect(auto.code.startsWith('STR')).toBe(true);
    expect(auto.organizationId).toBe('hq');
    expect(auto.lifecycle).toBe('ACTIVE');
  });

  it('加盟开通输入自动生成组织编码并规范化可选邮箱', () => {
    const input = provisionFranchiseInput({ name: 'Org1', ownerPhone: '13800000000', ownerDisplayName: 'Owner', ownerEmail: '' });
    expect(input.name).toBe('Org1');
    expect(input.code.startsWith('ORG')).toBe(true);
    expect(input.ownerPhone).toBe('13800000000');
    expect(input.ownerEmail).toBeNull();
  });

  it('一次性密码仅在新账号结果中展示，既有账号只显示待邀请', () => {
    expect(provisionOutcome({ temporaryPassword: 'Secret-Once-42!', invitationPending: false })).toEqual({ kind: 'temporary-password', value: 'Secret-Once-42!' });
    expect(provisionOutcome({ temporaryPassword: null, invitationPending: true })).toEqual({ kind: 'invitation-pending' });
  });

  it('一次性密码结果页不再暴露可重复提交按钮', () => {
    const dialog = ProvisionFranchiseDialog({
      open: true, onOpenChange: () => undefined, onSubmit: async () => undefined,
      outcome: { kind: 'temporary-password', value: 'Secret-Once-42!' },
    });
    expect(dialog.props.hideSubmit).toBe(true);
  });
});

describe('总部权限与审核', () => {
  it('直营店删除按钮只服从独立 delete 权限', () => {
    expect(canDeleteDirectStore(['hqStore:delete'])).toBe(true);
    expect(canDeleteDirectStore(['hqStore:update'])).toBe(false);
  });

  it('审核权限独立，退回必须填写原因而批准不发送原因', () => {
    expect(availableReviewActions(['store:approve'])).toEqual({ approve: true, reject: false });
    expect(availableReviewActions(['store:reject'])).toEqual({ approve: false, reject: true });
    expect(reviewStoreInput('store-1', true, 'ignored')).toEqual({ storeId: 'store-1', approved: true });
    expect(() => reviewStoreInput('store-1', false, '  ')).toThrow('VALIDATION_FAILED');
    expect(reviewStoreInput('store-1', false, 'missing license')).toEqual({ storeId: 'store-1', approved: false, rejectionReason: 'missing license' });
  });

  it('审核对话框同样使用统一表单壳', () => {
    const element = ReviewStoreDialog({ open: true, onOpenChange: () => undefined, storeId: 'store-1', approved: true, onSubmit: async () => undefined });
    expect(element.type).toBe(AdminFormDialogShell);
  });
});
