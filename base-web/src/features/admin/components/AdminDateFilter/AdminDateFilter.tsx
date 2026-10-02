import { Input } from '@/components/ui/input';
import { AdminFilterClearButton } from '../AdminFilterClearButton';

export function AdminDateFilter({ label, value, onChange }: { label: string; value: string; onChange(value: string): void }) {
  return <div className="group/filter relative w-full sm:w-40">
    <Input aria-label={label} type="date" className="h-10 w-full bg-card pr-9 shadow-none" value={value}
      onChange={(event) => onChange(event.target.value)} />
    {value && <AdminFilterClearButton label={label} onClear={() => onChange('')} />}
  </div>;
}
