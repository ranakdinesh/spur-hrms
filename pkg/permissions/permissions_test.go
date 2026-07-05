package permissions

import (
	"strings"
	"testing"
)

func TestCatalogKeysAreUniqueAndModuleLocal(t *testing.T) {
	seen := make(map[string]struct{}, len(Catalog))
	for _, permission := range Catalog {
		if permission.Key == "" {
			t.Fatal("permission key must not be empty")
		}
		if strings.HasPrefix(permission.Key, ModuleCode+".") {
			t.Fatalf("permission key %q must be module-local, not module-prefixed", permission.Key)
		}
		if permission.Description == "" {
			t.Fatalf("permission %q must have a description", permission.Key)
		}
		if _, ok := seen[permission.Key]; ok {
			t.Fatalf("duplicate permission key %q", permission.Key)
		}
		seen[permission.Key] = struct{}{}
	}
}

func TestRoleTemplatesReferenceCatalogPermissions(t *testing.T) {
	catalog := make(map[string]struct{}, len(Catalog))
	for _, permission := range Catalog {
		catalog[permission.Key] = struct{}{}
	}

	for _, template := range RoleTemplates {
		if template.Code == "" || template.Name == "" {
			t.Fatalf("role template code and name are required: %#v", template)
		}
		if len(template.Permissions) == 0 {
			t.Fatalf("role template %q must include at least one permission", template.Code)
		}
		seen := map[string]struct{}{}
		for _, permission := range template.Permissions {
			if _, ok := catalog[permission]; !ok {
				t.Fatalf("role template %q references unknown permission %q", template.Code, permission)
			}
			if _, ok := seen[permission]; ok {
				t.Fatalf("role template %q contains duplicate permission %q", template.Code, permission)
			}
			seen[permission] = struct{}{}
		}
	}
}

func TestManifestMatchesCatalogAndRoleTemplates(t *testing.T) {
	manifest := Manifest()
	if manifest.Code != ModuleCode {
		t.Fatalf("manifest code = %q, want %q", manifest.Code, ModuleCode)
	}
	if manifest.Name != ModuleName {
		t.Fatalf("manifest name = %q, want %q", manifest.Name, ModuleName)
	}
	if len(manifest.Permissions) != len(Catalog) {
		t.Fatalf("manifest has %d permissions, catalog has %d", len(manifest.Permissions), len(Catalog))
	}
	for i, permission := range Catalog {
		if manifest.Permissions[i].Slug != permission.Key {
			t.Fatalf("manifest permission %d slug = %q, want %q", i, manifest.Permissions[i].Slug, permission.Key)
		}
	}
	if len(manifest.RoleTemplates) != len(RoleTemplates) {
		t.Fatalf("manifest has %d role templates, catalog has %d", len(manifest.RoleTemplates), len(RoleTemplates))
	}
	for i, template := range RoleTemplates {
		if manifest.RoleTemplates[i].Code != template.Code {
			t.Fatalf("manifest template %d code = %q, want %q", i, manifest.RoleTemplates[i].Code, template.Code)
		}
	}
}

func TestEmployeeRoleTemplateIsTenantUserBaseline(t *testing.T) {
	employee := roleTemplateByCode(t, "EMPLOYEE")
	required := []string{
		DashboardEmployeeView,
		EmployeesView,
		LeavesApply,
		LeavesCancel,
		AttendanceCheckIn,
		AttendanceCheckOut,
		AttendanceRegularize,
		SalarySlipsView,
		SalarySlipsDownload,
		NotificationsRead,
		PoliciesView,
	}
	for _, permission := range required {
		if !containsPermission(employee.Permissions, permission) {
			t.Fatalf("EMPLOYEE template missing baseline permission %q", permission)
		}
	}

	for _, permission := range []string{EmployeesList, LeavesApprove, LeavesReject, AttendanceReport, EmployeesUpdate, DashboardHRView} {
		if containsPermission(employee.Permissions, permission) {
			t.Fatalf("EMPLOYEE template must not include elevated permission %q", permission)
		}
	}
}

func TestApplicantRoleTemplateCannotListEmployees(t *testing.T) {
	applicant := roleTemplateByCode(t, "APPLICANT")
	for _, permission := range []string{EmployeesList, EmployeesView, DashboardHRView} {
		if containsPermission(applicant.Permissions, permission) {
			t.Fatalf("APPLICANT template must not include employee admin permission %q", permission)
		}
	}
}

func TestManagerRoleTemplateIsTeamScopeAddon(t *testing.T) {
	manager := roleTemplateByCode(t, "MANAGER")
	required := []string{
		EmployeesList,
		EmployeesView,
		LeavesApprove,
		LeavesReject,
		LeavesReport,
		AttendanceReviewRequest,
		AttendanceReport,
	}
	for _, permission := range required {
		if !containsPermission(manager.Permissions, permission) {
			t.Fatalf("MANAGER template missing team permission %q", permission)
		}
	}

	for _, permission := range []string{DashboardEmployeeView, DashboardHRView, SalarySlipsView, NotificationsRead, PoliciesView} {
		if containsPermission(manager.Permissions, permission) {
			t.Fatalf("MANAGER template must stay additive and not include baseline permission %q", permission)
		}
	}
}

func roleTemplateByCode(t *testing.T, code string) RoleTemplate {
	t.Helper()
	for _, template := range RoleTemplates {
		if template.Code == code {
			return template
		}
	}
	t.Fatalf("role template %q not found", code)
	return RoleTemplate{}
}

func containsPermission(values []string, permission string) bool {
	for _, value := range values {
		if value == permission {
			return true
		}
	}
	return false
}
