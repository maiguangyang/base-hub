import {
  Select, SelectContent, SelectItem, SelectTrigger, SelectValue,
} from '@/components/ui/select';
import { AdminFilterClearButton } from '../AdminFilterClearButton';

/** 由业务目录提供的真实筛选项，不包含占位提示。 */
export interface AdminFilterOption {
  value: string;
  label: string;
}

/** 空值显示占位提示；支持单项清空与列表统一重置。 */
export interface AdminFilterSelectProps {
  value: string | undefined;
  placeholder: string;
  options: readonly AdminFilterOption[];
  disabled?: boolean;
  onValueChange: (value: string | undefined) => void;
}

/** 后台列表统一单选筛选；placeholder 只用于控件提示，不生成选项。 */
export function AdminFilterSelect(props: AdminFilterSelectProps) {
  return (
    <div className="group/filter relative w-full sm:w-40"><Select
      value={props.value ?? ''}
      disabled={props.disabled}
      onValueChange={(value) => props.onValueChange(value || undefined)}
    >
      <SelectTrigger aria-label={props.placeholder} className={`h-10 w-full bg-card shadow-none [&_[data-slot=select-value]]:truncate ${props.value ? '[&_[data-slot=select-value]]:pr-6' : ''}`}>
        <SelectValue placeholder={props.placeholder} />
      </SelectTrigger>
      <SelectContent>
        {props.options.map((option) => (
          <SelectItem key={option.value} value={option.value}>{option.label}</SelectItem>
        ))}
      </SelectContent>
    </Select>
      {props.value && <AdminFilterClearButton label={props.placeholder} disabled={props.disabled}
        className="right-7" onClear={() => props.onValueChange(undefined)} />}
    </div>
  );
}
