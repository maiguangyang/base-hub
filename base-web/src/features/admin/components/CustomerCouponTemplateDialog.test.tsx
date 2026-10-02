// @vitest-environment jsdom

import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { CustomerCouponTemplateDialog } from './CustomerCouponTemplateDialog';

afterEach(() => {
  cleanup();
  vi.useRealTimers();
});

beforeEach(() => {
  vi.useFakeTimers({ shouldAdvanceTime: true });
  vi.setSystemTime(new Date('2026-06-15T10:00:00.000Z'));
});

function renderDialog(onSave = vi.fn().mockResolvedValue(undefined), issuer: 'headquarters' | 'store' = 'headquarters') {
  render(<CustomerCouponTemplateDialog issuer={issuer} open busy={false} onOpenChange={vi.fn()} onSave={onSave} />);
  fireEvent.change(screen.getByLabelText('名称'), { target: { value: '欢迎券' } });
  fireEvent.change(screen.getByLabelText('面额（元）'), { target: { value: '5.00' } });
  return onSave;
}

function submit() {
  fireEvent.click(screen.getByRole('button', { name: /创建模板|保存模板/ }));
}

it('does not render redundant field descriptions', () => {
  renderDialog();
  expect(screen.queryByText('留空立即生效。')).toBeNull();
  expect(screen.queryByText('会员激活后 1～365 天内有效。')).toBeNull();
  expect(screen.queryByText('留空不限量；填写后按确定总量发放。')).toBeNull();
  expect(screen.queryByText('留空持续补发；截止后不再向新注册会员发放。')).toBeNull();
});

it.each(['', '0'])('submits blank effective time, cutoff, total, and threshold %j with their optional meanings', async (minimum) => {
  const onSave = renderDialog();
  fireEvent.change(screen.getByLabelText('消费门槛（元）'), { target: { value: minimum } });
  expect((screen.getByLabelText('总发放量') as HTMLInputElement).value).toBe('');
  expect(screen.getByTestId('coupon-distribution-cutoff-row').className).not.toContain('grid-cols-2');
  submit();
  await waitFor(() => expect(onSave).toHaveBeenCalledWith(expect.objectContaining({
    minSpendFen: 0,
    effectiveAt: null,
    distributionEndsAt: null,
    totalIssueLimit: null,
  })));
});

it('converts valid local dates to ISO and accepts a later cutoff', async () => {
  const onSave = renderDialog();
  fireEvent.change(screen.getByLabelText('生效日期时间'), { target: { value: '2026-06-16T09:30' } });
  fireEvent.change(screen.getByLabelText('发放截止日期时间'), { target: { value: '2026-06-17T09:30' } });
  submit();
  await waitFor(() => expect(onSave).toHaveBeenCalledWith(expect.objectContaining({
    effectiveAt: new Date('2026-06-16T09:30').toISOString(),
    distributionEndsAt: new Date('2026-06-17T09:30').toISOString(),
  })));
});

it('rejects a cutoff that is not after effective time or now when effective time is blank', async () => {
  renderDialog();
  fireEvent.change(screen.getByLabelText('生效日期时间'), { target: { value: '2026-06-17T09:30' } });
  fireEvent.change(screen.getByLabelText('发放截止日期时间'), { target: { value: '2026-06-17T09:30' } });
  submit();
  expect(await screen.findByText('发放截止时间须晚于生效时间')).toBeTruthy();

  fireEvent.change(screen.getByLabelText('生效日期时间'), { target: { value: '' } });
  fireEvent.change(screen.getByLabelText('发放截止日期时间'), { target: { value: '2026-06-15T10:00' } });
  submit();
  expect(await screen.findByText('发放截止时间须晚于当前时间')).toBeTruthy();
});

describe('numeric boundaries', () => {
  it.each([
    ['0', '有效天数须为 1～365 天'],
    ['366', '有效天数须为 1～365 天'],
  ])('rejects days value %s', async (value, message) => {
    const onSave = renderDialog();
    fireEvent.change(screen.getByLabelText('激活后有效天数'), { target: { value } });
    submit();
    expect(await screen.findByText(message)).toBeTruthy();
    expect(onSave).not.toHaveBeenCalled();
  });

  it.each(['0', '-1', '1.5'])('rejects explicit total %s and reserves blank for unlimited', async (value) => {
    const onSave = renderDialog();
    fireEvent.change(screen.getByLabelText('总发放量'), { target: { value } });
    submit();
    expect(await screen.findByText('总发放量须为正整数；不限量请留空')).toBeTruthy();
    expect(onSave).not.toHaveBeenCalled();
  });

  it.each(['1', '365'])('accepts days boundary %s', async (value) => {
    const onSave = renderDialog();
    fireEvent.change(screen.getByLabelText('激活后有效天数'), { target: { value } });
    submit();
    await waitFor(() => expect(onSave).toHaveBeenCalledWith(expect.objectContaining({ daysAfterActivation: Number(value) })));
  });
});

it('uses the same callback contract for headquarters and store dialogs', async () => {
  const hqSave = renderDialog();
  submit();
  await waitFor(() => expect(hqSave).toHaveBeenCalled());
  cleanup();

  const storeSave = renderDialog(vi.fn().mockResolvedValue(undefined), 'store');
  submit();
  await waitFor(() => expect(storeSave).toHaveBeenCalled());

  expect(Object.keys(storeSave.mock.calls[0][0]).sort()).toEqual(Object.keys(hqSave.mock.calls[0][0]).sort());
});
