import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { LoaderCircle } from 'lucide-react';
import { AdminPageHeader } from '@/features/admin/components/AdminPageHeader';
import { AdminStatusSwitch } from '@/features/admin/components/AdminStatusSwitch';
import { useAuthStore } from '@/features/auth/store/authStore';
import { useModelConfigPage } from './useModelConfigPage';

const statusLabels = {
  UNCONFIGURED: '尚未配置', DRAFT: '待测试', TESTED: '测试通过', ACTIVE: '已启用',
};

export function ModelConfigPage() {
  const viewer = useAuthStore((state) => state.viewer);
  const canRead = viewer?.currentWorkspace?.workspaceType === 'HEADQUARTERS' && viewer.permissions.includes('aiModelConfig:read');
  if (!canRead) return <p role="alert">当前工作区无权查看模型配置。</p>;
  return <ModelConfigForm canManage={viewer.permissions.includes('aiModelConfig:manage')} />;
}

function ModelConfigForm({ canManage }: { canManage: boolean }) {
  const page = useModelConfigPage();
  const editingDisabled = !canManage || page.loading || page.pending;
  return (
    <div className="flex flex-col gap-5">
      <AdminPageHeader title="AI 模型配置" description="配置总部与加盟商工作台共用的模型服务。" />
      <Card className="max-w-2xl">
        <CardHeader>
          <CardTitle>连接设置</CardTitle>
          <CardDescription>当前状态：{modelStatusLabel(page.status?.status)}。密钥保存后只显示是否已配置，不会回显。</CardDescription>
        </CardHeader>
        <CardContent className="flex flex-col gap-5">
          <div className="flex flex-col gap-2">
            <Label htmlFor="ai-model-name">模型名称</Label>
            <Input id="ai-model-name" value={page.modelName} onChange={(event) => page.setModelName(event.target.value)} placeholder="例如：gpt-4.1" disabled={editingDisabled} />
          </div>
          <div className="flex flex-col gap-2">
            <Label htmlFor="ai-model-url">模型地址</Label>
            <Input id="ai-model-url" type="url" value={page.baseUrl} onChange={(event) => page.setBaseUrl(event.target.value)} placeholder="例如：https://model.example/v1" disabled={editingDisabled} />
            <p className="text-xs text-muted-foreground">填写 OpenAI 兼容的 Chat Completions 地址前缀，不含 /chat/completions；不同服务商的前缀可能不同。</p>
            {page.baseUrl.trim().toLowerCase().startsWith('http://') && <p className="text-xs text-destructive">HTTP 会明文传输 API Key 和 AI 请求内容，建议为公网模型服务启用 HTTPS。</p>}
          </div>
          <div className="flex flex-col gap-2">
            <Label htmlFor="ai-model-key">API Key</Label>
            <Input id="ai-model-key" type="password" autoComplete="off" value={page.apiKey} onChange={(event) => page.setApiKey(event.target.value)} placeholder={page.status?.keyConfigured ? '已配置；留空表示沿用现有密钥' : '输入密钥'} disabled={editingDisabled} />
          </div>
          <p className="text-xs text-muted-foreground">保存配置后默认不启用；测试连接通过后，才可开启供两个工作台使用的 AI。修改并保存配置会停用当前模型。</p>
          <Feedback role="alert" message={page.error} />
          <Feedback role="status" message={page.notice} />
          {canManage && <ModelConfigActions page={page} />}
        </CardContent>
      </Card>
    </div>
  );
}

function ModelConfigActions({ page }: { page: ReturnType<typeof useModelConfigPage> }) {
  const state = page.status?.status ?? 'UNCONFIGURED';
  const ready = modelConfigReady(page);
  const busy = page.loading || page.pending;
  const unchanged = modelConfigUnchanged(page);
  return (
    <div className="flex flex-wrap items-center gap-4">
      <div className="flex gap-2">
        <Button type="button" disabled={busy || !ready} onClick={() => void page.save()}>保存</Button>
        <ProbeButton page={page} busy={busy} unchanged={unchanged} />
      </div>
      <div className="flex items-center gap-2">
        <span className="text-sm font-medium">是否启用</span>
        <AdminStatusSwitch
          checked={state === 'ACTIVE'}
          label="是否启用"
          disabled={busy || (state !== 'ACTIVE' && (state !== 'TESTED' || !unchanged))}
          onCheckedChange={(next) => void (next ? page.activate() : page.deactivate())}
        />
      </div>
    </div>
  );
}

function ProbeButton({ page, busy, unchanged }: { page: ReturnType<typeof useModelConfigPage>; busy: boolean; unchanged: boolean }) {
  return <Button type="button" variant="outline" disabled={busy || page.status?.status !== 'DRAFT' || !unchanged} aria-busy={page.testing} onClick={() => void page.probe()}>
    {page.testing && <LoaderCircle className="size-4 animate-spin" aria-hidden="true" />}
    {page.testing ? '测试中…' : '测试连接'}
  </Button>;
}

function Feedback({ role, message }: { role: 'alert' | 'status'; message?: string }) {
  if (!message) return null;
  return <p role={role} className={role === 'alert' ? 'text-sm text-destructive' : 'text-sm text-foreground'}>{message}</p>;
}

function modelStatusLabel(status?: keyof typeof statusLabels) {
  return statusLabels[status ?? 'UNCONFIGURED'];
}

function modelConfigReady(page: ReturnType<typeof useModelConfigPage>) {
  return Boolean(page.modelName.trim() && page.baseUrl.trim() && (page.apiKey || page.status?.keyConfigured));
}

function modelConfigUnchanged(page: ReturnType<typeof useModelConfigPage>) {
  return page.apiKey === '' && page.modelName.trim() === page.status?.modelName && page.baseUrl.trim() === page.status?.baseUrl;
}
