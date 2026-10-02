import { useRef, useState, type KeyboardEvent } from 'react';
import { Check, ChevronDown, Search } from 'lucide-react';
import { Input } from '@/components/ui/input';
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover';
import type { AdminFilterOption } from './AdminFilterSelect';
import { AdminFilterClearButton } from './AdminFilterClearButton';

interface AdminSearchSelectProps {
  value: string | undefined;
  options: readonly AdminFilterOption[];
  label: string;
  placeholder: string;
  searchPlaceholder: string;
  emptyMessage: string;
  disabled?: boolean;
  clearLabel?: string;
  compact?: boolean;
  clearable?: boolean;
  onValueChange(value: string | undefined): void;
}

/** 后台表单与列表共用的可搜索单选框。 */
export function AdminSearchSelect(props: AdminSearchSelectProps) {
  const { value, options, label, searchPlaceholder, emptyMessage, clearLabel, compact, onValueChange } = props;
  const [open, setOpen] = useState(false);
  const [query, setQuery] = useState('');
  const searchRef = useRef<HTMLInputElement>(null);
  const listRef = useRef<HTMLDivElement>(null);
  const selected = options.find((option) => option.value === value);
  const matches = options.filter((option) => option.label.toLocaleLowerCase().includes(query.trim().toLocaleLowerCase()));
  const choose = (next?: string) => { onValueChange(next); setOpen(false); setQuery(''); };
  const focusFirst = (event: KeyboardEvent<HTMLInputElement>) => {
    if (event.key !== 'ArrowDown') return;
    event.preventDefault();
    listRef.current?.querySelector<HTMLButtonElement>('[role="option"]')?.focus();
  };
  const moveFocus = (event: KeyboardEvent<HTMLButtonElement>) => {
    if (event.key !== 'ArrowDown' && event.key !== 'ArrowUp') return;
    event.preventDefault();
    const items = Array.from(listRef.current?.querySelectorAll<HTMLButtonElement>('[role="option"]') ?? []);
    const next = items.indexOf(event.currentTarget) + (event.key === 'ArrowDown' ? 1 : -1);
    if (next < 0) searchRef.current?.focus();
    else items[next]?.focus();
  };

  return <Popover open={open} onOpenChange={(next) => { setOpen(next); if (!next) setQuery(''); }}>
    <SearchSelectTrigger {...props} selectedLabel={selected?.label} open={open} onClear={() => choose()} />
    <PopoverContent data-admin-form-surface="" align="start" className={compact ? 'w-64 p-1' : 'w-[var(--radix-popover-trigger-width)] p-1'}>
      <div className="relative border-b pb-1">
        <Search aria-hidden="true" className="pointer-events-none absolute left-2.5 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
        <Input ref={searchRef} aria-label={searchPlaceholder} placeholder={searchPlaceholder} value={query} onChange={(event) => setQuery(event.target.value)} onKeyDown={focusFirst} className="border-0 pl-9 shadow-none focus-visible:ring-0" />
      </div>
      <div ref={listRef} role="listbox" aria-label={`${label}选项`} className="max-h-64 overflow-y-auto p-1">
        {clearLabel && !query.trim() && <button type="button" role="option" aria-selected={!value} className="flex w-full items-center justify-between rounded-sm px-2 py-2 text-left text-sm hover:bg-accent focus-visible:bg-accent focus-visible:outline-none aria-selected:bg-accent aria-selected:text-primary" onKeyDown={moveFocus} onClick={() => choose()}>
          {clearLabel}{!value && <Check aria-hidden="true" className="size-4 text-primary" />}
        </button>}
        {matches.map((option) => <button key={option.value} type="button" role="option" aria-selected={option.value === value} title={option.label} className="flex w-full items-center justify-between gap-2 rounded-sm px-2 py-2 text-left text-sm hover:bg-accent focus-visible:bg-accent focus-visible:outline-none aria-selected:bg-accent aria-selected:text-primary" onKeyDown={moveFocus} onClick={() => choose(option.value)}>
          <span className="truncate">{option.label}</span>{option.value === value && <Check aria-hidden="true" className="size-4 shrink-0 text-primary" />}
        </button>)}
        {matches.length === 0 && <p className="px-2 py-4 text-center text-sm text-muted-foreground">{emptyMessage}</p>}
      </div>
    </PopoverContent>
  </Popover>;
}

function SearchSelectTrigger({ label, placeholder, selectedLabel, value, compact, clearable, disabled, open, onClear }: AdminSearchSelectProps & {
  selectedLabel?: string;
  open: boolean;
  onClear(): void;
}) {
  const canClear = Boolean(clearable && value);
  return <div className={`group/filter relative w-full${compact ? ' sm:w-40' : ''}`}>
    <PopoverTrigger asChild>
      <button type="button" role="combobox" aria-label={label} aria-haspopup="listbox" aria-expanded={open} disabled={disabled} title={selectedLabel}
        data-admin-form-surface=""
        className="flex h-10 w-full items-center justify-between gap-2 rounded-md border border-input bg-card px-3 text-sm font-normal shadow-none outline-none focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:cursor-not-allowed disabled:opacity-50">
        <span className={`truncate${canClear ? ' pr-6' : ''}`}>{selectedLabel ?? placeholder}</span><ChevronDown className="size-4 shrink-0 text-muted-foreground" aria-hidden="true" />
      </button>
    </PopoverTrigger>
    {canClear && <AdminFilterClearButton label={label} disabled={disabled} className="right-7" onClear={onClear} />}
  </div>;
}
