package services

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
	"github.com/ranakdinesh/spur-hrms/pkg/pbkdf2verifier"
)

func (s *TenantService) VerifyLegacyPBKDF2Password(password string, encodedHash string) (bool, error) {
	ok, err := pbkdf2verifier.Verify(password, encodedHash)
	if err != nil {
		s.logError("verify legacy pbkdf2 password", err)
		return false, err
	}
	return ok, nil
}

func (s *TenantService) MigrateLegacyPassword(ctx context.Context, cmd ports.LegacyPasswordMigrationCommand) (*ports.LegacyPasswordMigrationResult, error) {
	if cmd.TenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate legacy password migration tenant", err)
		return nil, err
	}
	if strings.TrimSpace(cmd.Identifier) == "" {
		err := domain.ErrInvalidLegacyPasswordIdentifier
		s.logError("validate legacy password migration identifier", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	if strings.TrimSpace(cmd.Password) == "" {
		err := domain.ErrInvalidLegacyPassword
		s.logError("validate legacy password migration password", err, serviceTenantIDField(cmd.TenantID), serviceStringField("identifier", cmd.Identifier))
		return nil, err
	}
	if s.legacyPasswordMigration == nil {
		err := domain.ErrLegacyPasswordMigrationPortMissing
		s.logError("legacy password migration port missing", err, serviceTenantIDField(cmd.TenantID), serviceStringField("identifier", cmd.Identifier))
		return nil, err
	}
	result, err := s.legacyPasswordMigration.VerifyAndMigrateLegacyPassword(ctx, cmd)
	if err != nil {
		s.logError("migrate legacy password", err, serviceTenantIDField(cmd.TenantID), serviceStringField("identifier", cmd.Identifier))
		return nil, err
	}
	return result, nil
}
