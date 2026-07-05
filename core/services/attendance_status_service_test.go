package services

import (
	"testing"
	"time"

	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func TestAttendancePunchStats(t *testing.T) {
	base := time.Date(2026, 6, 18, 0, 0, 0, 0, time.UTC)
	in := base.Add(9 * time.Hour)
	mid := base.Add(13 * time.Hour)
	out := base.Add(18 * time.Hour)
	tests := []struct {
		name          string
		records       []*domain.Attendance
		wantFirstIn   *time.Time
		wantLastOut   *time.Time
		wantWorkedMin int32
	}{
		{"empty", nil, nil, nil, 0},
		{"single checkin", []*domain.Attendance{attendancePunch(domain.AttendanceCheckin, in)}, &in, nil, 0},
		{"checkin checkout", []*domain.Attendance{attendancePunch(domain.AttendanceCheckin, in), attendancePunch(domain.AttendanceCheckout, out)}, &in, &out, 540},
		{"uses earliest checkin and latest checkout", []*domain.Attendance{attendancePunch(domain.AttendanceCheckin, mid), attendancePunch(domain.AttendanceCheckout, out), attendancePunch(domain.AttendanceCheckin, in), attendancePunch(domain.AttendanceCheckout, mid)}, &in, &out, 540},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFirst, gotLast, gotWorked := attendancePunchStats(tt.records)
			if !sameTimePtr(gotFirst, tt.wantFirstIn) || !sameTimePtr(gotLast, tt.wantLastOut) || gotWorked != tt.wantWorkedMin {
				t.Fatalf("got first=%v last=%v worked=%d", gotFirst, gotLast, gotWorked)
			}
		})
	}
}

func TestAttendanceRuleOutcomePriority(t *testing.T) {
	start := "09:00"
	end := "18:00"
	absentAfter := int32(180)
	halfAfter := int32(60)
	workingHour := &domain.WorkingHour{IsWorkingDay: true, StartTime: start, EndTime: end}
	policy := &domain.AttendancePolicy{GraceLateMinutes: 10, GraceEarlyMinutes: 10, MinHalfDayMinutes: 240, MinFullDayMinutes: 420, HalfDayLateAfterMinutes: &halfAfter, AbsentLateAfterMinutes: &absentAfter}
	tests := []struct {
		name        string
		firstIn     time.Time
		lastOut     *time.Time
		worked      int32
		wantLate    int32
		wantEarly   int32
		wantOutcome string
	}{
		{"on time", atUTC("2026-06-18T09:05:00Z"), ptrTime(atUTC("2026-06-18T18:00:00Z")), 480, 5, 0, domain.AttendanceRuleOutcomeOnTime},
		{"late", atUTC("2026-06-18T09:30:00Z"), ptrTime(atUTC("2026-06-18T18:00:00Z")), 480, 30, 0, domain.AttendanceRuleOutcomeLate},
		{"half day late threshold wins", atUTC("2026-06-18T10:30:00Z"), ptrTime(atUTC("2026-06-18T18:00:00Z")), 420, 90, 0, domain.AttendanceRuleOutcomeHalfDay},
		{"absent late threshold wins", atUTC("2026-06-18T12:30:00Z"), ptrTime(atUTC("2026-06-18T18:00:00Z")), 300, 210, 0, domain.AttendanceRuleOutcomeAbsent},
		{"early exit", atUTC("2026-06-18T09:00:00Z"), ptrTime(atUTC("2026-06-18T17:30:00Z")), 480, 0, 30, domain.AttendanceRuleOutcomeEarlyExit},
		{"worked minutes absent", atUTC("2026-06-18T09:00:00Z"), ptrTime(atUTC("2026-06-18T11:00:00Z")), 120, 0, 420, domain.AttendanceRuleOutcomeAbsent},
		{"worked minutes half day", atUTC("2026-06-18T09:00:00Z"), ptrTime(atUTC("2026-06-18T15:00:00Z")), 360, 0, 180, domain.AttendanceRuleOutcomeHalfDay},
		{"missing checkout", atUTC("2026-06-18T09:00:00Z"), nil, 0, 0, 0, domain.AttendanceRuleOutcomeMissingCheckout},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			late, early, outcome := attendanceRuleOutcome(&tt.firstIn, tt.lastOut, tt.worked, workingHour, policy, nil)
			if late != tt.wantLate || early != tt.wantEarly || outcome != tt.wantOutcome {
				t.Fatalf("got late=%d early=%d outcome=%s", late, early, outcome)
			}
		})
	}
}

func TestAttendanceCheckoutRequired(t *testing.T) {
	start := "09:00"
	end := "18:00"
	firstIn := atUTC("2026-06-18T09:00:00Z")
	status := &domain.AttendanceDailyStatus{
		Date:         atUTC("2026-06-18T00:00:00Z"),
		FirstCheckIn: &firstIn,
		WorkingHour:  &domain.WorkingHour{IsWorkingDay: true, StartTime: start, EndTime: end},
		Policy:       &domain.AttendancePolicy{GraceEarlyMinutes: 10},
	}
	tests := []struct {
		name string
		day  time.Time
		now  time.Time
		want bool
	}{
		{"same day before shift end", atUTC("2026-06-18T00:00:00Z"), atUTC("2026-06-18T17:00:00Z"), false},
		{"same day after shift end and grace", atUTC("2026-06-18T00:00:00Z"), atUTC("2026-06-18T18:11:00Z"), true},
		{"past day", atUTC("2026-06-17T00:00:00Z"), atUTC("2026-06-18T10:00:00Z"), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status.Date = tt.day
			if got := attendanceCheckoutRequired(tt.day, status, tt.now); got != tt.want {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExplicitAttendanceStatus(t *testing.T) {
	tests := []struct {
		name    string
		records []*domain.Attendance
		want    string
	}{
		{"none", nil, ""},
		{"skips blank", []*domain.Attendance{{Status: nil}}, ""},
		{"returns first explicit", []*domain.Attendance{{Status: testStringPtr("present")}, {Status: testStringPtr("absent")}}, "present"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := explicitAttendanceStatus(tt.records); got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func sameTimePtr(got *time.Time, want *time.Time) bool {
	if got == nil || want == nil {
		return got == want
	}
	return got.Equal(*want)
}

func atUTC(value string) time.Time {
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		panic(err)
	}
	return parsed
}

func ptrTime(value time.Time) *time.Time { return &value }

func testStringPtr(value string) *string { return &value }

func attendancePunch(action string, value time.Time) *domain.Attendance {
	return &domain.Attendance{Type: &action, Time: &value}
}
