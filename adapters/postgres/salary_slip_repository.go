package postgres

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateSalarySlip(ctx context.Context, item *domain.SalarySlip, actorID *uuid.UUID) (*domain.SalarySlip, error) {
	row, err := s.getQueries(ctx).CreateSalarySlip(ctx, salarySlipCreateParams(item, actorID))
	if err != nil {
		return nil, s.logDBError(ctx, "create salary slip", err, tenantIDField(item.TenantID), stringField("user_id", item.UserID.String()))
	}
	return mapSalarySlip(row), nil
}

func (s *Store) UpdateSalarySlip(ctx context.Context, item *domain.SalarySlip, actorID *uuid.UUID) (*domain.SalarySlip, error) {
	row, err := s.getQueries(ctx).UpdateSalarySlipByID(ctx, sqlc.UpdateSalarySlipByIDParams{TenantID: item.TenantID, ID: item.ID, FyID: item.FYID, TemplateID: item.TemplateID, GrossSalary: numericFromFloat(item.GrossSalary), TotalEarnings: numericFromFloat(item.TotalEarnings), TotalDeductions: numericFromFloat(item.TotalDeductions), AbsentDeduction: numericFromFloat(item.AbsentDeduction), NetSalary: numericFromFloat(item.NetSalary), AbsentDays: item.AbsentDays, PresentDays: item.PresentDays, TotalDays: item.TotalDays, LwpDays: numericFromFloat(item.LWPDays), NoOfPhWeo: item.NoOfPHWEO, IsSpecial: item.IsSpecial, IsRegenerated: item.IsRegenerated, PdfPath: textFromPtr(item.PDFPath), UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "update salary slip", err, tenantIDField(item.TenantID), stringField("salary_slip_id", item.ID.String()))
	}
	return mapSalarySlip(row), nil
}

func salarySlipCreateParams(item *domain.SalarySlip, actorID *uuid.UUID) sqlc.CreateSalarySlipParams {
	return sqlc.CreateSalarySlipParams{TenantID: item.TenantID, UserID: item.UserID, FyID: item.FYID, TemplateID: item.TemplateID, Month: item.Month, Year: item.Year, GrossSalary: numericFromFloat(item.GrossSalary), TotalEarnings: numericFromFloat(item.TotalEarnings), TotalDeductions: numericFromFloat(item.TotalDeductions), AbsentDeduction: numericFromFloat(item.AbsentDeduction), NetSalary: numericFromFloat(item.NetSalary), AbsentDays: item.AbsentDays, PresentDays: item.PresentDays, TotalDays: item.TotalDays, LwpDays: numericFromFloat(item.LWPDays), NoOfPhWeo: item.NoOfPHWEO, IsSpecial: item.IsSpecial, IsRegenerated: item.IsRegenerated, PdfPath: textFromPtr(item.PDFPath), CreatedBy: uuidFromPtr(actorID)}
}

func (s *Store) ListSalarySlipsByUser(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) ([]*domain.SalarySlip, error) {
	rows, err := s.getQueries(ctx).ListSalarySlipsByUser(ctx, sqlc.ListSalarySlipsByUserParams{TenantID: tenantID, UserID: userID})
	if err != nil {
		return nil, s.logDBError(ctx, "list salary slips by user", err, tenantIDField(tenantID), stringField("user_id", userID.String()))
	}
	return mapSalarySlips(rows), nil
}

func (s *Store) ListRecentSalarySlipsByUser(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, limit int32) ([]*domain.SalarySlip, error) {
	rows, err := s.getQueries(ctx).ListRecentSalarySlipsByUser(ctx, sqlc.ListRecentSalarySlipsByUserParams{TenantID: tenantID, UserID: userID, Limit: limit})
	if err != nil {
		return nil, s.logDBError(ctx, "list recent salary slips by user", err, tenantIDField(tenantID), stringField("user_id", userID.String()))
	}
	return mapSalarySlips(rows), nil
}

func (s *Store) ListSalarySlipsByTenantPeriod(ctx context.Context, tenantID uuid.UUID, month int32, year int32) ([]*domain.SalarySlip, error) {
	rows, err := s.getQueries(ctx).ListSalarySlipsByTenantPeriod(ctx, sqlc.ListSalarySlipsByTenantPeriodParams{TenantID: tenantID, Month: month, Year: year})
	if err != nil {
		return nil, s.logDBError(ctx, "list salary slips by tenant period", err, tenantIDField(tenantID))
	}
	return mapSalarySlips(rows), nil
}

func (s *Store) GetSalarySlip(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.SalarySlip, error) {
	row, err := s.getQueries(ctx).GetSalarySlip(ctx, sqlc.GetSalarySlipParams{TenantID: tenantID, ID: id})
	if err != nil {
		return nil, s.logDBError(ctx, "get salary slip", err, tenantIDField(tenantID), stringField("salary_slip_id", id.String()))
	}
	return mapSalarySlip(row), nil
}

func (s *Store) GetSalarySlipByPeriod(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, month int32, year int32) (*domain.SalarySlip, error) {
	row, err := s.getQueries(ctx).GetSalarySlipByPeriod(ctx, sqlc.GetSalarySlipByPeriodParams{TenantID: tenantID, UserID: userID, Month: month, Year: year})
	if err != nil {
		return nil, s.logDBError(ctx, "get salary slip by period", err, tenantIDField(tenantID), stringField("user_id", userID.String()))
	}
	return mapSalarySlip(row), nil
}

func (s *Store) CreateSalarySlipItem(ctx context.Context, item *domain.SalarySlipItem, actorID *uuid.UUID) (*domain.SalarySlipItem, error) {
	row, err := s.getQueries(ctx).CreateSalarySlipItem(ctx, sqlc.CreateSalarySlipItemParams{TenantID: item.TenantID, SlipID: item.SlipID, ItemType: item.ItemType, Code: item.Code, Name: item.Name, Amount: numericFromFloat(item.Amount), SortOrder: item.SortOrder, CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create salary slip item", err, tenantIDField(item.TenantID), stringField("salary_slip_id", item.SlipID.String()))
	}
	return mapSalarySlipItem(row), nil
}

func (s *Store) CreateSalarySlipLeave(ctx context.Context, item *domain.SalarySlipLeave, actorID *uuid.UUID) (*domain.SalarySlipLeave, error) {
	row, err := s.getQueries(ctx).CreateSalarySlipLeave(ctx, sqlc.CreateSalarySlipLeaveParams{TenantID: item.TenantID, SlipID: item.SlipID, LeaveTypeID: item.LeaveTypeID, LeaveTypeName: textFromPtr(item.LeaveTypeName), TotalDays: numericFromFloat(item.TotalDays), UsedDays: numericFromFloat(item.UsedDays), BalanceDays: numericFromFloat(item.BalanceDays), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create salary slip leave", err, tenantIDField(item.TenantID), stringField("salary_slip_id", item.SlipID.String()))
	}
	return mapSalarySlipLeave(row), nil
}

func (s *Store) ListSalarySlipItems(ctx context.Context, tenantID uuid.UUID, slipID uuid.UUID) ([]*domain.SalarySlipItem, error) {
	rows, err := s.getQueries(ctx).ListSalarySlipItems(ctx, sqlc.ListSalarySlipItemsParams{TenantID: tenantID, SlipID: slipID})
	if err != nil {
		return nil, s.logDBError(ctx, "list salary slip items", err, tenantIDField(tenantID), stringField("salary_slip_id", slipID.String()))
	}
	return mapSalarySlipItems(rows), nil
}

func (s *Store) ListSalarySlipLeaves(ctx context.Context, tenantID uuid.UUID, slipID uuid.UUID) ([]*domain.SalarySlipLeave, error) {
	rows, err := s.getQueries(ctx).ListSalarySlipLeaves(ctx, sqlc.ListSalarySlipLeavesParams{TenantID: tenantID, SlipID: slipID})
	if err != nil {
		return nil, s.logDBError(ctx, "list salary slip leaves", err, tenantIDField(tenantID), stringField("salary_slip_id", slipID.String()))
	}
	return mapSalarySlipLeaves(rows), nil
}

func (s *Store) DeleteSalarySlipItemsBySlip(ctx context.Context, tenantID uuid.UUID, slipID uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteSalarySlipItemsBySlip(ctx, sqlc.SoftDeleteSalarySlipItemsBySlipParams{TenantID: tenantID, SlipID: slipID, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete salary slip items by slip", err, tenantIDField(tenantID), stringField("salary_slip_id", slipID.String()))
	}
	return nil
}

func (s *Store) DeleteSalarySlipLeavesBySlip(ctx context.Context, tenantID uuid.UUID, slipID uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteSalarySlipLeavesBySlip(ctx, sqlc.SoftDeleteSalarySlipLeavesBySlipParams{TenantID: tenantID, SlipID: slipID, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete salary slip leaves by slip", err, tenantIDField(tenantID), stringField("salary_slip_id", slipID.String()))
	}
	return nil
}

func (s *Store) GetSalarySlipFormat(ctx context.Context, tenantID uuid.UUID) (*domain.SalarySlipFormat, error) {
	row, err := s.getQueries(ctx).GetSalarySlipFormat(ctx, tenantID)
	if err != nil {
		return nil, s.logDBError(ctx, "get salary slip format", err, tenantIDField(tenantID))
	}
	item, err := mapSalarySlipFormat(row)
	if err != nil {
		return nil, s.logDBError(ctx, "map salary slip format", err, tenantIDField(tenantID))
	}
	return item, nil
}

func (s *Store) UpsertSalarySlipFormat(ctx context.Context, item *domain.SalarySlipFormat, actorID *uuid.UUID) (*domain.SalarySlipFormat, error) {
	customFields, err := json.Marshal(item.CustomFields)
	if err != nil {
		return nil, s.logDBError(ctx, "marshal salary slip custom fields", err, tenantIDField(item.TenantID))
	}
	row, err := s.getQueries(ctx).UpsertSalarySlipFormat(ctx, sqlc.UpsertSalarySlipFormatParams{TenantID: item.TenantID, Title: item.Title, Subtitle: textFromPtr(item.Subtitle), LogoPath: textFromPtr(item.LogoPath), PrimaryColor: item.PrimaryColor, AccentColor: item.AccentColor, ShowLeaveBalance: item.ShowLeaveBalance, ShowYtdSummary: item.ShowYTDSummary, ShowEmployeeBank: item.ShowEmployeeBank, ShowEmployerContributions: item.ShowEmployerContributions, FooterText: textFromPtr(item.FooterText), CustomFields: customFields, UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "upsert salary slip format", err, tenantIDField(item.TenantID))
	}
	return mapSalarySlipFormat(row)
}
