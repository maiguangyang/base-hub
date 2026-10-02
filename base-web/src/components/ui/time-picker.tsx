import * as React from 'react';
import { Clock, X } from 'lucide-react';
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover';
import { ScrollArea } from '@/components/ui/scroll-area';
import { Button } from '@/components/ui/button';
import { cn } from '@/lib/utils';

export const HOURS = Array.from({ length: 24 }, (_, i) => String(i).padStart(2, '0'));
export const END_HOURS = Array.from({ length: 25 }, (_, i) => String(i).padStart(2, '0'));
export const MINUTES = Array.from({ length: 12 }, (_, i) => String(i * 5).padStart(2, '0'));

export const DEFAULT_TIME_RANGE_PRESETS = [
  { label: '09:00 - 22:00', start: '09:00', end: '22:00' },
  { label: '10:00 - 22:00', start: '10:00', end: '22:00' },
  { label: '10:30 - 21:30', start: '10:30', end: '21:30' },
  { label: '11:00 - 23:00', start: '11:00', end: '23:00' },
  { label: '全天营业 (00:00 - 24:00)', start: '00:00', end: '24:00' },
];

export function parseTime(timeStr?: string | null): { hour: string; minute: string } {
  if (!timeStr) return { hour: '09', minute: '00' };
  const parts = timeStr.trim().split(':');
  const h = parts[0]?.padStart(2, '0') ?? '09';
  const m = parts[1]?.padStart(2, '0') ?? '00';
  return { hour: h, minute: m };
}

export function parseTimeRange(rangeStr?: string | null): { start: string; end: string } {
  if (!rangeStr) return { start: '09:00', end: '22:00' };
  const parts = rangeStr.split(/[-~至到]/).map((s) => s.trim());
  if (parts.length >= 2 && parts[0] && parts[1]) {
    return { start: parts[0], end: parts[1] };
  }
  return { start: '09:00', end: '22:00' };
}

export function formatTimeRange(start: string, end: string): string {
  return `${start} - ${end}`;
}

export interface TimeRangePickerProps {
  value?: string;
  onChange?(value: string): void;
  placeholder?: string;
  className?: string;
  disabled?: boolean;
  presets?: Array<{ label: string; start: string; end: string }>;
}

export function TimeRangePicker({
  value,
  onChange,
  placeholder = '请选择营业时间段',
  className,
  disabled = false,
  presets = DEFAULT_TIME_RANGE_PRESETS,
}: TimeRangePickerProps) {
  const [open, setOpen] = React.useState(false);

  const initialRange = React.useMemo(() => parseTimeRange(value), [value]);
  const [startHour, setStartHour] = React.useState(parseTime(initialRange.start).hour);
  const [startMinute, setStartMinute] = React.useState(parseTime(initialRange.start).minute);
  const [endHour, setEndHour] = React.useState(parseTime(initialRange.end).hour);
  const [endMinute, setEndMinute] = React.useState(parseTime(initialRange.end).minute);

  React.useEffect(() => {
    if (value) {
      const parsed = parseTimeRange(value);
      const s = parseTime(parsed.start);
      const e = parseTime(parsed.end);
      setStartHour(s.hour);
      setStartMinute(s.minute);
      setEndHour(e.hour);
      setEndMinute(e.minute);
    }
  }, [value]);

  const currentStart = `${startHour}:${startMinute}`;
  const currentEnd = `${endHour}:${endMinute}`;
  const currentRange = formatTimeRange(currentStart, currentEnd);

  function applyPreset(start: string, end: string) {
    const s = parseTime(start);
    const e = parseTime(end);
    setStartHour(s.hour);
    setStartMinute(s.minute);
    setEndHour(e.hour);
    setEndMinute(e.minute);
    onChange?.(formatTimeRange(start, end));
    setOpen(false);
  }

  function handleConfirm() {
    onChange?.(currentRange);
    setOpen(false);
  }

  function handleClear(e: React.MouseEvent) {
    e.stopPropagation();
    onChange?.('');
  }

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <button
          type="button"
          disabled={disabled}
          data-slot="time-range-picker-trigger"
          data-admin-form-surface=""
          className={cn(
            'flex h-9 w-full items-center justify-between rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-xs transition-colors hover:bg-accent/40 focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:cursor-not-allowed disabled:opacity-50 text-left',
            !value && 'text-muted-foreground',
            className,
          )}
        >
          <span className="flex items-center gap-2">
            <Clock className="size-4 shrink-0 text-muted-foreground" />
            <span>{value || placeholder}</span>
          </span>
          {value && !disabled && (
            <span
              role="button"
              tabIndex={0}
              aria-label="清空时间"
              className="rounded p-0.5 text-muted-foreground hover:text-foreground"
              onClick={handleClear}
              onKeyDown={(e) => {
                if (e.key === 'Enter' || e.key === ' ') {
                  e.stopPropagation();
                  onChange?.('');
                }
              }}
            >
              <X className="size-3.5" />
            </span>
          )}
        </button>
      </PopoverTrigger>
      <PopoverContent data-admin-form-surface="" className="w-auto p-4" align="start">
        <div className="flex flex-col gap-3">
          {/* 预设时段 */}
          {presets.length > 0 && (
            <div className="flex flex-wrap gap-1.5 border-b border-border pb-3">
              {presets.map((preset) => (
                <button
                  key={preset.label}
                  type="button"
                  className={cn(
                    'rounded-md border px-2 py-1 text-xs font-medium transition-colors hover:bg-accent',
                    value === formatTimeRange(preset.start, preset.end)
                      ? 'border-primary bg-primary/10 text-primary font-semibold'
                      : 'border-border text-muted-foreground',
                  )}
                  onClick={() => applyPreset(preset.start, preset.end)}
                >
                  {preset.label}
                </button>
              ))}
            </div>
          )}

          {/* 时间段选择器主体 */}
          <div className="flex items-center gap-4">
            {/* 开始时间 */}
            <div className="flex flex-col gap-1.5">
              <span className="text-xs font-semibold text-muted-foreground">
                开始时间: <strong className="text-foreground">{currentStart}</strong>
              </span>
              <div className="flex rounded-md border border-border bg-muted/20 p-1">
                {/* 时 */}
                <ScrollArea className="h-44 w-12 pr-1">
                  <div className="flex flex-col gap-0.5">
                    {HOURS.map((h) => (
                      <button
                        key={h}
                        type="button"
                        className={cn(
                          'rounded py-1 text-center text-xs transition-colors hover:bg-accent',
                          startHour === h && 'bg-primary text-primary-foreground font-semibold',
                        )}
                        onClick={() => setStartHour(h)}
                      >
                        {h}
                      </button>
                    ))}
                  </div>
                </ScrollArea>
                <div className="w-px bg-border my-1" />
                {/* 分 */}
                <ScrollArea className="h-44 w-12 pl-1">
                  <div className="flex flex-col gap-0.5">
                    {MINUTES.map((m) => (
                      <button
                        key={m}
                        type="button"
                        className={cn(
                          'rounded py-1 text-center text-xs transition-colors hover:bg-accent',
                          startMinute === m && 'bg-primary text-primary-foreground font-semibold',
                        )}
                        onClick={() => setStartMinute(m)}
                      >
                        {m}
                      </button>
                    ))}
                  </div>
                </ScrollArea>
              </div>
            </div>

            <div className="text-muted-foreground font-medium pt-5">至</div>

            {/* 结束时间 */}
            <div className="flex flex-col gap-1.5">
              <span className="text-xs font-semibold text-muted-foreground">
                结束时间: <strong className="text-foreground">{currentEnd}</strong>
              </span>
              <div className="flex rounded-md border border-border bg-muted/20 p-1">
                {/* 时 */}
                <ScrollArea className="h-44 w-12 pr-1">
                  <div className="flex flex-col gap-0.5">
                    {END_HOURS.map((h) => (
                      <button
                        key={h}
                        type="button"
                        className={cn(
                          'rounded py-1 text-center text-xs transition-colors hover:bg-accent',
                          endHour === h && 'bg-primary text-primary-foreground font-semibold',
                        )}
                        onClick={() => setEndHour(h)}
                      >
                        {h}
                      </button>
                    ))}
                  </div>
                </ScrollArea>
                <div className="w-px bg-border my-1" />
                {/* 分 */}
                <ScrollArea className="h-44 w-12 pl-1">
                  <div className="flex flex-col gap-0.5">
                    {MINUTES.map((m) => (
                      <button
                        key={m}
                        type="button"
                        className={cn(
                          'rounded py-1 text-center text-xs transition-colors hover:bg-accent',
                          endMinute === m && 'bg-primary text-primary-foreground font-semibold',
                        )}
                        onClick={() => setEndMinute(m)}
                      >
                        {m}
                      </button>
                    ))}
                  </div>
                </ScrollArea>
              </div>
            </div>
          </div>

          {/* 底部确认栏 */}
          <div className="flex items-center justify-between border-t border-border pt-3">
            <span className="text-xs text-muted-foreground">
              选中: <span className="font-medium text-foreground">{currentRange}</span>
            </span>
            <div className="flex gap-2">
              <Button type="button" variant="outline" size="sm" onClick={() => setOpen(false)}>
                取消
              </Button>
              <Button type="button" size="sm" onClick={handleConfirm}>
                确定
              </Button>
            </div>
          </div>
        </div>
      </PopoverContent>
    </Popover>
  );
}

export interface TimePickerProps {
  value?: string;
  onChange?(value: string): void;
  placeholder?: string;
  className?: string;
  disabled?: boolean;
}

export function TimePicker({
  value,
  onChange,
  placeholder = '请选择时间',
  className,
  disabled = false,
}: TimePickerProps) {
  const [open, setOpen] = React.useState(false);
  const parsed = React.useMemo(() => parseTime(value), [value]);
  const [hour, setHour] = React.useState(parsed.hour);
  const [minute, setMinute] = React.useState(parsed.minute);

  React.useEffect(() => {
    if (value) {
      const p = parseTime(value);
      setHour(p.hour);
      setMinute(p.minute);
    }
  }, [value]);

  const currentTime = `${hour}:${minute}`;

  function handleConfirm() {
    onChange?.(currentTime);
    setOpen(false);
  }

  function handleClear(e: React.MouseEvent) {
    e.stopPropagation();
    onChange?.('');
  }

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <button
          type="button"
          disabled={disabled}
          data-slot="time-picker-trigger"
          data-admin-form-surface=""
          className={cn(
            'flex h-9 w-full items-center justify-between rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-xs transition-colors hover:bg-accent/40 focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:cursor-not-allowed disabled:opacity-50 text-left',
            !value && 'text-muted-foreground',
            className,
          )}
        >
          <span className="flex items-center gap-2">
            <Clock className="size-4 shrink-0 text-muted-foreground" />
            <span>{value || placeholder}</span>
          </span>
          {value && !disabled && (
            <span
              role="button"
              tabIndex={0}
              aria-label="清空时间"
              className="rounded p-0.5 text-muted-foreground hover:text-foreground"
              onClick={handleClear}
              onKeyDown={(e) => {
                if (e.key === 'Enter' || e.key === ' ') {
                  e.stopPropagation();
                  onChange?.('');
                }
              }}
            >
              <X className="size-3.5" />
            </span>
          )}
        </button>
      </PopoverTrigger>
      <PopoverContent data-admin-form-surface="" className="w-auto p-4" align="start">
        <div className="flex flex-col gap-3">
          <div className="flex items-center justify-between border-b border-border pb-2">
            <span className="text-xs font-semibold text-muted-foreground">
              当前时间: <strong className="text-foreground">{currentTime}</strong>
            </span>
          </div>

          <div className="flex rounded-md border border-border bg-muted/20 p-1">
            {/* 时 */}
            <ScrollArea className="h-44 w-14 pr-1">
              <div className="flex flex-col gap-0.5">
                {HOURS.map((h) => (
                  <button
                    key={h}
                    type="button"
                    className={cn(
                      'rounded py-1 text-center text-xs transition-colors hover:bg-accent',
                      hour === h && 'bg-primary text-primary-foreground font-semibold',
                    )}
                    onClick={() => setHour(h)}
                  >
                    {h} 时
                  </button>
                ))}
              </div>
            </ScrollArea>
            <div className="w-px bg-border my-1" />
            {/* 分 */}
            <ScrollArea className="h-44 w-14 pl-1">
              <div className="flex flex-col gap-0.5">
                {MINUTES.map((m) => (
                  <button
                    key={m}
                    type="button"
                    className={cn(
                      'rounded py-1 text-center text-xs transition-colors hover:bg-accent',
                      minute === m && 'bg-primary text-primary-foreground font-semibold',
                    )}
                    onClick={() => setMinute(m)}
                  >
                    {m} 分
                  </button>
                ))}
              </div>
            </ScrollArea>
          </div>

          <div className="flex items-center justify-between border-t border-border pt-3">
            <Button type="button" variant="outline" size="sm" onClick={() => setOpen(false)}>
              取消
            </Button>
            <Button type="button" size="sm" onClick={handleConfirm}>
              确定
            </Button>
          </div>
        </div>
      </PopoverContent>
    </Popover>
  );
}
