package domain

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidPayGroup       = errors.New("pay group is invalid")
	ErrPayGroupNotFound      = errors.New("pay group not found")
	ErrInvalidPayGroupMember = errors.New("pay group member is invalid")
	ErrInvalidPayRun         = errors.New("pay run is invalid")
	ErrPayRunNotFound        = errors.New("pay run not found")
	ErrPayRunLocked          = errors.New("pay run is locked")
	ErrPayRunBlocked         = errors.New("pay run has blocked employees")
)

const (
	PayGroupAll            = "all"
	PayGroupBranch         = "branch"
	PayGroupDepartment     = "department"
	PayGroupEmploymentType = "employment_type"
	PayGroupReportingTag   = "reporting_tag"
	PayGroupManual         = "manual"
	PayGroupMixed          = "mixed"

	PayGroupManualInclude = "manual_include"
	PayGroupManualExclude = "manual_exclude"

	PayRunDraft          = "draft"
	PayRunReadinessReady = "readiness_ready"
	PayRunBlocked        = "blocked"
	PayRunFrozen         = "frozen"
	PayRunProcessing     = "processing"
	PayRunGenerated      = "generated"
	PayRunLocked         = "locked"
	PayRunUnlocked       = "unlocked"
	PayRunFailed         = "failed"

	PayRunEmployeePending   = "pending"
	PayRunEmployeeReady     = "ready"
	PayRunEmployeeBlocked   = "blocked"
	PayRunEmployeeGenerated = "generated"
	PayRunEmployeeSkipped   = "skipped"
	PayRunEmployeeFailed    = "failed"
)

type PayGroup struct {
	ID               uuid.UUID         `json:"id"`
	TenantID         uuid.UUID         `json:"tenant_id"`
	Code             string            `json:"code"`
	Name             string            `json:"name"`
	Description      *string           `json:"description,omitempty"`
	GroupingType     string            `json:"grouping_type"`
	BranchID         *uuid.UUID        `json:"branch_id,omitempty"`
	DepartmentID     *uuid.UUID        `json:"department_id,omitempty"`
	EmploymentTypeID *uuid.UUID        `json:"employment_type_id,omitempty"`
	ReportingTag     *string           `json:"reporting_tag,omitempty"`
	Rules            json.RawMessage   `json:"rules,omitempty"`
	IsActive         bool              `json:"is_active"`
	Inactive         bool              `json:"inactive"`
	CreatedAt        time.Time         `json:"created_at"`
	CreatedBy        *uuid.UUID        `json:"created_by,omitempty"`
	UpdatedAt        time.Time         `json:"updated_at"`
	UpdatedBy        *uuid.UUID        `json:"updated_by,omitempty"`
	Members          []*PayGroupMember `json:"members,omitempty"`
	EmployeeCount    int32             `json:"employee_count,omitempty"`
}

type PayGroupMember struct {
	ID             uuid.UUID  `json:"id"`
	TenantID       uuid.UUID  `json:"tenant_id"`
	PayGroupID     uuid.UUID  `json:"pay_group_id"`
	UserID         uuid.UUID  `json:"user_id"`
	MembershipType string     `json:"membership_type"`
	EffectiveFrom  *time.Time `json:"effective_from,omitempty"`
	EffectiveTo    *time.Time `json:"effective_to,omitempty"`
	Inactive       bool       `json:"inactive"`
	CreatedAt      time.Time  `json:"created_at"`
	CreatedBy      *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt      time.Time  `json:"updated_at"`
	UpdatedBy      *uuid.UUID `json:"updated_by,omitempty"`
}

type PayGroupEmployee struct {
	EmployeeID         uuid.UUID  `json:"employee_id"`
	UserID             uuid.UUID  `json:"user_id"`
	EmployeeCode       *string    `json:"employee_code,omitempty"`
	Firstname          string     `json:"firstname"`
	Lastname           *string    `json:"lastname,omitempty"`
	BranchID           *uuid.UUID `json:"branch_id,omitempty"`
	BranchName         *string    `json:"branch_name,omitempty"`
	DepartmentID       *uuid.UUID `json:"department_id,omitempty"`
	DepartmentName     *string    `json:"department_name,omitempty"`
	EmploymentTypeID   *uuid.UUID `json:"employment_type_id,omitempty"`
	EmploymentTypeName *string    `json:"employment_type_name,omitempty"`
	MatchSource        string     `json:"match_source"`
}

type PayRun struct {
	ID                  uuid.UUID         `json:"id"`
	TenantID            uuid.UUID         `json:"tenant_id"`
	PayGroupID          uuid.UUID         `json:"pay_group_id"`
	FYID                uuid.UUID         `json:"fy_id"`
	Month               int32             `json:"month"`
	Year                int32             `json:"year"`
	Status              string            `json:"status"`
	EmployeeCount       int32             `json:"employee_count"`
	ReadyCount          int32             `json:"ready_count"`
	BlockedCount        int32             `json:"blocked_count"`
	GeneratedCount      int32             `json:"generated_count"`
	AttendanceFrozenAt  *time.Time        `json:"attendance_frozen_at,omitempty"`
	LOPFrozenAt         *time.Time        `json:"lop_frozen_at,omitempty"`
	AdjustmentsFrozenAt *time.Time        `json:"adjustments_frozen_at,omitempty"`
	GeneratedAt         *time.Time        `json:"generated_at,omitempty"`
	LockedAt            *time.Time        `json:"locked_at,omitempty"`
	LockedBy            *uuid.UUID        `json:"locked_by,omitempty"`
	UnlockedAt          *time.Time        `json:"unlocked_at,omitempty"`
	UnlockedBy          *uuid.UUID        `json:"unlocked_by,omitempty"`
	Readiness           json.RawMessage   `json:"readiness,omitempty"`
	Notes               *string           `json:"notes,omitempty"`
	Inactive            bool              `json:"inactive"`
	CreatedAt           time.Time         `json:"created_at"`
	CreatedBy           *uuid.UUID        `json:"created_by,omitempty"`
	UpdatedAt           time.Time         `json:"updated_at"`
	UpdatedBy           *uuid.UUID        `json:"updated_by,omitempty"`
	Employees           []*PayRunEmployee `json:"employees,omitempty"`
	Events              []*PayRunEvent    `json:"events,omitempty"`
}

type PayRunEmployee struct {
	ID               uuid.UUID  `json:"id"`
	TenantID         uuid.UUID  `json:"tenant_id"`
	PayRunID         uuid.UUID  `json:"pay_run_id"`
	UserID           uuid.UUID  `json:"user_id"`
	ReadinessStatus  string     `json:"readiness_status"`
	BlockerReason    *string    `json:"blocker_reason,omitempty"`
	SalarySlipID     *uuid.UUID `json:"salary_slip_id,omitempty"`
	GeneratedAt      *time.Time `json:"generated_at,omitempty"`
	Inactive         bool       `json:"inactive"`
	CreatedAt        time.Time  `json:"created_at"`
	CreatedBy        *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt        time.Time  `json:"updated_at"`
	UpdatedBy        *uuid.UUID `json:"updated_by,omitempty"`
	EmployeeCode     *string    `json:"employee_code,omitempty"`
	Firstname        string     `json:"firstname,omitempty"`
	Lastname         *string    `json:"lastname,omitempty"`
	BranchName       *string    `json:"branch_name,omitempty"`
	DepartmentName   *string    `json:"department_name,omitempty"`
	GrossAmount      float64    `json:"gross_amount,omitempty"`
	EarningsAmount   float64    `json:"earnings_amount,omitempty"`
	DeductionsAmount float64    `json:"deductions_amount,omitempty"`
	NetAmount        float64    `json:"net_amount,omitempty"`
}

type PayRunEvent struct {
	ID         uuid.UUID       `json:"id"`
	TenantID   uuid.UUID       `json:"tenant_id"`
	PayRunID   uuid.UUID       `json:"pay_run_id"`
	Action     string          `json:"action"`
	FromStatus *string         `json:"from_status,omitempty"`
	ToStatus   *string         `json:"to_status,omitempty"`
	Remarks    *string         `json:"remarks,omitempty"`
	Metadata   json.RawMessage `json:"metadata,omitempty"`
	Inactive   bool            `json:"inactive"`
	CreatedAt  time.Time       `json:"created_at"`
	CreatedBy  *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt  time.Time       `json:"updated_at"`
	UpdatedBy  *uuid.UUID      `json:"updated_by,omitempty"`
}

type PayRunInput struct {
	ID           uuid.UUID       `json:"id"`
	TenantID     uuid.UUID       `json:"tenant_id"`
	PayRunID     uuid.UUID       `json:"pay_run_id"`
	UserID       uuid.UUID       `json:"user_id"`
	InputType    string          `json:"input_type"`
	SourceType   string          `json:"source_type"`
	SourceID     *uuid.UUID      `json:"source_id,omitempty"`
	Description  string          `json:"description"`
	Quantity     *float64        `json:"quantity,omitempty"`
	Amount       *float64        `json:"amount,omitempty"`
	Metadata     json.RawMessage `json:"metadata,omitempty"`
	Inactive     bool            `json:"inactive"`
	CreatedAt    time.Time       `json:"created_at"`
	CreatedBy    *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt    time.Time       `json:"updated_at"`
	UpdatedBy    *uuid.UUID      `json:"updated_by,omitempty"`
	EmployeeCode *string         `json:"employee_code,omitempty"`
	Firstname    string          `json:"firstname,omitempty"`
	Lastname     *string         `json:"lastname,omitempty"`
}

type PayRunComponent struct {
	ID               uuid.UUID       `json:"id"`
	TenantID         uuid.UUID       `json:"tenant_id"`
	PayRunID         uuid.UUID       `json:"pay_run_id"`
	UserID           uuid.UUID       `json:"user_id"`
	ComponentType    string          `json:"component_type"`
	Code             string          `json:"code"`
	Name             string          `json:"name"`
	Amount           float64         `json:"amount"`
	SourceInputID    *uuid.UUID      `json:"source_input_id,omitempty"`
	SalaryTemplateID *uuid.UUID      `json:"salary_template_id,omitempty"`
	Taxable          bool            `json:"taxable"`
	Statutory        bool            `json:"statutory"`
	EmployerCost     bool            `json:"employer_cost"`
	SortOrder        int32           `json:"sort_order"`
	Metadata         json.RawMessage `json:"metadata,omitempty"`
	Inactive         bool            `json:"inactive"`
	CreatedAt        time.Time       `json:"created_at"`
	CreatedBy        *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt        time.Time       `json:"updated_at"`
	UpdatedBy        *uuid.UUID      `json:"updated_by,omitempty"`
	EmployeeCode     *string         `json:"employee_code,omitempty"`
	Firstname        string          `json:"firstname,omitempty"`
	Lastname         *string         `json:"lastname,omitempty"`
}

type PayRunLedgerSummary struct {
	PayRunID           uuid.UUID `json:"pay_run_id"`
	EmployeeCount      int32     `json:"employee_count"`
	DraftEmployeeCount int32     `json:"draft_employee_count"`
	GrossAmount        float64   `json:"gross_amount"`
	TotalEarnings      float64   `json:"total_earnings"`
	TotalDeductions    float64   `json:"total_deductions"`
	NetAmount          float64   `json:"net_amount"`
	EmployerCostAmount float64   `json:"employer_cost_amount"`
	InputCount         int32     `json:"input_count"`
	ComponentCount     int32     `json:"component_count"`
}

type PayRunCommandCenter struct {
	Run        *PayRun              `json:"run"`
	Summary    *PayRunLedgerSummary `json:"summary"`
	Inputs     []*PayRunInput       `json:"inputs,omitempty"`
	Components []*PayRunComponent   `json:"components,omitempty"`
}

func NewPayGroup(item PayGroup) (*PayGroup, error) {
	if item.TenantID == uuid.Nil || strings.TrimSpace(item.Code) == "" || strings.TrimSpace(item.Name) == "" {
		return nil, ErrInvalidPayGroup
	}
	groupingType, ok := normalizePayGroupEnum(item.GroupingType, PayGroupManual, PayGroupAll, PayGroupBranch, PayGroupDepartment, PayGroupEmploymentType, PayGroupReportingTag, PayGroupManual, PayGroupMixed)
	if !ok {
		return nil, ErrInvalidPayGroup
	}
	item.Code = strings.TrimSpace(item.Code)
	item.Name = strings.TrimSpace(item.Name)
	item.GroupingType = groupingType
	item.Rules = normalizePayGroupJSON(item.Rules, "{}")
	if item.ID == uuid.Nil {
		item.ID = uuid.New()
	}
	now := time.Now().UTC()
	item.CreatedAt = now
	item.UpdatedAt = now
	return &item, nil
}

func NewPayGroupMember(item PayGroupMember) (*PayGroupMember, error) {
	if item.TenantID == uuid.Nil || item.PayGroupID == uuid.Nil || item.UserID == uuid.Nil {
		return nil, ErrInvalidPayGroupMember
	}
	membershipType, ok := normalizePayGroupEnum(item.MembershipType, PayGroupManualInclude, PayGroupManualInclude, PayGroupManualExclude)
	if !ok {
		return nil, ErrInvalidPayGroupMember
	}
	if item.EffectiveFrom != nil && item.EffectiveTo != nil && item.EffectiveTo.Before(*item.EffectiveFrom) {
		return nil, ErrInvalidPayGroupMember
	}
	item.MembershipType = membershipType
	if item.ID == uuid.Nil {
		item.ID = uuid.New()
	}
	now := time.Now().UTC()
	item.CreatedAt = now
	item.UpdatedAt = now
	return &item, nil
}

func NewPayRun(item PayRun) (*PayRun, error) {
	if item.TenantID == uuid.Nil || item.PayGroupID == uuid.Nil || item.FYID == uuid.Nil || ValidatePayrollPeriod(item.Month, item.Year) != nil {
		return nil, ErrInvalidPayRun
	}
	status, ok := ValidatePayRunStatus(item.Status)
	if !ok {
		return nil, ErrInvalidPayRun
	}
	item.Status = status
	item.Readiness = normalizePayGroupJSON(item.Readiness, "{}")
	if item.ID == uuid.Nil {
		item.ID = uuid.New()
	}
	now := time.Now().UTC()
	item.CreatedAt = now
	item.UpdatedAt = now
	return &item, nil
}

func ValidatePayRunStatus(value string) (string, bool) {
	return normalizePayGroupEnum(value, PayRunDraft, PayRunDraft, PayRunReadinessReady, PayRunBlocked, PayRunFrozen, PayRunProcessing, PayRunGenerated, PayRunLocked, PayRunUnlocked, PayRunFailed)
}

func ValidatePayRunEmployeeStatus(value string) (string, bool) {
	return normalizePayGroupEnum(value, PayRunEmployeePending, PayRunEmployeePending, PayRunEmployeeReady, PayRunEmployeeBlocked, PayRunEmployeeGenerated, PayRunEmployeeSkipped, PayRunEmployeeFailed)
}

func NewPayRunEmployee(item PayRunEmployee) (*PayRunEmployee, error) {
	if item.TenantID == uuid.Nil || item.PayRunID == uuid.Nil || item.UserID == uuid.Nil {
		return nil, ErrInvalidPayRun
	}
	status, ok := ValidatePayRunEmployeeStatus(item.ReadinessStatus)
	if !ok {
		return nil, ErrInvalidPayRun
	}
	item.ReadinessStatus = status
	if item.ID == uuid.Nil {
		item.ID = uuid.New()
	}
	now := time.Now().UTC()
	item.CreatedAt = now
	item.UpdatedAt = now
	return &item, nil
}

func NewPayRunEvent(item PayRunEvent) (*PayRunEvent, error) {
	if item.TenantID == uuid.Nil || item.PayRunID == uuid.Nil || strings.TrimSpace(item.Action) == "" {
		return nil, ErrInvalidPayRun
	}
	item.Action = strings.TrimSpace(item.Action)
	item.Metadata = normalizePayGroupJSON(item.Metadata, "{}")
	if item.ID == uuid.Nil {
		item.ID = uuid.New()
	}
	now := time.Now().UTC()
	item.CreatedAt = now
	item.UpdatedAt = now
	return &item, nil
}

func normalizePayGroupEnum(value string, fallback string, allowed ...string) (string, bool) {
	cleaned := strings.TrimSpace(strings.ToLower(value))
	if cleaned == "" {
		cleaned = fallback
	}
	for _, candidate := range allowed {
		if cleaned == candidate {
			return cleaned, true
		}
	}
	return "", false
}

func normalizePayGroupJSON(value json.RawMessage, fallback string) json.RawMessage {
	if len(value) == 0 || !json.Valid(value) {
		return json.RawMessage(fallback)
	}
	return value
}
