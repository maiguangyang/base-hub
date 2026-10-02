import { useRef, useState } from 'react';
import { useMutation } from '@apollo/client/react';
import { runAdminAction } from '@/features/admin/components/AdminFormDialogShell';
import { getGraphQLErrorCode } from '@/lib/graphql/errors';
import { PROVISION_FRANCHISE_MUTATION, RESET_FRANCHISE_INITIAL_PASSWORD_MUTATION, RESTORE_ORGANIZATION_MUTATION, SUSPEND_ORGANIZATION_MUTATION } from '../../graphql/franchises';
import type { FranchiseRow, VisibleTemporaryPasswords } from './FranchiseTable';
import { provisionFranchiseInput, provisionOutcome, type ProvisionOutcome, type ProvisionValues } from './ProvisionFranchiseDialog';

const emptyValues: ProvisionValues = { name: '', ownerPhone: '', ownerDisplayName: '', ownerEmail: '' };

export function useFranchiseListActions(refetch: () => Promise<unknown>) {
  const provisionActions = useFranchiseProvision(refetch);
  const [actionError, setActionError] = useState<string>();
  const [temporaryPasswords, setTemporaryPasswords] = useState<VisibleTemporaryPasswords>({});
  const [resettingOrganizationId, setResettingOrganizationId] = useState<string>();
  const [changingStatusOrganizationId, setChangingStatusOrganizationId] = useState<string>();
  const resetting = useRef(false);
  const changingStatus = useRef(false);
  const [suspend] = useMutation(SUSPEND_ORGANIZATION_MUTATION);
  const [restore] = useMutation(RESTORE_ORGANIZATION_MUTATION);
  const [resetInitialPassword] = useMutation(RESET_FRANCHISE_INITIAL_PASSWORD_MUTATION, { fetchPolicy: 'no-cache' });

  /** 接收对话框已确认的原因；成功返回 true 供界面关闭弹窗。 */
  async function suspendRow(row: FranchiseRow, reasonCode: string): Promise<boolean> {
    if (changingStatus.current) return false;
    return changeStatus(row.id, async () => { await suspend({ variables: { input: { organizationId: row.id, reasonCode } } }); await refetch(); });
  }
  async function restoreRow(row: FranchiseRow) {
    if (changingStatus.current) return;
    await changeStatus(row.id, async () => { await restore({ variables: { id: row.id } }); await refetch(); });
  }
  /** 状态请求与列表刷新期间阻止重复提交，并保留服务端最终状态。 */
  async function changeStatus(organizationId: string, action: () => Promise<void>): Promise<boolean> {
    changingStatus.current = true;
    setChangingStatusOrganizationId(organizationId);
    try { return await runAdminAction(action, setActionError); }
    finally { changingStatus.current = false; setChangingStatusOrganizationId(undefined); }
  }
  async function resetPasswordRow(row: FranchiseRow) {
    const accountId = row.initialAccountId;
    if (resetting.current || !accountId || !row.initialAccount) return;
    resetting.current = true;
    setResettingOrganizationId(row.id);
    setActionError(undefined);
    setTemporaryPasswords((current) => clearAccountPassword(current, accountId));
    try {
      const result = await resetInitialPassword({ variables: { organizationId: row.id } });
      const value = result.data?.resetFranchiseInitialPassword;
      if (!value?.temporaryPassword) throw new Error('EMPTY_TEMPORARY_PASSWORD');
      setTemporaryPasswords((current) => replaceAccountPassword(current, row.id, value.accountId, value.temporaryPassword));
    } catch (cause) {
      setActionError(resetPasswordErrorMessage(getGraphQLErrorCode(cause)));
    } finally {
      resetting.current = false;
      setResettingOrganizationId(undefined);
    }
  }
  return { ...provisionActions, actionError, temporaryPasswords, resettingOrganizationId, changingStatusOrganizationId, suspendRow, restoreRow, resetPasswordRow };
}

function useFranchiseProvision(refetch: () => Promise<unknown>) {
  const [open, setOpen] = useState(false);
  const [values, setValues] = useState(emptyValues);
  const [outcome, setOutcome] = useState<ProvisionOutcome>();
  const [provisionError, setProvisionError] = useState<string>();
  const [provision, provisionState] = useMutation(PROVISION_FRANCHISE_MUTATION);
  function changeOpen(next: boolean) {
    setOpen(next);
    if (!next) { setOutcome(undefined); setValues(emptyValues); setProvisionError(undefined); }
  }
  function changeValues(next: ProvisionValues) { setValues(next); setProvisionError(undefined); }
  async function submit() {
    setProvisionError(undefined);
    try {
      const result = await provision({ variables: { input: provisionFranchiseInput(values) } });
      if (!result.data?.provisionFranchise) throw new Error('EMPTY_PROVISION_RESULT');
      setOutcome(provisionOutcome(result.data.provisionFranchise));
    } catch (cause) {
      setProvisionError(provisionErrorMessage(getGraphQLErrorCode(cause)));
      return;
    }
    try { await refetch(); }
    catch { setProvisionError('加盟商已开通，但列表刷新失败。请刷新页面后查看。'); }
  }
  return { open, values, outcome, provisionError, provisionState, setOpen, changeValues, changeOpen, submit };
}

function provisionErrorMessage(code?: string): string {
  if (code === 'PERMISSION_DENIED') return '该手机号对应的账号无法用于开通加盟商，请确认账号未停用且不属于总部。';
  if (code === 'VALIDATION_FAILED') return '开通信息无效或组织编码重复，请核对后重试。';
  return '开通失败，请稍后重试。';
}

export function replaceAccountPassword(current: VisibleTemporaryPasswords, organizationId: string, accountId: string, password: string): VisibleTemporaryPasswords {
  const next = clearAccountPassword(current, accountId);
  delete next[organizationId];
  next[organizationId] = { accountId, password };
  return next;
}

function clearAccountPassword(current: VisibleTemporaryPasswords, accountId: string): VisibleTemporaryPasswords {
  return Object.fromEntries(Object.entries(current).filter(([, item]) => item.accountId !== accountId));
}

export function resetPasswordErrorMessage(code?: string): string {
  if (code === 'MEMBERSHIP_INACTIVE') return '初始账号尚未激活或已停用，暂时无法重置。';
  if (code === 'PERMISSION_DENIED') return '初始账号不可重置，请检查账号状态与操作权限。';
  if (code === 'CONFLICT') return '初始账号记录与当前绑定不一致，请暂停重置并联系管理员核查。';
  return '重置结果暂未确认，旧临时密码可能已失效。如未收到新密码，请再次重置。';
}
