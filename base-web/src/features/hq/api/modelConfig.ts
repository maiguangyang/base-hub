import { appConfig } from '@/config/app';

export interface ModelConfigStatus {
  modelName: string;
  baseUrl: string;
  keyConfigured: boolean;
  version: number;
  status: 'UNCONFIGURED' | 'DRAFT' | 'TESTED' | 'ACTIVE';
  testedAt?: string;
}

export interface ModelConfigUpdate {
  modelName: string;
  baseUrl: string;
  apiKey: string;
  version: number;
}

export class ModelConfigRequestError extends Error {
  constructor(readonly code: string) { super(code); }
}

const modelConfigURL = `${appConfig.engine.httpUrl}/api/ai/model-config`;

async function request(path: string, method: 'GET' | 'PUT' | 'POST', body?: object): Promise<ModelConfigStatus> {
  let response: Response;
  try {
    response = await fetch(`${modelConfigURL}${path}`, {
      method, credentials: 'include', cache: 'no-store',
      headers: { Accept: 'application/json', ...(body ? { 'Content-Type': 'application/json' } : {}) },
      ...(body ? { body: JSON.stringify(body) } : {}),
    });
  } catch {
    throw new ModelConfigRequestError('NETWORK_ERROR');
  }
  if (!response.ok) {
    const payload = await response.json().catch(() => null) as { code?: string } | null;
    throw new ModelConfigRequestError(payload?.code ?? 'UNKNOWN');
  }
  return response.json() as Promise<ModelConfigStatus>;
}

export const modelConfigAPI = {
  status: () => request('', 'GET'),
  save: (input: ModelConfigUpdate) => request('', 'PUT', input),
  probe: (version: number) => request('/probe', 'POST', { version }),
  activate: (version: number) => request('/activate', 'POST', { version }),
  deactivate: (version: number) => request('/deactivate', 'POST', { version }),
};
