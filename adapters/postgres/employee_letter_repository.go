package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateEmployeeLetterTemplate(ctx context.Context, item *domain.EmployeeLetterTemplate, actorID *uuid.UUID) (*domain.EmployeeLetterTemplate, error) {
	q := s.getQueries(ctx)
	if item.IsDefault {
		if err := q.ClearDefaultEmployeeLetterTemplates(ctx, sqlc.ClearDefaultEmployeeLetterTemplatesParams{TenantID: item.TenantID, LetterType: item.LetterType, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
			return nil, s.logDBError(ctx, "clear default employee letter templates", err, tenantIDField(item.TenantID), stringField("letter_type", item.LetterType))
		}
	}
	row, err := q.CreateEmployeeLetterTemplate(ctx, sqlc.CreateEmployeeLetterTemplateParams{TenantID: item.TenantID, LetterType: item.LetterType, Name: item.Name, Description: textFromPtr(item.Description), Subject: textFromPtr(item.Subject), BodyHtml: item.BodyHTML, FooterHtml: textFromPtr(item.FooterHTML), Locale: item.Locale, IsDefault: item.IsDefault, IsActive: item.IsActive, CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create employee letter template", fmt.Errorf("hrms: create employee letter template: %w", err), tenantIDField(item.TenantID), stringField("letter_type", item.LetterType))
	}
	return mapEmployeeLetterTemplate(row), nil
}

func (s *Store) ListEmployeeLetterTemplates(ctx context.Context, tenantID uuid.UUID, letterType *string) ([]*domain.EmployeeLetterTemplate, error) {
	rows, err := s.getQueries(ctx).ListEmployeeLetterTemplates(ctx, sqlc.ListEmployeeLetterTemplatesParams{TenantID: tenantID, LetterType: textFromPtr(letterType)})
	if err != nil {
		return nil, s.logDBError(ctx, "list employee letter templates", err, tenantIDField(tenantID))
	}
	return mapEmployeeLetterTemplates(rows), nil
}

func (s *Store) GetDefaultEmployeeLetterTemplate(ctx context.Context, tenantID uuid.UUID, letterType string) (*domain.EmployeeLetterTemplate, error) {
	row, err := s.getQueries(ctx).GetDefaultEmployeeLetterTemplate(ctx, sqlc.GetDefaultEmployeeLetterTemplateParams{TenantID: tenantID, LetterType: letterType})
	if err != nil {
		return nil, s.logDBError(ctx, "get default employee letter template", err, tenantIDField(tenantID), stringField("letter_type", letterType))
	}
	return mapEmployeeLetterTemplate(row), nil
}

func (s *Store) GetEmployeeLetterTemplate(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.EmployeeLetterTemplate, error) {
	row, err := s.getQueries(ctx).GetEmployeeLetterTemplate(ctx, sqlc.GetEmployeeLetterTemplateParams{TenantID: tenantID, ID: id})
	if err != nil {
		return nil, s.logDBError(ctx, "get employee letter template", err, tenantIDField(tenantID), stringField("employee_letter_template_id", id.String()))
	}
	return mapEmployeeLetterTemplate(row), nil
}

func (s *Store) UpdateEmployeeLetterTemplate(ctx context.Context, item *domain.EmployeeLetterTemplate, actorID *uuid.UUID) (*domain.EmployeeLetterTemplate, error) {
	q := s.getQueries(ctx)
	if item.IsDefault {
		if err := q.ClearDefaultEmployeeLetterTemplates(ctx, sqlc.ClearDefaultEmployeeLetterTemplatesParams{TenantID: item.TenantID, LetterType: item.LetterType, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
			return nil, s.logDBError(ctx, "clear default employee letter templates", err, tenantIDField(item.TenantID), stringField("letter_type", item.LetterType))
		}
	}
	row, err := q.UpdateEmployeeLetterTemplate(ctx, sqlc.UpdateEmployeeLetterTemplateParams{TenantID: item.TenantID, ID: item.ID, LetterType: item.LetterType, Name: item.Name, Description: textFromPtr(item.Description), Subject: textFromPtr(item.Subject), BodyHtml: item.BodyHTML, FooterHtml: textFromPtr(item.FooterHTML), Locale: item.Locale, IsDefault: item.IsDefault, IsActive: item.IsActive, UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "update employee letter template", err, tenantIDField(item.TenantID), stringField("employee_letter_template_id", item.ID.String()))
	}
	return mapEmployeeLetterTemplate(row), nil
}

func (s *Store) DeleteEmployeeLetterTemplate(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteEmployeeLetterTemplate(ctx, sqlc.SoftDeleteEmployeeLetterTemplateParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete employee letter template", err, tenantIDField(tenantID), stringField("employee_letter_template_id", id.String()))
	}
	return nil
}

func (s *Store) CreateEmployeeLetter(ctx context.Context, item *domain.EmployeeLetter, actorID *uuid.UUID) (*domain.EmployeeLetter, error) {
	q := s.getQueries(ctx)
	if err := q.ClearLatestEmployeeLetters(ctx, sqlc.ClearLatestEmployeeLettersParams{TenantID: item.TenantID, EmployeeID: item.EmployeeID, LetterType: item.LetterType, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return nil, s.logDBError(ctx, "clear latest employee letters", err, tenantIDField(item.TenantID), stringField("employee_id", item.EmployeeID.String()), stringField("letter_type", item.LetterType))
	}
	version, err := q.NextEmployeeLetterVersion(ctx, sqlc.NextEmployeeLetterVersionParams{TenantID: item.TenantID, EmployeeID: item.EmployeeID, LetterType: item.LetterType})
	if err != nil {
		return nil, s.logDBError(ctx, "next employee letter version", err, tenantIDField(item.TenantID), stringField("employee_id", item.EmployeeID.String()), stringField("letter_type", item.LetterType))
	}
	row, err := q.CreateEmployeeLetter(ctx, sqlc.CreateEmployeeLetterParams{TenantID: item.TenantID, EmployeeID: item.EmployeeID, UserID: item.UserID, TemplateID: uuidFromPtr(item.TemplateID), DocumentTypeID: uuidFromPtr(item.DocumentTypeID), EmployeeDocumentID: uuidFromPtr(item.EmployeeDocumentID), LetterType: item.LetterType, Subject: textFromPtr(item.Subject), RenderedHtml: textFromPtr(item.RenderedHTML), Status: item.Status, IssueDate: dateFromPtr(item.IssueDate), EffectiveDate: dateFromPtr(item.EffectiveDate), EndDate: dateFromPtr(item.EndDate), PdfPath: textFromPtr(item.PDFPath), Version: version, IsLatest: true, ApprovalRequestedAt: timestamptzFromPtr(item.ApprovalRequestedAt), SignatureToken: textFromPtr(item.SignatureToken), SignatureRequestedAt: timestamptzFromPtr(item.SignatureRequestedAt), SignerEmail: textFromPtr(item.SignerEmail), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create employee letter", fmt.Errorf("hrms: create employee letter: %w", err), tenantIDField(item.TenantID), stringField("employee_id", item.EmployeeID.String()), stringField("letter_type", item.LetterType))
	}
	return mapEmployeeLetter(row), nil
}

func (s *Store) ListEmployeeLetters(ctx context.Context, filter domain.EmployeeLetterFilter) ([]*domain.EmployeeLetter, error) {
	rows, err := s.getQueries(ctx).ListEmployeeLetters(ctx, sqlc.ListEmployeeLettersParams{TenantID: filter.TenantID, EmployeeID: uuidFromPtr(filter.EmployeeID), LetterType: textFromPtr(filter.LetterType), Status: textFromPtr(filter.Status), Search: textFromPtr(filter.Search), Limit: filter.Limit, Offset: filter.Offset})
	if err != nil {
		return nil, s.logDBError(ctx, "list employee letters", err, tenantIDField(filter.TenantID))
	}
	return mapEmployeeLetters(rows), nil
}

func (s *Store) CountEmployeeLetters(ctx context.Context, filter domain.EmployeeLetterFilter) (int64, error) {
	count, err := s.getQueries(ctx).CountEmployeeLetters(ctx, sqlc.CountEmployeeLettersParams{TenantID: filter.TenantID, EmployeeID: uuidFromPtr(filter.EmployeeID), LetterType: textFromPtr(filter.LetterType), Status: textFromPtr(filter.Status), Search: textFromPtr(filter.Search)})
	if err != nil {
		return 0, s.logDBError(ctx, "count employee letters", err, tenantIDField(filter.TenantID))
	}
	return count, nil
}

func (s *Store) GetEmployeeLetter(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.EmployeeLetter, error) {
	row, err := s.getQueries(ctx).GetEmployeeLetter(ctx, sqlc.GetEmployeeLetterParams{TenantID: tenantID, ID: id})
	if err != nil {
		return nil, s.logDBError(ctx, "get employee letter", err, tenantIDField(tenantID), stringField("employee_letter_id", id.String()))
	}
	return mapEmployeeLetter(row), nil
}

func (s *Store) GetEmployeeLetterBySignatureToken(ctx context.Context, token string) (*domain.EmployeeLetter, error) {
	row, err := s.getQueries(ctx).GetEmployeeLetterBySignatureToken(ctx, textFromPtr(&token))
	if err != nil {
		return nil, s.logDBError(ctx, "get employee letter by signature token", err)
	}
	return mapEmployeeLetter(row), nil
}

func (s *Store) UpdateEmployeeLetterPDF(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, pdfPath *string, employeeDocumentID *uuid.UUID, actorID *uuid.UUID) (*domain.EmployeeLetter, error) {
	row, err := s.getQueries(ctx).UpdateEmployeeLetterPDF(ctx, sqlc.UpdateEmployeeLetterPDFParams{TenantID: tenantID, ID: id, PdfPath: textFromPtr(pdfPath), EmployeeDocumentID: uuidFromPtr(employeeDocumentID), UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "update employee letter pdf", err, tenantIDField(tenantID), stringField("employee_letter_id", id.String()))
	}
	return mapEmployeeLetter(row), nil
}

func (s *Store) UpdateEmployeeLetterStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, remarks *string, actorID *uuid.UUID) (*domain.EmployeeLetter, error) {
	row, err := s.getQueries(ctx).UpdateEmployeeLetterStatus(ctx, sqlc.UpdateEmployeeLetterStatusParams{TenantID: tenantID, ID: id, Status: status, UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "update employee letter status", err, tenantIDField(tenantID), stringField("employee_letter_id", id.String()), stringField("status", status))
	}
	return mapEmployeeLetter(row), nil
}

func (s *Store) SignEmployeeLetter(ctx context.Context, token string, signerName string, signerEmail string, signerIP *string, userAgent *string, signatureHash string) (*domain.EmployeeLetter, error) {
	row, err := s.getQueries(ctx).SignEmployeeLetter(ctx, sqlc.SignEmployeeLetterParams{SignatureToken: textFromPtr(&token), SignerName: textFromPtr(&signerName), SignerEmail: textFromPtr(&signerEmail), SignerIp: addrFromStringPtr(signerIP), SignerUserAgent: textFromPtr(userAgent), SignatureHash: textFromPtr(&signatureHash)})
	if err != nil {
		return nil, s.logDBError(ctx, "sign employee letter", err)
	}
	return mapEmployeeLetter(row), nil
}

func (s *Store) DeleteEmployeeLetter(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteEmployeeLetter(ctx, sqlc.SoftDeleteEmployeeLetterParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete employee letter", err, tenantIDField(tenantID), stringField("employee_letter_id", id.String()))
	}
	return nil
}

func (s *Store) CreateEmployeeLetterEvent(ctx context.Context, event *domain.EmployeeLetterEvent, actorID *uuid.UUID) (*domain.EmployeeLetterEvent, error) {
	row, err := s.getQueries(ctx).CreateEmployeeLetterEvent(ctx, sqlc.CreateEmployeeLetterEventParams{TenantID: event.TenantID, EmployeeLetterID: event.EmployeeLetterID, FromStatus: textFromPtr(event.FromStatus), ToStatus: event.ToStatus, Action: event.Action, Remarks: textFromPtr(event.Remarks), ActorEmail: textFromPtr(event.ActorEmail), IpAddress: addrFromStringPtr(event.IPAddress), UserAgent: textFromPtr(event.UserAgent), Metadata: jsonBytesFromAnyMap(event.Metadata), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create employee letter event", err, tenantIDField(event.TenantID), stringField("employee_letter_id", event.EmployeeLetterID.String()))
	}
	return mapEmployeeLetterEvent(row), nil
}

func (s *Store) ListEmployeeLetterEvents(ctx context.Context, tenantID uuid.UUID, employeeLetterID uuid.UUID) ([]*domain.EmployeeLetterEvent, error) {
	rows, err := s.getQueries(ctx).ListEmployeeLetterEvents(ctx, sqlc.ListEmployeeLetterEventsParams{TenantID: tenantID, EmployeeLetterID: employeeLetterID})
	if err != nil {
		return nil, s.logDBError(ctx, "list employee letter events", err, tenantIDField(tenantID), stringField("employee_letter_id", employeeLetterID.String()))
	}
	return mapEmployeeLetterEvents(rows), nil
}
