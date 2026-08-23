package handlers

import "testing"

func TestIsValidUUID(t *testing.T) {
	valid := "550e8400-e29b-41d4-a716-446655440000"
	invalid := "not-a-uuid"

	if !isValidUUID(valid) {
		t.Errorf("expected %s to be valid", valid)
	}
	if isValidUUID(invalid) {
		t.Errorf("expected %s to be invalid", invalid)
	}
	if isValidUUID("") {
		t.Error("expected empty string to be invalid")
	}
}