import { appConfig } from '@/config/app';

interface Input {
  organizationId: string;
  accountId: string;
  evidenceReference: string;
  attested: boolean;
}

export class InitialAccountRequestError extends Error {
  constructor(readonly code: string) { super(code); }
}

export async function setFranchiseInitialAccount(input: Input): Promise<void> {
  let response: Response;
  try {
    response = await fetch(`${appConfig.engine.httpUrl}/api/franchise-initial-account`, {
      method: 'POST', credentials: 'include', headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(input),
    });
  } catch {
    throw new InitialAccountRequestError('RESULT_UNKNOWN');
  }
  if (!response.ok) {
    const payload = await response.json().catch(() => null) as { code?: string } | null;
    throw new InitialAccountRequestError(payload?.code ?? 'UNKNOWN');
  }
}
