import ts from 'typescript';

const shellName = 'AdminFormDialogShell';
const shellModule = '@/features/admin/components/AdminFormDialogShell';
const confirmationParts = ['AlertDialog', 'AlertDialogContent', 'AlertDialogHeader', 'AlertDialogTitle', 'AlertDialogDescription', 'AlertDialogFooter'];
const nativeFormControls = new Set(['form', 'input', 'select', 'textarea']);
const sharedFormControlModules = new Map([
  ['Checkbox', '@/components/ui/checkbox'],
  ['Input', '@/components/ui/input'],
  ['RadioGroup', '@/components/ui/radio-group'],
  ['Select', '@/components/ui/select'],
  ['Switch', '@/components/ui/switch'],
  ['Textarea', '@/components/ui/textarea'],
  ['TimePicker', '@/components/ui/time-picker'],
  ['TimeRangePicker', '@/components/ui/time-picker'],
]);

/** 分析后台 Dialog 源码，只对具有表单结构的新增/编辑弹窗施加统一壳约束。 */
export function analyzeAdminDialog(source, fileName = 'dialog.tsx') {
  const file = ts.createSourceFile(fileName, source, ts.ScriptTarget.Latest, true, ts.ScriptKind.TSX);
  const shellImports = new Set();
  const confirmationImports = new Map();
  const formControlImports = new Set();
  const dialogNamespaces = new Set();
  const dialogContentImports = new Set();
  const rendered = new Set();
  const jsxNodes = [];
  const visit = (node) => {
    if (ts.isImportDeclaration(node)) {
      const moduleName = ts.isStringLiteral(node.moduleSpecifier) ? node.moduleSpecifier.text : '';
      const bindings = node.importClause?.namedBindings;
      if (bindings && ts.isNamedImports(bindings)) {
        for (const element of bindings.elements) {
          const importedName = element.propertyName?.text ?? element.name.text;
          const localName = element.name.text;
          if (moduleName === shellModule && importedName === shellName) shellImports.add(localName);
          if (moduleName === '@/components/ui/alert-dialog' && confirmationParts.includes(importedName)) confirmationImports.set(importedName, localName);
          if (sharedFormControlModules.get(importedName) === moduleName) formControlImports.add(localName);
          if (moduleName === 'radix-ui' && importedName === 'Dialog') dialogNamespaces.add(localName);
          if (moduleName === '@/components/ui/dialog' && importedName === 'DialogContent') {
            dialogContentImports.add(localName);
          }
        }
      }
    } else if (ts.isJsxOpeningElement(node) || ts.isJsxSelfClosingElement(node)) {
      const tagName = node.tagName.getText(file);
      rendered.add(tagName);
      jsxNodes.push(node);
    }
    ts.forEachChild(node, visit);
  };
  visit(file);

  const shellRendered = [...shellImports].some((name) => rendered.has(name));
  // 带输入的确认动作仍使用完整的官方 AlertDialog；新增/编辑表单继续强制共享表单壳。
  const hasConfirmationShell = confirmationParts.every((name) => rendered.has(confirmationImports.get(name)));
  const rawDialogContentNames = new Set([
    ...[...dialogNamespaces].map((name) => `${name}.Content`),
    ...dialogContentImports,
  ]);
  const dialogContentNames = new Set([
    ...rawDialogContentNames,
    ...confirmationImports.entries()
      .filter(([name]) => name === 'AlertDialogContent')
      .map(([, localName]) => localName),
  ]);
  const dialogSurfaceNodes = jsxNodes.filter((node) => {
    const name = node.tagName.getText(file);
    return dialogContentNames.has(name) || /DialogShell$/.test(name);
  });
  const isInsideSurface = (node, surface) => {
    for (let parent = node.parent; parent; parent = parent.parent) {
      if (parent === surface.parent) return true;
    }
    return false;
  };
  const isFormControl = (node) => {
    const name = node.tagName.getText(file);
    return nativeFormControls.has(name) || formControlImports.has(name);
  };
  const hasExceptionAttribute = (node) => node.attributes.properties.some((item) => ts.isJsxAttribute(item)
    && item.name.getText() === 'data-admin-form-dialog-exception');
  const surfaceState = dialogSurfaceNodes.map((surface) => {
    const hasFields = jsxNodes.some((node) => node !== surface && isInsideSurface(node, surface)
      && (isFormControl(node) || /(?:Form|Fields|Editor)/.test(node.tagName.getText(file))));
    return {
      hasFields,
      isRaw: rawDialogContentNames.has(surface.tagName.getText(file)),
      hasException: hasExceptionAttribute(surface),
    };
  });
  const hasFormDialog = surfaceState.some((state) => state.hasFields);
  const hasRawFormDialog = surfaceState.some((state) => state.hasFields && state.isRaw);
  const hasInvalidFormDialogException = surfaceState.some((state) => state.hasException);
  const isFormDialog = shellImports.size > 0
    || hasFormDialog;
  return {
    isFormDialog,
    missingShellImport: isFormDialog && shellImports.size === 0 && !hasConfirmationShell,
    shellNotRendered: isFormDialog && shellImports.size > 0 && !shellRendered,
    hasRawDialogContent: isFormDialog && hasRawFormDialog,
    hasInvalidFormDialogException,
  };
}

/** 防止语义扫描范围变化让弹窗契约在零覆盖时静默通过。 */
export function dialogCoverageFailures(scannedFileCount, formDialogCount) {
  if (scannedFileCount === 0) return ['未扫描到任何后台 TSX 文件'];
  if (formDialogCount === 0) return ['未识别到任何表单弹窗'];
  return [];
}
