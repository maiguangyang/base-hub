import { useState } from 'react';
import type { AiEvent } from '@/features/admin/lib/aiStream';
import { CopyTemporaryPasswordButton } from '@/features/admin/components/CopyTemporaryPasswordButton';

export function AiSecretCard({ secret }: { secret: Extract<AiEvent, { type: 'secret' }> }) {
  const [revealedSecret, setRevealedSecret] = useState<typeof secret | null>(null);
  const visible = revealedSecret === secret;
  return <div className="rounded-lg border border-amber-300 bg-amber-50 p-3 text-sm text-amber-950">
    <p className="font-medium">临时密码 · 账号 {secret.targetAccountId}</p>
    <p className="mt-2 break-all font-mono">{visible ? secret.value : '••••••••'}</p>
    <button type="button" aria-pressed={visible} className="mr-3 text-sm text-primary underline" onClick={() => setRevealedSecret(visible ? null : secret)}>{visible ? '隐藏临时密码' : '显示临时密码'}</button>
    <CopyTemporaryPasswordButton password={secret.value} />
    <p className="mt-2 text-xs">仅在当前页面显示，离开工作区后会清除。</p>
  </div>;
}
