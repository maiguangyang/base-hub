import { useState, type SubmitEvent } from 'react';
import { useMutation } from '@apollo/client/react';
import { useFragment } from '@/__generated__';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { mobilePhoneInputProps } from '@/lib/mobilePhoneInput';
import { AuthPageShell } from '@/features/auth/components/AuthPageShell';
import { LOGIN_MUTATION, VIEWER_FIELDS_FRAGMENT } from '@/features/auth/graphql/auth';
import { useAuthStore } from '@/features/auth/store/authStore';
import { getGraphQLErrorCode } from '@/lib/graphql/errors';
import { disposeGraphQLRuntime } from '@/lib/graphql/client';
import { browserReturnDestination, completeLogin } from '../authFlow';
import { authErrorMessage, authMessages } from '../authMessages';

/** 管理端统一登录页；提交后立即从组件状态清除明文密码。 */
export function LoginPage() {
  const setViewer = useAuthStore((state) => state.setViewer);
  const [login, { loading }] = useMutation(LOGIN_MUTATION);
  const [phone, setPhone] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState<string>();

  async function submit(event: SubmitEvent<HTMLFormElement>) {
    event.preventDefault();
    const submittedPassword = password;
    setPassword('');
    setError(undefined);
    try {
      const result = await login({ variables: { input: { phone, password: submittedPassword } } });
      const viewer = useFragment(VIEWER_FIELDS_FRAGMENT, result.data?.login.viewer);
      if (!viewer) throw new Error('AUTH_REQUIRED');
			await completeLogin(viewer, setViewer, disposeGraphQLRuntime, (path) => window.location.replace(path), browserReturnDestination());
    } catch (cause) {
      setError(authErrorMessage(getGraphQLErrorCode(cause)));
    }
  }

  return (
    <AuthPageShell title={authMessages.login.title}>
      <form className="flex flex-col gap-5" autoComplete="off" onSubmit={submit}>
        <label className="flex flex-col gap-4 text-sm font-medium">
          <span>{authMessages.login.phone}</span>
          <Input className="h-11 px-4" {...mobilePhoneInputProps} autoComplete="off" placeholder={authMessages.login.phonePlaceholder} value={phone} onChange={(event) => setPhone(event.target.value)} required />
        </label>
        <label className="flex flex-col gap-4 text-sm font-medium">
          <span>{authMessages.login.password}</span>
          <Input className="h-11 px-4" autoComplete="off" placeholder={authMessages.login.passwordPlaceholder} type="password" value={password} onChange={(event) => setPassword(event.target.value)} required />
        </label>
        {error && <p role="alert" className="text-sm text-destructive">{error}</p>}
        <Button className="h-11 w-full" disabled={loading} type="submit">{loading ? authMessages.common.submitting : authMessages.login.submit}</Button>
      </form>
    </AuthPageShell>
  );
}
