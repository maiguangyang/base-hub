import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import type { PaymentChannel, ScopeRef, ScopeView } from '@/features/hq/api/paymentConfig';
import { PaymentChannelCard } from './PaymentChannelCard';
import { usePaymentConfigPage } from './usePaymentConfigPage';
import { usePaymentConfigDrafts } from './usePaymentConfigDrafts';

const channelLabels: Record<PaymentChannel, string> = { WECHAT: '微信支付', ALIPAY: '支付宝支付' };
export function PaymentConfigPanel({ target, canManage }: { target: ScopeRef; canManage: boolean }) {
  const page = usePaymentConfigPage(target);
  const scopeKey = `${target.scope}:${target.organizationId ?? ''}:${target.storeId ?? ''}`;
  const [activeChannel, setActiveChannel] = useState<PaymentChannel>('WECHAT');
  const drafts = usePaymentConfigDrafts(scopeKey, page.view, page.save);
  return <div className="flex flex-col gap-5">
    <PaymentFeedback page={page} partialNotice={drafts.partialNotice} />
    {page.view && <PaymentTabs view={page.view} page={page} target={target} canManage={canManage} scopeKey={scopeKey}
      activeChannel={activeChannel} setActiveChannel={setActiveChannel} drafts={drafts} />}
  </div>;
}

type Page = ReturnType<typeof usePaymentConfigPage>;
type Drafts = ReturnType<typeof usePaymentConfigDrafts>;

function PaymentFeedback({ page, partialNotice }: { page: Page; partialNotice?: string }) {
  return <>
    {page.error && <div role="alert" className="flex items-center gap-3 text-sm text-destructive"><span>{page.error}</span>{page.conflict && <Button type="button" variant="outline" size="sm" onClick={() => void page.reload()}>重新加载最新数据</Button>}</div>}
    {page.notice && <p role="status" className="text-sm">{page.notice}</p>}
    {partialNotice && <p role="status" className="text-sm">{partialNotice}</p>}
    {page.loading && <p role="status">正在加载配置…</p>}
  </>;
}

function PaymentTabs({ view, page, target, canManage, scopeKey, activeChannel, setActiveChannel, drafts }: {
  view: ScopeView; page: Page; target: ScopeRef; canManage: boolean; scopeKey: string;
  activeChannel: PaymentChannel; setActiveChannel(channel: PaymentChannel): void; drafts: Drafts;
}) {
  const selectedChannel = view.channels.some((item) => item.own.channel === activeChannel)
    ? activeChannel : view.channels[0]?.own.channel;
  const busy = page.loading || page.pending || drafts.savingAll;
  const hasChanges = view.channels.some((item) => drafts.formFor(item).dirty);
  return <Tabs value={selectedChannel} onValueChange={(value) => setActiveChannel(value as PaymentChannel)} className="w-full max-w-3xl gap-5">
    <TabsList className="h-auto w-full justify-start gap-1">
      {view.channels.map((item) => <TabsTrigger key={item.own.channel} value={item.own.channel}>{channelLabels[item.own.channel]}</TabsTrigger>)}
    </TabsList>
    {view.channels.map((item) => <div key={`${scopeKey}:${item.own.channel}`} hidden={selectedChannel !== item.own.channel}>
      <TabsContent value={item.own.channel} forceMount>
        <PaymentChannelCard channel={item.own.channel} item={item} scope={target.scope} canManage={canManage}
          busy={busy} form={drafts.formFor(item)} error={drafts.formErrors[item.own.channel]}
          onFormChange={(form) => drafts.changeForm(item, form)} onState={page.setState} onRestore={page.restoreInheritance} />
      </TabsContent>
    </div>)}
    {canManage && <div><Button type="button" disabled={busy || !hasChanges} onClick={() => void drafts.saveAll(setActiveChannel)}>保存配置</Button></div>}
  </Tabs>;
}
