import { describe, expect, it } from 'vitest';
import { permissionDisplayName, type PermissionOption } from './permissionLabels';

const permission = (action: string): PermissionOption => ({ id: action, name: `permission.${action.replace(':', '.')}`, action, module: action.split(':')[0], scope: 'SYSTEM' });

describe('角色权限中文名称', () => {
  it('用加盟商开通的业务动作名称避免重复措辞', () => {
    expect(permissionDisplayName(permission('franchise:provision'))).toBe('开通加盟商');
  });
});
