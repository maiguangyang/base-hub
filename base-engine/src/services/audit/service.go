/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package audit

import (
	"encoding/json"

	"github.com/gofrs/uuid"
	"base-engine/auth"
	"base-engine/gen"
	"gorm.io/gorm"
)

// Metadata 是允许写入审计 JSON 的非敏感字段集合。
type Metadata struct {
	ActorMembershipID       string   `json:"actorMembershipId,omitempty"`
	RequestID               string   `json:"requestId,omitempty"`
	ReasonCode              string   `json:"reasonCode,omitempty"`
	TargetStatus            string   `json:"targetStatus,omitempty"`
	RoleIDs                 []string `json:"roleIds,omitempty"`
	PermissionIDs           []string `json:"permissionIds,omitempty"`
	StoreIDs                []string `json:"storeIds,omitempty"`
	Source                  string   `json:"source,omitempty"`
	EvidenceReference       string   `json:"evidenceReference,omitempty"`
	ModelConfigVersion      uint64   `json:"modelConfigVersion,omitempty"`
	PaymentConfigOldVersion *uint64  `json:"paymentConfigOldVersion,omitempty"`
	PaymentConfigNewVersion *uint64  `json:"paymentConfigNewVersion,omitempty"`
	PaymentChannel          string   `json:"paymentChannel,omitempty"`
	PaymentScope            string   `json:"paymentScope,omitempty"`
	RunID                   string   `json:"runId,omitempty"`
	ToolID                  string   `json:"toolId,omitempty"`
	BeforeVersion           uint64   `json:"beforeVersion,omitempty"`
	AfterVersion            uint64   `json:"afterVersion,omitempty"`
	ExpectedPoints          int64    `json:"expectedPoints,omitempty"`
	LedgerPoints            int64    `json:"ledgerPoints,omitempty"`
	CountedQuantity         *int64   `json:"countedQuantity,omitempty"`
	CountedAt               string   `json:"countedAt,omitempty"`
}

// MetadataForPrincipal 使用请求操作者上下文补充白名单元数据。
func MetadataForPrincipal(principal *auth.WorkspacePrincipal, metadata Metadata) Metadata {
	if principal == nil {
		return metadata
	}
	metadata.RequestID = principal.RequestID
	if principal.MembershipID != nil {
		metadata.ActorMembershipID = *principal.MembershipID
	}
	return metadata
}

// Entry 描述一条结构化审计事件。
type Entry struct {
	ActorAccountID string
	SessionID      *string
	OrganizationID *string
	StoreID        *string
	Action         string
	ResourceType   string
	ResourceID     string
	ResultCode     string
	Metadata       Metadata
}

// Service 写入结构化不可变审计记录。
type Service struct{}

// NewService 创建审计服务。
func NewService() *Service { return &Service{} }

// Write 把 allowlist 元数据序列化到当前事务。
func (s *Service) Write(tx *gorm.DB, entry Entry) error {
	metadata, err := json.Marshal(withAIContext(tx.Statement.Context, entry.Metadata))
	if err != nil {
		return err
	}
	metadataJSON := string(metadata)
	var actorAccountID, resourceID *string
	if entry.ActorAccountID != "" {
		actorAccountID = &entry.ActorAccountID
	}
	if entry.ResourceID != "" {
		resourceID = &entry.ResourceID
	}
	record := gen.AuditLog{
		ID: uuid.Must(uuid.NewV4()).String(), ActorAccountID: actorAccountID,
		SessionID: entry.SessionID, OrganizationID: entry.OrganizationID, StoreID: entry.StoreID,
		Action: entry.Action, ResourceType: entry.ResourceType, ResourceID: resourceID,
		ResultCode: entry.ResultCode, MetadataJSON: &metadataJSON,
	}
	return tx.Create(&record).Error
}
