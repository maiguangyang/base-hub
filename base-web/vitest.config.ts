import { getViteConfig } from 'astro/config';
import { configDefaults } from 'vitest/config';

// 用 Astro 的 getViteConfig 而非裸 defineConfig：它继承 astro.config.mjs 与
// tsconfig 的 paths 解析，使测试中的 @/ 别名与构建时保持同一套规则，
// 避免出现"构建通过但测试解析不到模块"的分裂。
export default getViteConfig({
  test: {
    // 同时覆盖 .tsx：组件测试若被静默跳过，不报错也不提示，最难察觉。
    // scripts/ 下的契约检查脚本同样需要回归防护，一并纳入。
    include: ['src/**/*.test.{ts,tsx}', 'scripts/**/*.test.mjs'],
    // 这两个契约使用 node:test；由完整 test 入口交给 Node 执行。
    exclude: [...configDefaults.exclude, 'scripts/checkbox-contract.test.mjs', 'scripts/select-placeholder-contract.test.mjs'],
  },
});
