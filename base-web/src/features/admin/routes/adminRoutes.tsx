import { Navigate, useLocation, type RouteObject } from "react-router";
import type { BusinessWorkspaceType } from "../config/navigation";
import { AdminNotFoundPage } from "../pages/AdminNotFoundPage";
import { AuditLogListPage } from "@/features/hq/pages/AuditLogListPage";
import { HqStoreListPage } from "@/features/hq/pages/HqStoreListPage/HqStoreListPage";
import { FranchiseListPage } from "@/features/hq/pages/FranchiseListPage";
import { HqDashboardPage } from "@/features/hq/pages/HqDashboardPage";
import { ModelConfigPage } from "@/features/hq/pages/ModelConfigPage";
import { PaymentConfigPage } from "@/features/hq/pages/PaymentConfigPage";
import { StoreApprovalListPage } from "@/features/hq/pages/StoreApprovalListPage";
import { AdministratorListPage } from "@/features/hq/pages/AdministratorListPage";
import { RoleListPage as HqRoleListPage } from "@/features/hq/pages/RoleListPage";
import { AuditLogListPage as FranchiseAuditLogListPage } from "@/features/franchise/pages/AuditLogListPage";
import { FranchiseDashboardPage } from "@/features/franchise/pages/FranchiseDashboardPage";
import { RoleListPage } from "@/features/franchise/pages/RoleListPage";
import { StaffListPage } from "@/features/franchise/pages/StaffListPage";
import { StoreListPage } from "@/features/franchise/pages/StoreListPage";
import { StoreEditorPage } from "@/features/admin/pages/StoreEditorPage/StoreEditorPage";

const hqRoutes: RouteObject[] = [
  { path: "/admin/hq", element: <HqDashboardPage /> },
  { path: "/admin/hq/stores", element: <HqStoreListPage /> },
  { path: "/admin/hq/stores/manage", element: <StoreEditorPage workspace="HEADQUARTERS" /> },
  { path: "/admin/hq/direct-stores", element: <LegacyHqStoreRedirect type="HEADQUARTERS" /> },
  { path: "/admin/hq/franchises", element: <FranchiseListPage /> },
  { path: "/admin/hq/franchise-stores", element: <LegacyHqStoreRedirect type="FRANCHISE" /> },
  { path: "/admin/hq/store-approvals", element: <StoreApprovalListPage /> },
  { path: "/admin/hq/administrators", element: <AdministratorListPage /> },
  { path: "/admin/hq/roles", element: <HqRoleListPage /> },
  { path: "/admin/hq/audit", element: <AuditLogListPage /> },
  { path: "/admin/hq/ai-model", element: <ModelConfigPage /> },
  { path: "/admin/hq/payment-config", element: <PaymentConfigPage /> },
];

function LegacyHqStoreRedirect({ type }: { type: "HEADQUARTERS" | "FRANCHISE" }) {
  const { search } = useLocation();
  const params = new URLSearchParams(search);
  params.set("type", type);
  return <Navigate to={`/admin/hq/stores?${params}`} replace />;
}

const franchiseRoutes: RouteObject[] = [
  { path: "/admin/franchise", element: <FranchiseDashboardPage /> },
  { path: "/admin/franchise/stores", element: <StoreListPage /> },
  { path: "/admin/franchise/stores/manage", element: <StoreEditorPage workspace="FRANCHISE" /> },
  { path: "/admin/franchise/staff", element: <StaffListPage /> },
  { path: "/admin/franchise/roles", element: <RoleListPage /> },
  { path: "/admin/franchise/audit", element: <FranchiseAuditLogListPage /> },
];

/** 返回单一工作台路由，避免 keep-alive 栈挂载另一租户的页面。 */
export function routesForWorkspace(workspaceType: BusinessWorkspaceType): RouteObject[] {
  return workspaceType === "HEADQUARTERS" ? hqRoutes : franchiseRoutes;
}

export const adminRoutes: RouteObject[] = [
  ...hqRoutes,
  ...franchiseRoutes,
  { path: "*", element: <AdminNotFoundPage /> },
];
