package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateDocumentType(ctx context.Context, item *domain.DocumentType, actorID *uuid.UUID) (*domain.DocumentType, error) {
	row, err := s.getQueries(ctx).CreateDocumentType(ctx, sqlc.CreateDocumentTypeParams{
		TenantID:            item.TenantID,
		Name:                item.Name,
		Description:         textFromPtr(item.Description),
		IsRequired:          item.IsRequired,
		Instructions:        textFromPtr(item.Instructions),
		AllowedContentTypes: item.AllowedContentTypes,
		MaxFileSizeBytes:    item.MaxFileSizeBytes,
		DisplayOrder:        item.DisplayOrder,
		CreatedBy:           uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "create document type", err, tenantIDField(item.TenantID), stringField("document_type_name", item.Name))
	}
	return mapDocumentType(row), nil
}

func (s *Store) ListDocumentTypes(ctx context.Context, tenantID uuid.UUID) ([]*domain.DocumentType, error) {
	rows, err := s.getQueries(ctx).ListDocumentTypes(ctx, tenantID)
	if err != nil {
		return nil, s.logDBError(ctx, "list document types", err, tenantIDField(tenantID))
	}
	items := make([]*domain.DocumentType, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapDocumentType(row))
	}
	return items, nil
}

func (s *Store) GetDocumentType(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.DocumentType, error) {
	row, err := s.getQueries(ctx).GetDocumentType(ctx, sqlc.GetDocumentTypeParams{TenantID: tenantID, ID: id})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrDocumentTypeNotFound
		}
		return nil, s.logDBError(ctx, "get document type", err, tenantIDField(tenantID), stringField("document_type_id", id.String()))
	}
	return mapDocumentType(row), nil
}

func (s *Store) UpdateDocumentType(ctx context.Context, item *domain.DocumentType, actorID *uuid.UUID) (*domain.DocumentType, error) {
	row, err := s.getQueries(ctx).UpdateDocumentType(ctx, sqlc.UpdateDocumentTypeParams{
		TenantID:            item.TenantID,
		ID:                  item.ID,
		Name:                item.Name,
		Description:         textFromPtr(item.Description),
		IsRequired:          item.IsRequired,
		Instructions:        textFromPtr(item.Instructions),
		AllowedContentTypes: item.AllowedContentTypes,
		MaxFileSizeBytes:    item.MaxFileSizeBytes,
		DisplayOrder:        item.DisplayOrder,
		UpdatedBy:           uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "update document type", err, tenantIDField(item.TenantID), stringField("document_type_id", item.ID.String()))
	}
	return mapDocumentType(row), nil
}

func (s *Store) DeleteDocumentType(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteDocumentType(ctx, sqlc.SoftDeleteDocumentTypeParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete document type", err, tenantIDField(tenantID), stringField("document_type_id", id.String()))
	}
	return nil
}

func (s *Store) CreateEmployeeDocument(ctx context.Context, item *domain.EmployeeDocument, actorID *uuid.UUID) (*domain.EmployeeDocument, error) {
	row, err := s.getQueries(ctx).CreateEmployeeDocument(ctx, sqlc.CreateEmployeeDocumentParams{
		TenantID:            item.TenantID,
		UserID:              item.UserID,
		DocumentTypeID:      uuidFromPtr(item.DocumentTypeID),
		Title:               textFromPtr(item.Title),
		FilePath:            textFromPtr(item.FilePath),
		Status:              item.Status,
		OriginalFileName:    textFromPtr(item.OriginalFileName),
		ContentType:         textFromPtr(item.ContentType),
		FileSizeBytes:       int8FromPtr(item.FileSizeBytes),
		Encrypted:           item.Encrypted,
		EncryptionAlgorithm: item.EncryptionAlgorithm,
		CreatedBy:           uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "create employee document", err, tenantIDField(item.TenantID), stringField("user_id", item.UserID.String()))
	}
	return mapEmployeeDocumentRecord(row), nil
}

func (s *Store) GetEmployeeDocument(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.EmployeeDocument, error) {
	row, err := s.getQueries(ctx).GetEmployeeDocument(ctx, sqlc.GetEmployeeDocumentParams{TenantID: tenantID, ID: id})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrEmployeeDocumentNotFound
		}
		return nil, s.logDBError(ctx, "get employee document", err, tenantIDField(tenantID), stringField("document_id", id.String()))
	}
	return mapEmployeeDocumentRecord(row), nil
}

func (s *Store) UpdateEmployeeDocument(ctx context.Context, item *domain.EmployeeDocument, actorID *uuid.UUID) (*domain.EmployeeDocument, error) {
	row, err := s.getQueries(ctx).UpdateEmployeeDocument(ctx, sqlc.UpdateEmployeeDocumentParams{
		TenantID:         item.TenantID,
		ID:               item.ID,
		DocumentTypeID:   uuidFromPtr(item.DocumentTypeID),
		Title:            textFromPtr(item.Title),
		FilePath:         textFromPtr(item.FilePath),
		OriginalFileName: textFromPtr(item.OriginalFileName),
		ContentType:      textFromPtr(item.ContentType),
		FileSizeBytes:    int8FromPtr(item.FileSizeBytes),
		Status:           item.Status,
		UpdatedBy:        uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "update employee document", err, tenantIDField(item.TenantID), stringField("document_id", item.ID.String()))
	}
	return mapEmployeeDocumentRecord(row), nil
}

func (s *Store) DeleteEmployeeDocument(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteEmployeeDocument(ctx, sqlc.SoftDeleteEmployeeDocumentParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete employee document", err, tenantIDField(tenantID), stringField("document_id", id.String()))
	}
	return nil
}

func (s *Store) ReviewEmployeeDocument(ctx context.Context, tenantID uuid.UUID, documentID uuid.UUID, status string, remarks *string, actorID *uuid.UUID) (*domain.EmployeeDocument, error) {
	row, err := s.getQueries(ctx).ReviewEmployeeDocument(ctx, sqlc.ReviewEmployeeDocumentParams{
		TenantID:      tenantID,
		ID:            documentID,
		Status:        status,
		ReviewRemarks: textFromPtr(remarks),
		ReviewedBy:    uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "review employee document", err, tenantIDField(tenantID), stringField("document_id", documentID.String()), stringField("status", status))
	}
	return mapEmployeeDocumentRecord(row), nil
}

func int8FromPtr(value *int64) pgtype.Int8 {
	if value == nil {
		return pgtype.Int8{Valid: false}
	}
	return pgtype.Int8{Int64: *value, Valid: true}
}
