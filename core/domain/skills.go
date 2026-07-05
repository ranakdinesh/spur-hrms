package domain

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	SkillSourceGlobal = "global"
	SkillSourceTenant = "tenant"

	SkillTypeTechnical  = "technical"
	SkillTypeFunctional = "functional"
	SkillTypeBehavioral = "behavioral"
	SkillTypeCompliance = "compliance"
	SkillTypeTool       = "tool"
	SkillTypeLanguage   = "language"
	SkillTypeDomain     = "domain"
	SkillTypeCustom     = "custom"

	SkillProficiencyBeginner     = "beginner"
	SkillProficiencyIntermediate = "intermediate"
	SkillProficiencyAdvanced     = "advanced"
	SkillProficiencyExpert       = "expert"

	SkillVerificationSelfDeclared    = "self_declared"
	SkillVerificationManagerEndorsed = "manager_endorsed"
	SkillVerificationHRVerified      = "hr_verified"
	SkillVerificationExpired         = "expired"
	SkillVerificationRejected        = "rejected"

	SkillAssessmentSelf     = "self"
	SkillAssessmentManager  = "manager"
	SkillAssessmentHR       = "hr"
	SkillAssessmentExternal = "external"

	SkillAssessmentSubmitted = "submitted"
	SkillAssessmentObserved  = "observed"
	SkillAssessmentPassed    = "passed"
	SkillAssessmentFailed    = "failed"
)

var (
	ErrInvalidSkillCategory   = errors.New("skill category is invalid")
	ErrSkillCategoryNotFound  = errors.New("skill category not found")
	ErrInvalidSkill           = errors.New("skill is invalid")
	ErrSkillNotFound          = errors.New("skill not found")
	ErrInvalidWorkerSkill     = errors.New("worker skill is invalid")
	ErrWorkerSkillNotFound    = errors.New("worker skill not found")
	ErrInvalidSkillAssessment = errors.New("worker skill assessment is invalid")
)

type SkillCategory struct {
	ID          uuid.UUID       `json:"id"`
	TenantID    *uuid.UUID      `json:"tenant_id,omitempty"`
	ParentID    *uuid.UUID      `json:"parent_id,omitempty"`
	Code        string          `json:"code"`
	Name        string          `json:"name"`
	Description *string         `json:"description,omitempty"`
	SourceScope string          `json:"source_scope"`
	SortOrder   int32           `json:"sort_order"`
	Metadata    json.RawMessage `json:"metadata,omitempty"`
	Inactive    bool            `json:"inactive"`
	CreatedAt   time.Time       `json:"created_at"`
	CreatedBy   *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt   time.Time       `json:"updated_at"`
	UpdatedBy   *uuid.UUID      `json:"updated_by,omitempty"`
	ParentName  *string         `json:"parent_name,omitempty"`
}

type Skill struct {
	ID                  uuid.UUID       `json:"id"`
	TenantID            *uuid.UUID      `json:"tenant_id,omitempty"`
	CategoryID          *uuid.UUID      `json:"category_id,omitempty"`
	Code                string          `json:"code"`
	Name                string          `json:"name"`
	Description         *string         `json:"description,omitempty"`
	SkillType           string          `json:"skill_type"`
	SourceScope         string          `json:"source_scope"`
	CertificateRequired bool            `json:"certificate_required"`
	AssessmentRequired  bool            `json:"assessment_required"`
	IsActive            bool            `json:"is_active"`
	Metadata            json.RawMessage `json:"metadata,omitempty"`
	Inactive            bool            `json:"inactive"`
	CreatedAt           time.Time       `json:"created_at"`
	CreatedBy           *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt           time.Time       `json:"updated_at"`
	UpdatedBy           *uuid.UUID      `json:"updated_by,omitempty"`
	CategoryName        *string         `json:"category_name,omitempty"`
	CategoryCode        *string         `json:"category_code,omitempty"`
}

type WorkerSkill struct {
	ID                   uuid.UUID       `json:"id"`
	TenantID             uuid.UUID       `json:"tenant_id"`
	WorkerProfileID      uuid.UUID       `json:"worker_profile_id"`
	SkillID              uuid.UUID       `json:"skill_id"`
	SkillNameSnapshot    string          `json:"skill_name_snapshot"`
	Proficiency          string          `json:"proficiency"`
	YearsExperience      *float64        `json:"years_experience,omitempty"`
	LastUsedOn           *time.Time      `json:"last_used_on,omitempty"`
	VerificationStatus   string          `json:"verification_status"`
	CertificateURL       *string         `json:"certificate_url,omitempty"`
	CertificateExpiresOn *time.Time      `json:"certificate_expires_on,omitempty"`
	AssessmentScore      *float64        `json:"assessment_score,omitempty"`
	AssessedOn           *time.Time      `json:"assessed_on,omitempty"`
	EndorsedBy           *uuid.UUID      `json:"endorsed_by,omitempty"`
	EndorsedAt           *time.Time      `json:"endorsed_at,omitempty"`
	VerifiedBy           *uuid.UUID      `json:"verified_by,omitempty"`
	VerifiedAt           *time.Time      `json:"verified_at,omitempty"`
	Notes                *string         `json:"notes,omitempty"`
	Metadata             json.RawMessage `json:"metadata,omitempty"`
	Inactive             bool            `json:"inactive"`
	CreatedAt            time.Time       `json:"created_at"`
	CreatedBy            *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt            time.Time       `json:"updated_at"`
	UpdatedBy            *uuid.UUID      `json:"updated_by,omitempty"`
	WorkerDisplayName    *string         `json:"worker_display_name,omitempty"`
	WorkerCode           *string         `json:"worker_code,omitempty"`
	SkillCode            *string         `json:"skill_code,omitempty"`
	SkillName            *string         `json:"skill_name,omitempty"`
	SkillType            *string         `json:"skill_type,omitempty"`
	SkillSourceScope     *string         `json:"skill_source_scope,omitempty"`
	CertificateRequired  bool            `json:"certificate_required"`
	AssessmentRequired   bool            `json:"assessment_required"`
	CategoryName         *string         `json:"category_name,omitempty"`
	CategoryCode         *string         `json:"category_code,omitempty"`
}

type WorkerSkillAssessment struct {
	ID             uuid.UUID       `json:"id"`
	TenantID       uuid.UUID       `json:"tenant_id"`
	WorkerSkillID  uuid.UUID       `json:"worker_skill_id"`
	AssessmentType string          `json:"assessment_type"`
	ResultStatus   string          `json:"result_status"`
	Score          *float64        `json:"score,omitempty"`
	MaxScore       *float64        `json:"max_score,omitempty"`
	AssessedBy     *uuid.UUID      `json:"assessed_by,omitempty"`
	AssessedOn     time.Time       `json:"assessed_on"`
	EvidenceURL    *string         `json:"evidence_url,omitempty"`
	Notes          *string         `json:"notes,omitempty"`
	Metadata       json.RawMessage `json:"metadata,omitempty"`
	CreatedAt      time.Time       `json:"created_at"`
	CreatedBy      *uuid.UUID      `json:"created_by,omitempty"`
}

type SkillCategoryInput struct {
	TenantID    uuid.UUID
	ParentID    *uuid.UUID
	Code        string
	Name        string
	Description *string
	SortOrder   int32
	Metadata    json.RawMessage
}

type SkillInput struct {
	TenantID            uuid.UUID
	CategoryID          *uuid.UUID
	Code                string
	Name                string
	Description         *string
	SkillType           string
	CertificateRequired bool
	AssessmentRequired  bool
	IsActive            bool
	Metadata            json.RawMessage
}

type WorkerSkillInput struct {
	TenantID             uuid.UUID
	WorkerProfileID      uuid.UUID
	SkillID              uuid.UUID
	SkillNameSnapshot    string
	Proficiency          string
	YearsExperience      *float64
	LastUsedOn           *time.Time
	VerificationStatus   string
	CertificateURL       *string
	CertificateExpiresOn *time.Time
	AssessmentScore      *float64
	AssessedOn           *time.Time
	Notes                *string
	Metadata             json.RawMessage
}

type WorkerSkillAssessmentInput struct {
	TenantID       uuid.UUID
	WorkerSkillID  uuid.UUID
	AssessmentType string
	ResultStatus   string
	Score          *float64
	MaxScore       *float64
	AssessedBy     *uuid.UUID
	AssessedOn     *time.Time
	EvidenceURL    *string
	Notes          *string
	Metadata       json.RawMessage
}

type SkillCategoryFilter struct {
	TenantID    uuid.UUID
	SourceScope *string
	ParentID    *uuid.UUID
	Search      *string
}

type SkillFilter struct {
	TenantID    uuid.UUID
	CategoryID  *uuid.UUID
	SkillType   *string
	SourceScope *string
	IsActive    *bool
	Search      *string
}

type WorkerSkillFilter struct {
	TenantID                  uuid.UUID
	WorkerProfileID           *uuid.UUID
	SkillID                   *uuid.UUID
	CategoryID                *uuid.UUID
	Proficiency               *string
	VerificationStatus        *string
	CertificateExpiringBefore *time.Time
	Search                    *string
}

type SkillsSummaryRow struct {
	Status                   string `json:"status"`
	WorkerSkillCount         int32  `json:"worker_skill_count"`
	WorkerCount              int32  `json:"worker_count"`
	SkillCount               int32  `json:"skill_count"`
	ExpiringCertificateCount int32  `json:"expiring_certificate_count"`
}

func NewSkillCategory(input SkillCategoryInput) (*SkillCategory, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	code := normalizeSkillCode(input.Code)
	name := strings.TrimSpace(input.Name)
	if code == "" || name == "" {
		return nil, ErrInvalidSkillCategory
	}
	metadata, err := validSkillMetadata(input.Metadata)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	return &SkillCategory{TenantID: &input.TenantID, ParentID: cleanUUIDOptional(input.ParentID), Code: code, Name: name, Description: cleanOptional(input.Description), SourceScope: SkillSourceTenant, SortOrder: input.SortOrder, Metadata: metadata, CreatedAt: now, UpdatedAt: now}, nil
}

func NewSkill(input SkillInput) (*Skill, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	code := normalizeSkillCode(input.Code)
	name := strings.TrimSpace(input.Name)
	skillType := normalizeWorkerProfileEnum(input.SkillType, SkillTypeTechnical)
	if code == "" || name == "" || !containsString([]string{SkillTypeTechnical, SkillTypeFunctional, SkillTypeBehavioral, SkillTypeCompliance, SkillTypeTool, SkillTypeLanguage, SkillTypeDomain, SkillTypeCustom}, skillType) {
		return nil, ErrInvalidSkill
	}
	metadata, err := validSkillMetadata(input.Metadata)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	return &Skill{TenantID: &input.TenantID, CategoryID: cleanUUIDOptional(input.CategoryID), Code: code, Name: name, Description: cleanOptional(input.Description), SkillType: skillType, SourceScope: SkillSourceTenant, CertificateRequired: input.CertificateRequired, AssessmentRequired: input.AssessmentRequired, IsActive: input.IsActive, Metadata: metadata, CreatedAt: now, UpdatedAt: now}, nil
}

func NewWorkerSkill(input WorkerSkillInput) (*WorkerSkill, error) {
	if input.TenantID == uuid.Nil || input.WorkerProfileID == uuid.Nil || input.SkillID == uuid.Nil {
		return nil, ErrInvalidWorkerSkill
	}
	name := strings.TrimSpace(input.SkillNameSnapshot)
	proficiency := normalizeWorkerProfileEnum(input.Proficiency, SkillProficiencyBeginner)
	status := normalizeWorkerProfileEnum(input.VerificationStatus, SkillVerificationSelfDeclared)
	if name == "" || !containsString([]string{SkillProficiencyBeginner, SkillProficiencyIntermediate, SkillProficiencyAdvanced, SkillProficiencyExpert}, proficiency) || !containsString([]string{SkillVerificationSelfDeclared, SkillVerificationManagerEndorsed, SkillVerificationHRVerified, SkillVerificationExpired, SkillVerificationRejected}, status) {
		return nil, ErrInvalidWorkerSkill
	}
	if input.YearsExperience != nil && *input.YearsExperience < 0 {
		return nil, ErrInvalidWorkerSkill
	}
	if input.AssessmentScore != nil && (*input.AssessmentScore < 0 || *input.AssessmentScore > 100) {
		return nil, ErrInvalidWorkerSkill
	}
	metadata, err := validSkillMetadata(input.Metadata)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	return &WorkerSkill{TenantID: input.TenantID, WorkerProfileID: input.WorkerProfileID, SkillID: input.SkillID, SkillNameSnapshot: name, Proficiency: proficiency, YearsExperience: input.YearsExperience, LastUsedOn: datePtrUTC(input.LastUsedOn), VerificationStatus: status, CertificateURL: cleanOptional(input.CertificateURL), CertificateExpiresOn: datePtrUTC(input.CertificateExpiresOn), AssessmentScore: input.AssessmentScore, AssessedOn: datePtrUTC(input.AssessedOn), Notes: cleanOptional(input.Notes), Metadata: metadata, CreatedAt: now, UpdatedAt: now}, nil
}

func NewWorkerSkillAssessment(input WorkerSkillAssessmentInput) (*WorkerSkillAssessment, error) {
	if input.TenantID == uuid.Nil || input.WorkerSkillID == uuid.Nil {
		return nil, ErrInvalidSkillAssessment
	}
	assessmentType := normalizeWorkerProfileEnum(input.AssessmentType, SkillAssessmentManager)
	resultStatus := normalizeWorkerProfileEnum(input.ResultStatus, SkillAssessmentSubmitted)
	if !containsString([]string{SkillAssessmentSelf, SkillAssessmentManager, SkillAssessmentHR, SkillAssessmentExternal}, assessmentType) || !containsString([]string{SkillAssessmentSubmitted, SkillAssessmentObserved, SkillAssessmentPassed, SkillAssessmentFailed}, resultStatus) {
		return nil, ErrInvalidSkillAssessment
	}
	if input.Score != nil && *input.Score < 0 {
		return nil, ErrInvalidSkillAssessment
	}
	if input.MaxScore != nil && *input.MaxScore <= 0 {
		return nil, ErrInvalidSkillAssessment
	}
	metadata, err := validSkillMetadata(input.Metadata)
	if err != nil {
		return nil, err
	}
	assessedOn := time.Now().UTC()
	if input.AssessedOn != nil && !input.AssessedOn.IsZero() {
		assessedOn = *datePtrUTC(input.AssessedOn)
	}
	return &WorkerSkillAssessment{TenantID: input.TenantID, WorkerSkillID: input.WorkerSkillID, AssessmentType: assessmentType, ResultStatus: resultStatus, Score: input.Score, MaxScore: input.MaxScore, AssessedBy: cleanUUIDOptional(input.AssessedBy), AssessedOn: assessedOn, EvidenceURL: cleanOptional(input.EvidenceURL), Notes: cleanOptional(input.Notes), Metadata: metadata, CreatedAt: time.Now().UTC(), CreatedBy: cleanUUIDOptional(input.AssessedBy)}, nil
}

func normalizeSkillCode(value string) string {
	clean := strings.ToLower(strings.TrimSpace(value))
	clean = strings.ReplaceAll(clean, " ", "-")
	clean = strings.ReplaceAll(clean, "_", "-")
	return clean
}

func NormalizeSkillVerificationStatus(value string) string {
	status := normalizeWorkerProfileEnum(value, "")
	if !containsString([]string{SkillVerificationSelfDeclared, SkillVerificationManagerEndorsed, SkillVerificationHRVerified, SkillVerificationExpired, SkillVerificationRejected}, status) {
		return ""
	}
	return status
}

func validSkillMetadata(value json.RawMessage) (json.RawMessage, error) {
	metadata := normalizeWorkerJSONObject(value, "{}")
	if !json.Valid(metadata) || !jsonObject(metadata) {
		return nil, ErrInvalidSkill
	}
	return metadata, nil
}
