import sourceMessages from '../locales/zh-CN.json';

/** 官网与后台共同提供的语言；与现有翻译脚本启用列表保持一致。 */
export const localeCodes = ['zh-CN', 'en-US', 'hi-IN', 'es-ES', 'ar-SA', 'bn-BD', 'pt-BR', 'ru-RU', 'ja-JP'] as const;

/** URL 中允许的页面语言。 */
export type Locale = (typeof localeCodes)[number];

/** 两个入口共享 URL 构造规则，切换语言时保持当前入口。 */
export type Entry = 'home' | 'admin';

/** 下拉菜单显示语言的自称，译文缺失时仍能辨认目标语言。 */
export const localeNames: Record<Locale, string> = {
  'zh-CN': '简体中文',
  'en-US': 'English',
  'hi-IN': 'हिन्दी',
  'es-ES': 'Español',
  'ar-SA': 'العربية',
  'bn-BD': 'বাংলা',
  'pt-BR': 'Português',
  'ru-RU': 'Русский',
  'ja-JP': '日本語',
};

/** 页面文案始终具有完整的中文源文件形状。 */
export type Messages = typeof sourceMessages;

/** 构建时识别同目录中的目标语言 JSON；新增文件无需修改入口代码。 */
const localeFiles = import.meta.glob('../locales/*.json', { eager: true, import: 'default' });

/** 判断译文节点能否继续按 key 递归合并。 */
function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === 'object' && value !== null && !Array.isArray(value);
}

/** 只覆盖已有源 key 的字符串译文；缺失或非字符串目标值保留中文。 */
export function mergeMessages<T extends object>(source: T, target: unknown): T {
  if (!isRecord(target)) return source;
  const merged: Record<string, unknown> = {};
  for (const [key, value] of Object.entries(source)) {
    const translated = target[key];
    merged[key] = isRecord(value)
      ? mergeMessages(value, translated)
      : typeof value === 'string' && typeof translated === 'string' && translated.length > 0
        ? translated
        : value;
  }
  return merged as T;
}

/** 读取当前 URL 语言文案；文件或 key 缺失时回退简体中文。 */
export function messagesFor(locale: Locale): Messages {
  if (locale === 'zh-CN') return sourceMessages;
  return mergeMessages(sourceMessages, localeFiles[`../locales/${locale}.json`]);
}

/** 生成官网或后台路径。后台只走简体中文，不接受语言前缀；官网保留 URL 语言。 */
export function entryPath(locale: Locale, entry: Entry): string {
  if (entry === 'admin') return '/admin';
  return locale === 'zh-CN' ? '/' : `/${locale}/`;
}
