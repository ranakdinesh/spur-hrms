package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

var defaultLeaveTypes = []domain.LeaveTypeInput{
	{Name: "Annual Leave", Shortcode: strPtr("AL"), Description: strPtr("Planned vacation or annual paid time off."), IsPaid: true, IsCarryForward: true, MaxCarryForward: 15, IsSystem: true},
	{Name: "Earned Leave", Shortcode: strPtr("EL"), Description: strPtr("Privilege leave earned over time based on service rules."), IsPaid: true, IsCarryForward: true, MaxCarryForward: 30, IsSystem: true},
	{Name: "Casual Leave", Shortcode: strPtr("CL"), Description: strPtr("Short-notice personal leave for urgent needs."), IsPaid: true, IsConsecutiveLimit: true, ConsecutiveDaysLimit: 3, IsSystem: true},
	{Name: "Sick Leave", Shortcode: strPtr("SL"), Description: strPtr("Health-related leave for illness, treatment, or recovery."), IsPaid: true, IsSystem: true},
	{Name: "Medical Leave", Shortcode: strPtr("ML"), Description: strPtr("Extended medical leave requiring documentation as configured by policy."), IsPaid: true, IsSystem: true},
	{Name: "Unpaid Leave", Shortcode: strPtr("LOP"), Description: strPtr("Leave without pay when paid balances are unavailable or not applicable."), IsPaid: false, IsSystem: true},
	{Name: "Bereavement Leave", Shortcode: strPtr("BL"), Description: strPtr("Leave for death or serious family bereavement needs."), IsPaid: true, IsSystem: true},
	{Name: "Compensatory Off", Shortcode: strPtr("CO"), Description: strPtr("Time off granted against approved extra work or holidays worked."), IsPaid: true, IsCarryForward: true, MaxCarryForward: 5, IsSystem: true},
	{Name: "Time Off In Lieu", Shortcode: strPtr("TOIL"), Description: strPtr("Alternative time off in lieu of eligible overtime or extra work."), IsPaid: true, IsCarryForward: true, MaxCarryForward: 5, IsSystem: true},
	{Name: "Maternity Leave", Shortcode: strPtr("MAT"), Description: strPtr("Parental leave for childbirth recovery and newborn care."), IsPaid: true, IsSystem: true},
	{Name: "Paternity Leave", Shortcode: strPtr("PAT"), Description: strPtr("Parental leave for fathers or secondary caregivers."), IsPaid: true, IsSystem: true},
	{Name: "Adoption Leave", Shortcode: strPtr("ADP"), Description: strPtr("Leave for adoption placement and family bonding."), IsPaid: true, IsSystem: true},
	{Name: "Parental Leave", Shortcode: strPtr("PRL"), Description: strPtr("General parent or caregiver leave for family needs."), IsPaid: true, IsSystem: true},
	{Name: "Family and Medical Leave", Shortcode: strPtr("FML"), Description: strPtr("Family care or serious health condition leave category."), IsPaid: false, IsSystem: true},
	{Name: "Religious Leave", Shortcode: strPtr("RL"), Description: strPtr("Leave for religious observance or duties."), IsPaid: false, IsSystem: true},
	{Name: "Jury Duty Leave", Shortcode: strPtr("JDL"), Description: strPtr("Civic leave for jury duty or court-mandated service."), IsPaid: true, IsSystem: true},
	{Name: "Military Leave", Shortcode: strPtr("MIL"), Description: strPtr("Leave for military or reserve service obligations."), IsPaid: false, IsSystem: true},
	{Name: "Sabbatical Leave", Shortcode: strPtr("SAB"), Description: strPtr("Longer planned break for service, study, or renewal."), IsPaid: false, IsSystem: true},
	{Name: "Study Leave", Shortcode: strPtr("STL"), Description: strPtr("Leave for approved education, exams, or training."), IsPaid: true, IsSystem: true},
	{Name: "Other Leave", Shortcode: strPtr("OTH"), Description: strPtr("General fallback leave type for tenant-specific policy mapping."), IsPaid: false, IsSystem: true},
}

func (s *TenantService) CreateLeaveType(ctx context.Context, cmd ports.LeaveTypeCommand) (*domain.LeaveType, error) {
	item, err := s.leaveTypeFromCommand(cmd, false)
	if err != nil {
		s.logError("validate leave type create", err, serviceTenantIDField(cmd.TenantID), serviceStringField("leave_type_name", cmd.Name))
		return nil, err
	}
	result, err := s.leaveTypes.CreateLeaveType(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create leave type", err, serviceTenantIDField(cmd.TenantID), serviceStringField("leave_type_name", item.Name))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListLeaveTypes(ctx context.Context, tenantID uuid.UUID) ([]*domain.LeaveType, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate leave type list tenant", err)
		return nil, err
	}
	if err := s.ensureSystemLeaveTypes(ctx, tenantID, nil); err != nil {
		return nil, err
	}
	items, err := s.leaveTypes.ListLeaveTypes(ctx, tenantID)
	if err != nil {
		s.logError("list leave types", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) GetLeaveType(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.LeaveType, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate leave type get tenant", err)
		return nil, err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidLeaveTypeID
		s.logError("validate leave type get id", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	item, err := s.leaveTypes.GetLeaveType(ctx, tenantID, id)
	if err != nil {
		s.logError("get leave type", err, serviceTenantIDField(tenantID), serviceStringField("leave_type_id", id.String()))
		return nil, err
	}
	return item, nil
}

func (s *TenantService) UpdateLeaveType(ctx context.Context, cmd ports.LeaveTypeCommand) (*domain.LeaveType, error) {
	if cmd.ID == uuid.Nil {
		err := domain.ErrInvalidLeaveTypeID
		s.logError("validate leave type update id", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	existing, err := s.GetLeaveType(ctx, cmd.TenantID, cmd.ID)
	if err != nil {
		return nil, err
	}
	if existing.IsSystem {
		cmd.Name = existing.Name
		cmd.Shortcode = existing.Shortcode
		cmd.Description = existing.Description
		cmd.IsPaid = existing.IsPaid
		cmd.IsCarryForward = existing.IsCarryForward
		cmd.MaxCarryForward = existing.MaxCarryForward
		cmd.IsConsecutiveLimit = existing.IsConsecutiveLimit
		cmd.ConsecutiveDaysLimit = existing.ConsecutiveDaysLimit
	}
	item, err := s.leaveTypeFromCommand(cmd, false)
	if err != nil {
		s.logError("validate leave type update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("leave_type_id", cmd.ID.String()))
		return nil, err
	}
	item.ID = cmd.ID
	result, err := s.leaveTypes.UpdateLeaveType(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update leave type", err, serviceTenantIDField(cmd.TenantID), serviceStringField("leave_type_id", cmd.ID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) DeleteLeaveType(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate leave type delete tenant", err)
		return err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidLeaveTypeID
		s.logError("validate leave type delete id", err, serviceTenantIDField(tenantID))
		return err
	}
	existing, err := s.GetLeaveType(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if existing.IsSystem {
		err := domain.ErrSystemLeaveTypeLocked
		s.logError("validate leave type delete system", err, serviceTenantIDField(tenantID), serviceStringField("leave_type_id", id.String()))
		return err
	}
	if err := s.leaveTypes.DeleteLeaveType(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete leave type", err, serviceTenantIDField(tenantID), serviceStringField("leave_type_id", id.String()))
		return err
	}
	return nil
}

func (s *TenantService) ensureSystemLeaveTypes(ctx context.Context, tenantID uuid.UUID, actorID *uuid.UUID) error {
	for _, input := range defaultLeaveTypes {
		input.TenantID = tenantID
		input.IsEnabled = true
		item, err := domain.NewLeaveType(input)
		if err != nil {
			s.logError("validate system leave type seed", err, serviceTenantIDField(tenantID), serviceStringField("leave_type_name", input.Name))
			return err
		}
		if _, err := s.leaveTypes.UpsertSystemLeaveType(ctx, item, actorID); err != nil {
			s.logError("seed system leave type", err, serviceTenantIDField(tenantID), serviceStringField("leave_type_name", item.Name))
			return err
		}
	}
	return nil
}

func (s *TenantService) leaveTypeFromCommand(cmd ports.LeaveTypeCommand, isSystem bool) (*domain.LeaveType, error) {
	return domain.NewLeaveType(domain.LeaveTypeInput{
		TenantID:             cmd.TenantID,
		Name:                 cmd.Name,
		Shortcode:            cmd.Shortcode,
		Description:          cmd.Description,
		IsPaid:               cmd.IsPaid,
		IsCarryForward:       cmd.IsCarryForward,
		MaxCarryForward:      cmd.MaxCarryForward,
		IsConsecutiveLimit:   cmd.IsConsecutiveLimit,
		ConsecutiveDaysLimit: cmd.ConsecutiveDaysLimit,
		IsEnabled:            cmd.IsEnabled,
		IsSystem:             isSystem,
	})
}

func strPtr(value string) *string {
	return &value
}
