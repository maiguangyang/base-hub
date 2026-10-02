/**
 * 生成系统内部业务实体编码（用于门店、组织等，用户无需手动填写或维护）。
 * 格式：前缀（可选）+ 36进制时间戳 + 随机大写字母数字，总长度保证在 2~32 字符之间。
 */
export function generateEntityCode(prefix: string = ''): string {
  const timestamp = Date.now().toString(36).toUpperCase();
  const random = Math.random().toString(36).slice(2, 8).toUpperCase();
  return `${prefix}${timestamp}${random}`.slice(0, 32);
}
