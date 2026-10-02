import { useState } from 'react';
import type { ChannelView, PaymentChannel, ScopeView } from '@/features/hq/api/paymentConfig';
import { buildDraft, emptyForm, type ChannelForm } from './PaymentChannelCard';
import type { SaveDraft } from './usePaymentConfigPage';

type DraftEntry = { key: string; form: ChannelForm };
type Prepared = { item: ChannelView; draft: SaveDraft | string };

function draftKey(scopeKey: string, item: ChannelView): string {
  return `${scopeKey}:${item.own.channel}:${item.own.recordId ?? ''}:${item.own.version}:${item.own.ratePpm}`;
}

export function usePaymentConfigDrafts(scopeKey: string, view: ScopeView | undefined, save: (channel: PaymentChannel, draft: SaveDraft) => Promise<boolean>) {
  const [drafts, setDrafts] = useState<Partial<Record<PaymentChannel, DraftEntry>>>({});
  const [formErrors, setFormErrors] = useState<Partial<Record<PaymentChannel, string>>>({});
  const [savingAll, setSavingAll] = useState(false);
  const [partialNotice, setPartialNotice] = useState<string>();

  const formFor = (item: ChannelView) => {
    const entry = drafts[item.own.channel];
    return entry?.key === draftKey(scopeKey, item) ? entry.form : emptyForm(item);
  };
  const changeForm = (item: ChannelView, form: ChannelForm) => {
    setDrafts((previous) => ({ ...previous, [item.own.channel]: { key: draftKey(scopeKey, item), form } }));
    setFormErrors((previous) => ({ ...previous, [item.own.channel]: undefined }));
    setPartialNotice(undefined);
  };
  async function saveAll(onInvalid: (channel: PaymentChannel) => void) {
    if (!view || savingAll) return;
    const prepared: Prepared[] = view.channels.filter((item) => formFor(item).dirty)
      .map((item) => ({ item, draft: buildDraft(item.own.channel, item, formFor(item)) }));
    const invalid = prepared.find(({ draft }) => typeof draft === 'string');
    if (invalid) {
      onInvalid(invalid.item.own.channel);
      setFormErrors((previous) => ({ ...previous, [invalid.item.own.channel]: invalid.draft as string }));
      return;
    }
    setSavingAll(true); setPartialNotice(undefined);
    let completed = 0;
    try {
      for (const { item, draft } of prepared) {
        if (!await save(item.own.channel, draft as SaveDraft)) break;
        completed += 1;
        setDrafts((previous) => ({ ...previous, [item.own.channel]: { key: draftKey(scopeKey, item), form: emptyForm(item) } }));
      }
      if (completed > 0 && completed < prepared.length) setPartialNotice(`已保存 ${completed} 个支付方式，其余配置未保存，请检查错误后重试。`);
    } finally { setSavingAll(false); }
  }

  return { formFor, changeForm, formErrors, savingAll, partialNotice, saveAll };
}
