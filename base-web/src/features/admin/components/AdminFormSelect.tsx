import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';

/** Radix 不允许空字符串选项；仅为显式声明的业务空值项使用内部值。 */
const emptyValue = '__admin_form_empty__';

/** 表单的真实业务选项，可标记已停用项。 */
export interface AdminFormSelectOption {
  value: string;
  label: string;
  disabled?: boolean;
}

/** 占位文案不生成选项；可选关联通过 clearLabel 显式声明空值业务选项。 */
export function AdminFormSelect({ id, value, placeholder, clearLabel, options, disabled, onValueChange, onCloseAutoFocus }: {
  id: string;
  value: string;
  placeholder: string;
  clearLabel?: string;
  options: readonly AdminFormSelectOption[];
  disabled?: boolean;
  onValueChange(value: string): void;
  onCloseAutoFocus?(event: Event): void;
}) {
  return <Select value={value || (clearLabel ? emptyValue : '')} disabled={disabled}
    onValueChange={(next) => onValueChange(next === emptyValue ? '' : next)}>
    <SelectTrigger id={id} className="h-10 w-full"><SelectValue placeholder={placeholder} /></SelectTrigger>
    <SelectContent onCloseAutoFocus={onCloseAutoFocus}>
      {clearLabel && <SelectItem value={emptyValue}>{clearLabel}</SelectItem>}
      {options.map((option) => <SelectItem key={option.value} value={option.value} disabled={option.disabled}>{option.label}</SelectItem>)}
    </SelectContent>
  </Select>;
}
