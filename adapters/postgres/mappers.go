package postgres

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapTenantProfile(row sqlc.HrmsTenantProfile) *domain.TenantProfile {
	return &domain.TenantProfile{
		TenantID:             row.TenantID,
		Subdomain:            row.Subdomain,
		MobileActivationCode: row.MobileActivationCode,
		DisplayName:          ptrFromText(row.DisplayName),
		LogoObjectKey:        ptrFromText(row.LogoObjectKey),
		CreatedAt:            timeFromTimestamptz(row.CreatedAt),
		UpdatedAt:            timeFromTimestamptz(row.UpdatedAt),
	}
}

func mapTenantSetting(row sqlc.HrmsTenantSetting) (*domain.TenantSetting, error) {
	item := &domain.TenantSetting{
		TenantID:  row.TenantID,
		Key:       row.Key,
		CreatedAt: timeFromTimestamptz(row.CreatedAt),
		UpdatedAt: timeFromTimestamptz(row.UpdatedAt),
		Value:     map[string]any{},
	}
	if len(row.Value) > 0 {
		if err := json.Unmarshal(row.Value, &item.Value); err != nil {
			return nil, fmt.Errorf("hrms: unmarshal tenant setting: %w", err)
		}
	}
	return item, nil
}

func mapBranch(row sqlc.HrmsBranch) *domain.Branch {
	return &domain.Branch{
		ID:                  row.ID,
		TenantID:            row.TenantID,
		Name:                row.BranchName,
		Address:             ptrFromText(row.Address),
		City:                ptrFromText(row.City),
		State:               ptrFromText(row.State),
		Country:             ptrFromText(row.Country),
		Pincode:             ptrFromText(row.Pincode),
		Phone:               ptrFromText(row.Phone),
		BranchManagerUserID: ptrFromUUID(row.BranchManagerUserID),
		HRUserID:            ptrFromUUID(row.HrUserID),
		AccountsUserID:      ptrFromUUID(row.AccountsUserID),
		Inactive:            row.Inactive,
		CreatedAt:           timeFromTimestamptz(row.CreatedAt),
		CreatedBy:           ptrFromUUID(row.CreatedBy),
		UpdatedAt:           timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:           ptrFromUUID(row.UpdatedBy),
	}
}

func mapBranches(rows []sqlc.HrmsBranch) []*domain.Branch {
	items := make([]*domain.Branch, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapBranch(row))
	}
	return items
}

func mapDepartment(row sqlc.HrmsDepartment) *domain.Department {
	return &domain.Department{
		ID:          row.ID,
		TenantID:    row.TenantID,
		Name:        row.Name,
		ShortCode:   row.ShortCode,
		Description: ptrFromText(row.Description),
		Inactive:    row.Inactive,
		CreatedAt:   timeFromTimestamptz(row.CreatedAt),
		CreatedBy:   ptrFromUUID(row.CreatedBy),
		UpdatedAt:   timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:   ptrFromUUID(row.UpdatedBy),
	}
}

func mapDepartments(rows []sqlc.HrmsDepartment) []*domain.Department {
	items := make([]*domain.Department, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapDepartment(row))
	}
	return items
}

func mapDesignation(row sqlc.HrmsDesignation) *domain.Designation {
	return &domain.Designation{
		ID:                 row.ID,
		TenantID:           row.TenantID,
		Name:               row.Name,
		LevelCode:          row.LevelCode,
		SeniorityRank:      row.SeniorityRank,
		Description:        ptrFromText(row.Description),
		AttendanceRequired: row.AttendanceRequired,
		Inactive:           row.Inactive,
		CreatedAt:          timeFromTimestamptz(row.CreatedAt),
		CreatedBy:          ptrFromUUID(row.CreatedBy),
		UpdatedAt:          timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:          ptrFromUUID(row.UpdatedBy),
	}
}

func mapDesignations(rows []sqlc.HrmsDesignation) []*domain.Designation {
	items := make([]*domain.Designation, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapDesignation(row))
	}
	return items
}

func mapDesignationLevelCode(row sqlc.HrmsDesignationLevelCode) *domain.DesignationLevelCode {
	return &domain.DesignationLevelCode{
		ID:          row.ID,
		TenantID:    row.TenantID,
		Code:        row.Code,
		Label:       row.Label,
		Description: ptrFromText(row.Description),
		SortOrder:   row.SortOrder,
		Inactive:    row.Inactive,
		CreatedAt:   timeFromTimestamptz(row.CreatedAt),
		CreatedBy:   ptrFromUUID(row.CreatedBy),
		UpdatedAt:   timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:   ptrFromUUID(row.UpdatedBy),
	}
}

func mapDesignationLevelCodes(rows []sqlc.HrmsDesignationLevelCode) []*domain.DesignationLevelCode {
	items := make([]*domain.DesignationLevelCode, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapDesignationLevelCode(row))
	}
	return items
}

func mapDesignationSeniorityRank(row sqlc.HrmsDesignationSeniorityRank) *domain.DesignationSeniorityRank {
	return &domain.DesignationSeniorityRank{
		ID:          row.ID,
		TenantID:    row.TenantID,
		RankValue:   row.RankValue,
		Label:       row.Label,
		Description: ptrFromText(row.Description),
		SortOrder:   row.SortOrder,
		Inactive:    row.Inactive,
		CreatedAt:   timeFromTimestamptz(row.CreatedAt),
		CreatedBy:   ptrFromUUID(row.CreatedBy),
		UpdatedAt:   timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:   ptrFromUUID(row.UpdatedBy),
	}
}

func mapDesignationSeniorityRanks(rows []sqlc.HrmsDesignationSeniorityRank) []*domain.DesignationSeniorityRank {
	items := make([]*domain.DesignationSeniorityRank, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapDesignationSeniorityRank(row))
	}
	return items
}

func mapEmploymentType(row sqlc.HrmsEmploymentType) *domain.EmploymentType {
	return &domain.EmploymentType{
		ID:        row.ID,
		TenantID:  row.TenantID,
		Name:      row.Name,
		Inactive:  row.Inactive,
		CreatedAt: timeFromTimestamptz(row.CreatedAt),
		CreatedBy: ptrFromUUID(row.CreatedBy),
		UpdatedAt: timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy: ptrFromUUID(row.UpdatedBy),
	}
}

func mapEmploymentTypes(rows []sqlc.HrmsEmploymentType) []*domain.EmploymentType {
	items := make([]*domain.EmploymentType, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapEmploymentType(row))
	}
	return items
}

func mapMaritalStatus(row sqlc.HrmsMaritalStatus) *domain.MaritalStatus {
	return &domain.MaritalStatus{
		ID:        row.ID,
		TenantID:  row.TenantID,
		Name:      row.Name,
		Inactive:  row.Inactive,
		CreatedAt: timeFromTimestamptz(row.CreatedAt),
		CreatedBy: ptrFromUUID(row.CreatedBy),
		UpdatedAt: timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy: ptrFromUUID(row.UpdatedBy),
	}
}

func mapMaritalStatuses(rows []sqlc.HrmsMaritalStatus) []*domain.MaritalStatus {
	items := make([]*domain.MaritalStatus, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapMaritalStatus(row))
	}
	return items
}

func mapWorkingHour(row sqlc.HrmsWorkingHour) *domain.WorkingHour {
	item := &domain.WorkingHour{
		ID:           row.ID,
		TenantID:     row.TenantID,
		BranchID:     ptrFromUUID(row.BranchID),
		UserID:       ptrFromUUID(row.UserID),
		DayOfWeek:    row.DayOfWeek,
		IsWorkingDay: row.IsWorkingDay,
		StartTime:    clockStringFromTime(row.StartTime),
		EndTime:      clockStringFromTime(row.EndTime),
		BreakMinutes: row.BreakMinutes,
		Inactive:     row.Inactive,
		CreatedAt:    timeFromTimestamptz(row.CreatedAt),
		CreatedBy:    ptrFromUUID(row.CreatedBy),
		UpdatedAt:    timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:    ptrFromUUID(row.UpdatedBy),
	}
	item.Source = domain.WorkingHourScope(item)
	return item
}

func mapWorkingHours(rows []sqlc.HrmsWorkingHour) []*domain.WorkingHour {
	items := make([]*domain.WorkingHour, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapWorkingHour(row))
	}
	return items
}

func mapFinancialYear(row sqlc.HrmsFinancialYear) *domain.FinancialYear {
	return &domain.FinancialYear{
		ID:            row.ID,
		TenantID:      row.TenantID,
		Name:          row.Name,
		StartDate:     timeFromDate(row.StartDate),
		EndDate:       timeFromDate(row.EndDate),
		IsActive:      row.IsActive,
		PayrollYear:   row.PayrollYear,
		LeaveYear:     row.LeaveYear,
		HolidayYear:   row.HolidayYear,
		ReportingYear: row.ReportingYear,
		IsLocked:      row.IsLocked,
		LockedAt:      ptrFromTimestamptz(row.LockedAt),
		LockedBy:      ptrFromUUID(row.LockedBy),
		CloseNote:     ptrFromText(row.CloseNote),
		Inactive:      row.Inactive,
		CreatedAt:     timeFromTimestamptz(row.CreatedAt),
		CreatedBy:     ptrFromUUID(row.CreatedBy),
		UpdatedAt:     timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:     ptrFromUUID(row.UpdatedBy),
	}
}

func mapFinancialYears(rows []sqlc.HrmsFinancialYear) []*domain.FinancialYear {
	items := make([]*domain.FinancialYear, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapFinancialYear(row))
	}
	return items
}

func mapHoliday(row sqlc.HrmsHoliday) *domain.Holiday {
	return &domain.Holiday{
		ID:         row.ID,
		TenantID:   row.TenantID,
		BranchID:   ptrFromUUID(row.BranchID),
		FYID:       ptrFromUUID(row.FyID),
		Name:       row.Name,
		Date:       timeFromDate(row.Date),
		IsOptional: row.IsOptional,
		Inactive:   row.Inactive,
		CreatedAt:  timeFromTimestamptz(row.CreatedAt),
		CreatedBy:  ptrFromUUID(row.CreatedBy),
		UpdatedAt:  timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:  ptrFromUUID(row.UpdatedBy),
	}
}

func mapHolidays(rows []sqlc.HrmsHoliday) []*domain.Holiday {
	items := make([]*domain.Holiday, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapHoliday(row))
	}
	return items
}

func textFromString(value string) pgtype.Text {
	return pgtype.Text{String: value, Valid: true}
}

func textFromPtr(value *string) pgtype.Text {
	if value == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *value, Valid: true}
}

func ptrFromText(value pgtype.Text) *string {
	if !value.Valid {
		return nil
	}
	return &value.String
}

func uuidFromPtr(value *uuid.UUID) pgtype.UUID {
	if value == nil || *value == uuid.Nil {
		return pgtype.UUID{Valid: false}
	}
	return pgtype.UUID{Bytes: *value, Valid: true}
}

func ptrFromUUID(value pgtype.UUID) *uuid.UUID {
	if !value.Valid {
		return nil
	}
	id := uuid.UUID(value.Bytes)
	return &id
}

func timeFromClockString(value string) pgtype.Time {
	minutes, err := clockStringToMinutes(value)
	if err != nil {
		return pgtype.Time{Valid: false}
	}
	return pgtype.Time{Microseconds: int64(minutes) * int64(time.Minute/time.Microsecond), Valid: true}
}

func clockStringFromTime(value pgtype.Time) string {
	if !value.Valid {
		return "00:00"
	}
	minutes := int(value.Microseconds / int64(time.Minute/time.Microsecond))
	return fmt.Sprintf("%02d:%02d", minutes/60, minutes%60)
}

func clockStringToMinutes(value string) (int, error) {
	parsed, err := time.Parse("15:04", value)
	if err != nil {
		return 0, err
	}
	return parsed.Hour()*60 + parsed.Minute(), nil
}

func dateFromTime(value time.Time) pgtype.Date {
	if value.IsZero() {
		return pgtype.Date{Valid: false}
	}
	return pgtype.Date{Time: time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, time.UTC), Valid: true}
}

func timeFromDate(value pgtype.Date) time.Time {
	if !value.Valid {
		return time.Time{}
	}
	return time.Date(value.Time.Year(), value.Time.Month(), value.Time.Day(), 0, 0, 0, 0, time.UTC)
}

func dateFromPtr(value *time.Time) pgtype.Date {
	if value == nil {
		return pgtype.Date{Valid: false}
	}
	return dateFromTime(*value)
}

func ptrFromDate(value pgtype.Date) *time.Time {
	if !value.Valid {
		return nil
	}
	date := timeFromDate(value)
	return &date
}

func timestamptzFromPtr(value *time.Time) pgtype.Timestamptz {
	if value == nil || value.IsZero() {
		return pgtype.Timestamptz{Valid: false}
	}
	return pgtype.Timestamptz{Time: value.UTC(), Valid: true}
}

func timeFromTimestamptz(value pgtype.Timestamptz) time.Time {
	if !value.Valid {
		return time.Time{}
	}
	return value.Time
}

func ptrFromTimestamptz(value pgtype.Timestamptz) *time.Time {
	if !value.Valid {
		return nil
	}
	return &value.Time
}

func mapTenantBrandingByTenantID(row sqlc.GetTenantBrandingByTenantIDRow) *domain.TenantBranding {
	return tenantBrandingFromParts(
		row.TenantID,
		row.Subdomain,
		ptrFromText(row.DisplayName),
		ptrFromText(row.LogoPath),
		ptrFromText(row.FaviconPath),
		row.Layout,
		row.ColorMode,
		row.SidebarSize,
		row.LayoutWidth,
		row.CardLayout,
		row.ThemeColor,
		row.PrimaryColor,
		row.SecondaryColor,
		row.TertiaryColor,
		row.TopbarColor,
		row.SidebarColor,
		row.TopbarBackground,
		row.SidebarBackground,
		row.FontFamily,
		row.Preloader,
	)
}

func mapTenantBrandingBySubdomain(row sqlc.ResolveTenantBrandingBySubdomainRow) *domain.TenantBranding {
	return tenantBrandingFromParts(
		row.TenantID,
		row.Subdomain,
		ptrFromText(row.DisplayName),
		ptrFromText(row.LogoPath),
		ptrFromText(row.FaviconPath),
		row.Layout,
		row.ColorMode,
		row.SidebarSize,
		row.LayoutWidth,
		row.CardLayout,
		row.ThemeColor,
		row.PrimaryColor,
		row.SecondaryColor,
		row.TertiaryColor,
		row.TopbarColor,
		row.SidebarColor,
		row.TopbarBackground,
		row.SidebarBackground,
		row.FontFamily,
		row.Preloader,
	)
}

func mapTenantBrandingUpsert(row sqlc.UpsertTenantBrandingRow) *domain.TenantBranding {
	return tenantBrandingFromParts(
		row.TenantID,
		row.Subdomain,
		ptrFromText(row.DisplayName),
		ptrFromText(row.LogoPath),
		ptrFromText(row.FaviconPath),
		row.Layout,
		row.ColorMode,
		row.SidebarSize,
		row.LayoutWidth,
		row.CardLayout,
		row.ThemeColor,
		row.PrimaryColor,
		row.SecondaryColor,
		row.TertiaryColor,
		row.TopbarColor,
		row.SidebarColor,
		row.TopbarBackground,
		row.SidebarBackground,
		row.FontFamily,
		row.Preloader,
	)
}

func tenantBrandingFromParts(
	tenantID uuid.UUID,
	subdomain string,
	displayName *string,
	logoPath *string,
	faviconPath *string,
	layout string,
	colorMode string,
	sidebarSize string,
	layoutWidth string,
	cardLayout string,
	themeColor string,
	primaryColor string,
	secondaryColor string,
	tertiaryColor string,
	topbarColor string,
	sidebarColor string,
	topbarBackground string,
	sidebarBackground string,
	fontFamily string,
	preloader bool,
) *domain.TenantBranding {
	return &domain.TenantBranding{
		TenantID:          tenantID,
		Subdomain:         subdomain,
		DisplayName:       displayName,
		LogoPath:          logoPath,
		FaviconPath:       faviconPath,
		Layout:            layout,
		ColorMode:         colorMode,
		SidebarSize:       sidebarSize,
		LayoutWidth:       layoutWidth,
		CardLayout:        cardLayout,
		ThemeColor:        themeColor,
		PrimaryColor:      primaryColor,
		SecondaryColor:    secondaryColor,
		TertiaryColor:     tertiaryColor,
		TopbarColor:       topbarColor,
		SidebarColor:      sidebarColor,
		TopbarBackground:  topbarBackground,
		SidebarBackground: sidebarBackground,
		FontFamily:        fontFamily,
		Preloader:         preloader,
	}
}
