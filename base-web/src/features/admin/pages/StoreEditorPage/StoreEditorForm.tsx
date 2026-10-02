import { useState, type SubmitEvent } from 'react';
import type { AdminStoreEditorQuery } from '@/__generated__/graphql';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { AdminPageHeader } from '@/features/admin/components/AdminPageHeader';
import { StoreBasicFields, StoreFacilityFields, StoreLocationFields, StoreOperationFields } from './StoreFormSections';
import { storeValuesFromRecord, type StoreValues } from './storeInputs';
import type { StoreWorkspace } from './storeEditorPolicy';
import { StoreDocumentCard, type StoreDocumentChange } from './StoreDocumentCard';
import { useStoreDocumentSave, type DocumentChanges, type DocumentKind } from './useStoreDocumentSave';

export { storeSaveErrorMessage } from './useStoreDocumentSave';

type ExistingStore = NonNullable<AdminStoreEditorQuery['store']>;

interface StoreEditorFormProps {
  workspace: StoreWorkspace;
  organizationId: string;
  existing?: ExistingStore;
  onBack(): void;
  onSaved(): void;
}

function StoreEditorFields({ existing, values, change, partial }: {
  existing?: ExistingStore; values: StoreValues; partial: boolean;
  change<K extends keyof StoreValues>(field: K, value: StoreValues[K]): void;
}) {
  return <fieldset disabled={partial} className="flex min-w-0 max-w-5xl flex-col gap-5">
    {existing?.lifecycle === 'REJECTED' && existing.rejectionReason
      ? <p role="alert" className="rounded-md border border-destructive/30 bg-destructive/5 p-3 text-sm text-destructive">审核退回原因：{existing.rejectionReason}</p> : null}
    <StoreBasicFields values={values} change={change} />
    <StoreLocationFields values={values} change={change} />
    <StoreOperationFields values={values} change={change} />
    <StoreFacilityFields values={values} change={change} />
  </fieldset>;
}

function StoreEditorDocumentSlots({ existing, documents, changeDocument, saving, partial, completed, documentErrors }: {
  existing?: ExistingStore; documents: DocumentChanges; saving: boolean; partial: boolean; completed: Set<DocumentKind>;
  documentErrors: Partial<Record<DocumentKind, string>>;
  changeDocument(kind: DocumentKind, value: StoreDocumentChange): void;
}) {
  return <aside className="flex flex-col gap-5 xl:-ml-4 xl:border-l xl:border-border xl:pl-4" aria-label="门店证照">
    <StoreDocumentCard label="营业执照" storedPath={existing?.businessLicenseImageUrl} change={documents.BUSINESS_LICENSE}
      onChange={(next) => changeDocument('BUSINESS_LICENSE', next)} disabled={saving || (partial && completed.has('BUSINESS_LICENSE'))}
      saveError={documentErrors.BUSINESS_LICENSE} />
    <StoreDocumentCard label="其他" storedPath={existing?.otherDocumentImageUrl} change={documents.OTHER}
      onChange={(next) => changeDocument('OTHER', next)} disabled={saving || (partial && completed.has('OTHER'))}
      saveError={documentErrors.OTHER} />
  </aside>;
}

export function StoreEditorForm(props: StoreEditorFormProps) {
  const { workspace, existing, onBack } = props;
  const [values, setValues] = useState<StoreValues>(() => storeValuesFromRecord(existing));
  const [documents, setDocuments] = useState<DocumentChanges>({ BUSINESS_LICENSE: {}, OTHER: {} });
  const { save, saving, partial, documentErrors, completed, clearDocumentError } = useStoreDocumentSave(props);
  const change = <K extends keyof StoreValues>(field: K, value: StoreValues[K]) => setValues((current) => ({ ...current, [field]: value }));
  const changeDocument = (kind: DocumentKind, value: StoreDocumentChange) => {
    clearDocumentError(kind);
    setDocuments((current) => ({ ...current, [kind]: value }));
  };
  function handleSubmit(event: SubmitEvent<HTMLFormElement>) {
    event.preventDefault();
    void save(values, documents);
  }

  const title = existing ? '编辑门店' : workspace === 'HEADQUARTERS' ? '新增直营店' : '申请新门店';
  const description = workspace === 'HEADQUARTERS'
    ? '完善门店资料，创建后立即启用。' : '完善门店资料，保存后可提交总部审核。';
  return <div className="flex flex-col gap-5">
    <AdminPageHeader title={title} description={description} actions={<>
      <Button type="button" variant="outline" onClick={onBack} disabled={saving}>返回门店</Button>
      <Button type="submit" form="store-editor-form" disabled={saving}>保存</Button>
    </>} />
    <Card className="w-full rounded-lg"><CardContent>
      <form id="store-editor-form" className="flex w-full flex-col gap-5" onSubmit={handleSubmit}>
        <div className="grid w-full gap-8 xl:grid-cols-[minmax(0,1fr)_minmax(20rem,32rem)]">
          <StoreEditorFields existing={existing} values={values} change={change} partial={partial || saving} />
          <StoreEditorDocumentSlots existing={existing} documents={documents} changeDocument={changeDocument}
            saving={saving} partial={partial} completed={completed} documentErrors={documentErrors} />
        </div>
      </form>
    </CardContent></Card>
  </div>;
}
