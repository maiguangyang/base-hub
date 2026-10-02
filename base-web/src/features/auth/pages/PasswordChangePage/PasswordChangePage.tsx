import { useState, type SubmitEvent } from 'react';
import { useMutation } from '@apollo/client/react';
import { useFragment } from '@/__generated__';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { AuthPageShell } from '@/features/auth/components/AuthPageShell';
import { CHANGE_TEMPORARY_PASSWORD_MUTATION, VIEWER_FIELDS_FRAGMENT } from '@/features/auth/graphql/auth';
import { useAuthStore } from '@/features/auth/store/authStore';
import { getGraphQLErrorCode } from '@/lib/graphql/errors';
import { disposeGraphQLRuntime } from '@/lib/graphql/client';
import { browserReturnDestination, completePasswordChange, validateNewPassword } from '../authFlow';
import { authErrorMessage, authMessages } from '../authMessages';

/** 首次登录强制修改临时密码页。 */
export function PasswordChangePage() {
  const setViewer = useAuthStore((state) => state.setViewer);
  const [changePassword, { loading }] = useMutation(CHANGE_TEMPORARY_PASSWORD_MUTATION);
  const [values, setValues] = useState({ current: '', password: '', confirmation: '' });
  const [error, setError] = useState<string>();

  async function submit(event: SubmitEvent<HTMLFormElement>) {
    event.preventDefault();
    const submitted = values;
    setValues({ current: '', password: '', confirmation: '' });
    const validation = validateNewPassword(submitted.password, submitted.confirmation);
    if (validation) return setError(authErrorMessage(validation));
    setError(undefined);
    try {
      const result = await changePassword({ variables: { input: { currentPassword: submitted.current, newPassword: submitted.password } } });
      const viewer = useFragment(VIEWER_FIELDS_FRAGMENT, result.data?.changeTemporaryPassword);
      if (!viewer) throw new Error('AUTH_REQUIRED');
      await completePasswordChange(viewer, setViewer, disposeGraphQLRuntime, (path) => window.location.replace(path), browserReturnDestination());
    } catch (cause) {
      setError(authErrorMessage(getGraphQLErrorCode(cause)));
    }
  }

  function change(field: keyof typeof values, value: string) {
    setValues((current) => ({ ...current, [field]: value }));
  }

  return (
    <AuthPageShell title={authMessages.passwordChange.title} description={authMessages.passwordChange.description}>
      <form className="flex flex-col gap-5" onSubmit={submit}>
        <PasswordField label={authMessages.passwordChange.current} placeholder={authMessages.passwordChange.currentPlaceholder} value={values.current} autoComplete="current-password" onChange={(value) => change('current', value)} />
        <PasswordField label={authMessages.passwordChange.next} placeholder={authMessages.passwordChange.nextPlaceholder} value={values.password} autoComplete="new-password" onChange={(value) => change('password', value)} />
        <PasswordField label={authMessages.passwordChange.confirmation} placeholder={authMessages.passwordChange.confirmationPlaceholder} value={values.confirmation} autoComplete="new-password" onChange={(value) => change('confirmation', value)} />
        <p className="text-xs text-muted-foreground">{authMessages.passwordChange.policy}</p>
        {error && <p role="alert" className="text-sm text-destructive">{error}</p>}
        <Button className="w-full" disabled={loading} type="submit">{loading ? authMessages.common.submitting : authMessages.passwordChange.submit}</Button>
      </form>
    </AuthPageShell>
  );
}

/** 密码输入框字段属性。 */
interface PasswordFieldProps {
  label: string;
  placeholder: string;
  value: string;
  autoComplete: string;
  onChange(value: string): void;
}

/** 密码输入控件封装，统一绑定 label 与 placeholder。 */
function PasswordField({ label, placeholder, value, autoComplete, onChange }: PasswordFieldProps) {
  return (
    <label className="flex flex-col gap-2 text-sm font-medium">
      <span>{label}</span>
      <Input type="password" placeholder={placeholder} value={value} autoComplete={autoComplete} onChange={(event) => onChange(event.target.value)} required />
    </label>
  );
}
