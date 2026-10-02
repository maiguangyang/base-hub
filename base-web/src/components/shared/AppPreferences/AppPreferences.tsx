import { Globe2 } from 'lucide-react';
import { ThemeToggle } from '@/components/shared/ThemeToggle';
import { Button } from '@/components/ui/button';
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from '@/components/ui/dropdown-menu';
import { entryPath, localeCodes, localeNames, type Locale, type Messages } from '@/i18n';

/** 官网头部的偏好控件。文案与方向均由 URL 语言决定。
 *  后台不再使用本组件——它只需要主题切换，直接用共享的 ThemeToggle。 */
export interface AppPreferencesProps {
  locale: Locale;
  messages: Messages['common'];
}

/** 官网的语言与主题控件。 */
export function AppPreferences({ locale, messages }: AppPreferencesProps) {
  const direction = locale === 'ar-SA' ? 'rtl' : 'ltr';

  return (
    <div className="flex items-center gap-2">
      <LanguageMenu locale={locale} label={messages.language} direction={direction} />
      <ThemeToggle labels={messages.theme} dir={direction} />
    </div>
  );
}

/** 语言菜单只切换官网语言；后台恒为简体中文，不出现在此列表的目标中。 */
function LanguageMenu({ locale, label, direction }: Pick<AppPreferencesProps, 'locale'> & { label: string; direction: 'ltr' | 'rtl' }) {
  return (
    <DropdownMenu dir={direction}>
      <DropdownMenuTrigger asChild>
        <Button type="button" variant="outline" aria-label={label}>
          <Globe2 aria-hidden="true" />
          <span className="hidden sm:inline">{localeNames[locale]}</span>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        {localeCodes.map((code) => (
          <DropdownMenuItem key={code} asChild>
            <a href={entryPath(code, 'home')} lang={code} aria-current={code === locale ? 'page' : undefined}>
              {localeNames[code]}
            </a>
          </DropdownMenuItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
