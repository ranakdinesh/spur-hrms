package services

import (
	"context"
	"encoding/base64"
	"strings"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

const (
	defaultEmployeeDocumentAllowedContentTypes = "application/pdf,image/jpeg,image/png,image/webp,application/msword,application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	defaultEmployeeDocumentMaxFileSizeBytes    = int64(10 * 1024 * 1024)
)

func (s *TenantService) CreateDocumentType(ctx context.Context, cmd ports.DocumentTypeCommand) (*domain.DocumentType, error) {
	name := strings.TrimSpace(cmd.Name)
	if cmd.TenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate document type tenant", err)
		return nil, err
	}
	if name == "" {
		err := domain.ErrInvalidDocumentTypeName
		s.logError("validate document type name", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	item := &domain.DocumentType{
		TenantID:            cmd.TenantID,
		Name:                name,
		Description:         cleanCommandString(cmd.Description),
		IsRequired:          boolValue(cmd.IsRequired, true),
		Instructions:        cleanCommandString(cmd.Instructions),
		AllowedContentTypes: cleanContentTypes(cmd.AllowedContentTypes),
		MaxFileSizeBytes:    maxFileSizeOrDefault(cmd.MaxFileSizeBytes),
		DisplayOrder:        cmd.DisplayOrder,
	}
	result, err := s.employeeDocuments.CreateDocumentType(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create document type", err, serviceTenantIDField(cmd.TenantID), serviceStringField("document_type_name", name))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListDocumentTypes(ctx context.Context, tenantID uuid.UUID) ([]*domain.DocumentType, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate document type list tenant", err)
		return nil, err
	}
	items, err := s.employeeDocuments.ListDocumentTypes(ctx, tenantID)
	if err != nil {
		s.logError("list document types", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) GetDocumentType(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.DocumentType, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate document type get tenant", err)
		return nil, err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidDocumentTypeID
		s.logError("validate document type get id", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	item, err := s.employeeDocuments.GetDocumentType(ctx, tenantID, id)
	if err != nil {
		s.logError("get document type", err, serviceTenantIDField(tenantID), serviceStringField("document_type_id", id.String()))
		return nil, err
	}
	return item, nil
}

func (s *TenantService) UpdateDocumentType(ctx context.Context, cmd ports.DocumentTypeCommand) (*domain.DocumentType, error) {
	if cmd.ID == uuid.Nil {
		err := domain.ErrInvalidDocumentTypeID
		s.logError("validate document type update id", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	name := strings.TrimSpace(cmd.Name)
	if name == "" {
		err := domain.ErrInvalidDocumentTypeName
		s.logError("validate document type update name", err, serviceTenantIDField(cmd.TenantID), serviceStringField("document_type_id", cmd.ID.String()))
		return nil, err
	}
	existing, err := s.employeeDocuments.GetDocumentType(ctx, cmd.TenantID, cmd.ID)
	if err != nil {
		s.logError("get document type before update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("document_type_id", cmd.ID.String()))
		return nil, err
	}
	item := &domain.DocumentType{
		ID:                  cmd.ID,
		TenantID:            cmd.TenantID,
		Name:                name,
		Description:         cleanCommandString(cmd.Description),
		IsRequired:          boolValue(cmd.IsRequired, existing.IsRequired),
		Instructions:        cleanCommandString(cmd.Instructions),
		AllowedContentTypes: firstString(strings.TrimSpace(cmd.AllowedContentTypes), existing.AllowedContentTypes, defaultEmployeeDocumentAllowedContentTypes),
		MaxFileSizeBytes:    maxFileSizeOrExisting(cmd.MaxFileSizeBytes, existing.MaxFileSizeBytes),
		DisplayOrder:        cmd.DisplayOrder,
	}
	result, err := s.employeeDocuments.UpdateDocumentType(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update document type", err, serviceTenantIDField(cmd.TenantID), serviceStringField("document_type_id", cmd.ID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) DeleteDocumentType(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate document type delete tenant", err)
		return err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidDocumentTypeID
		s.logError("validate document type delete id", err, serviceTenantIDField(tenantID))
		return err
	}
	if err := s.employeeDocuments.DeleteDocumentType(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete document type", err, serviceTenantIDField(tenantID), serviceStringField("document_type_id", id.String()))
		return err
	}
	return nil
}

func (s *TenantService) CreateEmployeeDocument(ctx context.Context, cmd ports.EmployeeDocumentCommand) (*domain.EmployeeDocument, error) {
	profile, err := s.validateEmployeeDocumentCommand(ctx, cmd, false)
	if err != nil {
		return nil, err
	}
	documentType, err := s.documentTypeForCommand(ctx, cmd)
	if err != nil {
		return nil, err
	}
	if employeeHasApprovedDocumentType(profile, cmd.DocumentTypeID, uuid.Nil) {
		err := domain.ErrEmployeeDocumentApprovedLocked
		s.logError("validate approved employee document create lock", err, serviceTenantIDField(cmd.TenantID), serviceStringField("employee_id", cmd.EmployeeID.String()))
		return nil, err
	}
	fileContent, err := decodeEmployeeDocumentFile(cmd.FileContentBase64)
	if err != nil {
		s.logError("decode employee document file", err, serviceTenantIDField(cmd.TenantID), serviceStringField("employee_id", cmd.EmployeeID.String()))
		return nil, err
	}
	if len(fileContent) > 0 && s.documentStorage == nil {
		err := domain.ErrEmployeeDocumentStorageMissing
		s.logError("store employee document file", err, serviceTenantIDField(cmd.TenantID), serviceStringField("employee_id", cmd.EmployeeID.String()))
		return nil, err
	}
	if err := s.validateEmployeeDocumentFile(cmd, documentType, fileContent); err != nil {
		s.logError("validate employee document file", err, serviceTenantIDField(cmd.TenantID), serviceStringField("employee_id", cmd.EmployeeID.String()))
		return nil, err
	}
	if isBlankPtr(cmd.Title) && isBlankPtr(cmd.FilePath) && len(fileContent) == 0 {
		err := domain.ErrInvalidEmployeeDocument
		s.logError("validate employee document content", err, serviceTenantIDField(cmd.TenantID), serviceStringField("employee_id", cmd.EmployeeID.String()))
		return nil, err
	}
	item := &domain.EmployeeDocument{
		TenantID:            cmd.TenantID,
		UserID:              profile.Employee.UserID,
		DocumentTypeID:      nilIfZeroUUID(cmd.DocumentTypeID),
		Title:               cleanCommandString(cmd.Title),
		FilePath:            cleanCommandString(cmd.FilePath),
		Status:              domain.EmployeeDocumentStatusPendingReview,
		OriginalFileName:    cleanFileName(cmd.FileName),
		ContentType:         cleanFileName(cmd.FileContentType),
		FileSizeBytes:       fileSizePtr(fileContent),
		Encrypted:           true,
		EncryptionAlgorithm: "AES-256-GCM",
	}
	result, err := s.employeeDocuments.CreateEmployeeDocument(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create employee document", err, serviceTenantIDField(cmd.TenantID), serviceStringField("employee_id", cmd.EmployeeID.String()))
		return nil, err
	}
	if len(fileContent) > 0 {
		path, err := s.documentStorage.StoreEmployeeDocument(ctx, ports.StoreEmployeeDocumentInput{TenantID: cmd.TenantID, EmployeeID: cmd.EmployeeID, DocumentID: result.ID, FileName: cmd.FileName, ContentType: cmd.FileContentType, Content: fileContent})
		if err != nil {
			s.logError("store employee document file", err, serviceTenantIDField(cmd.TenantID), serviceStringField("document_id", result.ID.String()))
			return nil, err
		}
		result.FilePath = &path
		result, err = s.employeeDocuments.UpdateEmployeeDocument(ctx, result, cmd.ActorID)
		if err != nil {
			s.logError("attach employee document file", err, serviceTenantIDField(cmd.TenantID), serviceStringField("document_id", result.ID.String()))
			return nil, err
		}
	}
	return result, nil
}

func (s *TenantService) UpdateEmployeeDocument(ctx context.Context, cmd ports.EmployeeDocumentCommand) (*domain.EmployeeDocument, error) {
	if cmd.ID == uuid.Nil {
		err := domain.ErrInvalidEmployeeDocumentID
		s.logError("validate employee document update id", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	profile, err := s.validateEmployeeDocumentCommand(ctx, cmd, true)
	if err != nil {
		return nil, err
	}
	existing, err := s.employeeDocuments.GetEmployeeDocument(ctx, cmd.TenantID, cmd.ID)
	if err != nil {
		s.logError("get employee document before update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("document_id", cmd.ID.String()))
		return nil, err
	}
	if existing.UserID != profile.Employee.UserID {
		err := domain.ErrEmployeeDocumentNotFound
		s.logError("validate employee document owner", err, serviceTenantIDField(cmd.TenantID), serviceStringField("document_id", cmd.ID.String()), serviceStringField("employee_id", cmd.EmployeeID.String()))
		return nil, err
	}
	if existing.Status == domain.EmployeeDocumentStatusApproved || employeeHasApprovedDocumentType(profile, cmd.DocumentTypeID, cmd.ID) {
		err := domain.ErrEmployeeDocumentApprovedLocked
		s.logError("validate approved employee document update lock", err, serviceTenantIDField(cmd.TenantID), serviceStringField("document_id", cmd.ID.String()), serviceStringField("employee_id", cmd.EmployeeID.String()))
		return nil, err
	}
	fileContent, err := decodeEmployeeDocumentFile(cmd.FileContentBase64)
	if err != nil {
		s.logError("decode employee document update file", err, serviceTenantIDField(cmd.TenantID), serviceStringField("document_id", cmd.ID.String()))
		return nil, err
	}
	if len(fileContent) > 0 && s.documentStorage == nil {
		err := domain.ErrEmployeeDocumentStorageMissing
		s.logError("store employee document update file", err, serviceTenantIDField(cmd.TenantID), serviceStringField("document_id", cmd.ID.String()))
		return nil, err
	}
	documentType, err := s.documentTypeForCommand(ctx, cmd)
	if err != nil {
		return nil, err
	}
	if err := s.validateEmployeeDocumentFile(cmd, documentType, fileContent); err != nil {
		s.logError("validate employee document update file", err, serviceTenantIDField(cmd.TenantID), serviceStringField("document_id", cmd.ID.String()))
		return nil, err
	}
	item := &domain.EmployeeDocument{
		ID:                  cmd.ID,
		TenantID:            cmd.TenantID,
		UserID:              profile.Employee.UserID,
		DocumentTypeID:      nilIfZeroUUID(cmd.DocumentTypeID),
		Title:               cleanCommandString(cmd.Title),
		FilePath:            firstCleanString(cmd.FilePath, existing.FilePath),
		OriginalFileName:    firstCleanString(cleanFileName(cmd.FileName), existing.OriginalFileName),
		ContentType:         firstCleanString(cleanFileName(cmd.FileContentType), existing.ContentType),
		FileSizeBytes:       firstFileSize(fileSizePtr(fileContent), existing.FileSizeBytes),
		Status:              domain.EmployeeDocumentStatusPendingReview,
		Encrypted:           true,
		EncryptionAlgorithm: "AES-256-GCM",
	}
	if len(fileContent) > 0 {
		path, err := s.documentStorage.StoreEmployeeDocument(ctx, ports.StoreEmployeeDocumentInput{TenantID: cmd.TenantID, EmployeeID: cmd.EmployeeID, DocumentID: cmd.ID, FileName: cmd.FileName, ContentType: cmd.FileContentType, Content: fileContent})
		if err != nil {
			s.logError("store employee document update file", err, serviceTenantIDField(cmd.TenantID), serviceStringField("document_id", cmd.ID.String()))
			return nil, err
		}
		item.FilePath = &path
	}
	result, err := s.employeeDocuments.UpdateEmployeeDocument(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update employee document", err, serviceTenantIDField(cmd.TenantID), serviceStringField("document_id", cmd.ID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ReviewEmployeeDocument(ctx context.Context, cmd ports.EmployeeDocumentReviewCommand) (*domain.EmployeeDocument, error) {
	if cmd.DocumentID == uuid.Nil {
		err := domain.ErrInvalidEmployeeDocumentID
		s.logError("validate employee document review id", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	status, err := domain.ValidateEmployeeDocumentReviewStatus(cmd.Status)
	if err != nil {
		s.logError("validate employee document review status", err, serviceTenantIDField(cmd.TenantID), serviceStringField("document_id", cmd.DocumentID.String()), serviceStringField("status", cmd.Status))
		return nil, err
	}
	profile, err := s.validateEmployeeDocumentCommand(ctx, ports.EmployeeDocumentCommand{TenantID: cmd.TenantID, EmployeeID: cmd.EmployeeID, ID: cmd.DocumentID}, true)
	if err != nil {
		return nil, err
	}
	existing, err := s.employeeDocuments.GetEmployeeDocument(ctx, cmd.TenantID, cmd.DocumentID)
	if err != nil {
		s.logError("get employee document before review", err, serviceTenantIDField(cmd.TenantID), serviceStringField("document_id", cmd.DocumentID.String()))
		return nil, err
	}
	if existing.UserID != profile.Employee.UserID {
		err := domain.ErrEmployeeDocumentNotFound
		s.logError("validate employee document review owner", err, serviceTenantIDField(cmd.TenantID), serviceStringField("document_id", cmd.DocumentID.String()), serviceStringField("employee_id", cmd.EmployeeID.String()))
		return nil, err
	}
	result, err := s.employeeDocuments.ReviewEmployeeDocument(ctx, cmd.TenantID, cmd.DocumentID, status, cleanCommandString(cmd.Remarks), cmd.ActorID)
	if err != nil {
		s.logError("review employee document", err, serviceTenantIDField(cmd.TenantID), serviceStringField("document_id", cmd.DocumentID.String()), serviceStringField("status", status))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) DeleteEmployeeDocument(ctx context.Context, tenantID uuid.UUID, employeeID uuid.UUID, documentID uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate employee document delete tenant", err)
		return err
	}
	if employeeID == uuid.Nil {
		err := domain.ErrInvalidEmployeeID
		s.logError("validate employee document delete employee", err, serviceTenantIDField(tenantID))
		return err
	}
	if documentID == uuid.Nil {
		err := domain.ErrInvalidEmployeeDocumentID
		s.logError("validate employee document delete id", err, serviceTenantIDField(tenantID), serviceStringField("employee_id", employeeID.String()))
		return err
	}
	profile, err := s.employees.GetEmployeeProfile(ctx, tenantID, employeeID)
	if err != nil {
		s.logError("get employee before document delete", err, serviceTenantIDField(tenantID), serviceStringField("employee_id", employeeID.String()))
		return err
	}
	existing, err := s.employeeDocuments.GetEmployeeDocument(ctx, tenantID, documentID)
	if err != nil {
		s.logError("get employee document before delete", err, serviceTenantIDField(tenantID), serviceStringField("document_id", documentID.String()))
		return err
	}
	if existing.UserID != profile.Employee.UserID {
		err := domain.ErrEmployeeDocumentNotFound
		s.logError("validate employee document delete owner", err, serviceTenantIDField(tenantID), serviceStringField("document_id", documentID.String()), serviceStringField("employee_id", employeeID.String()))
		return err
	}
	if err := s.employeeDocuments.DeleteEmployeeDocument(ctx, tenantID, documentID, actorID); err != nil {
		s.logError("delete employee document", err, serviceTenantIDField(tenantID), serviceStringField("document_id", documentID.String()))
		return err
	}
	return nil
}

func (s *TenantService) validateEmployeeDocumentCommand(ctx context.Context, cmd ports.EmployeeDocumentCommand, requireDocumentID bool) (*domain.EmployeeProfile, error) {
	if cmd.TenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate employee document tenant", err)
		return nil, err
	}
	if cmd.EmployeeID == uuid.Nil {
		err := domain.ErrInvalidEmployeeID
		s.logError("validate employee document employee", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	if requireDocumentID && cmd.ID == uuid.Nil {
		err := domain.ErrInvalidEmployeeDocumentID
		s.logError("validate employee document id", err, serviceTenantIDField(cmd.TenantID), serviceStringField("employee_id", cmd.EmployeeID.String()))
		return nil, err
	}
	if cmd.DocumentTypeID != nil && *cmd.DocumentTypeID != uuid.Nil {
		if _, err := s.employeeDocuments.GetDocumentType(ctx, cmd.TenantID, *cmd.DocumentTypeID); err != nil {
			s.logError("validate employee document type", err, serviceTenantIDField(cmd.TenantID), serviceStringField("document_type_id", cmd.DocumentTypeID.String()))
			return nil, err
		}
	}
	profile, err := s.employees.GetEmployeeProfile(ctx, cmd.TenantID, cmd.EmployeeID)
	if err != nil {
		s.logError("get employee before document change", err, serviceTenantIDField(cmd.TenantID), serviceStringField("employee_id", cmd.EmployeeID.String()))
		return nil, err
	}
	return profile, nil
}

func decodeEmployeeDocumentFile(value string) ([]byte, error) {
	clean := strings.TrimSpace(value)
	if clean == "" {
		return nil, nil
	}
	if comma := strings.Index(clean, ","); comma >= 0 {
		clean = clean[comma+1:]
	}
	return base64.StdEncoding.DecodeString(clean)
}

func nilIfZeroUUID(value *uuid.UUID) *uuid.UUID {
	if value == nil || *value == uuid.Nil {
		return nil
	}
	return value
}

func isBlankPtr(value *string) bool {
	return value == nil || strings.TrimSpace(*value) == ""
}

func firstCleanString(preferred *string, fallback *string) *string {
	if clean := cleanCommandString(preferred); clean != nil {
		return clean
	}
	return cleanCommandString(fallback)
}

func cleanFileName(value string) *string {
	clean := strings.TrimSpace(value)
	if clean == "" {
		return nil
	}
	return &clean
}

func fileSizePtr(content []byte) *int64 {
	if len(content) == 0 {
		return nil
	}
	size := int64(len(content))
	return &size
}

func firstFileSize(preferred *int64, fallback *int64) *int64 {
	if preferred != nil {
		return preferred
	}
	return fallback
}

func boolValue(value *bool, fallback bool) bool {
	if value == nil {
		return fallback
	}
	return *value
}

func cleanContentTypes(value string) string {
	return firstString(strings.TrimSpace(value), defaultEmployeeDocumentAllowedContentTypes)
}

func maxFileSizeOrDefault(value int64) int64 {
	if value <= 0 {
		return defaultEmployeeDocumentMaxFileSizeBytes
	}
	return value
}

func maxFileSizeOrExisting(value int64, existing int64) int64 {
	if value > 0 {
		return value
	}
	if existing > 0 {
		return existing
	}
	return defaultEmployeeDocumentMaxFileSizeBytes
}

func firstString(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func (s *TenantService) documentTypeForCommand(ctx context.Context, cmd ports.EmployeeDocumentCommand) (*domain.DocumentType, error) {
	if cmd.DocumentTypeID == nil || *cmd.DocumentTypeID == uuid.Nil {
		return nil, nil
	}
	item, err := s.employeeDocuments.GetDocumentType(ctx, cmd.TenantID, *cmd.DocumentTypeID)
	if err != nil {
		s.logError("get employee document type for file validation", err, serviceTenantIDField(cmd.TenantID), serviceStringField("document_type_id", cmd.DocumentTypeID.String()))
		return nil, err
	}
	return item, nil
}

func (s *TenantService) validateEmployeeDocumentFile(cmd ports.EmployeeDocumentCommand, documentType *domain.DocumentType, content []byte) error {
	if len(content) == 0 {
		return nil
	}
	limit := defaultEmployeeDocumentMaxFileSizeBytes
	allowed := defaultEmployeeDocumentAllowedContentTypes
	if documentType != nil {
		if documentType.MaxFileSizeBytes > 0 {
			limit = documentType.MaxFileSizeBytes
		}
		if strings.TrimSpace(documentType.AllowedContentTypes) != "" {
			allowed = documentType.AllowedContentTypes
		}
	}
	if int64(len(content)) > limit {
		return domain.ErrEmployeeDocumentFileTooLarge
	}
	contentType := strings.ToLower(strings.TrimSpace(cmd.FileContentType))
	if contentType == "" {
		return domain.ErrEmployeeDocumentFileTypeDenied
	}
	for _, item := range strings.Split(allowed, ",") {
		if strings.ToLower(strings.TrimSpace(item)) == contentType {
			if hasExpectedFileSignature(contentType, content) {
				return nil
			}
			return domain.ErrEmployeeDocumentFileTypeDenied
		}
	}
	return domain.ErrEmployeeDocumentFileTypeDenied
}

func employeeHasApprovedDocumentType(profile *domain.EmployeeProfile, documentTypeID *uuid.UUID, exceptDocumentID uuid.UUID) bool {
	if profile == nil || documentTypeID == nil || *documentTypeID == uuid.Nil {
		return false
	}
	for _, document := range profile.Documents {
		if document == nil || document.DocumentTypeID == nil || *document.DocumentTypeID != *documentTypeID {
			continue
		}
		if exceptDocumentID != uuid.Nil && document.ID == exceptDocumentID {
			continue
		}
		if document.Status == domain.EmployeeDocumentStatusApproved {
			return true
		}
	}
	return false
}

func hasExpectedFileSignature(contentType string, content []byte) bool {
	if len(content) == 0 {
		return false
	}
	switch contentType {
	case "application/pdf":
		return hasPrefix(content, []byte("%PDF"))
	case "image/jpeg":
		return len(content) >= 3 && content[0] == 0xff && content[1] == 0xd8 && content[2] == 0xff
	case "image/png":
		return hasPrefix(content, []byte{0x89, 'P', 'N', 'G', 0x0d, 0x0a, 0x1a, 0x0a})
	case "image/webp":
		return len(content) >= 12 && string(content[0:4]) == "RIFF" && string(content[8:12]) == "WEBP"
	case "application/msword":
		return hasPrefix(content, []byte{0xd0, 0xcf, 0x11, 0xe0, 0xa1, 0xb1, 0x1a, 0xe1})
	case "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
		return hasPrefix(content, []byte{'P', 'K', 0x03, 0x04})
	default:
		return true
	}
}

func hasPrefix(content []byte, prefix []byte) bool {
	if len(content) < len(prefix) {
		return false
	}
	for i, value := range prefix {
		if content[i] != value {
			return false
		}
	}
	return true
}
