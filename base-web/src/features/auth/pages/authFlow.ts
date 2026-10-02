import type { ViewerSummary } from '../store/authStore';
import { reportNonBlockingError } from '@/lib/diagnostics';

export type WorkspaceSummary = ViewerSummary['workspaces'][number];
export type PasswordValidationCode = 'PASSWORD_WEAK' | 'PASSWORD_CONFIRM_MISMATCH';

/** 改密后用替代 Session 重建整个 GraphQL runtime，避免复用已撤销连接。 */
export async function completePasswordChange(
  viewer: ViewerSummary,
  setViewer: (viewer: ViewerSummary) => void,
  disposeRuntime: () => Promise<void>,
  replaceLocation: (path: string) => void,
  returnDestination?: string,
): Promise<void> {
	await completeAuthentication(viewer, setViewer, disposeRuntime, replaceLocation, returnDestination);
}

/** 登录后销毁可能属于上一账号的 GraphQL runtime，并以新 Cookie 完整导航。 */
export async function completeLogin(
	viewer: ViewerSummary,
	setViewer: (viewer: ViewerSummary) => void,
	disposeRuntime: () => Promise<void>,
	replaceLocation: (path: string) => void,
	returnDestination?: string,
): Promise<void> {
	await completeAuthentication(viewer, setViewer, disposeRuntime, replaceLocation, returnDestination);
}

async function completeAuthentication(
	viewer: ViewerSummary,
	setViewer: (viewer: ViewerSummary) => void,
	disposeRuntime: () => Promise<void>,
	replaceLocation: (path: string) => void,
	returnDestination?: string,
): Promise<void> {
  setViewer(viewer);
  const nextPath = nextAuthPath(viewer, returnDestination);
  // 服务端认证状态已经切换，本地缓存清理只能尽力而为，不能阻断安全导航。
  try {
    await disposeRuntime();
  } catch (cause) {
    reportNonBlockingError('认证切换后的 GraphQL runtime 清理失败，继续执行安全导航。', cause);
  }
  replaceLocation(nextPath);
}

/** 根据权威 Viewer 选择固定认证路由，绝不接受服务端外部跳转地址。 */
export function nextAuthPath(viewer: ViewerSummary, returnDestination?: string): string {
  const destination = safeReturnDestination(returnDestination);
  if (viewer.account.mustChangePassword) return authPath('/admin/change-password', destination);
  const workspace = viewer.currentWorkspace;
  if (!workspace || workspace.workspaceType === 'DISCOVERY') return authPath('/admin/workspaces', destination);
  const home = workspace.workspaceType === 'HEADQUARTERS' ? '/admin/hq' : '/admin/franchise';
  return destination && destinationMatchesWorkspace(destination, home) ? destination : home;
}

function destinationMatchesWorkspace(destination: string, home: string): boolean {
  const pathname = new URL(destination, 'https://return.invalid').pathname;
  return pathname === home || pathname.startsWith(`${home}/`);
}

type LocationParts = Pick<Location, 'pathname' | 'search' | 'hash'>;

/** 将浏览器当前地址压缩为可回跳的站内路径。 */
export function currentDestination(location: LocationParts): string {
  return `${location.pathname}${location.search}${location.hash}`;
}

/** 只允许总部或加盟商工作台内部地址，拒绝外部地址和认证循环。 */
export function safeReturnDestination(value: string | null | undefined): string | undefined {
  if (!value || !value.startsWith('/') || value.startsWith('//')) return undefined;
  try {
    const base = 'https://return.invalid';
    const parsed = new URL(value, base);
    if (parsed.origin !== base || !isBusinessDestination(parsed.pathname)) return undefined;
    return `${parsed.pathname}${parsed.search}${parsed.hash}`;
  } catch {
    return undefined;
  }
}

/** 从当前认证页查询参数读取并验证回跳地址。 */
export function returnDestinationFromSearch(search: string): string | undefined {
  return safeReturnDestination(new URLSearchParams(search).get('dt'));
}

/** 在浏览器认证页读取安全回跳地址；服务端渲染时返回空值。 */
export function browserReturnDestination(): string | undefined {
  return typeof window === 'undefined' ? undefined : returnDestinationFromSearch(window.location.search);
}

/** 构造携带安全回跳地址的登录 URL。 */
export function loginPath(destination?: string): string {
  return authPath('/admin/login', destination);
}

/** 给认证流程页面附加同一个安全回跳地址。 */
export function authPath(path: string, destination?: string): string {
  const safe = safeReturnDestination(destination);
  return safe ? `${path}?${new URLSearchParams({ dt: safe })}` : path;
}

function isBusinessDestination(pathname: string): boolean {
  return pathname === '/admin/hq' || pathname.startsWith('/admin/hq/')
    || pathname === '/admin/franchise' || pathname.startsWith('/admin/franchise/');
}

/** 过滤仅用于发现阶段的 DISCOVERY 工作台。 */
export function businessWorkspaces(workspaces: readonly WorkspaceSummary[]): WorkspaceSummary[] {
  return workspaces.filter((workspace) => workspace.workspaceType !== 'DISCOVERY');
}

/** 在提交前镜像 Engine 的密码长度策略。 */
export function validateNewPassword(password: string, confirmation: string): PasswordValidationCode | undefined {
  const length = [...password].length;
  if (length < 8 || length > 20) return 'PASSWORD_WEAK';
  if (password !== confirmation) return 'PASSWORD_CONFIRM_MISMATCH';
  return undefined;
}
