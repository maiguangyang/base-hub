import { AdminPageHeader } from '@/features/admin/components/AdminPageHeader';
import { useAuthStore } from '@/features/auth/store/authStore';
import { PaymentConfigPanel } from './PaymentConfigPanel';

export function PaymentConfigPage() {
  const viewer = useAuthStore((state) => state.viewer);
  const canRead = viewer?.currentWorkspace?.workspaceType === 'HEADQUARTERS' && viewer.permissions.includes('paymentConfig:read');
  if (!canRead) return <p role="alert">当前工作区无权查看支付配置。</p>;
  return <div className="flex flex-col gap-5">
    <AdminPageHeader title="全局支付配置" description="设置全局共用的微信与支付宝商户资料及预计平台费率；加盟商和门店的专属配置请在对应列表中打开。" />
    <PaymentConfigPanel target={{ scope: 'GLOBAL' }} canManage={viewer.permissions.includes('paymentConfig:manage')} />
  </div>;
}
