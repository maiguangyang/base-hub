import { gql } from '@/__generated__';

export const HQ_AUDIT_LOGS_QUERY = gql(`
  query HqAuditLogs($page: Int!, $pageSize: Int!, $q: String, $filter: AuditLogFilterType) {
    auditLogs(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {
      data { id action resourceType resourceId resultCode actorAccountId organizationId storeId createdAt }
      total current_page per_page total_page
    }
  }
`);
