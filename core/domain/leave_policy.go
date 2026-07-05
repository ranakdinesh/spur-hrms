package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

const (
	LeaveAllocationFixed   = "fixed"
	LeaveAllocationMonthly = "monthly"
)

var (
	ErrInvalidLeavePolicyID          = errors.New("leave_policy_id is required")
	ErrInvalidLeavePolicyType        = errors.New("leave policy leave_type_id is required")
	ErrInvalidLeavePolicyFY          = errors.New("leave policy fy_id is required")
	ErrInvalidLeavePolicyDays        = errors.New("leave policy total_days cannot be negative")
	ErrInvalidLeaveAllocationType    = errors.New("leave policy allocation_type must be fixed or monthly")
	ErrInvalidLeaveMonthlyAllocation = errors.New("leave policy monthly allocation cannot be negative")
	ErrInvalidLeaveAllocationMonth   = errors.New("leave allocation month must be between 1 and 12")
	ErrLeavePolicyNotFound           = errors.New("leave policy not found")
	ErrLeavePolicyAlreadyExists      = errors.New("leave policy already exists for this leave type and financial year")
)

type LeavePolicy struct {
	ID                   uuid.UUID  `json:"id"`
	TenantID             uuid.UUID  `json:"tenant_id"`
	LeaveTypeID          uuid.UUID  `json:"leave_type_id"`
	FYID                 uuid.UUID  `json:"fy_id"`
	TotalDays            float64    `json:"total_days"`
	AllocationType       string     `json:"allocation_type"`
	Jan                  int32      `json:"jan"`
	Feb                  int32      `json:"feb"`
	Mar                  int32      `json:"mar"`
	Apr                  int32      `json:"apr"`
	May                  int32      `json:"may"`
	Jun                  int32      `json:"jun"`
	Jul                  int32      `json:"jul"`
	Aug                  int32      `json:"aug"`
	Sep                  int32      `json:"sep"`
	Oct                  int32      `json:"oct"`
	Nov                  int32      `json:"nov"`
	Dec                  int32      `json:"dec"`
	IsSandwichApplicable bool       `json:"is_sandwich_applicable"`
	Inactive             bool       `json:"inactive"`
	CreatedAt            time.Time  `json:"created_at"`
	CreatedBy            *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt            time.Time  `json:"updated_at"`
	UpdatedBy            *uuid.UUID `json:"updated_by,omitempty"`
}

type MonthlyLeaveAllocation struct {
	TenantID    uuid.UUID `json:"tenant_id"`
	PolicyID    uuid.UUID `json:"policy_id"`
	LeaveTypeID uuid.UUID `json:"leave_type_id"`
	FYID        uuid.UUID `json:"fy_id"`
	Month       int32     `json:"month"`
	Days        int32     `json:"days"`
}

func (p *LeavePolicy) AllocationForMonth(month int32) (int32, error) {
	if month < 1 || month > 12 {
		return 0, ErrInvalidLeaveAllocationMonth
	}
	values := []int32{p.Jan, p.Feb, p.Mar, p.Apr, p.May, p.Jun, p.Jul, p.Aug, p.Sep, p.Oct, p.Nov, p.Dec}
	return values[month-1], nil
}

type LeavePolicyInput struct {
	TenantID             uuid.UUID
	LeaveTypeID          uuid.UUID
	FYID                 uuid.UUID
	TotalDays            float64
	AllocationType       string
	Monthly              [12]int
	IsSandwichApplicable bool
}

func NewLeavePolicy(input LeavePolicyInput) (*LeavePolicy, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	if input.LeaveTypeID == uuid.Nil {
		return nil, ErrInvalidLeavePolicyType
	}
	if input.FYID == uuid.Nil {
		return nil, ErrInvalidLeavePolicyFY
	}
	if input.TotalDays < 0 {
		return nil, ErrInvalidLeavePolicyDays
	}
	allocationType := input.AllocationType
	if allocationType == "" {
		allocationType = LeaveAllocationFixed
	}
	if allocationType != LeaveAllocationFixed && allocationType != LeaveAllocationMonthly {
		return nil, ErrInvalidLeaveAllocationType
	}
	monthly := input.Monthly
	if allocationType == LeaveAllocationMonthly {
		if monthly == [12]int{} {
			monthly = DistributeInto12Months(input.TotalDays)
		}
		for _, value := range monthly {
			if value < 0 {
				return nil, ErrInvalidLeaveMonthlyAllocation
			}
		}
	} else {
		monthly = [12]int{}
	}
	now := time.Now().UTC()
	return &LeavePolicy{
		TenantID:             input.TenantID,
		LeaveTypeID:          input.LeaveTypeID,
		FYID:                 input.FYID,
		TotalDays:            input.TotalDays,
		AllocationType:       allocationType,
		Jan:                  int32(monthly[0]),
		Feb:                  int32(monthly[1]),
		Mar:                  int32(monthly[2]),
		Apr:                  int32(monthly[3]),
		May:                  int32(monthly[4]),
		Jun:                  int32(monthly[5]),
		Jul:                  int32(monthly[6]),
		Aug:                  int32(monthly[7]),
		Sep:                  int32(monthly[8]),
		Oct:                  int32(monthly[9]),
		Nov:                  int32(monthly[10]),
		Dec:                  int32(monthly[11]),
		IsSandwichApplicable: input.IsSandwichApplicable,
		CreatedAt:            now,
		UpdatedAt:            now,
	}, nil
}
