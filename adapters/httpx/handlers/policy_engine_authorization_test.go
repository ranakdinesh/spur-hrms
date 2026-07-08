package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
	"github.com/ranakdinesh/spur-hrms/pkg/permissions"
	"github.com/ranakdinesh/spur-identity/adapters/http/httputil"
)

type policyEngineTenantService struct {
	ports.TenantService
	listCalled    bool
	createCalled  bool
	previewCalled bool
	gotTenantID   uuid.UUID
	gotKind       string
	gotSubject    domain.PolicyResolutionSubject
}

func (s *policyEngineTenantService) ListPolicySets(ctx context.Context, tenantID uuid.UUID, policyKind string) ([]*domain.PolicySet, error) {
	s.listCalled = true
	s.gotTenantID = tenantID
	s.gotKind = policyKind
	return []*domain.PolicySet{{ID: uuid.New(), TenantID: tenantID, PolicyKind: policyKind, Code: "default", Name: "Default", IsActive: true}}, nil
}

func (s *policyEngineTenantService) CreatePolicySet(ctx context.Context, cmd ports.PolicySetCommand) (*domain.PolicySet, error) {
	s.createCalled = true
	s.gotTenantID = cmd.TenantID
	s.gotKind = cmd.PolicyKind
	return &domain.PolicySet{ID: uuid.New(), TenantID: cmd.TenantID, PolicyKind: cmd.PolicyKind, Code: cmd.Code, Name: cmd.Name, IsActive: cmd.IsActive}, nil
}

func (s *policyEngineTenantService) ResolveEffectivePolicySet(ctx context.Context, subject domain.PolicyResolutionSubject, policyKind string) (*domain.PolicyResolutionResult, error) {
	s.previewCalled = true
	s.gotTenantID = subject.TenantID
	s.gotKind = policyKind
	s.gotSubject = subject
	return &domain.PolicyResolutionResult{
		Policy:     &domain.PolicySet{ID: uuid.New(), TenantID: subject.TenantID, PolicyKind: policyKind, Code: "effective", Name: "Effective", IsActive: true},
		Candidates: []domain.PolicyResolutionCandidate{},
	}, nil
}

func TestPolicyEngineListRequiresPolicyViewPermission(t *testing.T) {
	tenantID := uuid.New()
	tests := []struct {
		name        string
		grants      []string
		wantStatus  int
		wantService bool
	}{
		{
			name:        "employee self service permissions are denied",
			grants:      []string{permissions.ModuleCode + "." + permissions.AttendanceSelfView, permissions.ModuleCode + "." + permissions.LeaveSelfView},
			wantStatus:  http.StatusForbidden,
			wantService: false,
		},
		{
			name:        "attendance policy view is allowed",
			grants:      []string{permissions.ModuleCode + "." + permissions.AttendancePolicyView},
			wantStatus:  http.StatusOK,
			wantService: true,
		},
		{
			name:        "module local leave policy view is allowed",
			grants:      []string{permissions.LeavePolicyView},
			wantStatus:  http.StatusOK,
			wantService: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &policyEngineTenantService{}
			handler := New(
				svc,
				func(context.Context) string { return tenantID.String() },
				nil,
				func(context.Context) bool { return false },
				nil,
			)
			kind := domain.PolicyKindAttendance
			if tt.name == "module local leave policy view is allowed" {
				kind = domain.PolicyKindLeave
			}
			request := httptest.NewRequest(http.MethodGet, "/hrms/policy-engine/policy-sets?kind="+kind, nil)
			request = request.WithContext(httputil.SetPermissions(request.Context(), tt.grants))
			recorder := httptest.NewRecorder()

			handler.ListPolicyEnginePolicySets(recorder, request)

			if recorder.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d; body=%s", recorder.Code, tt.wantStatus, recorder.Body.String())
			}
			if svc.listCalled != tt.wantService {
				t.Fatalf("service called = %t, want %t", svc.listCalled, tt.wantService)
			}
			if tt.wantService && (svc.gotTenantID != tenantID || svc.gotKind != kind) {
				t.Fatalf("service tenant/kind = %s/%s, want %s/%s", svc.gotTenantID, svc.gotKind, tenantID, kind)
			}
		})
	}
}

func TestPolicyEngineCreateRequiresPolicyManagePermission(t *testing.T) {
	tenantID := uuid.New()
	body := []byte(`{"policy_kind":"attendance","code":"field","name":"Field Policy","is_active":true}`)
	tests := []struct {
		name        string
		grants      []string
		wantStatus  int
		wantService bool
	}{
		{
			name:        "view permission cannot create",
			grants:      []string{permissions.AttendancePolicyView},
			wantStatus:  http.StatusForbidden,
			wantService: false,
		},
		{
			name:        "manage permission can create",
			grants:      []string{permissions.AttendancePolicyManage},
			wantStatus:  http.StatusCreated,
			wantService: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &policyEngineTenantService{}
			handler := New(
				svc,
				func(context.Context) string { return tenantID.String() },
				nil,
				func(context.Context) bool { return false },
				nil,
			)
			request := httptest.NewRequest(http.MethodPost, "/hrms/policy-engine/policy-sets", bytes.NewReader(body))
			request = request.WithContext(httputil.SetPermissions(request.Context(), tt.grants))
			recorder := httptest.NewRecorder()

			handler.CreatePolicyEnginePolicySet(recorder, request)

			if recorder.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d; body=%s", recorder.Code, tt.wantStatus, recorder.Body.String())
			}
			if svc.createCalled != tt.wantService {
				t.Fatalf("service called = %t, want %t", svc.createCalled, tt.wantService)
			}
		})
	}
}

func TestPolicyEnginePreviewUsesTenantAndSubject(t *testing.T) {
	tenantID := uuid.New()
	employeeID := uuid.New()
	departmentID := uuid.New()
	requestBody := map[string]any{
		"policy_kind":      "leave",
		"employee_user_id": employeeID,
		"department_id":    departmentID,
		"role_codes":       []string{"Manager", "HR"},
		"date":             "2026-07-08",
	}
	body, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}
	svc := &policyEngineTenantService{}
	handler := New(
		svc,
		func(context.Context) string { return tenantID.String() },
		nil,
		func(context.Context) bool { return false },
		nil,
	)
	request := httptest.NewRequest(http.MethodPost, "/hrms/policy-engine/effective-preview", bytes.NewReader(body))
	request = request.WithContext(httputil.SetPermissions(request.Context(), []string{permissions.LeavePolicyView}))
	recorder := httptest.NewRecorder()

	handler.PreviewPolicyEngineEffectivePolicy(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusOK, recorder.Body.String())
	}
	if !svc.previewCalled {
		t.Fatal("expected preview service call")
	}
	if svc.gotTenantID != tenantID || svc.gotKind != domain.PolicyKindLeave {
		t.Fatalf("service tenant/kind = %s/%s, want %s/%s", svc.gotTenantID, svc.gotKind, tenantID, domain.PolicyKindLeave)
	}
	if svc.gotSubject.EmployeeUserID != employeeID {
		t.Fatalf("employee id = %s, want %s", svc.gotSubject.EmployeeUserID, employeeID)
	}
	if svc.gotSubject.DepartmentID == nil || *svc.gotSubject.DepartmentID != departmentID {
		t.Fatalf("department id = %v, want %s", svc.gotSubject.DepartmentID, departmentID)
	}
	if got := svc.gotSubject.Date.Format("2006-01-02"); got != "2026-07-08" {
		t.Fatalf("date = %s, want 2026-07-08", got)
	}
	if svc.gotSubject.Date.Location() != time.UTC {
		t.Fatalf("date location = %s, want UTC", svc.gotSubject.Date.Location())
	}
}
