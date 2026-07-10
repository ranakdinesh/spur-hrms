package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/pkg/permissions"
	"github.com/ranakdinesh/spur-template/pkg/authcontext"
)

func TestTenantWideHRPayrollReadsDenyEmployeeSelfServicePermissions(t *testing.T) {
	tenantID := uuid.New()
	handler := New(
		nil,
		func(context.Context) string { return tenantID.String() },
		nil,
		func(context.Context) bool { return false },
		nil,
		authcontext.Permissions,
	)
	employeeGrants := []string{
		permissions.ModuleCode + "." + permissions.DashboardEmployeeView,
		permissions.ModuleCode + "." + permissions.SalarySlipsView,
		permissions.ModuleCode + "." + permissions.SalarySlipsDownload,
		permissions.ModuleCode + "." + permissions.AttendanceView,
		permissions.ModuleCode + "." + permissions.LeavesView,
	}

	tests := []struct {
		name    string
		target  string
		handler func(http.ResponseWriter, *http.Request)
	}{
		{name: "hr dashboard", target: "/hrms/dashboard/hr?month=7&year=2026", handler: handler.GetHRDashboard},
		{name: "payroll reconciliation", target: "/hrms/payroll-reconciliation?month=7&year=2026", handler: handler.ListPayrollReconciliation},
		{name: "salary templates", target: "/hrms/salary-templates", handler: handler.ListSalaryTemplates},
		{name: "payroll locks", target: "/hrms/payroll-locks", handler: handler.ListPayrollPeriodLocks},
		{name: "payroll imports", target: "/hrms/payroll-imports?limit=20", handler: handler.ListPayrollImportBatches},
		{name: "pay groups", target: "/hrms/pay-groups", handler: handler.ListPayGroups},
		{name: "pay runs", target: "/hrms/pay-runs?month=7&year=2026", handler: handler.ListPayRuns},
		{name: "consolidated salary sheet", target: "/hrms/consolidated-salary-sheet?month=7&year=2026", handler: handler.ListConsolidatedSalarySheet},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.target, nil)
			request = request.WithContext(authcontext.SetPermissions(request.Context(), employeeGrants))
			recorder := httptest.NewRecorder()

			tt.handler(recorder, request)

			if recorder.Code != http.StatusForbidden {
				t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusForbidden, recorder.Body.String())
			}
		})
	}
}
