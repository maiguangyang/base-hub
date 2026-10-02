import { renderToStaticMarkup } from 'react-dom/server';
import { describe, expect, it } from 'vitest';
import {
  parseTime,
  parseTimeRange,
  formatTimeRange,
  TimePicker,
  TimeRangePicker,
  DEFAULT_TIME_RANGE_PRESETS,
} from './time-picker';

describe('time-picker utilities', () => {
  it('correctly parses time string', () => {
    expect(parseTime('09:30')).toEqual({ hour: '09', minute: '30' });
    expect(parseTime('9:5')).toEqual({ hour: '09', minute: '05' });
    expect(parseTime('')).toEqual({ hour: '09', minute: '00' });
    expect(parseTime(null)).toEqual({ hour: '09', minute: '00' });
  });

  it('correctly parses various time range string formats', () => {
    expect(parseTimeRange('09:00 - 22:00')).toEqual({ start: '09:00', end: '22:00' });
    expect(parseTimeRange('10:30 ~ 21:30')).toEqual({ start: '10:30', end: '21:30' });
    expect(parseTimeRange('08:00 至 20:00')).toEqual({ start: '08:00', end: '20:00' });
    expect(parseTimeRange('')).toEqual({ start: '09:00', end: '22:00' });
    expect(parseTimeRange(null)).toEqual({ start: '09:00', end: '22:00' });
  });

  it('correctly formats time range', () => {
    expect(formatTimeRange('09:00', '22:00')).toBe('09:00 - 22:00');
  });

  it('contains expected restaurant business hour presets', () => {
    expect(DEFAULT_TIME_RANGE_PRESETS.length).toBeGreaterThanOrEqual(4);
    expect(DEFAULT_TIME_RANGE_PRESETS.some((p) => p.start === '09:00' && p.end === '22:00')).toBe(true);
    expect(DEFAULT_TIME_RANGE_PRESETS.some((p) => p.start === '00:00' && p.end === '24:00')).toBe(true);
  });
});

describe('TimePicker and TimeRangePicker rendering', () => {
  it('renders TimePicker trigger with value or placeholder', () => {
    const markupWithVal = renderToStaticMarkup(<TimePicker value="14:30" />);
    expect(markupWithVal).toContain('14:30');
    expect(markupWithVal).toContain('data-slot="time-picker-trigger"');

    const markupPlaceholder = renderToStaticMarkup(<TimePicker placeholder="选择时间" />);
    expect(markupPlaceholder).toContain('选择时间');
  });

  it('renders TimeRangePicker trigger with value or placeholder', () => {
    const markupWithVal = renderToStaticMarkup(<TimeRangePicker value="09:00 - 22:00" />);
    expect(markupWithVal).toContain('09:00 - 22:00');
    expect(markupWithVal).toContain('data-slot="time-range-picker-trigger"');

    const markupPlaceholder = renderToStaticMarkup(<TimeRangePicker placeholder="请选择营业时间" />);
    expect(markupPlaceholder).toContain('请选择营业时间');
  });
});
