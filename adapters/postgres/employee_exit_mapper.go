package postgres

import (
	"encoding/json"

	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapEmployeeExitRequest(row sqlc.HrmsEmployeeExitRequest) *domain.EmployeeExitRequest {
	return &domain.EmployeeExitRequest{
		ID:                     row.ID,
		TenantID:               row.TenantID,
		EmployeeID:             row.EmployeeID,
		EmployeeUserID:         row.EmployeeUserID,
		InitiatedBy:            ptrFromUUID(row.InitiatedBy),
		ApprovedBy:             ptrFromUUID(row.ApprovedBy),
		ApprovedAt:             ptrFromTimestamptz(row.ApprovedAt),
		CompletedBy:            ptrFromUUID(row.CompletedBy),
		CompletedAt:            ptrFromTimestamptz(row.CompletedAt),
		Status:                 row.Status,
		ExitType:               row.ExitType,
		Reason:                 ptrFromText(row.Reason),
		ResignationDate:        ptrFromDate(row.ResignationDate),
		NoticeStartDate:        ptrFromDate(row.NoticeStartDate),
		LastWorkingDate:        timeFromDate(row.LastWorkingDate),
		RequestedRelievingDate: ptrFromDate(row.RequestedRelievingDate),
		ApprovedRelievingDate:  ptrFromDate(row.ApprovedRelievingDate),
		FinalSettlementStatus:  row.FinalSettlementStatus,
		AccessRevocationStatus: row.AccessRevocationStatus,
		AssetClearanceStatus:   row.AssetClearanceStatus,
		HandoverStatus:         row.HandoverStatus,
		ExitInterviewStatus:    row.ExitInterviewStatus,
		Notes:                  ptrFromText(row.Notes),
		RejectionReason:        ptrFromText(row.RejectionReason),
		CancelReason:           ptrFromText(row.CancelReason),
		Inactive:               row.Inactive,
		CreatedAt:              timeFromTimestamptz(row.CreatedAt),
		CreatedBy:              ptrFromUUID(row.CreatedBy),
		UpdatedAt:              timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:              ptrFromUUID(row.UpdatedBy),
	}
}

func mapEmployeeExitListRow(row sqlc.ListEmployeeExitRequestsRow) *domain.EmployeeExitRequest {
	item := &domain.EmployeeExitRequest{
		ID:                     row.ID,
		TenantID:               row.TenantID,
		EmployeeID:             row.EmployeeID,
		EmployeeUserID:         row.EmployeeUserID,
		EmployeeFirstname:      &row.EmployeeFirstname,
		EmployeeLastname:       ptrFromText(row.EmployeeLastname),
		EmployeeCode:           ptrFromText(row.EmployeeCode),
		EmployeeEmail:          ptrFromText(row.EmployeeEmail),
		DepartmentName:         ptrFromText(row.DepartmentName),
		BranchName:             ptrFromText(row.BranchName),
		InitiatedBy:            ptrFromUUID(row.InitiatedBy),
		ApprovedBy:             ptrFromUUID(row.ApprovedBy),
		ApprovedAt:             ptrFromTimestamptz(row.ApprovedAt),
		CompletedBy:            ptrFromUUID(row.CompletedBy),
		CompletedAt:            ptrFromTimestamptz(row.CompletedAt),
		Status:                 row.Status,
		ExitType:               row.ExitType,
		Reason:                 ptrFromText(row.Reason),
		ResignationDate:        ptrFromDate(row.ResignationDate),
		NoticeStartDate:        ptrFromDate(row.NoticeStartDate),
		LastWorkingDate:        timeFromDate(row.LastWorkingDate),
		RequestedRelievingDate: ptrFromDate(row.RequestedRelievingDate),
		ApprovedRelievingDate:  ptrFromDate(row.ApprovedRelievingDate),
		FinalSettlementStatus:  row.FinalSettlementStatus,
		AccessRevocationStatus: row.AccessRevocationStatus,
		AssetClearanceStatus:   row.AssetClearanceStatus,
		HandoverStatus:         row.HandoverStatus,
		ExitInterviewStatus:    row.ExitInterviewStatus,
		Notes:                  ptrFromText(row.Notes),
		RejectionReason:        ptrFromText(row.RejectionReason),
		CancelReason:           ptrFromText(row.CancelReason),
		TotalTasks:             row.TotalTasks,
		CompletedTasks:         row.CompletedTasks,
		BlockedTasks:           row.BlockedTasks,
		Inactive:               row.Inactive,
		CreatedAt:              timeFromTimestamptz(row.CreatedAt),
		CreatedBy:              ptrFromUUID(row.CreatedBy),
		UpdatedAt:              timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:              ptrFromUUID(row.UpdatedBy),
	}
	return item
}

func mapEmployeeExitDetailRow(row sqlc.GetEmployeeExitRequestRow) *domain.EmployeeExitRequest {
	item := mapEmployeeExitRequest(sqlc.HrmsEmployeeExitRequest{
		ID: row.ID, TenantID: row.TenantID, EmployeeID: row.EmployeeID, EmployeeUserID: row.EmployeeUserID,
		InitiatedBy: row.InitiatedBy, ApprovedBy: row.ApprovedBy, ApprovedAt: row.ApprovedAt, CompletedBy: row.CompletedBy, CompletedAt: row.CompletedAt,
		Status: row.Status, ExitType: row.ExitType, Reason: row.Reason, ResignationDate: row.ResignationDate, NoticeStartDate: row.NoticeStartDate,
		LastWorkingDate: row.LastWorkingDate, RequestedRelievingDate: row.RequestedRelievingDate, ApprovedRelievingDate: row.ApprovedRelievingDate,
		FinalSettlementStatus: row.FinalSettlementStatus, AccessRevocationStatus: row.AccessRevocationStatus, AssetClearanceStatus: row.AssetClearanceStatus,
		HandoverStatus: row.HandoverStatus, ExitInterviewStatus: row.ExitInterviewStatus, Notes: row.Notes, RejectionReason: row.RejectionReason,
		CancelReason: row.CancelReason, Inactive: row.Inactive, CreatedAt: row.CreatedAt, CreatedBy: row.CreatedBy, UpdatedAt: row.UpdatedAt, UpdatedBy: row.UpdatedBy,
	})
	item.EmployeeFirstname = &row.EmployeeFirstname
	item.EmployeeLastname = ptrFromText(row.EmployeeLastname)
	item.EmployeeCode = ptrFromText(row.EmployeeCode)
	item.EmployeeEmail = ptrFromText(row.EmployeeEmail)
	item.DepartmentName = ptrFromText(row.DepartmentName)
	item.BranchName = ptrFromText(row.BranchName)
	return item
}

func mapEmployeeExitTask(row sqlc.HrmsEmployeeExitTask) *domain.EmployeeExitTask {
	return &domain.EmployeeExitTask{
		ID:             row.ID,
		TenantID:       row.TenantID,
		ExitRequestID:  row.ExitRequestID,
		EmployeeUserID: row.EmployeeUserID,
		TaskKey:        row.TaskKey,
		Title:          row.Title,
		Description:    ptrFromText(row.Description),
		OwnerRole:      ptrFromText(row.OwnerRole),
		OwnerUserID:    ptrFromUUID(row.OwnerUserID),
		DueDate:        ptrFromDate(row.DueDate),
		Status:         row.Status,
		CompletedBy:    ptrFromUUID(row.CompletedBy),
		CompletedAt:    ptrFromTimestamptz(row.CompletedAt),
		Remarks:        ptrFromText(row.Remarks),
		SortOrder:      row.SortOrder,
		Inactive:       row.Inactive,
		CreatedAt:      timeFromTimestamptz(row.CreatedAt),
		CreatedBy:      ptrFromUUID(row.CreatedBy),
		UpdatedAt:      timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:      ptrFromUUID(row.UpdatedBy),
	}
}

func mapEmployeeExitTasks(rows []sqlc.HrmsEmployeeExitTask) []*domain.EmployeeExitTask {
	items := make([]*domain.EmployeeExitTask, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapEmployeeExitTask(row))
	}
	return items
}

func mapEmployeeExitEvent(row sqlc.HrmsEmployeeExitEvent) *domain.EmployeeExitEvent {
	metadata := json.RawMessage(row.Metadata)
	return &domain.EmployeeExitEvent{
		ID:            row.ID,
		TenantID:      row.TenantID,
		ExitRequestID: row.ExitRequestID,
		ExitTaskID:    ptrFromUUID(row.ExitTaskID),
		Action:        row.Action,
		FromStatus:    ptrFromText(row.FromStatus),
		ToStatus:      ptrFromText(row.ToStatus),
		Remarks:       ptrFromText(row.Remarks),
		Metadata:      metadata,
		Inactive:      row.Inactive,
		CreatedAt:     timeFromTimestamptz(row.CreatedAt),
		CreatedBy:     ptrFromUUID(row.CreatedBy),
		UpdatedAt:     timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:     ptrFromUUID(row.UpdatedBy),
	}
}

func mapEmployeeExitEvents(rows []sqlc.HrmsEmployeeExitEvent) []*domain.EmployeeExitEvent {
	items := make([]*domain.EmployeeExitEvent, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapEmployeeExitEvent(row))
	}
	return items
}
