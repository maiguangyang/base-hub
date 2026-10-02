import { useEffect, useRef, useState } from 'react';
import { useMutation } from '@apollo/client/react';
import { KeyRound, LoaderCircle, LogOut, Sparkles } from 'lucide-react';
import { useNavigate } from 'react-router';
import { ThemeToggle } from '@/components/shared/ThemeToggle';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';
import { Breadcrumb, BreadcrumbItem, BreadcrumbLink, BreadcrumbList, BreadcrumbPage, BreadcrumbSeparator } from '@/components/ui/breadcrumb';
import { Separator } from '@/components/ui/separator';
import { SidebarTrigger } from '@/components/ui/sidebar';
import { findAdminNavItem } from '@/features/admin/config/navigation';
import { useActiveTabPath } from '@/features/admin/hooks/useAdminTabs';
import { LOGOUT_MUTATION } from '@/features/auth/graphql/auth';
import { useAuthStore, type ViewerSummary } from '@/features/auth/store/authStore';
import { getGraphQLRuntime } from '@/lib/graphql/client';
import { reportNonBlockingError } from '@/lib/diagnostics';
import { resetAdminWorkspace } from '@/stores/slices/adminTabsSlice';
import { AdminPasswordChangeDialog } from './AdminPasswordChangeDialog';

/** 后台顶栏文案。后台只走简体中文，不接入 i18n。 */
const themeLabels = { label: '外观', light: '浅色', dark: '深色', system: '跟随系统' };

/** 后台顶栏。面包屑的数据源是标签 store 的 activePath 而非 useLocation——
 *  后者在 keep-alive 栈中不可信，且共用 activePath 能保证面包屑、标签栏、
 *  侧边栏选中态三者始终一致。 */
export function AdminHeader({ onOpenAi }: { onOpenAi: () => void }) {
  const activePath = useActiveTabPath();
  const current = findAdminNavItem(activePath);

  return (
    <header className="flex h-14 shrink-0 items-center gap-3 border-b border-border bg-card px-4">
      <SidebarTrigger />
      <Separator orientation="vertical" className="h-5" />
      <Breadcrumb>
        <BreadcrumbList>
          <BreadcrumbItem>
            <BreadcrumbLink href="/">首页</BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbPage>{current ? current.title : '未知页面'}</BreadcrumbPage>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>
      <div className="ml-auto flex items-center gap-3">
        <button type="button" onClick={onOpenAi} aria-label="打开 AI 助手" title="AI 助手" className="inline-flex size-9 items-center justify-center rounded-md text-foreground transition-colors hover:bg-accent focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring">
          <Sparkles className="size-4" aria-hidden="true" />
        </button>
        <ThemeToggle labels={themeLabels} />
        <AdminAccountMenu />
      </div>
    </header>
  );
}

/** 右上角账号菜单。注销成功前保留当前浏览器会话，避免网络错误造成假退出。 */
function AdminAccountMenu() {
  const accountMenuRef = useRef<HTMLDivElement>(null);
  const navigate = useNavigate();
  const viewer = useAuthStore((state) => state.viewer);
  const [passwordOpen, setPasswordOpen] = useState(false);
  const { displayName, secondary, initial } = accountSummary(viewer);
  const { error, loading, signOut } = useLogoutAction(() => navigate('/admin/login', { replace: true }));

  function openPasswordDialog() {
    setPasswordOpen(true);
  }

  useEffect(() => {
    function dismissFocus(event: PointerEvent) {
      const accountMenu = accountMenuRef.current;
      if (!accountMenu || accountMenu.contains(event.target as Node)) return;
      const focused = document.activeElement;
      if (focused instanceof HTMLElement && accountMenu.contains(focused)) focused.blur();
    }
    document.addEventListener('pointerdown', dismissFocus);
    return () => document.removeEventListener('pointerdown', dismissFocus);
  }, []);

  return (
    <>
      <div
        ref={accountMenuRef}
        data-account-menu
        role="group"
        aria-label={`账号操作：${displayName}`}
        tabIndex={0}
        className="group/account relative flex size-8 cursor-pointer items-center justify-center rounded-full outline-none transition-colors after:absolute after:-inset-1.5 hover:bg-accent focus-visible:ring-[3px] focus-visible:ring-ring/50"
      >
        <Avatar className="size-8">
          <AvatarFallback>{initial}</AvatarFallback>
        </Avatar>
        <AccountMenuContent
          displayName={displayName} secondary={secondary} loading={loading} error={error}
          onChangePassword={openPasswordDialog} onSignOut={signOut}
        />
      </div>
      <AdminPasswordChangeDialog open={passwordOpen} onOpenChange={setPasswordOpen} />
    </>
  );
}

interface AccountMenuContentProps {
  displayName: string;
  secondary?: string | null;
  loading: boolean;
  error?: string;
  onChangePassword(): void;
  onSignOut(): Promise<void>;
}

function AccountMenuContent(props: AccountMenuContentProps) {
  return (
    <div className="pointer-events-none invisible absolute top-full right-0 z-50 w-64 pt-2 opacity-0 transition-opacity group-hover/account:pointer-events-auto group-hover/account:visible group-hover/account:opacity-100 group-focus-within/account:pointer-events-auto group-focus-within/account:visible group-focus-within/account:opacity-100">
      <div id="admin-account-menu" role="region" aria-label="账号操作" className="rounded-md border bg-popover p-1 text-popover-foreground shadow-lg">
        <div className="px-2 py-2"><p className="truncate text-sm font-medium">{props.displayName}</p>{props.secondary && <p className="truncate text-xs text-muted-foreground">{props.secondary}</p>}</div>
        <button type="button" onClick={props.onChangePassword} className="flex w-full items-center gap-2 rounded-sm px-2 py-2 text-sm outline-none hover:bg-accent focus:bg-accent focus:text-accent-foreground">
          <KeyRound className="size-4" aria-hidden="true" />修改密码
        </button>
        <button type="button" disabled={props.loading} onClick={() => void props.onSignOut()} className="mt-1 flex w-full items-center gap-2 border-t border-border px-2 py-2 text-sm text-destructive outline-none hover:bg-destructive/10 focus:bg-destructive/10 focus:text-destructive disabled:pointer-events-none disabled:opacity-50">
          {props.loading ? <LoaderCircle className="size-4 animate-spin" aria-hidden="true" /> : <LogOut className="size-4" aria-hidden="true" />}
          {props.loading ? '正在退出…' : '退出登录'}
        </button>
        {props.error && <p role="alert" className="px-2 py-1.5 text-xs text-destructive">{props.error}</p>}
      </div>
    </div>
  );
}

function accountSummary(viewer: ViewerSummary | null) {
  const account = viewer?.account;
  const displayName = account?.displayName || '管理员';
  return {
    displayName,
    secondary: account?.email || account?.phone,
    initial: [...displayName][0] || '管',
  };
}

function useLogoutAction(onSuccess: () => void) {
  const resetAuth = useAuthStore((state) => state.reset);
  const [logout, { loading }] = useMutation(LOGOUT_MUTATION);
  const [error, setError] = useState<string>();

  async function signOut() {
    if (loading) return;
    setError(undefined);
    try {
      const result = await logout();
      if (result.data?.logout !== true) throw new Error('LOGOUT_FAILED');
    } catch {
      setError('退出失败，请重试。');
      return;
    }
    const runtime = getGraphQLRuntime({ resetAuth, resetWorkspace: resetAdminWorkspace });
    try {
      await runtime.resetSession('LOGOUT');
    } catch (cause) {
      reportNonBlockingError('后台会话清理失败，继续执行安全跳转。', cause);
    }
    onSuccess();
  }

  return { error, loading, signOut };
}
