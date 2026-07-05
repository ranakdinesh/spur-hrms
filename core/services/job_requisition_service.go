package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) CreateJobRequisition(ctx context.Context, cmd ports.JobRequisitionCommand) (*domain.JobRequisition, error) {
	if cmd.RequestedBy == uuid.Nil && cmd.ActorID != nil {
		cmd.RequestedBy = *cmd.ActorID
	}
	if err := s.applyJobPositionDefaults(ctx, &cmd); err != nil {
		return nil, err
	}
	item, err := domain.NewJobRequisition(jobRequisitionInput(cmd))
	if err != nil {
		s.logError("validate job requisition create", err, serviceTenantIDField(cmd.TenantID), serviceStringField("title", cmd.Title))
		return nil, err
	}
	if err := s.validateJobRequisitionReferences(ctx, item); err != nil {
		return nil, err
	}
	result, err := s.jobRequisitions.CreateJobRequisition(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create job requisition", err, serviceTenantIDField(cmd.TenantID), serviceStringField("title", item.Title))
		return nil, err
	}
	if _, err := s.writeJobRequisitionLog(ctx, result.TenantID, result.ID, nil, result.Status, "Created", cmd.Notes, cmd.ActorID); err != nil {
		return nil, err
	}
	return s.GetJobRequisition(ctx, result.TenantID, result.ID)
}

func (s *TenantService) ListJobRequisitions(ctx context.Context, filter domain.JobRequisitionFilter) (*domain.JobRequisitionPage, error) {
	if filter.TenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate job requisition list tenant", err)
		return nil, err
	}
	filter.Search = cleanStringPtr(filter.Search)
	filter.Status = cleanStringPtr(filter.Status)
	if filter.Status != nil {
		status, err := domain.ValidateRequisitionStatus(*filter.Status)
		if err != nil {
			s.logError("validate job requisition list status", err, serviceTenantIDField(filter.TenantID), serviceStringField("status", *filter.Status))
			return nil, domain.ErrInvalidJobRequisitionStatus
		}
		filter.Status = &status
	}
	limit, offset := normalizeListWindow(filter.Limit, filter.Offset)
	filter.Limit = limit
	filter.Offset = offset
	items, err := s.jobRequisitions.ListJobRequisitions(ctx, filter)
	if err != nil {
		s.logError("list job requisitions", err, serviceTenantIDField(filter.TenantID))
		return nil, err
	}
	total, err := s.jobRequisitions.CountJobRequisitions(ctx, filter)
	if err != nil {
		s.logError("count job requisitions", err, serviceTenantIDField(filter.TenantID))
		return nil, err
	}
	page := &domain.JobRequisitionPage{Items: items, Total: total, Limit: limit, Offset: offset}
	if int64(offset)+int64(len(items)) < total {
		next := offset + limit
		page.NextOffset = &next
	}
	return page, nil
}

func (s *TenantService) GetJobRequisition(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.JobRequisition, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate job requisition get tenant", err)
		return nil, err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidJobRequisitionID
		s.logError("validate job requisition get id", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	item, err := s.jobRequisitions.GetJobRequisition(ctx, tenantID, id)
	if err != nil {
		s.logError("get job requisition", err, serviceTenantIDField(tenantID), serviceStringField("job_requisition_id", id.String()))
		return nil, err
	}
	logs, err := s.jobRequisitions.ListJobRequisitionLogs(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	item.Logs = logs
	item.LogCount = int32(len(logs))
	return item, nil
}

func (s *TenantService) UpdateJobRequisition(ctx context.Context, cmd ports.JobRequisitionCommand) (*domain.JobRequisition, error) {
	if cmd.ID == uuid.Nil {
		err := domain.ErrInvalidJobRequisitionID
		s.logError("validate job requisition update id", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	existing, err := s.GetJobRequisition(ctx, cmd.TenantID, cmd.ID)
	if err != nil {
		return nil, err
	}
	if existing.Status != domain.ReqStatusDraft && existing.Status != domain.ReqStatusRejected {
		err := domain.ErrInvalidJobRequisitionAction
		s.logError("validate job requisition update status", err, serviceTenantIDField(cmd.TenantID), serviceStringField("job_requisition_id", cmd.ID.String()), serviceStringField("status", existing.Status))
		return nil, err
	}
	if cmd.RequestedBy == uuid.Nil {
		cmd.RequestedBy = existing.RequestedBy
	}
	if cmd.Status == "" {
		cmd.Status = existing.Status
	}
	if err := s.applyJobPositionDefaults(ctx, &cmd); err != nil {
		return nil, err
	}
	item, err := domain.NewJobRequisition(jobRequisitionInput(cmd))
	if err != nil {
		s.logError("validate job requisition update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("job_requisition_id", cmd.ID.String()))
		return nil, err
	}
	item.ID = cmd.ID
	if err := s.validateJobRequisitionReferences(ctx, item); err != nil {
		return nil, err
	}
	result, err := s.jobRequisitions.UpdateJobRequisition(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update job requisition", err, serviceTenantIDField(cmd.TenantID), serviceStringField("job_requisition_id", cmd.ID.String()))
		return nil, err
	}
	if _, err := s.writeJobRequisitionLog(ctx, result.TenantID, result.ID, &existing.Status, result.Status, "Updated", cmd.Notes, cmd.ActorID); err != nil {
		return nil, err
	}
	return s.GetJobRequisition(ctx, result.TenantID, result.ID)
}

func (s *TenantService) SubmitJobRequisition(ctx context.Context, cmd ports.JobRequisitionActionCommand) (*domain.JobRequisition, error) {
	return s.transitionJobRequisition(ctx, cmd, []string{domain.ReqStatusDraft, domain.ReqStatusRejected}, domain.ReqStatusPending, "Submitted")
}

func (s *TenantService) ApproveJobRequisition(ctx context.Context, cmd ports.JobRequisitionActionCommand) (*domain.JobRequisition, error) {
	return s.transitionJobRequisition(ctx, cmd, []string{domain.ReqStatusPending}, domain.ReqStatusApproved, "Approved")
}

func (s *TenantService) RejectJobRequisition(ctx context.Context, cmd ports.JobRequisitionActionCommand) (*domain.JobRequisition, error) {
	return s.transitionJobRequisition(ctx, cmd, []string{domain.ReqStatusPending}, domain.ReqStatusRejected, "Rejected")
}

func (s *TenantService) CloseJobRequisition(ctx context.Context, cmd ports.JobRequisitionActionCommand) (*domain.JobRequisition, error) {
	return s.transitionJobRequisition(ctx, cmd, []string{domain.ReqStatusDraft, domain.ReqStatusPending, domain.ReqStatusApproved, domain.ReqStatusRejected}, domain.ReqStatusClosed, "Closed")
}

func (s *TenantService) DeleteJobRequisition(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if _, err := s.GetJobRequisition(ctx, tenantID, id); err != nil {
		return err
	}
	if err := s.jobRequisitions.DeleteJobRequisition(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete job requisition", err, serviceTenantIDField(tenantID), serviceStringField("job_requisition_id", id.String()))
		return err
	}
	return nil
}

func (s *TenantService) ListJobRequisitionLogs(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) ([]*domain.JobRequisitionLog, error) {
	if _, err := s.GetJobRequisition(ctx, tenantID, id); err != nil {
		return nil, err
	}
	return s.jobRequisitions.ListJobRequisitionLogs(ctx, tenantID, id)
}

func (s *TenantService) transitionJobRequisition(ctx context.Context, cmd ports.JobRequisitionActionCommand, allowed []string, nextStatus string, action string) (*domain.JobRequisition, error) {
	if cmd.ActorID == nil || *cmd.ActorID == uuid.Nil {
		err := domain.ErrInvalidJobRequisitionUser
		s.logError("validate job requisition action actor", err, serviceTenantIDField(cmd.TenantID), serviceStringField("job_requisition_id", cmd.ID.String()))
		return nil, err
	}
	existing, err := s.GetJobRequisition(ctx, cmd.TenantID, cmd.ID)
	if err != nil {
		return nil, err
	}
	if !containsString(allowed, existing.Status) {
		err := domain.ErrInvalidJobRequisitionAction
		s.logError("validate job requisition transition", err, serviceTenantIDField(cmd.TenantID), serviceStringField("job_requisition_id", cmd.ID.String()), serviceStringField("from_status", existing.Status), serviceStringField("to_status", nextStatus))
		return nil, err
	}
	result, err := s.jobRequisitions.UpdateJobRequisitionStatus(ctx, cmd.TenantID, cmd.ID, nextStatus, cmd.Remarks, cmd.ActorID)
	if err != nil {
		s.logError("transition job requisition", err, serviceTenantIDField(cmd.TenantID), serviceStringField("job_requisition_id", cmd.ID.String()), serviceStringField("to_status", nextStatus))
		return nil, err
	}
	if _, err := s.writeJobRequisitionLog(ctx, cmd.TenantID, cmd.ID, &existing.Status, nextStatus, action, cmd.Remarks, cmd.ActorID); err != nil {
		return nil, err
	}
	return s.GetJobRequisition(ctx, result.TenantID, result.ID)
}

func (s *TenantService) writeJobRequisitionLog(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, fromStatus *string, toStatus string, action string, remarks *string, actorID *uuid.UUID) (*domain.JobRequisitionLog, error) {
	item := &domain.JobRequisitionLog{TenantID: tenantID, JobRequisitionID: id, FromStatus: fromStatus, ToStatus: toStatus, Action: action, Remarks: cleanStringPtr(remarks)}
	result, err := s.jobRequisitions.CreateJobRequisitionLog(ctx, item, actorID)
	if err != nil {
		s.logError("write job requisition log", err, serviceTenantIDField(tenantID), serviceStringField("job_requisition_id", id.String()), serviceStringField("action", action))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) applyJobPositionDefaults(ctx context.Context, cmd *ports.JobRequisitionCommand) error {
	position, err := s.GetJobPosition(ctx, cmd.TenantID, cmd.JobPositionID)
	if err != nil {
		return err
	}
	if cmd.Title == "" {
		cmd.Title = position.Title
	}
	if cmd.Level == nil {
		cmd.Level = position.Level
	}
	if cmd.Category == nil {
		cmd.Category = position.Category
	}
	if cmd.DepartmentID == nil {
		cmd.DepartmentID = position.DepartmentID
	}
	if cmd.EmploymentTypeID == nil {
		cmd.EmploymentTypeID = position.EmploymentTypeID
	}
	if cmd.WorkMode == nil {
		cmd.WorkMode = position.WorkMode
	}
	return nil
}

func (s *TenantService) validateJobRequisitionReferences(ctx context.Context, item *domain.JobRequisition) error {
	if _, err := s.GetJobPosition(ctx, item.TenantID, item.JobPositionID); err != nil {
		return err
	}
	if item.DepartmentID != nil {
		if _, err := s.GetDepartment(ctx, item.TenantID, *item.DepartmentID); err != nil {
			return err
		}
	}
	if item.EmploymentTypeID != nil {
		if _, err := s.GetEmploymentType(ctx, item.TenantID, *item.EmploymentTypeID); err != nil {
			return err
		}
	}
	return nil
}

func jobRequisitionInput(cmd ports.JobRequisitionCommand) domain.JobRequisitionInput {
	return domain.JobRequisitionInput{TenantID: cmd.TenantID, JobPositionID: cmd.JobPositionID, Code: cmd.Code, Title: cmd.Title, Level: cmd.Level, Category: cmd.Category, DepartmentID: cmd.DepartmentID, EmploymentTypeID: cmd.EmploymentTypeID, Description: cmd.Description, WorkMode: cmd.WorkMode, TotalOpenings: cmd.TotalOpenings, ReasonForHire: cmd.ReasonForHire, MinSalary: cmd.MinSalary, MaxSalary: cmd.MaxSalary, Currency: cmd.Currency, TargetHireDate: cmd.TargetHireDate, ExpectedClosureDate: cmd.ExpectedClosureDate, RequestedBy: cmd.RequestedBy, RequestedDate: cmd.RequestedDate, Priority: cmd.Priority, Status: cmd.Status, Notes: cmd.Notes}
}
