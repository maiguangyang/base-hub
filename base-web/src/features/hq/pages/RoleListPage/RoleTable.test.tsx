import { renderToStaticMarkup } from 'react-dom/server';
import { describe, expect, it } from 'vitest';
import { RoleTable, type RoleRow } from './RoleTable';

describe('总部角色列表', () => {
  it('将总部超级管理员稳定键显示为中文名称', () => {
    const rows: RoleRow[] = [{
      id: 'hq-super-admin',
      name: 'role.hqSuperAdministrator',
      kind: 'HQ_SUPER_ADMIN',
      organizationId: 'hq',
      permissions: [],
    }];

    const markup = renderToStaticMarkup(
      <RoleTable
        rows={rows}
        permissions={[]}
        canUpdate={false}
        canDelete={false}
        onEdit={() => undefined}
        onDelete={() => undefined}
      />,
    );

    expect(markup).toContain('总部超级管理员');
    expect(markup).not.toContain('role.hqSuperAdministrator');
  });
});

describe('总部角色权限摘要', () => {
  it('将系统角色的大量权限压缩为范围和数量摘要', () => {
    const rows: RoleRow[] = [{
      id: 'hq-super-admin',
      name: 'role.hqSuperAdministrator',
      kind: 'HQ_SUPER_ADMIN',
      organizationId: 'hq',
      permissions: Array.from({ length: 40 }, (_, index) => ({
        id: `permission-${index}`,
        name: `permission.resource${index}.read`,
        action: `resource${index}:read`,
        module: `resource${index}`,
        scope: 'SYSTEM' as const,
      })),
    }];

    const markup = renderToStaticMarkup(
      <RoleTable rows={rows} permissions={[]} canUpdate={false} canDelete={false} onEdit={() => undefined} onDelete={() => undefined} />,
    );

    expect(markup).toContain('权限范围');
    expect(markup).toContain('全部系统权限');
    expect(markup).toContain('40 项');
    expect(markup).not.toContain('resource0:read');
  });

  it('将自定义角色权限显示为数量摘要', () => {
    const rows: RoleRow[] = [{
      id: 'custom-role',
      name: '运营人员',
      kind: 'CUSTOM',
      organizationId: 'hq',
      permissions: ['read', 'update'].map((action) => ({
        id: action,
        name: `permission.hqStore.${action}`,
        action: `hqStore:${action}`,
        module: 'hqStore',
        scope: 'SYSTEM' as const,
      })),
    }];

    const markup = renderToStaticMarkup(
      <RoleTable rows={rows} permissions={[]} canUpdate={false} canDelete={false} onEdit={() => undefined} onDelete={() => undefined} />,
    );

    expect(markup).toContain('自定义权限');
    expect(markup).toContain('2 项');
    expect(markup).not.toContain('hqStore:read');
  });
});

describe('总部角色表格层级', () => {
  it('使用分层表头和一致的行操作样式', () => {
    const permission = {
      id: 'read',
      name: '查看角色',
      action: 'hqRole:read',
      module: 'hqRole',
      scope: 'SYSTEM' as const,
    };
    const rows: RoleRow[] = [{
      id: 'custom-role',
      name: '运营人员',
      kind: 'CUSTOM',
      organizationId: 'hq',
      permissions: [permission],
    }];

    const markup = renderToStaticMarkup(
      <RoleTable rows={rows} permissions={[permission.action]} canUpdate canDelete onEdit={() => undefined} onDelete={() => undefined} />,
    );

    expect(markup).not.toContain('bg-muted/30');
    expect(markup).not.toContain('h-12');
    expect(markup).toContain('text-right');
    expect(markup).toContain('flex justify-end gap-2');
    expect(markup).toContain('data-size="sm"');
    expect(markup).toContain('data-variant="destructive"');
  });
});
