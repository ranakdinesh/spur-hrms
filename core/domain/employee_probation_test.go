package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestEmployeeProbationValidationAndStatus(t *testing.T) {
	joining := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC)
	mobile := "9999999999"

	employee, err := NewEmployee(EmployeeInput{
		TenantID:              uuid.New(),
		UserID:                uuid.New(),
		Firstname:             "Single",
		Mobile:                &mobile,
		JoiningDate:           &joining,
		ProbationStatus:       EmployeeProbationProbation,
		ProbationStartDate:    &joining,
		ProbationEndDate:      &end,
		ProbationDurationDays: PayrollStaffProbationDurationDays,
		IsPayrollStaff:        true,
	})
	if err != nil {
		t.Fatalf("NewEmployee() unexpected error = %v", err)
	}
	if !employee.IsOnProbation(time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)) {
		t.Fatal("employee should be on probation before probation end date")
	}
	if employee.IsOnProbation(time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)) {
		t.Fatal("employee should not be on probation after probation end date")
	}
}

func TestEmployeePayrollProbationRequiresSixMonths(t *testing.T) {
	mobile := "9999999999"
	_, err := NewEmployee(EmployeeInput{
		TenantID:              uuid.New(),
		UserID:                uuid.New(),
		Firstname:             "Payroll",
		Mobile:                &mobile,
		ProbationStatus:       EmployeeProbationProbation,
		ProbationDurationDays: 90,
		IsPayrollStaff:        true,
	})
	if err != ErrInvalidEmployeeProbation {
		t.Fatalf("NewEmployee() error = %v, want %v", err, ErrInvalidEmployeeProbation)
	}
}
