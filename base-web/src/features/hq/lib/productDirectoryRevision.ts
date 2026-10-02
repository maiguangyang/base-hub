import { useSyncExternalStore } from 'react';
import type { ApolloClient } from '@apollo/client';

type Directory = 'categories' | 'brands' | 'specifications' | 'specificationValues' | 'packageTemplates';

const fields: Record<Directory, string> = {
  categories: 'hqProductCategories', brands: 'hqProductBrands',
  specifications: 'hqSpecifications', specificationValues: 'hqSpecificationValues', packageTemplates: 'hqProductPackageTemplates',
};
const revisions: Record<Directory, number> = { categories: 0, brands: 0, specifications: 0, specificationValues: 0, packageTemplates: 0 };
const listeners = new Set<() => void>();

function subscribe(listener: () => void) {
  listeners.add(listener);
  return () => { listeners.delete(listener); };
}

export function useProductDirectoryRevision(directory: Directory) {
  return useSyncExternalStore(subscribe, () => revisions[directory], () => 0);
}

export function invalidateProductDirectory(client: ApolloClient, directory: Directory) {
  client.cache.evict({ id: 'ROOT_QUERY', fieldName: fields[directory] });
  client.cache.gc();
  revisions[directory] += 1;
  listeners.forEach((listener) => listener());
}
