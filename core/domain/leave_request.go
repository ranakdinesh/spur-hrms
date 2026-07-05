package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidLeaveID           = errors.New("leave_id is required")
	ErrInvalidLeaveUser         = errors.New("leave user_id is required")
	ErrInvalidLeaveType         = errors.New("leave_type_id is required")
	ErrInvalidLeaveDays         = errors.New("leave days must be greater than zero")
	ErrLeaveBalanceInsufficient = errors.New("insufficient leave balance")
	ErrLeaveOverlap             = errors.New("leave overlaps an existing request")
	ErrLeaveDatesOutsideFY      = errors.New("leave dates must be within the selected financial year")
	ErrLeaveRequestBelowMinimum = errors.New("leave request is below minimum days configured for this policy")
	ErrLeaveRequestAboveMaximum = errors.New("leave request exceeds maximum days configured for this policy")
	ErrLeaveHalfDayNotAllowed   = errors.New("half-day leave is not allowed for this policy")
	ErrLeaveProbationRestricted = errors.New("earned leave is not available during probation")
	ErrLeaveNotFound            = errors.New("leave not found")
)

type Leave struct {
	ID            uuid.UUID  `json:"id"`
	TenantID      uuid.UUID  `json:"tenant_id"`
	UserID        uuid.UUID  `json:"user_id"`
	LeaveTypeID   uuid.UUID  `json:"leave_type_id"`
	FYID          uuid.UUID  `json:"fy_id"`
	StartDate     time.Time  `json:"start_date"`
	EndDate       time.Time  `json:"end_date"`
	StartDayType  string     `json:"start_day_type"`
	EndDayType    string     `json:"end_day_type"`
	Days          float64    `json:"days"`
	Reason        *string    `json:"reason,omitempty"`
	Status        string     `json:"status"`
	AppliedDate   time.Time  `json:"applied_date"`
	FromLeaveType *uuid.UUID `json:"from_leave_type,omitempty"`
	ToLeaveType   *uuid.UUID `json:"to_leave_type,omitempty"`
	IsSandwich    bool       `json:"is_sandwich"`
	Inactive      bool       `json:"inactive"`
	CreatedAt     time.Time  `json:"created_at"`
	CreatedBy     *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt     time.Time  `json:"updated_at"`
	UpdatedBy     *uuid.UUID `json:"updated_by,omitempty"`
}

type LeaveApproval struct {
	ID                uuid.UUID  `json:"id"`
	TenantID          uuid.UUID  `json:"tenant_id"`
	LeaveID           uuid.UUID  `json:"leave_id"`
	ApproverID        uuid.UUID  `json:"approver_id"`
	Status            string     `json:"status"`
	Remarks           *string    `json:"remarks,omitempty"`
	ActionDate        *time.Time `json:"action_date,omitempty"`
	WorkflowID        *uuid.UUID `json:"workflow_id,omitempty"`
	WorkflowStepID    *uuid.UUID `json:"workflow_step_id,omitempty"`
	StepOrder         int32      `json:"step_order"`
	DecisionRule      string     `json:"decision_rule"`
	RequiredApprovals int32      `json:"required_approvals"`
	Inactive          bool       `json:"inactive"`
	CreatedAt         time.Time  `json:"created_at"`
	CreatedBy         *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt         time.Time  `json:"updated_at"`
	UpdatedBy         *uuid.UUID `json:"updated_by,omitempty"`
}

type LeaveApplication struct {
	Leave    *Leave         `json:"leave"`
	Approval *LeaveApproval `json:"approval"`
	Balance  *LeaveBalance  `json:"balance"`
}

type LeaveReportFilter struct {
	TenantID     uuid.UUID  `json:"tenant_id"`
	ManagerID    *uuid.UUID `json:"manager_id,omitempty"`
	FYID         *uuid.UUID `json:"fy_id,omitempty"`
	UserID       *uuid.UUID `json:"user_id,omitempty"`
	DepartmentID *uuid.UUID `json:"department_id,omitempty"`
	LeaveTypeID  *uuid.UUID `json:"leave_type_id,omitempty"`
	Status       *string    `json:"status,omitempty"`
	StartDate    *time.Time `json:"start_date,omitempty"`
	EndDate      *time.Time `json:"end_date,omitempty"`
}

type LeaveReportRow struct {
	ID                 uuid.UUID  `json:"id"`
	TenantID           uuid.UUID  `json:"tenant_id"`
	UserID             uuid.UUID  `json:"user_id"`
	EmployeeCode       *string    `json:"employee_code,omitempty"`
	Firstname          string     `json:"firstname"`
	Lastname           *string    `json:"lastname,omitempty"`
	ReportingManagerID *uuid.UUID `json:"reporting_manager_id,omitempty"`
	DepartmentID       *uuid.UUID `json:"department_id,omitempty"`
	DepartmentName     *string    `json:"department_name,omitempty"`
	DesignationID      *uuid.UUID `json:"designation_id,omitempty"`
	DesignationName    *string    `json:"designation_name,omitempty"`
	LeaveTypeID        uuid.UUID  `json:"leave_type_id"`
	LeaveTypeName      *string    `json:"leave_type_name,omitempty"`
	LeaveTypeShortcode *string    `json:"leave_type_shortcode,omitempty"`
	FYID               uuid.UUID  `json:"fy_id"`
	FinancialYearName  *string    `json:"financial_year_name,omitempty"`
	StartDate          time.Time  `json:"start_date"`
	EndDate            time.Time  `json:"end_date"`
	StartDayType       string     `json:"start_day_type"`
	EndDayType         string     `json:"end_day_type"`
	Days               float64    `json:"days"`
	Reason             *string    `json:"reason,omitempty"`
	Status             string     `json:"status"`
	IsSandwich         bool       `json:"is_sandwich"`
	AppliedDate        time.Time  `json:"applied_date"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

type LeaveReportSummary struct {
	TotalRequests int32   `json:"total_requests"`
	TotalDays     float64 `json:"total_days"`
	EmployeeCount int32   `json:"employee_count"`
	PendingCount  int32   `json:"pending_count"`
	ApprovedCount int32   `json:"approved_count"`
	RejectedCount int32   `json:"rejected_count"`
	CanceledCount int32   `json:"canceled_count"`
	PendingDays   float64 `json:"pending_days"`
	ApprovedDays  float64 `json:"approved_days"`
	RejectedDays  float64 `json:"rejected_days"`
	CanceledDays  float64 `json:"canceled_days"`
}

func NewLeaveApplication(tenantID, userID, leaveTypeID, fyID uuid.UUID, startDate, endDate time.Time, startDayType, endDayType string, reason *string, days float64, isSandwich bool) (*Leave, error) {
	if tenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	if userID == uuid.Nil {
		return nil, ErrInvalidLeaveUser
	}
	if leaveTypeID == uuid.Nil {
		return nil, ErrInvalidLeaveType
	}
	if fyID == uuid.Nil {
		return nil, ErrInvalidLeavePolicyFY
	}
	if endDate.Before(startDate) {
		return nil, ErrInvalidDateRange
	}
	if days <= 0 {
		return nil, ErrInvalidLeaveDays
	}
	normalizedStartDayType, err := ValidateLeaveDayType(startDayType)
	if err != nil {
		return nil, err
	}
	normalizedEndDayType, err := ValidateLeaveDayType(endDayType)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	return &Leave{TenantID: tenantID, UserID: userID, LeaveTypeID: leaveTypeID, FYID: fyID, StartDate: dateOnly(startDate), EndDate: dateOnly(endDate), StartDayType: normalizedStartDayType, EndDayType: normalizedEndDayType, Days: days, Reason: cleanOptional(reason), Status: LeaveStatusPending, IsSandwich: isSandwich, AppliedDate: now, CreatedAt: now, UpdatedAt: now}, nil
}

func LeaveDayPartMaskForDate(leaveStart, leaveEnd time.Time, startDayType, endDayType string, day time.Time) (int, error) {
	leaveStart = dateOnly(leaveStart)
	leaveEnd = dateOnly(leaveEnd)
	day = dateOnly(day)
	if day.Before(leaveStart) || day.After(leaveEnd) {
		return 0, nil
	}
	if leaveStart.Equal(leaveEnd) {
		startMask, err := leaveDayTypeMask(startDayType)
		if err != nil {
			return 0, err
		}
		endMask, err := leaveDayTypeMask(endDayType)
		if err != nil {
			return 0, err
		}
		return startMask | endMask, nil
	}
	if day.Equal(leaveStart) {
		return leaveDayTypeMask(startDayType)
	}
	if day.Equal(leaveEnd) {
		return leaveDayTypeMask(endDayType)
	}
	return 3, nil
}

func LeavesOverlap(a *Leave, b *Leave) (bool, error) {
	if a == nil || b == nil {
		return false, nil
	}
	start := a.StartDate
	if b.StartDate.After(start) {
		start = b.StartDate
	}
	end := a.EndDate
	if b.EndDate.Before(end) {
		end = b.EndDate
	}
	if end.Before(start) {
		return false, nil
	}
	for day := dateOnly(start); !day.After(end); day = day.AddDate(0, 0, 1) {
		aMask, err := LeaveDayPartMaskForDate(a.StartDate, a.EndDate, a.StartDayType, a.EndDayType, day)
		if err != nil {
			return false, err
		}
		bMask, err := LeaveDayPartMaskForDate(b.StartDate, b.EndDate, b.StartDayType, b.EndDayType, day)
		if err != nil {
			return false, err
		}
		if aMask&bMask != 0 {
			return true, nil
		}
	}
	return false, nil
}

func leaveDayTypeMask(dayType string) (int, error) {
	normalized, err := ValidateLeaveDayType(dayType)
	if err != nil {
		return 0, err
	}
	switch normalized {
	case LeaveDayFirstHalf:
		return 1, nil
	case LeaveDaySecondHalf:
		return 2, nil
	default:
		return 3, nil
	}
}
