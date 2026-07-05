package domain

import (
	"encoding/json"
	"errors"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	CompCycleDraft     = "draft"
	CompCycleOpen      = "open"
	CompCycleSubmitted = "submitted"
	CompCycleApproved  = "approved"
	CompCycleFinalized = "finalized"
	CompCycleCancelled = "cancelled"

	CompRecommendationDraft           = "draft"
	CompRecommendationSubmitted       = "submitted"
	CompRecommendationApproved        = "approved"
	CompRecommendationRejected        = "rejected"
	CompRecommendationFinalized       = "finalized"
	CompRecommendationHandedToPayroll = "handed_to_payroll"

	CompEquityOpen         = "open"
	CompEquityAcknowledged = "acknowledged"
	CompEquityResolved     = "resolved"
	CompEquityWaived       = "waived"
)

var (
	ErrInvalidCompensationPayBand         = errors.New("invalid compensation pay band")
	ErrCompensationPayBandNotFound        = errors.New("compensation pay band not found")
	ErrInvalidCompensationCycle           = errors.New("invalid compensation cycle")
	ErrCompensationCycleNotFound          = errors.New("compensation cycle not found")
	ErrInvalidCompensationBudgetPool      = errors.New("invalid compensation budget pool")
	ErrCompensationBudgetPoolNotFound     = errors.New("compensation budget pool not found")
	ErrInvalidCompensationRecommendation  = errors.New("invalid compensation recommendation")
	ErrCompensationRecommendationNotFound = errors.New("compensation recommendation not found")
	ErrInvalidCompensationEquityCheck     = errors.New("invalid compensation equity check")
	ErrCompensationEquityCheckNotFound    = errors.New("compensation equity check not found")
)

type CompensationPayBand struct {
	ID            uuid.UUID       `json:"id"`
	TenantID      uuid.UUID       `json:"tenant_id"`
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
	Inactive      bool            `json:"inactive"`
	CreatedAt     time.Time       `json:"created_at"`
	CreatedBy     *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt     time.Time       `json:"updated_at"`
	UpdatedBy     *uuid.UUID      `json:"updated_by,omitempty"`
}

type CompensationCycle struct {
	ID               uuid.UUID       `json:"id"`
	TenantID         uuid.UUID       `json:"tenant_id"`
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
	FinalizedAt      *time.Time      `json:"finalized_at,omitempty"`
	Metadata         json.RawMessage `json:"metadata,omitempty"`
	Inactive         bool            `json:"inactive"`
	CreatedAt        time.Time       `json:"created_at"`
	CreatedBy        *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt        time.Time       `json:"updated_at"`
	UpdatedBy        *uuid.UUID      `json:"updated_by,omitempty"`
}

type CompensationBudgetPool struct {
	ID                  uuid.UUID  `json:"id"`
	TenantID            uuid.UUID  `json:"tenant_id"`
	CycleID             uuid.UUID  `json:"cycle_id"`
	Name                string     `json:"name"`
	PoolType            string     `json:"pool_type"`
	OwnerUserID         *uuid.UUID `json:"owner_user_id,omitempty"`
	DepartmentID        *uuid.UUID `json:"department_id,omitempty"`
	BranchID            *uuid.UUID `json:"branch_id,omitempty"`
	BudgetAmount        float64    `json:"budget_amount"`
	AllocatedAmount     float64    `json:"allocated_amount"`
	CommittedAmount     float64    `json:"committed_amount"`
	RecommendationCount int64      `json:"recommendation_count"`
	Notes               *string    `json:"notes,omitempty"`
	Inactive            bool       `json:"inactive"`
	CreatedAt           time.Time  `json:"created_at"`
	CreatedBy           *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt           time.Time  `json:"updated_at"`
	UpdatedBy           *uuid.UUID `json:"updated_by,omitempty"`
}

type CompensationRecommendation struct {
	ID                          uuid.UUID       `json:"id"`
	TenantID                    uuid.UUID       `json:"tenant_id"`
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
	PayrollHandoffAt            *time.Time      `json:"payroll_handoff_at,omitempty"`
	Metadata                    json.RawMessage `json:"metadata,omitempty"`
	Inactive                    bool            `json:"inactive"`
	CreatedAt                   time.Time       `json:"created_at"`
	CreatedBy                   *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt                   time.Time       `json:"updated_at"`
	UpdatedBy                   *uuid.UUID      `json:"updated_by,omitempty"`
	WorkerDisplayName           *string         `json:"worker_display_name,omitempty"`
	WorkerCode                  *string         `json:"worker_code,omitempty"`
	PayBandCode                 *string         `json:"pay_band_code,omitempty"`
	PayBandName                 *string         `json:"pay_band_name,omitempty"`
	BudgetPoolName              *string         `json:"budget_pool_name,omitempty"`
}

type CompensationEquityCheck struct {
	ID                uuid.UUID  `json:"id"`
	TenantID          uuid.UUID  `json:"tenant_id"`
	CycleID           uuid.UUID  `json:"cycle_id"`
	WorkerProfileID   uuid.UUID  `json:"worker_profile_id"`
	PayBandID         *uuid.UUID `json:"pay_band_id,omitempty"`
	CheckType         string     `json:"check_type"`
	Severity          string     `json:"severity"`
	CurrentSalary     float64    `json:"current_salary"`
	BandMin           float64    `json:"band_min"`
	BandMidpoint      float64    `json:"band_midpoint"`
	BandMax           float64    `json:"band_max"`
	VariancePercent   float64    `json:"variance_percent"`
	Finding           string     `json:"finding"`
	Recommendation    *string    `json:"recommendation,omitempty"`
	Status            string     `json:"status"`
	Inactive          bool       `json:"inactive"`
	CreatedAt         time.Time  `json:"created_at"`
	CreatedBy         *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt         time.Time  `json:"updated_at"`
	UpdatedBy         *uuid.UUID `json:"updated_by,omitempty"`
	WorkerDisplayName *string    `json:"worker_display_name,omitempty"`
	WorkerCode        *string    `json:"worker_code,omitempty"`
	PayBandCode       *string    `json:"pay_band_code,omitempty"`
	PayBandName       *string    `json:"pay_band_name,omitempty"`
}

type CompensationEvent struct {
	ID         uuid.UUID       `json:"id"`
	TenantID   uuid.UUID       `json:"tenant_id"`
	CycleID    *uuid.UUID      `json:"cycle_id,omitempty"`
	SourceType string          `json:"source_type"`
	SourceID   *uuid.UUID      `json:"source_id,omitempty"`
	Action     string          `json:"action"`
	FromStatus *string         `json:"from_status,omitempty"`
	ToStatus   *string         `json:"to_status,omitempty"`
	Remarks    *string         `json:"remarks,omitempty"`
	Metadata   json.RawMessage `json:"metadata,omitempty"`
	CreatedAt  time.Time       `json:"created_at"`
	CreatedBy  *uuid.UUID      `json:"created_by,omitempty"`
}

type CompensationSummaryRow struct {
	Metric       string  `json:"metric"`
	MetricCount  int64   `json:"metric_count"`
	MetricAmount float64 `json:"metric_amount"`
}

type CompensationFilter struct {
	TenantID uuid.UUID
	CycleID  *uuid.UUID
	Status   *string
	Search   *string
	Limit    int32
	Offset   int32
}

type CompensationPayBandFilter struct {
	TenantID     uuid.UUID
	IsActive     *bool
	CurrencyCode *string
	Search       *string
	Limit        int32
	Offset       int32
}

func NormalizeCompensationCycleStatus(value string) string {
	status := normalizeWorkerProfileEnum(value, "")
	for _, allowed := range []string{CompCycleDraft, CompCycleOpen, CompCycleSubmitted, CompCycleApproved, CompCycleFinalized, CompCycleCancelled} {
		if status == allowed {
			return status
		}
	}
	return ""
}

func NormalizeCompensationRecommendationStatus(value string) string {
	status := normalizeWorkerProfileEnum(value, "")
	for _, allowed := range []string{CompRecommendationDraft, CompRecommendationSubmitted, CompRecommendationApproved, CompRecommendationRejected, CompRecommendationFinalized, CompRecommendationHandedToPayroll} {
		if status == allowed {
			return status
		}
	}
	return ""
}

func NormalizeCompensationEquityStatus(value string) string {
	status := normalizeWorkerProfileEnum(value, "")
	for _, allowed := range []string{CompEquityOpen, CompEquityAcknowledged, CompEquityResolved, CompEquityWaived} {
		if status == allowed {
			return status
		}
	}
	return ""
}

func ValidateCompensationPayBand(item *CompensationPayBand) error {
	if item == nil || item.TenantID == uuid.Nil || strings.TrimSpace(item.Code) == "" || strings.TrimSpace(item.Name) == "" || strings.TrimSpace(item.CurrencyCode) == "" {
		return ErrInvalidCompensationPayBand
	}
	if invalidMoney(item.MinPay) || invalidMoney(item.MidpointPay) || invalidMoney(item.MaxPay) || item.MidpointPay < item.MinPay || item.MaxPay < item.MidpointPay {
		return ErrInvalidCompensationPayBand
	}
	item.Code = strings.TrimSpace(item.Code)
	item.Name = strings.TrimSpace(item.Name)
	item.CurrencyCode = strings.ToUpper(strings.TrimSpace(item.CurrencyCode))
	item.MinPay, item.MidpointPay, item.MaxPay = roundMoney(item.MinPay), roundMoney(item.MidpointPay), roundMoney(item.MaxPay)
	return nil
}

func ValidateCompensationCycle(item *CompensationCycle) error {
	if item == nil || item.TenantID == uuid.Nil || strings.TrimSpace(item.Code) == "" || strings.TrimSpace(item.Name) == "" || strings.TrimSpace(item.CurrencyCode) == "" || invalidMoney(item.BudgetAmount) {
		return ErrInvalidCompensationCycle
	}
	item.Status = NormalizeCompensationCycleStatus(item.Status)
	if item.Status == "" {
		item.Status = CompCycleDraft
	}
	if strings.TrimSpace(item.CycleType) == "" {
		item.CycleType = "annual"
	}
	item.Code = strings.TrimSpace(item.Code)
	item.Name = strings.TrimSpace(item.Name)
	item.CurrencyCode = strings.ToUpper(strings.TrimSpace(item.CurrencyCode))
	item.BudgetAmount = roundMoney(item.BudgetAmount)
	return nil
}

func ValidateCompensationBudgetPool(item *CompensationBudgetPool) error {
	if item == nil || item.TenantID == uuid.Nil || item.CycleID == uuid.Nil || strings.TrimSpace(item.Name) == "" || invalidMoney(item.BudgetAmount) || invalidMoney(item.AllocatedAmount) {
		return ErrInvalidCompensationBudgetPool
	}
	if strings.TrimSpace(item.PoolType) == "" {
		item.PoolType = "merit"
	}
	item.Name = strings.TrimSpace(item.Name)
	item.BudgetAmount = roundMoney(item.BudgetAmount)
	item.AllocatedAmount = roundMoney(item.AllocatedAmount)
	return nil
}

func ValidateCompensationRecommendation(item *CompensationRecommendation) error {
	if item == nil || item.TenantID == uuid.Nil || item.CycleID == uuid.Nil || item.WorkerProfileID == uuid.Nil || invalidMoney(item.CurrentSalary) || invalidMoney(item.RecommendedSalary) {
		return ErrInvalidCompensationRecommendation
	}
	item.Status = NormalizeCompensationRecommendationStatus(item.Status)
	if item.Status == "" {
		item.Status = CompRecommendationDraft
	}
	item.CurrentSalary = roundMoney(item.CurrentSalary)
	item.RecommendedSalary = roundMoney(item.RecommendedSalary)
	if item.RecommendedSalary >= item.CurrentSalary {
		item.RecommendedIncrementAmount = roundMoney(item.RecommendedSalary - item.CurrentSalary)
		if item.CurrentSalary > 0 {
			item.RecommendedIncrementPercent = math.Round(((item.RecommendedIncrementAmount/item.CurrentSalary)*100)*10000) / 10000
		}
	}
	if item.CurrentCompaRatio < 0 || math.IsNaN(item.CurrentCompaRatio) || math.IsInf(item.CurrentCompaRatio, 0) {
		return ErrInvalidCompensationRecommendation
	}
	return nil
}

func invalidMoney(value float64) bool {
	return value < 0 || math.IsNaN(value) || math.IsInf(value, 0)
}
