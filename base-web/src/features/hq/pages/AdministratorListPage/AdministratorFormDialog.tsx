import { AdminValidatedField } from '@/features/admin/components/AdminValidatedField';
import { CopyTemporaryPasswordButton } from '@/features/admin/components/CopyTemporaryPasswordButton';
import { AdminActionError, AdminFormDialogShell } from '@/features/admin/components/AdminFormDialogShell';
import { AdminRoleSelect } from '@/features/admin/components/AdminRoleSelect/AdminRoleSelect';

export interface AdministratorValues {
  phone: string;
  displayName: string;
  email: string;
  roleIds: string[];
}

export type AdministratorInviteOutcome = { kind: 'temporary-password'; value: string } | { kind: 'invitation-pending' };

interface AdministratorFormDialogProps {
  open: boolean;
  mode: 'create' | 'edit';
  onOpenChange(open: boolean): void;
  values?: AdministratorValues;
  onValuesChange?(values: AdministratorValues): void;
  onSubmit(): void | Promise<void>;
  isSubmitting?: boolean;
  roles: { id: string; name: string; kind: string }[];
  account?: { displayName: string; phone: string; email?: string | null };
  outcome?: AdministratorInviteOutcome;
  error?: string;
}

export function administratorInviteInput(values: AdministratorValues) {
  return {
    phone: values.phone, displayName: values.displayName, email: values.email || null,
    roleIds: values.roleIds, storeAccessMode: 'ALL_STORES' as const, storeIds: [],
  };
}

export function administratorInviteOutcome(result: { temporaryPassword?: string | null; invitationPending: boolean }): AdministratorInviteOutcome {
  return result.temporaryPassword ? { kind: 'temporary-password', value: result.temporaryPassword } : { kind: 'invitation-pending' };
}

export function emptyAdministratorValues(): AdministratorValues {
  return { phone: '', displayName: '', email: '', roleIds: [] };
}

export function assignableAdministratorRoles<T extends { kind: string }>(roles: T[], actorIsSuper: boolean): T[] {
  return roles.filter((role) => role.kind === 'CUSTOM' || actorIsSuper && role.kind === 'HQ_SUPER_ADMIN');
}

export function AdministratorFormDialog(props: AdministratorFormDialogProps) {
  const values = props.values ?? emptyAdministratorValues();
  const change = (field: 'phone' | 'displayName' | 'email', value: string) => props.onValuesChange?.({ ...values, [field]: value });
  return (
    <AdminFormDialogShell open={props.open} onOpenChange={props.onOpenChange} title={props.mode === 'create' ? '新增管理员' : '编辑管理员角色'} description={props.mode === 'create' ? '创建总部管理员账号并分配一个或多个角色。' : '账号资料为全局信息，此处仅调整总部角色。'} submitLabel={props.mode === 'create' ? '创建管理员' : '保存角色'} cancelLabel={props.outcome ? '完成' : '取消'} isSubmitting={props.isSubmitting} hideSubmit={Boolean(props.outcome)} onSubmit={props.onSubmit}>
      {props.outcome ? <InviteResult outcome={props.outcome} /> : <div className="flex flex-col gap-5">
        <AdminActionError message={props.error} />
        {props.mode === 'create' ? <>
          <AdminValidatedField label="手机号" placeholder="请输入手机号" value={values.phone} onChange={(value) => change('phone', value)} kind="phone" />
          <AdminValidatedField label="姓名" placeholder="请输入姓名" value={values.displayName} onChange={(value) => change('displayName', value)} maxLength={64} />
          <AdminValidatedField label="邮箱（选填）" placeholder="请输入邮箱" value={values.email} onChange={(value) => change('email', value)} kind="email" required={false} />
        </> : <AccountSummary account={props.account} />}
        <AdminRoleSelect id="administrator-roles" roles={props.roles} value={values.roleIds}
          onValueChange={(roleIds) => props.onValuesChange?.({ ...values, roleIds })} />
      </div>}
    </AdminFormDialogShell>
  );
}

function AccountSummary({ account }: { account?: AdministratorFormDialogProps['account'] }) {
  if (!account) return null;
  return <div className="rounded-lg border bg-muted/30 p-4 text-sm"><p className="font-medium">{account.displayName}</p><p className="mt-1 text-muted-foreground">{account.phone}{account.email ? ` · ${account.email}` : ''}</p></div>;
}

function InviteResult({ outcome }: { outcome: AdministratorInviteOutcome }) {
  if (outcome.kind === 'invitation-pending') return <p className="text-sm">该手机号已有账号，已发送待接受的总部邀请。</p>;
  return <div className="flex flex-col gap-2"><p className="text-sm font-medium">临时密码仅显示一次，请立即安全交付。</p><code className="block rounded bg-muted p-3">{outcome.value}</code><CopyTemporaryPasswordButton password={outcome.value} /></div>;
}
