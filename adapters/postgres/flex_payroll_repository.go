package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateFlexPayRun(ctx context.Context, item *domain.FlexPayRun, actorID *uuid.UUID) (*domain.FlexPayRun, error) {
	row, err := s.getQueries(ctx).CreateFlexPayRun(ctx, sqlc.CreateFlexPayRunParams{
		TenantID: item.TenantID, RunCode: item.RunCode, Title: item.Title, RunType: item.RunType, Status: item.Status,
		PeriodStart: dateFromTime(item.PeriodStart), PeriodEnd: dateFromTime(item.PeriodEnd), PayoutDate: dateFromPtr(item.PayoutDate),
		CurrencyCode: item.CurrencyCode, SourcePolicy: item.SourcePolicy, Notes: textFromPtr(item.Notes), Metadata: jsonBytesFromRaw(item.Metadata),
		CreatedBy: uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "create flex pay run", err, tenantIDField(item.TenantID), stringField("run_code", item.RunCode))
	}
	return mapFlexPayRun(row), nil
}

func (s *Store) GetFlexPayRun(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.FlexPayRun, error) {
	row, err := s.getQueries(ctx).GetFlexPayRun(ctx, sqlc.GetFlexPayRunParams{TenantID: tenantID, ID: id})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrFlexPayRunNotFound
		}
		return nil, s.logDBError(ctx, "get flex pay run", err, tenantIDField(tenantID), stringField("flex_pay_run_id", id.String()))
	}
	return mapFlexPayRun(row), nil
}

func (s *Store) ListFlexPayRuns(ctx context.Context, filter domain.FlexPayRunFilter) ([]*domain.FlexPayRun, error) {
	rows, err := s.getQueries(ctx).ListFlexPayRuns(ctx, sqlc.ListFlexPayRunsParams{
		TenantID: filter.TenantID, Status: textFromPtr(filter.Status), RunType: textFromPtr(filter.RunType),
		DateFrom: dateFromPtr(filter.DateFrom), DateTo: dateFromPtr(filter.DateTo), Search: textFromPtr(filter.Search),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "list flex pay runs", err, tenantIDField(filter.TenantID))
	}
	return mapFlexPayRuns(rows), nil
}

func (s *Store) UpdateFlexPayRunStatus(ctx context.Context, item *domain.FlexPayRun, actorID *uuid.UUID) (*domain.FlexPayRun, error) {
	row, err := s.getQueries(ctx).UpdateFlexPayRunStatus(ctx, sqlc.UpdateFlexPayRunStatusParams{
		TenantID: item.TenantID, ID: item.ID, Status: item.Status, InvoiceCount: item.InvoiceCount, ItemCount: item.ItemCount,
		GrossAmount: numericFromFlexAmount(item.GrossAmount), TdsAmount: numericFromFlexAmount(item.TDSAmount), GstAmount: numericFromFlexAmount(item.GSTAmount), NetAmount: numericFromFlexAmount(item.NetAmount),
		UpdatedBy: uuidFromPtr(actorID), PaymentReference: textFromPtr(item.PaymentReference), ExportBatchRef: textFromPtr(item.ExportBatchRef),
		Notes: textFromPtr(item.Notes), Metadata: jsonBytesFromRaw(item.Metadata),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrFlexPayRunNotFound
		}
		return nil, s.logDBError(ctx, "update flex pay run status", err, tenantIDField(item.TenantID), stringField("flex_pay_run_id", item.ID.String()))
	}
	return mapFlexPayRun(row), nil
}

func (s *Store) DeleteFlexPayRun(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteFlexPayRun(ctx, sqlc.SoftDeleteFlexPayRunParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete flex pay run", err, tenantIDField(tenantID), stringField("flex_pay_run_id", id.String()))
	}
	return nil
}

func (s *Store) CreateContractorInvoice(ctx context.Context, item *domain.ContractorInvoice, actorID *uuid.UUID) (*domain.ContractorInvoice, error) {
	row, err := s.getQueries(ctx).CreateContractorInvoice(ctx, contractorInvoiceCreateParams(item, actorID))
	if err != nil {
		return nil, s.logDBError(ctx, "create contractor invoice", err, tenantIDField(item.TenantID), stringField("invoice_number", item.InvoiceNumber))
	}
	return mapContractorInvoice(row), nil
}

func (s *Store) UpdateContractorInvoice(ctx context.Context, item *domain.ContractorInvoice, actorID *uuid.UUID) (*domain.ContractorInvoice, error) {
	params := contractorInvoiceUpdateParams(item, actorID)
	row, err := s.getQueries(ctx).UpdateContractorInvoice(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrContractorInvoiceNotFound
		}
		return nil, s.logDBError(ctx, "update contractor invoice", err, tenantIDField(item.TenantID), stringField("invoice_id", item.ID.String()))
	}
	return mapContractorInvoice(row), nil
}

func (s *Store) GetContractorInvoice(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.ContractorInvoice, error) {
	row, err := s.getQueries(ctx).GetContractorInvoice(ctx, sqlc.GetContractorInvoiceParams{TenantID: tenantID, ID: id})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrContractorInvoiceNotFound
		}
		return nil, s.logDBError(ctx, "get contractor invoice", err, tenantIDField(tenantID), stringField("invoice_id", id.String()))
	}
	return mapContractorInvoice(row), nil
}

func (s *Store) ListContractorInvoices(ctx context.Context, filter domain.ContractorInvoiceFilter) ([]*domain.ContractorInvoice, error) {
	rows, err := s.getQueries(ctx).ListContractorInvoices(ctx, sqlc.ListContractorInvoicesParams{
		TenantID: filter.TenantID, FlexPayRunID: uuidFromPtr(filter.FlexPayRunID), WorkerProfileID: uuidFromPtr(filter.WorkerProfileID),
		Status: textFromPtr(filter.Status), Search: textFromPtr(filter.Search),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "list contractor invoices", err, tenantIDField(filter.TenantID))
	}
	return mapContractorInvoices(rows), nil
}

func (s *Store) UpdateContractorInvoiceStatus(ctx context.Context, item *domain.ContractorInvoice, actorID *uuid.UUID) (*domain.ContractorInvoice, error) {
	row, err := s.getQueries(ctx).UpdateContractorInvoiceStatus(ctx, sqlc.UpdateContractorInvoiceStatusParams{
		TenantID: item.TenantID, ID: item.ID, Status: item.Status, UpdatedBy: uuidFromPtr(actorID),
		RejectionReason: textFromPtr(item.RejectionReason), PaymentReference: textFromPtr(item.PaymentReference),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrContractorInvoiceNotFound
		}
		return nil, s.logDBError(ctx, "update contractor invoice status", err, tenantIDField(item.TenantID), stringField("invoice_id", item.ID.String()))
	}
	return mapContractorInvoice(row), nil
}

func (s *Store) DeleteContractorInvoice(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteContractorInvoice(ctx, sqlc.SoftDeleteContractorInvoiceParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete contractor invoice", err, tenantIDField(tenantID), stringField("invoice_id", id.String()))
	}
	return nil
}

func (s *Store) CreateFlexPayRunItem(ctx context.Context, item *domain.FlexPayRunItem, actorID *uuid.UUID) (*domain.FlexPayRunItem, error) {
	row, err := s.getQueries(ctx).CreateFlexPayRunItem(ctx, sqlc.CreateFlexPayRunItemParams{
		TenantID: item.TenantID, FlexPayRunID: item.FlexPayRunID, ContractorInvoiceID: uuidFromPtr(item.ContractorInvoiceID),
		WorkerProfileID: item.WorkerProfileID, EngagementID: uuidFromPtr(item.EngagementID), SourceType: item.SourceType, SourceID: uuidFromPtr(item.SourceID),
		Description: item.Description, Quantity: numericFromFlexAmount(item.Quantity), RateAmount: numericFromFlexAmount(item.RateAmount),
		GrossAmount: numericFromFlexAmount(item.GrossAmount), TdsSection: textFromPtr(item.TDSSection), TdsRate: numericFromFlexAmount(item.TDSRate),
		TdsAmount: numericFromFlexAmount(item.TDSAmount), GstRate: numericFromFlexAmount(item.GSTRate), GstAmount: numericFromFlexAmount(item.GSTAmount),
		NetAmount: numericFromFlexAmount(item.NetAmount), Status: item.Status, Metadata: jsonBytesFromRaw(item.Metadata), CreatedBy: uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "create flex pay run item", err, tenantIDField(item.TenantID), stringField("flex_pay_run_id", item.FlexPayRunID.String()))
	}
	return mapFlexPayRunItem(row), nil
}

func (s *Store) ListFlexPayRunItems(ctx context.Context, tenantID uuid.UUID, flexPayRunID uuid.UUID) ([]*domain.FlexPayRunItem, error) {
	rows, err := s.getQueries(ctx).ListFlexPayRunItems(ctx, sqlc.ListFlexPayRunItemsParams{TenantID: tenantID, FlexPayRunID: flexPayRunID})
	if err != nil {
		return nil, s.logDBError(ctx, "list flex pay run items", err, tenantIDField(tenantID), stringField("flex_pay_run_id", flexPayRunID.String()))
	}
	return mapFlexPayRunItems(rows), nil
}

func (s *Store) UpdateFlexPayRunItemsStatusByRun(ctx context.Context, tenantID uuid.UUID, flexPayRunID uuid.UUID, status string, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).UpdateFlexPayRunItemsStatusByRun(ctx, sqlc.UpdateFlexPayRunItemsStatusByRunParams{TenantID: tenantID, FlexPayRunID: flexPayRunID, Status: status, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "update flex pay run item statuses", err, tenantIDField(tenantID), stringField("flex_pay_run_id", flexPayRunID.String()))
	}
	return nil
}

func (s *Store) GetFlexPayRunTotals(ctx context.Context, tenantID uuid.UUID, flexPayRunID uuid.UUID) (*domain.FlexPayRunTotals, error) {
	row, err := s.getQueries(ctx).GetFlexPayRunTotals(ctx, sqlc.GetFlexPayRunTotalsParams{TenantID: tenantID, FlexPayRunID: flexPayRunID})
	if err != nil {
		return nil, s.logDBError(ctx, "get flex pay run totals", err, tenantIDField(tenantID), stringField("flex_pay_run_id", flexPayRunID.String()))
	}
	return &domain.FlexPayRunTotals{InvoiceCount: row.InvoiceCount, ItemCount: row.ItemCount, GrossAmount: floatFromNumeric(row.GrossAmount), TDSAmount: floatFromNumeric(row.TdsAmount), GSTAmount: floatFromNumeric(row.GstAmount), NetAmount: floatFromNumeric(row.NetAmount)}, nil
}

func (s *Store) CreateFlexPayRunEvent(ctx context.Context, item *domain.FlexPayRunEvent) (*domain.FlexPayRunEvent, error) {
	row, err := s.getQueries(ctx).CreateFlexPayRunEvent(ctx, sqlc.CreateFlexPayRunEventParams{
		TenantID: item.TenantID, FlexPayRunID: uuidFromPtr(item.FlexPayRunID), ContractorInvoiceID: uuidFromPtr(item.ContractorInvoiceID),
		EventType: item.EventType, FromStatus: textFromPtr(item.FromStatus), ToStatus: textFromPtr(item.ToStatus), Comment: textFromPtr(item.Comment),
		ActorID: uuidFromPtr(item.ActorID), Metadata: jsonBytesFromRaw(item.Metadata),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "create flex pay run event", err, tenantIDField(item.TenantID), stringField("event_type", item.EventType))
	}
	return mapFlexPayRunEvent(row), nil
}

func (s *Store) ListFlexPayRunEvents(ctx context.Context, tenantID uuid.UUID, flexPayRunID *uuid.UUID, contractorInvoiceID *uuid.UUID) ([]*domain.FlexPayRunEvent, error) {
	rows, err := s.getQueries(ctx).ListFlexPayRunEvents(ctx, sqlc.ListFlexPayRunEventsParams{TenantID: tenantID, FlexPayRunID: uuidFromPtr(flexPayRunID), ContractorInvoiceID: uuidFromPtr(contractorInvoiceID)})
	if err != nil {
		return nil, s.logDBError(ctx, "list flex pay run events", err, tenantIDField(tenantID))
	}
	return mapFlexPayRunEvents(rows), nil
}

func (s *Store) ListApprovedWorkLogPaymentCandidates(ctx context.Context, tenantID uuid.UUID, startDate string, endDate string) ([]*domain.WorkLogPaymentCandidate, error) {
	start, end, err := parseFlexDateRange(startDate, endDate)
	if err != nil {
		return nil, err
	}
	rows, err := s.getQueries(ctx).ListApprovedWorkLogPaymentCandidates(ctx, sqlc.ListApprovedWorkLogPaymentCandidatesParams{TenantID: tenantID, LogDate: dateFromTime(start), LogDate_2: dateFromTime(end)})
	if err != nil {
		return nil, s.logDBError(ctx, "list approved work log payment candidates", err, tenantIDField(tenantID))
	}
	return mapWorkLogPaymentCandidates(rows), nil
}

func (s *Store) ListAcceptedMilestonePaymentCandidates(ctx context.Context, tenantID uuid.UUID, startDate string, endDate string) ([]*domain.MilestonePaymentCandidate, error) {
	start, end, err := parseFlexDateRange(startDate, endDate)
	if err != nil {
		return nil, err
	}
	rows, err := s.getQueries(ctx).ListAcceptedMilestonePaymentCandidates(ctx, sqlc.ListAcceptedMilestonePaymentCandidatesParams{TenantID: tenantID, AcceptedAt: pgtype.Timestamptz{Time: start, Valid: true}, AcceptedAt_2: pgtype.Timestamptz{Time: end.Add(24*time.Hour - time.Nanosecond), Valid: true}})
	if err != nil {
		return nil, s.logDBError(ctx, "list accepted milestone payment candidates", err, tenantIDField(tenantID))
	}
	return mapMilestonePaymentCandidates(rows), nil
}

func contractorInvoiceCreateParams(item *domain.ContractorInvoice, actorID *uuid.UUID) sqlc.CreateContractorInvoiceParams {
	return sqlc.CreateContractorInvoiceParams{
		TenantID: item.TenantID, FlexPayRunID: uuidFromPtr(item.FlexPayRunID), WorkerProfileID: item.WorkerProfileID, EngagementID: uuidFromPtr(item.EngagementID),
		InvoiceNumber: item.InvoiceNumber, InvoiceDate: dateFromTime(item.InvoiceDate), DueDate: dateFromPtr(item.DueDate), Status: item.Status,
		VendorName: item.VendorName, VendorGstin: textFromPtr(item.VendorGSTIN), PlaceOfSupply: textFromPtr(item.PlaceOfSupply), ReverseCharge: item.ReverseCharge,
		CurrencyCode: item.CurrencyCode, GrossAmount: numericFromFlexAmount(item.GrossAmount), TdsSection: textFromPtr(item.TDSSection), TdsRate: numericFromFlexAmount(item.TDSRate),
		TdsAmount: numericFromFlexAmount(item.TDSAmount), GstRate: numericFromFlexAmount(item.GSTRate), GstAmount: numericFromFlexAmount(item.GSTAmount), NetAmount: numericFromFlexAmount(item.NetAmount),
		AttachmentPath: textFromPtr(item.AttachmentPath), Notes: textFromPtr(item.Notes), Metadata: jsonBytesFromRaw(item.Metadata), CreatedBy: uuidFromPtr(actorID),
	}
}

func contractorInvoiceUpdateParams(item *domain.ContractorInvoice, actorID *uuid.UUID) sqlc.UpdateContractorInvoiceParams {
	params := contractorInvoiceCreateParams(item, actorID)
	return sqlc.UpdateContractorInvoiceParams{
		TenantID: params.TenantID, ID: item.ID, FlexPayRunID: params.FlexPayRunID, WorkerProfileID: params.WorkerProfileID, EngagementID: params.EngagementID,
		InvoiceNumber: params.InvoiceNumber, InvoiceDate: params.InvoiceDate, DueDate: params.DueDate, VendorName: params.VendorName, VendorGstin: params.VendorGstin,
		PlaceOfSupply: params.PlaceOfSupply, ReverseCharge: params.ReverseCharge, CurrencyCode: params.CurrencyCode, GrossAmount: params.GrossAmount,
		TdsSection: params.TdsSection, TdsRate: params.TdsRate, TdsAmount: params.TdsAmount, GstRate: params.GstRate, GstAmount: params.GstAmount,
		NetAmount: params.NetAmount, AttachmentPath: params.AttachmentPath, Notes: params.Notes, Metadata: params.Metadata, UpdatedBy: params.CreatedBy,
	}
}

func parseFlexDateRange(startDate string, endDate string) (time.Time, time.Time, error) {
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	return start, end, nil
}
