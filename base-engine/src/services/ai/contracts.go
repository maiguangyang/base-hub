package ai

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"

	"github.com/gorilla/mux"
	"github.com/vektah/gqlparser/v2/ast"
)

type ContractRecord struct {
	OperationID       string   `json:"operationId"`
	Protocol          string   `json:"protocol"`
	Method            string   `json:"method,omitempty"`
	Path              string   `json:"path,omitempty"`
	Fingerprint       string   `json:"fingerprint"`
	Classifications   []string `json:"classifications"`
	Availability      string   `json:"availability"`
	NonCallableReason string   `json:"nonCallableReason,omitempty"`
}

func ContractHash(parts ...string) string {
	hash := sha256.New()
	for _, part := range parts {
		_, _ = fmt.Fprintf(hash, "%d:%s;", len(part), part)
	}
	return hex.EncodeToString(hash.Sum(nil))
}

func DiscoverContracts(schema *ast.Schema, router *mux.Router) ([]ContractRecord, error) {
	if schema == nil || router == nil {
		return nil, errors.New("missing schema or router")
	}
	records := graphQLContracts(schema)
	httpRecords, err := httpContracts(schema, router)
	if err != nil {
		return nil, err
	}
	records = append(records, httpRecords...)
	sort.Slice(records, func(i, j int) bool { return records[i].OperationID < records[j].OperationID })
	return records, nil
}

func CompareContracts(schema *ast.Schema, router *mux.Router, inventory []ContractRecord) error {
	actual, err := DiscoverContracts(schema, router)
	if err != nil {
		return err
	}
	baseline := make(map[string]ContractRecord, len(inventory))
	for _, record := range inventory {
		if record.OperationID == "" || baseline[record.OperationID].OperationID != "" {
			return fmt.Errorf("invalid or duplicate operationId %q", record.OperationID)
		}
		baseline[record.OperationID] = record
	}
	for _, record := range actual {
		previous, exists := baseline[record.OperationID]
		if !exists || !sameContractShape(previous, record) {
			return fmt.Errorf("contract drift: %s", record.OperationID)
		}
		delete(baseline, record.OperationID)
	}
	for id := range baseline {
		return fmt.Errorf("contract removed: %s", id)
	}
	return nil
}

func sameContractShape(left, right ContractRecord) bool {
	return left.Fingerprint == right.Fingerprint && left.Protocol == right.Protocol &&
		left.Method == right.Method && left.Path == right.Path
}
