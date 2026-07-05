package domain

import "errors"

var (
	ErrLegacyPasswordMigrationPortMissing = errors.New("legacy password migration port is not configured")
	ErrInvalidLegacyPasswordIdentifier    = errors.New("legacy password identifier is required")
	ErrInvalidLegacyPassword              = errors.New("legacy password is required")
)
