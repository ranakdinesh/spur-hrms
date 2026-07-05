package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidWorkingHourID       = errors.New("working_hour_id is required")
	ErrInvalidWorkingHourDay      = errors.New("day_of_week must be Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, or Sunday")
	ErrInvalidWorkingHourTime     = errors.New("working hour time must use HH:MM format")
	ErrInvalidWorkingHourRange    = errors.New("working hour end_time must be after start_time for working days")
	ErrInvalidWorkingHourBreak    = errors.New("break_minutes must be zero or greater")
	ErrInvalidWorkingHourScope    = errors.New("working hour scope can target tenant, branch, or user, not branch and user together")
	ErrInvalidWorkingHourBranch   = errors.New("branch_id is required")
	ErrInvalidWorkingHourUser     = errors.New("user_id is required")
	ErrWorkingHourNotFound        = errors.New("working hour not found")
	ErrNoTenantWorkingHoursToCopy = errors.New("tenant default working hours are required before copying to a branch")
)

var dayOrder = map[string]int{
	"Monday":    1,
	"Tuesday":   2,
	"Wednesday": 3,
	"Thursday":  4,
	"Friday":    5,
	"Saturday":  6,
	"Sunday":    7,
}

type WorkingHour struct {
	ID           uuid.UUID  `json:"id"`
	TenantID     uuid.UUID  `json:"tenant_id"`
	BranchID     *uuid.UUID `json:"branch_id,omitempty"`
	UserID       *uuid.UUID `json:"user_id,omitempty"`
	DayOfWeek    string     `json:"day_of_week"`
	IsWorkingDay bool       `json:"is_working_day"`
	StartTime    string     `json:"start_time"`
	EndTime      string     `json:"end_time"`
	BreakMinutes int32      `json:"break_minutes"`
	Source       string     `json:"source,omitempty"`
	Inactive     bool       `json:"inactive"`
	CreatedAt    time.Time  `json:"created_at"`
	CreatedBy    *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt    time.Time  `json:"updated_at"`
	UpdatedBy    *uuid.UUID `json:"updated_by,omitempty"`
}

type WorkingHourInput struct {
	TenantID     uuid.UUID
	BranchID     *uuid.UUID
	UserID       *uuid.UUID
	DayOfWeek    string
	IsWorkingDay bool
	StartTime    string
	EndTime      string
	BreakMinutes int32
}

func NewWorkingHour(input WorkingHourInput) (*WorkingHour, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	if cleanUUIDOptional(input.BranchID) != nil && cleanUUIDOptional(input.UserID) != nil {
		return nil, ErrInvalidWorkingHourScope
	}
	day, err := NormalizeDayOfWeek(input.DayOfWeek)
	if err != nil {
		return nil, err
	}
	startTime, endTime, err := normalizeWorkingHourTimes(input.StartTime, input.EndTime)
	if err != nil {
		return nil, err
	}
	if input.IsWorkingDay && endTime <= startTime {
		return nil, ErrInvalidWorkingHourRange
	}
	if input.BreakMinutes < 0 {
		return nil, ErrInvalidWorkingHourBreak
	}
	now := time.Now().UTC()
	return &WorkingHour{
		TenantID:     input.TenantID,
		BranchID:     cleanUUIDOptional(input.BranchID),
		UserID:       cleanUUIDOptional(input.UserID),
		DayOfWeek:    day,
		IsWorkingDay: input.IsWorkingDay,
		StartTime:    formatMinutes(startTime),
		EndTime:      formatMinutes(endTime),
		BreakMinutes: input.BreakMinutes,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

func NormalizeDayOfWeek(value string) (string, error) {
	clean := strings.ToLower(strings.TrimSpace(value))
	for day := range dayOrder {
		if strings.ToLower(day) == clean {
			return day, nil
		}
	}
	return "", ErrInvalidWorkingHourDay
}

func WorkingHourScope(item *WorkingHour) string {
	if item == nil {
		return "tenant"
	}
	if item.UserID != nil {
		return "user"
	}
	if item.BranchID != nil {
		return "branch"
	}
	return "tenant"
}

func DefaultWorkingHour(tenantID uuid.UUID, dayOfWeek string) (*WorkingHour, error) {
	day, err := NormalizeDayOfWeek(dayOfWeek)
	if err != nil {
		return nil, err
	}
	isWorkingDay := dayOrder[day] <= 5
	input := WorkingHourInput{
		TenantID:     tenantID,
		DayOfWeek:    day,
		IsWorkingDay: isWorkingDay,
		StartTime:    "09:00",
		EndTime:      "18:00",
		BreakMinutes: 60,
	}
	if !isWorkingDay {
		input.StartTime = "00:00"
		input.EndTime = "00:00"
		input.BreakMinutes = 0
	}
	item, err := NewWorkingHour(input)
	if err != nil {
		return nil, err
	}
	item.Source = "system_default"
	return item, nil
}

func DefaultWorkingHourInputs(tenantID uuid.UUID) []WorkingHourInput {
	items := make([]WorkingHourInput, 0, len(dayOrder))
	for _, day := range []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"} {
		isWorkingDay := dayOrder[day] <= 5
		input := WorkingHourInput{TenantID: tenantID, DayOfWeek: day, IsWorkingDay: isWorkingDay, StartTime: "09:00", EndTime: "18:00", BreakMinutes: 60}
		if !isWorkingDay {
			input.StartTime = "00:00"
			input.EndTime = "00:00"
			input.BreakMinutes = 0
		}
		items = append(items, input)
	}
	return items
}

func normalizeWorkingHourTimes(startValue string, endValue string) (int, int, error) {
	start, err := parseHourMinute(startValue)
	if err != nil {
		return 0, 0, err
	}
	end, err := parseHourMinute(endValue)
	if err != nil {
		return 0, 0, err
	}
	return start, end, nil
}

func parseHourMinute(value string) (int, error) {
	parsed, err := time.Parse("15:04", strings.TrimSpace(value))
	if err != nil {
		return 0, ErrInvalidWorkingHourTime
	}
	return parsed.Hour()*60 + parsed.Minute(), nil
}

func formatMinutes(value int) string {
	return fmt.Sprintf("%02d:%02d", value/60, value%60)
}
