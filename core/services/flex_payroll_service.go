package services

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) CreateFlexPayRun(ctx context.Context, cmd ports.FlexPayRunCommand) (*domain.FlexPayRun, error) {
	item, err := flexPayRunFromCommand(cmd)
	if err != nil {
		s.logError("validate flex pay run", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	result, err := s.flexPayroll.CreateFlexPayRun(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create flex pay run", err, serviceTenantIDField(cmd.TenantID), serviceStringField("run_code", cmd.RunCode))
		return nil, err
	}
	return s.hydrateFlexPayRun(ctx, result)
}

func (s *TenantService) ListFlexPayRuns(ctx context.Context, filter domain.FlexPayRunFilter) ([]*domain.FlexPayRun, error) {
	items, err := s.flexPayroll.ListFlexPayRuns(ctx, filter)
	if err != nil {
		s.logError("list flex pay runs", err, serviceTenantIDField(filter.TenantID))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) GetFlexPayRun(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.FlexPayRun, error) {
	item, err := s.flexPayroll.GetFlexPayRun(ctx, tenantID, id)
	if err != nil {
		s.logError("get flex pay run", err, serviceTenantIDField(tenantID), serviceStringField("flex_pay_run_id", id.String()))
		return nil, err
	}
	return s.hydrateFlexPayRun(ctx, item)
}

func (s *TenantService) DeleteFlexPayRun(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	run, err := s.flexPayroll.GetFlexPayRun(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if !containsFlexStatus(run.Status, domain.FlexPayStatusDraft, domain.FlexPayStatusRejected, domain.FlexPayStatusGenerated) {
		return domain.ErrFlexPayRunLocked
	}
	return s.flexPayroll.DeleteFlexPayRun(ctx, tenantID, id, actorID)
}

func (s *TenantService) GenerateFlexPayRun(ctx context.Context, cmd ports.FlexPayRunGenerateCommand) (*domain.FlexPayRun, error) {
	run, err := s.flexPayroll.GetFlexPayRun(ctx, cmd.TenantID, cmd.FlexPayRunID)
	if err != nil {
		return nil, err
	}
	if !containsFlexStatus(run.Status, domain.FlexPayStatusDraft, domain.FlexPayStatusGenerated, domain.FlexPayStatusRejected) {
		return nil, domain.ErrFlexPayRunLocked
	}
	workLogs, err := s.flexPayroll.ListApprovedWorkLogPaymentCandidates(ctx, run.TenantID, flexDate(run.PeriodStart), flexDate(run.PeriodEnd))
	if err != nil {
		return nil, err
	}
	milestones, err := s.flexPayroll.ListAcceptedMilestonePaymentCandidates(ctx, run.TenantID, flexDate(run.PeriodStart), flexDate(run.PeriodEnd))
	if err != nil {
		return nil, err
	}
	groups := flexSourceGroups(workLogs, milestones, cmd)
	existingInvoices, err := s.flexPayroll.ListContractorInvoices(ctx, domain.ContractorInvoiceFilter{TenantID: run.TenantID, FlexPayRunID: &run.ID})
	if err != nil {
		return nil, err
	}
	err = s.runSystem(ctx, func(tx context.Context) error {
		invoiceIndex := len(existingInvoices)
		for _, group := range groups {
			invoiceIndex++
			invoice, err := domain.NewContractorInvoice(domain.ContractorInvoice{
				TenantID: run.TenantID, FlexPayRunID: &run.ID, WorkerProfileID: group.WorkerProfileID, EngagementID: group.EngagementID,
				InvoiceNumber: fmt.Sprintf("%s-%03d", run.RunCode, invoiceIndex), InvoiceDate: time.Now().UTC(), Status: domain.FlexPayStatusDraft,
				VendorName: group.WorkerName, CurrencyCode: group.CurrencyCode, GrossAmount: group.GrossAmount, TDSSection: group.TDSSection,
				TDSRate: group.TDSRate, TDSAmount: group.TDSAmount, GSTRate: group.GSTRate, GSTAmount: group.GSTAmount, NetAmount: group.NetAmount,
				Metadata: json.RawMessage(`{"source":"generated"}`),
			})
			if err != nil {
				return err
			}
			createdInvoice, err := s.flexPayroll.CreateContractorInvoice(tx, invoice, cmd.ActorID)
			if err != nil {
				return err
			}
			for _, source := range group.Sources {
				source.Item.ContractorInvoiceID = &createdInvoice.ID
				if _, err := s.flexPayroll.CreateFlexPayRunItem(tx, source.Item, cmd.ActorID); err != nil {
					return err
				}
			}
		}
		return s.refreshFlexPayRunStatus(tx, run, domain.FlexPayStatusGenerated, cmd.ActorID, nil, nil, flexPtrString(fmt.Sprintf("Generated %d new contractor invoice(s).", len(groups))))
	})
	if err != nil {
		s.logError("generate flex pay run", err, serviceTenantIDField(cmd.TenantID), serviceStringField("flex_pay_run_id", cmd.FlexPayRunID.String()))
		return nil, err
	}
	return s.GetFlexPayRun(ctx, cmd.TenantID, cmd.FlexPayRunID)
}

func (s *TenantService) SubmitFlexPayRun(ctx context.Context, cmd ports.FlexPayRunActionCommand) (*domain.FlexPayRun, error) {
	return s.transitionFlexPayRun(ctx, cmd, domain.FlexPayStatusSubmitted, domain.FlexPayStatusDraft, domain.FlexPayStatusGenerated, domain.FlexPayStatusRejected)
}

func (s *TenantService) ApproveFlexPayRun(ctx context.Context, cmd ports.FlexPayRunActionCommand) (*domain.FlexPayRun, error) {
	return s.transitionFlexPayRun(ctx, cmd, domain.FlexPayStatusApproved, domain.FlexPayStatusSubmitted)
}

func (s *TenantService) RejectFlexPayRun(ctx context.Context, cmd ports.FlexPayRunActionCommand) (*domain.FlexPayRun, error) {
	return s.transitionFlexPayRun(ctx, cmd, domain.FlexPayStatusRejected, domain.FlexPayStatusSubmitted)
}

func (s *TenantService) MarkFlexPayRunPaymentPending(ctx context.Context, cmd ports.FlexPayRunActionCommand) (*domain.FlexPayRun, error) {
	return s.transitionFlexPayRun(ctx, cmd, domain.FlexPayStatusPaymentPending, domain.FlexPayStatusApproved)
}

func (s *TenantService) MarkFlexPayRunPaid(ctx context.Context, cmd ports.FlexPayRunActionCommand) (*domain.FlexPayRun, error) {
	return s.transitionFlexPayRun(ctx, cmd, domain.FlexPayStatusPaid, domain.FlexPayStatusApproved, domain.FlexPayStatusPaymentPending)
}

func (s *TenantService) CreateContractorInvoice(ctx context.Context, cmd ports.ContractorInvoiceCommand) (*domain.ContractorInvoice, error) {
	item, err := contractorInvoiceFromCommand(cmd)
	if err != nil {
		s.logError("validate contractor invoice", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	result, err := s.flexPayroll.CreateContractorInvoice(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create contractor invoice", err, serviceTenantIDField(cmd.TenantID), serviceStringField("invoice_number", cmd.InvoiceNumber))
		return nil, err
	}
	if result.FlexPayRunID != nil {
		_ = s.refreshFlexPayRunByID(ctx, result.TenantID, *result.FlexPayRunID, domain.FlexPayStatusGenerated, cmd.ActorID, nil, nil, nil)
	}
	return result, nil
}

func (s *TenantService) UpdateContractorInvoice(ctx context.Context, cmd ports.ContractorInvoiceCommand) (*domain.ContractorInvoice, error) {
	item, err := contractorInvoiceFromCommand(cmd)
	if err != nil {
		return nil, err
	}
	item.ID = cmd.ID
	result, err := s.flexPayroll.UpdateContractorInvoice(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update contractor invoice", err, serviceTenantIDField(cmd.TenantID), serviceStringField("invoice_id", cmd.ID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListContractorInvoices(ctx context.Context, filter domain.ContractorInvoiceFilter) ([]*domain.ContractorInvoice, error) {
	return s.flexPayroll.ListContractorInvoices(ctx, filter)
}

func (s *TenantService) GetContractorInvoice(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.ContractorInvoice, error) {
	return s.flexPayroll.GetContractorInvoice(ctx, tenantID, id)
}

func (s *TenantService) DeleteContractorInvoice(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	invoice, err := s.flexPayroll.GetContractorInvoice(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if !containsFlexStatus(invoice.Status, domain.FlexPayStatusDraft, domain.FlexPayStatusRejected) {
		return domain.ErrFlexPayRunLocked
	}
	return s.flexPayroll.DeleteContractorInvoice(ctx, tenantID, id, actorID)
}

func (s *TenantService) SubmitContractorInvoice(ctx context.Context, cmd ports.ContractorInvoiceActionCommand) (*domain.ContractorInvoice, error) {
	return s.transitionContractorInvoice(ctx, cmd, domain.FlexPayStatusSubmitted, domain.FlexPayStatusDraft, domain.FlexPayStatusRejected)
}

func (s *TenantService) ApproveContractorInvoice(ctx context.Context, cmd ports.ContractorInvoiceActionCommand) (*domain.ContractorInvoice, error) {
	return s.transitionContractorInvoice(ctx, cmd, domain.FlexPayStatusApproved, domain.FlexPayStatusSubmitted)
}

func (s *TenantService) RejectContractorInvoice(ctx context.Context, cmd ports.ContractorInvoiceActionCommand) (*domain.ContractorInvoice, error) {
	return s.transitionContractorInvoice(ctx, cmd, domain.FlexPayStatusRejected, domain.FlexPayStatusSubmitted)
}

func (s *TenantService) MarkContractorInvoicePaid(ctx context.Context, cmd ports.ContractorInvoiceActionCommand) (*domain.ContractorInvoice, error) {
	return s.transitionContractorInvoice(ctx, cmd, domain.FlexPayStatusPaid, domain.FlexPayStatusApproved, domain.FlexPayStatusPaymentPending)
}

func (s *TenantService) CreateFlexPayRunItem(ctx context.Context, cmd ports.FlexPayRunItemCommand) (*domain.FlexPayRunItem, error) {
	tds, gst, net := domain.FlexTaxAmounts(cmd.Quantity*cmd.RateAmount, cmd.TDSRate, cmd.GSTRate)
	item, err := domain.NewFlexPayRunItem(domain.FlexPayRunItem{
		TenantID: cmd.TenantID, FlexPayRunID: cmd.FlexPayRunID, ContractorInvoiceID: cmd.ContractorInvoiceID, WorkerProfileID: cmd.WorkerProfileID,
		EngagementID: cmd.EngagementID, SourceType: cmd.SourceType, Description: cmd.Description, Quantity: cmd.Quantity, RateAmount: cmd.RateAmount,
		GrossAmount: roundFlex(cmd.Quantity * cmd.RateAmount), TDSSection: cmd.TDSSection, TDSRate: cmd.TDSRate, TDSAmount: tds, GSTRate: cmd.GSTRate, GSTAmount: gst, NetAmount: net,
		Status: domain.FlexPayStatusDraft, Metadata: cmd.Metadata,
	})
	if err != nil {
		return nil, err
	}
	result, err := s.flexPayroll.CreateFlexPayRunItem(ctx, item, cmd.ActorID)
	if err != nil {
		return nil, err
	}
	_ = s.refreshFlexPayRunByID(ctx, cmd.TenantID, cmd.FlexPayRunID, domain.FlexPayStatusGenerated, cmd.ActorID, nil, nil, nil)
	return result, nil
}

func (s *TenantService) ListFlexPayRunItems(ctx context.Context, tenantID uuid.UUID, flexPayRunID uuid.UUID) ([]*domain.FlexPayRunItem, error) {
	return s.flexPayroll.ListFlexPayRunItems(ctx, tenantID, flexPayRunID)
}

func (s *TenantService) ListFlexPayRunEvents(ctx context.Context, tenantID uuid.UUID, flexPayRunID *uuid.UUID, invoiceID *uuid.UUID) ([]*domain.FlexPayRunEvent, error) {
	return s.flexPayroll.ListFlexPayRunEvents(ctx, tenantID, flexPayRunID, invoiceID)
}

func (s *TenantService) ExportFlexPayRunCSV(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) ([]byte, string, error) {
	run, err := s.GetFlexPayRun(ctx, tenantID, id)
	if err != nil {
		return nil, "", err
	}
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	_ = writer.Write([]string{"Run Code", "Run Title", "Invoice", "Worker", "Worker Code", "Engagement", "Source", "Description", "Quantity", "Rate", "Gross", "TDS Section", "TDS Rate", "TDS", "GST Rate", "GST", "Net", "Status", "Payment Reference"})
	for _, item := range run.Items {
		_ = writer.Write([]string{
			run.RunCode, run.Title, flexStringValue(item.InvoiceNumber), flexStringValue(item.WorkerDisplayName), flexStringValue(item.WorkerCode), flexStringValue(item.EngagementTitle),
			item.SourceType, item.Description, flexMoneyString(item.Quantity), flexMoneyString(item.RateAmount), flexMoneyString(item.GrossAmount), flexStringValue(item.TDSSection),
			flexMoneyString(item.TDSRate), flexMoneyString(item.TDSAmount), flexMoneyString(item.GSTRate), flexMoneyString(item.GSTAmount), flexMoneyString(item.NetAmount), item.Status, flexStringValue(run.PaymentReference),
		})
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, "", err
	}
	return buf.Bytes(), fmt.Sprintf("%s-flex-payroll.csv", strings.ToLower(run.RunCode)), nil
}

func (s *TenantService) transitionFlexPayRun(ctx context.Context, cmd ports.FlexPayRunActionCommand, next string, allowed ...string) (*domain.FlexPayRun, error) {
	run, err := s.flexPayroll.GetFlexPayRun(ctx, cmd.TenantID, cmd.FlexPayRunID)
	if err != nil {
		return nil, err
	}
	if !containsFlexStatus(run.Status, allowed...) {
		return nil, domain.ErrFlexPayRunLocked
	}
	previous := run.Status
	err = s.runSystem(ctx, func(tx context.Context) error {
		return s.refreshFlexPayRunStatus(tx, run, next, cmd.ActorID, cmd.PaymentReference, cmd.ExportBatchRef, cmd.Comment, previous)
	})
	if err != nil {
		return nil, err
	}
	return s.GetFlexPayRun(ctx, cmd.TenantID, cmd.FlexPayRunID)
}

func (s *TenantService) refreshFlexPayRunStatus(ctx context.Context, run *domain.FlexPayRun, next string, actorID *uuid.UUID, paymentRef *string, exportRef *string, comment *string, previousOverride ...string) error {
	totals, err := s.flexPayroll.GetFlexPayRunTotals(ctx, run.TenantID, run.ID)
	if err != nil {
		return err
	}
	previous := run.Status
	if len(previousOverride) > 0 && previousOverride[0] != "" {
		previous = previousOverride[0]
	}
	run.Status = next
	run.InvoiceCount = totals.InvoiceCount
	run.ItemCount = totals.ItemCount
	run.GrossAmount = totals.GrossAmount
	run.TDSAmount = totals.TDSAmount
	run.GSTAmount = totals.GSTAmount
	run.NetAmount = totals.NetAmount
	run.PaymentReference = paymentRef
	run.ExportBatchRef = exportRef
	if comment != nil {
		run.Notes = comment
	}
	if run.Metadata == nil {
		run.Metadata = json.RawMessage(`{}`)
	}
	if _, err := s.flexPayroll.UpdateFlexPayRunStatus(ctx, run, actorID); err != nil {
		return err
	}
	if err := s.flexPayroll.UpdateFlexPayRunItemsStatusByRun(ctx, run.TenantID, run.ID, next, actorID); err != nil {
		return err
	}
	invoices, err := s.flexPayroll.ListContractorInvoices(ctx, domain.ContractorInvoiceFilter{TenantID: run.TenantID, FlexPayRunID: &run.ID})
	if err != nil {
		return err
	}
	for _, invoice := range invoices {
		invoice.Status = next
		invoice.PaymentReference = paymentRef
		if next == domain.FlexPayStatusRejected {
			invoice.RejectionReason = comment
		}
		if _, err := s.flexPayroll.UpdateContractorInvoiceStatus(ctx, invoice, actorID); err != nil {
			return err
		}
	}
	runID := run.ID
	_, err = s.flexPayroll.CreateFlexPayRunEvent(ctx, &domain.FlexPayRunEvent{TenantID: run.TenantID, FlexPayRunID: &runID, EventType: "status_change", FromStatus: &previous, ToStatus: &next, Comment: comment, ActorID: actorID, Metadata: json.RawMessage(`{}`)})
	return err
}

func (s *TenantService) refreshFlexPayRunByID(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, actorID *uuid.UUID, paymentRef *string, exportRef *string, comment *string) error {
	run, err := s.flexPayroll.GetFlexPayRun(ctx, tenantID, id)
	if err != nil {
		return err
	}
	return s.refreshFlexPayRunStatus(ctx, run, status, actorID, paymentRef, exportRef, comment)
}

func (s *TenantService) transitionContractorInvoice(ctx context.Context, cmd ports.ContractorInvoiceActionCommand, next string, allowed ...string) (*domain.ContractorInvoice, error) {
	invoice, err := s.flexPayroll.GetContractorInvoice(ctx, cmd.TenantID, cmd.InvoiceID)
	if err != nil {
		return nil, err
	}
	if !containsFlexStatus(invoice.Status, allowed...) {
		return nil, domain.ErrFlexPayRunLocked
	}
	previous := invoice.Status
	invoice.Status = next
	invoice.PaymentReference = cmd.PaymentReference
	if next == domain.FlexPayStatusRejected {
		invoice.RejectionReason = cmd.Comment
	}
	result, err := s.flexPayroll.UpdateContractorInvoiceStatus(ctx, invoice, cmd.ActorID)
	if err != nil {
		return nil, err
	}
	invoiceID := invoice.ID
	_, _ = s.flexPayroll.CreateFlexPayRunEvent(ctx, &domain.FlexPayRunEvent{TenantID: invoice.TenantID, FlexPayRunID: invoice.FlexPayRunID, ContractorInvoiceID: &invoiceID, EventType: "invoice_status_change", FromStatus: &previous, ToStatus: &next, Comment: cmd.Comment, ActorID: cmd.ActorID, Metadata: json.RawMessage(`{}`)})
	return result, nil
}

func (s *TenantService) hydrateFlexPayRun(ctx context.Context, run *domain.FlexPayRun) (*domain.FlexPayRun, error) {
	items, err := s.flexPayroll.ListFlexPayRunItems(ctx, run.TenantID, run.ID)
	if err != nil {
		return nil, err
	}
	invoices, err := s.flexPayroll.ListContractorInvoices(ctx, domain.ContractorInvoiceFilter{TenantID: run.TenantID, FlexPayRunID: &run.ID})
	if err != nil {
		return nil, err
	}
	events, err := s.flexPayroll.ListFlexPayRunEvents(ctx, run.TenantID, &run.ID, nil)
	if err != nil {
		return nil, err
	}
	run.Items = items
	run.Invoices = invoices
	run.Events = events
	return run, nil
}

func flexPayRunFromCommand(cmd ports.FlexPayRunCommand) (*domain.FlexPayRun, error) {
	start, err := parseFlexCommandDate(cmd.PeriodStart)
	if err != nil {
		return nil, domain.ErrInvalidFlexPayRun
	}
	end, err := parseFlexCommandDate(cmd.PeriodEnd)
	if err != nil {
		return nil, domain.ErrInvalidFlexPayRun
	}
	payout, err := parseOptionalFlexCommandDate(cmd.PayoutDate)
	if err != nil {
		return nil, domain.ErrInvalidFlexPayRun
	}
	return domain.NewFlexPayRun(domain.FlexPayRun{ID: cmd.ID, TenantID: cmd.TenantID, RunCode: cmd.RunCode, Title: cmd.Title, RunType: cmd.RunType, Status: domain.FlexPayStatusDraft, PeriodStart: start, PeriodEnd: end, PayoutDate: payout, CurrencyCode: cmd.CurrencyCode, SourcePolicy: cmd.SourcePolicy, Notes: cmd.Notes, Metadata: cmd.Metadata})
}

func contractorInvoiceFromCommand(cmd ports.ContractorInvoiceCommand) (*domain.ContractorInvoice, error) {
	invoiceDate, err := parseFlexCommandDate(cmd.InvoiceDate)
	if err != nil {
		return nil, domain.ErrInvalidContractorInvoice
	}
	dueDate, err := parseOptionalFlexCommandDate(cmd.DueDate)
	if err != nil {
		return nil, domain.ErrInvalidContractorInvoice
	}
	tds, gst, net := domain.FlexTaxAmounts(cmd.GrossAmount, cmd.TDSRate, cmd.GSTRate)
	return domain.NewContractorInvoice(domain.ContractorInvoice{ID: cmd.ID, TenantID: cmd.TenantID, FlexPayRunID: cmd.FlexPayRunID, WorkerProfileID: cmd.WorkerProfileID, EngagementID: cmd.EngagementID, InvoiceNumber: cmd.InvoiceNumber, InvoiceDate: invoiceDate, DueDate: dueDate, Status: domain.FlexPayStatusDraft, VendorName: cmd.VendorName, VendorGSTIN: cmd.VendorGSTIN, PlaceOfSupply: cmd.PlaceOfSupply, ReverseCharge: cmd.ReverseCharge, CurrencyCode: cmd.CurrencyCode, GrossAmount: cmd.GrossAmount, TDSSection: cmd.TDSSection, TDSRate: cmd.TDSRate, TDSAmount: tds, GSTRate: cmd.GSTRate, GSTAmount: gst, NetAmount: net, AttachmentPath: cmd.AttachmentPath, Notes: cmd.Notes, Metadata: cmd.Metadata})
}

func parseFlexCommandDate(value string) (time.Time, error) {
	return time.Parse("2006-01-02", strings.TrimSpace(value))
}

func parseOptionalFlexCommandDate(value *string) (*time.Time, error) {
	if value == nil || strings.TrimSpace(*value) == "" {
		return nil, nil
	}
	parsed, err := parseFlexCommandDate(*value)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

func (s *TenantService) runSystem(ctx context.Context, fn func(context.Context) error) error {
	if s.system == nil {
		return fn(ctx)
	}
	return s.system.RunAsSystem(ctx, fn)
}

type flexSourceGroup struct {
	WorkerProfileID uuid.UUID
	EngagementID    *uuid.UUID
	WorkerName      string
	CurrencyCode    string
	TDSSection      *string
	TDSRate         float64
	GSTRate         float64
	GrossAmount     float64
	TDSAmount       float64
	GSTAmount       float64
	NetAmount       float64
	Sources         []flexSource
}

type flexSource struct {
	Item *domain.FlexPayRunItem
}

func flexSourceGroups(workLogs []*domain.WorkLogPaymentCandidate, milestones []*domain.MilestonePaymentCandidate, cmd ports.FlexPayRunGenerateCommand) []*flexSourceGroup {
	groups := map[string]*flexSourceGroup{}
	for _, row := range workLogs {
		engagementID := row.EngagementID
		quantity := row.BillableHours
		if row.RateUnit == domain.EngagementRateDay {
			quantity = row.BillableHours / 8
		}
		gross := roundFlex(quantity * row.RateAmount)
		section := defaultTDSSection(cmd.TDSSection, domain.FlexPayTDS194J)
		tdsRate, gstRate := defaultRates(cmd)
		tds, gst, net := domain.FlexTaxAmounts(gross, tdsRate, gstRate)
		key := flexGroupKey(row.WorkerProfileID, &engagementID)
		group := getFlexGroup(groups, key, row.WorkerProfileID, &engagementID, row.WorkerDisplayName, row.CurrencyCode, section, tdsRate, gstRate)
		sourceID := row.WorkLogID
		group.Sources = append(group.Sources, flexSource{Item: mustFlexItem(domain.FlexPayRunItem{TenantID: row.TenantID, FlexPayRunID: cmd.FlexPayRunID, WorkerProfileID: row.WorkerProfileID, EngagementID: &engagementID, SourceType: domain.FlexPaySourceWorkLog, SourceID: &sourceID, Description: fmt.Sprintf("Approved work log %s - %s", flexDate(row.LogDate), row.EngagementTitle), Quantity: quantity, RateAmount: row.RateAmount, GrossAmount: gross, TDSSection: section, TDSRate: tdsRate, TDSAmount: tds, GSTRate: gstRate, GSTAmount: gst, NetAmount: net, Status: domain.FlexPayStatusDraft, Metadata: json.RawMessage(`{"source":"work_log"}`)})})
		addFlexGroupTotals(group, gross, tds, gst, net)
	}
	for _, row := range milestones {
		section := defaultTDSSection(cmd.TDSSection, domain.FlexPayTDS194C)
		tdsRate, gstRate := defaultRates(cmd)
		tds, gst, net := domain.FlexTaxAmounts(row.Amount, tdsRate, gstRate)
		key := flexGroupKey(row.WorkerProfileID, row.EngagementID)
		group := getFlexGroup(groups, key, row.WorkerProfileID, row.EngagementID, row.WorkerDisplayName, row.CurrencyCode, section, tdsRate, gstRate)
		sourceID := row.MilestoneID
		group.Sources = append(group.Sources, flexSource{Item: mustFlexItem(domain.FlexPayRunItem{TenantID: row.TenantID, FlexPayRunID: cmd.FlexPayRunID, WorkerProfileID: row.WorkerProfileID, EngagementID: row.EngagementID, SourceType: domain.FlexPaySourceMilestone, SourceID: &sourceID, Description: fmt.Sprintf("Accepted milestone %s - %s", flexStringValue(row.MilestoneCode), row.Title), Quantity: 1, RateAmount: row.Amount, GrossAmount: row.Amount, TDSSection: section, TDSRate: tdsRate, TDSAmount: tds, GSTRate: gstRate, GSTAmount: gst, NetAmount: net, Status: domain.FlexPayStatusDraft, Metadata: json.RawMessage(`{"source":"milestone"}`)})})
		addFlexGroupTotals(group, row.Amount, tds, gst, net)
	}
	out := make([]*flexSourceGroup, 0, len(groups))
	for _, group := range groups {
		out = append(out, group)
	}
	return out
}

func getFlexGroup(groups map[string]*flexSourceGroup, key string, workerID uuid.UUID, engagementID *uuid.UUID, workerName string, currency string, section *string, tdsRate float64, gstRate float64) *flexSourceGroup {
	if group, ok := groups[key]; ok {
		return group
	}
	group := &flexSourceGroup{WorkerProfileID: workerID, EngagementID: engagementID, WorkerName: workerName, CurrencyCode: currency, TDSSection: section, TDSRate: tdsRate, GSTRate: gstRate}
	groups[key] = group
	return group
}

func addFlexGroupTotals(group *flexSourceGroup, gross float64, tds float64, gst float64, net float64) {
	group.GrossAmount = roundFlex(group.GrossAmount + gross)
	group.TDSAmount = roundFlex(group.TDSAmount + tds)
	group.GSTAmount = roundFlex(group.GSTAmount + gst)
	group.NetAmount = roundFlex(group.NetAmount + net)
}

func mustFlexItem(item domain.FlexPayRunItem) *domain.FlexPayRunItem {
	result, err := domain.NewFlexPayRunItem(item)
	if err != nil {
		panic(err)
	}
	return result
}

func defaultTDSSection(value *string, fallback string) *string {
	if value != nil && strings.TrimSpace(*value) != "" {
		clean := strings.TrimSpace(*value)
		return &clean
	}
	return &fallback
}

func defaultRates(cmd ports.FlexPayRunGenerateCommand) (float64, float64) {
	tdsRate := 0.0
	gstRate := 0.0
	if cmd.TDSRate != nil {
		tdsRate = *cmd.TDSRate
	}
	if cmd.GSTRate != nil {
		gstRate = *cmd.GSTRate
	}
	return tdsRate, gstRate
}

func containsFlexStatus(value string, allowed ...string) bool {
	for _, status := range allowed {
		if value == status {
			return true
		}
	}
	return false
}

func flexGroupKey(workerID uuid.UUID, engagementID *uuid.UUID) string {
	if engagementID == nil {
		return workerID.String()
	}
	return workerID.String() + ":" + engagementID.String()
}

func flexDate(value time.Time) string {
	return value.Format("2006-01-02")
}

func flexMoneyString(value float64) string {
	return fmt.Sprintf("%.2f", value)
}

func flexStringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func flexPtrString(value string) *string {
	return &value
}

func roundFlex(value float64) float64 {
	return math.Round(value*100) / 100
}
