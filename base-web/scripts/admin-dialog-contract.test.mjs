import { describe, expect, it } from 'vitest';
import { analyzeAdminDialog, dialogCoverageFailures } from './admin-dialog-contract.mjs';

const shellImport = "import { AdminFormDialogShell } from '@/features/admin/components/AdminFormDialogShell';";

describe('WEB-UI-009 表单弹窗契约', () => {
  it('表单控件装配统一弹窗壳时通过', () => {
    const source = `${shellImport}\nexport function ProbeDialog() { return <AdminFormDialogShell><input /></AdminFormDialogShell>; }`;
    expect(analyzeAdminDialog(source)).toEqual({
      isFormDialog: true,
      missingShellImport: false,
      shellNotRendered: false,
      hasRawDialogContent: false,
      hasInvalidFormDialogException: false,
    });
  });

  it('按 DialogContent 内的实际 Input 识别箭头函数表单弹窗', () => {
    const source = "import { DialogContent } from '@/components/ui/dialog';\nimport { Input } from '@/components/ui/input';\nexport const PriceDialog = () => <DialogContent><Input /><button>保存</button></DialogContent>;";
    expect(analyzeAdminDialog(source, 'PriceDialog.tsx')).toMatchObject({
      isFormDialog: true,
      missingShellImport: true,
      hasRawDialogContent: true,
    });
  });

  it('不允许延迟保存声明绕过统一表单弹窗壳', () => {
    const source = "import { DialogContent } from '@/components/ui/dialog';\nimport { Input } from '@/components/ui/input';\nexport const BarcodeDialog = () => <DialogContent data-admin-form-dialog-exception=\"deferred-parent-save\"><Input /></DialogContent>;";
    expect(analyzeAdminDialog(source, 'BarcodeDialog.tsx')).toMatchObject({
      isFormDialog: true,
      missingShellImport: true,
      hasRawDialogContent: true,
      hasInvalidFormDialogException: true,
    });
  });

  it('未知或空白的表单弹窗豁免理由不能通过', () => {
    const source = "import { DialogContent } from '@/components/ui/dialog';\nimport { Input } from '@/components/ui/input';\nexport const PriceDialog = () => <DialogContent data-admin-form-dialog-exception=\"temporary\"><Input /></DialogContent>;";
    expect(analyzeAdminDialog(source, 'PriceDialog.tsx').hasInvalidFormDialogException).toBe(true);
  });

  it('注释和字符串中的壳组件名不算装配', () => {
    const source = `import { DialogContent } from '@/components/ui/dialog';\n// AdminFormDialogShell\nexport function ProbeDialog() { const hint = '<AdminFormDialogShell />'; return <DialogContent><form><input aria-label={hint} /></form></DialogContent>; }`;
    expect(analyzeAdminDialog(source)).toMatchObject({
      isFormDialog: true,
      missingShellImport: true,
      shellNotRendered: false,
    });
  });

  it('只引入但未渲染统一弹窗壳时失败', () => {
    const source = `${shellImport}\nexport function ProbeDialog() { return <textarea />; }`;
    expect(analyzeAdminDialog(source)).toMatchObject({
      isFormDialog: true,
      missingShellImport: false,
      shellNotRendered: true,
    });
  });

  it('字段委托给子组件时仍拦截只引入但未渲染的统一弹窗壳', () => {
    const source = `${shellImport}\nimport { FormFields } from './FormFields';\nexport function ProbeDialog() { return <FormFields />; }`;
    expect(analyzeAdminDialog(source)).toMatchObject({
      isFormDialog: true,
      missingShellImport: false,
      shellNotRendered: true,
    });
  });

  it('不依赖表单控件 import 在源码中的先后顺序', () => {
    const source = "import { DialogContent } from '@/components/ui/dialog';\nexport function ProbeDialog() { return <DialogContent><form><Input /></form></DialogContent>; }\nimport { Input } from '@/components/ui/input';";
    expect(analyzeAdminDialog(source)).toMatchObject({
      isFormDialog: true,
      missingShellImport: true,
    });
  });

  it('拒绝表单弹窗直接使用 Dialog.Content', () => {
    const source = "import { Dialog } from 'radix-ui';\nexport function ProbeDialog() { return <Dialog.Content><form /></Dialog.Content>; }";
    expect(analyzeAdminDialog(source)).toMatchObject({
      isFormDialog: true,
      hasRawDialogContent: true,
    });
  });

  it('带原因输入的确认操作可使用完整的 shadcn AlertDialog', () => {
    const source = "import { AlertDialog, AlertDialogContent, AlertDialogHeader, AlertDialogTitle, AlertDialogDescription, AlertDialogFooter } from '@/components/ui/alert-dialog';\nexport function SuspendDialog() { return <AlertDialog><AlertDialogContent><AlertDialogHeader><AlertDialogTitle>暂停</AlertDialogTitle><AlertDialogDescription>填写原因</AlertDialogDescription></AlertDialogHeader><form><input /><AlertDialogFooter><button>确认</button></AlertDialogFooter></form></AlertDialogContent></AlertDialog>; }";
    expect(analyzeAdminDialog(source, 'SuspendDialog.tsx')).toMatchObject({
      isFormDialog: true,
      missingShellImport: false,
      shellNotRendered: false,
    });
  });

  it('只渲染不完整的 AlertDialog 不豁免表单弹窗壳', () => {
    const source = "import { AlertDialog, AlertDialogContent } from '@/components/ui/alert-dialog';\nexport function SuspendDialog() { return <AlertDialog><AlertDialogContent><form><input /></form></AlertDialogContent></AlertDialog>; }";
    expect(analyzeAdminDialog(source, 'SuspendDialog.tsx').missingShellImport).toBe(true);
  });

  it('不含表单控件的确认弹窗不受本契约约束', () => {
    const source = "import { Dialog } from 'radix-ui';\nexport function ConfirmDialog() { return <Dialog.Content><button>确认</button></Dialog.Content>; }";
    expect(analyzeAdminDialog(source)).toMatchObject({
      isFormDialog: false,
      missingShellImport: false,
      shellNotRendered: false,
      hasRawDialogContent: false,
    });
  });

  it('只有 onSubmit 回调类型的确认弹窗不误判为表单', () => {
    const source = "import { Dialog } from 'radix-ui';\ninterface Props { onSubmit(): void }\nexport function ConfirmDialog(props: Props) { return <Dialog.Content><button onClick={props.onSubmit}>确认</button></Dialog.Content>; }";
    expect(analyzeAdminDialog(source, 'ConfirmDialog.tsx')).toMatchObject({
      isFormDialog: false,
      missingShellImport: false,
      hasRawDialogContent: false,
    });
  });

  it('字段委托给子组件的 FormDialog 仍视为表单弹窗', () => {
    const source = "import { Dialog } from 'radix-ui';\nimport { DirectStoreBasicFields } from './DirectStoreFormSections';\nexport function DirectStoreFormDialog() { return <Dialog.Content><DirectStoreBasicFields /></Dialog.Content>; }";
    expect(analyzeAdminDialog(source, 'DirectStoreFormDialog.tsx')).toMatchObject({
      isFormDialog: true,
      missingShellImport: true,
      hasRawDialogContent: true,
    });
  });

  it('未使用辅助组件上的 onSubmit 不把确认弹窗误判为表单', () => {
    const source = "import { Dialog } from 'radix-ui';\nfunction UnusedSubmit() { return <Widget onSubmit={() => undefined} />; }\nexport function ConfirmDialog() { return <Dialog.Content><button>确认</button></Dialog.Content>; }";
    expect(analyzeAdminDialog(source, 'ConfirmDialog.tsx')).toMatchObject({
      isFormDialog: false,
      missingShellImport: false,
      hasRawDialogContent: false,
    });
  });

  it('同名弹窗壳来自错误模块时仍判定为缺少共享壳', () => {
    const source = "import { AdminFormDialogShell } from './fake';\nexport function ProbeDialog() { return <AdminFormDialogShell><form><input /></form></AdminFormDialogShell>; }";
    expect(analyzeAdminDialog(source, 'ProbeDialog.tsx')).toMatchObject({
      isFormDialog: true,
      missingShellImport: true,
    });
  });
});

describe('弹窗覆盖率保护', () => {
  it('未发现 Dialog 文件时失败', () => {
    expect(dialogCoverageFailures(0, 0)).toContain('未扫描到任何后台 TSX 文件');
  });

  it('扫描到 Dialog 但未识别表单弹窗时失败', () => {
    expect(dialogCoverageFailures(2, 0)).toContain('未识别到任何表单弹窗');
  });

  it('扫描到表单弹窗时通过', () => {
    expect(dialogCoverageFailures(2, 1)).toEqual([]);
  });
});
