import { useEffect, useState } from 'react';
import { Search } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { AdminFilterClearButton } from '../AdminFilterClearButton';

export interface AdminListSearchProps {
  value: string;
  placeholder: string;
  submitLabel?: string;
  /** 嵌入已有 form 时避免生成无效的嵌套 form。 */
  embedded?: boolean;
  onSearch: (value: string) => void;
}

/** 后台列表统一搜索框：编辑草稿不触发查询，仅在提交时应用。 */
export function AdminListSearch({ value, placeholder, submitLabel = '搜索', embedded = false, onSearch }: AdminListSearchProps) {
  const [draft, setDraft] = useState(value);

  useEffect(() => setDraft(value), [value]);

  const submitDraft = () => onSearch(draft.trim());
  const fields = <>
    <div className="group/filter relative w-full sm:w-80">
      <Search aria-hidden="true" className="pointer-events-none absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
      <Input
        aria-label={placeholder}
        className="h-10 bg-card pl-9 pr-9 shadow-none"
        value={draft}
        placeholder={placeholder}
        onKeyDown={embedded ? (event) => {
          if (event.key !== 'Enter') return;
          event.preventDefault();
          event.stopPropagation();
          submitDraft();
        } : undefined}
        onChange={(event) => setDraft(event.target.value)}
      />
      {(draft || value) && <AdminFilterClearButton label={placeholder} onClear={() => { setDraft(''); onSearch(''); }} />}
    </div>
    <Button type={embedded ? 'button' : 'submit'} onClick={embedded ? submitDraft : undefined}
      variant="secondary" className="h-10 border border-border bg-card px-5 shadow-none hover:bg-accent">
      {submitLabel}
    </Button>
  </>;

  if (embedded) return <div className="flex w-full flex-col gap-3 sm:w-auto sm:flex-row sm:items-center">{fields}</div>;

  return (
    <form
      className="flex w-full flex-col gap-3 sm:w-auto sm:flex-row sm:items-center"
      onSubmit={(event) => {
        event.preventDefault();
        submitDraft();
      }}
    >
      {fields}
    </form>
  );
}
