import { useEffect, useState } from 'react';
import { Monitor, Moon, Sun } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { DropdownMenu, DropdownMenuContent, DropdownMenuRadioGroup, DropdownMenuRadioItem, DropdownMenuTrigger } from '@/components/ui/dropdown-menu';

/** 同一浏览器在官网与后台共享的主题偏好。键名为既有契约，首帧脚本依赖它。 */
const THEME_STORAGE_KEY = 'korean-theme';

/** 主题模式仅表达用户选择，system 的实际颜色由系统媒体查询决定。 */
export type ThemeMode = 'light' | 'dark' | 'system';

/** 三态主题切换控件。官网与后台共用同一偏好存储与同一套事件契约。 */
export interface ThemeToggleProps {
  /** 菜单文案。官网从 i18n 传入，后台传入中文常量。 */
  labels: { label: string; light: string; dark: string; system: string };
  /** 文字方向，供 Radix 菜单在阿拉伯语下正确排布。后台恒为 ltr。 */
  dir?: 'ltr' | 'rtl';
}

export function ThemeToggle({ labels, dir = 'ltr' }: ThemeToggleProps) {
  const [theme, setTheme] = useState<ThemeMode>('system');

  useEffect(() => {
    try {
      const saved = localStorage.getItem(THEME_STORAGE_KEY);
      if (saved === 'light' || saved === 'dark' || saved === 'system') setTheme(saved);
    } catch {
      // 禁用本地存储时仍允许当前页面切换主题。
    }
  }, []);

  /** 将用户选择同步到首帧主题脚本和后续页面访问。 */
  function selectTheme(mode: ThemeMode) {
    setTheme(mode);
    try {
      localStorage.setItem(THEME_STORAGE_KEY, mode);
    } catch {
      // 本次选择仍由布局脚本应用，但无法跨刷新持久化。
    }
    window.dispatchEvent(new CustomEvent('korean-theme-change', { detail: mode }));
  }

  const options = [
    { mode: 'light', label: labels.light, Icon: Sun },
    { mode: 'dark', label: labels.dark, Icon: Moon },
    { mode: 'system', label: labels.system, Icon: Monitor },
  ] as const;

  return (
    <DropdownMenu dir={dir}>
      <DropdownMenuTrigger asChild>
        <Button type="button" variant="outline" aria-label={labels.label}>
          {theme === 'dark' ? <Moon aria-hidden="true" /> : theme === 'light' ? <Sun aria-hidden="true" /> : <Monitor aria-hidden="true" />}
          <span className="hidden sm:inline">{labels[theme]}</span>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        <DropdownMenuRadioGroup value={theme} onValueChange={(value) => {
          if (value === 'light' || value === 'dark' || value === 'system') selectTheme(value);
        }}>
          {options.map(({ mode, label, Icon }) => (
            <DropdownMenuRadioItem key={mode} value={mode}>
              <Icon aria-hidden="true" />
              {label}
            </DropdownMenuRadioItem>
          ))}
        </DropdownMenuRadioGroup>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
