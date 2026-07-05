package services

import (
	"context"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) CreateCelebration(ctx context.Context, cmd ports.CelebrationCommand) (*domain.CelebrationListItem, error) {
	item, err := s.celebrationFromCommand(ctx, cmd)
	if err != nil {
		return nil, err
	}
	if err := s.validateCelebrationUnique(ctx, item, uuid.Nil); err != nil {
		return nil, err
	}
	created, err := s.celebrations.CreateCelebration(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create celebration", err, serviceTenantIDField(cmd.TenantID), serviceStringField("celebration_type_id", cmd.CelebrationTypeID.String()))
		return nil, err
	}
	return s.enrichCelebration(ctx, cmd.TenantID, created)
}

func (s *TenantService) ListCelebrations(ctx context.Context, tenantID uuid.UUID) ([]*domain.CelebrationListItem, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate celebration list tenant", err)
		return nil, err
	}
	items, err := s.celebrations.ListCelebrations(ctx, tenantID)
	if err != nil {
		s.logError("list celebrations", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	return s.enrichCelebrations(ctx, tenantID, items)
}

func (s *TenantService) GetCelebration(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.CelebrationListItem, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate celebration get tenant", err)
		return nil, err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidCelebrationID
		s.logError("validate celebration get id", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	item, err := s.celebrations.GetCelebration(ctx, tenantID, id)
	if err != nil {
		s.logError("get celebration", err, serviceTenantIDField(tenantID), serviceStringField("celebration_id", id.String()))
		return nil, err
	}
	return s.enrichCelebration(ctx, tenantID, item)
}

func (s *TenantService) UpdateCelebration(ctx context.Context, cmd ports.CelebrationCommand) (*domain.CelebrationListItem, error) {
	if cmd.ID == uuid.Nil {
		err := domain.ErrInvalidCelebrationID
		s.logError("validate celebration update id", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	item, err := s.celebrationFromCommand(ctx, cmd)
	if err != nil {
		return nil, err
	}
	item.ID = cmd.ID
	if err := s.validateCelebrationUnique(ctx, item, cmd.ID); err != nil {
		return nil, err
	}
	updated, err := s.celebrations.UpdateCelebration(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update celebration", err, serviceTenantIDField(cmd.TenantID), serviceStringField("celebration_id", cmd.ID.String()))
		return nil, err
	}
	return s.enrichCelebration(ctx, cmd.TenantID, updated)
}

func (s *TenantService) DeleteCelebration(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate celebration delete tenant", err)
		return err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidCelebrationID
		s.logError("validate celebration delete id", err, serviceTenantIDField(tenantID))
		return err
	}
	if err := s.celebrations.DeleteCelebration(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete celebration", err, serviceTenantIDField(tenantID), serviceStringField("celebration_id", id.String()))
		return err
	}
	return nil
}

func (s *TenantService) celebrationFromCommand(ctx context.Context, cmd ports.CelebrationCommand) (*domain.Celebration, error) {
	date, err := parseOptionalDate(cmd.CelebrationDate)
	if err != nil || date == nil {
		s.logError("validate celebration date", domain.ErrInvalidCelebrationDate, serviceTenantIDField(cmd.TenantID), serviceStringField("celebration_date", cmd.CelebrationDate))
		return nil, domain.ErrInvalidCelebrationDate
	}
	celebrationType, err := s.GetCelebrationType(ctx, cmd.TenantID, cmd.CelebrationTypeID)
	if err != nil {
		return nil, err
	}
	if celebrationType.IsUserCelebration && (cmd.UserID == nil || *cmd.UserID == uuid.Nil) {
		err := domain.ErrInvalidCelebrationUser
		s.logError("validate celebration user", err, serviceTenantIDField(cmd.TenantID), serviceStringField("celebration_type_id", cmd.CelebrationTypeID.String()))
		return nil, err
	}
	if cmd.UserID != nil && *cmd.UserID != uuid.Nil {
		if _, err := s.employees.GetEmployeeByUserID(ctx, cmd.TenantID, *cmd.UserID); err != nil {
			s.logError("validate celebration employee", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", (*cmd.UserID).String()))
			return nil, err
		}
	}
	if cmd.BranchID != nil && *cmd.BranchID != uuid.Nil {
		if _, err := s.branches.GetBranch(ctx, cmd.TenantID, *cmd.BranchID); err != nil {
			s.logError("validate celebration branch", err, serviceTenantIDField(cmd.TenantID), serviceStringField("branch_id", (*cmd.BranchID).String()))
			return nil, err
		}
	}
	item, err := domain.NewCelebration(domain.CelebrationInput{TenantID: cmd.TenantID, BranchID: cmd.BranchID, UserID: cmd.UserID, CelebrationTypeID: cmd.CelebrationTypeID, CelebrationDate: date, CustomTitle: cmd.CustomTitle, Description: cmd.Description})
	if err != nil {
		s.logError("validate celebration", err, serviceTenantIDField(cmd.TenantID), serviceStringField("celebration_type_id", cmd.CelebrationTypeID.String()))
		return nil, err
	}
	return item, nil
}

func (s *TenantService) validateCelebrationUnique(ctx context.Context, item *domain.Celebration, currentID uuid.UUID) error {
	if item == nil || item.UserID == nil || *item.UserID == uuid.Nil {
		return nil
	}
	existing, err := s.celebrations.ListCelebrationsByUser(ctx, item.TenantID, *item.UserID)
	if err != nil {
		s.logError("list celebration uniqueness rows", err, serviceTenantIDField(item.TenantID), serviceStringField("user_id", (*item.UserID).String()))
		return err
	}
	for _, row := range existing {
		if row != nil && row.CelebrationTypeID == item.CelebrationTypeID && row.ID != currentID {
			err := domain.ErrCelebrationAlreadyExists
			s.logError("validate celebration duplicate user type", err, serviceTenantIDField(item.TenantID), serviceStringField("user_id", (*item.UserID).String()), serviceStringField("celebration_type_id", item.CelebrationTypeID.String()))
			return err
		}
	}
	return nil
}

func (s *TenantService) enrichCelebration(ctx context.Context, tenantID uuid.UUID, item *domain.Celebration) (*domain.CelebrationListItem, error) {
	items, err := s.enrichCelebrations(ctx, tenantID, []*domain.Celebration{item})
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, domain.ErrCelebrationNotFound
	}
	return items[0], nil
}

func (s *TenantService) enrichCelebrations(ctx context.Context, tenantID uuid.UUID, items []*domain.Celebration) ([]*domain.CelebrationListItem, error) {
	types, err := s.celebrations.ListCelebrationTypes(ctx, tenantID)
	if err != nil {
		s.logError("load celebration list types", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	typeByID := map[uuid.UUID]*domain.CelebrationType{}
	for _, item := range types {
		if item != nil {
			typeByID[item.ID] = item
		}
	}
	employees, err := s.ListEmployees(ctx, tenantID)
	if err != nil {
		s.logError("load celebration list employees", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	employeeByUser := map[uuid.UUID]*domain.EmployeeListItem{}
	for _, employee := range employees {
		if employee != nil {
			employeeByUser[employee.UserID] = employee
		}
	}
	branches, err := s.ListBranches(ctx, tenantID)
	if err != nil {
		s.logError("load celebration list branches", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	branchByID := map[uuid.UUID]*domain.Branch{}
	for _, branch := range branches {
		if branch != nil {
			branchByID[branch.ID] = branch
		}
	}
	now := time.Now().UTC()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	result := make([]*domain.CelebrationListItem, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		celebrationType := typeByID[item.CelebrationTypeID]
		row := &domain.CelebrationListItem{Celebration: *item, CelebrationTypeName: "Celebration", IsYearly: true}
		if celebrationType != nil {
			row.CelebrationTypeName = celebrationType.Name
			row.IsYearly = celebrationType.IsYearly
			row.IsUserCelebration = celebrationType.IsUserCelebration
		}
		if item.UserID != nil {
			if employee := employeeByUser[*item.UserID]; employee != nil {
				name := employeeDisplayName(employee.Firstname, employee.MiddleName, employee.Lastname)
				row.EmployeeName = &name
				row.EmployeeCode = employee.EmployeeCode
			}
		}
		if item.BranchID != nil {
			if branch := branchByID[*item.BranchID]; branch != nil {
				row.BranchName = &branch.Name
			}
		}
		if item.CelebrationDate != nil {
			next := nextDashboardCelebrationDate(*item.CelebrationDate, today, row.IsYearly)
			row.NextOccurrenceDate = &next
			days := int(next.Sub(today).Hours() / 24)
			row.DaysUntilNextOccurrence = &days
		}
		result = append(result, row)
	}
	sort.SliceStable(result, func(i, j int) bool {
		left := result[i].NextOccurrenceDate
		right := result[j].NextOccurrenceDate
		if left == nil && right == nil {
			return result[i].CreatedAt.After(result[j].CreatedAt)
		}
		if left == nil {
			return false
		}
		if right == nil {
			return true
		}
		if left.Equal(*right) {
			return result[i].CelebrationTypeName < result[j].CelebrationTypeName
		}
		return left.Before(*right)
	})
	return result, nil
}
