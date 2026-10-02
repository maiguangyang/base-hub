/** 客户端导航到未登记路径时的兜底页。渲染在标签栈之外，无标签上下文。
 *  浏览器直接输入未枚举的 URL 时由宿主返回 404，请求到不了这里。 */
export function AdminNotFoundPage() {
  return (
    <div className="flex h-full flex-col items-center justify-center gap-3 p-6 text-center">
      <p className="text-lg font-semibold">页面不存在</p>
      <p className="text-sm text-muted-foreground">该路径未在后台菜单中登记。</p>
      <a className="text-sm text-primary underline" href="/admin">返回仪表盘</a>
    </div>
  );
}
