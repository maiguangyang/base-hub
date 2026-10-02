#!/usr/bin/env node

import { readFileSync, writeFileSync } from 'node:fs';
import { resolve } from 'node:path';

const query = process.argv[2];

if (!query) {
  throw new Error('缺少参数，请传入形如 "?baseUri=https%3A%2F%2Fexample.com" 的查询字符串');
}

const parameters = new Map(
  new URLSearchParams(query.startsWith('?') ? query.slice(1) : query),
);

if (parameters.size === 0) {
  throw new Error('查询字符串中没有可写入 Env 的参数');
}

const envPath = resolve(process.cwd(), 'src/config/env.ts');
const source = readFileSync(envPath, 'utf8');
const envObjectPattern = /const\s+Env\s*=\s*\{([\s\S]*?)\n\};/;
const envObjectMatch = envObjectPattern.exec(source);

if (!envObjectMatch) {
  throw new Error(`无法在 ${envPath} 中找到 const Env = { ... };`);
}

const formatKey = (key) =>
  /^[A-Za-z_$][\w$]*$/.test(key) ? key : JSON.stringify(key);

let body = envObjectMatch[1];

for (const [key, value] of parameters) {
  const escapedKey = key.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
  const existingPropertyPattern = new RegExp(
    `^(\\s*)(?:${escapedKey}|["']${escapedKey}["'])\\s*:[^\\n]*$`,
    'm',
  );
  const replacement = `$1${formatKey(key)}: ${JSON.stringify(value)},`;

  if (existingPropertyPattern.test(body)) {
    body = body.replace(existingPropertyPattern, replacement);
  } else {
    body = `${body.replace(/\s*$/, '')}\n  ${formatKey(key)}: ${JSON.stringify(value)},`;
  }
}

const updatedSource = source.replace(
  envObjectPattern,
  `const Env = {${body}\n};`,
);

writeFileSync(envPath, updatedSource, 'utf8');
console.log(`已更新 ${envPath}：${[...parameters.keys()].join(', ')}`);
