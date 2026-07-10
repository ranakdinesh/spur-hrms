package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
	"github.com/ranakdinesh/spur-hrms/pkg/permissions"
	"github.com/ranakdinesh/spur-template/pkg/authcontext"
)

type employeeListTenantService struct {
	ports.TenantService
	called      bool
	gotTenantID uuid.UUID
}

func (s *employeeListTenantService) ListEmployees(ctx context.Context, tenantID uuid.UUID) ([]*domain.EmployeeListItem, error) {
	s.called = true
	s.gotTenantID = tenantID
	return []*domain.EmployeeListItem{
		{
			Employee: domain.Employee{
				ID:        uuid.New(),
				TenantID:  tenantID,
				UserID:    uuid.New(),
				Firstname: "Priya",
			},
		},
	}, nil
}

func TestListEmployeesRequiresEmployeesListPermission(t *testing.T) {
	tenantID := uuid.New()
	tests := []struct {
		name        string
		grants      []string
		wantStatus  int
		wantService bool
	}{
		{
			name:        "employee self service permissions are denied",
			grants:      []string{permissions.ModuleCode + "." + permissions.DashboardEmployeeView, permissions.ModuleCode + "." + permissions.EmployeesView},
			wantStatus:  http.StatusForbidden,
			wantService: false,
		},
		{
			name:        "applicant portal permissions are denied",
			grants:      []string{permissions.ModuleCode + "." + permissions.ApplicantPortalView},
			wantStatus:  http.StatusForbidden,
			wantService: false,
		},
		{
			name:        "tenant admin full permission key is allowed",
			grants:      []string{permissions.ModuleCode + "." + permissions.EmployeesList},
			wantStatus:  http.StatusOK,
			wantService: true,
		},
		{
			name:        "module local permission key is allowed",
			grants:      []string{permissions.EmployeesList},
			wantStatus:  http.StatusOK,
			wantService: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &employeeListTenantService{}
			handler := New(
				svc,
				func(context.Context) string { return tenantID.String() },
				nil,
				func(context.Context) bool { return false },
				nil,
				authcontext.Permissions,
			)
			request := httptest.NewRequest(http.MethodGet, "/hrms/employees", nil)
			request = request.WithContext(authcontext.SetPermissions(request.Context(), tt.grants))
			recorder := httptest.NewRecorder()

			handler.ListEmployees(recorder, request)

			if recorder.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d; body=%s", recorder.Code, tt.wantStatus, recorder.Body.String())
			}
			if svc.called != tt.wantService {
				t.Fatalf("service called = %t, want %t", svc.called, tt.wantService)
			}
			if tt.wantService {
				if svc.gotTenantID != tenantID {
					t.Fatalf("tenantID = %s, want %s", svc.gotTenantID, tenantID)
				}
				var body []domain.EmployeeListItem
				if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
					t.Fatalf("response is not employee list json: %v", err)
				}
				if len(body) != 1 || body[0].TenantID != tenantID {
					t.Fatalf("unexpected body: %#v", body)
				}
			}
		})
	}
}

func TestListEmployeesAllowsSuperAdminSupportContextWithoutTenantPermission(t *testing.T) {
	tenantID := uuid.New()
	svc := &employeeListTenantService{}
	handler := New(
		svc,
		func(context.Context) string { return tenantID.String() },
		nil,
		func(context.Context) bool { return true },
		nil,
		authcontext.Permissions,
	)
	request := httptest.NewRequest(http.MethodGet, "/hrms/employees", nil)
	recorder := httptest.NewRecorder()

	handler.ListEmployees(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusOK, recorder.Body.String())
	}
	if !svc.called {
		t.Fatal("expected service to be called")
	}
	if svc.gotTenantID != tenantID {
		t.Fatalf("tenantID = %s, want %s", svc.gotTenantID, tenantID)
	}
}
