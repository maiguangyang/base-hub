import { useState } from 'react';
import { Input } from '@/components/ui/input';
import { mobilePhoneInputProps } from '@/lib/mobilePhoneInput';

const emailPattern = '[a-z0-9+_\\-]+(\\.[a-z0-9+_\\-]+)*@([a-z0-9\\-]+\\.)+[a-z]{2,6}';
const emailRule = new RegExp(`^${emailPattern}$`);

interface AdminValidatedFieldProps {
  label: string;
  placeholder: string;
  value: string;
  onChange(value: string): void;
  kind?: 'phone' | 'email' | 'text';
  maxLength?: number;
  required?: boolean;
}

function validationMessage(kind: AdminValidatedFieldProps['kind'], value: string): string | undefined {
  if (!value) return undefined;
  if (kind === 'phone' && !/^1[34578][0-9]{9}$/.test(value)) return '请输入11位有效手机号';
  if (kind === 'email' && !emailRule.test(value)) return '请输入有效的邮箱地址';
  return undefined;
}

export function AdminValidatedField({ label, placeholder, value, onChange, kind = 'text', maxLength, required = true }: AdminValidatedFieldProps) {
  const [touched, setTouched] = useState(false);
  const error = touched ? validationMessage(kind, value) : undefined;
  const inputProps = kind === 'phone' ? mobilePhoneInputProps : kind === 'email' ? { type: 'email', maxLength: 128, pattern: emailPattern, title: '请输入有效的邮箱地址' } : { maxLength };

  return <label className="flex flex-col gap-2 text-sm">
    <span>{label}</span>
    <Input
      placeholder={placeholder}
      value={value}
      onChange={(event) => onChange(event.target.value)}
      onBlur={() => setTouched(true)}
      required={required}
      {...inputProps}
      aria-invalid={Boolean(error)}
    />
    {error && <span role="alert" className="text-xs text-destructive">{error}</span>}
  </label>;
}
