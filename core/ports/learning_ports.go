package ports

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type LearningRepo interface {
	CreateLearningCourse(ctx context.Context, item *domain.LearningCourse, actorID *uuid.UUID) (*domain.LearningCourse, error)
	UpdateLearningCourse(ctx context.Context, item *domain.LearningCourse, actorID *uuid.UUID) (*domain.LearningCourse, error)
	GetLearningCourse(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.LearningCourse, error)
	ListLearningCourses(ctx context.Context, filter domain.LearningCourseFilter) ([]*domain.LearningCourse, error)
	DeleteLearningCourse(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateLearningPath(ctx context.Context, item *domain.LearningPath, actorID *uuid.UUID) (*domain.LearningPath, error)
	UpdateLearningPath(ctx context.Context, item *domain.LearningPath, actorID *uuid.UUID) (*domain.LearningPath, error)
	GetLearningPath(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.LearningPath, error)
	ListLearningPaths(ctx context.Context, filter domain.LearningPathFilter) ([]*domain.LearningPath, error)
	DeleteLearningPath(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	UpsertLearningPathCourse(ctx context.Context, item *domain.LearningPathCourse, actorID *uuid.UUID) (*domain.LearningPathCourse, error)
	ListLearningPathCourses(ctx context.Context, tenantID uuid.UUID, pathID uuid.UUID) ([]*domain.LearningPathCourse, error)
	DeleteLearningPathCourse(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateLearningEnrollment(ctx context.Context, item *domain.LearningEnrollment, actorID *uuid.UUID) (*domain.LearningEnrollment, error)
	UpdateLearningEnrollmentStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, score *float64, certificateURL *string, certificateFileName *string, certificateContentType *string, notes *string, actorID *uuid.UUID) (*domain.LearningEnrollment, error)
	GetLearningEnrollment(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.LearningEnrollment, error)
	ListLearningEnrollments(ctx context.Context, filter domain.LearningEnrollmentFilter) ([]*domain.LearningEnrollment, error)
	DeleteLearningEnrollment(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateLearningRecommendation(ctx context.Context, item *domain.LearningRecommendation, actorID *uuid.UUID) (*domain.LearningRecommendation, error)
	UpdateLearningRecommendationStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, actorID *uuid.UUID) (*domain.LearningRecommendation, error)
	ListLearningRecommendations(ctx context.Context, filter domain.LearningRecommendationFilter) ([]*domain.LearningRecommendation, error)
	GenerateSkillGapLearningRecommendations(ctx context.Context, tenantID uuid.UUID, actorID *uuid.UUID) ([]*domain.LearningRecommendation, error)
	GetLearningSummary(ctx context.Context, tenantID uuid.UUID) ([]*domain.LearningSummaryRow, error)
}

type LearningCertificateStorage interface {
	StoreLearningCertificate(ctx context.Context, input StoreLearningCertificateInput) (string, error)
}

type StoreLearningCertificateInput struct {
	TenantID        uuid.UUID
	EnrollmentID    uuid.UUID
	WorkerProfileID uuid.UUID
	FileName        string
	ContentType     string
	Content         []byte
}

type LearningCourseCommand struct {
	TenantID            uuid.UUID       `json:"tenant_id"`
	ID                  uuid.UUID       `json:"id,omitempty"`
	Code                string          `json:"code"`
	Title               string          `json:"title"`
	Description         *string         `json:"description,omitempty"`
	CourseType          string          `json:"course_type"`
	DeliveryMode        string          `json:"delivery_mode"`
	Provider            *string         `json:"provider,omitempty"`
	DurationMinutes     int32           `json:"duration_minutes"`
	SkillID             *uuid.UUID      `json:"skill_id,omitempty"`
	ComplianceRuleID    *uuid.UUID      `json:"compliance_rule_id,omitempty"`
	Mandatory           bool            `json:"mandatory"`
	AIReadiness         bool            `json:"ai_readiness"`
	CertificateRequired bool            `json:"certificate_required"`
	BudgetAmount        *float64        `json:"budget_amount,omitempty"`
	CurrencyCode        string          `json:"currency_code"`
	IsActive            bool            `json:"is_active"`
	Metadata            json.RawMessage `json:"metadata,omitempty"`
	ActorID             *uuid.UUID      `json:"-"`
}

type LearningPathCommand struct {
	TenantID    uuid.UUID       `json:"tenant_id"`
	ID          uuid.UUID       `json:"id,omitempty"`
	Code        string          `json:"code"`
	Title       string          `json:"title"`
	Description *string         `json:"description,omitempty"`
	PathType    string          `json:"path_type"`
	TargetRole  *string         `json:"target_role,omitempty"`
	SkillID     *uuid.UUID      `json:"skill_id,omitempty"`
	AIReadiness bool            `json:"ai_readiness"`
	IsActive    bool            `json:"is_active"`
	Metadata    json.RawMessage `json:"metadata,omitempty"`
	ActorID     *uuid.UUID      `json:"-"`
}

type LearningPathCourseCommand struct {
	TenantID  uuid.UUID  `json:"tenant_id"`
	ID        uuid.UUID  `json:"id,omitempty"`
	PathID    uuid.UUID  `json:"path_id"`
	CourseID  uuid.UUID  `json:"course_id"`
	SortOrder int32      `json:"sort_order"`
	Required  bool       `json:"required"`
	ActorID   *uuid.UUID `json:"-"`
}

type LearningEnrollmentCommand struct {
	TenantID         uuid.UUID       `json:"tenant_id"`
	CourseID         uuid.UUID       `json:"course_id"`
	PathID           *uuid.UUID      `json:"path_id,omitempty"`
	WorkerProfileID  uuid.UUID       `json:"worker_profile_id"`
	AssignmentSource string          `json:"assignment_source"`
	Status           string          `json:"status"`
	NominatedBy      *uuid.UUID      `json:"nominated_by,omitempty"`
	AssignedBy       *uuid.UUID      `json:"assigned_by,omitempty"`
	DueDate          *time.Time      `json:"due_date,omitempty"`
	Notes            *string         `json:"notes,omitempty"`
	Metadata         json.RawMessage `json:"metadata,omitempty"`
	ActorID          *uuid.UUID      `json:"-"`
}

type LearningEnrollmentStatusCommand struct {
	TenantID               uuid.UUID  `json:"tenant_id"`
	ID                     uuid.UUID  `json:"id"`
	Status                 string     `json:"status"`
	Score                  *float64   `json:"score,omitempty"`
	CertificateURL         *string    `json:"certificate_url,omitempty"`
	CertificateFileName    *string    `json:"certificate_file_name,omitempty"`
	CertificateContentType *string    `json:"certificate_content_type,omitempty"`
	Notes                  *string    `json:"notes,omitempty"`
	ActorID                *uuid.UUID `json:"-"`
}

type LearningCertificateCommand struct {
	TenantID      uuid.UUID  `json:"tenant_id"`
	EnrollmentID  uuid.UUID  `json:"enrollment_id"`
	FileName      string     `json:"file_name"`
	ContentType   string     `json:"content_type"`
	ContentBase64 string     `json:"content_base64"`
	Score         *float64   `json:"score,omitempty"`
	Notes         *string    `json:"notes,omitempty"`
	ActorID       *uuid.UUID `json:"-"`
}

type LearningRecommendationCommand struct {
	TenantID        uuid.UUID       `json:"tenant_id"`
	WorkerProfileID *uuid.UUID      `json:"worker_profile_id,omitempty"`
	SkillID         *uuid.UUID      `json:"skill_id,omitempty"`
	CourseID        *uuid.UUID      `json:"course_id,omitempty"`
	PathID          *uuid.UUID      `json:"path_id,omitempty"`
	SourceType      string          `json:"source_type"`
	Reason          string          `json:"reason"`
	Priority        string          `json:"priority"`
	ConfidenceScore *float64        `json:"confidence_score,omitempty"`
	Status          string          `json:"status"`
	Metadata        json.RawMessage `json:"metadata,omitempty"`
	ActorID         *uuid.UUID      `json:"-"`
}

type LearningRecommendationStatusCommand struct {
	TenantID uuid.UUID  `json:"tenant_id"`
	ID       uuid.UUID  `json:"id"`
	Status   string     `json:"status"`
	ActorID  *uuid.UUID `json:"-"`
}
