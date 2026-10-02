import { renderToStaticMarkup } from 'react-dom/server';
import { describe, expect, it } from 'vitest';
import { Input } from './input';

describe('Input', () => {
  it('键盘聚焦时使用主题边框和外圈', () => {
    const markup = renderToStaticMarkup(<Input />);

    expect(markup).toContain('focus-visible:border-ring');
    expect(markup).toContain('focus-visible:ring-[3px]');
    expect(markup).toContain('focus-visible:ring-ring/50');
  });
});
