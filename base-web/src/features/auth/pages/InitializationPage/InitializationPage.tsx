import { useState, type SubmitEvent } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { mobilePhoneInputProps } from '@/lib/mobilePhoneInput';
import { AuthPageShell } from '@/features/auth/components/AuthPageShell';
import {
  initializeSystem, SystemInitializationError, type SystemInitializationInput,
} from '@/features/auth/api/systemInitialization';
import { browserReturnDestination, loginPath, validateNewPassword } from '../authFlow';
import { authErrorMessage, authMessages } from '../authMessages';

interface InitializationActions {
  initialize(input: SystemInitializationInput): Promise<void>;
  navigate(path: string): void;
}

const emptyValues: SystemInitializationInput = { phone: '', password: '', passwordConfirmation: '' };

/** 校验并提交一次性初始化；冲突表示另一请求已完成，直接进入登录流程。 */
export async function submitSystemInitialization(
  input: SystemInitializationInput,
  destination: string | undefined,
  actions: InitializationActions,
): Promise<string | undefined> {
  const validation = validateNewPassword(input.password, input.passwordConfirmation);
  if (validation) return validation;
  try {
    await actions.initialize(input);
  } catch (cause) {
    if (cause instanceof SystemInitializationError && cause.code === 'HQ_ALREADY_BOOTSTRAPPED') {
      actions.navigate(loginPath(destination));
      return undefined;
    }
    return cause instanceof SystemInitializationError ? cause.code : 'INTERNAL_ERROR';
  }
  actions.navigate(loginPath(destination));
  return undefined;
}

/** 首次部署时创建最高级总部管理员。 */
export function InitializationPage() {
  const [values, setValues] = useState(emptyValues);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string>();

  async function submit(event: SubmitEvent<HTMLFormElement>) {
    event.preventDefault();
    setSubmitting(true);
    setError(undefined);
    const code = await submitSystemInitialization(values, browserReturnDestination(), {
      initialize: initializeSystem,
      navigate: (path) => window.location.replace(path),
    });
    if (code) {
      setError(authErrorMessage(code));
      setValues((current) => ({ ...current, password: '', passwordConfirmation: '' }));
    }
    setSubmitting(false);
  }

  function change(field: keyof SystemInitializationInput, value: string) {
    setValues((current) => ({ ...current, [field]: value }));
  }

  return (
    <AuthPageShell title={authMessages.initialization.title} description={authMessages.initialization.description}>
      <form className="flex flex-col gap-5" autoComplete="off" onSubmit={submit}>
        <InitializationField label={authMessages.initialization.phone} placeholder={authMessages.initialization.phonePlaceholder} value={values.phone} inputMode="tel" onChange={(value) => change('phone', value)} />
        <InitializationField label={authMessages.initialization.password} placeholder={authMessages.initialization.passwordPlaceholder} value={values.password} type="password" onChange={(value) => change('password', value)} />
        <InitializationField label={authMessages.initialization.confirmation} placeholder={authMessages.initialization.confirmationPlaceholder} value={values.passwordConfirmation} type="password" onChange={(value) => change('passwordConfirmation', value)} />
        <p className="text-xs text-muted-foreground">{authMessages.passwordChange.policy}</p>
        {error && <p role="alert" className="text-sm text-destructive">{error}</p>}
        <Button className="h-11 w-full" disabled={submitting} type="submit">
          {submitting ? authMessages.common.submitting : authMessages.initialization.submit}
        </Button>
      </form>
    </AuthPageShell>
  );
}

interface InitializationFieldProps {
  label: string;
  placeholder: string;
  value: string;
  type?: string;
  inputMode?: 'tel';
  onChange(value: string): void;
}

function InitializationField({ label, placeholder, value, type, inputMode, onChange }: InitializationFieldProps) {
  return (
    <label className="flex flex-col gap-4 text-sm font-medium">
      <span>{label}</span>
      <Input className="h-11 px-4" {...(inputMode ? mobilePhoneInputProps : {})} autoComplete="off" placeholder={placeholder} type={type} value={value} onChange={(event) => onChange(event.target.value)} required />
    </label>
  );
}
