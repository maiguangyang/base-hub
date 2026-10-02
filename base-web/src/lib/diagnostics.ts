/** 记录不应阻断用户安全流程的异常；诊断通道自身失败时保持静默。 */
export function reportNonBlockingError(message: string, cause: unknown): void {
  try {
    console.error(message, cause);
  } catch {
    // 日志设施不能反过来破坏退出、改密等终态流程。
  }
}
