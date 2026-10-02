// @vitest-environment jsdom

import { cleanup, render, screen } from '@testing-library/react';
import { afterEach, describe, expect, it } from 'vitest';
import { Input } from '@/components/ui/input';
import { Select, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Textarea } from '@/components/ui/textarea';
import { TimePicker, TimeRangePicker } from '@/components/ui/time-picker';
import { AdminActionError, AdminFormDialogShell, runAdminAction } from './AdminFormDialogShell';

afterEach(cleanup);

describe('AdminFormDialogShell', () => {
  it('保留结构化错误对应的原因，不显示内部异常原文', async () => {
    let message: string | undefined;
    await runAdminAction(async () => { throw { errors: [{ extensions: { code: 'PERMISSION_DENIED' } }] }; }, (next) => { message = next; });
    expect(message).toContain('没有执行此操作的权限');
    await runAdminAction(async () => { throw new Error('SQL password=secret'); }, (next) => { message = next; });
    expect(message).not.toContain('secret');
  });
  it('将 rejected submit 转换为用户可见错误且不继续抛出', async () => {
    let message: string | undefined;
    await expect(runAdminAction(
      async () => { throw new Error('network'); },
      (next) => { message = next; },
    )).resolves.toBe(false);
    expect(message).toBe('操作失败：操作未完成，请稍后重试；如涉及新增或发放，请先核对结果。');
  });

  it('执行新动作前清除旧错误，成功时返回 true', async () => {
    const messages: Array<string | undefined> = ['旧错误'];
    await expect(runAdminAction(async () => undefined, (next) => messages.push(next))).resolves.toBe(true);
    expect(messages).toEqual(['旧错误', undefined]);
  });

  it('以 alert 呈现行操作错误', () => {
    const element = AdminActionError({ message: '操作失败，请稍后重试。' });
    expect(element?.props.role).toBe('alert');
    expect(element?.props.children).toBe('操作失败，请稍后重试。');
    expect(AdminActionError({})).toBeNull();
  });
});

describe('AdminFormDialogShell 视觉结构', () => {
  it('弹窗内共享输入控件统一暴露后台表单表面语义', () => {
    render(
      <AdminFormDialogShell open onOpenChange={() => undefined} title="表单" submitLabel="保存" cancelLabel="取消" onSubmit={() => undefined}>
        <Input aria-label="文本框" />
        <Textarea aria-label="多行文本框" />
        <Select><SelectTrigger aria-label="选择框"><SelectValue placeholder="请选择" /></SelectTrigger></Select>
        <TimePicker />
        <TimeRangePicker />
      </AdminFormDialogShell>,
    );

    const surfaces = document.querySelectorAll('[data-admin-form-surface]');
    expect(surfaces).toHaveLength(5);
    expect(screen.getByRole('dialog').className).not.toContain('[&_input]:!bg-white');
  });

  it('固定头尾并让表单主体独立滚动', () => {
    render(
      <AdminFormDialogShell open onOpenChange={() => undefined} title="表单" submitLabel="保存" cancelLabel="取消" onSubmit={() => undefined}>
        <Input aria-label="文本框" />
      </AdminFormDialogShell>,
    );

    expect(screen.getByRole('dialog').className).toContain('overflow-hidden');
    expect(screen.getByText('表单').closest('header')?.className).toContain('shrink-0');
    expect(screen.getByTestId('admin-form-dialog-body').className).toContain('overflow-y-auto');
    expect(screen.getByTestId('admin-form-dialog-footer').className).toContain('shrink-0');
    expect(screen.getByRole('dialog').className).toContain('max-w-2xl');
  });
});
