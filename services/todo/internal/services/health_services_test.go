package services

import (
	"net/http"
	"testing"
	"time"
)

func TestHealthService_GetHealth(t *testing.T) {
	before := time.Now().UTC()

	resp, err := NewHealthService().GetHealth()

	after := time.Now().UTC()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Status != http.StatusOK {
		t.Errorf("status = %d, want %d", resp.Status, http.StatusOK)
	}
	if resp.Message != "todo service up and running" {
		t.Errorf("message = %q, want %q", resp.Message, "todo service up and running")
	}
	if resp.Time.Before(before) || resp.Time.After(after) {
		t.Errorf("time %v not within window [%v, %v]", resp.Time, before, after)
	}
}