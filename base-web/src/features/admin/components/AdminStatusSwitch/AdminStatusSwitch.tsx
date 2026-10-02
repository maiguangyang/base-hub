import { Switch } from '@/components/ui/switch';

/** 列表主状态开关。§4.4 规定启用/禁用类主状态必须用它，
 *  页面禁止直接使用原始 Switch，也禁止手写状态颜色。 */
export interface AdminStatusSwitchProps {
  checked: boolean;
  /** 无障碍名称，说明这个开关控制的是什么。 */
  label: string;
  onCheckedChange: (next: boolean) => void;
  disabled?: boolean;
}

export function AdminStatusSwitch({ checked, label, onCheckedChange, disabled = false }: AdminStatusSwitchProps) {
  return <Switch checked={checked} aria-label={label} disabled={disabled} onCheckedChange={onCheckedChange} className="data-[state=checked]:bg-success-fg" />;
}
