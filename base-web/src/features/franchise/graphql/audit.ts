import { gql } from '@/__generated__';

export const FRANCHISE_AUDIT_LOGS_QUERY = gql(`
  query FranchiseAuditLogs($page: Int!, $pageSize: Int!, $q: String, $filter: AuditLogFilterType) {
    auditLogs(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {
      data { id action resourceType resourceId resultCode actorAccountId storeId createdAt }
      total current_page per_page total_page
    }
  }
`);
