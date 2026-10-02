package ai

import (
	"encoding/json"
	"errors"
)

const productMainImageToolID = "HqSetProductMainImage"

func validateImageAttachmentStep(step ApprovedStep, expectedID string) error {
	if step.ToolID != productMainImageToolID {
		return nil
	}
	if expectedID == "" {
		return errors.New("IMAGE_ATTACHMENT_NOT_SELECTED")
	}
	var args struct {
		AttachmentID string `json:"attachmentId"`
	}
	if err := json.Unmarshal(step.Arguments, &args); err != nil || args.AttachmentID != expectedID {
		return errors.New("IMAGE_ATTACHMENT_MISMATCH")
	}
	return nil
}
