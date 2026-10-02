// @ts-check

import react from '@astrojs/react';
import tailwindcss from '@tailwindcss/vite';
import { defineConfig } from 'astro/config';
import checker from 'vite-plugin-checker';

export default defineConfig({
  integrations: [react()],
  i18n: {
    defaultLocale: 'zh-CN',
    locales: ['zh-CN', 'en-US', 'hi-IN', 'es-ES', 'ar-SA', 'bn-BD', 'pt-BR', 'ru-RU', 'ja-JP'],
    routing: { prefixDefaultLocale: false },
  },
  vite: {
    plugins: [tailwindcss(), checker({ typescript: true })],
  },
});
