import type { HqFranchisesQuery } from '@/__generated__/graphql';
import { Button } from '@/components/ui/button';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { CopyTemporaryPasswordButton } from '@/features/admin/components/CopyTemporaryPasswordButton';
import { AdminStatusSwitch } from '@/features/admin/components/AdminStatusSwitch';
import { Link } from 'react-router';

export type FranchiseRow = NonNullable<HqFranchisesQuery['organizations']>['data'][number];
export type VisibleTemporaryPasswords = Record<string, { accountId: string; password: string }>;

interface FranchiseTableProps { rows: FranchiseRow[]; canViewStores: boolean; canSuspend: boolean; canRestore: boolean; canResetPassword: boolean; canConfirmInitialAccount: boolean; resettingOrganizationId?: string; changingStatusOrganizationId?: string; temporaryPasswords: VisibleTemporaryPasswords; onSuspend(row: FranchiseRow): void; onRestore(row: FranchiseRow): void; onResetPassword(row: FranchiseRow): void; onConfirmInitialAccount(row: FranchiseRow): void; onConfigurePayment?(row: FranchiseRow): void }

export function franchiseStorePath(organizationId: string): string {
  return `/admin/hq/stores?${new URLSearchParams({ organizationId, type: 'FRANCHISE' })}`;
}

export function maskTemporaryPassword(password: string): string {
  return '•'.repeat(Math.max(0, password.length - 2)) + password.slice(-2);
}

export function FranchiseTable(props: FranchiseTableProps) {
  return (
    <Table><TableHeader><TableRow><TableHead>编码</TableHead><TableHead>名称</TableHead><TableHead>登录账户</TableHead>{props.canResetPassword && <TableHead>临时密码</TableHead>}<TableHead className="text-center">状态</TableHead><TableHead className="admin-operation-column">操作</TableHead></TableRow></TableHeader>
      <TableBody>{props.rows.map((row) => <FranchiseTableRow key={row.id} row={row} props={props} />)}</TableBody>
    </Table>
  );
}

function FranchiseTableRow({ row, props }: { row: FranchiseRow; props: FranchiseTableProps }) {
  return <TableRow><TableCell>{row.code}</TableCell><TableCell>{row.name}</TableCell><TableCell>{props.canResetPassword ? row.initialAccount?.phone ?? (row.initialAccountId ? '账号已失效' : '未初始化') : '无权限查看'}</TableCell>
    {props.canResetPassword && <TableCell><TemporaryPasswordCell row={row} values={props.temporaryPasswords} /></TableCell>}
    <TableCell><div className="flex justify-center"><AdminStatusSwitch checked={row.status === 'ACTIVE'} label={`${row.name}状态`} disabled={Boolean(props.changingStatusOrganizationId) || (row.status === 'ACTIVE' ? !props.canSuspend : !props.canRestore)} onCheckedChange={(next) => next ? props.onRestore(row) : props.onSuspend(row)} /></div></TableCell>
    <TableCell className="text-right"><FranchiseActions row={row} props={props} /></TableCell>
  </TableRow>;
}

function TemporaryPasswordCell({ row, values }: { row: FranchiseRow; values: VisibleTemporaryPasswords }) {
  if (!row.initialAccountId) return <span className="text-muted-foreground">未初始化</span>;
  if (!row.initialAccount) return <span className="text-muted-foreground">账号已失效</span>;
  const visible = values[row.id];
  const password = visible && visible.accountId === row.initialAccountId ? visible.password : undefined;
  if (!password) return <span className="text-muted-foreground">—</span>;
  return <div className="flex items-center gap-2"><code aria-label="脱敏临时密码">{maskTemporaryPassword(password)}</code><CopyTemporaryPasswordButton key={password} password={password} failureMessage="复制失败，请重试" iconOnly /></div>;
}

function FranchiseActions({ row, props }: { row: FranchiseRow; props: FranchiseTableProps }) {
  return <div className="flex flex-wrap justify-end gap-2">
    {props.canViewStores && <Button asChild variant="outline" size="sm"><Link to={franchiseStorePath(row.id)}>查看门店</Link></Button>}
    {props.onConfigurePayment && <Button type="button" variant="outline" size="sm" onClick={() => props.onConfigurePayment?.(row)}>支付配置</Button>}
    {props.canResetPassword && <InitialAccountActions row={row} props={props} />}
  </div>;
}

function InitialAccountActions({ row, props }: { row: FranchiseRow; props: FranchiseTableProps }) {
  if (row.initialAccountId) return <>
    <Button variant="outline" size="sm" disabled={Boolean(props.resettingOrganizationId) || !row.initialAccount} onClick={() => props.onResetPassword(row)}>{props.resettingOrganizationId === row.id ? '重置中…' : '重置密码'}</Button>
  </>;
  return <>
    {props.canConfirmInitialAccount && <Button variant="outline" size="sm" onClick={() => props.onConfirmInitialAccount(row)}>初始化账户</Button>}
    <Button variant="outline" size="sm" disabled aria-describedby={`initial-account-${row.id}`}>重置密码</Button>
    <span id={`initial-account-${row.id}`} className="sr-only">账户未初始化，暂不可重置密码</span>
  </>;
}
