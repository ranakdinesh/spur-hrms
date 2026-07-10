package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
)

func TestHandlerTenantAndActorExtraction(t *testing.T) {
	tenantID := uuid.New()
	userID := uuid.New()
	handler := New(nil,
		func(context.Context) string { return tenantID.String() },
		func(context.Context) uuid.UUID { return userID },
		func(context.Context) bool { return true },
		nil,
		nil,
	)
	request := httptest.NewRequest(http.MethodGet, "/hrms/tenant-profile", nil)
	gotTenantID, err := handler.tenantIDFromRequest(request)
	if err != nil {
		t.Fatalf("tenantIDFromRequest returned error: %v", err)
	}
	if gotTenantID != tenantID {
		t.Fatalf("got tenant %s, want %s", gotTenantID, tenantID)
	}
	gotActorID := handler.actorIDFromRequest(request)
	if gotActorID == nil || *gotActorID != userID {
		t.Fatalf("got actor %v, want %s", gotActorID, userID)
	}
	if !handler.isSuperAdminRequest(request) {
		t.Fatal("expected super admin request")
	}
}

func TestHandlerTenantExtractionErrors(t *testing.T) {
	tests := []struct {
		name    string
		handler *Handler
		request *http.Request
	}{
		{"nil request", New(nil, func(context.Context) string { return uuid.NewString() }, nil, nil, nil, nil), nil},
		{"missing resolver", New(nil, nil, nil, nil, nil, nil), httptest.NewRequest(http.MethodGet, "/", nil)},
		{"invalid uuid", New(nil, func(context.Context) string { return "bad" }, nil, nil, nil, nil), httptest.NewRequest(http.MethodGet, "/", nil)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.handler.tenantIDFromRequest(tt.request); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}

func TestRespondErrorWritesJSON(t *testing.T) {
	recorder := httptest.NewRecorder()
	respondError(recorder, http.StatusBadRequest, "bad request")
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("got status %d", recorder.Code)
	}
	var body map[string]string
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("response is not json: %v", err)
	}
	if body["error"] != "bad request" {
		t.Fatalf("got body %#v", body)
	}
}

func TestHandlerRespondErrorLogsAndWrites(t *testing.T) {
	handler := New(nil, func(context.Context) string { return uuid.NewString() }, nil, nil, nil, nil)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/hrms/test", nil)
	handler.respondError(recorder, request, http.StatusInternalServerError, "test operation", errors.New("boom"), "failed")
	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("got status %d", recorder.Code)
	}
	var body map[string]string
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("response is not json: %v", err)
	}
	if body["error"] != "failed" {
		t.Fatalf("got body %#v", body)
	}
}
