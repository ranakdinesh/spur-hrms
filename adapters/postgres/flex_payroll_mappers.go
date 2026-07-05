package postgres

import (
	"encoding/json"
	"math"
	"math/big"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapFlexPayRun(row sqlc.HrmsFlexPayRun) *domain.FlexPayRun {
	return &domain.FlexPayRun{
		ID:               row.ID,
		TenantID:         row.TenantID,
		RunCode:          row.RunCode,
		Title:            row.Title,
		RunType:          row.RunType,
		Status:           row.Status,
		PeriodStart:      timeFromDate(row.PeriodStart),
		PeriodEnd:        timeFromDate(row.PeriodEnd),
		PayoutDate:       ptrFromDate(row.PayoutDate),
		CurrencyCode:     row.CurrencyCode,
		SourcePolicy:     row.SourcePolicy,
		InvoiceCount:     row.InvoiceCount,
		ItemCount:        row.ItemCount,
		GrossAmount:      floatFromNumeric(row.GrossAmount),
		TDSAmount:        floatFromNumeric(row.TdsAmount),
		GSTAmount:        floatFromNumeric(row.GstAmount),
		NetAmount:        floatFromNumeric(row.NetAmount),
		GeneratedAt:      ptrFromTimestamptz(row.GeneratedAt),
		SubmittedAt:      ptrFromTimestamptz(row.SubmittedAt),
		SubmittedBy:      ptrFromUUID(row.SubmittedBy),
		ApprovedAt:       ptrFromTimestamptz(row.ApprovedAt),
		ApprovedBy:       ptrFromUUID(row.ApprovedBy),
		RejectedAt:       ptrFromTimestamptz(row.RejectedAt),
		RejectedBy:       ptrFromUUID(row.RejectedBy),
		PaidAt:           ptrFromTimestamptz(row.PaidAt),
		PaidBy:           ptrFromUUID(row.PaidBy),
		PaymentReference: ptrFromText(row.PaymentReference),
		ExportBatchRef:   ptrFromText(row.ExportBatchRef),
		Notes:            ptrFromText(row.Notes),
		Metadata:         json.RawMessage(row.Metadata),
		Inactive:         row.Inactive,
		CreatedAt:        timeFromTimestamptz(row.CreatedAt),
		CreatedBy:        ptrFromUUID(row.CreatedBy),
		UpdatedAt:        timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:        ptrFromUUID(row.UpdatedBy),
	}
}

func mapFlexPayRuns(rows []sqlc.HrmsFlexPayRun) []*domain.FlexPayRun {
	items := make([]*domain.FlexPayRun, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapFlexPayRun(row))
	}
	return items
}

func mapContractorInvoice(row sqlc.HrmsContractorInvoice) *domain.ContractorInvoice {
	return &domain.ContractorInvoice{
		ID:               row.ID,
		TenantID:         row.TenantID,
		FlexPayRunID:     ptrFromUUID(row.FlexPayRunID),
		WorkerProfileID:  row.WorkerProfileID,
		EngagementID:     ptrFromUUID(row.EngagementID),
		InvoiceNumber:    row.InvoiceNumber,
		InvoiceDate:      timeFromDate(row.InvoiceDate),
		DueDate:          ptrFromDate(row.DueDate),
		Status:           row.Status,
		VendorName:       row.VendorName,
		VendorGSTIN:      ptrFromText(row.VendorGstin),
		PlaceOfSupply:    ptrFromText(row.PlaceOfSupply),
		ReverseCharge:    row.ReverseCharge,
		CurrencyCode:     row.CurrencyCode,
		GrossAmount:      floatFromNumeric(row.GrossAmount),
		TDSSection:       ptrFromText(row.TdsSection),
		TDSRate:          floatFromNumeric(row.TdsRate),
		TDSAmount:        floatFromNumeric(row.TdsAmount),
		GSTRate:          floatFromNumeric(row.GstRate),
		GSTAmount:        floatFromNumeric(row.GstAmount),
		NetAmount:        floatFromNumeric(row.NetAmount),
		SubmittedAt:      ptrFromTimestamptz(row.SubmittedAt),
		SubmittedBy:      ptrFromUUID(row.SubmittedBy),
		ApprovedAt:       ptrFromTimestamptz(row.ApprovedAt),
		ApprovedBy:       ptrFromUUID(row.ApprovedBy),
		RejectedAt:       ptrFromTimestamptz(row.RejectedAt),
		RejectedBy:       ptrFromUUID(row.RejectedBy),
		RejectionReason:  ptrFromText(row.RejectionReason),
		PaidAt:           ptrFromTimestamptz(row.PaidAt),
		PaidBy:           ptrFromUUID(row.PaidBy),
		PaymentReference: ptrFromText(row.PaymentReference),
		AttachmentPath:   ptrFromText(row.AttachmentPath),
		Notes:            ptrFromText(row.Notes),
		Metadata:         json.RawMessage(row.Metadata),
		Inactive:         row.Inactive,
		CreatedAt:        timeFromTimestamptz(row.CreatedAt),
		CreatedBy:        ptrFromUUID(row.CreatedBy),
		UpdatedAt:        timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:        ptrFromUUID(row.UpdatedBy),
	}
}

func mapContractorInvoiceList(row sqlc.ListContractorInvoicesRow) *domain.ContractorInvoice {
	item := mapContractorInvoice(sqlc.HrmsContractorInvoice{
		ID: row.ID, TenantID: row.TenantID, FlexPayRunID: row.FlexPayRunID, WorkerProfileID: row.WorkerProfileID,
		EngagementID: row.EngagementID, InvoiceNumber: row.InvoiceNumber, InvoiceDate: row.InvoiceDate, DueDate: row.DueDate,
		Status: row.Status, VendorName: row.VendorName, VendorGstin: row.VendorGstin, PlaceOfSupply: row.PlaceOfSupply,
		ReverseCharge: row.ReverseCharge, CurrencyCode: row.CurrencyCode, GrossAmount: row.GrossAmount, TdsSection: row.TdsSection,
		TdsRate: row.TdsRate, TdsAmount: row.TdsAmount, GstRate: row.GstRate, GstAmount: row.GstAmount, NetAmount: row.NetAmount,
		SubmittedAt: row.SubmittedAt, SubmittedBy: row.SubmittedBy, ApprovedAt: row.ApprovedAt, ApprovedBy: row.ApprovedBy,
		RejectedAt: row.RejectedAt, RejectedBy: row.RejectedBy, RejectionReason: row.RejectionReason, PaidAt: row.PaidAt,
		PaidBy: row.PaidBy, PaymentReference: row.PaymentReference, AttachmentPath: row.AttachmentPath, Notes: row.Notes,
		Metadata: row.Metadata, Inactive: row.Inactive, CreatedAt: row.CreatedAt, CreatedBy: row.CreatedBy, UpdatedAt: row.UpdatedAt, UpdatedBy: row.UpdatedBy,
	})
	item.WorkerDisplayName = &row.WorkerDisplayName
	item.WorkerCode = ptrFromText(row.WorkerCode)
	item.EngagementTitle = ptrFromText(row.EngagementTitle)
	item.EngagementCode = ptrFromText(row.EngagementCode)
	return item
}

func mapContractorInvoices(rows []sqlc.ListContractorInvoicesRow) []*domain.ContractorInvoice {
	items := make([]*domain.ContractorInvoice, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapContractorInvoiceList(row))
	}
	return items
}

func mapFlexPayRunItem(row sqlc.HrmsFlexPayRunItem) *domain.FlexPayRunItem {
	return &domain.FlexPayRunItem{
		ID:                  row.ID,
		TenantID:            row.TenantID,
		FlexPayRunID:        row.FlexPayRunID,
		ContractorInvoiceID: ptrFromUUID(row.ContractorInvoiceID),
		WorkerProfileID:     row.WorkerProfileID,
		EngagementID:        ptrFromUUID(row.EngagementID),
		SourceType:          row.SourceType,
		SourceID:            ptrFromUUID(row.SourceID),
		Description:         row.Description,
		Quantity:            floatFromNumeric(row.Quantity),
		RateAmount:          floatFromNumeric(row.RateAmount),
		GrossAmount:         floatFromNumeric(row.GrossAmount),
		TDSSection:          ptrFromText(row.TdsSection),
		TDSRate:             floatFromNumeric(row.TdsRate),
		TDSAmount:           floatFromNumeric(row.TdsAmount),
		GSTRate:             floatFromNumeric(row.GstRate),
		GSTAmount:           floatFromNumeric(row.GstAmount),
		NetAmount:           floatFromNumeric(row.NetAmount),
		Status:              row.Status,
		Metadata:            json.RawMessage(row.Metadata),
		Inactive:            row.Inactive,
		CreatedAt:           timeFromTimestamptz(row.CreatedAt),
		CreatedBy:           ptrFromUUID(row.CreatedBy),
		UpdatedAt:           timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:           ptrFromUUID(row.UpdatedBy),
	}
}

func mapFlexPayRunItemList(row sqlc.ListFlexPayRunItemsRow) *domain.FlexPayRunItem {
	item := mapFlexPayRunItem(sqlc.HrmsFlexPayRunItem{
		ID: row.ID, TenantID: row.TenantID, FlexPayRunID: row.FlexPayRunID, ContractorInvoiceID: row.ContractorInvoiceID,
		WorkerProfileID: row.WorkerProfileID, EngagementID: row.EngagementID, SourceType: row.SourceType, SourceID: row.SourceID,
		Description: row.Description, Quantity: row.Quantity, RateAmount: row.RateAmount, GrossAmount: row.GrossAmount,
		TdsSection: row.TdsSection, TdsRate: row.TdsRate, TdsAmount: row.TdsAmount, GstRate: row.GstRate, GstAmount: row.GstAmount,
		NetAmount: row.NetAmount, Status: row.Status, Metadata: row.Metadata, Inactive: row.Inactive, CreatedAt: row.CreatedAt,
		CreatedBy: row.CreatedBy, UpdatedAt: row.UpdatedAt, UpdatedBy: row.UpdatedBy,
	})
	item.WorkerDisplayName = &row.WorkerDisplayName
	item.WorkerCode = ptrFromText(row.WorkerCode)
	item.EngagementTitle = ptrFromText(row.EngagementTitle)
	item.EngagementCode = ptrFromText(row.EngagementCode)
	item.InvoiceNumber = ptrFromText(row.InvoiceNumber)
	return item
}

func mapFlexPayRunItems(rows []sqlc.ListFlexPayRunItemsRow) []*domain.FlexPayRunItem {
	items := make([]*domain.FlexPayRunItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapFlexPayRunItemList(row))
	}
	return items
}

func mapFlexPayRunEvent(row sqlc.HrmsFlexPayRunEvent) *domain.FlexPayRunEvent {
	return &domain.FlexPayRunEvent{
		ID:                  row.ID,
		TenantID:            row.TenantID,
		FlexPayRunID:        ptrFromUUID(row.FlexPayRunID),
		ContractorInvoiceID: ptrFromUUID(row.ContractorInvoiceID),
		EventType:           row.EventType,
		FromStatus:          ptrFromText(row.FromStatus),
		ToStatus:            ptrFromText(row.ToStatus),
		Comment:             ptrFromText(row.Comment),
		ActorID:             ptrFromUUID(row.ActorID),
		Metadata:            json.RawMessage(row.Metadata),
		CreatedAt:           timeFromTimestamptz(row.CreatedAt),
	}
}

func mapFlexPayRunEvents(rows []sqlc.HrmsFlexPayRunEvent) []*domain.FlexPayRunEvent {
	items := make([]*domain.FlexPayRunEvent, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapFlexPayRunEvent(row))
	}
	return items
}

func mapWorkLogPaymentCandidates(rows []sqlc.ListApprovedWorkLogPaymentCandidatesRow) []*domain.WorkLogPaymentCandidate {
	items := make([]*domain.WorkLogPaymentCandidate, 0, len(rows))
	for _, row := range rows {
		items = append(items, &domain.WorkLogPaymentCandidate{
			WorkLogID: row.WorkLogID, TenantID: row.TenantID, EngagementID: row.EngagementID, WorkerProfileID: row.WorkerProfileID,
			LogDate: timeFromDate(row.LogDate), BillableHours: floatFromNumeric(row.BillableHours), RateAmount: floatFromNumeric(row.RateAmount),
			RateUnit: row.RateUnit, CurrencyCode: row.CurrencyCode, EngagementTitle: row.EngagementTitle, EngagementCode: ptrFromText(row.EngagementCode),
			WorkerDisplayName: row.WorkerDisplayName, WorkerCode: ptrFromText(row.WorkerCode),
		})
	}
	return items
}

func mapMilestonePaymentCandidates(rows []sqlc.ListAcceptedMilestonePaymentCandidatesRow) []*domain.MilestonePaymentCandidate {
	items := make([]*domain.MilestonePaymentCandidate, 0, len(rows))
	for _, row := range rows {
		items = append(items, &domain.MilestonePaymentCandidate{
			MilestoneID: row.MilestoneID, TenantID: row.TenantID, ProjectID: row.ProjectID, EngagementID: ptrFromUUID(row.EngagementID),
			WorkerProfileID: row.WorkerProfileID, Title: row.Title, MilestoneCode: ptrFromText(row.MilestoneCode), Amount: floatFromNumeric(row.Amount),
			CurrencyCode: row.CurrencyCode, AcceptedAt: ptrFromTimestamptz(row.AcceptedAt), ProjectName: row.ProjectName, ProjectCode: ptrFromText(row.ProjectCode),
			WorkerDisplayName: row.WorkerDisplayName, WorkerCode: ptrFromText(row.WorkerCode),
		})
	}
	return items
}

func numericFromFlexAmount(value float64) pgtype.Numeric {
	scaled := int64(math.Round(value * 100))
	return pgtype.Numeric{Int: big.NewInt(scaled), Exp: -2, Valid: true}
}
