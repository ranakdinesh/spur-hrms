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
	OKRCycleStatusDraft    = "draft"
	OKRCycleStatusActive   = "active"
	OKRCycleStatusClosed   = "closed"
	OKRCycleStatusArchived = "archived"

	OKRReviewCadenceWeekly   = "weekly"
	OKRReviewCadenceBiweekly = "biweekly"
	OKRReviewCadenceMonthly  = "monthly"

	ObjectiveOwnerCompany    = "company"
	ObjectiveOwnerDepartment = "department"
	ObjectiveOwnerProject    = "project"
	ObjectiveOwnerWorker     = "worker"

	ObjectiveStatusDraft     = "draft"
	ObjectiveStatusActive    = "active"
	ObjectiveStatusAtRisk    = "at_risk"
	ObjectiveStatusCompleted = "completed"
	ObjectiveStatusClosed    = "closed"
	ObjectiveStatusCancelled = "cancelled"

	ObjectivePriorityLow      = "low"
	ObjectivePriorityNormal   = "normal"
	ObjectivePriorityHigh     = "high"
	ObjectivePriorityCritical = "critical"

	KeyResultMetricNumber   = "number"
	KeyResultMetricPercent  = "percent"
	KeyResultMetricCurrency = "currency"
	KeyResultMetricBoolean  = "boolean"

	KeyResultStatusNotStarted = "not_started"
	KeyResultStatusOnTrack    = "on_track"
	KeyResultStatusAtRisk     = "at_risk"
	KeyResultStatusBehind     = "behind"
	KeyResultStatusCompleted  = "completed"
	KeyResultStatusClosed     = "closed"
	KeyResultStatusCancelled  = "cancelled"

	OKRConfidenceLow    = "low"
	OKRConfidenceMedium = "medium"
	OKRConfidenceHigh   = "high"
)

var (
	ErrInvalidOKRCycle         = errors.New("okr cycle is invalid")
	ErrOKRCycleNotFound        = errors.New("okr cycle not found")
	ErrInvalidObjective        = errors.New("objective is invalid")
	ErrObjectiveNotFound       = errors.New("objective not found")
	ErrInvalidKeyResult        = errors.New("key result is invalid")
	ErrKeyResultNotFound       = errors.New("key result not found")
	ErrInvalidKeyResultCheckIn = errors.New("key result check-in is invalid")
)

type OKRCycle struct {
	ID            uuid.UUID       `json:"id"`
	TenantID      uuid.UUID       `json:"tenant_id"`
	Name          string          `json:"name"`
	CycleCode     string          `json:"cycle_code"`
	Description   *string         `json:"description,omitempty"`
	StartDate     time.Time       `json:"start_date"`
	EndDate       time.Time       `json:"end_date"`
	Status        string          `json:"status"`
	ReviewCadence string          `json:"review_cadence"`
	Metadata      json.RawMessage `json:"metadata,omitempty"`
	Inactive      bool            `json:"inactive"`
	CreatedAt     time.Time       `json:"created_at"`
	CreatedBy     *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt     time.Time       `json:"updated_at"`
	UpdatedBy     *uuid.UUID      `json:"updated_by,omitempty"`
}

type OKRCycleInput struct {
	TenantID      uuid.UUID
	Name          string
	CycleCode     string
	Description   *string
	StartDate     *time.Time
	EndDate       *time.Time
	Status        string
	ReviewCadence string
	Metadata      json.RawMessage
}

type OKRCycleFilter struct {
	TenantID uuid.UUID
	Status   *string
	Search   *string
}

type Objective struct {
	ID                       uuid.UUID       `json:"id"`
	TenantID                 uuid.UUID       `json:"tenant_id"`
	CycleID                  uuid.UUID       `json:"cycle_id"`
	ParentObjectiveID        *uuid.UUID      `json:"parent_objective_id,omitempty"`
	OwnerType                string          `json:"owner_type"`
	OwnerWorkerProfileID     *uuid.UUID      `json:"owner_worker_profile_id,omitempty"`
	OwnerDepartmentID        *uuid.UUID      `json:"owner_department_id,omitempty"`
	OwnerProjectID           *uuid.UUID      `json:"owner_project_id,omitempty"`
	Title                    string          `json:"title"`
	Description              *string         `json:"description,omitempty"`
	Status                   string          `json:"status"`
	Priority                 string          `json:"priority"`
	ProgressPercent          float64         `json:"progress_percent"`
	Weight                   float64         `json:"weight"`
	StartDate                *time.Time      `json:"start_date,omitempty"`
	DueDate                  *time.Time      `json:"due_date,omitempty"`
	Metadata                 json.RawMessage `json:"metadata,omitempty"`
	Inactive                 bool            `json:"inactive"`
	CreatedAt                time.Time       `json:"created_at"`
	CreatedBy                *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt                time.Time       `json:"updated_at"`
	UpdatedBy                *uuid.UUID      `json:"updated_by,omitempty"`
	CycleName                *string         `json:"cycle_name,omitempty"`
	ParentObjectiveTitle     *string         `json:"parent_objective_title,omitempty"`
	OwnerWorkerName          *string         `json:"owner_worker_name,omitempty"`
	OwnerWorkerCode          *string         `json:"owner_worker_code,omitempty"`
	OwnerDepartmentName      *string         `json:"owner_department_name,omitempty"`
	OwnerProjectName         *string         `json:"owner_project_name,omitempty"`
	OwnerProjectCode         *string         `json:"owner_project_code,omitempty"`
	KeyResultCount           int32           `json:"key_result_count"`
	AverageKeyResultProgress float64         `json:"average_key_result_progress"`
}

type ObjectiveInput struct {
	TenantID             uuid.UUID
	CycleID              uuid.UUID
	ParentObjectiveID    *uuid.UUID
	OwnerType            string
	OwnerWorkerProfileID *uuid.UUID
	OwnerDepartmentID    *uuid.UUID
	OwnerProjectID       *uuid.UUID
	Title                string
	Description          *string
	Status               string
	Priority             string
	ProgressPercent      *float64
	Weight               *float64
	StartDate            *time.Time
	DueDate              *time.Time
	Metadata             json.RawMessage
}

type ObjectiveFilter struct {
	TenantID  uuid.UUID
	CycleID   *uuid.UUID
	OwnerType *string
	Status    *string
	Search    *string
}

type KeyResult struct {
	ID                uuid.UUID       `json:"id"`
	TenantID          uuid.UUID       `json:"tenant_id"`
	ObjectiveID       uuid.UUID       `json:"objective_id"`
	Title             string          `json:"title"`
	Description       *string         `json:"description,omitempty"`
	MetricType        string          `json:"metric_type"`
	StartValue        float64         `json:"start_value"`
	TargetValue       float64         `json:"target_value"`
	CurrentValue      float64         `json:"current_value"`
	ProgressPercent   float64         `json:"progress_percent"`
	Confidence        string          `json:"confidence"`
	Status            string          `json:"status"`
	Weight            float64         `json:"weight"`
	UnitLabel         *string         `json:"unit_label,omitempty"`
	DueDate           *time.Time      `json:"due_date,omitempty"`
	Metadata          json.RawMessage `json:"metadata,omitempty"`
	Inactive          bool            `json:"inactive"`
	CreatedAt         time.Time       `json:"created_at"`
	CreatedBy         *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt         time.Time       `json:"updated_at"`
	UpdatedBy         *uuid.UUID      `json:"updated_by,omitempty"`
	ObjectiveTitle    *string         `json:"objective_title,omitempty"`
	CycleName         *string         `json:"cycle_name,omitempty"`
	LatestCheckinDate *time.Time      `json:"latest_checkin_date,omitempty"`
	LatestNote        *string         `json:"latest_note,omitempty"`
}

type KeyResultInput struct {
	TenantID        uuid.UUID
	ObjectiveID     uuid.UUID
	Title           string
	Description     *string
	MetricType      string
	StartValue      float64
	TargetValue     float64
	CurrentValue    float64
	ProgressPercent *float64
	Confidence      string
	Status          string
	Weight          *float64
	UnitLabel       *string
	DueDate         *time.Time
	Metadata        json.RawMessage
}

type KeyResultFilter struct {
	TenantID    uuid.UUID
	ObjectiveID *uuid.UUID
	Status      *string
}

type KeyResultCheckIn struct {
	ID              uuid.UUID       `json:"id"`
	TenantID        uuid.UUID       `json:"tenant_id"`
	KeyResultID     uuid.UUID       `json:"key_result_id"`
	CheckInDate     time.Time       `json:"checkin_date"`
	Value           float64         `json:"value"`
	ProgressPercent float64         `json:"progress_percent"`
	Confidence      string          `json:"confidence"`
	Status          string          `json:"status"`
	Note            *string         `json:"note,omitempty"`
	Metadata        json.RawMessage `json:"metadata,omitempty"`
	Inactive        bool            `json:"inactive"`
	CreatedAt       time.Time       `json:"created_at"`
	CreatedBy       *uuid.UUID      `json:"created_by,omitempty"`
	KeyResultTitle  *string         `json:"key_result_title,omitempty"`
	ObjectiveID     *uuid.UUID      `json:"objective_id,omitempty"`
	ObjectiveTitle  *string         `json:"objective_title,omitempty"`
}

type KeyResultCheckInInput struct {
	TenantID        uuid.UUID
	KeyResultID     uuid.UUID
	CheckInDate     *time.Time
	Value           float64
	ProgressPercent *float64
	Confidence      string
	Status          string
	Note            *string
	Metadata        json.RawMessage
}

type KeyResultCheckInFilter struct {
	TenantID    uuid.UUID
	KeyResultID *uuid.UUID
	ObjectiveID *uuid.UUID
}

type OKRSummaryRow struct {
	OwnerType       string  `json:"owner_type"`
	ObjectiveCount  int32   `json:"objective_count"`
	KeyResultCount  int32   `json:"key_result_count"`
	AverageProgress float64 `json:"average_progress"`
	AtRiskCount     int32   `json:"at_risk_count"`
	CompletedCount  int32   `json:"completed_count"`
}

func NewOKRCycle(input OKRCycleInput) (*OKRCycle, error) {
	if input.TenantID == uuid.Nil || input.StartDate == nil || input.EndDate == nil || input.EndDate.Before(*input.StartDate) {
		return nil, ErrInvalidOKRCycle
	}
	name := strings.TrimSpace(input.Name)
	code := strings.TrimSpace(input.CycleCode)
	if name == "" || code == "" {
		return nil, ErrInvalidOKRCycle
	}
	status := normalizeWorkerProfileEnum(input.Status, OKRCycleStatusDraft)
	if !containsString([]string{OKRCycleStatusDraft, OKRCycleStatusActive, OKRCycleStatusClosed, OKRCycleStatusArchived}, status) {
		return nil, ErrInvalidOKRCycle
	}
	cadence := normalizeWorkerProfileEnum(input.ReviewCadence, OKRReviewCadenceWeekly)
	if !containsString([]string{OKRReviewCadenceWeekly, OKRReviewCadenceBiweekly, OKRReviewCadenceMonthly}, cadence) {
		return nil, ErrInvalidOKRCycle
	}
	metadata := cleanOKRObjectJSON(input.Metadata)
	return &OKRCycle{TenantID: input.TenantID, Name: name, CycleCode: code, Description: cleanOKRStringPtr(input.Description), StartDate: *input.StartDate, EndDate: *input.EndDate, Status: status, ReviewCadence: cadence, Metadata: metadata}, nil
}

func NewObjective(input ObjectiveInput) (*Objective, error) {
	if input.TenantID == uuid.Nil || input.CycleID == uuid.Nil {
		return nil, ErrInvalidObjective
	}
	title := strings.TrimSpace(input.Title)
	if title == "" {
		return nil, ErrInvalidObjective
	}
	ownerType := normalizeWorkerProfileEnum(input.OwnerType, ObjectiveOwnerCompany)
	if !containsString([]string{ObjectiveOwnerCompany, ObjectiveOwnerDepartment, ObjectiveOwnerProject, ObjectiveOwnerWorker}, ownerType) {
		return nil, ErrInvalidObjective
	}
	if !validObjectiveOwner(ownerType, input.OwnerWorkerProfileID, input.OwnerDepartmentID, input.OwnerProjectID) {
		return nil, ErrInvalidObjective
	}
	status := normalizeWorkerProfileEnum(input.Status, ObjectiveStatusDraft)
	if !containsString([]string{ObjectiveStatusDraft, ObjectiveStatusActive, ObjectiveStatusAtRisk, ObjectiveStatusCompleted, ObjectiveStatusClosed, ObjectiveStatusCancelled}, status) {
		return nil, ErrInvalidObjective
	}
	priority := normalizeWorkerProfileEnum(input.Priority, ObjectivePriorityNormal)
	if !containsString([]string{ObjectivePriorityLow, ObjectivePriorityNormal, ObjectivePriorityHigh, ObjectivePriorityCritical}, priority) {
		return nil, ErrInvalidObjective
	}
	if input.StartDate != nil && input.DueDate != nil && input.DueDate.Before(*input.StartDate) {
		return nil, ErrInvalidObjective
	}
	progress := cleanOKRPercent(input.ProgressPercent)
	weight := cleanOKRPositiveFloat(input.Weight, 1)
	return &Objective{TenantID: input.TenantID, CycleID: input.CycleID, ParentObjectiveID: cleanOKRUUIDPtr(input.ParentObjectiveID), OwnerType: ownerType, OwnerWorkerProfileID: cleanOKRUUIDPtr(input.OwnerWorkerProfileID), OwnerDepartmentID: cleanOKRUUIDPtr(input.OwnerDepartmentID), OwnerProjectID: cleanOKRUUIDPtr(input.OwnerProjectID), Title: title, Description: cleanOKRStringPtr(input.Description), Status: status, Priority: priority, ProgressPercent: progress, Weight: weight, StartDate: input.StartDate, DueDate: input.DueDate, Metadata: cleanOKRObjectJSON(input.Metadata)}, nil
}

func NewKeyResult(input KeyResultInput) (*KeyResult, error) {
	if input.TenantID == uuid.Nil || input.ObjectiveID == uuid.Nil {
		return nil, ErrInvalidKeyResult
	}
	title := strings.TrimSpace(input.Title)
	if title == "" {
		return nil, ErrInvalidKeyResult
	}
	metricType := normalizeWorkerProfileEnum(input.MetricType, KeyResultMetricNumber)
	if !containsString([]string{KeyResultMetricNumber, KeyResultMetricPercent, KeyResultMetricCurrency, KeyResultMetricBoolean}, metricType) {
		return nil, ErrInvalidKeyResult
	}
	if metricType == KeyResultMetricBoolean {
		input.StartValue = 0
		input.TargetValue = 1
		if input.CurrentValue != 0 {
			input.CurrentValue = 1
		}
	}
	if input.TargetValue == input.StartValue {
		return nil, ErrInvalidKeyResult
	}
	confidence := normalizeWorkerProfileEnum(input.Confidence, OKRConfidenceMedium)
	if !containsString([]string{OKRConfidenceLow, OKRConfidenceMedium, OKRConfidenceHigh}, confidence) {
		return nil, ErrInvalidKeyResult
	}
	status := normalizeWorkerProfileEnum(input.Status, deriveOKRKeyResultStatus(cleanOKRPercent(input.ProgressPercent)))
	if !containsString([]string{KeyResultStatusNotStarted, KeyResultStatusOnTrack, KeyResultStatusAtRisk, KeyResultStatusBehind, KeyResultStatusCompleted, KeyResultStatusClosed, KeyResultStatusCancelled}, status) {
		return nil, ErrInvalidKeyResult
	}
	progress := cleanOKRPercent(input.ProgressPercent)
	if input.ProgressPercent == nil {
		progress = progressFromOKRValues(input.StartValue, input.TargetValue, input.CurrentValue)
	}
	return &KeyResult{TenantID: input.TenantID, ObjectiveID: input.ObjectiveID, Title: title, Description: cleanOKRStringPtr(input.Description), MetricType: metricType, StartValue: input.StartValue, TargetValue: input.TargetValue, CurrentValue: input.CurrentValue, ProgressPercent: progress, Confidence: confidence, Status: status, Weight: cleanOKRPositiveFloat(input.Weight, 1), UnitLabel: cleanOKRStringPtr(input.UnitLabel), DueDate: input.DueDate, Metadata: cleanOKRObjectJSON(input.Metadata)}, nil
}

func NewKeyResultCheckIn(input KeyResultCheckInInput) (*KeyResultCheckIn, error) {
	if input.TenantID == uuid.Nil || input.KeyResultID == uuid.Nil {
		return nil, ErrInvalidKeyResultCheckIn
	}
	checkInDate := time.Now()
	if input.CheckInDate != nil {
		checkInDate = *input.CheckInDate
	}
	progress := cleanOKRPercent(input.ProgressPercent)
	confidence := normalizeWorkerProfileEnum(input.Confidence, OKRConfidenceMedium)
	if !containsString([]string{OKRConfidenceLow, OKRConfidenceMedium, OKRConfidenceHigh}, confidence) {
		return nil, ErrInvalidKeyResultCheckIn
	}
	status := normalizeWorkerProfileEnum(input.Status, deriveOKRKeyResultStatus(progress))
	if !containsString([]string{KeyResultStatusNotStarted, KeyResultStatusOnTrack, KeyResultStatusAtRisk, KeyResultStatusBehind, KeyResultStatusCompleted, KeyResultStatusClosed, KeyResultStatusCancelled}, status) {
		return nil, ErrInvalidKeyResultCheckIn
	}
	return &KeyResultCheckIn{TenantID: input.TenantID, KeyResultID: input.KeyResultID, CheckInDate: checkInDate, Value: input.Value, ProgressPercent: progress, Confidence: confidence, Status: status, Note: cleanOKRStringPtr(input.Note), Metadata: cleanOKRObjectJSON(input.Metadata)}, nil
}

func NormalizeOKRSearch(value *string) *string {
	if value == nil {
		return nil
	}
	clean := strings.TrimSpace(*value)
	if clean == "" {
		return nil
	}
	return &clean
}

func CalculateKeyResultProgress(metricType string, startValue float64, targetValue float64, currentValue float64) float64 {
	if metricType == KeyResultMetricBoolean {
		if currentValue >= 1 {
			return 100
		}
		return 0
	}
	return progressFromOKRValues(startValue, targetValue, currentValue)
}

func validObjectiveOwner(ownerType string, workerID *uuid.UUID, departmentID *uuid.UUID, projectID *uuid.UUID) bool {
	worker := workerID != nil && *workerID != uuid.Nil
	department := departmentID != nil && *departmentID != uuid.Nil
	project := projectID != nil && *projectID != uuid.Nil
	switch ownerType {
	case ObjectiveOwnerCompany:
		return !worker && !department && !project
	case ObjectiveOwnerDepartment:
		return department && !worker && !project
	case ObjectiveOwnerProject:
		return project && !worker && !department
	case ObjectiveOwnerWorker:
		return worker && !department && !project
	default:
		return false
	}
}

func cleanOKRStringPtr(value *string) *string {
	if value == nil {
		return nil
	}
	clean := strings.TrimSpace(*value)
	if clean == "" {
		return nil
	}
	return &clean
}

func cleanOKRUUIDPtr(value *uuid.UUID) *uuid.UUID {
	if value == nil || *value == uuid.Nil {
		return nil
	}
	return value
}

func cleanOKRObjectJSON(value json.RawMessage) json.RawMessage {
	if len(value) == 0 || !json.Valid(value) {
		return json.RawMessage(`{}`)
	}
	var parsed any
	if err := json.Unmarshal(value, &parsed); err != nil {
		return json.RawMessage(`{}`)
	}
	if _, ok := parsed.(map[string]any); !ok {
		return json.RawMessage(`{}`)
	}
	return value
}

func cleanOKRPercent(value *float64) float64 {
	if value == nil || math.IsNaN(*value) || math.IsInf(*value, 0) {
		return 0
	}
	if *value < 0 {
		return 0
	}
	if *value > 100 {
		return 100
	}
	return math.Round(*value*100) / 100
}

func cleanOKRPositiveFloat(value *float64, fallback float64) float64 {
	if value == nil || math.IsNaN(*value) || math.IsInf(*value, 0) || *value <= 0 {
		return fallback
	}
	return math.Round(*value*100) / 100
}

func progressFromOKRValues(startValue float64, targetValue float64, currentValue float64) float64 {
	span := targetValue - startValue
	if span == 0 {
		return 0
	}
	progress := ((currentValue - startValue) / span) * 100
	if span < 0 {
		progress = ((startValue - currentValue) / (startValue - targetValue)) * 100
	}
	return cleanOKRPercent(&progress)
}

func deriveOKRKeyResultStatus(progress float64) string {
	switch {
	case progress >= 100:
		return KeyResultStatusCompleted
	case progress <= 0:
		return KeyResultStatusNotStarted
	default:
		return KeyResultStatusOnTrack
	}
}
