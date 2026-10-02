package tools

import (
	"fmt"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"
)

var inputFieldDescriptions = map[string]string{
	"accountId":                "账号的唯一 ID；必须来自已授权的账号查询结果",
	"action":                   "审计动作代码；按权限或审计记录中的实际动作筛选",
	"actorAccountId":           "执行审计动作的账号 ID",
	"address":                  "门店详细地址，需与省、市、区字段一致",
	"amountFen":                "优惠券面额，单位为整数分",
	"batchId":                  "门店库存批次 ID，必须来自已授权的批次查询",
	"batchIds":                 "本次盘点选择的已有批次 ID 列表，每项须属于目标门店",
	"batchNumber":              "供应商批号或生产批号；无外部批号时可省略",
	"barcode":                  "包装条码；有外部可扫描条码时填写",
	"attachmentId":             "本轮操作员上传并选择的图片附件 ID；只能使用对话提供的 Attachment-ID",
	"brand":                    "商品品牌；仅按总部已确认的资料填写",
	"brandId":                  "总部品牌 ID，须先通过品牌查询核实",
	"categoryId":               "总部商品分类 ID，须先通过分类查询核实",
	"listingEnabled":           "只查询该门店已选用或未选用的规格；空值表示全部",
	"published":                "只查询总部仍发布或已停用的规格；空值表示全部",
	"containsPackageId":        "此包装直接包含的下一级包装 ID",
	"containsQuantity":         "此包装包含下一级包装的正整数数量",
	"defaultPackageTemplateId": "新 SKU 默认复制的包装换算模板 ID；可留空",
	"disableDefaultPackage":    "新 SKU 不继承默认包装；已有 SKU 选择后停用当前包装及其门店报价，保留历史记录",
	"delta":                    "库存盘点差异，正数增加、负数扣减，单位为完整包装",
	"description":              "总部商品描述，按确认过的资料填写",
	"discountFen":              "门店活动减免金额，单位为整数分",
	"endsAt":                   "门店活动结束时间，结束时刻不再生效",
	"expiresAt":                "批次失效时间；要求有效期的商品入库必须填写",
	"factorSnapshot":           "库存拆包发生时保存的换算数量快照",
	"hasDifference":            "仅在复核或已入账盘点单中筛选是否存在实盘差异；盘点中单据不参与此筛选",
	"from":                     "时间范围起点，按用户明确给出的本地日期换算为完整时间",
	"fixedPriceFen":            "活动固定价格，单位为整数分",
	"imageUrl":                 "已登记的商品图片引用地址",
	"identifier":               "顾客主动提供的完整手机号或已核实的会员 ID，禁止前缀查询",
	"ingredients":              "商品配料信息，按总部确认的资料填写",
	"allergens":                "商品过敏原信息，按总部确认的资料填写",
	"listingId":                "门店已引用的销售规格 ID；批次库存查询可省略以读取当前授权门店全部批次，其它操作须先查询核实",
	"listingIds":               "本次盘点选择的门店商品 ID 列表，每项须属于目标门店",
	"lineId":                   "目标盘点行 ID，须先读取所属门店盘点单确认",
	"loss":                     "此库存调整是否属于报损，报损只能减少库存",
	"origin":                   "商品产地信息，按总部确认的资料填写",
	"offerId":                  "门店包装销售项 ID，须先通过授权门店查询核实",
	"orderedIds":               "该规格下全部规格值 ID，按用户要求的最终展示顺序排列；不得遗漏或重复",
	"packageId":                "商品包装 ID，须先确认属于目标销售规格",
	"packageTemplateId":        "当前 SKU 单独使用的包装换算模板 ID；已有 SKU 会生成并发布新包装版本",
	"packageTemplateIds":       "需要复制到商品销售规格的包装换算模板 ID 列表",
	"packageSetVersion":        "包装换算集合版本，历史库存沿用原版本",
	"parentId":                 "上级商品分类 ID；顶层分类可省略",
	"priceFen":                 "门店包装销售价格，单位为整数分",
	"producedAt":               "商品批次生产时间；有明确生产日期时填写",
	"productId":                "总部商品 ID，须先通过商品查询核实",
	"promotionId":              "待改版活动 ID；新建活动时可省略",
	"quantity":                 "完整包装整数数量；盘点实数可为零，入库与拆包必须大于零",
	"requiredQuantity":         "组合活动中该目标包装所需数量",
	"ruleKey":                  "活动各版本共享的规则标识",
	"shelfLifeDays":            "商品 SKU 的保质期天数，使用正整数",
	"sellable":                 "只查询当前可售或已过期的库存批次；空值表示全部",
	"skuId":                    "总部商品 SKU ID，须先通过商品 SKU 查询核实",
	"skuOverrides":             "填写需改变启停状态或包装选择的 SKU 组合；未填写的既有 SKU 保留原状态和包装，新组合默认启用",
	"sourcePackageId":          "本次拆包的源包装 ID",
	"sourceReference":          "商品入库来源单据或凭证引用",
	"specificationId":          "商品规格 ID，须先查询总部规格目录确认",
	"selections":               "商品选择的规格及各自已勾选的值；更新时先查询商品原有 selectedValueIds 并按规格归组，避免清除原选择；未选择时传空数组",
	"valueIds":                 "规格值 ID 列表；须属于对应规格，新组合只能使用启用的规格和值",
	"startsAt":                 "门店活动开始时间，包含开始时刻",
	"stackWithHqCoupon":        "门店是否请求活动与总部券叠加，仍受总部规则限制",
	"stackWithStoreCoupon":     "门店是否请求活动与门店券叠加，仍受总部规则限制",
	"stackWithMemberPrice":     "门店是否请求会员价与活动叠加，仍受总部规则限制",
	"storageInstructions":      "商品储存要求，按总部确认的资料填写",
	"suggestedPriceFen":        "总部建议销售价，单位为整数分",
	"targetPackageId":          "本次拆包直接得到的目标包装 ID",
	"targets":                  "门店活动适用的包装销售项及所需数量",
	"thresholdFen":             "满额活动的金额门槛，单位为整数分",
	"thresholdQuantity":        "满件活动的完整包装件数门槛",
	"timeZone":                 "门店活动使用的 IANA 时区名称",
	"to":                       "时间范围终点，按用户明确给出的本地日期换算到当日末尾",
	"basisCode":                "会员注销处理依据代码；只从本轮操作员明确标注为‘依据代码’的内容提取，缺少时先询问",
	"businessHours":            "门店营业时间，使用业务接口接受的时间格式",
	"businessStatus":           "门店营业状态，使用接口定义的状态值",
	"channel":                  "支付渠道，只能选择微信或支付宝",
	"city":                     "门店所在城市",
	"code":                     "组织或门店的业务编码；创建时按用户提供的编码填写",
	"contactPhone":             "门店对外联系电话",
	"createdAt":                "审计记录的创建时间筛选值，使用接口要求的整数时间值",
	"daysAfterActivation":      "优惠券激活后的有效天数，使用 1 至 365 的整数",
	"effectiveAt":              "优惠券开始生效的日期时间；省略表示创建后立即生效",
	"distributionEndsAt":       "优惠券自动派发截止日期时间；省略表示不限制，到达该时刻不再向尚未获券的会员派发",
	"discountBasisPoints":      "折扣基点数，按所属权益或活动规则允许的范围填写",
	"discountEnabled":          "是否启用整单会员折扣",
	"dispositionReference":     "注销前权益处置凭证编号；只从本轮操作员明确标注为‘凭证编号’或‘工单号’的内容提取，缺少时先询问",
	"displayName":              "被邀请管理员或员工的显示名称",
	"district":                 "门店所在区县",
	"email":                    "被邀请人员的电子邮箱；只有用户提供时才填写",
	"earnAmountFen":            "未来消费赠分门槛金额，单位为整数分",
	"earnPoints":               "达到消费赠分门槛时赠送的整数积分",
	"enabled":                  "是否启用目标商品、活动或券模板",
	"entryId":                  "待冲正积分流水的唯一 ID，须先查询核实",
	"evidenceReference":        "会员积分校正凭证编号；只从本轮操作员明确标注为‘凭证编号’或‘工单号’的内容提取，缺少时先询问",
	"expectedVersion":          "当前会员权益规则版本，用于并发修改校验",
	"filter":                   "列表筛选对象；只填写有明确条件的字段，未提供的条件保持省略",
	"id":                       "目标记录的唯一 ID；先通过有权限的只读查询取得",
	"identityEvidence":         "身份核验凭证编号；只从本轮操作员明确标注为‘凭证编号’或‘工单号’的内容提取，缺少时先询问",
	"ids":                      "要操作的记录 ID 列表；每个 ID 均需来自已核实的目标",
	"input":                    "执行当前业务操作的参数对象；完整填写必填字段，不添加未声明的字段",
	"kind":                     "角色或业务对象的类型枚举，选择与当前工作区匹配的值",
	"lifecycle":                "门店生命周期状态，使用接口定义的枚举值",
	"managerName":              "门店负责人的姓名",
	"managerPhone":             "门店负责人的联系电话",
	"manualGrantMaxDaily":      "总部自然日人工赠分总上限，使用整数积分",
	"manualGrantMaxSingle":     "单次人工赠分上限，使用整数积分",
	"maxRedemptionBasisPoints": "未来单次消费积分抵扣比例上限，使用基点数",
	"maxRedemptionPoints":      "未来单次消费可抵扣的最大积分",
	"memberId":                 "目标顾客会员的唯一 ID，须先查询核实",
	"minSpendFen":              "优惠券最低消费门槛，单位为整数分",
	"membershipId":             "员工或管理员成员关系的唯一 ID",
	"name":                     "组织、门店或角色的名称，按用户明确提供的名称填写",
	"note":                     "积分赠送或冲正时仅填写凭证编号，不填说明；编号须为 4 至 64 位字母、数字、_ 或 -，且含数字或分隔符；只从本轮操作员明确标注为‘凭证编号’或‘工单号’的内容提取，缺少时先询问",
	"organizationId":           "组织的唯一 ID；加盟商工作区必须保持在当前授权组织内",
	"ownerDisplayName":         "加盟商初始负责人显示名称",
	"ownerEmail":               "加盟商初始负责人的电子邮箱；可省略",
	"ownerPhone":               "加盟商初始负责人的手机号",
	"page":                     "从 1 开始的页码；需要更多结果时逐页递增",
	"pageSize":                 "每页记录数，范围为 1 至 50",
	"perMemberLimit":           "每位会员可累计获发该券的最大数量",
	"perPage":                  "每页记录数，范围为 1 至 50",
	"permissionsIds":           "授予角色的权限 ID 列表；先读取权限目录并逐项核实",
	"phone":                    "目标人员或顾客会员的手机号；建档时须为 11 位中国大陆手机号，查询时只使用操作员给出的号码或前缀",
	"points":                   "本次人工赠送的整数积分",
	"promotionWithHqCoupon":    "总部是否允许门店活动与总部券叠加",
	"promotionWithStoreCoupon": "总部是否允许门店活动与门店券叠加",
	"memberPriceWithPromotion": "总部是否允许门店会员价与其他门店活动叠加",
	"pointsWithPromotion":      "总部是否允许积分与门店活动叠加",
	"pointsWithCoupon":         "总部是否允许积分与优惠券叠加",
	"purchaseEarnEnabled":      "是否启用未来消费赠分规则",
	"province":                 "门店所在省份",
	"q":                        "列表关键词；使用用户给出的名称、编码或其他明确线索",
	"reasonCode":               "停用组织的业务原因代码，按用户确认的理由填写",
	"redeemAmountFen":          "未来积分抵扣对应金额，单位为整数分",
	"redeemPoints":             "未来抵扣换算所需的整数积分",
	"redemptionEnabled":        "是否启用未来积分抵扣规则",
	"receiptFooter":            "门店小票底部展示文字",
	"recordId":                 "当前范围与渠道的支付配置记录 ID；创建停用覆盖时填写空字符串",
	"rejectionReason":          "驳回门店申请的原因；需具体说明给门店方",
	"resourceId":               "审计对象的唯一 ID",
	"resourceType":             "审计对象的资源类型",
	"requestKey":               "幂等请求键，由受信代码生成，模型不得自行填写",
	"resultCode":               "审计结果代码，用于筛选成功或失败记录",
	"roleIds":                  "待分配给人员的角色 ID 列表；先核实角色权限",
	"rolesIds":                 "人员关联的角色 ID 列表；先核实角色权限",
	"status":                   "目标对象的状态，选择接口列出的枚举值",
	"state":                    "支付渠道目标状态，只能为启用或停用",
	"storeAccessMode":          "人员的门店访问范围模式，选择接口定义的枚举值",
	"scope":                    "支付配置范围，全局、加盟商或加盟门店",
	"storeArea":                "门店面积数值，按用户提供的面积填写",
	"storeId":                  "目标门店的唯一 ID，需属于当前授权范围",
	"storeIds":                 "人员可访问的门店 ID 列表，需逐项核实授权范围",
	"storesIds":                "人员关联的门店 ID 列表，需逐项核实授权范围",
	"supportDineIn":            "是否支持堂食，使用布尔值",
	"supportTakeout":           "是否支持外带，使用布尔值",
	"tableCount":               "门店餐桌数量，使用整数",
	"templateId":               "目标优惠券模板的唯一 ID，须先查询核实",
	"title":                    "优惠券模板对外展示标题",
	"totalIssueLimit":          "券模板允许的累计发放总量；省略表示不限制",
	"grantId":                  "待作废发券记录的唯一 ID，须先查询核实",
	"version":                  "当前支付配置的并发版本；未建记录时填写 0",
}

func describeInputSchema(name string, schema *jsonschema.Schema, context string) {
	if schema == nil {
		return
	}
	meaning := inputFieldDescriptions[name]
	if meaning == "" {
		meaning = fmt.Sprintf("%s 字段；仅按已确认的业务资料填写", name)
	}
	schema.Description = fmt.Sprintf("%s。用于%s。", strings.TrimSuffix(meaning, "。"), context)
	applyCouponInputConstraints(name, schema, context)
	appendInputEnumDescription(schema)
	if schema.Items != nil && schema.Items.Description == "" {
		schema.Items.Description = "列表中的单个值；类型须符合此字段的元素约束。"
	}
}

func applyCouponInputConstraints(name string, schema *jsonschema.Schema, context string) {
	if strings.Contains(context, "CouponTemplateInput") &&
		(name == "effectiveAt" || name == "distributionEndsAt" || name == "totalIssueLimit") {
		schema.Types = []string{schema.Type, "null"}
		schema.Type = ""
	}
	if name == "daysAfterActivation" {
		minimum, maximum := float64(1), float64(365)
		schema.Minimum, schema.Maximum = &minimum, &maximum
	}
}

func appendInputEnumDescription(schema *jsonschema.Schema) {
	if len(schema.Enum) > 0 {
		var values []string
		for _, value := range schema.Enum {
			values = append(values, fmt.Sprint(value))
		}
		schema.Description += " 可选值：" + strings.Join(values, "、") + "。"
	}
}
