import { CircleAlert, X } from 'lucide-react';
import { cn } from '@/lib/utils';
import type { AdminTab } from '@/stores/slices/tabOperations';

/** 单个标签。errored 为真时显示错误标识——隐藏标签的错误 fallback 不可见，
 *  不在标签栏标识，用户在切过去之前完全无感知。 */
export interface TabItemProps {
  tab: AdminTab;
  isActive: boolean;
  errored: boolean;
  onActivate: (path: string) => void;
  onClose: (path: string) => void;
}

/**
 * 渲染单个标签。
 *
 * 标签主体是原生 <button> 而非带 onClick 的 div：中台是高频键盘操作场景，
 * div 无法聚焦、无法用 Enter/Space 激活，屏幕阅读器也读不出可交互语义。
 * 关闭按钮是并列的兄弟节点而非嵌套 button——HTML 不允许 button 嵌套。
 * 外层 div 只负责布局与激活态下划线，不承担交互。
 */
export function TabItem({ tab, isActive, errored, onActivate, onClose }: TabItemProps) {
  return (
    <div
      data-tab-active={isActive}
      className={cn(
        // 用 h-full 而非 h-10：父条 h-10 含 1px 下边框，内容盒只有 39px，
        // 写死 40px 会高出 1px 并撑出滚动条。
        'group relative flex h-full shrink-0 items-center border-r border-border pr-2 text-sm',
        isActive ? 'bg-card font-medium text-primary' : 'text-muted-foreground hover:bg-accent hover:text-accent-foreground',
      )}
    >
      <button
        type="button"
        aria-current={isActive ? 'page' : undefined}
        className="flex h-full cursor-pointer items-center gap-1.5 px-3 outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-inset"
        onClick={() => onActivate(tab.path)}
        onAuxClick={(event) => {
          // 中键关闭。固定标签由 onClose 侧的 store 逻辑拒绝。
          if (event.button === 1) {
            event.preventDefault();
            onClose(tab.path);
          }
        }}
      >
        {errored ? <CircleAlert className="size-3.5 text-danger-fg" aria-label="此标签渲染失败" /> : null}
        <span>{tab.title}</span>
      </button>
      {tab.affix ? null : (
        <button
          type="button"
          aria-label={`关闭 ${tab.title}`}
          className={cn(
            'rounded-sm p-0.5 outline-none hover:bg-accent focus-visible:ring-2 focus-visible:ring-ring',
            // 非活动标签的关闭按钮平时隐藏，但获得焦点时必须可见，否则键盘用户看不到焦点在哪。
            isActive ? 'opacity-100' : 'opacity-0 focus-visible:opacity-100 group-hover:opacity-100',
          )}
          onClick={() => onClose(tab.path)}
        >
          <X className="size-3.5" />
        </button>
      )}
      {isActive ? <span className="absolute inset-x-0 bottom-0 h-0.5 bg-primary" /> : null}
    </div>
  );
}
