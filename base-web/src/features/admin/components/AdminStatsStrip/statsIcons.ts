import { CheckCircle2, ListChecks, Timer, UserCheck, Users, type LucideIcon } from 'lucide-react';

/** 统计项图标键到 lucide 组件的映射。与 navIcons 同理，
 *  让调用方只传字符串，避免业务配置对象依赖 React。 */
export const statsIcons: Record<string, LucideIcon> = {
  users: Users,
  'user-check': UserCheck,
  'list-checks': ListChecks,
  timer: Timer,
  'check-circle': CheckCircle2,
};
