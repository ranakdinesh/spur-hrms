package domain

import "github.com/google/uuid"

type LeavePreview struct {
	TenantID                    uuid.UUID        `json:"tenant_id"`
	UserID                      uuid.UUID        `json:"user_id"`
	LeaveTypeID                 uuid.UUID        `json:"leave_type_id"`
	FYID                        uuid.UUID        `json:"fy_id"`
	StartDate                   string           `json:"start_date"`
	EndDate                     string           `json:"end_date"`
	StartDayType                string           `json:"start_day_type"`
	EndDayType                  string           `json:"end_day_type"`
	BaseDays                    float64          `json:"base_days"`
	SandwichDays                float64          `json:"sandwich_days"`
	TotalDays                   float64          `json:"total_days"`
	IsSandwich                  bool             `json:"is_sandwich"`
	BalanceBefore               float64          `json:"balance_before"`
	PendingBefore               float64          `json:"pending_before"`
	UsedBefore                  float64          `json:"used_before"`
	BalanceAfter                float64          `json:"balance_after"`
	PendingAfter                float64          `json:"pending_after"`
	PaidLeave                   bool             `json:"paid_leave"`
	Allowed                     bool             `json:"allowed"`
	BlockingReasons             []string         `json:"blocking_reasons,omitempty"`
	Warnings                    []string         `json:"warnings,omitempty"`
	RequiresAttachment          bool             `json:"requires_attachment"`
	AttachmentRequiredAfterDays *float64         `json:"attachment_required_after_days,omitempty"`
	NoticeRequired              bool             `json:"notice_required"`
	NoticeDays                  int32            `json:"notice_days"`
	PayrollImpact               string           `json:"payroll_impact,omitempty"`
	EffectivePolicy             *PolicySet       `json:"effective_policy,omitempty"`
	EffectiveLeaveRule          *LeavePolicyRule `json:"effective_leave_rule,omitempty"`
}
