// @vitest-environment jsdom

import { cleanup, fireEvent, render, screen } from '@testing-library/react';
import { renderToStaticMarkup } from 'react-dom/server';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { HqStoreTable, type HqStoreRow } from './HqStoreTable';

const direct = {
  id: 'direct', code: 'D1', name: '直营一店', organizationId: 'hq',
  organization: { id: 'hq', name: '总部', type: 'HEADQUARTERS' },
  lifecycle: 'ACTIVE', businessStatus: 'OPEN', contactPhone: null,
  city: '上海', district: '徐汇', address: '漕溪路1号', businessHours: '09:00 - 22:00',
} as HqStoreRow;
const franchise = {
  ...direct, id: 'franchise', code: 'F1', name: '加盟一店', organizationId: 'f1',
  organization: { id: 'f1', name: '甲加盟商', type: 'FRANCHISE' },
} as HqStoreRow;

describe('统一门店列表操作边界', () => {
  afterEach(cleanup);
  it.each(['DRAFT', 'PENDING_APPROVAL', 'REJECTED', 'ACTIVE'] as const)('%s 门店的支付入口与启用条件一致', (lifecycle) => {
    const onConfigurePayment = vi.fn();
    const row = { ...franchise, lifecycle };
    render(<HqStoreTable rows={[row]} hqOrganizationId="hq" canUpdate canDelete statusBusy={false}
      onEdit={() => undefined} onDelete={() => undefined} onBusinessStatus={() => undefined} onConfigurePayment={onConfigurePayment} />);
    const button = screen.getByRole('button', { name: '支付配置' }) as HTMLButtonElement;
    expect(button.disabled).toBe(lifecycle !== 'ACTIVE');
    fireEvent.click(button);
    if (lifecycle === 'ACTIVE') expect(onConfigurePayment).toHaveBeenCalledWith(row);
    else expect(onConfigurePayment).not.toHaveBeenCalled();
  });
  it('一张表展示直营和加盟，并只向总部自己的门店提供编辑及营业开关', () => {
    const markup = renderToStaticMarkup(<HqStoreTable rows={[direct, franchise]} hqOrganizationId="hq"
      canUpdate canDelete statusBusy={false} onEdit={() => undefined} onDelete={() => undefined}
      onBusinessStatus={() => undefined} />);
    expect(markup).toContain('直营一店');
    expect(markup).toContain('加盟一店');
    expect(markup).toContain('甲加盟商');
    expect(markup).toContain('漕溪路1号');
    expect(markup).toContain('09:00 - 22:00');
    expect(markup.match(/>编辑</g)).toHaveLength(1);
    expect(markup.match(/>删除</g)).toHaveLength(1);
    expect(markup).toContain('aria-label="加盟一店营业状态"');
    expect(markup).toContain('aria-label="直营一店营业状态"');
  });

  it('加盟门店营业开关只读，直营门店可操作', () => {
    render(<HqStoreTable rows={[direct, franchise]} hqOrganizationId="hq" canUpdate canDelete
      statusBusy={false} onEdit={() => undefined} onDelete={() => undefined} onBusinessStatus={() => undefined} />);
    expect((screen.getByRole('switch', { name: '直营一店营业状态' }) as HTMLButtonElement).disabled).toBe(false);
    expect((screen.getByRole('switch', { name: '加盟一店营业状态' }) as HTMLButtonElement).disabled).toBe(true);
  });

  it('其他总部组织的门店不出现直营操作或加盟支付配置', () => {
    const otherHq = { ...direct, id: 'other-hq', organizationId: 'other',
      organization: { id: 'other', name: '其他总部', type: 'HEADQUARTERS' } } as HqStoreRow;
    const markup = renderToStaticMarkup(<HqStoreTable rows={[otherHq]} hqOrganizationId="hq" canUpdate canDelete
      statusBusy={false} onEdit={() => undefined} onDelete={() => undefined} onBusinessStatus={() => undefined}
      onConfigurePayment={() => undefined} />);
    expect(markup).not.toContain('>编辑<');
    expect(markup).not.toContain('>删除<');
    expect(markup).not.toContain('>支付配置<');
  });
});
