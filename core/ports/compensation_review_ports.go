package ports

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type CompensationReviewRepo interface {
	CreateCompensationPayBand(ctx context.Context, item *domain.CompensationPayBand, actorID *uuid.UUID) (*domain.CompensationPayBand, error)
	UpdateCompensationPayBand(ctx context.Context, item *domain.CompensationPayBand, actorID *uuid.UUID) (*domain.CompensationPayBand, error)
	GetCompensationPayBand(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.CompensationPayBand, error)
	ListCompensationPayBands(ctx context.Context, filter domain.CompensationPayBandFilter) ([]*domain.CompensationPayBand, error)
	DeleteCompensationPayBand(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateCompensationCycle(ctx context.Context, item *domain.CompensationCycle, actorID *uuid.UUID) (*domain.CompensationCycle, error)
	UpdateCompensationCycle(ctx context.Context, item *domain.CompensationCycle, actorID *uuid.UUID) (*domain.CompensationCycle, error)
	UpdateCompensationCycleStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, actorID *uuid.UUID) (*domain.CompensationCycle, error)
	GetCompensationCycle(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.CompensationCycle, error)
	ListCompensationCycles(ctx context.Context, filter domain.CompensationFilter) ([]*domain.CompensationCycle, error)
	CreateCompensationBudgetPool(ctx context.Context, item *domain.CompensationBudgetPool, actorID *uuid.UUID) (*domain.CompensationBudgetPool, error)
	UpdateCompensationBudgetPool(ctx context.Context, item *domain.CompensationBudgetPool, actorID *uuid.UUID) (*domain.CompensationBudgetPool, error)
	ListCompensationBudgetPools(ctx context.Context, tenantID uuid.UUID, cycleID uuid.UUID) ([]*domain.CompensationBudgetPool, error)
	DeleteCompensationBudgetPool(ctx context.Context, tenantID uuid.UUID, cycleID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateCompensationRecommendation(ctx context.Context, item *domain.CompensationRecommendation, actorID *uuid.UUID) (*domain.CompensationRecommendation, error)
	UpdateCompensationRecommendation(ctx context.Context, item *domain.CompensationRecommendation, actorID *uuid.UUID) (*domain.CompensationRecommendation, error)
	UpdateCompensationRecommendationStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, actorID *uuid.UUID) (*domain.CompensationRecommendation, error)
	GetCompensationRecommendation(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.CompensationRecommendation, error)
	ListCompensationRecommendations(ctx context.Context, filter domain.CompensationFilter) ([]*domain.CompensationRecommendation, error)
	DeleteCompensationRecommendation(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	GenerateCompensationEquityChecks(ctx context.Context, tenantID uuid.UUID, cycleID uuid.UUID, actorID *uuid.UUID) ([]*domain.CompensationEquityCheck, error)
	ListCompensationEquityChecks(ctx context.Context, filter domain.CompensationFilter) ([]*domain.CompensationEquityCheck, error)
	UpdateCompensationEquityCheckStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, actorID *uuid.UUID) (*domain.CompensationEquityCheck, error)
	CreateCompensationEvent(ctx context.Context, item *domain.CompensationEvent, actorID *uuid.UUID) (*domain.CompensationEvent, error)
	ListCompensationEvents(ctx context.Context, filter domain.CompensationFilter, sourceType *string, sourceID *uuid.UUID) ([]*domain.CompensationEvent, error)
	GetCompensationSummary(ctx context.Context, tenantID uuid.UUID) ([]*domain.CompensationSummaryRow, error)
}

type CompensationPayBandCommand struct {
	TenantID      uuid.UUID       `json:"tenant_id"`
	ID            uuid.UUID       `json:"id,omitempty"`
	Code          string          `json:"code"`
	Name          string          `json:"name"`
	JobFamily     *string         `json:"job_family,omitempty"`
	LevelCode     *string         `json:"level_code,omitempty"`
	LocationLabel *string         `json:"location_label,omitempty"`
	CurrencyCode  string          `json:"currency_code"`
	MinPay        float64         `json:"min_pay"`
	MidpointPay   float64         `json:"midpoint_pay"`
	MaxPay        float64         `json:"max_pay"`
	EffectiveFrom *time.Time      `json:"effective_from,omitempty"`
	EffectiveTo   *time.Time      `json:"effective_to,omitempty"`
	IsActive      bool            `json:"is_active"`
	Notes         *string         `json:"notes,omitempty"`
	Metadata      json.RawMessage `json:"metadata,omitempty"`
	ActorID       *uuid.UUID      `json:"-"`
}

type CompensationCycleCommand struct {
	TenantID         uuid.UUID       `json:"tenant_id"`
	ID               uuid.UUID       `json:"id,omitempty"`
	Code             string          `json:"code"`
	Name             string          `json:"name"`
	FiscalYearID     *uuid.UUID      `json:"fiscal_year_id,omitempty"`
	Status           string          `json:"status"`
	CycleType        string          `json:"cycle_type"`
	StartsOn         *time.Time      `json:"starts_on,omitempty"`
	EndsOn           *time.Time      `json:"ends_on,omitempty"`
	EffectiveDate    *time.Time      `json:"effective_date,omitempty"`
	CurrencyCode     string          `json:"currency_code"`
	BudgetAmount     float64         `json:"budget_amount"`
	PlanningGuidance *string         `json:"planning_guidance,omitempty"`
	ApprovalPolicy   *string         `json:"approval_policy,omitempty"`
	Metadata         json.RawMessage `json:"metadata,omitempty"`
	ActorID          *uuid.UUID      `json:"-"`
}

type CompensationStatusCommand struct {
	TenantID uuid.UUID  `json:"tenant_id"`
	ID       uuid.UUID  `json:"id"`
	Status   string     `json:"status"`
	Remarks  *string    `json:"remarks,omitempty"`
	ActorID  *uuid.UUID `json:"-"`
}

type CompensationBudgetPoolCommand struct {
	TenantID        uuid.UUID  `json:"tenant_id"`
	ID              uuid.UUID  `json:"id,omitempty"`
	CycleID         uuid.UUID  `json:"cycle_id"`
	Name            string     `json:"name"`
	PoolType        string     `json:"pool_type"`
	OwnerUserID     *uuid.UUID `json:"owner_user_id,omitempty"`
	DepartmentID    *uuid.UUID `json:"department_id,omitempty"`
	BranchID        *uuid.UUID `json:"branch_id,omitempty"`
	BudgetAmount    float64    `json:"budget_amount"`
	AllocatedAmount float64    `json:"allocated_amount"`
	Notes           *string    `json:"notes,omitempty"`
	ActorID         *uuid.UUID `json:"-"`
}

type CompensationRecommendationCommand struct {
	TenantID                    uuid.UUID       `json:"tenant_id"`
	ID                          uuid.UUID       `json:"id,omitempty"`
	CycleID                     uuid.UUID       `json:"cycle_id"`
	WorkerProfileID             uuid.UUID       `json:"worker_profile_id"`
	PayBandID                   *uuid.UUID      `json:"pay_band_id,omitempty"`
	BudgetPoolID                *uuid.UUID      `json:"budget_pool_id,omitempty"`
	CurrentSalary               float64         `json:"current_salary"`
	CurrentCompaRatio           float64         `json:"current_compa_ratio"`
	RecommendedSalary           float64         `json:"recommended_salary"`
	RecommendedIncrementAmount  float64         `json:"recommended_increment_amount"`
	RecommendedIncrementPercent float64         `json:"recommended_increment_percent"`
	PromotionRecommended        bool            `json:"promotion_recommended"`
	RecommendedDesignationID    *uuid.UUID      `json:"recommended_designation_id,omitempty"`
	Reason                      *string         `json:"reason,omitempty"`
	PerformanceRating           *string         `json:"performance_rating,omitempty"`
	EquityFlag                  bool            `json:"equity_flag"`
	EquityNotes                 *string         `json:"equity_notes,omitempty"`
	Status                      string          `json:"status"`
	EffectiveDate               *time.Time      `json:"effective_date,omitempty"`
	Metadata                    json.RawMessage `json:"metadata,omitempty"`
	ActorID                     *uuid.UUID      `json:"-"`
}
