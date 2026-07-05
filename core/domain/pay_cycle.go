package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	PayCycleMonthly     = "monthly"
	PayCycleSemiMonthly = "semi_monthly"
	PayCycleWeekly      = "weekly"
	PayCycleBiWeekly    = "bi_weekly"
	PayCycleCustom      = "custom"

	PayrollAttendanceSourceAttendance = "attendance"
	PayrollAttendanceSourceManual     = "manual"
	PayrollAttendanceSourceImport     = "import"
	PayrollAttendanceSourceNone       = "none"

	PayrollAttendancePeriodCurrentMonth  = "current_month"
	PayrollAttendancePeriodPreviousMonth = "previous_month"
	PayrollAttendancePeriodCustomDays    = "custom_days"

	PayrollPayoutSameMonth   = "same_month"
	PayrollPayoutNextMonth   = "next_month"
	PayrollPayoutFixedOffset = "fixed_offset"

	PayrollProrationCalendarDays = "calendar_days"
	PayrollProrationWorkingDays  = "working_days"
	PayrollProrationFixed26      = "fixed_26"
	PayrollProrationFixed30      = "fixed_30"

	PayrollArrearsSameCycle = "same_cycle"
	PayrollArrearsNextCycle = "next_cycle"
	PayrollArrearsManual    = "manual"

	PayrollRoundingNone         = "none"
	PayrollRoundingNearestRupee = "nearest_rupee"
	PayrollRoundingCeilRupee    = "ceil_rupee"
	PayrollRoundingFloorRupee   = "floor_rupee"
	PayrollRoundingTwoDecimals  = "two_decimals"

	PayrollPaymentBankTransfer = "bank_transfer"
	PayrollPaymentCash         = "cash"
	PayrollPaymentCheque       = "cheque"
	PayrollPaymentUPI          = "upi"
	PayrollPaymentMixed        = "mixed"

	PayrollPaymentFileBankCSV  = "bank_csv"
	PayrollPaymentFileBankXLSX = "bank_xlsx"
	PayrollPaymentFileNACH     = "nach"
	PayrollPaymentFileNone     = "none"
	PayrollPaymentFileCustom   = "custom"
)

var (
	ErrInvalidPayCycleConfig = errors.New("invalid pay cycle configuration")
	ErrPayCycleNotFound      = errors.New("pay cycle configuration not found")
)

type PayCycle struct {
	ID                     uuid.UUID  `json:"id"`
	TenantID               uuid.UUID  `json:"tenant_id"`
	Name                   string     `json:"name"`
	CycleType              string     `json:"cycle_type"`
	PayDay                 *int32     `json:"pay_day,omitempty"`
	StartDay               *int32     `json:"start_day,omitempty"`
	EndDay                 *int32     `json:"end_day,omitempty"`
	AttendanceSource       string     `json:"attendance_source"`
	AttendancePeriodType   string     `json:"attendance_period_type"`
	AttendanceCutoffDay    int32      `json:"attendance_cutoff_day"`
	PayoutTiming           string     `json:"payout_timing"`
	PayoutOffsetDays       int32      `json:"payout_offset_days"`
	IncludeWeeklyOffs      bool       `json:"include_weekly_offs"`
	IncludeHolidays        bool       `json:"include_holidays"`
	ProrateJoiningExit     bool       `json:"prorate_joining_exit"`
	ProrationBasis         string     `json:"proration_basis"`
	AllowArrears           bool       `json:"allow_arrears"`
	ArrearsMode            string     `json:"arrears_mode"`
	AllowNegativeNetPay    bool       `json:"allow_negative_net_pay"`
	OvertimeComponentCode  *string    `json:"overtime_component_code,omitempty"`
	LWPComponentCode       string     `json:"lwp_component_code"`
	RoundingMode           string     `json:"rounding_mode"`
	PaymentMode            string     `json:"payment_mode"`
	PaymentFileFormat      string     `json:"payment_file_format"`
	RequiresApproval       bool       `json:"requires_approval"`
	AutoLockAfterApproval  bool       `json:"auto_lock_after_approval"`
	PayrollLockDay         *int32     `json:"payroll_lock_day,omitempty"`
	PFEnabled              bool       `json:"pf_enabled"`
	PFEmployeeRate         float64    `json:"pf_employee_rate"`
	PFEmployerRate         float64    `json:"pf_employer_rate"`
	PFWageCeiling          float64    `json:"pf_wage_ceiling"`
	PFApplyCeiling         bool       `json:"pf_apply_ceiling"`
	ESIEnabled             bool       `json:"esi_enabled"`
	ESIEmployeeRate        float64    `json:"esi_employee_rate"`
	ESIEmployerRate        float64    `json:"esi_employer_rate"`
	ESIWageCeiling         float64    `json:"esi_wage_ceiling"`
	ProfessionalTaxEnabled bool       `json:"professional_tax_enabled"`
	TDSEnabled             bool       `json:"tds_enabled"`
	CountryCode            string     `json:"country_code"`
	StateCode              *string    `json:"state_code,omitempty"`
	Notes                  *string    `json:"notes,omitempty"`
	Inactive               bool       `json:"inactive"`
	CreatedAt              time.Time  `json:"created_at"`
	CreatedBy              *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt              time.Time  `json:"updated_at"`
	UpdatedBy              *uuid.UUID `json:"updated_by,omitempty"`
}

type PayCyclePeriod struct {
	PeriodStart      time.Time  `json:"period_start"`
	PeriodEnd        time.Time  `json:"period_end"`
	AttendanceStart  time.Time  `json:"attendance_start"`
	AttendanceEnd    time.Time  `json:"attendance_end"`
	AttendanceCutoff time.Time  `json:"attendance_cutoff"`
	PayoutDate       time.Time  `json:"payout_date"`
	LockDate         *time.Time `json:"lock_date,omitempty"`
	Month            int        `json:"month"`
	Year             int        `json:"year"`
}

type PayCycleInput struct {
	TenantID               uuid.UUID
	Name                   string
	CycleType              string
	PayDay                 *int32
	StartDay               *int32
	EndDay                 *int32
	AttendanceSource       string
	AttendancePeriodType   string
	AttendanceCutoffDay    int32
	PayoutTiming           string
	PayoutOffsetDays       int32
	IncludeWeeklyOffs      bool
	IncludeHolidays        bool
	ProrateJoiningExit     bool
	ProrationBasis         string
	AllowArrears           bool
	ArrearsMode            string
	AllowNegativeNetPay    bool
	OvertimeComponentCode  *string
	LWPComponentCode       string
	RoundingMode           string
	PaymentMode            string
	PaymentFileFormat      string
	RequiresApproval       bool
	AutoLockAfterApproval  bool
	PayrollLockDay         *int32
	PFEnabled              bool
	PFEmployeeRate         float64
	PFEmployerRate         float64
	PFWageCeiling          float64
	PFApplyCeiling         bool
	ESIEnabled             bool
	ESIEmployeeRate        float64
	ESIEmployerRate        float64
	ESIWageCeiling         float64
	ProfessionalTaxEnabled bool
	TDSEnabled             bool
	CountryCode            string
	StateCode              *string
	Notes                  *string
}

func NewPayCycle(input PayCycleInput) (*PayCycle, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	item := defaultPayCycle(input.TenantID)
	item.Name = firstClean(input.Name, item.Name)
	item.CycleType = firstClean(input.CycleType, item.CycleType)
	item.PayDay = cleanDayPtr(input.PayDay)
	item.StartDay = cleanDayPtr(input.StartDay)
	item.EndDay = cleanDayPtr(input.EndDay)
	item.AttendanceSource = firstClean(input.AttendanceSource, item.AttendanceSource)
	item.AttendancePeriodType = firstClean(input.AttendancePeriodType, item.AttendancePeriodType)
	item.AttendanceCutoffDay = defaultInt32(input.AttendanceCutoffDay, item.AttendanceCutoffDay)
	item.PayoutTiming = firstClean(input.PayoutTiming, item.PayoutTiming)
	item.PayoutOffsetDays = input.PayoutOffsetDays
	item.IncludeWeeklyOffs = input.IncludeWeeklyOffs
	item.IncludeHolidays = input.IncludeHolidays
	item.ProrateJoiningExit = input.ProrateJoiningExit
	item.ProrationBasis = firstClean(input.ProrationBasis, item.ProrationBasis)
	item.AllowArrears = input.AllowArrears
	item.ArrearsMode = firstClean(input.ArrearsMode, item.ArrearsMode)
	item.AllowNegativeNetPay = input.AllowNegativeNetPay
	item.OvertimeComponentCode = cleanOptional(input.OvertimeComponentCode)
	item.LWPComponentCode = firstClean(input.LWPComponentCode, item.LWPComponentCode)
	item.RoundingMode = firstClean(input.RoundingMode, item.RoundingMode)
	item.PaymentMode = firstClean(input.PaymentMode, item.PaymentMode)
	item.PaymentFileFormat = firstClean(input.PaymentFileFormat, item.PaymentFileFormat)
	item.RequiresApproval = input.RequiresApproval
	item.AutoLockAfterApproval = input.AutoLockAfterApproval
	item.PayrollLockDay = cleanDayPtr(input.PayrollLockDay)
	item.PFEnabled = input.PFEnabled
	item.PFEmployeeRate = defaultFloat(input.PFEmployeeRate, item.PFEmployeeRate)
	item.PFEmployerRate = defaultFloat(input.PFEmployerRate, item.PFEmployerRate)
	item.PFWageCeiling = defaultFloat(input.PFWageCeiling, item.PFWageCeiling)
	item.PFApplyCeiling = input.PFApplyCeiling
	item.ESIEnabled = input.ESIEnabled
	item.ESIEmployeeRate = defaultFloat(input.ESIEmployeeRate, item.ESIEmployeeRate)
	item.ESIEmployerRate = defaultFloat(input.ESIEmployerRate, item.ESIEmployerRate)
	item.ESIWageCeiling = defaultFloat(input.ESIWageCeiling, item.ESIWageCeiling)
	item.ProfessionalTaxEnabled = input.ProfessionalTaxEnabled
	item.TDSEnabled = input.TDSEnabled
	item.CountryCode = strings.ToUpper(firstClean(input.CountryCode, item.CountryCode))
	item.StateCode = cleanOptional(input.StateCode)
	item.Notes = cleanOptional(input.Notes)
	if !validPayCycle(item) {
		return nil, ErrInvalidPayCycleConfig
	}
	return item, nil
}

func defaultPayCycle(tenantID uuid.UUID) *PayCycle {
	now := time.Now().UTC()
	payDay := int32(30)
	startDay := int32(1)
	endDay := int32(31)
	return &PayCycle{TenantID: tenantID, Name: "Monthly Payroll", CycleType: PayCycleMonthly, PayDay: &payDay, StartDay: &startDay, EndDay: &endDay, AttendanceSource: PayrollAttendanceSourceAttendance, AttendancePeriodType: PayrollAttendancePeriodCurrentMonth, AttendanceCutoffDay: 25, PayoutTiming: PayrollPayoutSameMonth, IncludeWeeklyOffs: true, IncludeHolidays: true, ProrateJoiningExit: true, ProrationBasis: PayrollProrationCalendarDays, AllowArrears: true, ArrearsMode: PayrollArrearsNextCycle, LWPComponentCode: "lwp", RoundingMode: PayrollRoundingNearestRupee, PaymentMode: PayrollPaymentBankTransfer, PaymentFileFormat: PayrollPaymentFileBankCSV, RequiresApproval: true, AutoLockAfterApproval: true, PFEnabled: true, PFEmployeeRate: 12, PFEmployerRate: 12, PFWageCeiling: 15000, PFApplyCeiling: true, ESIEnabled: true, ESIEmployeeRate: 0.75, ESIEmployerRate: 3.25, ESIWageCeiling: 21000, ProfessionalTaxEnabled: true, TDSEnabled: true, CountryCode: "IN", CreatedAt: now, UpdatedAt: now}
}

func validPayCycle(item *PayCycle) bool {
	return item != nil && item.Name != "" && containsString([]string{PayCycleMonthly, PayCycleSemiMonthly, PayCycleWeekly, PayCycleBiWeekly, PayCycleCustom}, item.CycleType) && containsString([]string{PayrollAttendanceSourceAttendance, PayrollAttendanceSourceManual, PayrollAttendanceSourceImport, PayrollAttendanceSourceNone}, item.AttendanceSource) && containsString([]string{PayrollAttendancePeriodCurrentMonth, PayrollAttendancePeriodPreviousMonth, PayrollAttendancePeriodCustomDays}, item.AttendancePeriodType) && containsString([]string{PayrollPayoutSameMonth, PayrollPayoutNextMonth, PayrollPayoutFixedOffset}, item.PayoutTiming) && containsString([]string{PayrollProrationCalendarDays, PayrollProrationWorkingDays, PayrollProrationFixed26, PayrollProrationFixed30}, item.ProrationBasis) && containsString([]string{PayrollArrearsSameCycle, PayrollArrearsNextCycle, PayrollArrearsManual}, item.ArrearsMode) && containsString([]string{PayrollRoundingNone, PayrollRoundingNearestRupee, PayrollRoundingCeilRupee, PayrollRoundingFloorRupee, PayrollRoundingTwoDecimals}, item.RoundingMode) && containsString([]string{PayrollPaymentBankTransfer, PayrollPaymentCash, PayrollPaymentCheque, PayrollPaymentUPI, PayrollPaymentMixed}, item.PaymentMode) && containsString([]string{PayrollPaymentFileBankCSV, PayrollPaymentFileBankXLSX, PayrollPaymentFileNACH, PayrollPaymentFileNone, PayrollPaymentFileCustom}, item.PaymentFileFormat) && validDay(item.AttendanceCutoffDay) && item.PayoutOffsetDays >= 0 && item.PayoutOffsetDays <= 31 && item.PFEmployeeRate >= 0 && item.PFEmployerRate >= 0 && item.PFWageCeiling >= 0 && item.ESIEmployeeRate >= 0 && item.ESIEmployerRate >= 0 && item.ESIWageCeiling >= 0
}

func containsString(items []string, value string) bool {
	for _, item := range items {
		if item == value {
			return true
		}
	}
	return false
}

func validDay(value int32) bool { return value >= 1 && value <= 31 }

func cleanDayPtr(value *int32) *int32 {
	if value == nil || *value < 1 || *value > 31 {
		return nil
	}
	clean := *value
	return &clean
}

func firstClean(value string, fallback string) string {
	clean := strings.TrimSpace(value)
	if clean == "" {
		return fallback
	}
	return clean
}

func defaultInt32(value int32, fallback int32) int32 {
	if value == 0 {
		return fallback
	}
	return value
}

func defaultFloat(value float64, fallback float64) float64 {
	if value == 0 {
		return fallback
	}
	return value
}
