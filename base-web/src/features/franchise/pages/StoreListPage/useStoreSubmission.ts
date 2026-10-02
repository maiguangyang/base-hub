import { useRef, useState } from 'react';
import { useMutation } from '@apollo/client/react';
import { runAdminAction } from '@/features/admin/components/AdminFormDialogShell';
import { SUBMIT_FRANCHISE_STORE_MUTATION } from '../../graphql/stores';
import type { FranchiseStoreRow } from './StoreTable';

export function useStoreSubmission(refetch: () => Promise<unknown>, setError: (message?: string) => void) {
  const [pendingSubmit, setPendingSubmit] = useState<FranchiseStoreRow>();
  const [submitting, setSubmitting] = useState(false);
  const busy = useRef(false);
  const [submitStore] = useMutation(SUBMIT_FRANCHISE_STORE_MUTATION);
  const openSubmit = (row: FranchiseStoreRow) => { setError(undefined); setPendingSubmit(row); };
  const closeSubmit = () => { if (!busy.current) { setPendingSubmit(undefined); setError(undefined); } };

  async function confirmSubmit() {
    if (!pendingSubmit || busy.current) return;
    busy.current = true;
    setSubmitting(true);
    try {
      const submitted = await runAdminAction(async () => {
        await submitStore({ variables: { id: pendingSubmit.id } });
      }, setError);
      if (!submitted) return;
      setPendingSubmit(undefined);
      try { await refetch(); }
      catch { setError('门店已提交审核，但列表刷新失败，请刷新页面。'); }
    } finally {
      busy.current = false;
      setSubmitting(false);
    }
  }

  return { pendingSubmit, submitting, openSubmit, closeSubmit, confirmSubmit };
}
