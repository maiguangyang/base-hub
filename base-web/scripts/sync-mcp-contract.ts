import { copyFileSync, existsSync, mkdirSync } from 'node:fs';
import { dirname, resolve } from 'node:path';

const sourcePath = resolve(process.cwd(), '../mcp-prerequisite-candidate-contract-v1.md');
const destinationPath = resolve(
  process.cwd(),
  'src/features/public-home/content/mcp-integration/zh-CN.md',
);

if (!existsSync(sourcePath)) {
  throw new Error(`MCP contract source not found: ${sourcePath}`);
}

mkdirSync(dirname(destinationPath), { recursive: true });
copyFileSync(sourcePath, destinationPath);

console.log(`Synced MCP contract snapshot to ${destinationPath}`);
