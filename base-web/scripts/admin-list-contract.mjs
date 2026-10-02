import ts from 'typescript';

/** §4.6 规定的页面层 + 工作台层骨架组件，缺一不可。 */
export const requiredShell = ['AdminPageHeader', 'AdminStatsStrip', 'AdminTableShell', 'AdminToolbar'];

const requiredComponentModules = new Map(
  requiredShell.map((name) => [name, `@/features/admin/components/${name}`]),
);
const filterContainerModule = '@/features/admin/components/AdminListFilters';
const canonicalFilterModules = new Map([
  ['AdminListSearch', '@/features/admin/components/AdminListSearch'],
  ['AdminFilterSelect', '@/features/admin/components/AdminFilterSelect'],
]);
const formSurfaceModules = new Map([
  ['@/components/ui/input', new Set(['Input'])],
  ['@/components/ui/textarea', new Set(['Textarea'])],
  ['@/components/ui/select', new Set(['SelectTrigger', 'SelectContent'])],
  ['@/components/ui/time-picker', new Set(['TimePicker', 'TimeRangePicker'])],
  ['@/components/ui/checkbox', new Set(['Checkbox'])],
  ['@/components/ui/radio-group', new Set(['RadioGroupItem'])],
  ['@/features/admin/components/AdminFormSelect', new Set(['AdminFormSelect'])],
  ['@/features/admin/components/AdminRoleSelect', new Set(['AdminRoleSelect'])],
  ['@/features/admin/components/AdminSearchSelect', new Set(['AdminSearchSelect'])],
  ['@/features/admin/components/AdminStatusSwitch', new Set(['AdminStatusSwitch'])],
]);
const customSurfaceModules = new Map([
  ['@/components/ui/popover', {
    surfaces: new Set(['PopoverContent']),
    triggers: new Set(['PopoverTrigger']),
  }],
  ['@/components/ui/dropdown-menu', {
    surfaces: new Set(['DropdownMenuContent']),
    triggers: new Set(['DropdownMenuTrigger']),
  }],
]);
const compositeFormSurfaceModules = new Map([
  ['@/features/admin/components/AdminStatusSwitch', new Set(['AdminStatusSwitch'])],
]);
const allowedFormSurfaceExceptions = new Set(['non-form-action-menu']);

/** §4.4 禁止在后台页面手写状态颜色 class；tone 必须来自共享映射。 */
const forbiddenColor = /\b(?:bg|text|border)-(?:red|green|yellow|amber|emerald|rose|lime|orange)-\d{2,3}\b|oklch\(|#[0-9a-fA-F]{6}\b/;
const localFormColor = /(?:^|\s)!?(?:bg-white|text-black)(?=\s|$)/;
const descendantFormColor = /\[&_[^\]]*(?:input|textarea|select|data-slot)[^\]]*\]:!?(?:bg-white|text-black)/;

/** 字符串与模板字面量节点——类名住在这里，颜色检查以此为准。
 *  不含 JsxText：页面文案不是 class，纳入只会制造误报。 */
function isClassBearingLiteral(node) {
  return ts.isStringLiteral(node)
    || ts.isNoSubstitutionTemplateLiteral(node)
    || ts.isTemplateHead(node)
    || ts.isTemplateMiddle(node)
    || ts.isTemplateTail(node);
}

function isInsideJsxTag(node, tagName, file) {
  for (let parent = node.parent; parent; parent = parent.parent) {
    if (ts.isJsxElement(parent) && parent.openingElement.tagName.getText(file) === tagName) return true;
  }
  return false;
}

/** 判断 JSX 节点是否显式接入后台表单表面语义。 */
function hasFormSurfaceAttribute(node) {
  return node.attributes.properties.some((attribute) => ts.isJsxAttribute(attribute)
    && attribute.name.getText() === 'data-admin-form-surface');
}

/** 读取可审查的静态 JSX 字符串属性；动态值和空属性按无效值处理。 */
function jsxStringAttribute(node, name) {
  const attribute = node.attributes.properties.find((item) => ts.isJsxAttribute(item)
    && item.name.getText() === name);
  if (!attribute || !ts.isJsxAttribute(attribute)) return undefined;
  return ts.isStringLiteral(attribute.initializer) ? attribute.initializer.text : '';
}

/** 收集一个 className 属性内的静态字符串片段。 */
function classLiterals(node) {
  const attribute = node.attributes.properties.find((item) => ts.isJsxAttribute(item)
    && item.name.getText() === 'className');
  if (!attribute || !ts.isJsxAttribute(attribute) || !attribute.initializer) return [];
  const values = [];
  const collect = (child) => {
    if (isClassBearingLiteral(child)) values.push(child.text);
    ts.forEachChild(child, collect);
  };
  collect(attribute.initializer);
  return values;
}

/** 判断一个 JSX 外壳下是否包含表单表面控件。 */
function containsFormSurface(node, file, formSurfaceNames) {
  let found = false;
  const visit = (child) => {
    if (found) return;
    if ((ts.isJsxOpeningElement(child) || ts.isJsxSelfClosingElement(child)) && child !== node) {
      const tagName = child.tagName.getText(file);
      if (['input', 'textarea', 'select'].includes(tagName)
        || formSurfaceNames.has(tagName)
        || hasFormSurfaceAttribute(child)) {
        found = true;
        return;
      }
    }
    ts.forEachChild(child, visit);
  };
  ts.forEachChild(node.parent, visit);
  return found;
}

/**
 * 分析一个后台列表页源码，判定它是否满足 §4.6 结构契约与 §4.4 颜色约束。
 *
 * 用 TypeScript 编译器 API 而非正则或手写状态机：注释、字符串、JSX 文本、
 * 正则字面量、转义引号全部由编译器正确区分。手写扫描器在
 * `<p>don't</p>` 这类 JSX 文本上会把撇号当成字符串起点，吞掉后续全部 JSX，
 * 让正确的页面报出「引入了但未渲染」这种错误结论。
 * 仓库既有的 check-translate-baseline.mjs 用的是同一套 API。
 *
 * @returns {{ missingImports: string[], notRendered: string[], hasForbiddenColor: boolean, hasRawListInput: boolean, hasInvalidFilterControl: boolean, missingFilterContainer: boolean, hasFormSurface: boolean, hasMissingFormSurfaceHook: boolean, hasInvalidFormSurfaceException: boolean, hasPageLocalFormSurfaceColor: boolean }}
 */
export function analyzeListPage(source) {
  const file = ts.createSourceFile('page.tsx', source, ts.ScriptTarget.Latest, true, ts.ScriptKind.TSX);
  const requiredImports = new Map();
  const rendered = new Set();
  const rawInputNames = new Set();
  const rawSelectNames = new Set();
  const canonicalFilterNames = new Set();
  const formSurfaceNames = new Set();
  const customSurfaceNames = new Set();
  const customSurfaceTriggerNames = new Set();
  const compositeFormSurfaceNames = new Set();
  const jsxNodes = [];
  const sourceLiterals = [];
  let hasForbiddenColor = false;

  const visit = (node) => {
    if (ts.isImportDeclaration(node)) {
      const moduleName = ts.isStringLiteral(node.moduleSpecifier) ? node.moduleSpecifier.text : '';
      const bindings = node.importClause?.namedBindings;
      if (bindings && ts.isNamedImports(bindings)) {
        for (const element of bindings.elements) {
          const importedName = element.propertyName?.text ?? element.name.text;
          const localName = element.name.text;
          if (requiredComponentModules.get(importedName) === moduleName) {
            requiredImports.set(importedName, localName);
          }
          if (moduleName === '@/components/ui/input' && importedName === 'Input') {
            rawInputNames.add(localName);
          }
          if (moduleName === '@/components/ui/select' && importedName === 'Select') {
            rawSelectNames.add(localName);
          }
          if (canonicalFilterModules.get(importedName) === moduleName) {
            canonicalFilterNames.add(localName);
          }
          if (moduleName === filterContainerModule && importedName === 'AdminListFilters') {
            requiredImports.set('AdminListFilters', localName);
          }
          if (formSurfaceModules.get(moduleName)?.has(importedName)) {
            formSurfaceNames.add(localName);
          }
          const customSurfaceModule = customSurfaceModules.get(moduleName);
          if (customSurfaceModule?.surfaces.has(importedName)) customSurfaceNames.add(localName);
          if (customSurfaceModule?.triggers.has(importedName)) customSurfaceTriggerNames.add(localName);
          if (compositeFormSurfaceModules.get(moduleName)?.has(importedName)) {
            compositeFormSurfaceNames.add(localName);
          }
        }
      } else if (bindings && ts.isNamespaceImport(bindings)) {
        if (moduleName === '@/components/ui/input') rawInputNames.add(`${bindings.name.text}.Input`);
        if (moduleName === '@/components/ui/select') rawSelectNames.add(`${bindings.name.text}.Select`);
      }
    } else if (ts.isJsxOpeningElement(node) || ts.isJsxSelfClosingElement(node)) {
      const tagName = node.tagName.getText(file);
      rendered.add(tagName);
      jsxNodes.push(node);
    } else if (isClassBearingLiteral(node) && forbiddenColor.test(node.text)) {
      hasForbiddenColor = true;
    }
    if (isClassBearingLiteral(node)) sourceLiterals.push(node.text);
    ts.forEachChild(node, visit);
  };
  visit(file);
  const isFormSurface = (node) => {
    const tagName = node.tagName.getText(file);
    return ['input', 'textarea', 'select'].includes(tagName)
      || formSurfaceNames.has(tagName)
      || hasFormSurfaceAttribute(node);
  };
  const isCustomSurfaceCandidate = (node) => {
    const tagName = node.tagName.getText(file);
    if (customSurfaceNames.has(tagName)) return true;
    if (tagName !== 'button') return false;
    if (jsxStringAttribute(node, 'role') === 'combobox') return true;
    return [...customSurfaceTriggerNames].some((triggerName) => isInsideJsxTag(node, triggerName, file));
  };
  const customSurfaceCandidates = jsxNodes.filter(isCustomSurfaceCandidate);
  const hasValidFormSurfaceException = (node) => {
    const exception = jsxStringAttribute(node, 'data-admin-form-surface-exception');
    return exception !== undefined && allowedFormSurfaceExceptions.has(exception);
  };
  const hasFormSurface = jsxNodes.some(isFormSurface) || customSurfaceCandidates.length > 0;
  const hasMissingFormSurfaceHook = customSurfaceCandidates.some((node) => !hasFormSurfaceAttribute(node)
    && !hasValidFormSurfaceException(node));
  const hasInvalidFormSurfaceException = customSurfaceCandidates.some((node) => {
    const exception = jsxStringAttribute(node, 'data-admin-form-surface-exception');
    return exception !== undefined && !allowedFormSurfaceExceptions.has(exception);
  });
  const hasDirectFormColor = jsxNodes.some((node) => isFormSurface(node)
    && classLiterals(node).some((value) => localFormColor.test(value)));
  const hasCompositeFormColor = jsxNodes.some((node) => ['div', 'label'].includes(node.tagName.getText(file))
    && classLiterals(node).some((value) => localFormColor.test(value))
    && containsFormSurface(node, file, compositeFormSurfaceNames));
  const hasPageLocalFormSurfaceColor = hasDirectFormColor
    || hasCompositeFormColor
    || sourceLiterals.some((value) => descendantFormColor.test(value));
  const toolbarLocalName = requiredImports.get('AdminToolbar');
  const filterContainerLocalName = requiredImports.get('AdminListFilters');
  const hasRawListInput = toolbarLocalName !== undefined && jsxNodes.some((node) => {
    const tagName = node.tagName.getText(file);
    const isRawInput = tagName === 'input' || rawInputNames.has(tagName);
    return isRawInput && isInsideJsxTag(node, toolbarLocalName, file);
  });
  const hasFilterContainer = toolbarLocalName !== undefined && filterContainerLocalName !== undefined
    && jsxNodes.some((node) => node.tagName.getText(file) === filterContainerLocalName
      && isInsideJsxTag(node, toolbarLocalName, file));
  const hasInvalidFilterControl = toolbarLocalName !== undefined && jsxNodes.some((node) => {
    if (!isInsideJsxTag(node, toolbarLocalName, file)) return false;
    const tagName = node.tagName.getText(file);
    if (tagName === 'select' || rawSelectNames.has(tagName)) return true;
    return canonicalFilterNames.has(tagName)
      && (filterContainerLocalName === undefined || !isInsideJsxTag(node, filterContainerLocalName, file));
  });

  return {
    hasTableShell: requiredImports.has('AdminTableShell') && rendered.has(requiredImports.get('AdminTableShell')),
    missingImports: requiredShell.filter((name) => !requiredImports.has(name)),
    notRendered: requiredShell.filter((name) => {
      const localName = requiredImports.get(name);
      return localName !== undefined && !rendered.has(localName);
    }),
    hasForbiddenColor,
    hasRawListInput,
    hasInvalidFilterControl,
    missingFilterContainer: !hasFilterContainer,
    hasFormSurface,
    hasMissingFormSurfaceHook,
    hasInvalidFormSurfaceException,
    hasPageLocalFormSurfaceColor,
  };
}
