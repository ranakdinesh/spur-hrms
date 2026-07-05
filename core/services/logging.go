package services

import (
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type serviceLogField func(*zerolog.Event)

func serviceTenantIDField(tenantID uuid.UUID) serviceLogField {
	return func(event *zerolog.Event) {
		if tenantID != uuid.Nil {
			event.Str("tenant_id", tenantID.String())
		}
	}
}

func serviceStringField(key string, value string) serviceLogField {
	return func(event *zerolog.Event) {
		if value != "" {
			event.Str(key, value)
		}
	}
}

func (s *TenantService) logError(operation string, err error, fields ...serviceLogField) {
	if s != nil && s.log != nil && err != nil {
		event := s.log.Error().Err(err).Str("operation", operation)
		for _, field := range fields {
			if field != nil {
				field(event)
			}
		}
		event.Msg("hrms service error")
	}
}
