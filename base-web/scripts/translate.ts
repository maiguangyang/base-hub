#!/usr/bin/env node
/**
 * AI 辅助翻译脚本
 * 
 * 功能：
 * - 读取 zh-CN.json 作为翻译源
 * - 自动检测目标语言文件中缺失的键
 * - 调用本地 AI API 进行翻译
 * - 增量翻译（不覆盖已存在的翻译）
 * 
 * 使用方式：
 * - pnpm translate          # 翻译所有目标语言
 * - pnpm translate --lang en  # 仅翻译英语
 */


import * as fs from 'fs';
import * as path from 'path';
import * as crypto from 'crypto';

// ============================================================================
// 配置
// ============================================================================

const AI_API_URL = process.env.AI_API_URL || 'http://localhost:1234/v1/chat/completions';
const AI_MODEL = process.env.AI_MODEL || 'HY-MT1.5-7B-4bit';

// 源语言文件路径
const LOCALES_DIR = path.join(process.cwd(), 'src/locales');
const SOURCE_FILE = path.join(LOCALES_DIR, 'zh-CN.json');
const CACHE_FILE = path.join(LOCALES_DIR, '.translation-cache.json');

// 目标语言配置（按使用热度排序）
// 文件名格式与 locales 目录一致（如 en-US.json, es-ES.json）
const TARGET_LOCALES = [
  // 第一梯队：全球主要语言（10亿+用户）
  'en-US',    // 英语 - 全球通用
  'hi-IN',    // 印地语 - 6亿+用户
  'es-ES',    // 西班牙语 - 5亿+用户
  'ar-SA',    // 阿拉伯语 - 4亿+用户
  'bn-BD',    // 孟加拉语 - 3亿+用户
  'pt-BR',    // 葡萄牙语 - 2.5亿+用户
  'ru-RU',    // 俄语 - 2.5亿+用户
  'ja-JP',    // 日语 - 1.2亿+用户

  // // 第二梯队：亚太地区热门语言
  // 'ko-KR',    // 韩语 - 8千万用户
  // 'vi-VN',    // 越南语 - 8千万用户
  // 'th-TH',    // 泰语 - 6千万用户
  // 'id-ID',    // 印尼语 - 4千万用户（2.7亿人口）
  // 'ms-MY',    // 马来语 - 3千万用户
  // 'tl-PH',    // 菲律宾语 - 2.8千万用户
  // 'zh-TW',    // 繁体中文 - 台港澳

  // // 第三梯队：欧洲主要语言
  // 'de-DE',    // 德语 - 1亿+用户
  // 'fr-FR',    // 法语 - 2.8亿用户
  // 'it-IT',    // 意大利语 - 6千万用户
  // 'pl-PL',    // 波兰语 - 4千万用户
  // 'nl-NL',    // 荷兰语 - 2.5千万用户
  // 'uk-UA',    // 乌克兰语 - 4千万用户
  // 'tr-TR',    // 土耳其语 - 8千万用户
  // 'cs-CZ',    // 捷克语 - 1千万用户

  // // 第四梯队：南亚及中东语言
  // 'ur-PK',    // 乌尔都语 - 2.3亿用户
  // 'ta-IN',    // 泰米尔语 - 8千万用户
  // 'te-IN',    // 泰卢固语 - 8千万用户
  // 'mr-IN',    // 马拉地语 - 8千万用户
  // 'gu-IN',    // 古吉拉特语 - 5千万用户
  // 'fa-IR',    // 波斯语 - 1.1亿用户
  // 'he-IL',    // 希伯来语 - 9百万用户

  // // 第五梯队：东南亚语言
  // 'km-KH',    // 高棉语 - 1.6千万用户
  // 'my-MM',    // 缅甸语 - 3.3千万用户
];

// 语言代码 → Chinese Names 映射（用于 AI 提示词）
const LANGUAGE_MAP: Record<string, string> = {
  'zh-CN': '中文',
  'en-US': '英语',
  'es-ES': '西班牙语',
  'fr-FR': '法语',
  'pt-BR': '葡萄牙语',
  'pt-PT': '葡萄牙语',
  'ja-JP': '日语',
  'ko-KR': '韩语',
  'zh-TW': '繁体中文',
  'zh-HK': '粤语',
  'de-DE': '德语',
  'it-IT': '意大利语',
  'ru-RU': '俄语',
  'ar-SA': '阿拉伯语',
  'th-TH': '泰语',
  'vi-VN': '越南语',
  'id-ID': '印尼语',
  'ms-MY': '马来语',
  'tr-TR': '土耳其语',
  'pl-PL': '波兰语',
  'nl-NL': '荷兰语',
  'uk-UA': '乌克兰语',
  'hi-IN': '印地语',
  'cs-CZ': '捷克语',
  'km-KH': '高棉语',
  'my-MM': '缅甸语',
  'fa-IR': '波斯语',
  'gu-IN': '古吉拉特语',
  'ur-PK': '乌尔都语',
  'te-IN': '泰卢固语',
  'mr-IN': '马拉地语',
  'he-IL': '希伯来语',
  'bn-BD': '孟加拉语',
  'ta-IN': '泰米尔语',
  'tl-PH': '菲律宾语',
  'bo-CN': '藏语',
  'kk-KZ': '哈萨克语',
  'mn-MN': '蒙古语',
  'ug-CN': '维吾尔语',
};

// ============================================================================
// JSON 工具函数
// ============================================================================

/**
 * 将嵌套 JSON 展平为 path -> value 映射
 */
function flattenJSON(obj: Record<string, unknown>, prefix = ''): Record<string, string> {
  const result: Record<string, string> = {};

  for (const key in obj) {
    const newKey = prefix ? `${prefix}.${key}` : key;
    const value = obj[key];

    if (typeof value === 'object' && value !== null && !Array.isArray(value)) {
      Object.assign(result, flattenJSON(value as Record<string, unknown>, newKey));
    } else if (typeof value === 'string') {
      result[newKey] = value;
    }
  }

  return result;
}

// ============================================================================
// 哈希缓存工具
// ============================================================================

// 按语言存储缓存，每个语言有自己的哈希记录
interface LanguageCache {
  hashes: Record<string, string>;  // key -> hash of source value
  updatedAt: string;
}

interface TranslationCache {
  // locale code -> language cache
  languages: Record<string, LanguageCache>;
}

/**
 * 计算字符串的 MD5 哈希值
 */
function hashString(str: string): string {
  return crypto.createHash('md5').update(str).digest('hex');
}

/**
 * 加载翻译缓存
 */
function loadCache(): TranslationCache {
  if (fs.existsSync(CACHE_FILE)) {
    try {
      const data = JSON.parse(fs.readFileSync(CACHE_FILE, 'utf-8'));
      // 兼容旧格式
      if (data.hashes && !data.languages) {
        return { languages: {} };
      }
      return data as TranslationCache;
    } catch {
      // 缓存文件损坏，返回空缓存
    }
  }
  return { languages: {} };
}

/**
 * 保存翻译缓存
 */
function saveCache(cache: TranslationCache): void {
  fs.writeFileSync(CACHE_FILE, JSON.stringify(cache, null, 2) + '\n', 'utf-8');
}

/**
 * 获取指定语言的缓存哈希
 */
function getLanguageHashes(cache: TranslationCache, locale: string): Record<string, string> {
  return cache.languages[locale]?.hashes || {};
}

/**
 * 更新指定语言的缓存
 */
function updateLanguageCache(cache: TranslationCache, locale: string, hashes: Record<string, string>): void {
  cache.languages[locale] = {
    hashes,
    updatedAt: new Date().toISOString()
  };
}

/**
 * 根据源文件生成键的哈希映射
 */
function generateSourceHashes(sourceFlat: Record<string, string>): Record<string, string> {
  const hashes: Record<string, string> = {};
  for (const key in sourceFlat) {
    hashes[key] = hashString(sourceFlat[key]);
  }
  return hashes;
}

/**
 * 找出需要翻译的键（新增 + 修改）
 */
function findKeysToTranslate(
  sourceFlat: Record<string, string>,
  targetFlat: Record<string, string>,
  cachedHashes: Record<string, string>,
  currentHashes: Record<string, string>
): { newKeys: string[]; changedKeys: string[] } {
  const newKeys: string[] = [];
  const changedKeys: string[] = [];

  for (const key in sourceFlat) {
    if (!(key in targetFlat)) {
      // 目标文件中不存在的键 → 新增
      newKeys.push(key);
    } else {
      const cached = cachedHashes[key];
      if (cached !== undefined && cached !== currentHashes[key]) {
        // 有缓存记录且哈希不一致 → 源文件已修改 → 需要重新翻译
        changedKeys.push(key);
      }
      // cached === undefined：未曾被跟踪 → 保守跳过（用 --force 可强制重翻）
      // cached === currentHashes[key]：哈希一致 → 无需翻译
    }
  }

  return { newKeys, changedKeys };
}


/**
 * 将展平的映射重建为嵌套 JSON
 */
function unflattenJSON(flat: Record<string, string>): Record<string, unknown> {
  const result: Record<string, unknown> = {};

  for (const key in flat) {
    const parts = key.split('.');
    let current = result;

    for (let i = 0; i < parts.length - 1; i++) {
      const part = parts[i];
      if (!(part in current)) {
        current[part] = {};
      }
      current = current[part] as Record<string, unknown>;
    }

    current[parts[parts.length - 1]] = flat[key];
  }

  return result;
}

/**
 * 深度合并两个对象（target 中已存在的值不会被覆盖）
 */
function deepMerge(target: Record<string, unknown>, source: Record<string, unknown>): Record<string, unknown> {
  const result = { ...target };

  for (const key in source) {
    if (key in result) {
      if (typeof result[key] === 'object' && typeof source[key] === 'object' &&
        result[key] !== null && source[key] !== null &&
        !Array.isArray(result[key]) && !Array.isArray(source[key])) {
        result[key] = deepMerge(result[key] as Record<string, unknown>, source[key] as Record<string, unknown>);
      }
      // 已存在的值不覆盖
    } else {
      result[key] = source[key];
    }
  }

  return result;
}

// ============================================================================
// AI 翻译
// ============================================================================

/**
 * 调用 AI API 进行翻译
 */
async function callAI(text: string, targetLang: string): Promise<string> {
  const chineseName = LANGUAGE_MAP[targetLang] || targetLang;
  const systemPrompt = `将以下文本翻译为${chineseName}，注意只需要输出翻译后的结果，不要额外解释：`;

  try {
    const response = await fetch(AI_API_URL, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer sk-1234567890'
      },
      body: JSON.stringify({
        model: AI_MODEL,
        messages: [
          { role: 'system', content: systemPrompt },
          { role: 'user', content: text },
        ],
        temperature: 0.3,
      }),
    });

    if (!response.ok) {
      throw new Error(`API 请求失败: ${response.status} ${response.statusText}`);
    }

    const data = await response.json() as {
      choices: Array<{ message: { content: string } }>;
    };

    const result = data.choices?.[0]?.message?.content?.trim();
    if (!result) {
      throw new Error('AI 返回空结果');
    }

    return result;
  } catch (error) {
    if (error instanceof Error) {
      throw new Error(`翻译失败: ${error.message}`);
    }
    throw error;
  }
}



// ============================================================================
// 主流程
// ============================================================================

/**
 * 翻译单个语言文件（使用哈希检测变化）
 */
async function translateFile(
  targetLang: string,
  cachedHashes: Record<string, string>,
  currentHashes: Record<string, string>,
  forceAll = false
): Promise<{ added: number; updated: number; skipped: number; translatedHashes: Record<string, string> }> {
  const targetFile = path.join(LOCALES_DIR, `${targetLang}.json`);

  // 读取源文件
  if (!fs.existsSync(SOURCE_FILE)) {
    throw new Error(`源文件不存在: ${SOURCE_FILE}`);
  }
  const sourceData = JSON.parse(fs.readFileSync(SOURCE_FILE, 'utf-8')) as Record<string, unknown>;
  const sourceFlat = flattenJSON(sourceData);

  // 读取目标文件（如果存在）
  let targetData: Record<string, unknown> = {};
  let targetFlat: Record<string, string> = {};
  if (fs.existsSync(targetFile)) {
    targetData = JSON.parse(fs.readFileSync(targetFile, 'utf-8')) as Record<string, unknown>;
    targetFlat = flattenJSON(targetData);
  }

  // 确定需要翻译的键
  let newKeys: string[];
  let changedKeys: string[];

  if (forceAll) {
    newKeys = Object.keys(sourceFlat).filter(k => !(k in targetFlat));
    changedKeys = Object.keys(sourceFlat).filter(k => k in targetFlat);
  } else {
    const result = findKeysToTranslate(sourceFlat, targetFlat, cachedHashes, currentHashes);
    newKeys = result.newKeys;
    changedKeys = result.changedKeys;
  }

  const keysToTranslate = [...newKeys, ...changedKeys];

  if (keysToTranslate.length === 0) {
    console.log(`  ✓ ${targetLang}: 无需翻译`);
    return { added: 0, updated: 0, skipped: Object.keys(targetFlat).length, translatedHashes: {} };
  }

  const actionDesc = changedKeys.length > 0
    ? `翻译 ${newKeys.length} 新增 + ${changedKeys.length} 修改`
    : `翻译 ${newKeys.length} 新增`;
  console.log(`  ${targetLang}: ${actionDesc}...`);

  // 翻译键，记录成功翻译的 key 的哈希（用于选择性更新缓存）
  const newTranslations: Record<string, string> = {};
  const translatedHashes: Record<string, string> = {};
  for (const key of keysToTranslate) {
    const sourceText = sourceFlat[key];
    const isChanged = changedKeys.includes(key);
    try {
      const translated = await callAI(sourceText, targetLang);
      newTranslations[key] = translated;
      translatedHashes[key] = currentHashes[key]; // 仅记录成功翻译 key 的最新哈希
      const prefix = isChanged ? '🔄' : '✓';
      console.log(`    ${prefix} ${key}: ${sourceText} → ${translated}`);
    } catch (error) {
      console.error(`    ✗ ${key}: 翻译失败 - ${error instanceof Error ? error.message : error}`);
      // 使用原文作为回退，但不更新缓存（下次还会尝试翻译）
      newTranslations[key] = sourceText;
    }
  }

  // 合并并写入目标文件
  const newFlat = { ...targetFlat, ...newTranslations };
  const newData = unflattenJSON(newFlat);
  const mergedData = deepMerge(targetData, newData);

  fs.writeFileSync(targetFile, JSON.stringify(mergedData, null, 2) + '\n', 'utf-8');

  return {
    added: newKeys.length,
    updated: changedKeys.length,
    skipped: Object.keys(targetFlat).length - changedKeys.length,
    translatedHashes,
  };
}

/**
 * 解析命令行参数
 */
function parseArgs(): { lang?: string; help: boolean; force: boolean } {
  const args = process.argv.slice(2);
  const result: { lang?: string; help: boolean; force: boolean } = { help: false, force: false };

  for (let i = 0; i < args.length; i++) {
    if (args[i] === '--lang' && args[i + 1]) {
      result.lang = args[i + 1];
      i++;
    } else if (args[i] === '--help' || args[i] === '-h') {
      result.help = true;
    } else if (args[i] === '--force' || args[i] === '-f') {
      result.force = true;
    }
  }

  return result;
}

/**
 * 显示帮助信息
 */
function showHelp(): void {
  console.log(`
AI 辅助翻译脚本

使用方式:
  pnpm translate              翻译所有目标语言（仅缺失的键）
  pnpm translate --force      强制重新翻译所有键
  pnpm translate --lang en-US 仅翻译指定语言

选项:
  --lang <code>  指定目标语言代码 (如 en-US, es-ES, ja-JP)
  --force, -f    强制重新翻译所有键（包括已存在的）
  --help, -h     显示帮助信息

配置:
  TARGET_LOCALES: [${TARGET_LOCALES.slice(0, 5).join(', ')}, ... (${TARGET_LOCALES.length} 种语言)]
  AI_API_URL: ${AI_API_URL}
`);
}

// ============================================================================
// 自动更新 locales/index.ts
// ============================================================================

const INDEX_FILE = path.join(LOCALES_DIR, 'index.ts');

/**
 * 将 locale code 转换为变量名（如 en-US → enUS）
 */
function localeToVarName(locale: string): string {
  return locale.replace(/-/g, '');
}

/**
 * 自动更新 locales/index.ts，添加新语言的导入和配置
 */
function updateLocalesIndex(): void {
  // 扫描 locales 目录下的所有 JSON 文件（排除缓存文件）
  const files = fs.readdirSync(LOCALES_DIR)
    .filter(f => f.endsWith('.json') && !f.startsWith('.'));  // 排除以.开头的文件
  const locales = files.map(f => f.replace('.json', '')).sort();

  // 生成 fallback 映射
  const fallbackMap: Record<string, string[]> = {
    'default': ['zh-CN']
  };

  // 自动根据基础语言生成映射 (e.g. 'zh' -> 'zh-CN', 'es' -> 'es-ES')
  // 只有当基础语言不存在时才映射 (防止覆盖如 'pt' 和 'pt-BR' 共存的情况)
  for (const locale of locales) {
    const parts = locale.split('-');
    if (parts.length === 2) {
      const baseLang = parts[0];
      // 如果基础语言本身不是一个支持的 locale，则建立映射
      if (!locales.includes(baseLang)) {
        // 如果还没有映射，或者当前 locale 是该语言的第一个变体，则添加
        // 这里简化处理：直接映射到第一个遇到的变体
        // 改进：优先映射到同名的大语种变体 (如 zh -> zh-CN, es -> es-ES)
        if (!fallbackMap[baseLang]) {
          fallbackMap[baseLang] = [locale];
        }
      }
    }
  }

  // 生成导入语句
  const imports = locales.map(locale => {
    const varName = localeToVarName(locale);
    return `import ${varName} from './${locale}.json';`;
  }).join('\n');

  // 生成 resources 配置
  const resources = locales.map(locale => {
    const varName = localeToVarName(locale);
    return `  '${locale}': { translation: ${varName} },`;
  }).join('\n');

  // 生成 availableLocales
  const availableLocales = locales.map(locale => {
    const name = LANGUAGE_MAP[locale] || locale;

    // 生成原生名称
    let nativeName = name;
    try {
      // 使用 Intl.DisplayNames 获取原生语言名称
      // e.g. new Intl.DisplayNames(['de-DE'], { type: 'language' }).of('de-DE') -> "Deutsch"
      const dn = new Intl.DisplayNames([locale], { type: 'language' });
      const n = dn.of(locale);
      if (n) {
        nativeName = n;
        // 首字母大写优化 (对于某些语言可能需要)
        // nativeName = nativeName.charAt(0).toUpperCase() + nativeName.slice(1);
      }
    } catch (e) {
      console.warn(`  ⚠️ 无法获取 ${locale} 的原生名称，回退到中文名`);
    }

    return `  { code: '${locale}', name: '${name}', nativeName: '${nativeName}' },`;
  }).join('\n');

  // 生成完整的 index.ts 内容
  const content = `/**
 * i18next 初始化配置
 * 
 * 集中式语言管理：所有翻译文件统一存放在此目录
 * ⚠️ 此文件由 translate.ts 脚本自动生成，请勿手动编辑
 */

import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import LanguageDetector from 'i18next-browser-languagedetector';
import { LS_I18N_LOCALE } from '@/lib/storage-keys';

// 导入翻译文件 (自动生成)
${imports}

// 语言资源配置 (自动生成)
export const resources = {
${resources}
};

// 可用语言列表 (自动生成)
export const availableLocales = [
${availableLocales}
] as const;

export type Locale = typeof availableLocales[number]['code'];

// 默认语言
export const defaultLocale: Locale = 'zh-CN';

// 回退策略 (自动生成)
export const fallbackLng = ${JSON.stringify(fallbackMap, null, 2)};

// i18next 初始化（guard: 避免 Vite HMR 重复 init 导致语言状态错误）
if (!i18n.isInitialized) {
  i18n
    .use(LanguageDetector)
    .use(initReactI18next)
    .init({
      resources,
      fallbackLng,
      supportedLngs: availableLocales.map(l => l.code),

      // 只加载精确匹配的语言码（避免 i18next 尝试加载 'en' / 'en-us' 等父级变体导致 reject 警告）
      load: 'currentOnly',

      // 检测顺序：LocalStorage → 浏览器语言 → 默认语言
      detection: {
        order: ['localStorage', 'navigator'],
        lookupLocalStorage: LS_I18N_LOCALE,
        caches: ['localStorage'],
      },

      interpolation: {
        escapeValue: false, // React 已有 XSS 防护
      },
    });
}

export default i18n;
`;

  fs.writeFileSync(INDEX_FILE, content, 'utf-8');
  console.log(`  ✓ 已更新 ${INDEX_FILE}`);
}

/**
 * 清理不在 TARGET_LOCALES 中的语言文件和缓存
 */
function cleanupOrphanedLocales(): number {
  const sourceLocale = 'zh-CN'; // 源语言不能删除
  const allowedLocales = new Set([sourceLocale, ...TARGET_LOCALES]);

  // 扫描 locales 目录下的所有 JSON 文件（排除缓存文件）
  const files = fs.readdirSync(LOCALES_DIR)
    .filter(f => f.endsWith('.json') && !f.startsWith('.'));  // 排除以.开头的文件
  const orphanedFiles: string[] = [];

  for (const file of files) {
    const locale = file.replace('.json', '');
    if (!allowedLocales.has(locale)) {
      orphanedFiles.push(file);
    }
  }

  // 清理文件
  if (orphanedFiles.length > 0) {
    console.log(`  🗑️ 清理 ${orphanedFiles.length} 个不再需要的语言文件...`);
    for (const file of orphanedFiles) {
      const filePath = path.join(LOCALES_DIR, file);
      fs.unlinkSync(filePath);
      console.log(`    ✓ 已删除 ${file}`);
    }
  }

  // 清理缓存中废弃的语言条目
  const cache = loadCache();
  const orphanedCacheLocales = Object.keys(cache.languages).filter(
    locale => !allowedLocales.has(locale)
  );

  if (orphanedCacheLocales.length > 0) {
    console.log(`  🧹 清理 ${orphanedCacheLocales.length} 个缓存条目...`);
    for (const locale of orphanedCacheLocales) {
      delete cache.languages[locale];
      console.log(`    ✓ 已清理缓存: ${locale}`);
    }
    saveCache(cache);
  }

  return orphanedFiles.length + orphanedCacheLocales.length;
}

/**
 * 主函数
 */
async function main(): Promise<void> {
  const args = parseArgs();

  if (args.help) {
    showHelp();
    return;
  }

  console.log('🌐 AI 辅助翻译脚本');
  console.log(`   源文件: ${SOURCE_FILE}`);
  console.log(`   API: ${AI_API_URL}`);
  if (args.force) {
    console.log('   模式: 强制重新翻译所有键');
  } else {
    console.log('   模式: 智能检测（新增 + 修改）');
  }
  console.log('');

  // 加载源文件并计算哈希
  const sourceData = JSON.parse(fs.readFileSync(SOURCE_FILE, 'utf-8')) as Record<string, unknown>;
  const sourceFlat = flattenJSON(sourceData);
  const currentHashes = generateSourceHashes(sourceFlat);

  // 加载缓存
  const cache = loadCache();

  const locales = args.lang ? [args.lang] : TARGET_LOCALES;

  let totalAdded = 0;
  let totalUpdated = 0;
  let totalSkipped = 0;

  for (const locale of locales) {
    if (!LANGUAGE_MAP[locale]) {
      console.warn(`  ⚠ 未知语言代码: ${locale}`);
      continue;
    }

    // 获取该语言的缓存哈希
    const cachedHashes = getLanguageHashes(cache, locale);

    try {
      const { added, updated, skipped, translatedHashes } = await translateFile(locale, cachedHashes, currentHashes, args.force);
      totalAdded += added;
      totalUpdated += updated;
      totalSkipped += skipped;

      // 仅将成功翻译的 key 的哈希合并入缓存（未翻译的 key 保留旧哈希，确保下次仍能检测到变化）
      if (Object.keys(translatedHashes).length > 0) {
        const existing = getLanguageHashes(cache, locale);
        updateLanguageCache(cache, locale, { ...existing, ...translatedHashes });
        saveCache(cache);
      }
    } catch (error) {
      console.error(`  ✗ ${locale}: ${error instanceof Error ? error.message : error}`);
      // 翻译失败时不更新该语言的缓存
    }
  }

  console.log('');
  console.log(`📊 统计: 新增 ${totalAdded} 个，更新 ${totalUpdated} 个，跳过 ${totalSkipped} 个`);

  // 清理不在 TARGET_LOCALES 中的语言文件
  console.log('');
  console.log('🧹 检查废弃的语言文件...');
  const deletedCount = cleanupOrphanedLocales();
  if (deletedCount === 0) {
    console.log('  ✓ 无需清理');
  }

  // 自动更新 locales/index.ts
  console.log('');
  console.log('🔄 更新语言配置...');
  updateLocalesIndex();
}

// 运行主函数
main().catch((error) => {
  console.error('❌ 脚本执行失败:', error);
  process.exit(1);
});

