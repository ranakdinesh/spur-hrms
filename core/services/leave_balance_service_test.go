package services

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func TestLeaveRuleMatchesEmployeeScope(t *testing.T) {
	departmentID := uuid.New()
	otherDepartmentID := uuid.New()
	employmentTypeID := uuid.New()
	designationID := uuid.New()
	asOfDate := time.Date(2026, 7, 31, 0, 0, 0, 0, time.UTC)
	rule := &domain.LeavePolicyTemplateRule{
		DepartmentID:      &departmentID,
		EmploymentTypeID:  &employmentTypeID,
		DesignationID:     &designationID,
		ProbationStatus:   strPtr(domain.LeaveProbationConfirmed),
		CalculationConfig: map[string]any{"auto_apply_by_scope": true},
	}
	employee := &domain.Employee{
		UserID:           uuid.New(),
		DepartmentID:     &departmentID,
		EmploymentTypeID: &employmentTypeID,
		DesignationID:    &designationID,
		ProbationStatus:  domain.EmployeeProbationConfirmed,
	}
	if !leaveRuleMatchesEmployee(rule, employee, asOfDate) {
		t.Fatal("expected employee to match scoped confirmed rule")
	}
	employee.DepartmentID = &otherDepartmentID
	if leaveRuleMatchesEmployee(rule, employee, asOfDate) {
		t.Fatal("expected different department to be excluded")
	}
}

func TestLeaveRuleMatchesEmployeeProbationStatus(t *testing.T) {
	asOfDate := time.Date(2026, 7, 31, 0, 0, 0, 0, time.UTC)
	probationEnd := time.Date(2026, 8, 31, 0, 0, 0, 0, time.UTC)
	rule := &domain.LeavePolicyTemplateRule{ProbationStatus: strPtr(domain.LeaveProbationConfirmed)}
	employee := &domain.Employee{
		UserID:             uuid.New(),
		ProbationStatus:    domain.EmployeeProbationProbation,
		ProbationStartDate: &asOfDate,
		ProbationEndDate:   &probationEnd,
	}
	if leaveRuleMatchesEmployee(rule, employee, asOfDate) {
		t.Fatal("expected probation employee to be excluded from confirmed rule")
	}
}

func TestBoolFromLeaveRuleConfig(t *testing.T) {
	if !boolFromLeaveRuleConfig(map[string]any{"auto_apply_by_scope": "yes"}, "auto_apply_by_scope", false) {
		t.Fatal("expected yes string to resolve true")
	}
	if boolFromLeaveRuleConfig(map[string]any{}, "auto_apply_by_scope", false) {
		t.Fatal("expected missing value to use false fallback")
	}
}
