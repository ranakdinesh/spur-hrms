package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) ProvisionTenant(ctx context.Context, cmd ports.ProvisionTenantCommand) (*domain.TenantProvisioning, error) {
	provisioning, err := domain.NewTenantProvisioning(domain.TenantProvisioningInput{
		TenantID:             cmd.TenantID,
		CompanyName:          cmd.CompanyName,
		Subdomain:            cmd.Subdomain,
		MobileActivationCode: cmd.MobileActivationCode,
		AdminEmail:           cmd.AdminEmail,
		AdminName:            cmd.AdminName,
		TenantURL:            cmd.TenantURL,
	})
	if err != nil {
		s.logError("validate tenant provisioning", err, serviceTenantIDField(cmd.TenantID), serviceStringField("company_name", cmd.CompanyName), serviceStringField("subdomain", cmd.Subdomain))
		return nil, err
	}

	if err := s.system.RunAsSystem(ctx, func(txCtx context.Context) error {
		if err := s.validateTenantSubdomainAvailability(txCtx, provisioning.TenantID, provisioning.Subdomain, &provisioning.CompanyName); err != nil {
			return err
		}
		profile, err := domain.NewTenantProfile(provisioning.TenantID, provisioning.Subdomain, provisioning.MobileActivationCode, &provisioning.CompanyName, nil)
		if err != nil {
			return err
		}
		if _, err := s.profiles.UpsertTenantProfile(txCtx, profile); err != nil {
			return err
		}
		provisioning.Seeded["tenant_profile"] = 1

		branding, err := domain.NewTenantBranding(domain.TenantBrandingInput{TenantID: provisioning.TenantID, Subdomain: provisioning.Subdomain, DisplayName: &provisioning.CompanyName})
		if err != nil {
			return err
		}
		if _, err := s.branding.UpsertTenantBranding(txCtx, branding, cmd.ActorID, provisioning.MobileActivationCode); err != nil {
			return err
		}
		provisioning.Seeded["tenant_branding"] = 1

		seeded, err := s.seedTenantProvisioningDefaults(txCtx, provisioning.TenantID, cmd.ActorID)
		if err != nil {
			return err
		}
		for key, count := range seeded {
			provisioning.Seeded[key] = count
		}
		return nil
	}); err != nil {
		s.logError("provision tenant defaults", err, serviceTenantIDField(cmd.TenantID), serviceStringField("subdomain", provisioning.Subdomain))
		return nil, err
	}

	if err := s.sendTenantRegistrationEmail(ctx, provisioning); err != nil {
		s.logError("send tenant registration email", err, serviceTenantIDField(provisioning.TenantID), serviceStringField("subdomain", provisioning.Subdomain))
		return nil, err
	}
	s.log.Info().Str("tenant_id", provisioning.TenantID.String()).Str("subdomain", provisioning.Subdomain).Msg("hrms: tenant provisioning completed")
	return provisioning, nil
}

func (s *TenantService) GetTenantProvisioningStatus(ctx context.Context, tenantID uuid.UUID) (*ports.TenantProvisioningStatus, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate tenant provisioning status tenant", err)
		return nil, err
	}
	status := &ports.TenantProvisioningStatus{TenantID: tenantID}
	if err := s.system.RunAsSystem(ctx, func(txCtx context.Context) error {
		if profile, err := s.profiles.GetTenantProfile(txCtx, tenantID); err == nil {
			status.ProfileReady = true
			status.Subdomain = profile.Subdomain
			status.MobileActivationCode = profile.MobileActivationCode
		}
		if items, err := s.workingHours.ListWorkingHours(txCtx, tenantID); err == nil {
			status.WorkingHoursCount = len(items)
		}
		if items, err := s.financialYears.ListFinancialYears(txCtx, tenantID); err == nil {
			status.FinancialYearsCount = len(items)
		}
		if _, err := s.financialYears.GetActiveFinancialYear(txCtx, tenantID); err == nil {
			status.ActiveFinancialYearReady = true
		}
		if items, err := s.lookups.ListEmploymentTypes(txCtx, tenantID); err == nil {
			status.EmploymentTypesCount = len(items)
		}
		if items, err := s.lookups.ListMaritalStatuses(txCtx, tenantID); err == nil {
			status.MaritalStatusesCount = len(items)
		}
		if items, err := s.designationMasters.ListDesignationLevelCodes(txCtx, tenantID); err == nil {
			status.LevelCodesCount = len(items)
		}
		if items, err := s.designationMasters.ListDesignationSeniorityRanks(txCtx, tenantID); err == nil {
			status.SeniorityRanksCount = len(items)
		}
		if items, err := s.notifications.ListNotificationMasters(txCtx, tenantID); err == nil {
			status.NotificationMastersCount = len(items)
		}
		if _, err := s.payCycles.GetPayCycle(txCtx, tenantID); err == nil {
			status.PayCycleReady = true
		}
		if item, err := s.subscriptions.GetCurrentTenantSubscription(txCtx, tenantID); err == nil {
			status.CurrentSubscription = item
		}
		return nil
	}); err != nil {
		s.logError("get tenant provisioning status", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	return status, nil
}

func (s *TenantService) seedTenantProvisioningDefaults(ctx context.Context, tenantID uuid.UUID, actorID *uuid.UUID) (map[string]int, error) {
	seeded := map[string]int{}
	workingHours, err := s.workingHours.ListWorkingHours(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	if !hasTenantDefaultWorkingHours(workingHours) {
		for _, input := range domain.DefaultWorkingHourInputs(tenantID) {
			item, itemErr := domain.NewWorkingHour(input)
			if itemErr != nil {
				return nil, itemErr
			}
			if _, itemErr = s.workingHours.CreateWorkingHour(ctx, item, actorID); itemErr != nil {
				return nil, itemErr
			}
			seeded["working_hours"]++
		}
	}

	if count, err := s.seedActiveFinancialYear(ctx, tenantID, actorID); err != nil {
		return nil, err
	} else {
		seeded["financial_years"] = count
	}

	if count, err := s.seedEmploymentTypes(ctx, tenantID, actorID); err != nil {
		return nil, err
	} else {
		seeded["employment_types"] = count
	}
	if count, err := s.seedMaritalStatuses(ctx, tenantID, actorID); err != nil {
		return nil, err
	} else {
		seeded["marital_statuses"] = count
	}
	if count, err := s.seedDesignationLevelCodes(ctx, tenantID, actorID); err != nil {
		return nil, err
	} else {
		seeded["designation_level_codes"] = count
	}
	if count, err := s.seedDesignationSeniorityRanks(ctx, tenantID, actorID); err != nil {
		return nil, err
	} else {
		seeded["designation_seniority_ranks"] = count
	}
	if count, err := s.seedCelebrationTypes(ctx, tenantID, actorID); err != nil {
		return nil, err
	} else {
		seeded["celebration_types"] = count
	}
	if count, err := s.seedMissingNotificationMasters(ctx, tenantID, actorID); err != nil {
		return nil, err
	} else {
		seeded["notification_masters"] = count
	}
	if count, err := s.seedDefaultPayCycle(ctx, tenantID, actorID); err != nil {
		return nil, err
	} else {
		seeded["pay_cycle"] = count
	}
	return seeded, nil
}

func (s *TenantService) seedActiveFinancialYear(ctx context.Context, tenantID uuid.UUID, actorID *uuid.UUID) (int, error) {
	financialYears, err := s.financialYears.ListFinancialYears(ctx, tenantID)
	if err != nil {
		return 0, err
	}
	for _, item := range financialYears {
		if item != nil && item.IsActive {
			return 0, nil
		}
	}
	if len(financialYears) > 0 {
		for _, item := range financialYears {
			if item == nil || item.ID == uuid.Nil {
				continue
			}
			if _, err := s.financialYears.SetActiveFinancialYear(ctx, tenantID, item.ID, actorID); err != nil {
				return 0, err
			}
			return 1, nil
		}
	}
	fy, err := domain.NewFinancialYear(defaultProvisioningFinancialYearInput(tenantID))
	if err != nil {
		return 0, err
	}
	created, err := s.financialYears.CreateFinancialYear(ctx, fy, actorID)
	if err != nil {
		return 0, err
	}
	if _, err = s.financialYears.SetActiveFinancialYear(ctx, tenantID, created.ID, actorID); err != nil {
		return 0, err
	}
	return 1, nil
}

func (s *TenantService) seedDefaultPayCycle(ctx context.Context, tenantID uuid.UUID, actorID *uuid.UUID) (int, error) {
	if _, err := s.payCycles.GetPayCycle(ctx, tenantID); err == nil {
		return 0, nil
	} else if !errors.Is(err, domain.ErrPayCycleNotFound) {
		return 0, err
	}
	item, err := domain.NewPayCycle(domain.PayCycleInput{TenantID: tenantID})
	if err != nil {
		return 0, err
	}
	if _, err = s.payCycles.UpsertPayCycle(ctx, item, actorID); err != nil {
		return 0, err
	}
	return 1, nil
}

func (s *TenantService) seedEmploymentTypes(ctx context.Context, tenantID uuid.UUID, actorID *uuid.UUID) (int, error) {
	items, err := s.lookups.ListEmploymentTypes(ctx, tenantID)
	if err != nil || len(items) > 0 {
		return 0, err
	}
	count := 0
	for _, input := range domain.DefaultEmploymentTypeInputs(tenantID) {
		item, itemErr := domain.NewEmploymentType(input)
		if itemErr != nil {
			return count, itemErr
		}
		if _, itemErr = s.lookups.CreateEmploymentType(ctx, item, actorID); itemErr != nil {
			return count, itemErr
		}
		count++
	}
	return count, nil
}

func (s *TenantService) seedMaritalStatuses(ctx context.Context, tenantID uuid.UUID, actorID *uuid.UUID) (int, error) {
	items, err := s.lookups.ListMaritalStatuses(ctx, tenantID)
	if err != nil || len(items) > 0 {
		return 0, err
	}
	count := 0
	for _, input := range domain.DefaultMaritalStatusInputs(tenantID) {
		item, itemErr := domain.NewMaritalStatus(input)
		if itemErr != nil {
			return count, itemErr
		}
		if _, itemErr = s.lookups.CreateMaritalStatus(ctx, item, actorID); itemErr != nil {
			return count, itemErr
		}
		count++
	}
	return count, nil
}

func (s *TenantService) seedDesignationLevelCodes(ctx context.Context, tenantID uuid.UUID, actorID *uuid.UUID) (int, error) {
	items, err := s.designationMasters.ListDesignationLevelCodes(ctx, tenantID)
	if err != nil || len(items) > 0 {
		return 0, err
	}
	count := 0
	for _, input := range domain.DefaultDesignationLevelCodeInputs(tenantID) {
		item, itemErr := domain.NewDesignationLevelCode(input)
		if itemErr != nil {
			return count, itemErr
		}
		if _, itemErr = s.designationMasters.CreateDesignationLevelCode(ctx, item, actorID); itemErr != nil {
			return count, itemErr
		}
		count++
	}
	return count, nil
}

func (s *TenantService) seedDesignationSeniorityRanks(ctx context.Context, tenantID uuid.UUID, actorID *uuid.UUID) (int, error) {
	items, err := s.designationMasters.ListDesignationSeniorityRanks(ctx, tenantID)
	if err != nil || len(items) > 0 {
		return 0, err
	}
	count := 0
	for _, input := range domain.DefaultDesignationSeniorityRankInputs(tenantID) {
		item, itemErr := domain.NewDesignationSeniorityRank(input)
		if itemErr != nil {
			return count, itemErr
		}
		if _, itemErr = s.designationMasters.CreateDesignationSeniorityRank(ctx, item, actorID); itemErr != nil {
			return count, itemErr
		}
		count++
	}
	return count, nil
}

func defaultProvisioningFinancialYearInput(tenantID uuid.UUID) domain.FinancialYearInput {
	now := time.Now().UTC()
	startYear := now.Year()
	if now.Month() < time.April {
		startYear--
	}
	start := time.Date(startYear, time.April, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(startYear+1, time.March, 31, 0, 0, 0, 0, time.UTC)
	return domain.FinancialYearInput{TenantID: tenantID, StartDate: start, EndDate: end, PayrollYear: true, LeaveYear: true, HolidayYear: true, ReportingYear: true}
}

func (s *TenantService) sendTenantRegistrationEmail(ctx context.Context, provisioning *domain.TenantProvisioning) error {
	if s.registrationEmail == nil {
		s.log.Warn().Str("tenant_id", provisioning.TenantID.String()).Str("subdomain", provisioning.Subdomain).Msg("hrms: tenant registration email sender not configured")
		return nil
	}
	if provisioning.AdminEmail == nil || strings.TrimSpace(*provisioning.AdminEmail) == "" {
		s.log.Warn().Str("tenant_id", provisioning.TenantID.String()).Str("subdomain", provisioning.Subdomain).Msg("hrms: tenant registration email skipped because admin email is missing")
		return nil
	}
	message := ports.TenantRegistrationEmail{
		TenantID:             provisioning.TenantID,
		CompanyName:          provisioning.CompanyName,
		Subdomain:            provisioning.Subdomain,
		MobileActivationCode: provisioning.MobileActivationCode,
		AdminEmail:           *provisioning.AdminEmail,
	}
	if provisioning.AdminName != nil {
		message.AdminName = *provisioning.AdminName
	}
	if provisioning.TenantURL != nil {
		message.TenantURL = *provisioning.TenantURL
	}
	if err := s.registrationEmail.SendTenantRegistration(ctx, message); err != nil {
		return err
	}
	provisioning.EmailSent = true
	return nil
}
