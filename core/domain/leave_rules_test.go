package domain

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCalculateLeaveDays(t *testing.T) {
	day := time.Date(2026, 6, 11, 15, 30, 0, 0, time.FixedZone("IST", 19800))
	tests := []struct {
		name      string
		start     time.Time
		end       time.Time
		startType string
		endType   string
		expected  float64
		wantErr   error
	}{
		{"single full", day, day, LeaveDayFullDay, LeaveDayFullDay, 1, nil},
		{"single first half", day, day, LeaveDayFirstHalf, LeaveDayFirstHalf, 0.5, nil},
		{"single second half", day, day, LeaveDaySecondHalf, LeaveDaySecondHalf, 0.5, nil},
		{"single two halves", day, day, LeaveDayFirstHalf, LeaveDaySecondHalf, 1, nil},
		{"single one full endpoint", day, day, LeaveDayFirstHalf, LeaveDayFullDay, 1, nil},
		{"two full days", day, day.AddDate(0, 0, 1), LeaveDayFullDay, LeaveDayFullDay, 2, nil},
		{"multi day half endpoints", day, day.AddDate(0, 0, 2), LeaveDaySecondHalf, LeaveDayFirstHalf, 2, nil},
		{"normalizes time zones to date", day, time.Date(2026, 6, 12, 1, 0, 0, 0, time.FixedZone("GST", 14400)), LeaveDayFullDay, LeaveDayFirstHalf, 1.5, nil},
		{"invalid range", day.AddDate(0, 0, 1), day, LeaveDayFullDay, LeaveDayFullDay, 0, ErrInvalidDateRange},
		{"invalid day type", day, day, "morning", LeaveDayFullDay, 0, ErrInvalidEnumValue},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculateLeaveDays(tt.start, tt.end, tt.startType, tt.endType)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected %v, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.expected {
				t.Fatalf("got %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsSandwich(t *testing.T) {
	tests := []struct {
		name     string
		start    time.Time
		end      time.Time
		holidays []time.Time
		weekoffs []time.Weekday
		wantOK   bool
		wantGap  float64
	}{
		{
			name:     "weekend gap",
			start:    time.Date(2026, 6, 12, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC),
			weekoffs: []time.Weekday{time.Saturday, time.Sunday},
			wantOK:   true,
			wantGap:  2,
		},
		{
			name:     "holiday plus weekend gap",
			start:    time.Date(2026, 8, 14, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2026, 8, 17, 0, 0, 0, 0, time.UTC),
			holidays: []time.Time{time.Date(2026, 8, 15, 12, 0, 0, 0, time.UTC)},
			weekoffs: []time.Weekday{time.Sunday},
			wantOK:   true,
			wantGap:  2,
		},
		{
			name:     "working day breaks gap",
			start:    time.Date(2026, 6, 12, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2026, 6, 16, 0, 0, 0, 0, time.UTC),
			weekoffs: []time.Weekday{time.Saturday, time.Sunday},
			wantOK:   false,
			wantGap:  0,
		},
		{
			name:    "same date is never sandwich",
			start:   time.Date(2026, 6, 12, 0, 0, 0, 0, time.UTC),
			end:     time.Date(2026, 6, 12, 0, 0, 0, 0, time.UTC),
			wantOK:  false,
			wantGap: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, gap := IsSandwich(tt.start, tt.end, tt.holidays, tt.weekoffs)
			if ok != tt.wantOK || gap != tt.wantGap {
				t.Fatalf("got ok=%v gap=%v, want ok=%v gap=%v", ok, gap, tt.wantOK, tt.wantGap)
			}
		})
	}
}

func TestDistributeInto12Months(t *testing.T) {
	tests := []struct {
		name string
		days float64
		want [12]int
	}{
		{"zero", 0, [12]int{}},
		{"less than year", 6, [12]int{1, 1, 1, 1, 1, 1}},
		{"even year", 12, [12]int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}},
		{"remainder goes first", 14, [12]int{2, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}},
		{"fractional days use whole days", 14.75, [12]int{2, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}},
		{"more than year", 26, [12]int{3, 3, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DistributeInto12Months(tt.days)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestCalculateLeaveAccrualProratesFixedYearlyFromConfirmation(t *testing.T) {
	confirmed := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
	rule := &LeavePolicyTemplateRule{
		TenantID:         uuidForLeaveRuleTest("tenant"),
		TemplateID:       uuidForLeaveRuleTest("template"),
		LeaveTypeID:      uuidForLeaveRuleTest("leave-type"),
		AccrualMethod:    LeaveAccrualFixedYearly,
		AccrualFrequency: LeaveAccrualFrequencyYearly,
		CreditDays:       7,
		CalculationConfig: map[string]any{
			"prorate":       true,
			"prorate_basis": "confirmation_date",
			"rounding":      "nearest_half",
		},
	}

	got, err := CalculateLeaveAccrual(rule, LeaveAccrualInput{
		JoiningDate:      time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC),
		ConfirmationDate: &confirmed,
		AsOfDate:         time.Date(2027, 3, 31, 0, 0, 0, 0, time.UTC),
		PeriodStart:      time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC),
		PeriodEnd:        time.Date(2027, 3, 31, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.SourceType != LeaveLedgerSourceYearlyAccrual {
		t.Fatalf("source type got %q, want %q", got.SourceType, LeaveLedgerSourceYearlyAccrual)
	}
	if got.Days != 3.5 {
		t.Fatalf("days got %v, want 3.5", got.Days)
	}
}

func TestCalculateLeaveAccrualKeepsDecimalMonthlyCredits(t *testing.T) {
	rule := &LeavePolicyTemplateRule{
		TenantID:         uuidForLeaveRuleTest("tenant"),
		TemplateID:       uuidForLeaveRuleTest("template"),
		LeaveTypeID:      uuidForLeaveRuleTest("leave-type"),
		AccrualMethod:    LeaveAccrualMonthlyFixed,
		AccrualFrequency: LeaveAccrualFrequencyMonthly,
		CreditDays:       2.5,
	}

	got, err := CalculateLeaveAccrual(rule, LeaveAccrualInput{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Days != 2.5 {
		t.Fatalf("days got %v, want 2.5", got.Days)
	}
}

func uuidForLeaveRuleTest(value string) uuid.UUID {
	return uuid.NewSHA1(uuid.NameSpaceOID, []byte("leave-rule-test:"+value))
}
