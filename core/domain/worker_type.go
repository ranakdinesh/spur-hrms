package domain

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	WorkerTypePermanentFullTime  = "permanent_fulltime"
	WorkerTypePermanentPartTime  = "permanent_parttime"
	WorkerTypeFixedTermContract  = "fixed_term_contract"
	WorkerTypeProjectBased       = "project_based"
	WorkerTypeFreelancerGig      = "freelancer_gig"
	WorkerTypeIntern             = "intern"
	WorkerTypeConsultantRetainer = "consultant_retainer"
	WorkerTypeAgencyStaff        = "agency_staff"

	WorkerClassEmployee   = "employee"
	WorkerClassContractor = "contractor"
	WorkerClassTrainee    = "trainee"
	WorkerClassAgency     = "agency"

	WorkerAttendanceCheckInOut = "checkin_checkout"
	WorkerAttendanceHours      = "hours_logged"
	WorkerAttendanceMilestone  = "milestone_only"
	WorkerAttendanceNone       = "none"

	WorkerPayMonthlySalary = "monthly_salary"
	WorkerPayHourly        = "hourly"
	WorkerPayProject       = "project_milestone"
	WorkerPayInvoice       = "invoice"
	WorkerPayRetainer      = "retainer"
	WorkerPayStipend       = "stipend"

	WorkerTDS192  = "192"
	WorkerTDS194C = "194C"
	WorkerTDS194J = "194J"
	WorkerTDS194I = "194I"
	WorkerTDSNone = "none"

	WorkerRuleManualGuidance = "manual_guidance"
	WorkerRuleCompliance     = "compliance"
	WorkerRulePayroll        = "payroll"
	WorkerRuleAttendance     = "attendance"
)

var (
	ErrInvalidWorkerTypeID                     = errors.New("worker_type_id is required")
	ErrInvalidWorkerTypeCode                   = errors.New("worker type code is invalid")
	ErrInvalidWorkerTypeName                   = errors.New("worker type name is required")
	ErrInvalidWorkerClassificationGroup        = errors.New("worker type classification_group is invalid")
	ErrInvalidWorkerAttendanceMode             = errors.New("worker type attendance_mode is invalid")
	ErrInvalidWorkerPayMode                    = errors.New("worker type pay_mode is invalid")
	ErrInvalidWorkerTDSSection                 = errors.New("worker type tds_section is invalid")
	ErrInvalidWorkerSortOrder                  = errors.New("worker type sort_order must be between 0 and 9999")
	ErrInvalidWorkerStatutoryDefaults          = errors.New("worker type statutory_defaults must be valid JSON object")
	ErrInvalidWorkerClassificationRuleID       = errors.New("worker_classification_rule_id is required")
	ErrInvalidWorkerClassificationRuleName     = errors.New("worker classification rule name is required")
	ErrInvalidWorkerClassificationRuleType     = errors.New("worker classification rule type is invalid")
	ErrInvalidWorkerClassificationRulePriority = errors.New("worker classification rule priority must be between 1 and 9999")
	ErrInvalidWorkerClassificationRuleJSON     = errors.New("worker classification rule conditions and outcome must be valid JSON objects")
)

type WorkerType struct {
	ID                  uuid.UUID       `json:"id"`
	TenantID            uuid.UUID       `json:"tenant_id"`
	Code                string          `json:"code"`
	Name                string          `json:"name"`
	ClassificationGroup string          `json:"classification_group"`
	Description         *string         `json:"description,omitempty"`
	AttendanceMode      string          `json:"attendance_mode"`
	PayMode             string          `json:"pay_mode"`
	TDSSection          string          `json:"tds_section"`
	PFApplicable        bool            `json:"pf_applicable"`
	ESICApplicable      bool            `json:"esic_applicable"`
	PTApplicable        bool            `json:"pt_applicable"`
	LWFApplicable       bool            `json:"lwf_applicable"`
	CLRAApplicable      bool            `json:"clra_applicable"`
	LeaveApplicable     bool            `json:"leave_applicable"`
	OvertimeApplicable  bool            `json:"overtime_applicable"`
	RequiresAgreement   bool            `json:"requires_agreement"`
	RequiresInvoice     bool            `json:"requires_invoice"`
	RequiresAttendance  bool            `json:"requires_attendance"`
	StatutoryDefaults   json.RawMessage `json:"statutory_defaults,omitempty"`
	ComplianceNotes     *string         `json:"compliance_notes,omitempty"`
	IsSystemDefault     bool            `json:"is_system_default"`
	SortOrder           int32           `json:"sort_order"`
	Inactive            bool            `json:"inactive"`
	CreatedAt           time.Time       `json:"created_at"`
	CreatedBy           *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt           time.Time       `json:"updated_at"`
	UpdatedBy           *uuid.UUID      `json:"updated_by,omitempty"`
}

type WorkerTypeInput struct {
	TenantID            uuid.UUID
	Code                string
	Name                string
	ClassificationGroup string
	Description         *string
	AttendanceMode      string
	PayMode             string
	TDSSection          string
	PFApplicable        bool
	ESICApplicable      bool
	PTApplicable        bool
	LWFApplicable       bool
	CLRAApplicable      bool
	LeaveApplicable     bool
	OvertimeApplicable  bool
	RequiresAgreement   bool
	RequiresInvoice     bool
	RequiresAttendance  bool
	StatutoryDefaults   json.RawMessage
	ComplianceNotes     *string
	IsSystemDefault     bool
	SortOrder           int32
}

type WorkerClassificationRule struct {
	ID           uuid.UUID       `json:"id"`
	TenantID     uuid.UUID       `json:"tenant_id"`
	WorkerTypeID uuid.UUID       `json:"worker_type_id"`
	RuleName     string          `json:"rule_name"`
	RuleType     string          `json:"rule_type"`
	Priority     int32           `json:"priority"`
	Conditions   json.RawMessage `json:"conditions,omitempty"`
	Outcome      json.RawMessage `json:"outcome,omitempty"`
	Notes        *string         `json:"notes,omitempty"`
	Inactive     bool            `json:"inactive"`
	CreatedAt    time.Time       `json:"created_at"`
	CreatedBy    *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt    time.Time       `json:"updated_at"`
	UpdatedBy    *uuid.UUID      `json:"updated_by,omitempty"`
}

type WorkerClassificationRuleInput struct {
	TenantID     uuid.UUID
	WorkerTypeID uuid.UUID
	RuleName     string
	RuleType     string
	Priority     int32
	Conditions   json.RawMessage
	Outcome      json.RawMessage
	Notes        *string
}

func NewWorkerType(input WorkerTypeInput) (*WorkerType, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	code := normalizeWorkerEnum(input.Code)
	if !containsString(workerTypeCodes(), code) {
		return nil, ErrInvalidWorkerTypeCode
	}
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, ErrInvalidWorkerTypeName
	}
	group := normalizeWorkerEnum(input.ClassificationGroup)
	if !containsString([]string{WorkerClassEmployee, WorkerClassContractor, WorkerClassTrainee, WorkerClassAgency}, group) {
		return nil, ErrInvalidWorkerClassificationGroup
	}
	attendanceMode := normalizeWorkerEnum(input.AttendanceMode)
	if !containsString([]string{WorkerAttendanceCheckInOut, WorkerAttendanceHours, WorkerAttendanceMilestone, WorkerAttendanceNone}, attendanceMode) {
		return nil, ErrInvalidWorkerAttendanceMode
	}
	payMode := normalizeWorkerEnum(input.PayMode)
	if !containsString([]string{WorkerPayMonthlySalary, WorkerPayHourly, WorkerPayProject, WorkerPayInvoice, WorkerPayRetainer, WorkerPayStipend}, payMode) {
		return nil, ErrInvalidWorkerPayMode
	}
	tdsSection := strings.TrimSpace(input.TDSSection)
	if tdsSection == "" {
		tdsSection = WorkerTDSNone
	}
	if !containsString([]string{WorkerTDS192, WorkerTDS194C, WorkerTDS194J, WorkerTDS194I, WorkerTDSNone}, tdsSection) {
		return nil, ErrInvalidWorkerTDSSection
	}
	if input.SortOrder < 0 || input.SortOrder > 9999 {
		return nil, ErrInvalidWorkerSortOrder
	}
	statutoryDefaults := normalizeWorkerJSONObject(input.StatutoryDefaults, "{}")
	if !json.Valid(statutoryDefaults) || !jsonObject(statutoryDefaults) {
		return nil, ErrInvalidWorkerStatutoryDefaults
	}
	now := time.Now().UTC()
	return &WorkerType{
		TenantID:            input.TenantID,
		Code:                code,
		Name:                name,
		ClassificationGroup: group,
		Description:         cleanOptional(input.Description),
		AttendanceMode:      attendanceMode,
		PayMode:             payMode,
		TDSSection:          tdsSection,
		PFApplicable:        input.PFApplicable,
		ESICApplicable:      input.ESICApplicable,
		PTApplicable:        input.PTApplicable,
		LWFApplicable:       input.LWFApplicable,
		CLRAApplicable:      input.CLRAApplicable,
		LeaveApplicable:     input.LeaveApplicable,
		OvertimeApplicable:  input.OvertimeApplicable,
		RequiresAgreement:   input.RequiresAgreement,
		RequiresInvoice:     input.RequiresInvoice,
		RequiresAttendance:  input.RequiresAttendance,
		StatutoryDefaults:   statutoryDefaults,
		ComplianceNotes:     cleanOptional(input.ComplianceNotes),
		IsSystemDefault:     input.IsSystemDefault,
		SortOrder:           input.SortOrder,
		CreatedAt:           now,
		UpdatedAt:           now,
	}, nil
}

func NewWorkerClassificationRule(input WorkerClassificationRuleInput) (*WorkerClassificationRule, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	if input.WorkerTypeID == uuid.Nil {
		return nil, ErrInvalidWorkerTypeID
	}
	ruleName := strings.TrimSpace(input.RuleName)
	if ruleName == "" {
		return nil, ErrInvalidWorkerClassificationRuleName
	}
	ruleType := normalizeWorkerEnum(input.RuleType)
	if !containsString([]string{WorkerRuleManualGuidance, WorkerRuleCompliance, WorkerRulePayroll, WorkerRuleAttendance}, ruleType) {
		return nil, ErrInvalidWorkerClassificationRuleType
	}
	if input.Priority < 1 || input.Priority > 9999 {
		return nil, ErrInvalidWorkerClassificationRulePriority
	}
	conditions := normalizeWorkerJSONObject(input.Conditions, "{}")
	outcome := normalizeWorkerJSONObject(input.Outcome, "{}")
	if !json.Valid(conditions) || !jsonObject(conditions) || !json.Valid(outcome) || !jsonObject(outcome) {
		return nil, ErrInvalidWorkerClassificationRuleJSON
	}
	now := time.Now().UTC()
	return &WorkerClassificationRule{
		TenantID:     input.TenantID,
		WorkerTypeID: input.WorkerTypeID,
		RuleName:     ruleName,
		RuleType:     ruleType,
		Priority:     input.Priority,
		Conditions:   conditions,
		Outcome:      outcome,
		Notes:        cleanOptional(input.Notes),
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

func DefaultWorkerTypeInputs(tenantID uuid.UUID) []WorkerTypeInput {
	return []WorkerTypeInput{
		workerTypeInput(tenantID, WorkerTypePermanentFullTime, "Permanent Full-Time", WorkerClassEmployee, WorkerAttendanceCheckInOut, WorkerPayMonthlySalary, WorkerTDS192, true, true, true, true, false, true, true, true, false, true, 10, "Standard payroll employee with PF, ESIC, PT and leave defaults enabled where applicable."),
		workerTypeInput(tenantID, WorkerTypePermanentPartTime, "Permanent Part-Time", WorkerClassEmployee, WorkerAttendanceCheckInOut, WorkerPayMonthlySalary, WorkerTDS192, true, true, true, true, false, true, true, true, false, true, 20, "Part-time employee classification; statutory coverage should follow wage, hours and establishment rules."),
		workerTypeInput(tenantID, WorkerTypeFixedTermContract, "Fixed-Term Contract", WorkerClassEmployee, WorkerAttendanceCheckInOut, WorkerPayMonthlySalary, WorkerTDS192, true, true, true, true, false, true, true, true, false, true, 30, "Employee on a fixed tenure contract; agreement dates and statutory coverage must be reviewed."),
		workerTypeInput(tenantID, WorkerTypeProjectBased, "Project-Based Contractor", WorkerClassContractor, WorkerAttendanceMilestone, WorkerPayProject, WorkerTDS194C, false, false, false, false, true, false, false, true, true, false, 40, "Contractor paid against project milestones or work orders; CLRA review may be required."),
		workerTypeInput(tenantID, WorkerTypeFreelancerGig, "Freelancer / Gig Worker", WorkerClassContractor, WorkerAttendanceHours, WorkerPayInvoice, WorkerTDS194J, false, false, false, false, false, false, false, true, true, true, 50, "Independent professional or gig engagement; capture invoice, agreement and service category."),
		workerTypeInput(tenantID, WorkerTypeIntern, "Intern", WorkerClassTrainee, WorkerAttendanceCheckInOut, WorkerPayStipend, WorkerTDSNone, false, false, false, false, false, false, false, true, false, true, 60, "Internship or trainee arrangement; stipend and attendance rules depend on program policy."),
		workerTypeInput(tenantID, WorkerTypeConsultantRetainer, "Consultant Retainer", WorkerClassContractor, WorkerAttendanceNone, WorkerPayRetainer, WorkerTDS194J, false, false, false, false, false, false, false, true, true, false, 70, "Retainer consultant; attendance normally not required but deliverables and invoice controls are."),
		workerTypeInput(tenantID, WorkerTypeAgencyStaff, "Agency Staff", WorkerClassAgency, WorkerAttendanceCheckInOut, WorkerPayInvoice, WorkerTDS194C, false, false, false, false, true, false, true, true, true, true, 80, "Manpower supplied through an agency; gate attendance and CLRA documentation should be tracked."),
	}
}

func DefaultWorkerClassificationRuleInput(tenantID uuid.UUID, workerType *WorkerType) WorkerClassificationRuleInput {
	conditions := map[string]any{"worker_type_code": workerType.Code, "classification_group": workerType.ClassificationGroup}
	outcome := map[string]any{
		"attendance_mode": workerType.AttendanceMode,
		"pay_mode":        workerType.PayMode,
		"tds_section":     workerType.TDSSection,
		"pf_applicable":   workerType.PFApplicable,
		"esic_applicable": workerType.ESICApplicable,
		"clra_applicable": workerType.CLRAApplicable,
	}
	return WorkerClassificationRuleInput{
		TenantID:     tenantID,
		WorkerTypeID: workerType.ID,
		RuleName:     workerType.Name + " default classification",
		RuleType:     WorkerRuleManualGuidance,
		Priority:     maxInt32(1, workerType.SortOrder),
		Conditions:   mustJSONRaw(conditions),
		Outcome:      mustJSONRaw(outcome),
		Notes:        workerType.ComplianceNotes,
	}
}

func workerTypeInput(tenantID uuid.UUID, code string, name string, group string, attendanceMode string, payMode string, tds string, pf bool, esic bool, pt bool, lwf bool, clra bool, leave bool, overtime bool, agreement bool, invoice bool, attendance bool, sortOrder int32, complianceNotes string) WorkerTypeInput {
	defaults := map[string]any{
		"country":         "IN",
		"requires_review": true,
		"source":          "system_default",
	}
	return WorkerTypeInput{
		TenantID:            tenantID,
		Code:                code,
		Name:                name,
		ClassificationGroup: group,
		AttendanceMode:      attendanceMode,
		PayMode:             payMode,
		TDSSection:          tds,
		PFApplicable:        pf,
		ESICApplicable:      esic,
		PTApplicable:        pt,
		LWFApplicable:       lwf,
		CLRAApplicable:      clra,
		LeaveApplicable:     leave,
		OvertimeApplicable:  overtime,
		RequiresAgreement:   agreement,
		RequiresInvoice:     invoice,
		RequiresAttendance:  attendance,
		StatutoryDefaults:   mustJSONRaw(defaults),
		ComplianceNotes:     &complianceNotes,
		IsSystemDefault:     true,
		SortOrder:           sortOrder,
	}
}

func workerTypeCodes() []string {
	return []string{WorkerTypePermanentFullTime, WorkerTypePermanentPartTime, WorkerTypeFixedTermContract, WorkerTypeProjectBased, WorkerTypeFreelancerGig, WorkerTypeIntern, WorkerTypeConsultantRetainer, WorkerTypeAgencyStaff}
}

func normalizeWorkerEnum(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func normalizeWorkerJSONObject(value json.RawMessage, fallback string) json.RawMessage {
	if len(value) == 0 {
		return json.RawMessage(fallback)
	}
	return value
}

func jsonObject(value json.RawMessage) bool {
	var decoded map[string]any
	return json.Unmarshal(value, &decoded) == nil
}

func mustJSONRaw(value any) json.RawMessage {
	data, err := json.Marshal(value)
	if err != nil {
		return json.RawMessage(`{}`)
	}
	return data
}

func maxInt32(a int32, b int32) int32 {
	if a > b {
		return a
	}
	return b
}
