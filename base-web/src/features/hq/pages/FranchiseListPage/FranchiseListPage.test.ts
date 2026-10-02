import { describe, expect, it } from 'vitest';
import { replaceAccountPassword, resetPasswordErrorMessage } from './useFranchiseListActions';

describe('加盟商初始账号临时密码', () => {
  it('重置共享账号后清除其他加盟商的旧密码', () => {
    const current = {
      'org-a': { accountId: 'shared', password: 'OldPass1' },
      'org-b': { accountId: 'other', password: 'Other123' },
    };
    const next = replaceAccountPassword(current, 'org-c', 'shared', 'NewPass2');
    expect(next).toEqual({
      'org-b': { accountId: 'other', password: 'Other123' },
      'org-c': { accountId: 'shared', password: 'NewPass2' },
    });
    expect(current['org-a'].password).toBe('OldPass1');
  });

  it('针对未激活初始账号给出可理解的失败反馈', () => {
    expect(resetPasswordErrorMessage('MEMBERSHIP_INACTIVE')).toContain('尚未激活');
  });
});
