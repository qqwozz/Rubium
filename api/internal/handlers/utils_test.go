package handlers

import (
	"api/internal/validation"
	"testing"
)

func TestIsValidUUID(t *testing.T) {
	valid := "550e8400-e29b-41d4-a716-446655440000"
	invalid := "not-a-uuid"

	if !validation.IsValidUUID(valid) {
		t.Errorf("expected %s to be valid", valid)
	}
	if validation.IsValidUUID(invalid) {
		t.Errorf("expected %s to be invalid", invalid)
	}
	if validation.IsValidUUID("") {
		t.Error("expected empty string to be invalid")
	}
}