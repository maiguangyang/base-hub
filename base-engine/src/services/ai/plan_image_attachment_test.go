package ai

import (
	"encoding/json"
	"testing"
)

func TestProductImagePlanRequiresUserSelectedAttachment(t *testing.T) {
	step := ApprovedStep{ToolID: productMainImageToolID, Arguments: json.RawMessage(`{"productId":"product-a","attachmentId":"selected"}`)}
	if err := validateImageAttachmentStep(step, "selected"); err != nil {
		t.Fatal(err)
	}
	if err := validateImageAttachmentStep(step, ""); err == nil {
		t.Fatal("missing user attachment was accepted")
	}
	if err := validateImageAttachmentStep(step, "other"); err == nil {
		t.Fatal("foreign attachment was accepted")
	}
	if err := validateImageAttachmentStep(ApprovedStep{ToolID: "HqUpdateProduct", Arguments: step.Arguments}, ""); err != nil {
		t.Fatal(err)
	}
}
