import { Component, type ErrorInfo, type ReactNode } from 'react';
import { useAdminTabsStore } from '@/stores/slices/adminTabsSlice';

/** 标签级错误隔离的入参。tabPath 用于在 store 中留痕并驱动重试。 */
export interface AdminTabErrorBoundaryProps {
  tabPath: string;
  children: ReactNode;
}

interface AdminTabErrorBoundaryState {
  error: Error | null;
}

/** 单个标签的错误隔离。**必须是 class 组件**——React 错误边界只能由
 *  getDerivedStateFromError / componentDidCatch 实现，没有 hooks 等价物；
 *  写成函数组件不会捕获任何错误，且编译通过、无警告、无运行时报错。
 *  不隔离的话，隐藏标签抛出的异常会连带整个外壳白屏，
 *  而用户无从判断是哪个标签出的问题。 */
export class AdminTabErrorBoundary extends Component<AdminTabErrorBoundaryProps, AdminTabErrorBoundaryState> {
  state: AdminTabErrorBoundaryState = { error: null };

  static getDerivedStateFromError(error: Error): AdminTabErrorBoundaryState {
    return { error };
  }

  componentDidCatch(error: Error, info: ErrorInfo): void {
    useAdminTabsStore.getState().markTabError(this.props.tabPath);
    console.error(`后台标签 ${this.props.tabPath} 渲染失败`, error, info.componentStack);
  }

  /** 重试复用刷新机制：version 递增改变 React key，整棵子树连同本边界一起重建，
   *  因此不需要单独的边界重置逻辑。 */
  private readonly handleRetry = (): void => {
    useAdminTabsStore.getState().refreshTab(this.props.tabPath);
  };

  render(): ReactNode {
    const { error } = this.state;
    if (!error) return this.props.children;
    return (
      <div className="flex h-full flex-col items-center justify-center gap-4 p-6 text-center">
        <p className="text-sm font-medium text-danger-fg">此标签页渲染失败</p>
        <p className="max-w-lg text-sm text-muted-foreground">{error.message}</p>
        <button
          type="button"
          onClick={this.handleRetry}
          className="h-9 rounded-md border border-border px-4 text-sm hover:bg-accent"
        >
          重试
        </button>
      </div>
    );
  }
}
