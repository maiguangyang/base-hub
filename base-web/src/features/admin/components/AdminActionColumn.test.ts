import { readdirSync, readFileSync, statSync } from 'node:fs';
import { join, relative } from 'node:path';
import { fileURLToPath } from 'node:url';
import { expect, it } from 'vitest';

const pages = fileURLToPath(new URL('../../', import.meta.url));
const styles = fileURLToPath(new URL('../../../styles/global.css', import.meta.url));

function tsxFiles(directory: string): string[] {
  return readdirSync(directory).flatMap((entry) => {
    const path = join(directory, entry);
    if (statSync(path).isDirectory()) return tsxFiles(path);
    return entry.endsWith('.tsx') && !entry.endsWith('.test.tsx') ? [path] : [];
  });
}

it('keeps every admin action column in the normal horizontal scroll flow', () => {
  const fixedHeaders: string[] = [];
  let actionColumns = 0;
  for (const file of tsxFiles(pages)) {
    const source = readFileSync(file, 'utf8');
    const headers = source.matchAll(/<(th|TableHead)\b([^>]*)>\s*操作\s*<\/\1>/g);
    for (const match of headers) {
      actionColumns += 1;
      if (/data-admin-action-column|\bsticky\b/.test(match[2])) fixedHeaders.push(relative(pages, file));
    }
  }
  expect(actionColumns).toBeGreaterThan(5);
  expect(fixedHeaders).toEqual([]);
  expect(readFileSync(styles, 'utf8')).not.toContain('data-admin-action-column');
});
