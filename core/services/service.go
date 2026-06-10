package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-platform/logger"
	"y/core/domain"
	"y/core/ports"
)

type HrmsService struct {
	repo ports.HrmsRepo
	log  *logger.Loggerx
}

func NewHrmsService(repo ports.HrmsRepo, log *logger.Loggerx) *HrmsService {
	return &HrmsService{repo: repo, log: log}
}

func (s *HrmsService) Create(ctx context.Context, cmd ports.CreateHrmsCmd) (*domain.Hrms, error) {
	entity, err := domain.NewHrms(cmd.TenantID)
	if err != nil {
		return nil, err
	}
	result, err := s.repo.Create(ctx, entity)
	if err != nil {
		return nil, err
	}
	s.log.Info(ctx).Str("id", result.ID.String()).Msg("hrms: created")
	return result, nil
}

func (s *HrmsService) Get(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) (*domain.Hrms, error) {
	return s.repo.GetByID(ctx, id, tenantID)
}

func (s *HrmsService) List(ctx context.Context, tenantID uuid.UUID) ([]*domain.Hrms, error) {
	return s.repo.List(ctx, tenantID)
}

func (s *HrmsService) Delete(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) error {
	if err := s.repo.Delete(ctx, id, tenantID); err != nil {
		return err
	}
	s.log.Info(ctx).Str("id", id.String()).Msg("hrms: deleted")
	return nil
}
