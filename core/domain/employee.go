package domain

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidEmployeeID              = errors.New("employee_id is required")
	ErrInvalidEmployeeUserID          = errors.New("employee user_id is required")
	ErrInvalidEmployeeFirstName       = errors.New("employee first name is required")
	ErrInvalidEmployeeEmail           = errors.New("employee email is invalid")
	ErrInvalidEmployeeMobile          = errors.New("employee mobile is required")
	ErrInvalidEmployeeExperience      = errors.New("employee experience month must be between 0 and 11 and year cannot be negative")
	ErrInvalidEmployeeDate            = errors.New("employee date is invalid")
	ErrInvalidEmployeeResignation     = errors.New("resignation date must be on or after joining date")
	ErrInvalidEmployeeProbation       = errors.New("employee probation settings are invalid")
	ErrEmployeeIdentityPortMissing    = errors.New("employee identity port is not configured")
	ErrEmployeeCodeAlreadyExists      = errors.New("employee code already exists")
	ErrInvalidDocumentTypeID          = errors.New("document_type_id is required")
	ErrInvalidDocumentTypeName        = errors.New("document type name is required")
	ErrDocumentTypeNotFound           = errors.New("document type not found")
	ErrInvalidEmployeeDocumentID      = errors.New("employee document id is required")
	ErrInvalidEmployeeDocument        = errors.New("employee document title or file is required")
	ErrEmployeeDocumentNotFound       = errors.New("employee document not found")
	ErrEmployeeDocumentStorageMissing = errors.New("employee document storage is not configured")
	ErrInvalidEmployeeDocumentStatus  = errors.New("employee document status is invalid")
	ErrEmployeeDocumentFileTooLarge   = errors.New("employee document file exceeds configured size limit")
	ErrEmployeeDocumentFileTypeDenied = errors.New("employee document file type is not allowed")
	ErrEmployeeDocumentApprovedLocked = errors.New("approved employee document cannot be uploaded again")
)

var employeeEmailPattern = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

const (
	EmployeeDocumentStatusPendingReview         = "pending_review"
	EmployeeDocumentStatusApproved              = "approved"
	EmployeeDocumentStatusRejected              = "rejected"
	EmployeeDocumentStatusResubmissionRequested = "resubmission_requested"

	EmployeeProbationNotApplicable = "not_applicable"
	EmployeeProbationProbation     = "probation"
	EmployeeProbationConfirmed     = "confirmed"
	EmployeeProbationExtended      = "extended"

	PayrollStaffProbationDurationDays = 180
)

func ValidateEmployeeDocumentStatus(value string) (string, error) {
	status := strings.TrimSpace(value)
	switch status {
	case EmployeeDocumentStatusPendingReview, EmployeeDocumentStatusApproved, EmployeeDocumentStatusRejected, EmployeeDocumentStatusResubmissionRequested:
		return status, nil
	default:
		return "", ErrInvalidEmployeeDocumentStatus
	}
}

func ValidateEmployeeDocumentReviewStatus(value string) (string, error) {
	status := strings.TrimSpace(value)
	switch status {
	case EmployeeDocumentStatusApproved, EmployeeDocumentStatusRejected, EmployeeDocumentStatusResubmissionRequested:
		return status, nil
	default:
		return "", ErrInvalidEmployeeDocumentStatus
	}
}

type Employee struct {
	ID                    uuid.UUID  `json:"id"`
	TenantID              uuid.UUID  `json:"tenant_id"`
	UserID                uuid.UUID  `json:"user_id"`
	EmployeeCode          *string    `json:"employee_code,omitempty"`
	Firstname             string     `json:"firstname"`
	MiddleName            *string    `json:"middle_name,omitempty"`
	Lastname              *string    `json:"lastname,omitempty"`
	Email                 *string    `json:"email,omitempty"`
	Mobile                *string    `json:"mobile,omitempty"`
	DOB                   *time.Time `json:"dob,omitempty"`
	Gender                *string    `json:"gender,omitempty"`
	MaritalStatus         *string    `json:"marital_status,omitempty"`
	BloodGroup            *string    `json:"blood_group,omitempty"`
	ProfilePhotoPath      *string    `json:"profile_photo_path,omitempty"`
	Address               *string    `json:"address,omitempty"`
	City                  *string    `json:"city,omitempty"`
	State                 *string    `json:"state,omitempty"`
	Country               *string    `json:"country,omitempty"`
	Pincode               *string    `json:"pincode,omitempty"`
	EmergencyContact      *string    `json:"emergency_contact,omitempty"`
	JoiningDate           *time.Time `json:"joining_date,omitempty"`
	ResignationDate       *time.Time `json:"resignation_date,omitempty"`
	DepartmentID          *uuid.UUID `json:"department_id,omitempty"`
	BranchID              *uuid.UUID `json:"branch_id,omitempty"`
	DesignationID         *uuid.UUID `json:"designation_id,omitempty"`
	ReportingManagerID    *uuid.UUID `json:"reporting_manager_id,omitempty"`
	EmploymentTypeID      *uuid.UUID `json:"employment_type_id,omitempty"`
	Role                  *string    `json:"role,omitempty"`
	Grade                 *string    `json:"grade,omitempty"`
	ExperienceYear        int32      `json:"experience_year"`
	ExperienceMonth       int32      `json:"experience_month"`
	ProbationStatus       string     `json:"probation_status"`
	ProbationStartDate    *time.Time `json:"probation_start_date,omitempty"`
	ProbationEndDate      *time.Time `json:"probation_end_date,omitempty"`
	ProbationDurationDays int32      `json:"probation_duration_days"`
	ProbationConfirmedAt  *time.Time `json:"probation_confirmed_at,omitempty"`
	IsPayrollStaff        bool       `json:"is_payroll_staff"`
	Inactive              bool       `json:"inactive"`
	CreatedAt             time.Time  `json:"created_at"`
	CreatedBy             *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt             time.Time  `json:"updated_at"`
	UpdatedBy             *uuid.UUID `json:"updated_by,omitempty"`
}

type EmployeeListItem struct {
	Employee
	DepartmentName     *string `json:"department_name,omitempty"`
	BranchName         *string `json:"branch_name,omitempty"`
	DesignationName    *string `json:"designation_name,omitempty"`
	AttendanceRequired bool    `json:"attendance_required"`
	EmploymentTypeName *string `json:"employment_type_name,omitempty"`
}

type EmployeeStatutory struct {
	ID             uuid.UUID  `json:"id"`
	TenantID       uuid.UUID  `json:"tenant_id"`
	UserID         uuid.UUID  `json:"user_id"`
	PFNo           *string    `json:"pf_no,omitempty"`
	UANNo          *string    `json:"uan_no,omitempty"`
	ESICNo         *string    `json:"esic_no,omitempty"`
	PAN            *string    `json:"pan,omitempty"`
	Aadhaar        *string    `json:"aadhaar,omitempty"`
	PTApplicable   bool       `json:"pt_applicable"`
	PFApplicable   bool       `json:"pf_applicable"`
	ESICApplicable bool       `json:"esic_applicable"`
	LWFApplicable  bool       `json:"lwf_applicable"`
	Inactive       bool       `json:"inactive"`
	CreatedAt      time.Time  `json:"created_at"`
	CreatedBy      *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt      time.Time  `json:"updated_at"`
	UpdatedBy      *uuid.UUID `json:"updated_by,omitempty"`
}

type EmployeeBank struct {
	ID            uuid.UUID  `json:"id"`
	TenantID      uuid.UUID  `json:"tenant_id"`
	UserID        uuid.UUID  `json:"user_id"`
	BankName      *string    `json:"bank_name,omitempty"`
	AccountNumber *string    `json:"account_number,omitempty"`
	IFSCCode      *string    `json:"ifsc_code,omitempty"`
	AccountType   *string    `json:"account_type,omitempty"`
	BranchName    *string    `json:"branch_name,omitempty"`
	IsPrimary     bool       `json:"is_primary"`
	Inactive      bool       `json:"inactive"`
	CreatedAt     time.Time  `json:"created_at"`
	CreatedBy     *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt     time.Time  `json:"updated_at"`
	UpdatedBy     *uuid.UUID `json:"updated_by,omitempty"`
}

type DocumentType struct {
	ID                  uuid.UUID  `json:"id"`
	TenantID            uuid.UUID  `json:"tenant_id"`
	Name                string     `json:"name"`
	Description         *string    `json:"description,omitempty"`
	IsRequired          bool       `json:"is_required"`
	Instructions        *string    `json:"instructions,omitempty"`
	AllowedContentTypes string     `json:"allowed_content_types"`
	MaxFileSizeBytes    int64      `json:"max_file_size_bytes"`
	DisplayOrder        int32      `json:"display_order"`
	Inactive            bool       `json:"inactive"`
	CreatedAt           time.Time  `json:"created_at"`
	CreatedBy           *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt           time.Time  `json:"updated_at"`
	UpdatedBy           *uuid.UUID `json:"updated_by,omitempty"`
}

type EmployeeDocument struct {
	ID                  uuid.UUID  `json:"id"`
	TenantID            uuid.UUID  `json:"tenant_id"`
	UserID              uuid.UUID  `json:"user_id"`
	DocumentTypeID      *uuid.UUID `json:"document_type_id,omitempty"`
	DocumentTypeName    *string    `json:"document_type_name,omitempty"`
	Title               *string    `json:"title,omitempty"`
	FilePath            *string    `json:"file_path,omitempty"`
	Status              string     `json:"status"`
	ReviewRemarks       *string    `json:"review_remarks,omitempty"`
	ReviewedBy          *uuid.UUID `json:"reviewed_by,omitempty"`
	ReviewedAt          *time.Time `json:"reviewed_at,omitempty"`
	OriginalFileName    *string    `json:"original_file_name,omitempty"`
	ContentType         *string    `json:"content_type,omitempty"`
	FileSizeBytes       *int64     `json:"file_size_bytes,omitempty"`
	Encrypted           bool       `json:"encrypted"`
	EncryptionAlgorithm string     `json:"encryption_algorithm"`
	Inactive            bool       `json:"inactive"`
	CreatedAt           time.Time  `json:"created_at"`
	CreatedBy           *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt           time.Time  `json:"updated_at"`
	UpdatedBy           *uuid.UUID `json:"updated_by,omitempty"`
}

type EmployeeProfileLookups struct {
	Branches        []*Branch         `json:"branches"`
	Departments     []*Department     `json:"departments"`
	Designations    []*Designation    `json:"designations"`
	EmploymentTypes []*EmploymentType `json:"employment_types"`
	DocumentTypes   []*DocumentType   `json:"document_types"`
}

type EmployeeProfile struct {
	Employee   *EmployeeListItem        `json:"employee"`
	Statutory  *EmployeeStatutory       `json:"statutory,omitempty"`
	Banks      []*EmployeeBank          `json:"banks"`
	Documents  []*EmployeeDocument      `json:"documents"`
	Onboarding EmployeeOnboardingStatus `json:"onboarding"`
	Lookups    EmployeeProfileLookups   `json:"lookups"`
}

type EmployeeOnboardingStatus struct {
	Status                    string   `json:"status"`
	IsComplete                bool     `json:"is_complete"`
	RequiredDocuments         int      `json:"required_documents"`
	UploadedRequiredDocuments int      `json:"uploaded_required_documents"`
	ApprovedRequiredDocuments int      `json:"approved_required_documents"`
	PendingReviewDocuments    int      `json:"pending_review_documents"`
	RejectedDocuments         int      `json:"rejected_documents"`
	MissingRequiredDocuments  []string `json:"missing_required_documents"`
}

type EmployeeInput struct {
	TenantID              uuid.UUID
	UserID                uuid.UUID
	EmployeeCode          *string
	Firstname             string
	MiddleName            *string
	Lastname              *string
	Email                 *string
	Mobile                *string
	DOB                   *time.Time
	Gender                *string
	MaritalStatus         *string
	BloodGroup            *string
	ProfilePhotoPath      *string
	Address               *string
	City                  *string
	State                 *string
	Country               *string
	Pincode               *string
	EmergencyContact      *string
	JoiningDate           *time.Time
	ResignationDate       *time.Time
	DepartmentID          *uuid.UUID
	BranchID              *uuid.UUID
	DesignationID         *uuid.UUID
	ReportingManagerID    *uuid.UUID
	EmploymentTypeID      *uuid.UUID
	Role                  *string
	Grade                 *string
	ExperienceYear        int32
	ExperienceMonth       int32
	ProbationStatus       string
	ProbationStartDate    *time.Time
	ProbationEndDate      *time.Time
	ProbationDurationDays int32
	ProbationConfirmedAt  *time.Time
	IsPayrollStaff        bool
}

func NewEmployee(input EmployeeInput) (*Employee, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	if input.UserID == uuid.Nil {
		return nil, ErrInvalidEmployeeUserID
	}
	firstname := strings.TrimSpace(input.Firstname)
	if firstname == "" {
		return nil, ErrInvalidEmployeeFirstName
	}
	email := cleanOptional(input.Email)
	if email != nil && !employeeEmailPattern.MatchString(*email) {
		return nil, ErrInvalidEmployeeEmail
	}
	mobile := cleanOptional(input.Mobile)
	if mobile == nil {
		return nil, ErrInvalidEmployeeMobile
	}
	if input.ExperienceYear < 0 || input.ExperienceMonth < 0 || input.ExperienceMonth > 11 {
		return nil, ErrInvalidEmployeeExperience
	}
	joiningDate := datePtr(input.JoiningDate)
	resignationDate := datePtr(input.ResignationDate)
	if joiningDate != nil && resignationDate != nil && resignationDate.Before(*joiningDate) {
		return nil, ErrInvalidEmployeeResignation
	}
	probationStatus := normalizeEmployeeProbationStatus(input.ProbationStatus)
	probationStartDate := datePtr(input.ProbationStartDate)
	probationEndDate := datePtr(input.ProbationEndDate)
	probationConfirmedAt := datePtr(input.ProbationConfirmedAt)
	if probationEndDate != nil && probationStartDate != nil && probationEndDate.Before(*probationStartDate) {
		return nil, ErrInvalidEmployeeProbation
	}
	if input.ProbationDurationDays < 0 {
		return nil, ErrInvalidEmployeeProbation
	}
	if input.IsPayrollStaff && input.ProbationDurationDays < PayrollStaffProbationDurationDays {
		return nil, ErrInvalidEmployeeProbation
	}
	now := time.Now().UTC()
	return &Employee{
		TenantID:              input.TenantID,
		UserID:                input.UserID,
		EmployeeCode:          cleanOptional(input.EmployeeCode),
		Firstname:             firstname,
		MiddleName:            cleanOptional(input.MiddleName),
		Lastname:              cleanOptional(input.Lastname),
		Email:                 email,
		Mobile:                mobile,
		DOB:                   datePtr(input.DOB),
		Gender:                cleanOptional(input.Gender),
		MaritalStatus:         cleanOptional(input.MaritalStatus),
		BloodGroup:            cleanOptional(input.BloodGroup),
		ProfilePhotoPath:      cleanOptional(input.ProfilePhotoPath),
		Address:               cleanOptional(input.Address),
		City:                  cleanOptional(input.City),
		State:                 cleanOptional(input.State),
		Country:               cleanOptional(input.Country),
		Pincode:               cleanOptional(input.Pincode),
		EmergencyContact:      cleanOptional(input.EmergencyContact),
		JoiningDate:           joiningDate,
		ResignationDate:       resignationDate,
		DepartmentID:          cleanUUIDOptional(input.DepartmentID),
		BranchID:              cleanUUIDOptional(input.BranchID),
		DesignationID:         cleanUUIDOptional(input.DesignationID),
		ReportingManagerID:    cleanUUIDOptional(input.ReportingManagerID),
		EmploymentTypeID:      cleanUUIDOptional(input.EmploymentTypeID),
		Role:                  cleanOptional(input.Role),
		Grade:                 cleanOptional(input.Grade),
		ExperienceYear:        input.ExperienceYear,
		ExperienceMonth:       input.ExperienceMonth,
		ProbationStatus:       probationStatus,
		ProbationStartDate:    probationStartDate,
		ProbationEndDate:      probationEndDate,
		ProbationDurationDays: input.ProbationDurationDays,
		ProbationConfirmedAt:  probationConfirmedAt,
		IsPayrollStaff:        input.IsPayrollStaff,
		CreatedAt:             now,
		UpdatedAt:             now,
	}, nil
}

func normalizeEmployeeProbationStatus(value string) string {
	switch strings.TrimSpace(value) {
	case EmployeeProbationNotApplicable, EmployeeProbationProbation, EmployeeProbationConfirmed, EmployeeProbationExtended:
		return strings.TrimSpace(value)
	case "":
		return EmployeeProbationConfirmed
	default:
		return EmployeeProbationConfirmed
	}
}

func (e *Employee) IsOnProbation(asOf time.Time) bool {
	if e == nil {
		return false
	}
	if e.ProbationStatus != EmployeeProbationProbation && e.ProbationStatus != EmployeeProbationExtended {
		return false
	}
	if e.ProbationEndDate == nil {
		return true
	}
	return !employeeDateOnly(asOf).After(employeeDateOnly(*e.ProbationEndDate))
}

func employeeDateOnly(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, time.UTC)
}

func DefaultEmployeeCode(userID uuid.UUID) string {
	clean := strings.ToUpper(strings.ReplaceAll(userID.String(), "-", ""))
	if len(clean) > 8 {
		clean = clean[:8]
	}
	return "EMP-" + clean
}
