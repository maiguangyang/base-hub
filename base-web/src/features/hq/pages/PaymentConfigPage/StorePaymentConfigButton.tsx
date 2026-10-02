import type { StoreLifecycle } from '@/__generated__/graphql';
import { Button } from '@/components/ui/button';

export function StorePaymentConfigButton({ lifecycle, onClick }: {
  lifecycle: StoreLifecycle;
  onClick(): void;
}) {
  const disabled = lifecycle !== 'ACTIVE';
  return <span title={disabled ? '门店审核通过并启用后才能配置支付。' : undefined}>
    <Button type="button" variant="outline" size="sm" disabled={disabled} onClick={onClick}>支付配置</Button>
  </span>;
}
