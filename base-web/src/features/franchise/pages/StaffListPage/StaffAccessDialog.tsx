import { AdminFormDialogShell } from '@/features/admin/components/AdminFormDialogShell';
import { roleDisplayName } from '@/features/admin/config/roleDisplayName';
import { AdminFormSelect } from '@/features/admin/components/AdminFormSelect';
import { Checkbox } from '@/components/ui/checkbox';
import { Label } from '@/components/ui/label';
import type { StaffAccessMode } from './InviteStaffDialog';

export interface StaffAccessValues {
  roleIds: string[];
  storeAccessMode: StaffAccessMode;
  storeIds: string[];
}

interface StaffAccessDialogProps {
  open: boolean;
  onOpenChange(open: boolean): void;
  values?: StaffAccessValues;
  onValuesChange?(values: StaffAccessValues): void;
  onSubmit(): void | Promise<void>;
  isSubmitting?: boolean;
  roles?: { id: string; name: string; kind?: string }[];
  stores?: { id: string; name: string }[];
}

interface StaffAccessSource {
  storeAccessMode: StaffAccessMode;
  roles: { id: string }[];
  stores: { id: string }[];
}

export function staffAccessValues(source: StaffAccessSource): StaffAccessValues {
  return {
    roleIds: source.roles.map((role) => role.id),
    storeAccessMode: source.storeAccessMode,
    storeIds: source.stores.map((store) => store.id),
  };
}

export function StaffAccessDialog(props: StaffAccessDialogProps) {
  const values = props.values ?? emptyStaffAccessValues();
  const toggle = (field: 'roleIds' | 'storeIds', id: string) => {
    const selected = values[field];
    props.onValuesChange?.({ ...values, [field]: selected.includes(id) ? selected.filter((item) => item !== id) : [...selected, id] });
  };
  const changeMode = (mode: StaffAccessMode) => props.onValuesChange?.({
    ...values, storeAccessMode: mode, storeIds: mode === 'ALL_STORES' ? [] : values.storeIds,
  });
  return (
    <AdminFormDialogShell open={props.open} onOpenChange={props.onOpenChange} title="编辑员工权限" description="调整员工的多角色授权和可操作门店范围。" submitLabel="保存" cancelLabel="取消" isSubmitting={props.isSubmitting} onSubmit={props.onSubmit}>
      <div className="flex flex-col gap-5">
        <fieldset className="flex flex-col gap-5 text-sm"><legend>角色（可多选）</legend>{(props.roles ?? []).map((role) => <label key={role.id} className="flex items-center gap-2"><Checkbox aria-label={roleDisplayName(role)} checked={values.roleIds.includes(role.id)} onCheckedChange={() => toggle('roleIds', role.id)} />{roleDisplayName(role)}</label>)}</fieldset>
        <div className="flex flex-col gap-2 text-sm"><Label htmlFor="staff-store-access">门店范围</Label>
          <AdminFormSelect id="staff-store-access" value={values.storeAccessMode} placeholder="请选择门店范围"
            options={[{ value: 'ALL_STORES', label: '全部门店' }, { value: 'SELECTED_STORES', label: '指定门店' }]}
            onValueChange={(value) => changeMode(value as StaffAccessMode)} /></div>
        {values.storeAccessMode === 'SELECTED_STORES' && <fieldset className="flex flex-col gap-5 text-sm"><legend>指定门店</legend>{(props.stores ?? []).map((store) => <label key={store.id} className="flex items-center gap-2"><Checkbox aria-label={store.name} checked={values.storeIds.includes(store.id)} onCheckedChange={() => toggle('storeIds', store.id)} />{store.name}</label>)}</fieldset>}
      </div>
    </AdminFormDialogShell>
  );
}

export function emptyStaffAccessValues(): StaffAccessValues {
  return { roleIds: [], storeAccessMode: 'ALL_STORES', storeIds: [] };
}
