package postgres

import (
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapCompOffRequest(row sqlc.HrmsCompOffRequest) *domain.CompOffRequest {
	return &domain.CompOffRequest{
		ID:                 row.ID,
		TenantID:           row.TenantID,
		UserID:             row.UserID,
		LeaveTypeID:        row.LeaveTypeID,
		FYID:               row.FyID,
		WorkDate:           timeFromDate(row.WorkDate),
		WorkedMinutes:      row.WorkedMinutes,
		RequestedDays:      floatFromNumeric(row.RequestedDays),
		ApprovedDays:       floatPtrFromNumeric(row.ApprovedDays),
		ExpiryDate:         ptrFromDate(row.ExpiryDate),
		Reason:             ptrFromText(row.Reason),
		Status:             row.Status,
		ReviewedBy:         ptrFromUUID(row.ReviewedBy),
		ReviewedAt:         ptrFromTimestamptz(row.ReviewedAt),
		ReviewRemarks:      ptrFromText(row.ReviewRemarks),
		PayrollImpact:      row.PayrollImpact,
		SourceAttendanceID: ptrFromUUID(row.SourceAttendanceID),
		SourceSegmentID:    ptrFromUUID(row.SourceSegmentID),
		Metadata:           mapFromJSONBytes(row.Metadata),
		Inactive:           row.Inactive,
		CreatedAt:          timeFromTimestamptz(row.CreatedAt),
		CreatedBy:          ptrFromUUID(row.CreatedBy),
		UpdatedAt:          timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:          ptrFromUUID(row.UpdatedBy),
	}
}

func mapCompOffRequests(rows []sqlc.HrmsCompOffRequest) []*domain.CompOffRequest {
	items := make([]*domain.CompOffRequest, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapCompOffRequest(row))
	}
	return items
}
