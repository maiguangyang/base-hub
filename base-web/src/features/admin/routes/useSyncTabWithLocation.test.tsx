// @vitest-environment jsdom

import { cleanup, render, waitFor } from '@testing-library/react';
import { afterEach, describe, expect, it } from 'vitest';
import { MemoryRouter } from 'react-router';
import { createAdminTabManifest, createInitialSnapshot } from '@/stores/slices/tabOperations';
import { useAdminTabsStore } from '@/stores/slices/adminTabsSlice';
import { useSyncTabWithLocation } from './useSyncTabWithLocation';

function Probe() {
  const result = useSyncTabWithLocation();
  return <span>{String(result.isKnownPath)}</span>;
}

afterEach(() => cleanup());

describe('地址与后台标签同步', () => {
  it('把加盟门店页的查询上下文写入当前标签', async () => {
    const manifest = createAdminTabManifest('HEADQUARTERS', ['organization:read', 'store:read_all']);
    useAdminTabsStore.setState({ ...createInitialSnapshot(manifest), manifest, erroredPaths: [] });

    render(
      <MemoryRouter initialEntries={['/admin/hq/franchise-stores?organizationId=org-a']}>
        <Probe />
      </MemoryRouter>,
    );

    await waitFor(() => expect(useAdminTabsStore.getState().tabs.find(
      (tab) => tab.path === '/admin/hq/stores',
    )?.search).toBe('?organizationId=org-a&type=FRANCHISE'));
  });
});
