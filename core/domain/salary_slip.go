package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidSalarySlip       = errors.New("invalid salary slip")
	ErrInvalidSalarySlipID     = errors.New("salary_slip_id is required")
	ErrSalarySlipExists        = errors.New("salary slip already exists")
	ErrSalarySlipNotFound      = errors.New("salary slip not found")
	ErrSalarySlipFormatMissing = errors.New("salary slip format is not configured")
)

type SalarySlip struct {
	ID              uuid.UUID          `json:"id"`
	TenantID        uuid.UUID          `json:"tenant_id"`
	UserID          uuid.UUID          `json:"user_id"`
	FYID            uuid.UUID          `json:"fy_id"`
	TemplateID      uuid.UUID          `json:"template_id"`
	Month           int32              `json:"month"`
	Year            int32              `json:"year"`
	GrossSalary     float64            `json:"gross_salary"`
	TotalEarnings   float64            `json:"total_earnings"`
	TotalDeductions float64            `json:"total_deductions"`
	AbsentDeduction float64            `json:"absent_deduction"`
	NetSalary       float64            `json:"net_salary"`
	AbsentDays      int32              `json:"absent_days"`
	PresentDays     int32              `json:"present_days"`
	TotalDays       int32              `json:"total_days"`
	LWPDays         float64            `json:"lwp_days"`
	NoOfPHWEO       int32              `json:"no_of_ph_weo"`
	IsSpecial       bool               `json:"is_special"`
	IsRegenerated   bool               `json:"is_regenerated"`
	PDFPath         *string            `json:"pdf_path,omitempty"`
	Inactive        bool               `json:"inactive"`
	CreatedAt       time.Time          `json:"created_at"`
	CreatedBy       *uuid.UUID         `json:"created_by,omitempty"`
	UpdatedAt       time.Time          `json:"updated_at"`
	UpdatedBy       *uuid.UUID         `json:"updated_by,omitempty"`
	Items           []*SalarySlipItem  `json:"items,omitempty"`
	Leaves          []*SalarySlipLeave `json:"leaves,omitempty"`
}

type SalarySlipItem struct {
	ID        uuid.UUID  `json:"id"`
	TenantID  uuid.UUID  `json:"tenant_id"`
	SlipID    uuid.UUID  `json:"slip_id"`
	ItemType  string     `json:"item_type"`
	Code      string     `json:"code"`
	Name      string     `json:"name"`
	Amount    float64    `json:"amount"`
	SortOrder int32      `json:"sort_order"`
	Inactive  bool       `json:"inactive"`
	CreatedAt time.Time  `json:"created_at"`
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt time.Time  `json:"updated_at"`
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
}

type SalarySlipLeave struct {
	ID            uuid.UUID  `json:"id"`
	TenantID      uuid.UUID  `json:"tenant_id"`
	SlipID        uuid.UUID  `json:"slip_id"`
	LeaveTypeID   uuid.UUID  `json:"leave_type_id"`
	LeaveTypeName *string    `json:"leave_type_name,omitempty"`
	TotalDays     float64    `json:"total_days"`
	UsedDays      float64    `json:"used_days"`
	BalanceDays   float64    `json:"balance_days"`
	Inactive      bool       `json:"inactive"`
	CreatedAt     time.Time  `json:"created_at"`
	CreatedBy     *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt     time.Time  `json:"updated_at"`
	UpdatedBy     *uuid.UUID `json:"updated_by,omitempty"`
}

type SalarySlipFormat struct {
	ID                        uuid.UUID      `json:"id"`
	TenantID                  uuid.UUID      `json:"tenant_id"`
	Title                     string         `json:"title"`
	Subtitle                  *string        `json:"subtitle,omitempty"`
	LogoPath                  *string        `json:"logo_path,omitempty"`
	PrimaryColor              string         `json:"primary_color"`
	AccentColor               string         `json:"accent_color"`
	ShowLeaveBalance          bool           `json:"show_leave_balance"`
	ShowYTDSummary            bool           `json:"show_ytd_summary"`
	ShowEmployeeBank          bool           `json:"show_employee_bank"`
	ShowEmployerContributions bool           `json:"show_employer_contributions"`
	FooterText                *string        `json:"footer_text,omitempty"`
	CustomFields              map[string]any `json:"custom_fields"`
	Inactive                  bool           `json:"inactive"`
	CreatedAt                 time.Time      `json:"created_at"`
	CreatedBy                 *uuid.UUID     `json:"created_by,omitempty"`
	UpdatedAt                 time.Time      `json:"updated_at"`
	UpdatedBy                 *uuid.UUID     `json:"updated_by,omitempty"`
}

type SalarySlipFormatInput struct {
	TenantID                  uuid.UUID
	Title                     string
	Subtitle                  *string
	LogoPath                  *string
	PrimaryColor              string
	AccentColor               string
	ShowLeaveBalance          bool
	ShowYTDSummary            bool
	ShowEmployeeBank          bool
	ShowEmployerContributions bool
	FooterText                *string
	CustomFields              map[string]any
}

func NewSalarySlipFormat(input SalarySlipFormatInput) (*SalarySlipFormat, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	title := strings.TrimSpace(input.Title)
	if title == "" {
		title = "Salary Slip"
	}
	primary := strings.TrimSpace(input.PrimaryColor)
	if primary == "" {
		primary = "#111827"
	}
	accent := strings.TrimSpace(input.AccentColor)
	if accent == "" {
		accent = "#588368"
	}
	fields := input.CustomFields
	if fields == nil {
		fields = map[string]any{}
	}
	return &SalarySlipFormat{TenantID: input.TenantID, Title: title, Subtitle: cleanOptional(input.Subtitle), LogoPath: cleanOptional(input.LogoPath), PrimaryColor: primary, AccentColor: accent, ShowLeaveBalance: input.ShowLeaveBalance, ShowYTDSummary: input.ShowYTDSummary, ShowEmployeeBank: input.ShowEmployeeBank, ShowEmployerContributions: input.ShowEmployerContributions, FooterText: cleanOptional(input.FooterText), CustomFields: fields, CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()}, nil
}
