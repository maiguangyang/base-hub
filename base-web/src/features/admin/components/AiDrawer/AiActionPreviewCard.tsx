import type { AiApprovedStep } from '@/features/admin/lib/aiStream';

export function AiActionPreviewCard({ step }: { step: AiApprovedStep }) {
  return <li className="rounded-lg border border-border bg-card p-3 text-sm">
    <div className="flex items-center justify-between gap-2">
      <span className="font-medium">{step.sequence}. {step.title}</span>
      <span className={step.risk === 'HIGH' ? 'text-destructive' : 'text-muted-foreground'}>{step.risk === 'HIGH' ? '高风险' : '需确认'}</span>
    </div>
    {step.targetIds?.length ? <p className="mt-2 text-muted-foreground">目标：{step.targetIds.join('、')}</p> : null}
    {step.scopeIds?.length ? <p className="mt-2 text-muted-foreground">范围：{step.scopeIds.join('、')}</p> : null}
    <PolicyVersionNotice step={step} />
    <dl className="mt-2 grid grid-cols-[auto_1fr] gap-x-3 gap-y-1 break-all">
      {flattenArguments(step.arguments).map(([key, value]) => <div key={key} className="col-span-2 grid grid-cols-subgrid"><dt className="text-muted-foreground">{argumentLabel(key)}</dt><dd>{value}</dd></div>)}
    </dl>
    <p className="mt-2 text-xs text-muted-foreground">最多调用 {step.maxCalls} 次</p>
  </li>;
}

function PolicyVersionNotice({ step }: { step: AiApprovedStep }) {
  const version = step.arguments.expectedVersion;
  if (step.toolId !== 'HqSaveCustomerBenefitPolicy' || typeof version !== 'number' || !Number.isInteger(version) || version < 0) return null;
  return <p className="mt-2 text-muted-foreground">原版本 {version}，预计新版本 {version + 1}</p>;
}

const argumentLabels: Record<string, string> = {
  amountFen: '券面额（分）', minSpendFen: '使用门槛（分）',
  perMemberLimit: '每人发券上限（张）', totalIssueLimit: '总发券上限（张）',
  daysAfterActivation: '激活后有效期（天）', points: '积分数量',
  earnAmountFen: '消费赠分门槛（分）', earnPoints: '赠分数量',
  redeemAmountFen: '抵扣金额（分）', redeemPoints: '抵扣积分数量',
  manualGrantMaxSingle: '单次赠分上限', manualGrantMaxDaily: '每日赠分上限',
  maxRedemptionPoints: '单次抵扣积分上限', maxRedemptionBasisPoints: '单次抵扣比例上限（基点）',
  discountBasisPoints: '会员折扣（基点）', expectedVersion: '原版本',
  evidenceReference: '证明编号', identityEvidence: '身份核验依据',
  dispositionReference: '权益处置凭证', processingBasisCode: '建档依据代码',
  basisCode: '注销依据代码', note: '积分凭证编号',
};

function argumentLabel(path: string): string {
  const field = path.split('.').at(-1) ?? path;
  return argumentLabels[field] ?? path;
}

function flattenArguments(value: Record<string, unknown>, prefix = ''): [string, string][] {
  return Object.entries(value).flatMap(([key, item]) => {
    const path = prefix ? `${prefix}.${key}` : key;
    if (item && typeof item === 'object' && !Array.isArray(item)) return flattenArguments(item as Record<string, unknown>, path);
    return [[path, Array.isArray(item) ? item.map(String).join('、') : String(item ?? '')]];
  });
}
