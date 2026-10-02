import { describe, expect, it } from 'vitest';
import { analyzeListPage, requiredShell } from './admin-list-contract.mjs';

const pageComponents = [...requiredShell, 'AdminListFilters', 'AdminListSearch'];

/** 四个骨架组件和共享搜索的具名导入，各用例共用。 */
const imports = pageComponents
  .map((name) => `import { ${name} } from '@/features/admin/components/${name}';`)
  .join('\n');

/** 列表共享组件的完整装配，代表一个合规页面的主体。 */
const shell = '<AdminPageHeader title="x" /><AdminStatsStrip items={[]} />'
  + '<AdminTableShell toolbar={<AdminToolbar><AdminListFilters hasActiveFilters={false} onReset={() => undefined}><AdminListSearch value="" placeholder="搜索" onSearch={() => undefined} /></AdminListFilters></AdminToolbar>}>{null}</AdminTableShell>';

/** 组装一个页面源码：默认合规，body 可插入待测内容。 */
function page(body = '') {
  return `${imports}\nexport function ProbeListPage() {\n  return <div>${body}${shell}</div>;\n}\n`;
}

/** 合规断言：三项结果全部为空。 */
function expectClean(source) {
  const r = analyzeListPage(source);
  expect({ ...r, missingImports: r.missingImports, notRendered: r.notRendered })
    .toEqual({
      missingImports: [],
      notRendered: [],
      hasTableShell: true,
      hasForbiddenColor: false,
      hasRawListInput: false,
      hasInvalidFilterControl: false,
      missingFilterContainer: false,
      hasFormSurface: false,
      hasPageLocalFormSurfaceColor: false,
      hasMissingFormSurfaceHook: false,
      hasInvalidFormSurfaceException: false,
    });
}

describe('§4.6 结构契约', () => {
  it('完整装配四个骨架组件时通过', () => {
    expectClean(page());
  });

  it('仅在注释中提及组件名不算装配', () => {
    const source = `// ${pageComponents.join(' ')}\nexport function ProbeListPage() { return null; }\n`;
    expect(analyzeListPage(source).missingImports).toEqual(requiredShell);
  });

  it('引入但未渲染时逐个报出', () => {
    const source = `${imports}\nexport function ProbeListPage() { return null; }\n`;
    const r = analyzeListPage(source);
    expect(r.missingImports).toEqual([]);
    expect(r.notRendered).toEqual(requiredShell);
  });

  it('组件名出现在字符串里不算渲染', () => {
    const source = `${imports}\nexport function ProbeListPage() {\n  const hint = '<AdminToolbar />';\n  return <div>{hint}</div>;\n}\n`;
    expect(analyzeListPage(source).notRendered).toEqual(requiredShell);
  });
});

describe('§4.4 颜色约束', () => {
  it('className 中手写状态色被拦下', () => {
    expect(analyzeListPage(page('<span className="bg-green-500" />')).hasForbiddenColor).toBe(true);
  });

  it('模板字面量中手写状态色被拦下', () => {
    expect(analyzeListPage(page('<span className={`p-2 text-red-600`} />')).hasForbiddenColor).toBe(true);
  });

  it('裸十六进制与 oklch 被拦下', () => {
    expect(analyzeListPage(page('<span style={{ color: "#ff0000" }} />')).hasForbiddenColor).toBe(true);
    expect(analyzeListPage(page('<span style={{ color: "oklch(0.5 0.2 30)" }} />')).hasForbiddenColor).toBe(true);
  });

  it('注释里的状态色不算违规', () => {
    expectClean(page('{/* className="bg-green-500" 是禁止写法 */}'));
  });

  it('页面文案不会被误判为 class', () => {
    expectClean(page('<p>配色说明见设计文档</p>'));
  });
});

describe('WEB-UI-009 列表交互契约', () => {
  const search = '<AdminListSearch value="" placeholder="搜索" onSearch={() => undefined} />';

  it('允许列表正文渲染与搜索无关的 Input', () => {
    const source = `import { Input } from '@/components/ui/input';\n${page('<Input aria-label="行内值" />')}`;
    expect(analyzeListPage(source).hasRawListInput).toBe(false);
  });

  it('拦截工具栏直接渲染基础 Input', () => {
    const source = `import { Input } from '@/components/ui/input';\n${page().replace(search, '<Input aria-label="搜索" />')}`;
    expect(analyzeListPage(source).hasRawListInput).toBe(true);
  });

  it('拦截工具栏直接渲染原生 input', () => {
    expect(analyzeListPage(page().replace(search, '<input aria-label="搜索" />')).hasRawListInput).toBe(true);
  });

  it('拦截工具栏通过命名空间导入渲染基础 Input', () => {
    const source = `import * as InputUI from '@/components/ui/input';\n${page().replace(search, '<InputUI.Input />')}`;
    expect(analyzeListPage(source).hasRawListInput).toBe(true);
  });

  it('不依赖 import 在源码中的先后顺序', () => {
    const source = `${page().replace(search, '<Input aria-label="搜索" />')}\nimport { Input } from '@/components/ui/input';`;
    expect(analyzeListPage(source).hasRawListInput).toBe(true);
  });

  it('使用共享搜索时通过', () => {
    const source = page();
    expect(analyzeListPage(source)).toMatchObject({
      hasRawListInput: false,
    });
  });

  it('同名组件来自错误模块时不算满足共享骨架导入', () => {
    const source = page().replace(
      "import { AdminToolbar } from '@/features/admin/components/AdminToolbar';",
      "import { AdminToolbar } from './fake-toolbar';",
    );
    expect(analyzeListPage(source).missingImports).toContain('AdminToolbar');
  });

  it('业务没有搜索能力时不强制渲染搜索框', () => {
    const source = page()
      .replace("import { AdminListSearch } from '@/features/admin/components/AdminListSearch';\n", '')
      .replace('<AdminListSearch value="" placeholder="搜索" onSearch={() => undefined} />', '<span>无筛选</span>');
    expectClean(source);
  });
});

describe('WEB-UI-010 结构化筛选契约', () => {
  it('识别文件名不含 List 的实际列表工作台', () => {
    const source = page().replaceAll('ProbeListPage', 'ProductBrandPage');
    expect(analyzeListPage(source).hasTableShell).toBe(true);
  });
  it('拦截空共享容器旁边直接渲染的基础 Select', () => {
    const source = `import { Select } from '@/components/ui/select';\n${page().replace(
      '<AdminListFilters hasActiveFilters={false} onReset={() => undefined}><AdminListSearch value="" placeholder="搜索" onSearch={() => undefined} /></AdminListFilters>',
      '<AdminListFilters hasActiveFilters={false} onReset={() => undefined}></AdminListFilters><Select />',
    )}`;
    expect(analyzeListPage(source).hasInvalidFilterControl).toBe(true);
  });

  it('拦截共享搜索渲染在筛选容器之外', () => {
    const source = page().replace(
      '<AdminListFilters hasActiveFilters={false} onReset={() => undefined}><AdminListSearch value="" placeholder="搜索" onSearch={() => undefined} /></AdminListFilters>',
      '<AdminListFilters hasActiveFilters={false} onReset={() => undefined}></AdminListFilters><AdminListSearch value="" placeholder="搜索" onSearch={() => undefined} />',
    );
    expect(analyzeListPage(source).hasInvalidFilterControl).toBe(true);
  });

  it('缺少共享筛选容器时被拦截', () => {
    const source = page()
      .replace("import { AdminListFilters } from '@/features/admin/components/AdminListFilters';\n", '')
      .replace('<AdminListFilters hasActiveFilters={false} onReset={() => undefined}>', '')
      .replace('</AdminListFilters></AdminToolbar>', '</AdminToolbar>');
    expect(analyzeListPage(source).missingFilterContainer).toBe(true);
  });

  it('同名容器来自错误模块时被拦截', () => {
    const source = page().replace(
      "import { AdminListFilters } from '@/features/admin/components/AdminListFilters';",
      "import { AdminListFilters } from './fake-filters';",
    );
    expect(analyzeListPage(source).missingFilterContainer).toBe(true);
  });

  it('引入但未渲染共享筛选容器时被拦截', () => {
    const source = page()
      .replace('<AdminListFilters hasActiveFilters={false} onReset={() => undefined}>', '')
      .replace('</AdminListFilters></AdminToolbar>', '</AdminToolbar>');
    expect(analyzeListPage(source).missingFilterContainer).toBe(true);
  });

  it('共享筛选容器未放在 AdminToolbar 内时被拦截', () => {
    const source = page()
      .replace(
        '<AdminToolbar><AdminListFilters hasActiveFilters={false} onReset={() => undefined}>',
        '<AdminListFilters hasActiveFilters={false} onReset={() => undefined}><AdminToolbar>',
      )
      .replace('</AdminListFilters></AdminToolbar>', '</AdminToolbar></AdminListFilters>');
    expect(analyzeListPage(source).missingFilterContainer).toBe(true);
  });
});

describe('WEB-UI-009 表单表面契约', () => {
  it('按控件语义发现共享 Input，而不依赖文件名', () => {
    const source = "import { Input as Field } from '@/components/ui/input';\nexport function PolicyPage() { return <Field />; }";
    expect(analyzeListPage(source)).toMatchObject({
      hasFormSurface: true,
      hasPageLocalFormSurfaceColor: false,
    });
  });

  it('拦截直接写在表单控件上的局部白底黑字', () => {
    const source = "import { Input } from '@/components/ui/input';\nexport function Page() { return <Input className=\"!bg-white !text-black\" />; }";
    expect(analyzeListPage(source).hasPageLocalFormSurfaceColor).toBe(true);
  });

  it('拦截页面祖先选择器覆盖表单控件颜色', () => {
    const source = "import { Input } from '@/components/ui/input';\nconst patch = '[&_input]:!bg-white [&_input]:!text-black';\nexport function Page() { return <div className={patch}><Input /></div>; }";
    expect(analyzeListPage(source).hasPageLocalFormSurfaceColor).toBe(true);
  });

  it('拦截可见复合字段外壳的局部颜色补丁', () => {
    const source = "import { AdminStatusSwitch as Status } from '@/features/admin/components/AdminStatusSwitch';\nexport function Page() { return <div className=\"bg-white\"><Status /></div>; }";
    expect(analyzeListPage(source).hasPageLocalFormSurfaceColor).toBe(true);
  });

  it('拦截未接入共享语义的自定义 Popover 表单面板', () => {
    const source = "import { PopoverContent as Panel } from '@/components/ui/popover';\nexport function PricePicker() { return <Panel><button>选择价格</button></Panel>; }";
    expect(analyzeListPage(source)).toMatchObject({
      hasFormSurface: true,
      hasMissingFormSurfaceHook: true,
    });
  });

  it('拦截未接入共享语义的 combobox 触发器', () => {
    const source = 'export function StorePicker() { return <button role="combobox">选择门店</button>; }';
    expect(analyzeListPage(source).hasMissingFormSurfaceHook).toBe(true);
  });

  it('拦截 DropdownMenuTrigger 内未接入共享语义的自定义触发按钮', () => {
    const source = "import { DropdownMenuTrigger as Trigger } from '@/components/ui/dropdown-menu';\nexport function RolePicker() { return <Trigger asChild><button>选择角色</button></Trigger>; }";
    expect(analyzeListPage(source).hasMissingFormSurfaceHook).toBe(true);
  });

  it('非表单动作浮层必须使用已登记的显式例外', () => {
    const source = "import { DropdownMenuContent as Menu } from '@/components/ui/dropdown-menu';\nexport function ActionMenu() { return <Menu data-admin-form-surface-exception=\"non-form-action-menu\" />; }";
    expect(analyzeListPage(source)).toMatchObject({
      hasMissingFormSurfaceHook: false,
      hasInvalidFormSurfaceException: false,
    });
  });

  it('未知的表单表面例外不能通过', () => {
    const source = "import { PopoverContent } from '@/components/ui/popover';\nexport function Help() { return <PopoverContent data-admin-form-surface-exception=\"temporary\" />; }";
    expect(analyzeListPage(source).hasInvalidFormSurfaceException).toBe(true);
  });

  it('认可自定义控件暴露共享语义标记', () => {
    const source = 'export function CustomControl() { return <button data-admin-form-surface="" />; }';
    expect(analyzeListPage(source)).toMatchObject({
      hasFormSurface: true,
      hasPageLocalFormSurfaceColor: false,
      hasMissingFormSurfaceHook: false,
    });
  });
});

// 以下六例对应评审阶段用一次性探针发现的真实缺陷。
// 手写扫描器在这些输入上会把后续代码整段吞掉，让合规页面报出
// 「引入了但未渲染」这种错误结论；正则替换方案则会让 https:// 同行的
// 违规颜色逃过检查。固化为测试，防止再退回到那两种实现。
describe('解析边界（回归防护）', () => {
  it('JSX 文本中的英文撇号不吞掉后续 JSX', () => {
    expectClean(page("<p>Don't worry</p>"));
  });

  it('字符串中的 // 不被当作注释起点', () => {
    expectClean(page("<a href=\"https://docs.example.com\">doc</a>"));
  });

  it('违规颜色与含 // 的字符串同行时仍被拦下', () => {
    const source = `${imports}\nexport function ProbeListPage() {\n  const doc = 'https://x'; const bad = 'bg-green-500';\n  return <div>{doc}{bad}${shell}</div>;\n}\n`;
    expect(analyzeListPage(source).hasForbiddenColor).toBe(true);
  });

  it('正则字面量中的引号不吞掉后续 JSX', () => {
    expectClean(page('{String(/[\'"]/g)}'));
  });

  it('转义引号不导致字符串提前结束', () => {
    expectClean(page("{'a\\'b'}"));
  });

  it('模板字面量中的 ${} 插值不影响判定', () => {
    expectClean(page('{`a ${1 + 1} b`}'));
  });
});
