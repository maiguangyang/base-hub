import { ArrowUpRight, Languages } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { entryPath, type Locale, type Messages } from '@/i18n';

/** 官网骨架只接收已解析的语言与文案，不访问后台数据。 */
export interface HomePageProps {
  locale: Locale;
  messages: Messages['publicHome'];
}

/** 面向访客的首页占位内容。 */
export function HomePage({ locale, messages }: HomePageProps) {
  return (
    <div className="mx-auto grid min-h-[calc(100vh-9rem)] max-w-6xl items-center gap-12 px-5 py-16 lg:grid-cols-[1.2fr_0.8fr] lg:px-8">
      <section className="space-y-8">
        <p className="inline-flex rounded-full border border-primary/25 bg-primary/10 px-4 py-1.5 text-sm font-medium text-primary">{messages.eyebrow}</p>
        <div className="space-y-5">
          <h1 className="max-w-3xl text-5xl leading-tight font-semibold tracking-tight sm:text-6xl">{messages.title}</h1>
          <p className="max-w-2xl text-lg leading-relaxed text-muted-foreground">{messages.description}</p>
        </div>
        <Button asChild size="default">
          <a href={entryPath(locale, 'admin')}>{messages.action}<ArrowUpRight aria-hidden="true" /></a>
        </Button>
      </section>
      <aside className="relative overflow-hidden rounded-3xl border bg-card p-8 text-card-foreground shadow-xl shadow-primary/5 sm:p-10">
        <div className="absolute -end-16 -top-16 size-48 rounded-full bg-primary/10 blur-3xl" aria-hidden="true" />
        <div className="relative flex min-h-72 flex-col justify-between">
          <div className="flex size-14 items-center justify-center rounded-2xl bg-primary text-primary-foreground"><Languages className="size-7" aria-hidden="true" /></div>
          <div>
            <p className="mb-3 text-sm font-medium text-muted-foreground">{messages.status}</p>
            <p className="text-2xl font-semibold tracking-tight">Korean Hub</p>
          </div>
        </div>
      </aside>
    </div>
  );
}
