import { useEffect, useState, type ReactNode } from 'react';
import { Navigate, useLocation } from 'react-router';
import { Button } from '@/components/ui/button';
import { getSystemInitializationStatus } from '@/features/auth/api/systemInitialization';
import { InitializationPage } from '@/features/auth/pages/InitializationPage';
import { AuthPageShell } from './AuthPageShell';
import {
  authPath, currentDestination, loginPath, returnDestinationFromSearch, safeReturnDestination,
} from '../pages/authFlow';
import { authMessages } from '../pages/authMessages';

type InitializationPhase = 'checking' | 'initialized' | 'uninitialized' | 'error';
type LocationParts = Pick<Location, 'pathname' | 'search' | 'hash'>;

/** 根据权威初始化状态决定是否切换认证入口。 */
export function initializationRedirect(initialized: boolean, location: LocationParts): string | undefined {
  const isInitializationPage = location.pathname === '/admin/initialize';
  const existingDestination = returnDestinationFromSearch(location.search);
  if (initialized) return isInitializationPage ? loginPath(existingDestination) : undefined;
  if (isInitializationPage) return undefined;
  const destination = existingDestination ?? safeReturnDestination(currentDestination(location));
  return authPath('/admin/initialize', destination);
}

/** 初始化状态读取失败时保持关闭，并允许用户显式重试。 */
export function InitializationStatusError({ onRetry }: { onRetry(): void }) {
  return (
    <AuthPageShell title={authMessages.initialization.checkFailed}>
      <Button className="h-11 w-full" onClick={onRetry} type="button">
        {authMessages.initialization.retry}
      </Button>
    </AuthPageShell>
  );
}

/** 在挂载 Viewer 查询前确认系统已经初始化。 */
export function SystemInitializationGate({ children }: { children: ReactNode }) {
  const location = useLocation();
  const [phase, setPhase] = useState<InitializationPhase>('checking');
  const [attempt, setAttempt] = useState(0);

  useEffect(() => {
    const controller = new AbortController();
    setPhase('checking');
    void getSystemInitializationStatus(controller.signal)
      .then((initialized) => setPhase(initialized ? 'initialized' : 'uninitialized'))
      .catch(() => { if (!controller.signal.aborted) setPhase('error'); });
    return () => controller.abort();
  }, [attempt]);

  if (phase === 'checking') {
    return <div className="p-6 text-sm text-muted-foreground">{authMessages.initialization.checking}</div>;
  }
  if (phase === 'error') return <InitializationStatusError onRetry={() => setAttempt((value) => value + 1)} />;
  const redirect = initializationRedirect(phase === 'initialized', location);
  if (redirect) return <Navigate to={redirect} replace />;
  if (phase === 'uninitialized') return <InitializationPage />;
  return children;
}
