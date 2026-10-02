package ai

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"
	"base-engine/auth"
)

type Phase string

const (
	PhasePreview Phase = "PREVIEW"
	PhaseRun     Phase = "RUN"
)

type ToolMode string

const (
	ModeReadOnly ToolMode = "READ_ONLY"
	ModeWrite    ToolMode = "WRITE"
)

type WriteKind string

const (
	WriteCreate   WriteKind = "CREATE"
	WriteExisting WriteKind = "EXISTING"
)

type ToolSpec struct {
	ID, Name, Title, OperationID, Document, Path, Method, Description string
	Mode                                                              ToolMode
	Risk, Permission                                                  string
	AdditionalPermissions                                             []string
	Workspaces                                                        []auth.WorkspaceType
	OutputFields                                                      []string
	WriteKind                                                         WriteKind
	TargetFields, ScopeFields, ArgumentFields                         []string
	RequiredAttestation                                               bool
	SingleCallApproval                                                bool
	EvidencePaths                                                     []string
	RequestKeyPath                                                    string
	GeneratedCodePath                                                 string
	InputSchema                                                       *jsonschema.Resolved
}

type Catalog struct {
	specs []ToolSpec
	byID  map[string]ToolSpec
}

type contractWorkspace struct {
	operationID string
	workspace   auth.WorkspaceType
}

// NewCatalog 根据接口清单注册工具，并拒绝缺少工具的可调用接口。
func NewCatalog(inventory []ContractRecord, specs []ToolSpec) (*Catalog, error) {
	contracts := make(map[string]ContractRecord, len(inventory))
	for _, record := range inventory {
		if err := validateInventoryRecord(record, contracts); err != nil {
			return nil, err
		}
		contracts[record.OperationID] = record
	}
	catalog := &Catalog{byID: make(map[string]ToolSpec, len(specs))}
	covered := make(map[contractWorkspace]struct{}, len(specs))
	for _, spec := range specs {
		if spec.Name == "" {
			spec.Name = spec.ID
		}
		if spec.Title == "" {
			spec.Title = spec.Description
		}
		if err := validateToolSpec(spec, contracts[spec.OperationID]); err != nil {
			return nil, fmt.Errorf("%s: %w", spec.ID, err)
		}
		if _, exists := catalog.byID[spec.ID]; exists {
			return nil, fmt.Errorf("duplicate tool id: %s", spec.ID)
		}
		catalog.specs = append(catalog.specs, spec)
		catalog.byID[spec.ID] = spec
		for _, workspace := range spec.Workspaces {
			covered[contractWorkspace{spec.OperationID, workspace}] = struct{}{}
		}
	}
	if err := validateCallableCoverage(inventory, covered); err != nil {
		return nil, err
	}
	return catalog, nil
}

func validateInventoryRecord(record ContractRecord, known map[string]ContractRecord) error {
	if record.OperationID == "" || known[record.OperationID].OperationID != "" {
		return fmt.Errorf("invalid or duplicate operation id: %s", record.OperationID)
	}
	if record.Availability == "NON_CALLABLE" && (record.NonCallableReason == "" || record.NonCallableReason == "NOT_REVIEWED_FOR_AI") {
		return fmt.Errorf("unreviewed non-callable operation: %s", record.OperationID)
	}
	return nil
}

// validateCallableCoverage 确保每个可调用接口在其每个后台工作区都有工具。
func validateCallableCoverage(inventory []ContractRecord, covered map[contractWorkspace]struct{}) error {
	for _, record := range inventory {
		switch record.Availability {
		case "NON_CALLABLE":
			continue
		case "CALLABLE":
			if err := requireWorkspaceTools(record, covered); err != nil {
				return err
			}
		default:
			return fmt.Errorf("invalid availability for %s: %s", record.OperationID, record.Availability)
		}
	}
	return nil
}

func requireWorkspaceTools(record ContractRecord, covered map[contractWorkspace]struct{}) error {
	if len(record.Classifications) == 0 {
		return fmt.Errorf("callable operation has no workspace classification: %s", record.OperationID)
	}
	for _, classification := range record.Classifications {
		workspace, ok := classifiedWorkspace(classification)
		if !ok {
			return fmt.Errorf("unsupported callable classification: %s %s", record.OperationID, classification)
		}
		if _, ok := covered[contractWorkspace{record.OperationID, workspace}]; !ok {
			return fmt.Errorf("callable operation has no registered tool in %s: %s", workspace, record.OperationID)
		}
	}
	return nil
}

func classifiedWorkspace(classification string) (auth.WorkspaceType, bool) {
	switch classification {
	case "HEADQUARTERS_ADMIN":
		return auth.WorkspaceTypeHeadquarters, true
	case "FRANCHISE_ADMIN":
		return auth.WorkspaceTypeFranchise, true
	default:
		return "", false
	}
}

func validateToolSpec(spec ToolSpec, contract ContractRecord) error {
	if !hasRequiredToolMetadata(spec, contract) {
		return errors.New("invalid tool contract")
	}
	if err := validateToolTransport(spec, contract); err != nil {
		return err
	}
	for _, workspace := range spec.Workspaces {
		if !contractIncludesWorkspace(contract, workspace) {
			return errors.New("workspace outside contract classification")
		}
	}
	return nil
}

func hasRequiredToolMetadata(spec ToolSpec, contract ContractRecord) bool {
	return spec.ID != "" && contract.OperationID != "" && contract.Availability == "CALLABLE" &&
		spec.Name == spec.ID && spec.Title != "" && spec.Description != "" && spec.Permission != "" && spec.Risk != "" &&
		len(spec.Workspaces) > 0 && len(spec.OutputFields) > 0
}

func validateToolTransport(spec ToolSpec, contract ContractRecord) error {
	if spec.Mode != ModeReadOnly && spec.Mode != ModeWrite {
		return errors.New("invalid tool mode")
	}
	if contract.Protocol == "HTTP_ROUTE" {
		return validateHTTPToolTransport(spec, contract)
	}
	if !strings.HasPrefix(contract.Protocol, "GRAPHQL_") || spec.Document == "" || spec.Path != "" || spec.Method != "" {
		return errors.New("invalid GraphQL template")
	}
	if (spec.Mode == ModeReadOnly) != (contract.Protocol == "GRAPHQL_QUERY") {
		return errors.New("tool mode does not match protocol")
	}
	return validateToolDocument(spec, contract)
}

func validateHTTPToolTransport(spec ToolSpec, contract ContractRecord) error {
	if spec.Path != contract.Path || spec.Method != contract.Method || spec.Document != "" {
		return errors.New("invalid HTTP template")
	}
	return nil
}

func contractIncludesWorkspace(record ContractRecord, workspace auth.WorkspaceType) bool {
	category := ""
	switch workspace {
	case auth.WorkspaceTypeHeadquarters:
		category = "HEADQUARTERS_ADMIN"
	case auth.WorkspaceTypeFranchise:
		category = "FRANCHISE_ADMIN"
	default:
		return false
	}
	for _, value := range record.Classifications {
		if value == category {
			return true
		}
	}
	return false
}

func (c *Catalog) Lookup(id string) (ToolSpec, bool) {
	if c == nil {
		return ToolSpec{}, false
	}
	spec, exists := c.byID[id]
	return spec, exists
}

func (c *Catalog) Visible(principal *auth.WorkspacePrincipal, phase Phase, approvedIDs map[string]struct{}) []ToolSpec {
	if c == nil || principal == nil {
		return nil
	}
	visible := make([]ToolSpec, 0)
	for _, spec := range c.specs {
		if !specAllowed(spec, principal, phase, approvedIDs) {
			continue
		}
		visible = append(visible, spec)
	}
	return visible
}

func specAllowed(spec ToolSpec, principal *auth.WorkspacePrincipal, phase Phase, approvedIDs map[string]struct{}) bool {
	if !principal.Has(spec.Permission) {
		return false
	}
	for _, action := range spec.AdditionalPermissions {
		if !principal.Has(action) {
			return false
		}
	}
	workspaceAllowed := false
	for _, workspace := range spec.Workspaces {
		workspaceAllowed = workspaceAllowed || workspace == principal.WorkspaceType
	}
	if !workspaceAllowed {
		return false
	}
	if phase == PhasePreview {
		return spec.Mode == ModeReadOnly
	}
	if phase != PhaseRun {
		return false
	}
	_, approved := approvedIDs[spec.ID]
	return approved
}
