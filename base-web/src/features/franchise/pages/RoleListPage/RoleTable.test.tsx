import { renderToStaticMarkup } from 'react-dom/server';
import { describe, expect, it } from 'vitest';
import { RoleTable, type RoleRow } from './RoleTable';

describe('加盟商角色列表', () => {
  it('将加盟商负责人权限压缩为范围和数量摘要', () => {
    const rows: RoleRow[] = [{
      id: 'owner-role',
      name: 'role.franchiseOwner',
      kind: 'FRANCHISE_OWNER',
      organizationId: 'franchise',
      permissions: Array.from({ length: 6 }, (_, index) => ({
        id: `permission-${index}`,
        name: `permission.resource${index}.read`,
        action: `resource${index}:read`,
        module: `resource${index}`,
        scope: 'TENANT' as const,
      })),
    }];

    const markup = renderToStaticMarkup(
      <RoleTable rows={rows} canUpdate={false} canDelete={false} onEdit={() => undefined} onDelete={() => undefined} />,
    );

    expect(markup).toContain('权限范围');
    expect(markup).toContain('全部加盟权限');
    expect(markup).toContain('6 项');
    expect(markup).not.toContain('resource0:read');
  });

  it('使用分层表头和一致的行操作样式', () => {
    const rows: RoleRow[] = [{
      id: 'custom-role',
      name: '店长',
      kind: 'CUSTOM',
      organizationId: 'franchise',
      permissions: [{
        id: 'read',
        name: '查看门店',
        action: 'store:read',
        module: 'store',
        scope: 'TENANT' as const,
      }],
    }];

    const markup = renderToStaticMarkup(
      <RoleTable rows={rows} canUpdate canDelete onEdit={() => undefined} onDelete={() => undefined} />,
    );

    expect(markup).not.toContain('bg-muted/30');
    expect(markup).not.toContain('h-12');
    expect(markup).toContain('text-right');
    expect(markup).toContain('flex justify-end gap-2');
    expect(markup).toContain('data-size="sm"');
    expect(markup).toContain('data-variant="destructive"');
  });
});
