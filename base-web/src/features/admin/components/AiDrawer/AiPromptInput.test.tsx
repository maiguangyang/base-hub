// @vitest-environment jsdom

import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, expect, it, vi } from 'vitest';
import { useState, type ComponentProps } from 'react';
import { AiPromptInput } from './AiPromptInput';

afterEach(() => cleanup());

function PromptHarness(props: Omit<ComponentProps<typeof AiPromptInput>, 'prompt' | 'onPromptChange'>) {
  const [prompt, setPrompt] = useState('');
  return <AiPromptInput {...props} prompt={prompt} onPromptChange={setPrompt} />;
}

it('聊天输入区可选择图片附件并在发送后清空', async () => {
  const onSubmit = vi.fn();
  const onAttachment = vi.fn();
  render(<PromptHarness disabled={false} busy={false} requiredInputs={[]} onSubmit={onSubmit} onAttachment={onAttachment} onStop={vi.fn()} />);
  const file = new File(['image'], 'product.png', { type: 'image/png' });
  fireEvent.change(screen.getByLabelText('选择商品图片'), { target: { files: [file] } });
  expect(onAttachment).toHaveBeenCalledWith(file);
  const field = screen.getByRole('textbox', { name: '向 AI 助手发送消息' }) as HTMLTextAreaElement;
  fireEvent.change(field, { target: { value: '你好' } });
  fireEvent.click(screen.getByRole('button', { name: '发送消息' }));
  expect(onSubmit).toHaveBeenCalledWith('你好', undefined);
  await waitFor(() => expect(field.value).toBe(''));
});

it('空输入显示浅色不可发送按钮，处理中改为终止按钮', () => {
  const onStop = vi.fn();
  const props = { disabled: false, requiredInputs: [], onSubmit: vi.fn(), onAttachment: vi.fn(), onStop };
  const { rerender } = render(<PromptHarness {...props} busy={false} />);
  const send = screen.getByRole('button', { name: '发送消息' });
  expect(send.hasAttribute('disabled')).toBe(true);
  expect(send.className).toContain('bg-primary/40');
  rerender(<PromptHarness {...props} disabled busy />);
  expect(screen.queryByRole('button', { name: '发送消息' })).toBeNull();
  fireEvent.click(screen.getByRole('button', { name: '终止生成' }));
  expect(onStop).toHaveBeenCalledTimes(1);
});
