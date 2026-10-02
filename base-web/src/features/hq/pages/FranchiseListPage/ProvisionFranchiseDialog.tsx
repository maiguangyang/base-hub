import { AdminValidatedField } from '@/features/admin/components/AdminValidatedField';
import { CopyTemporaryPasswordButton } from '@/features/admin/components/CopyTemporaryPasswordButton';
import { AdminActionError, AdminFormDialogShell } from '@/features/admin/components/AdminFormDialogShell';
import { generateEntityCode } from '@/lib/codeGenerator';

export interface ProvisionValues { name: string; ownerPhone: string; ownerDisplayName: string; ownerEmail: string; code?: string }
export type ProvisionOutcome = { kind: 'temporary-password'; value: string } | { kind: 'invitation-pending' };

interface ProvisionFranchiseDialogProps {
  open: boolean;
  onOpenChange(open: boolean): void;
  values?: ProvisionValues;
  onValuesChange?(values: ProvisionValues): void;
  onSubmit(): void | Promise<void>;
  isSubmitting?: boolean;
  outcome?: ProvisionOutcome;
  error?: string;
}

export function provisionOutcome(result: { temporaryPassword?: string | null; invitationPending: boolean }): ProvisionOutcome {
  return result.temporaryPassword
    ? { kind: 'temporary-password', value: result.temporaryPassword }
    : { kind: 'invitation-pending' };
}

export function provisionFranchiseInput(values: ProvisionValues) {
  return {
    code: values.code?.trim() || generateEntityCode('ORG'),
    name: values.name,
    ownerPhone: values.ownerPhone,
    ownerDisplayName: values.ownerDisplayName,
    ownerEmail: values.ownerEmail || null,
  };
}

export function ProvisionFranchiseDialog(props: ProvisionFranchiseDialogProps) {
  const values = props.values ?? { name: '', ownerPhone: '', ownerDisplayName: '', ownerEmail: '' };
  const change = (field: keyof ProvisionValues, value: string) => props.onValuesChange?.({ ...values, [field]: value });
  return (
    <AdminFormDialogShell open={props.open} onOpenChange={props.onOpenChange} title="开通加盟商" description="创建组织并配置首位加盟商老板。" submitLabel="确认开通" cancelLabel={props.outcome ? '完成' : '关闭'} isSubmitting={props.isSubmitting} hideSubmit={Boolean(props.outcome)} onSubmit={props.onSubmit}>
      <AdminActionError message={props.error} />
      {props.outcome ? <ProvisionResult outcome={props.outcome} /> : (
        <div className="flex flex-col gap-5">
          <AdminValidatedField label="组织名称" placeholder="请输入组织名称" value={values.name} onChange={(value) => change('name', value)} maxLength={128} />
          <AdminValidatedField label="老板手机号" placeholder="请输入老板手机号" value={values.ownerPhone} onChange={(value) => change('ownerPhone', value)} kind="phone" />
          <AdminValidatedField label="老板姓名" placeholder="请输入老板姓名" value={values.ownerDisplayName} onChange={(value) => change('ownerDisplayName', value)} maxLength={64} />
          <AdminValidatedField label="老板邮箱（选填）" placeholder="请输入老板邮箱" value={values.ownerEmail} onChange={(value) => change('ownerEmail', value)} kind="email" required={false} />
        </div>
      )}
    </AdminFormDialogShell>
  );
}

function ProvisionResult({ outcome }: { outcome: ProvisionOutcome }) {
  if (outcome.kind === 'invitation-pending') return <p className="text-sm">该账号已存在，邀请已等待对方接受。</p>;
  return (
    <div className="flex flex-col gap-2 rounded-md border border-border p-3">
      <p className="text-sm font-medium">临时密码仅显示一次，请通过安全渠道交付。</p>
      <code className="block break-all rounded bg-muted p-2">{outcome.value}</code>
      <CopyTemporaryPasswordButton password={outcome.value} />
    </div>
  );
}
