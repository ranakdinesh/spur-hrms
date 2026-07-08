package postgres

import (
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapOvertimeRequest(row sqlc.HrmsOvertimeRequest) *domain.OvertimeRequest {
	return &domain.OvertimeRequest{
		ID:                   row.ID,
		TenantID:             row.TenantID,
		UserID:               row.UserID,
		WorkDate:             timeFromDate(row.WorkDate),
		RequestedMinutes:     row.RequestedMinutes,
		ApprovedMinutes:      int32PtrFromInt4(row.ApprovedMinutes),
		Reason:               ptrFromText(row.Reason),
		Status:               row.Status,
		ReviewedBy:           ptrFromUUID(row.ReviewedBy),
		ReviewedAt:           ptrFromTimestamptz(row.ReviewedAt),
		ReviewRemarks:        ptrFromText(row.ReviewRemarks),
		CalculationType:      row.CalculationType,
		RateMultiplier:       floatFromNumeric(row.RateMultiplier),
		PayrollComponentCode: ptrFromText(row.PayrollComponentCode),
		PayrollExportStatus:  row.PayrollExportStatus,
		PayrollExportedAt:    ptrFromTimestamptz(row.PayrollExportedAt),
		PayrollExportedBy:    ptrFromUUID(row.PayrollExportedBy),
		SourceAttendanceID:   ptrFromUUID(row.SourceAttendanceID),
		SourceSegmentID:      ptrFromUUID(row.SourceSegmentID),
		Metadata:             mapFromJSONBytes(row.Metadata),
		Inactive:             row.Inactive,
		CreatedAt:            timeFromTimestamptz(row.CreatedAt),
		CreatedBy:            ptrFromUUID(row.CreatedBy),
		UpdatedAt:            timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:            ptrFromUUID(row.UpdatedBy),
	}
}

func mapOvertimeRequests(rows []sqlc.HrmsOvertimeRequest) []*domain.OvertimeRequest {
	items := make([]*domain.OvertimeRequest, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapOvertimeRequest(row))
	}
	return items
}
