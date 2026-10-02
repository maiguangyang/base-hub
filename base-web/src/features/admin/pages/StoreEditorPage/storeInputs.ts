import type { StoreBusinessStatus } from '@/__generated__/graphql';
import { generateEntityCode } from '@/lib/codeGenerator';
import type { StoreWorkspace } from './storeEditorPolicy';

export interface StoreValues {
  name: string;
  code?: string;
  contactPhone?: string;
  managerName?: string;
  managerPhone?: string;
  province?: string;
  city?: string;
  district?: string;
  address?: string;
  businessHours?: string;
  businessStatus?: StoreBusinessStatus;
  supportDineIn?: boolean;
  supportTakeout?: boolean;
  storeArea?: number | string;
  tableCount?: number | string;
  receiptFooter?: string;
}

type NullableStoreField = Exclude<keyof StoreValues, 'name' | 'code' | 'storeArea' | 'tableCount'>;
export type StoreRecord = { name: string; code?: string | null } & {
  [K in NullableStoreField]?: StoreValues[K] | null;
} & { storeArea?: number | string | null; tableCount?: number | string | null };

export function defaultStoreValues(): StoreValues {
  return {
    name: '', contactPhone: '', managerName: '', managerPhone: '',
    province: '', city: '', district: '', address: '',
    businessHours: '09:00 - 22:00', businessStatus: 'OPEN',
    supportDineIn: true, supportTakeout: true,
    storeArea: '', tableCount: '', receiptFooter: '',
  };
}

function trimStr(value?: string): string { return value?.trim() ?? ''; }
function optionalStr(value?: string): string | null { return trimStr(value) || null; }
function parseNumber(value?: number | string): number | null {
  if (value === undefined || value === '') return null;
  const parsed = Number(value);
  return Number.isNaN(parsed) ? null : parsed;
}

export function storeUpdateInput(values: StoreValues) {
  return {
    name: trimStr(values.name), contactPhone: trimStr(values.contactPhone),
    managerName: optionalStr(values.managerName), managerPhone: optionalStr(values.managerPhone),
    province: trimStr(values.province), city: trimStr(values.city),
    district: trimStr(values.district), address: trimStr(values.address),
    businessHours: trimStr(values.businessHours),
    businessStatus: values.businessStatus || 'OPEN',
    supportDineIn: values.supportDineIn ?? true, supportTakeout: values.supportTakeout ?? true,
    storeArea: parseNumber(values.storeArea), tableCount: parseNumber(values.tableCount),
    receiptFooter: optionalStr(values.receiptFooter),
  };
}

export function storeCreateInput(values: StoreValues, organizationId: string, workspace: StoreWorkspace) {
  return {
    ...storeUpdateInput(values),
    code: trimStr(values.code) || generateEntityCode('STR'),
    organizationId,
    lifecycle: workspace === 'HEADQUARTERS' ? 'ACTIVE' as const : 'DRAFT' as const,
  };
}

export function storeBusinessStatusInput(open: boolean) {
  return { businessStatus: open ? 'OPEN' as StoreBusinessStatus : 'CLOSED' as StoreBusinessStatus };
}

function safeText(value?: string | null): string { return value ?? ''; }
function safeBool(value?: boolean | null): boolean { return value ?? true; }
function safeNumberText(value?: number | string | null): string { return value == null ? '' : String(value); }

export function storeValuesFromRecord(record?: StoreRecord): StoreValues {
  if (!record) return defaultStoreValues();
  return {
    name: record.name,
    code: record.code ?? undefined,
    contactPhone: safeText(record.contactPhone), managerName: safeText(record.managerName),
    managerPhone: safeText(record.managerPhone), province: safeText(record.province),
    city: safeText(record.city), district: safeText(record.district), address: safeText(record.address),
    businessHours: record.businessHours || '09:00 - 22:00', businessStatus: record.businessStatus || 'OPEN',
    supportDineIn: safeBool(record.supportDineIn), supportTakeout: safeBool(record.supportTakeout),
    storeArea: safeNumberText(record.storeArea), tableCount: safeNumberText(record.tableCount),
    receiptFooter: safeText(record.receiptFooter),
  };
}
