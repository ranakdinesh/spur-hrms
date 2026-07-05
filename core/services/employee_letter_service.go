package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) CreateEmployeeLetterTemplate(ctx context.Context, cmd ports.EmployeeLetterTemplateCommand) (*domain.EmployeeLetterTemplate, error) {
	item, err := domain.NewEmployeeLetterTemplate(employeeLetterTemplateInput(cmd))
	if err != nil {
		s.logError("validate employee letter template create", err, serviceTenantIDField(cmd.TenantID), serviceStringField("letter_type", cmd.LetterType))
		return nil, err
	}
	result, err := s.employeeLetters.CreateEmployeeLetterTemplate(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create employee letter template", err, serviceTenantIDField(cmd.TenantID), serviceStringField("letter_type", item.LetterType))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListEmployeeLetterTemplates(ctx context.Context, tenantID uuid.UUID, letterType *string) ([]*domain.EmployeeLetterTemplate, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate employee letter template list tenant", err)
		return nil, err
	}
	cleanType := cleanStringPtr(letterType)
	if cleanType != nil {
		normalized, err := domain.ValidateEmployeeLetterType(*cleanType)
		if err != nil {
			s.logError("validate employee letter template list type", err, serviceTenantIDField(tenantID), serviceStringField("letter_type", *cleanType))
			return nil, err
		}
		cleanType = &normalized
	}
	items, err := s.employeeLetters.ListEmployeeLetterTemplates(ctx, tenantID, cleanType)
	if err != nil {
		s.logError("list employee letter templates", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) GetEmployeeLetterTemplate(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.EmployeeLetterTemplate, error) {
	if tenantID == uuid.Nil || id == uuid.Nil {
		err := domain.ErrInvalidEmployeeLetterTemplate
		s.logError("validate employee letter template get", err)
		return nil, err
	}
	item, err := s.employeeLetters.GetEmployeeLetterTemplate(ctx, tenantID, id)
	if err != nil {
		s.logError("get employee letter template", err, serviceTenantIDField(tenantID), serviceStringField("employee_letter_template_id", id.String()))
		return nil, err
	}
	return item, nil
}

func (s *TenantService) UpdateEmployeeLetterTemplate(ctx context.Context, cmd ports.EmployeeLetterTemplateCommand) (*domain.EmployeeLetterTemplate, error) {
	if cmd.ID == uuid.Nil {
		err := domain.ErrInvalidEmployeeLetterTemplate
		s.logError("validate employee letter template update id", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	if _, err := s.GetEmployeeLetterTemplate(ctx, cmd.TenantID, cmd.ID); err != nil {
		return nil, err
	}
	item, err := domain.NewEmployeeLetterTemplate(employeeLetterTemplateInput(cmd))
	if err != nil {
		s.logError("validate employee letter template update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("employee_letter_template_id", cmd.ID.String()))
		return nil, err
	}
	item.ID = cmd.ID
	result, err := s.employeeLetters.UpdateEmployeeLetterTemplate(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update employee letter template", err, serviceTenantIDField(cmd.TenantID), serviceStringField("employee_letter_template_id", cmd.ID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) DeleteEmployeeLetterTemplate(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if _, err := s.GetEmployeeLetterTemplate(ctx, tenantID, id); err != nil {
		return err
	}
	if err := s.employeeLetters.DeleteEmployeeLetterTemplate(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete employee letter template", err, serviceTenantIDField(tenantID), serviceStringField("employee_letter_template_id", id.String()))
		return err
	}
	return nil
}

func (s *TenantService) GenerateEmployeeLetter(ctx context.Context, cmd ports.EmployeeLetterCommand) (*domain.EmployeeLetter, error) {
	profile, err := s.employees.GetEmployeeProfile(ctx, cmd.TenantID, cmd.EmployeeID)
	if err != nil {
		s.logError("get employee before letter generation", err, serviceTenantIDField(cmd.TenantID), serviceStringField("employee_id", cmd.EmployeeID.String()))
		return nil, err
	}
	letterType, err := domain.ValidateEmployeeLetterType(cmd.LetterType)
	if err != nil {
		s.logError("validate employee letter type", err, serviceTenantIDField(cmd.TenantID), serviceStringField("letter_type", cmd.LetterType))
		return nil, err
	}
	cmd.LetterType = letterType
	if cmd.TemplateID == nil {
		if tmpl, err := s.employeeLetters.GetDefaultEmployeeLetterTemplate(ctx, cmd.TenantID, letterType); err == nil && tmpl != nil {
			cmd.TemplateID = &tmpl.ID
		}
	}
	if cmd.DocumentTypeID != nil && *cmd.DocumentTypeID != uuid.Nil {
		if _, err := s.employeeDocuments.GetDocumentType(ctx, cmd.TenantID, *cmd.DocumentTypeID); err != nil {
			s.logError("validate employee letter document type", err, serviceTenantIDField(cmd.TenantID), serviceStringField("document_type_id", cmd.DocumentTypeID.String()))
			return nil, err
		}
	}
	item, err := domain.NewEmployeeLetter(employeeLetterInput(cmd, profile.Employee.UserID))
	if err != nil {
		s.logError("validate employee letter create", err, serviceTenantIDField(cmd.TenantID), serviceStringField("employee_id", cmd.EmployeeID.String()))
		return nil, err
	}
	if item.TemplateID != nil {
		tmpl, err := s.GetEmployeeLetterTemplate(ctx, cmd.TenantID, *item.TemplateID)
		if err != nil {
			return nil, err
		}
		renderEmployeeLetter(item, tmpl, profile)
	}
	if cmd.SignatureRequired {
		token, err := domain.NewEmployeeLetterSignatureToken()
		if err != nil {
			s.logError("generate employee letter signature token", err, serviceTenantIDField(cmd.TenantID))
			return nil, err
		}
		item.SignatureToken = &token
	}
	var result *domain.EmployeeLetter
	err = s.system.RunAsSystem(ctx, func(txCtx context.Context) error {
		saved, err := s.employeeLetters.CreateEmployeeLetter(txCtx, item, cmd.ActorID)
		if err != nil {
			return err
		}
		result = saved
		_, _ = s.employeeLetters.CreateEmployeeLetterEvent(txCtx, &domain.EmployeeLetterEvent{TenantID: saved.TenantID, EmployeeLetterID: saved.ID, ToStatus: saved.Status, Action: "generated", Remarks: saved.Subject, Metadata: map[string]any{"version": saved.Version, "letter_type": saved.LetterType}}, cmd.ActorID)
		return nil
	})
	if err != nil {
		s.logError("create employee letter", err, serviceTenantIDField(cmd.TenantID), serviceStringField("employee_id", cmd.EmployeeID.String()), serviceStringField("letter_type", letterType))
		return nil, err
	}
	if _, err := s.ensureEmployeeLetterPDF(ctx, result, cmd.LinkDocument, cmd.ActorID); err != nil {
		s.logError("store generated employee letter pdf", err, serviceTenantIDField(cmd.TenantID), serviceStringField("employee_letter_id", result.ID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListEmployeeLetters(ctx context.Context, filter domain.EmployeeLetterFilter) (*domain.EmployeeLetterPage, error) {
	if filter.TenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate employee letter list tenant", err)
		return nil, err
	}
	filter.Search = cleanStringPtr(filter.Search)
	filter.Status = cleanStringPtr(filter.Status)
	filter.LetterType = cleanStringPtr(filter.LetterType)
	if filter.Status != nil {
		if _, err := domain.ValidateEmployeeLetterStatus(*filter.Status); err != nil {
			s.logError("validate employee letter list status", err, serviceTenantIDField(filter.TenantID), serviceStringField("status", *filter.Status))
			return nil, err
		}
	}
	if filter.LetterType != nil {
		normalized, err := domain.ValidateEmployeeLetterType(*filter.LetterType)
		if err != nil {
			s.logError("validate employee letter list type", err, serviceTenantIDField(filter.TenantID), serviceStringField("letter_type", *filter.LetterType))
			return nil, err
		}
		filter.LetterType = &normalized
	}
	limit, offset := normalizeListWindow(filter.Limit, filter.Offset)
	filter.Limit = limit
	filter.Offset = offset
	items, err := s.employeeLetters.ListEmployeeLetters(ctx, filter)
	if err != nil {
		s.logError("list employee letters", err, serviceTenantIDField(filter.TenantID))
		return nil, err
	}
	total, err := s.employeeLetters.CountEmployeeLetters(ctx, filter)
	if err != nil {
		s.logError("count employee letters", err, serviceTenantIDField(filter.TenantID))
		return nil, err
	}
	page := &domain.EmployeeLetterPage{Items: items, Total: total, Limit: limit, Offset: offset}
	if int64(offset)+int64(len(items)) < total {
		next := offset + limit
		page.NextOffset = &next
	}
	return page, nil
}

func (s *TenantService) GetEmployeeLetter(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.EmployeeLetter, error) {
	if tenantID == uuid.Nil || id == uuid.Nil {
		err := domain.ErrInvalidEmployeeLetterID
		s.logError("validate employee letter get", err)
		return nil, err
	}
	item, err := s.employeeLetters.GetEmployeeLetter(ctx, tenantID, id)
	if err != nil {
		s.logError("get employee letter", err, serviceTenantIDField(tenantID), serviceStringField("employee_letter_id", id.String()))
		return nil, err
	}
	return item, nil
}

func (s *TenantService) GetEmployeeLetterBySignatureToken(ctx context.Context, token string) (*domain.EmployeeLetter, error) {
	if strings.TrimSpace(token) == "" {
		return nil, domain.ErrEmployeeLetterTokenMissing
	}
	return s.employeeLetters.GetEmployeeLetterBySignatureToken(ctx, token)
}

func (s *TenantService) UpdateEmployeeLetterStatus(ctx context.Context, cmd ports.EmployeeLetterStatusCommand) (*domain.EmployeeLetter, error) {
	existing, err := s.GetEmployeeLetter(ctx, cmd.TenantID, cmd.ID)
	if err != nil {
		return nil, err
	}
	status, err := domain.ValidateEmployeeLetterStatus(cmd.Status)
	if err != nil {
		s.logError("validate employee letter status", err, serviceTenantIDField(cmd.TenantID), serviceStringField("employee_letter_id", cmd.ID.String()), serviceStringField("status", cmd.Status))
		return nil, err
	}
	if existing.Status == domain.EmployeeLetterStatusRevoked && status != domain.EmployeeLetterStatusRevoked {
		err := domain.ErrEmployeeLetterCannotBeChanged
		s.logError("validate revoked employee letter status", err, serviceTenantIDField(cmd.TenantID), serviceStringField("employee_letter_id", cmd.ID.String()))
		return nil, err
	}
	result, err := s.employeeLetters.UpdateEmployeeLetterStatus(ctx, cmd.TenantID, cmd.ID, status, cmd.Remarks, cmd.ActorID)
	if err != nil {
		s.logError("update employee letter status", err, serviceTenantIDField(cmd.TenantID), serviceStringField("employee_letter_id", cmd.ID.String()), serviceStringField("status", status))
		return nil, err
	}
	_, _ = s.employeeLetters.CreateEmployeeLetterEvent(ctx, &domain.EmployeeLetterEvent{TenantID: result.TenantID, EmployeeLetterID: result.ID, FromStatus: &existing.Status, ToStatus: result.Status, Action: strings.ToLower(result.Status), Remarks: cmd.Remarks, Metadata: map[string]any{"version": result.Version, "letter_type": result.LetterType}}, cmd.ActorID)
	return result, nil
}

func (s *TenantService) SignEmployeeLetter(ctx context.Context, cmd ports.EmployeeLetterSignatureCommand) (*domain.EmployeeLetter, error) {
	if strings.TrimSpace(cmd.Token) == "" || strings.TrimSpace(cmd.SignerName) == "" || strings.TrimSpace(cmd.SignerEmail) == "" {
		err := domain.ErrEmployeeLetterSignature
		s.logError("validate employee letter signature", err)
		return nil, err
	}
	existing, err := s.employeeLetters.GetEmployeeLetterBySignatureToken(ctx, cmd.Token)
	if err != nil {
		s.logError("get employee letter by signature token", err)
		return nil, err
	}
	hash := employeeLetterSignatureHash(existing, cmd)
	result, err := s.employeeLetters.SignEmployeeLetter(ctx, cmd.Token, cmd.SignerName, cmd.SignerEmail, cmd.IPAddress, cmd.UserAgent, hash)
	if err != nil {
		s.logError("sign employee letter", err, serviceTenantIDField(existing.TenantID), serviceStringField("employee_letter_id", existing.ID.String()))
		return nil, err
	}
	_, _ = s.employeeLetters.CreateEmployeeLetterEvent(ctx, &domain.EmployeeLetterEvent{TenantID: result.TenantID, EmployeeLetterID: result.ID, FromStatus: &existing.Status, ToStatus: result.Status, Action: "signed", ActorEmail: &cmd.SignerEmail, IPAddress: cmd.IPAddress, UserAgent: cmd.UserAgent, Metadata: map[string]any{"signature_hash": hash}}, nil)
	if _, err := s.ensureEmployeeLetterPDF(ctx, result, result.EmployeeDocumentID != nil, nil); err != nil {
		s.logError("refresh signed employee letter pdf", err, serviceTenantIDField(result.TenantID), serviceStringField("employee_letter_id", result.ID.String()))
	}
	return result, nil
}

func (s *TenantService) DeleteEmployeeLetter(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if _, err := s.GetEmployeeLetter(ctx, tenantID, id); err != nil {
		return err
	}
	if err := s.employeeLetters.DeleteEmployeeLetter(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete employee letter", err, serviceTenantIDField(tenantID), serviceStringField("employee_letter_id", id.String()))
		return err
	}
	return nil
}

func (s *TenantService) ListEmployeeLetterEvents(ctx context.Context, tenantID uuid.UUID, employeeLetterID uuid.UUID) ([]*domain.EmployeeLetterEvent, error) {
	if _, err := s.GetEmployeeLetter(ctx, tenantID, employeeLetterID); err != nil {
		return nil, err
	}
	return s.employeeLetters.ListEmployeeLetterEvents(ctx, tenantID, employeeLetterID)
}

func (s *TenantService) RenderEmployeeLetterPDF(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) ([]byte, string, error) {
	letter, err := s.GetEmployeeLetter(ctx, tenantID, id)
	if err != nil {
		return nil, "", err
	}
	data, err := s.ensureEmployeeLetterPDF(ctx, letter, letter.EmployeeDocumentID != nil, actorID)
	if err != nil {
		return nil, "", err
	}
	_, _ = s.employeeLetters.CreateEmployeeLetterEvent(ctx, &domain.EmployeeLetterEvent{TenantID: letter.TenantID, EmployeeLetterID: letter.ID, FromStatus: &letter.Status, ToStatus: letter.Status, Action: "downloaded", Metadata: map[string]any{"version": letter.Version, "letter_type": letter.LetterType}}, actorID)
	return data, employeeLetterPDFName(letter), nil
}

func (s *TenantService) ensureEmployeeLetterPDF(ctx context.Context, letter *domain.EmployeeLetter, linkDocument bool, actorID *uuid.UUID) ([]byte, error) {
	data, err := s.renderEmployeeLetterPDF(ctx, letter)
	if err != nil {
		return nil, err
	}
	if s.employeeLetterStorage == nil {
		return nil, domain.ErrEmployeeLetterPDFMissing
	}
	path, err := s.employeeLetterStorage.StoreEmployeeLetterPDF(ctx, ports.StoreEmployeeLetterPDFInput{TenantID: letter.TenantID, EmployeeID: letter.EmployeeID, LetterID: letter.ID, LetterType: letter.LetterType, FileName: employeeLetterPDFName(letter), ContentType: "application/pdf", Content: data})
	if err != nil {
		return nil, err
	}
	var documentID *uuid.UUID
	if linkDocument {
		doc, err := s.ensureEmployeeLetterDocument(ctx, letter, path, int64(len(data)), actorID)
		if err != nil {
			return nil, err
		}
		documentID = &doc.ID
	}
	updated, err := s.employeeLetters.UpdateEmployeeLetterPDF(ctx, letter.TenantID, letter.ID, &path, documentID, actorID)
	if err != nil {
		return nil, err
	}
	*letter = *updated
	return data, nil
}

func (s *TenantService) ensureEmployeeLetterDocument(ctx context.Context, letter *domain.EmployeeLetter, path string, size int64, actorID *uuid.UUID) (*domain.EmployeeDocument, error) {
	if letter.EmployeeDocumentID != nil && *letter.EmployeeDocumentID != uuid.Nil {
		doc, err := s.employeeDocuments.GetEmployeeDocument(ctx, letter.TenantID, *letter.EmployeeDocumentID)
		if err != nil {
			return nil, err
		}
		doc.FilePath = &path
		doc.FileSizeBytes = &size
		doc.Status = domain.EmployeeDocumentStatusApproved
		return s.employeeDocuments.UpdateEmployeeDocument(ctx, doc, actorID)
	}
	title := fmt.Sprintf("%s letter v%d", strings.Title(letter.LetterType), letter.Version)
	fileName := employeeLetterPDFName(letter)
	item := &domain.EmployeeDocument{TenantID: letter.TenantID, UserID: letter.UserID, DocumentTypeID: letter.DocumentTypeID, Title: &title, FilePath: &path, Status: domain.EmployeeDocumentStatusApproved, OriginalFileName: &fileName, ContentType: employeeLetterStringPtr("application/pdf"), FileSizeBytes: &size, Encrypted: true, EncryptionAlgorithm: "AES-256-GCM"}
	return s.employeeDocuments.CreateEmployeeDocument(ctx, item, actorID)
}

func (s *TenantService) renderEmployeeLetterPDF(ctx context.Context, letter *domain.EmployeeLetter) ([]byte, error) {
	if s.employeeLetterPDF == nil {
		return nil, domain.ErrEmployeeLetterPDFMissing
	}
	employee, _ := s.employees.GetEmployeeByUserID(ctx, letter.TenantID, letter.UserID)
	return s.employeeLetterPDF.RenderEmployeeLetterPDF(ctx, ports.EmployeeLetterDocument{Letter: letter, Employee: employee})
}

func employeeLetterTemplateInput(cmd ports.EmployeeLetterTemplateCommand) domain.EmployeeLetterTemplateInput {
	return domain.EmployeeLetterTemplateInput{TenantID: cmd.TenantID, LetterType: cmd.LetterType, Name: cmd.Name, Description: cmd.Description, Subject: cmd.Subject, BodyHTML: cmd.BodyHTML, FooterHTML: cmd.FooterHTML, Locale: cmd.Locale, IsDefault: cmd.IsDefault, IsActive: cmd.IsActive}
}

func employeeLetterInput(cmd ports.EmployeeLetterCommand, userID uuid.UUID) domain.EmployeeLetterInput {
	return domain.EmployeeLetterInput{TenantID: cmd.TenantID, EmployeeID: cmd.EmployeeID, UserID: userID, TemplateID: cmd.TemplateID, DocumentTypeID: cmd.DocumentTypeID, LetterType: cmd.LetterType, Subject: cmd.Subject, RenderedHTML: cmd.RenderedHTML, IssueDate: cmd.IssueDate, EffectiveDate: cmd.EffectiveDate, EndDate: cmd.EndDate, SignatureRequired: cmd.SignatureRequired, SignerEmail: cmd.SignerEmail}
}

func renderEmployeeLetter(item *domain.EmployeeLetter, tmpl *domain.EmployeeLetterTemplate, profile *domain.EmployeeProfile) {
	subject := valueOrDefault(tmpl.Subject, strings.Title(item.LetterType)+" Letter")
	body := tmpl.BodyHTML
	employee := profile.Employee
	replacements := map[string]string{
		"employee_name":    strings.TrimSpace(strings.Join([]string{employee.Firstname, ptrString(employee.Lastname)}, " ")),
		"first_name":       employee.Firstname,
		"last_name":        ptrString(employee.Lastname),
		"employee_code":    ptrString(employee.EmployeeCode),
		"employee_email":   ptrString(employee.Email),
		"employee_mobile":  ptrString(employee.Mobile),
		"department":       ptrString(employee.DepartmentName),
		"branch":           ptrString(employee.BranchName),
		"designation":      ptrString(employee.DesignationName),
		"joining_date":     formatDate(employee.JoiningDate),
		"resignation_date": formatDate(employee.ResignationDate),
		"issue_date":       formatDate(item.IssueDate),
		"effective_date":   formatDate(item.EffectiveDate),
		"end_date":         formatDate(item.EndDate),
		"letter_type":      strings.Title(item.LetterType),
	}
	for key, value := range replacements {
		body = strings.ReplaceAll(body, "{{"+key+"}}", value)
		subject = strings.ReplaceAll(subject, "{{"+key+"}}", value)
	}
	if tmpl.FooterHTML != nil {
		body += "\n" + *tmpl.FooterHTML
	}
	item.Subject = &subject
	item.RenderedHTML = &body
}

func employeeLetterSignatureHash(letter *domain.EmployeeLetter, cmd ports.EmployeeLetterSignatureCommand) string {
	sum := sha256.Sum256([]byte(strings.Join([]string{letter.ID.String(), ptrString(letter.SignatureToken), cmd.SignerName, cmd.SignerEmail, ptrString(cmd.IPAddress), ptrString(cmd.UserAgent)}, "|")))
	return hex.EncodeToString(sum[:])
}

func employeeLetterPDFName(letter *domain.EmployeeLetter) string {
	return fmt.Sprintf("%s-letter-%s-v%d.pdf", letter.LetterType, letter.EmployeeID.String(), letter.Version)
}

func employeeLetterStringPtr(value string) *string {
	return &value
}
