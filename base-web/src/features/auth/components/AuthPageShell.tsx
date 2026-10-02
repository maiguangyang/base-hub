import type { ReactNode } from 'react';

interface AuthPageShellProps {
  title: string;
  description?: string;
  children: ReactNode;
}

/** 认证流程共享的聚焦式页面外壳。 */
export function AuthPageShell({ title, description, children }: AuthPageShellProps) {
  return (
    <main className="flex min-h-screen items-center justify-center bg-muted/30 p-4">
      <section className="w-full max-w-md rounded-xl border border-border bg-card p-7 shadow-sm sm:p-8">
        <div className="mb-8 space-y-2">
          <h1 className="text-2xl font-semibold tracking-tight text-foreground">{title}</h1>
          {description && <p className="text-sm text-muted-foreground">{description}</p>}
        </div>
        {children}
      </section>
    </main>
  );
}
