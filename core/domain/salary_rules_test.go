package domain

import (
	"math"
	"testing"
)

func TestCalculateUserSalary(t *testing.T) {
	structure := []SalaryItem{
		{ItemType: SalaryItemEarning, Code: SalaryCodeBasic, Amount: 30000},
		{ItemType: SalaryItemEarning, Code: SalaryCodeHRA, Amount: 10000},
		{ItemType: SalaryItemDeduction, Code: "pf", Amount: 2000},
	}
	tests := []struct {
		name              string
		gross             float64
		presentDays       int
		absentDays        int
		daysInMonth       int
		special           bool
		wantEarnings      float64
		wantDeductions    float64
		wantAbsent        float64
		wantNet           float64
		wantLWP           bool
		wantOriginalItems int
	}{
		{"regular full month", 40000, 30, 0, 30, false, 40000, 2000, 0, 38000, false, len(structure)},
		{"absent deduction", 40000, 28, 2, 30, false, 40000, roundMoney(2000 + 40000.0/30.0*2.0), roundMoney(40000.0 / 30.0 * 2.0), roundMoney(40000 - 2000 - 40000.0/30.0*2.0), true, len(structure)},
		{"special prorata with absent", 40000, 28, 2, 30, true, roundMoney(40000.0 / 30.0 * 28.0), roundMoney(2000 + 40000.0/30.0*2.0), roundMoney(40000.0 / 30.0 * 2.0), roundMoney(40000.0/30.0*28.0 - (2000 + 40000.0/30.0*2.0)), true, len(structure)},
		{"zero days in month avoids divide by zero", 40000, 28, 2, 0, true, 40000, 2000, 0, 38000, false, len(structure)},
		{"special with no present days is not prorated", 40000, 0, 0, 30, true, 40000, 2000, 0, 38000, false, len(structure)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateUserSalary(tt.gross, structure, tt.presentDays, tt.absentDays, tt.daysInMonth, tt.special)
			if !near(got.TotalEarnings, tt.wantEarnings) || !near(got.TotalDeductions, tt.wantDeductions) || !near(got.AbsentDeduction, tt.wantAbsent) || !near(got.NetSalary, tt.wantNet) {
				t.Fatalf("unexpected salary result: %#v", got)
			}
			if !tt.wantLWP && len(got.Items) != tt.wantOriginalItems {
				t.Fatalf("expected no LWP item: %#v", got.Items)
			}
			if tt.wantLWP {
				last := got.Items[len(got.Items)-1]
				if last.Code != SalaryCodeLWP || !near(last.Amount, tt.wantAbsent) {
					t.Fatalf("expected LWP deduction item, got %#v", got.Items)
				}
			}
		})
	}
}

func near(got, want float64) bool {
	return math.Abs(got-want) < 0.000001
}
