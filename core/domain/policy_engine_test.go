package domain

import (
	"errors"
	"testing"

	"github.com/google/uuid"
)

func TestValidatePolicyScope(t *testing.T) {
	scopeID := uuid.New()
	roleCode := "Manager"

	tests := []struct {
		name      string
		scopeType string
		scopeID   *uuid.UUID
		roleCode  *string
		want      string
		wantErr   error
	}{
		{name: "tenant default", scopeType: "tenant", want: PolicyScopeTenant},
		{name: "branch scope requires id", scopeType: "branch", scopeID: &scopeID, want: PolicyScopeBranch},
		{name: "department scope requires id", scopeType: "department", scopeID: &scopeID, want: PolicyScopeDepartment},
		{name: "designation scope requires id", scopeType: "designation", scopeID: &scopeID, want: PolicyScopeDesignation},
		{name: "workforce scope requires id", scopeType: "workforce_type", scopeID: &scopeID, want: PolicyScopeWorkforce},
		{name: "employee scope requires id", scopeType: "employee", scopeID: &scopeID, want: PolicyScopeEmployee},
		{name: "role group requires role code", scopeType: "role_group", roleCode: &roleCode, want: PolicyScopeRoleGroup},
		{name: "branch rejects missing id", scopeType: "branch", wantErr: ErrInvalidPolicyScope},
		{name: "role group rejects missing code", scopeType: "role_group", wantErr: ErrInvalidPolicyScope},
		{name: "tenant rejects id", scopeType: "tenant", scopeID: &scopeID, wantErr: ErrInvalidPolicyScope},
		{name: "unknown scope rejected", scopeType: "site", wantErr: ErrInvalidPolicyScope},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidatePolicyScope(tt.scopeType, tt.scopeID, tt.roleCode)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("ValidatePolicyScope() error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("ValidatePolicyScope() unexpected error = %v", err)
			}
			if got != tt.want {
				t.Fatalf("ValidatePolicyScope() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNormalizeRoleCodes(t *testing.T) {
	got := NormalizeRoleCodes([]string{" HR ", "manager", "hr", "", "MANAGER", "Employee"})
	want := []string{"hr", "manager", "employee"}
	if len(got) != len(want) {
		t.Fatalf("NormalizeRoleCodes() length = %d, want %d: %#v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("NormalizeRoleCodes()[%d] = %q, want %q; full result %#v", i, got[i], want[i], got)
		}
	}
}
