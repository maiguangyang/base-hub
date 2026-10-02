import { useRef, useState } from 'react';
import { useQuery } from '@apollo/client/react';
import { AdminActionError, AdminFormDialogShell } from '@/features/admin/components/AdminFormDialogShell';
import { AdminSearchSelect } from '@/features/admin/components/AdminSearchSelect';
import { AdminValidatedField } from '@/features/admin/components/AdminValidatedField';
import { Checkbox } from '@/components/ui/checkbox';
import { HQ_FRANCHISE_INITIAL_ACCOUNT_CANDIDATES_QUERY } from '../../graphql/franchises';
import type { FranchiseRow } from './FranchiseTable';
import { InitialAccountRequestError, setFranchiseInitialAccount } from './setFranchiseInitialAccount';
import { publishActionFeedback } from '@/lib/actionFeedback';

interface Props { row: FranchiseRow; onClose(): void; onConfirmed(): Promise<unknown> }

export function ConfirmInitialAccountDialog({ row, onClose, onConfirmed }: Props) {
  const [accountId, setAccountId] = useState('');
  const [evidenceReference, setEvidenceReference] = useState('');
  const [attested, setAttested] = useState(false);
  const { data, loading, error: loadError } = useQuery(HQ_FRANCHISE_INITIAL_ACCOUNT_CANDIDATES_QUERY, { variables: { id: row.id }, fetchPolicy: 'network-only' });
  const candidates = eligibleCandidates(data?.organization?.memberships);
  const { pending, confirmed, uncertain, error, setError, submit } = useInitialAccountSubmit({ row, onClose, onConfirmed, accountId, evidenceReference, attested, candidates });
  const view = initialAccountDialogView(row, loading, Boolean(loadError), accountId, evidenceReference, attested, confirmed || uncertain);
  return <AdminFormDialogShell open onOpenChange={(open) => { if (!open && !pending) onClose(); }} {...view} isSubmitting={pending} onSubmit={submit}>
    <InitialAccountFields candidates={candidates} accountId={accountId} evidenceReference={evidenceReference} attested={attested} error={error} loading={loading} loadError={Boolean(loadError)} pending={pending} confirmed={confirmed || uncertain} onAccountChange={(value) => { setAccountId(value); setError(undefined); }} onEvidenceChange={setEvidenceReference} onAttestedChange={setAttested} />
  </AdminFormDialogShell>;
}

function InitialAccountFields({ candidates, accountId, evidenceReference, attested, error, loading, loadError, pending, confirmed, onAccountChange, onEvidenceChange, onAttestedChange }: { candidates: Candidate[]; accountId: string; evidenceReference: string; attested: boolean; error?: string; loading: boolean; loadError: boolean; pending: boolean; confirmed: boolean; onAccountChange(value: string): void; onEvidenceChange(value: string): void; onAttestedChange(value: boolean): void }) {
  return <div className="flex flex-col gap-5">
      {loadError && <AdminActionError message="成员加载失败，请稍后重试。" />}
      <AdminActionError message={error} />
      <div className="flex flex-col gap-2 text-sm font-medium"><span>初始账号</span>
        <AdminSearchSelect label="初始账号" value={accountId || undefined} options={candidates.map((member) => ({ value: member.account!.id, label: `${member.account!.displayName} · ${member.account!.phone}` }))} placeholder="请选择并核对姓名、手机号" searchPlaceholder="搜索姓名或手机号" emptyMessage="没有匹配的成员" disabled={loading || loadError || pending || confirmed} onValueChange={(value) => onAccountChange(value ?? '')} />
      </div>
      <AdminValidatedField label="开通记录编号" placeholder="请输入已核对的开通记录编号" value={evidenceReference} onChange={onEvidenceChange} maxLength={128} />
      <label className="flex items-start gap-2 text-sm text-foreground"><Checkbox checked={attested} disabled={pending || confirmed} onCheckedChange={(checked) => onAttestedChange(checked === true)} className="mt-0.5" /><span>我已人工核对原始开通记录，确认所选账号与记录编号属于该加盟商。</span></label>
      {!loading && !loadError && candidates.length === 0 && <p className="text-sm text-muted-foreground">没有可确认的有效成员，请先核对成员状态。</p>}
      <p className="text-sm text-muted-foreground">系统无法独立核验历史开通凭证。保存后，密码重置将针对选定账号。</p>
    </div>;
}

type Candidate = { id: string; status: string; account?: { id: string; phone: string; displayName: string; status: string } | null };

function initialAccountDialogView(row: FranchiseRow, loading: boolean, loadError: boolean, accountId: string, evidenceReference: string, attested: boolean, confirmed: boolean) {
  return {
    title: '初始化账户',
    description: `请根据开通记录核对 ${row.name} 的原始账号，并填写该记录的编号。角色和成员创建顺序不能作为依据。`,
    submitLabel: '确认绑定',
    cancelLabel: confirmed ? '关闭' : '取消',
    submitDisabled: loading || loadError || !accountId || !evidenceReference.trim() || !attested || confirmed,
  };
}

function useInitialAccountSubmit(props: Props & { accountId: string; evidenceReference: string; attested: boolean; candidates: Candidate[] }) {
  const [error, setError] = useState<string>();
  const [pending, setPending] = useState(false);
  const [confirmed, setConfirmed] = useState(false);
  const [uncertain, setUncertain] = useState(false);
  const submitting = useRef(false);

  async function save(): Promise<boolean> {
    try {
      await setFranchiseInitialAccount({ organizationId: props.row.id, accountId: props.accountId, evidenceReference: props.evidenceReference.trim(), attested: props.attested });
      setConfirmed(true);
      publishActionFeedback('success', '初始化账户成功。');
      return true;
    } catch (cause) {
      if (cause instanceof InitialAccountRequestError && cause.code === 'RESULT_UNKNOWN') setUncertain(true);
      setError(confirmationErrorMessage(cause instanceof InitialAccountRequestError ? cause.code : undefined));
      publishActionFeedback('error', `初始化账户失败：${confirmationErrorMessage(cause instanceof InitialAccountRequestError ? cause.code : undefined)}`);
      return false;
    }
  }

  async function refresh(): Promise<void> {
    try {
      await props.onConfirmed();
      props.onClose();
    } catch {
      setError('初始账户已保存，但列表刷新失败。请刷新页面后查看，勿重复提交。');
    }
  }

  async function submit() {
    if (submitting.current || confirmed || uncertain) return;
    if (!props.accountId || !props.candidates.some((item) => item.account?.id === props.accountId) || !props.evidenceReference.trim() || !props.attested) {
      setError('请选择初始账号、填写开通记录编号，并确认已人工核对原始凭证。');
      return;
    }
    submitting.current = true;
    setPending(true);
    setError(undefined);
    if (await save()) await refresh();
    submitting.current = false;
    setPending(false);
  }

  return { pending, confirmed, uncertain, error, setError, submit };
}

function eligibleCandidates(memberships?: readonly Candidate[] | null): Candidate[] {
  return memberships?.filter((member) => member.status === 'ACTIVE' && member.account?.status === 'ACTIVE' && Boolean(member.account.id) && Boolean(member.account.phone)) ?? [];
}

function confirmationErrorMessage(code?: string): string {
	if (code === 'OPENING_RECORD_NUMBER_CONFLICT') return '开通记录编号已被使用，请核对并填写该加盟商的正确编号。';
  if (code === 'CONFLICT') return '初始账户已完成初始化，请刷新列表查看。';
  if (code === 'VALIDATION_FAILED') return '请填写有效的开通记录编号，并核对账号。';
  if (code === 'PERMISSION_DENIED') return '该成员不符合初始化条件，或当前账号没有操作权限。';
  if (code === 'AUTH_REQUIRED') return '登录已过期，请重新登录后操作。';
  if (code === 'RESULT_UNKNOWN') return '请求结果暂未确认，请刷新列表核对初始账户后再操作。';
  return '保存失败，请核对账号和成员状态后重试。';
}
