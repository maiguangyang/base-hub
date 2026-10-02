import { describe, expect, it } from 'vitest';
import translationSource from '../locales/zh-CN.json';
import { entryPath, localeCodes, mergeMessages, messagesFor } from './index';

describe('supported locale routes', () => {
  it('includes exactly the source and eight active target locales', () => {
    expect(localeCodes).toEqual(['zh-CN', 'en-US', 'hi-IN', 'es-ES', 'ar-SA', 'bn-BD', 'pt-BR', 'ru-RU', 'ja-JP']);
  });

  it('keeps the current entry when switching language', () => {
    expect(entryPath('zh-CN', 'home')).toBe('/');
    expect(entryPath('ar-SA', 'home')).toBe('/ar-SA/');
  });

  it('后台只走简体中文，语言前缀不参与后台路径', () => {
    for (const locale of localeCodes) {
      expect(entryPath(locale, 'admin')).toBe('/admin');
    }
  });
});

describe('translation fallback', () => {
  it('uses the translation script source as the page source', () => {
    expect(messagesFor('zh-CN')).toEqual(translationSource);
  });

  it('keeps Chinese text when a target file is absent', () => {
    expect(messagesFor('en-US').publicHome.title).toBe(messagesFor('zh-CN').publicHome.title);
  });

  it('uses translated nested keys and Chinese for missing siblings', () => {
    const source = { common: { theme: { light: '浅色', dark: '深色' } } };
    const result = mergeMessages(source, { common: { theme: { dark: 'Dark' } } });
    expect(result).toEqual({ common: { theme: { light: '浅色', dark: 'Dark' } } });
  });
});
