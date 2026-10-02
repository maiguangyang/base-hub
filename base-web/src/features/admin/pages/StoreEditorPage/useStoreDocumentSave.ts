import { useRef, useState } from 'react';
import { useMutation } from '@apollo/client/react';
import type { AdminStoreEditorQuery } from '@/__generated__/graphql';
import { useAdminToast } from '@/features/admin/components/AdminToast';
import { REMOVE_STORE_DOCUMENT_MUTATION, SET_STORE_DOCUMENT_MUTATION } from '@/features/admin/graphql/storeEditor';
import { CREATE_FRANCHISE_STORE_MUTATION, UPDATE_FRANCHISE_STORE_MUTATION } from '@/features/franchise/graphql/stores';
import { CREATE_DIRECT_STORE_MUTATION, UPDATE_DIRECT_STORE_MUTATION } from '@/features/hq/graphql/directStores';
import { getGraphQLErrorCode } from '@/lib/graphql/errors';
import { actionErrorReason } from '@/lib/graphql/actionErrors';
import { uploadStoreDocument } from './storeDocumentImages';
import { storeCreateInput, storeUpdateInput, type StoreValues } from './storeInputs';
import type { StoreWorkspace } from './storeEditorPolicy';
import type { StoreDocumentChange } from './StoreDocumentCard';

type ExistingStore = NonNullable<AdminStoreEditorQuery['store']>;
export type DocumentKind = 'BUSINESS_LICENSE' | 'OTHER';
export type DocumentChanges = Record<DocumentKind, StoreDocumentChange>;
const kinds: DocumentKind[] = ['BUSINESS_LICENSE', 'OTHER'];

class DocumentStepError extends Error {
  constructor(readonly kind: DocumentKind, message: string) { super(message); }
}

type Progress = {
  saved: boolean; id?: string; completed: Set<DocumentKind>;
  attachments: Map<DocumentKind, { file: File; id: string }>;
};

export function storeSaveErrorMessage(code?: string): string {
  if (code === 'PERMISSION_DENIED' || code === 'STORE_SCOPE_DENIED') return '当前账号无权维护该门店。';
  if (code === 'CONFLICT') return '门店状态已变化，请返回列表刷新后重试。';
  if (code === 'VALIDATION_FAILED') return '门店资料不符合要求，请检查后重试。';
  if (code) return `门店保存失败：${actionErrorReason({ extensions: { code } })}`;
  return '门店保存失败，请重试。';
}

function saveFailureMessage(progress: Progress, cause: unknown): string | undefined {
  if (progress.saved) {
    return progress.id
      ? '门店已保存，证照未完成，请重试；不会重复新增门店。'
      : '门店可能已保存，但未获得门店编号。请返回列表核对，避免重复新增。';
  }
  if (cause instanceof DocumentStepError) return undefined;
  return storeSaveErrorMessage(getGraphQLErrorCode(cause));
}

interface Props { workspace: StoreWorkspace; organizationId: string; existing?: ExistingStore; onSaved(): void }

function useStoreWrite({ workspace, organizationId, existing }: Props) {
  const [createHq] = useMutation(CREATE_DIRECT_STORE_MUTATION, { context: { adminFeedback: false } });
  const [updateHq] = useMutation(UPDATE_DIRECT_STORE_MUTATION, { context: { adminFeedback: false } });
  const [createFranchise] = useMutation(CREATE_FRANCHISE_STORE_MUTATION, { context: { adminFeedback: false } });
  const [updateFranchise] = useMutation(UPDATE_FRANCHISE_STORE_MUTATION, { context: { adminFeedback: false } });
  return async (values: StoreValues): Promise<string | undefined> => {
    if (existing) {
      const variables = { id: existing.id, input: storeUpdateInput(values) };
      await (workspace === 'HEADQUARTERS' ? updateHq({ variables }) : updateFranchise({ variables }));
      return existing.id;
    }
    const variables = { input: storeCreateInput(values, organizationId, workspace) };
    const result = await (workspace === 'HEADQUARTERS' ? createHq({ variables }) : createFranchise({ variables }));
    return result.data?.createStore?.id;
  };
}

function useDocumentWrite() {
  const [bind] = useMutation(SET_STORE_DOCUMENT_MUTATION, { context: { adminFeedback: false } });
  const [remove] = useMutation(REMOVE_STORE_DOCUMENT_MUTATION, { context: { adminFeedback: false } });
  return async (storeId: string, kind: DocumentKind, attachmentId?: string) => {
    if (attachmentId) await bind({ variables: { storeId, kind, attachmentId } });
    else await remove({ variables: { storeId, kind } });
  };
}

async function uploadChanged(changes: DocumentChanges, progress: Progress) {
  for (const kind of kinds) {
    const file = changes[kind].file;
    if (!file || progress.completed.has(kind) || progress.attachments.get(kind)?.file === file) continue;
    try {
      progress.attachments.set(kind, { file, id: await uploadStoreDocument(file) });
    } catch (cause) {
      throw new DocumentStepError(kind, cause instanceof Error ? cause.message : '证照上传失败，请重试。');
    }
  }
}

async function applyChanged(changes: DocumentChanges, progress: Progress,
  write: (storeId: string, kind: DocumentKind, attachmentId?: string) => Promise<void>) {
  for (const kind of kinds) {
    if (progress.completed.has(kind)) continue;
    const change = changes[kind];
    if (!change.file && !change.remove) continue;
    if (!progress.id) throw new Error('STORE_ID_UNAVAILABLE');
    try {
      if (change.file) await bindChangedFile(progress, kind, change.file, write);
      else await write(progress.id, kind);
    } catch (cause) {
      if (cause instanceof DocumentStepError) throw cause;
      throw new DocumentStepError(kind, getGraphQLErrorCode(cause) === 'PERMISSION_DENIED'
        ? '当前账号无权维护该门店。' : '证照保存失败，请重试。');
    }
    progress.completed.add(kind);
  }
}

async function bindChangedFile(progress: Progress, kind: DocumentKind, file: File,
  write: (storeId: string, kind: DocumentKind, attachmentId?: string) => Promise<void>) {
  const storeId = progress.id!;
  const previous = progress.attachments.get(kind);
  let attachmentId = previous?.file === file ? previous.id : undefined;
  if (!attachmentId) {
    attachmentId = await uploadForBind(file, kind);
    progress.attachments.set(kind, { file, id: attachmentId });
  }
  try {
    await write(storeId, kind, attachmentId);
  } catch (cause) {
    if (getGraphQLErrorCode(cause) !== 'VALIDATION_FAILED') throw cause;
    const refreshedId = await uploadForBind(file, kind);
    progress.attachments.set(kind, { file, id: refreshedId });
    await write(storeId, kind, refreshedId);
  }
}

async function uploadForBind(file: File, kind: DocumentKind): Promise<string> {
  try {
    return await uploadStoreDocument(file);
  } catch (cause) {
    throw new DocumentStepError(kind, cause instanceof Error ? cause.message : '证照上传失败，请重试。');
  }
}

export function useStoreDocumentSave(props: Props) {
  const showToast = useAdminToast();
  const [busy, setBusy] = useState(false);
  const [partial, setPartial] = useState(false);
  const [documentErrors, setDocumentErrors] = useState<Partial<Record<DocumentKind, string>>>({});
  const pending = useRef(false);
  const progress = useRef<Progress>({ saved: false, id: props.existing?.id, completed: new Set(), attachments: new Map() });
  const writeStore = useStoreWrite(props);
  const writeDocument = useDocumentWrite();
  async function save(values: StoreValues, changes: DocumentChanges) {
    if (pending.current) return;
    if (!values.businessHours?.trim()) { showToast('error', '请选择营业时间。'); return; }
    pending.current = true; setBusy(true); setDocumentErrors({});
    try {
      if (!progress.current.saved) {
        await uploadChanged(changes, progress.current);
        progress.current.id = await writeStore(values);
        progress.current.saved = true;
      }
      if (!progress.current.id) {
        throw new Error('STORE_ID_UNAVAILABLE');
      }
      await applyChanged(changes, progress.current, writeDocument);
      props.onSaved();
      showToast('success', '门店保存成功。');
    } catch (cause) {
      setPartial(progress.current.saved);
      if (cause instanceof DocumentStepError) setDocumentErrors({ [cause.kind]: cause.message });
      showToast('error', saveFailureMessage(progress.current, cause)
        ?? (cause instanceof Error ? cause.message : '门店保存失败，请重试。'));
    } finally { pending.current = false; setBusy(false); }
  }
  function clearDocumentError(kind: DocumentKind) {
    setDocumentErrors((current) => ({ ...current, [kind]: undefined }));
  }
  return { save, saving: busy, partial, documentErrors, completed: progress.current.completed, clearDocumentError };
}
