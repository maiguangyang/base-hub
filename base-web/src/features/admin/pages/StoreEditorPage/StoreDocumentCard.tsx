import { useEffect, useRef, useState, type DragEvent } from 'react';
import { Eye, ImagePlus, Trash2 } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { cn } from '@/lib/utils';
import { loadStoreDocumentImage, validateStoreDocument } from './storeDocumentImages';

export interface StoreDocumentChange { file?: File; remove?: boolean }

interface Props {
  label: string;
  storedPath?: string | null;
  change: StoreDocumentChange;
  onChange(change: StoreDocumentChange): void;
  disabled?: boolean;
  saveError?: string;
}

function useDocumentPreview(file?: File, storedPath?: string | null, removed?: boolean) {
  const [url, setURL] = useState<string>();
  const [error, setError] = useState<string>();
  const [retry, setRetry] = useState(0);
  useEffect(() => {
    if (removed || (!file && !storedPath)) { setURL(undefined); setError(undefined); return; }
    if (file) {
      const local = URL.createObjectURL(file);
      setURL(local); setError(undefined);
      return () => URL.revokeObjectURL(local);
    }
    let active = true;
    let source: string | undefined;
    loadStoreDocumentImage(storedPath!).then((blob) => {
      if (!active) return;
      source = URL.createObjectURL(blob);
      setURL(source); setError(undefined);
    }).catch(() => { if (active) { setURL(undefined); setError('证照预览加载失败，请重试。'); } });
    return () => { active = false; if (source) URL.revokeObjectURL(source); };
  }, [file, storedPath, removed, retry]);
  return { url, error, retry: () => setRetry((current) => current + 1) };
}

function PreviewError({ validationError, error, retry }: { validationError?: string; error?: string; retry(): void }) {
  if (!validationError && !error) return null;
  return <div role="alert" className="flex items-center gap-2 text-xs text-destructive">{validationError || error}
    {error && <Button type="button" variant="outline" size="sm" onClick={retry}>重试预览</Button>}
  </div>;
}

function DocumentActions({ label, url, disabled, onRemove }: { label: string; url?: string; disabled?: boolean; onRemove(): void }) {
  return <div className="absolute right-2 top-2 flex gap-2">
    <Dialog><DialogTrigger asChild><Button type="button" variant="outline" size="icon" aria-label={`预览${label}`}
      title={`预览${label}`} disabled={disabled || !url} className="size-11 bg-white text-black"><Eye className="size-4" /></Button></DialogTrigger>
      <DialogContent className="max-h-[90vh] bg-white text-black sm:max-w-3xl"><DialogHeader><DialogTitle>{label}</DialogTitle></DialogHeader>
        {url && <img src={url} alt={label} className="max-h-[75vh] w-full object-contain" />}
      </DialogContent></Dialog>
    <Button type="button" variant="outline" size="icon" aria-label={`删除${label}`} title={`删除${label}`}
      disabled={disabled} className="size-11 bg-white text-black" onClick={onRemove}><Trash2 className="size-4" /></Button>
  </div>;
}

function useDocumentDrop(disabled: boolean | undefined, choose: (file?: File) => void) {
  const [dragging, setDragging] = useState(false);
  const depth = useRef(0);
  return {
    dragging,
    onDragEnter: () => { if (!disabled) { depth.current += 1; setDragging(true); } },
    onDragOver: (event: DragEvent<HTMLDivElement>) => { event.preventDefault(); if (!disabled) setDragging(true); },
    onDragLeave: () => { depth.current = Math.max(0, depth.current - 1); if (depth.current === 0) setDragging(false); },
    onDrop: (event: DragEvent<HTMLDivElement>) => {
      event.preventDefault();
      depth.current = 0;
      setDragging(false);
      if (!disabled) choose(event.dataTransfer.files[0]);
    },
  };
}

export function StoreDocumentCard({ label, storedPath, change, onChange, disabled, saveError }: Props) {
  const [validationError, setValidationError] = useState<string>();
  const inputRef = useRef<HTMLInputElement>(null);
  const { url, error, retry } = useDocumentPreview(change.file, storedPath, change.remove);
  const selected = Boolean(change.file || (storedPath && !change.remove));
  function choose(file?: File) {
    if (!file) return;
    const failure = validateStoreDocument(file);
    setValidationError(failure);
    if (!failure) onChange({ file });
  }
  const drop = useDocumentDrop(disabled, choose);
  return <section aria-label={`${label}图片`} className="flex min-w-0 flex-col gap-3 text-black">
    <div className="text-sm font-medium">{label}</div>
    <Input ref={inputRef} aria-label={`上传${label}`} type="file" accept="image/jpeg,image/png,image/webp" className="sr-only" tabIndex={-1} disabled={disabled}
      onChange={(event) => { choose(event.target.files?.[0]); event.target.value = ''; }} />
    <div className={cn('relative h-36 overflow-hidden rounded-md border border-dashed border-border bg-white sm:h-40',
      drop.dragging && 'border-primary bg-primary/5')}
      onDragEnter={drop.onDragEnter} onDragOver={drop.onDragOver} onDragLeave={drop.onDragLeave} onDrop={drop.onDrop}>
      <button type="button" aria-label={`${selected ? '更换' : '选择'}${label}图片`} disabled={disabled}
        className="flex h-full w-full cursor-pointer items-center justify-center focus-visible:outline-2 focus-visible:outline-offset-[-2px] focus-visible:outline-primary disabled:cursor-not-allowed"
        onClick={() => inputRef.current?.click()}>
        {url ? <img src={url} alt={`${label}缩略图`} className="h-full w-full object-contain" />
          : <span className="flex flex-col items-center gap-2 text-sm text-muted-foreground"><ImagePlus className="size-7" />
            {error ? '预览暂不可用，点击重新上传' : selected ? '正在加载预览' : '点击选择图片，或拖拽到此处'}
          </span>}
      </button>
      {selected && <DocumentActions label={label} url={url} disabled={disabled}
        onRemove={() => { setValidationError(undefined); onChange(storedPath ? { remove: true } : {}); }} />}
    </div>
    <p className="text-xs text-muted-foreground">支持 JPG、PNG、WebP，单张不超过 5 MB</p>
    <PreviewError validationError={validationError} error={error} retry={retry} />
    {saveError && <div role="alert" className="flex flex-wrap items-center gap-2 text-xs text-destructive">
      <span>{saveError}</span>
      <Button type="submit" variant="outline" size="sm" disabled={disabled}>重试保存{label}</Button>
    </div>}
  </section>;
}
