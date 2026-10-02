import { AdminValidatedField } from '@/features/admin/components/AdminValidatedField';
import { CopyTemporaryPasswordButton } from '@/features/admin/components/CopyTemporaryPasswordButton';
import { AdminFormDialogShell } from '@/features/admin/components/AdminFormDialogShell';
import { AdminFormSelect } from '@/features/admin/components/AdminFormSelect';
import { AdminRoleSelect } from '@/features/admin/components/AdminRoleSelect/AdminRoleSelect';
import { Checkbox } from '@/components/ui/checkbox';
import { Label } from '@/components/ui/label';

export type StaffAccessMode = 'ALL_STORES' | 'SELECTED_STORES';
export interface InviteStaffValues { phone: string; displayName: string; email: string; roleIds: string[]; storeAccessMode: StaffAccessMode; storeIds: string[] }
export type StaffInviteOutcome = { kind: 'temporary-password'; value: string } | { kind: 'invitation-pending' };

interface InviteStaffDialogProps {
  open: boolean;
  onOpenChange(open: boolean): void;
  values?: InviteStaffValues;
  onValuesChange?(values: InviteStaffValues): void;
  onSubmit(): void | Promise<void>;
  isSubmitting?: boolean;
  outcome?: StaffInviteOutcome;
  roles?: { id: string; name: string; kind?: string }[];
  stores?: { id: string; name: string }[];
}

export function staffAccessInput(storeAccessMode: StaffAccessMode, storeIds: string[]) {
  if (storeAccessMode === 'ALL_STORES') return { storeAccessMode, storeIds: [] };
  if (storeIds.length === 0) throw new Error('VALIDATION_FAILED');
  return { storeAccessMode, storeIds };
}

export function staffInviteOutcome(result: { temporaryPassword?: string | null; invitationPending: boolean }): StaffInviteOutcome {
  return result.temporaryPassword ? { kind: 'temporary-password', value: result.temporaryPassword } : { kind: 'invitation-pending' };
}

export function InviteStaffDialog(props: InviteStaffDialogProps) {
  const values = props.values ?? emptyInviteValues();
  const change = (field: keyof InviteStaffValues, value: string) => props.onValuesChange?.({ ...values, [field]: value });
  const toggleStore = (id: string) => props.onValuesChange?.({ ...values,
    storeIds: values.storeIds.includes(id) ? values.storeIds.filter((item) => item !== id) : [...values.storeIds, id],
  });
  return (
    <AdminFormDialogShell open={props.open} onOpenChange={props.onOpenChange} title="新增员工" description="可为同一员工分配多个角色和门店范围。" submitLabel="新增员工" cancelLabel={props.outcome ? '完成' : '关闭'} isSubmitting={props.isSubmitting} hideSubmit={Boolean(props.outcome)} onSubmit={props.onSubmit}>
      {props.outcome ? <InviteResult outcome={props.outcome} /> : <div className="flex flex-col gap-5">
        <AdminValidatedField label="手机号" placeholder="请输入手机号" value={values.phone} onChange={(value) => change('phone', value)} kind="phone" />
        <AdminValidatedField label="姓名" placeholder="请输入姓名" value={values.displayName} onChange={(value) => change('displayName', value)} maxLength={64} />
        <AdminValidatedField label="邮箱（选填）" placeholder="请输入邮箱" value={values.email} onChange={(value) => change('email', value)} kind="email" required={false} />
        <AdminRoleSelect id="invite-staff-roles" roles={props.roles ?? []} value={values.roleIds}
          onValueChange={(roleIds) => props.onValuesChange?.({ ...values, roleIds })} />
        <div className="flex flex-col gap-2 text-sm"><Label htmlFor="invite-store-access">门店范围</Label>
          <AdminFormSelect id="invite-store-access" value={values.storeAccessMode} placeholder="请选择门店范围"
            options={[{ value: 'ALL_STORES', label: '全部门店' }, { value: 'SELECTED_STORES', label: '指定门店' }]}
            onValueChange={(value) => props.onValuesChange?.({ ...values, storeAccessMode: value as StaffAccessMode, storeIds: value === 'ALL_STORES' ? [] : values.storeIds })} /></div>
        {values.storeAccessMode === 'SELECTED_STORES' && <fieldset className="flex flex-col gap-5 text-sm"><legend>指定门店</legend>{(props.stores ?? []).map((store) => <label key={store.id} className="flex items-center gap-2"><Checkbox aria-label={store.name} checked={values.storeIds.includes(store.id)} onCheckedChange={() => toggleStore(store.id)} />{store.name}</label>)}</fieldset>}
      </div>}
    </AdminFormDialogShell>
  );
}

export function emptyInviteValues(): InviteStaffValues {
  return { phone: '', displayName: '', email: '', roleIds: [], storeAccessMode: 'ALL_STORES', storeIds: [] };
}

function InviteResult({ outcome }: { outcome: StaffInviteOutcome }) {
  if (outcome.kind === 'invitation-pending') return <p className="text-sm">既有账号已收到待接受邀请。</p>;
  return <div className="flex flex-col gap-2"><p className="text-sm font-medium">临时密码仅显示一次。</p><code className="block rounded bg-muted p-2">{outcome.value}</code><CopyTemporaryPasswordButton password={outcome.value} /></div>;
}
