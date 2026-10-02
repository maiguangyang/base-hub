/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package model

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

// TestFranchiseSchemaSourceContract 验证加盟组织 Schema 源文件的实体、API 与安全边界。
func TestFranchiseSchemaSourceContract(t *testing.T) {
	modelSchema := readSchemaFile(t, "model.graphql")
	extendSchema := readSchemaFile(t, "extend.graphql")

	assertContainsAll(t, modelSchema, []string{
		"type Account @entity", "type Organization @entity",
		"type OperatorMembership @entity", "type Permission @entity",
		"type OperatorRole @entity", "type Store @entity",
		"type Session @entity", "type MembershipInvitation @entity",
		"type AuditLog @entity",
	})
	assertContainsAll(t, modelSchema, []string{
		"enum OrganizationType", "enum PermissionScope", "enum WorkspaceType",
		"@validator(required: \"true\", type: \"phone\", unique: \"true\")",
		"uniqueScope: \"organizationId\"", "@relationship(inverse:",
	})
	assertContainsAll(t, extendSchema, []string{
		"viewer: Viewer!", "login(input: LoginInput!): LoginPayload!",
		"provisionFranchise(input: ProvisionFranchiseInput!)",
		"reviewStore(input: ReviewStoreInput!): Store!",
		"sessionEvents: SessionEvent!", "enum SessionEventCode",
	})
	assertContainsNone(t, modelSchema, []string{
		"type " + "User @entity", "type " + "Task @entity", "passwordHash", "passwordDigest",
		"tokenHash", "temporaryPassword", "type MembershipRole @entity",
		"type MembershipStore @entity", "type RolePermission @entity",
	})
	assertContainsNone(t, modelSchema+extendSchema, []string{"webSocket" + ": Any"})
}

// TestFranchiseInitialAccountSchema 确保初始账号独立于成员角色，并以组织 ID 重置密码。
func TestFranchiseInitialAccountSchema(t *testing.T) {
	modelSchema := readSchemaFile(t, "model.graphql")
	extendSchema := readSchemaFile(t, "extend.graphql")
	assertContainsAll(t, modelSchema, []string{
		`initializedOrganizations: [Organization!]! @relationship(inverse: "initialAccount")`,
		`initialAccount: Account @relationship(inverse: "initializedOrganizations")`,
	})
	assertContainsAll(t, extendSchema, []string{
		`resetFranchiseInitialPassword(organizationId: ID!): TemporaryPasswordPayload! @hasPermission(action: "account:update")`,
	})
	assertNoScalarEntityRelationshipIDs(t, modelSchema)
}

func TestFranchiseOpeningRecordSchema(t *testing.T) {
	modelSchema := readSchemaFile(t, "model.graphql")
	assertContainsAll(t, modelSchema, []string{
		`type FranchiseOpeningRecord @entity(title: "加盟商开通记录")`,
		`openingRecords: [FranchiseOpeningRecord!]! @relationship(inverse: "organization")`,
		`initialAccount: Account! @relationship(inverse: "openingRecords")`,
		`recordedByAccount: Account! @relationship(inverse: "recordedOpeningRecords")`,
		`organization: Organization! @relationship(inverse: "openingRecords")`,
		`recordNumber: String! @column(gorm: "type:varchar(128);NOT NULL;uniqueIndex;")`,
	})
	assertNoScalarEntityRelationshipIDs(t, modelSchema)
}

func TestCustomerSchemaUsesTypedRelationships(t *testing.T) {
	modelSchema := readSchemaFile(t, "model.graphql")
	assertContainsAll(t, modelSchema, []string{
		"type CustomerMember @entity",
		"type CustomerBenefitPolicy @entity",
		"type CustomerDailyPointGrantBudget @entity",
		"type CustomerPointEntry @entity",
		"type CustomerCouponTemplate @entity",
		"type CustomerCouponGrant @entity",
		`customerMembers: [CustomerMember!]! @relationship(inverse: "organization")`,
		`customerBenefitPolicies: [CustomerBenefitPolicy!]! @relationship(inverse: "organization")`,
		`customerDailyPointGrantBudgets: [CustomerDailyPointGrantBudget!]! @relationship(inverse: "organization")`,
		`customerCouponTemplates: [CustomerCouponTemplate!]! @relationship(inverse: "organization")`,
		`customerPointEntries: [CustomerPointEntry!]! @relationship(inverse: "sourceOrganization")`,
		`organization: Organization! @relationship(inverse: "customerMembers")`,
		`member: CustomerMember! @relationship(inverse: "pointEntries")`,
		`sourceOrganization: Organization! @relationship(inverse: "customerPointEntries")`,
		`reverses: CustomerPointEntry @relationship(inverse: "reversedBy")`,
		`reversedBy: CustomerPointEntry @relationship(inverse: "reverses")`,
		`template: CustomerCouponTemplate! @relationship(inverse: "grants")`,
		`member: CustomerMember! @relationship(inverse: "couponGrants")`,
		"phone: String @column",
	})
	if !regexp.MustCompile(`(?s)type CustomerMember @entity[^{]*\{[^}]*requestKey: String! @column`).MatchString(modelSchema) {
		t.Error("customer member creation request key missing")
	}
	assertNoScalarEntityRelationshipIDs(t, modelSchema)
}

func TestCustomerMemberCustomSchema(t *testing.T) {
	extendSchema := readSchemaFile(t, "extend.graphql")
	assertContainsAll(t, extendSchema, []string{
		"type HqCustomerMemberView",
		"type HqCustomerMemberPage",
		"phoneMasked: String!",
		`hqCustomerMembers(phone: String, status: CustomerMemberStatus, page: Int!, perPage: Int!): HqCustomerMemberPage! @hasPermission(action: "hqCustomer:read")`,
		`hqCustomerMember(id: ID!): HqCustomerMemberView @hasPermission(action: "hqCustomer:read")`,
		`hqCustomerSensitivePhone(id: ID!): String! @hasPermission(action: "customer:read_sensitive")`,
		`hqCreateCustomerMember(input: HqCreateCustomerMemberInput!): HqCustomerMemberView! @hasPermission(action: "hqCustomer:create")`,
		`hqSetCustomerMemberStatus(id: ID!, status: CustomerMemberStatus!): HqCustomerMemberView! @hasPermission(action: "hqCustomer:update")`,
		`hqRequestCustomerCancellation(id: ID!, identityEvidence: String!, basisCode: String!): HqCustomerMemberView! @hasPermission(action: "hqCustomer:cancel")`,
		`hqCompleteCustomerCancellation(id: ID!, dispositionReference: String!): HqCustomerMemberView! @hasPermission(action: "hqCustomer:cancel")`,
	})
}

func TestFranchiseInitialAccountGeneratedProjectionContract(t *testing.T) {
	generated := readSchemaFile(t, "../gen/generated.go")
	assertContainsAll(t, generated, []string{
		"InitialAccountID(ctx context.Context, obj *Organization) (*string, error)",
		"return ec.Resolvers.Organization().InitialAccountID(ctx, obj)",
	})
}

// TestEntityTitlesUseChineseLabels 验证 Dolphin 文档元数据直接使用中文实体标题。
func TestEntityTitlesUseChineseLabels(t *testing.T) {
	modelSchema := readSchemaFile(t, "model.graphql")

	assertContainsAll(t, modelSchema, []string{
		`type Account @entity(title: "账号")`,
		`type Organization @entity(title: "加盟商")`,
		`type OperatorMembership @entity(title: "加盟商成员")`,
		`type Permission @entity(title: "权限")`,
		`type OperatorRole @entity(title: "角色")`,
		`type Store @entity(title: "门店")`,
		`type Session @entity(title: "登录会话")`,
		`type MembershipInvitation @entity(title: "成员邀请")`,
		`type AuditLog @entity(title: "审计日志")`,
	})
	assertContainsNone(t, modelSchema, []string{`@entity(title: "entity.`})
}

// TestEntityReferencesUseRelationships 验证实体间引用必须由 typed relationship 表达。
func TestEntityReferencesUseRelationships(t *testing.T) {
	modelSchema := readSchemaFile(t, "model.graphql")
	assertContainsAll(t, modelSchema, []string{
		`reviewedByAccount: Account @relationship(inverse: "reviewedStores")`,
		`account: Account! @relationship(inverse: "sessions")`,
		`organization: Organization @relationship(inverse: "sessions")`,
		`membership: OperatorMembership! @relationship(inverse: "invitations")`,
		`invitedByAccount: Account! @relationship(inverse: "sentMembershipInvitations")`,
		`actorAccount: Account @relationship(inverse: "auditLogs")`,
		`session: Session @relationship(inverse: "auditLogs")`,
		`store: Store @relationship(inverse: "auditLogs")`,
		"resourceId: String",
	})
	assertNoScalarEntityRelationshipIDs(t, modelSchema)
}

// TestScalarEntityRelationshipIDsFailClosed 防止用 ID 等其他标量类型绕过关系约束。
func TestScalarEntityRelationshipIDsFailClosed(t *testing.T) {
	schema := `type Example @entity { ownerId: ID! @column }`
	if fields := scalarEntityRelationshipIDs(schema); len(fields) != 1 || fields[0] != "Example.ownerId" {
		t.Fatalf("scalar relationship ID was not rejected: %v", fields)
	}
	schema = `type StorePaymentConfig @entity {
merchantId: String
keyId: String
storeId: ID
}`
	if fields := scalarEntityRelationshipIDs(schema); len(fields) != 1 || fields[0] != "StorePaymentConfig.storeId" {
		t.Fatalf("payment identifiers bypassed relationship guard: %v", fields)
	}
}

// TestFranchiseGeneratedContract 验证 Dolphin 输出与已批准的 Schema 源保持一致。
func TestFranchiseGeneratedContract(t *testing.T) {
	generatedSchema := readSchemaFile(t, "../gen/schema.graphqls")
	generatedModels := readSchemaFile(t, "../gen/models.go")
	generatedHandlers := readSchemaFile(t, "../gen/resolver.go")
	databaseSource := readSchemaFile(t, "../gen/database.go")
	gqlgenConfig := readSchemaFile(t, "../gqlgen.yml")

	assertContainsAll(t, generatedSchema, []string{
		"type Account", "type Organization", "type OperatorMembership",
		"type Permission", "type OperatorRole", "type Store", "type Session",
		"type MembershipInvitation", "type AuditLog", "sessionEvents: SessionEvent!",
	})
	assertContainsAll(t, generatedHandlers, []string{
		"CreateAccount", "CreateOrganization", "CreateOperatorMembership",
		"CreatePermission", "CreateOperatorRole", "CreateStore", "CreateSession",
		"CreateMembershipInvitation", "CreateAuditLog",
	})
	assertContainsAll(t, databaseSource, []string{
		`"accounts":`, `"organizations":`, `"operator_memberships":`,
		`"permissions":`, `"operator_roles":`, `"stores":`, `"sessions":`,
		`"membership_invitations":`, `"audit_logs":`,
	})
	assertContainsAll(t, gqlgenConfig, []string{
		"  Account:", "  Organization:", "  OperatorMembership:",
		"  Permission:", "  OperatorRole:", "  Store:", "  Session:",
		"  MembershipInvitation:", "  AuditLog:",
	})
	assertContainsNone(t, generatedSchema+generatedHandlers+databaseSource+gqlgenConfig, []string{
		"type User", "type Task", "CreateUser", "CreateTask", `"users":`, `"tasks":`,
		"  User:", "  Task:",
	})
	assertContainsNone(t, generatedModels, []string{"InvitationID    *string"})
}

// readSchemaFile 读取同目录 Schema，并在读取失败时立即终止测试。
func readSchemaFile(t *testing.T, name string) string {
	t.Helper()
	data, err := os.ReadFile(name)
	if err != nil {
		t.Fatalf("read %s: %v", name, err)
	}
	return string(data)
}

// assertContainsAll 验证文本包含全部必需契约片段。
func assertContainsAll(t *testing.T, text string, snippets []string) {
	t.Helper()
	for _, snippet := range snippets {
		if !strings.Contains(text, snippet) {
			t.Errorf("missing schema contract %q", snippet)
		}
	}
}

// assertContainsNone 验证文本不包含废弃或敏感契约片段。
func assertContainsNone(t *testing.T, text string, snippets []string) {
	t.Helper()
	for _, snippet := range snippets {
		if strings.Contains(text, snippet) {
			t.Errorf("forbidden schema contract %q", snippet)
		}
	}
}

func assertNoScalarEntityRelationshipIDs(t *testing.T, schema string) {
	t.Helper()
	if fields := scalarEntityRelationshipIDs(schema); len(fields) > 0 {
		t.Fatalf("entity relationships must use @relationship; scalar ID fields found: %v", fields)
	}
}

func scalarEntityRelationshipIDs(schema string) []string {
	entityPattern := regexp.MustCompile(`(?s)type\s+(\w+)\s+@entity(?:\([^)]*\))?\s*\{(.*?)\}`)
	fieldPattern := regexp.MustCompile(`(?m)^\s*(\w+(?:Id|Ids))\s*:\s*[\[\]!\w]+`)
	var fields []string
	for _, entity := range entityPattern.FindAllStringSubmatch(schema, -1) {
		for _, field := range fieldPattern.FindAllStringSubmatch(entity[2], -1) {
			if entity[1] == "AuditLog" && field[1] == "resourceId" || isPaymentExternalID(entity[1], field[1]) {
				continue
			}
			fields = append(fields, entity[1]+"."+field[1])
		}
	}
	return fields
}

func isPaymentExternalID(entity, field string) bool {
	if field != "merchantId" && field != "keyId" {
		return false
	}
	switch entity {
	case "GlobalPaymentConfig", "FranchisePaymentConfig", "StorePaymentConfig":
		return true
	default:
		return false
	}
}
