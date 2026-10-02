import { useState } from 'react';
import { useMutation } from '@apollo/client/react';
import { useFragment } from '@/__generated__';
import { Input } from '@/components/ui/input';
import { AdminFormDialogShell } from '@/features/admin/components/AdminFormDialogShell';
import { CHANGE_TEMPORARY_PASSWORD_MUTATION, SELECT_WORKSPACE_MUTATION, VIEWER_FIELDS_FRAGMENT } from '@/features/auth/graphql/auth';
import { completePasswordChange, currentDestination, validateNewPassword } from '@/features/auth/pages/authFlow';
import { authErrorMessage, authMessages } from '@/features/auth/pages/authMessages';
import { useAuthStore } from '@/features/auth/store/authStore';
import { reportNonBlockingError } from '@/lib/diagnostics';
import { disposeGraphQLRuntime } from '@/lib/graphql/client';
import { getGraphQLErrorCode } from '@/lib/graphql/errors';

interface AdminPasswordChangeDialogProps {
  open: boolean;
  onOpenChange(open: boolean): void;
}

/** 已登录用户的改密入口保持在后台壳内。 */
export function AdminPasswordChangeDialog({ open, onOpenChange }: AdminPasswordChangeDialogProps) {
  const form = usePasswordDialogState(onOpenChange);

  return (
    <AdminFormDialogShell
      open={open} onOpenChange={form.changeOpen} title="修改密码"
      description="修改后，其他已登录设备将退出当前账号。"
      submitLabel={form.submitting ? authMessages.common.submitting : '保存新密码'}
      cancelLabel="取消" isSubmitting={form.submitting} onSubmit={form.submit}
    >
      <PasswordFields values={form.values} error={form.error} onChange={form.change} />
    </AdminFormDialogShell>
  );
}

function usePasswordDialogState(onOpenChange: (open: boolean) => void) {
  const [values, setValues] = useState({ current: '', password: '', confirmation: '' });
  const [error, setError] = useState<string>();
  const action = usePasswordAction();

  function change(field: keyof typeof values, value: string) {
    setValues((current) => ({ ...current, [field]: value }));
  }

  function changeOpen(nextOpen: boolean) {
    if (!nextOpen) {
      setValues({ current: '', password: '', confirmation: '' });
      setError(undefined);
    }
    onOpenChange(nextOpen);
  }

  function submit() {
    const validation = validateNewPassword(values.password, values.confirmation);
    if (validation) return setError(authErrorMessage(validation));
    return action.submit(values, () => setValues({ current: '', password: '', confirmation: '' }), setError);
  }

  return { values, error, submitting: action.submitting, change, changeOpen, submit };
}

type PasswordValues = { current: string; password: string; confirmation: string };

function usePasswordAction() {
  const viewer = useAuthStore((state) => state.viewer);
  const setViewer = useAuthStore((state) => state.setViewer);
  const [changePassword, changeState] = useMutation(CHANGE_TEMPORARY_PASSWORD_MUTATION);
  const [selectWorkspace, selectState] = useMutation(SELECT_WORKSPACE_MUTATION);

  async function submit(values: PasswordValues, clear: () => void, setError: (message?: string) => void) {
    const workspace = viewer?.currentWorkspace;
    if (!workspace || workspace.workspaceType === 'DISCOVERY') return setError(authErrorMessage('AUTH_REQUIRED'));
    setError(undefined);
    try {
      const changed = await changePassword({ variables: { input: { currentPassword: values.current, newPassword: values.password } } });
      const discoveryViewer = useFragment(VIEWER_FIELDS_FRAGMENT, changed.data?.changeTemporaryPassword);
      if (!discoveryViewer) throw new Error('AUTH_REQUIRED');
      let finalViewer = discoveryViewer;
      try {
        const selected = await selectWorkspace({ variables: { input: { workspaceType: workspace.workspaceType, organizationId: workspace.organizationId } } });
        finalViewer = useFragment(VIEWER_FIELDS_FRAGMENT, selected.data?.selectWorkspace) ?? discoveryViewer;
      } catch (cause) {
        reportNonBlockingError('改密后的 workspace 恢复失败，继续进入工作台选择流程。', cause);
      }
      clear();
      await completePasswordChange(finalViewer, setViewer, disposeGraphQLRuntime, (path) => window.location.replace(path), currentDestination(window.location));
    } catch (cause) {
      setError(authErrorMessage(getGraphQLErrorCode(cause)));
    }
  }

  return { submitting: changeState.loading || selectState.loading, submit };
}

function PasswordFields({ values, error, onChange }: { values: PasswordValues; error?: string; onChange(field: keyof PasswordValues, value: string): void }) {
  return <div className="flex flex-col gap-5">
    <PasswordField label="当前密码" placeholder="请输入当前密码" value={values.current} autoComplete="current-password" onChange={(value) => onChange('current', value)} />
    <PasswordField label="新密码" placeholder="请输入新密码" value={values.password} autoComplete="new-password" onChange={(value) => onChange('password', value)} />
    <PasswordField label="确认新密码" placeholder="请再次输入新密码" value={values.confirmation} autoComplete="new-password" onChange={(value) => onChange('confirmation', value)} />
    <p className="text-xs text-muted-foreground">{authMessages.passwordChange.policy}</p>
    {error && <p role="alert" className="text-sm text-destructive">{error}</p>}
  </div>;
}

interface PasswordFieldProps {
  label: string;
  placeholder: string;
  value: string;
  autoComplete: string;
  onChange(value: string): void;
}

function PasswordField({ label, placeholder, value, autoComplete, onChange }: PasswordFieldProps) {
  return (
    <label className="flex flex-col gap-2 text-sm font-medium">
      <span>{label}</span>
      <Input className="h-11" type="password" placeholder={placeholder} value={value} autoComplete={autoComplete} onChange={(event) => onChange(event.target.value)} required />
    </label>
  );
}
