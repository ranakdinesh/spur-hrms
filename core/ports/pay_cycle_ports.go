package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type PayCycleRepo interface {
	GetPayCycle(ctx context.Context, tenantID uuid.UUID) (*domain.PayCycle, error)
	UpsertPayCycle(ctx context.Context, item *domain.PayCycle, actorID *uuid.UUID) (*domain.PayCycle, error)
	DeletePayCycle(ctx context.Context, tenantID uuid.UUID, actorID *uuid.UUID) error
}

type PayCycleCommand struct {
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
	ActorID                *uuid.UUID `json:"-"`
}

type PayCyclePeriodQuery struct {
	TenantID uuid.UUID
	Month    int `json:"month"`
	Year     int `json:"year"`
}
