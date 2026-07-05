package services

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) ResendEmployeeCredentials(ctx context.Context, cmd ports.EmployeeCredentialActionCommand) (*domain.EmployeeCredentialEvent, error) {
	profile, err := s.credentialEmployeeProfile(ctx, cmd)
	if err != nil {
		return nil, err
	}
	email := strings.ToLower(strings.TrimSpace(valueFromPtr(profile.Employee.Email)))
	if email == "" {
		err := domain.ErrEmployeeCredentialTarget
		s.logError("validate employee credential target", err, serviceTenantIDField(cmd.TenantID), serviceStringField("employee_id", cmd.EmployeeID.String()))
		return s.recordCredentialEvent(ctx, profile.Employee, domain.CredentialEventResendCredentials, email, domain.CredentialDeliveryFailed, err.Error(), cmd.ActorID)
	}
	if s.employeeIdentity == nil {
		err := domain.ErrEmployeeCredentialPortMissing
		s.logError("resend employee credentials identity port missing", err, serviceTenantIDField(cmd.TenantID), serviceStringField("employee_id", cmd.EmployeeID.String()))
		return s.recordCredentialEvent(ctx, profile.Employee, domain.CredentialEventResendCredentials, email, domain.CredentialDeliveryFailed, err.Error(), cmd.ActorID)
	}
	err = s.employeeIdentity.SendEmployeePasswordReset(ctx, ports.EmployeeCredentialResetCommand{
		TenantID: cmd.TenantID, UserID: profile.Employee.UserID, Email: email, EmployeeID: profile.Employee.ID,
		Employee: employeeCredentialDisplayName(profile.Employee), EmployeeCode: valueFromPtr(profile.Employee.EmployeeCode), ActorID: cmd.ActorID,
	})
	if err != nil {
		s.logError("send employee credential reset", err, serviceTenantIDField(cmd.TenantID), serviceStringField("employee_id", cmd.EmployeeID.String()), serviceStringField("user_id", profile.Employee.UserID.String()))
		return s.recordCredentialEvent(ctx, profile.Employee, domain.CredentialEventResendCredentials, email, domain.CredentialDeliveryFailed, err.Error(), cmd.ActorID)
	}
	return s.recordCredentialEvent(ctx, profile.Employee, domain.CredentialEventResendCredentials, email, domain.CredentialDeliverySent, "", cmd.ActorID)
}

func (s *TenantService) ResetEmployeeTemporaryPassword(ctx context.Context, cmd ports.EmployeeCredentialActionCommand) (*domain.EmployeeCredentialEvent, error) {
	if len(strings.TrimSpace(cmd.TemporaryPassword)) < 8 {
		err := domain.ErrInvalidTemporaryPassword
		s.logError("validate employee temporary password", err, serviceTenantIDField(cmd.TenantID), serviceStringField("employee_id", cmd.EmployeeID.String()))
		return nil, err
	}
	profile, err := s.credentialEmployeeProfile(ctx, cmd)
	if err != nil {
		return nil, err
	}
	email := strings.ToLower(strings.TrimSpace(valueFromPtr(profile.Employee.Email)))
	if email == "" {
		err := domain.ErrEmployeeCredentialTarget
		s.logError("validate employee temporary password target", err, serviceTenantIDField(cmd.TenantID), serviceStringField("employee_id", cmd.EmployeeID.String()))
		return s.recordCredentialEvent(ctx, profile.Employee, domain.CredentialEventResetTemporaryPassword, email, domain.CredentialDeliveryFailed, err.Error(), cmd.ActorID)
	}
	if s.employeeIdentity == nil {
		err := domain.ErrEmployeeCredentialPortMissing
		s.logError("reset employee temporary password identity port missing", err, serviceTenantIDField(cmd.TenantID), serviceStringField("employee_id", cmd.EmployeeID.String()))
		return s.recordCredentialEvent(ctx, profile.Employee, domain.CredentialEventResetTemporaryPassword, email, domain.CredentialDeliveryFailed, err.Error(), cmd.ActorID)
	}
	err = s.employeeIdentity.SetEmployeeTemporaryPassword(ctx, ports.EmployeeTemporaryPasswordCommand{
		TenantID: cmd.TenantID, UserID: profile.Employee.UserID, Email: email, TemporaryPassword: strings.TrimSpace(cmd.TemporaryPassword),
		EmployeeID: profile.Employee.ID, Employee: employeeCredentialDisplayName(profile.Employee), EmployeeCode: valueFromPtr(profile.Employee.EmployeeCode), ActorID: cmd.ActorID,
	})
	if err != nil {
		s.logError("set employee temporary password", err, serviceTenantIDField(cmd.TenantID), serviceStringField("employee_id", cmd.EmployeeID.String()), serviceStringField("user_id", profile.Employee.UserID.String()))
		return s.recordCredentialEvent(ctx, profile.Employee, domain.CredentialEventResetTemporaryPassword, email, domain.CredentialDeliveryFailed, err.Error(), cmd.ActorID)
	}
	return s.recordCredentialEvent(ctx, profile.Employee, domain.CredentialEventResetTemporaryPassword, email, domain.CredentialDeliverySent, "", cmd.ActorID)
}

func (s *TenantService) ListEmployeeCredentialEvents(ctx context.Context, tenantID uuid.UUID, employeeID uuid.UUID, limit int32) ([]*domain.EmployeeCredentialEvent, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate employee credential event tenant", err)
		return nil, err
	}
	if employeeID == uuid.Nil {
		err := domain.ErrInvalidEmployeeID
		s.logError("validate employee credential event employee", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	items, err := s.employeeCredentialEvents.ListEmployeeCredentialEvents(ctx, tenantID, employeeID, limit)
	if err != nil {
		s.logError("list employee credential events", err, serviceTenantIDField(tenantID), serviceStringField("employee_id", employeeID.String()))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) credentialEmployeeProfile(ctx context.Context, cmd ports.EmployeeCredentialActionCommand) (*domain.EmployeeProfile, error) {
	if cmd.TenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate employee credential tenant", err)
		return nil, err
	}
	if cmd.EmployeeID == uuid.Nil {
		err := domain.ErrInvalidEmployeeID
		s.logError("validate employee credential employee", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	profile, err := s.employees.GetEmployeeProfile(ctx, cmd.TenantID, cmd.EmployeeID)
	if err != nil {
		s.logError("get employee for credential action", err, serviceTenantIDField(cmd.TenantID), serviceStringField("employee_id", cmd.EmployeeID.String()))
		return nil, err
	}
	if profile == nil || profile.Employee == nil || profile.Employee.UserID == uuid.Nil {
		err := domain.ErrInvalidEmployeeUserID
		s.logError("validate employee credential user", err, serviceTenantIDField(cmd.TenantID), serviceStringField("employee_id", cmd.EmployeeID.String()))
		return nil, err
	}
	return profile, nil
}

func (s *TenantService) recordCredentialEvent(ctx context.Context, employee *domain.EmployeeListItem, eventType string, target string, status string, reason string, actorID *uuid.UUID) (*domain.EmployeeCredentialEvent, error) {
	if employee == nil {
		return nil, domain.ErrInvalidEmployeeID
	}
	var failure *string
	if strings.TrimSpace(reason) != "" {
		clean := strings.TrimSpace(reason)
		failure = &clean
	}
	event := &domain.EmployeeCredentialEvent{
		TenantID: employee.TenantID, EmployeeID: employee.ID, UserID: employee.UserID, EventType: eventType,
		DeliveryChannel: domain.CredentialDeliveryEmail, DeliveryTarget: target, Status: status, FailureReason: failure, CreatedBy: actorID,
	}
	saved, err := s.employeeCredentialEvents.CreateEmployeeCredentialEvent(ctx, event)
	if err != nil {
		s.logError("record employee credential event", err, serviceTenantIDField(employee.TenantID), serviceStringField("employee_id", employee.ID.String()), serviceStringField("event_type", eventType))
		return nil, err
	}
	return saved, nil
}

func employeeCredentialDisplayName(employee *domain.EmployeeListItem) string {
	if employee == nil {
		return ""
	}
	return strings.TrimSpace(strings.Join([]string{employee.Firstname, valueFromPtr(employee.MiddleName), valueFromPtr(employee.Lastname)}, " "))
}
