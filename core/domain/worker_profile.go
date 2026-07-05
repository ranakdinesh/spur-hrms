package domain

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	WorkerProfileStatusDraft       = "draft"
	WorkerProfileStatusActive      = "active"
	WorkerProfileStatusPaused      = "paused"
	WorkerProfileStatusEnded       = "ended"
	WorkerProfileStatusBlacklisted = "blacklisted"

	WorkerCompliancePending    = "pending"
	WorkerComplianceReady      = "ready"
	WorkerComplianceReview     = "review_required"
	WorkerComplianceBlocked    = "blocked"
	WorkerPayrollNotApplicable = "not_applicable"
	WorkerPayrollPending       = "pending"
	WorkerPayrollReady         = "ready"
	WorkerPayrollBlocked       = "blocked"
)

var (
	ErrInvalidWorkerProfileID               = errors.New("worker_profile_id is required")
	ErrInvalidWorkerProfileName             = errors.New("worker profile display_name is required")
	ErrInvalidWorkerProfileStatus           = errors.New("worker profile status is invalid")
	ErrInvalidWorkerProfileComplianceStatus = errors.New("worker profile compliance_status is invalid")
	ErrInvalidWorkerProfilePayrollStatus    = errors.New("worker profile payroll_status is invalid")
	ErrInvalidWorkerProfileDates            = errors.New("worker profile end_date cannot be before start_date")
	ErrInvalidWorkerProfileEmployeeLink     = errors.New("worker profile employee_id and employee_user_id must be provided together")
	ErrInvalidWorkerProfileMetadata         = errors.New("worker profile metadata must be a valid JSON object")
)

type WorkerProfile struct {
	ID                 uuid.UUID       `json:"id"`
	TenantID           uuid.UUID       `json:"tenant_id"`
	WorkerTypeID       uuid.UUID       `json:"worker_type_id"`
	EmployeeID         *uuid.UUID      `json:"employee_id,omitempty"`
	EmployeeUserID     *uuid.UUID      `json:"employee_user_id,omitempty"`
	WorkerCode         *string         `json:"worker_code,omitempty"`
	DisplayName        string          `json:"display_name"`
	LegalName          *string         `json:"legal_name,omitempty"`
	Email              *string         `json:"email,omitempty"`
	Mobile             *string         `json:"mobile,omitempty"`
	ProfileStatus      string          `json:"profile_status"`
	StartDate          *time.Time      `json:"start_date,omitempty"`
	EndDate            *time.Time      `json:"end_date,omitempty"`
	BranchID           *uuid.UUID      `json:"branch_id,omitempty"`
	DepartmentID       *uuid.UUID      `json:"department_id,omitempty"`
	ReportingManagerID *uuid.UUID      `json:"reporting_manager_id,omitempty"`
	WorkLocationLabel  *string         `json:"work_location_label,omitempty"`
	SourcePartner      *string         `json:"source_partner,omitempty"`
	ExternalReference  *string         `json:"external_reference,omitempty"`
	ComplianceStatus   string          `json:"compliance_status"`
	PayrollStatus      string          `json:"payroll_status"`
	Notes              *string         `json:"notes,omitempty"`
	Metadata           json.RawMessage `json:"metadata,omitempty"`
	Inactive           bool            `json:"inactive"`
	CreatedAt          time.Time       `json:"created_at"`
	CreatedBy          *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt          time.Time       `json:"updated_at"`
	UpdatedBy          *uuid.UUID      `json:"updated_by,omitempty"`
}

type WorkerProfileListItem struct {
	WorkerProfile
	WorkerTypeCode      string  `json:"worker_type_code"`
	WorkerTypeName      string  `json:"worker_type_name"`
	ClassificationGroup string  `json:"classification_group"`
	AttendanceMode      string  `json:"attendance_mode"`
	PayMode             string  `json:"pay_mode"`
	EmployeeCode        *string `json:"employee_code,omitempty"`
	BranchName          *string `json:"branch_name,omitempty"`
	DepartmentName      *string `json:"department_name,omitempty"`
}

type WorkerProfileInput struct {
	TenantID           uuid.UUID
	WorkerTypeID       uuid.UUID
	EmployeeID         *uuid.UUID
	EmployeeUserID     *uuid.UUID
	WorkerCode         *string
	DisplayName        string
	LegalName          *string
	Email              *string
	Mobile             *string
	ProfileStatus      string
	StartDate          *time.Time
	EndDate            *time.Time
	BranchID           *uuid.UUID
	DepartmentID       *uuid.UUID
	ReportingManagerID *uuid.UUID
	WorkLocationLabel  *string
	SourcePartner      *string
	ExternalReference  *string
	ComplianceStatus   string
	PayrollStatus      string
	Notes              *string
	Metadata           json.RawMessage
}

type WorkerProfileFilter struct {
	TenantID            uuid.UUID
	WorkerTypeID        *uuid.UUID
	ClassificationGroup *string
	ProfileStatus       *string
	Search              *string
}

func NewWorkerProfile(input WorkerProfileInput) (*WorkerProfile, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	if input.WorkerTypeID == uuid.Nil {
		return nil, ErrInvalidWorkerTypeID
	}
	if (input.EmployeeID == nil) != (input.EmployeeUserID == nil) {
		return nil, ErrInvalidWorkerProfileEmployeeLink
	}
	if input.EmployeeID != nil && (*input.EmployeeID == uuid.Nil || *input.EmployeeUserID == uuid.Nil) {
		return nil, ErrInvalidWorkerProfileEmployeeLink
	}
	displayName := strings.TrimSpace(input.DisplayName)
	if displayName == "" {
		return nil, ErrInvalidWorkerProfileName
	}
	status := normalizeWorkerProfileEnum(input.ProfileStatus, WorkerProfileStatusDraft)
	if !containsString([]string{WorkerProfileStatusDraft, WorkerProfileStatusActive, WorkerProfileStatusPaused, WorkerProfileStatusEnded, WorkerProfileStatusBlacklisted}, status) {
		return nil, ErrInvalidWorkerProfileStatus
	}
	complianceStatus := normalizeWorkerProfileEnum(input.ComplianceStatus, WorkerCompliancePending)
	if !containsString([]string{WorkerCompliancePending, WorkerComplianceReady, WorkerComplianceReview, WorkerComplianceBlocked}, complianceStatus) {
		return nil, ErrInvalidWorkerProfileComplianceStatus
	}
	payrollStatus := normalizeWorkerProfileEnum(input.PayrollStatus, WorkerPayrollNotApplicable)
	if !containsString([]string{WorkerPayrollNotApplicable, WorkerPayrollPending, WorkerPayrollReady, WorkerPayrollBlocked}, payrollStatus) {
		return nil, ErrInvalidWorkerProfilePayrollStatus
	}
	if input.StartDate != nil && input.EndDate != nil && input.EndDate.Before(*input.StartDate) {
		return nil, ErrInvalidWorkerProfileDates
	}
	metadata := normalizeWorkerJSONObject(input.Metadata, "{}")
	if !json.Valid(metadata) || !jsonObject(metadata) {
		return nil, ErrInvalidWorkerProfileMetadata
	}
	now := time.Now().UTC()
	return &WorkerProfile{
		TenantID:           input.TenantID,
		WorkerTypeID:       input.WorkerTypeID,
		EmployeeID:         cleanUUIDOptional(input.EmployeeID),
		EmployeeUserID:     cleanUUIDOptional(input.EmployeeUserID),
		WorkerCode:         cleanOptional(input.WorkerCode),
		DisplayName:        displayName,
		LegalName:          cleanOptional(input.LegalName),
		Email:              cleanOptional(input.Email),
		Mobile:             cleanOptional(input.Mobile),
		ProfileStatus:      status,
		StartDate:          datePtrUTC(input.StartDate),
		EndDate:            datePtrUTC(input.EndDate),
		BranchID:           cleanUUIDOptional(input.BranchID),
		DepartmentID:       cleanUUIDOptional(input.DepartmentID),
		ReportingManagerID: cleanUUIDOptional(input.ReportingManagerID),
		WorkLocationLabel:  cleanOptional(input.WorkLocationLabel),
		SourcePartner:      cleanOptional(input.SourcePartner),
		ExternalReference:  cleanOptional(input.ExternalReference),
		ComplianceStatus:   complianceStatus,
		PayrollStatus:      payrollStatus,
		Notes:              cleanOptional(input.Notes),
		Metadata:           metadata,
		CreatedAt:          now,
		UpdatedAt:          now,
	}, nil
}

func normalizeWorkerProfileEnum(value string, fallback string) string {
	clean := strings.ToLower(strings.TrimSpace(value))
	if clean == "" {
		return fallback
	}
	return clean
}

func datePtrUTC(value *time.Time) *time.Time {
	if value == nil || value.IsZero() {
		return nil
	}
	clean := time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, time.UTC)
	return &clean
}
