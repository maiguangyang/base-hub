export interface PermissionOption {
  id: string;
  name: string;
  action: string;
  module: string;
  scope: 'SYSTEM' | 'TENANT';
}

// CRUD 名称对应 korean-engine/model/model.graphql 的 @entity(title)。
const entityTitles: Record<string, string> = {
  account: '账号', organization: '加盟商', operatorMembership: '加盟商成员', permission: '权限',
  operatorRole: '角色', store: '门店', session: '登录会话', membershipInvitation: '成员邀请',
  auditLog: '审计日志', franchiseOpeningRecord: '加盟商开通记录',
  globalPaymentConfig: '全局支付配置', franchisePaymentConfig: '加盟商支付配置',
  storePaymentConfig: '门店支付配置', customerMember: '顾客会员',
  customerBenefitPolicy: '顾客权益规则', customerDailyPointGrantBudget: '顾客积分日额度',
  customerPointEntry: '顾客积分流水', customerCouponTemplate: '顾客优惠券模板',
  customerCouponGrant: '顾客优惠券', productCategory: '商品分类', productBrand: '商品品牌',
  product: '商品', productSku: '商品 SKU', productPackage: '商品包装',
  specificationDefinition: '商品规格', specificationValue: '商品规格值',
  productSpecificationChoice: '商品已选规格值', productSkuSpecificationValue: 'SKU 规格值',
  productPackageTemplate: '包装换算模板', storeListing: '门店经营规格',
  storePackageOffer: '门店包装销售项', storePriceRevision: '门店价格历史',
  storeInventoryBatch: '门店库存批次', storeStockBalance: '门店包装库存余额',
  storeStockMovement: '门店库存流水', storeStocktake: '门店盘点单',
  storeStocktakeLine: '门店盘点行', storePromotion: '门店促销活动',
  storePromotionTarget: '促销活动适用包装',
};

const domainTitles: Record<string, string> = {
  franchise: '加盟商开通', hqRole: '总部角色', hqMembership: '管理员',
  hqStore: '直营门店', customer: '会员', report: '经营报表',
  aiModelConfig: 'AI 模型配置', paymentConfig: '支付配置',
  hqCustomer: '顾客会员', hqCustomerPolicy: '顾客权益规则',
  hqCustomerPoints: '顾客积分', hqCustomerCoupon: '顾客优惠券',
  hqProductCatalog: '商品目录', tenantAudit: '加盟商审计日志',
  franchiseProduct: '加盟商品', franchiseStock: '加盟库存',
  franchiseStocktake: '加盟盘点', franchisePromotion: '加盟促销',
  franchiseCoupon: '加盟优惠券',
};

const actionLabels: Record<string, string> = {
  read: '查看', create: '创建', update: '编辑', delete: '删除',
  restore: '恢复', suspend: '停用', provision: '开通',
  read_all: '查看全部', approve: '批准', reject: '退回',
  read_sensitive: '查看敏感信息', export: '导出', manage: '管理',
  cancel: '取消', grant: '发放', reverse: '冲正', correct: '更正',
  revoke: '撤销', record: '登记', post: '过账', submit: '提交',
};

export function permissionGroupLabel(resource: string): string {
  return entityTitles[resource] ?? domainTitles[resource] ?? resource;
}

export function permissionDisplayName(permission: PermissionOption): string {
  if (permission.action === 'franchise:provision') return '开通加盟商';
  const [resource, action] = permission.action.split(':');
  const title = permissionGroupLabel(resource);
  const verb = actionLabels[action];
  if (title !== resource && verb) return `${verb}${title}`;
  return permission.name && !permission.name.startsWith('permission.') ? permission.name : permission.action;
}
