package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidFinancialYearID     = errors.New("financial_year_id is required")
	ErrInvalidFinancialYearName   = errors.New("financial year name is required")
	ErrInvalidFinancialYearPeriod = errors.New("financial year end date must be on or after start date")
	ErrFinancialYearLocked        = errors.New("financial year is locked")
)

type FinancialYear struct {
	ID            uuid.UUID  `json:"id"`
	TenantID      uuid.UUID  `json:"tenant_id"`
	Name          string     `json:"name"`
	StartDate     time.Time  `json:"start_date"`
	EndDate       time.Time  `json:"end_date"`
	IsActive      bool       `json:"is_active"`
	PayrollYear   bool       `json:"payroll_year"`
	LeaveYear     bool       `json:"leave_year"`
	HolidayYear   bool       `json:"holiday_year"`
	ReportingYear bool       `json:"reporting_year"`
	IsLocked      bool       `json:"is_locked"`
	LockedAt      *time.Time `json:"locked_at,omitempty"`
	LockedBy      *uuid.UUID `json:"locked_by,omitempty"`
	CloseNote     *string    `json:"close_note,omitempty"`
	Inactive      bool       `json:"inactive"`
	CreatedAt     time.Time  `json:"created_at"`
	CreatedBy     *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt     time.Time  `json:"updated_at"`
	UpdatedBy     *uuid.UUID `json:"updated_by,omitempty"`
}

type FinancialYearInput struct {
	TenantID      uuid.UUID
	Name          string
	StartDate     time.Time
	EndDate       time.Time
	PayrollYear   bool
	LeaveYear     bool
	HolidayYear   bool
	ReportingYear bool
	CloseNote     *string
}

func NewFinancialYear(input FinancialYearInput) (*FinancialYear, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	startDate := dateOnly(input.StartDate)
	endDate := dateOnly(input.EndDate)
	if startDate.IsZero() || endDate.IsZero() || endDate.Before(startDate) {
		return nil, ErrInvalidFinancialYearPeriod
	}
	name := strings.TrimSpace(input.Name)
	if name == "" {
		name = defaultFinancialYearName(startDate, endDate)
	}
	if name == "" {
		return nil, ErrInvalidFinancialYearName
	}
	now := time.Now().UTC()
	return &FinancialYear{
		TenantID:      input.TenantID,
		Name:          name,
		StartDate:     startDate,
		EndDate:       endDate,
		PayrollYear:   input.PayrollYear,
		LeaveYear:     input.LeaveYear,
		HolidayYear:   input.HolidayYear,
		ReportingYear: input.ReportingYear,
		CloseNote:     cleanOptional(input.CloseNote),
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}

func defaultFinancialYearName(startDate time.Time, endDate time.Time) string {
	if startDate.IsZero() || endDate.IsZero() {
		return ""
	}
	return fmt.Sprintf("FY %d-%d", startDate.Year(), endDate.Year())
}
