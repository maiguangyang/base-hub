import astro from 'eslint-plugin-astro';
import tseslint from 'typescript-eslint';

export default [
  { ignores: ['dist/**', '.astro/**', 'node_modules/**', 'src/__generated__/**'] },
  ...tseslint.configs.recommended,
  ...astro.configs['flat/recommended'],
  {
    files: ['**/*.{ts,tsx,astro}'],
    rules: {
      'max-lines': ['error', { max: 250, skipBlankLines: true, skipComments: true }],
      'max-lines-per-function': ['error', { max: 50, skipBlankLines: true, skipComments: true }],
      complexity: ['error', 10],
      'max-depth': ['error', 3],
      // WEB-UI-013：占位文案只负责控件提示，不能渲染为可选或禁用选项。
      'no-restricted-syntax': ['error', {
        selector: 'JSXOpeningElement[name.name="input"]:has(> JSXAttribute[name.name="type"][value.value="checkbox"])',
        message: 'WEB-UI-014：复选框必须使用 @/components/ui/checkbox，未勾选统一白底；禁止原生 input checkbox。',
      }, {
        selector: 'JSXElement:has(> JSXOpeningElement[name.name=/^(SelectItem|option)$/]) > JSXExpressionContainer > :matches(Identifier[name="placeholder"], MemberExpression[property.name="placeholder"])',
        message: 'placeholder 只能用于未选择提示，不能作为下拉选项；业务空值选项请独立声明 clearLabel。',
      }, {
        selector: 'JSXElement:has(> JSXOpeningElement:has(> JSXAttribute[name.name="role"][value.value="option"])) > JSXExpressionContainer > :matches(Identifier[name="placeholder"], MemberExpression[property.name="placeholder"])',
        message: 'placeholder 不能作为 role="option" 选项；请只渲染真实业务项或显式 clearLabel。',
      }],
    },
  },
  // WEB-UI-014：业务页面只能使用统一的白底复选框组件。
  {
    files: ['src/features/**/*.{ts,tsx}'],
    rules: {
      'no-restricted-imports': ['error', {
        paths: [{ name: 'radix-ui', importNames: ['Checkbox'], message: '请使用 @/components/ui/checkbox，确保未勾选复选框统一白底。' },
          { name: '@radix-ui/react-checkbox', message: '请使用 @/components/ui/checkbox，确保未勾选复选框统一白底。' }],
      }],
    },
  },
  // 后台 keep-alive 渲染栈：所有已打开标签同时挂载，非活动标签用 hidden 隐藏。
  // 实测（react-router 8.4.0）：useLocation / useSearchParams / useMatches 返回的是
  // 活动标签的值；useParams 虽按标签隔离，但依赖 renderMatches 的 RouteContext 实现细节。
  // 页面一律通过 useAdminTab() 读取路由上下文。见 rules/korean-web.md §4.8。
  {
    files: ['src/features/admin/**/*.{ts,tsx}'],
    ignores: [
      'src/features/admin/routes/**',
      'src/features/admin/components/AdminLayout/**',
      'src/features/admin/components/TabStack/**',
      'src/features/admin/components/AdminTabBar/**',
      'src/features/admin/hooks/useAdminTab.ts',
    ],
    rules: {
      'no-restricted-imports': ['error', {
        paths: [{
          name: 'react-router',
          importNames: ['useParams', 'useLocation', 'useSearchParams', 'useMatch', 'useMatches'],
          message: 'keep-alive 渲染栈下这些 hook 不可信，请改用 @/features/admin/hooks/useAdminTab。',
        }, {
          name: 'radix-ui', importNames: ['Checkbox'], message: '请使用 @/components/ui/checkbox，确保未勾选复选框统一白底。',
        }, {
          name: '@radix-ui/react-checkbox', message: '请使用 @/components/ui/checkbox，确保未勾选复选框统一白底。',
        }],
      }],
    },
  },
  // src/components/ui/** 是 shadcn CLI 拉取的 vendored 源码，非本项目手写代码。
  // WEB-COHESION-005 约束的是 handwritten code；此处比照该规则对 generated output 的
  // 处理方式豁免。仅关闭规模类三条规则，max-depth 与其余检查保持生效。
  // 实测：sidebar.tsx 触发 max-lines(661 有效行) / max-lines-per-function(82, 97) / complexity(12)。
  {
    files: ['src/components/ui/**/*.{ts,tsx}'],
    rules: {
      'max-lines': 'off',
      'max-lines-per-function': 'off',
      complexity: 'off',
    },
  },
  // 既有翻译脚本的 catch 变量不参与逻辑，保留脚本字节内容。
  {
    files: ['scripts/translate.ts'],
    rules: {
      '@typescript-eslint/no-unused-vars': ['error', { caughtErrors: 'none' }],
      // 用户提供的既有翻译脚本按原字节保留；仅该文件使用当前规模基线。
      'max-lines': ['error', { max: 486, skipBlankLines: true, skipComments: true }],
      'max-lines-per-function': ['error', { max: 92, skipBlankLines: true, skipComments: true }],
      complexity: ['error', 11],
      'max-depth': ['error', 4],
    },
  },
];
