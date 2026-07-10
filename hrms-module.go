// Package hrms provides the Hrms Spur module.
//
// Install:
//
//	spur add module hrms
//
// Or manually:
//
//	go get github.com/ranakdinesh/spur-hrms@latest
//
// Wire in app.go:
//
//	hrmsModule, err := hrms.New(ctx, hrms.Options{DB: dbPool, Log: log, Cfg: cfg.Hrms})
//	hrmsModule.RegisterRoutes(r)
package hrms

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"
	"sort"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ranakdinesh/spur-hrms/adapters/communication"
	"github.com/ranakdinesh/spur-hrms/adapters/email"
	"github.com/ranakdinesh/spur-hrms/adapters/httpx"
	"github.com/ranakdinesh/spur-hrms/adapters/httpx/handlers"
	"github.com/ranakdinesh/spur-hrms/adapters/pdf"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres"
	"github.com/ranakdinesh/spur-hrms/adapters/push"
	storageadapter "github.com/ranakdinesh/spur-hrms/adapters/storage"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
	"github.com/ranakdinesh/spur-hrms/core/services"
	"github.com/ranakdinesh/spur-hrms/internal/logging"
	"github.com/ranakdinesh/spur-hrms/pkg/permissions"
	"github.com/ranakdinesh/spur-hrms/sql/migrations"
	templatemodule "github.com/ranakdinesh/spur-template/pkg/module"
	"github.com/rs/zerolog"
)

// Config holds all Hrms configuration from environment variables.
type Config struct {
	EmailProvider              string
	EmailFromName              string
	EmailFromEmail             string
	EmailReplyToEmail          string
	SMTPHost                   string
	SMTPPort                   int32
	SMTPUsername               string
	SMTPPassword               string
	SMTPEncryption             string
	SendGridAPIKey             string
	SendGridSandboxMode        bool
	EmailWebhookSigningSecret  string
	SMSProvider                string
	SMSEnabled                 bool
	SMSSenderID                string
	SMSAuthKey                 string
	SMSTemplateID              string
	SMSRoute                   string
	SMSCountryCode             string
	SMSBaseURL                 string
	WhatsAppProvider           string
	WhatsAppEnabled            bool
	WhatsAppAuthKey            string
	WhatsAppAppName            string
	WhatsAppSourceNumber       string
	WhatsAppTemplateID         string
	WhatsAppTemplateName       string
	WhatsAppNamespace          string
	WhatsAppBaseURL            string
	CommunicationWebhookSecret string
	StorageProvider            string
	StorageEnabled             bool
	StorageBucket              string
	StorageRegion              string
	StorageEndpoint            string
	StorageAccessKeyID         string
	StorageSecretAccessKey     string
	StorageUseSSL              bool
	StorageForcePathStyle      bool
	StorageObjectPrefix        string
	StoragePublicBaseURL       string
	StorageMaxFileSizeBytes    int64
	StorageAllowedContentTypes string
	PushProvider               string
	PushEnabled                bool
	PushProjectID              string
	PushClientEmail            string
	PushPrivateKey             string
	PushPrivateKeyID           string
	PushAuthURI                string
	PushTokenURI               string
	PushAndroidEnabled         bool
	PushIOSEnabled             bool
	PushWebEnabled             bool
	PushDefaultClickAction     string
	PushDefaultImageURL        string
	PushTTLSeconds             int32
	PushCollapseKey            string
}

// Options is passed by app.go when constructing this module.
type Options struct {
	DB                     *pgxpool.Pool
	Log                    *zerolog.Logger
	Cfg                    Config
	TenantIDFromContext    func(context.Context) string
	UserIDFromContext      func(context.Context) uuid.UUID
	IsSuperAdmin           func(context.Context) bool
	RolesFromContext       func(context.Context) []string
	PermissionsFromContext func(context.Context) []string
	EmployeeIdentity       ports.EmployeeIdentityPort
	LegacyPassword         ports.LegacyPasswordMigrationPort
	PolicyStorage          ports.PolicyFileStorage
	DocumentStorage        ports.EmployeeDocumentStorage
	SalarySlipPDF          ports.SalarySlipPDFRenderer
	SalarySlipStorage      ports.SalarySlipStorage
	EmployeeLetterPDF      ports.EmployeeLetterPDFRenderer
	EmployeeLetterStorage  ports.EmployeeLetterStorage
	AgreementPDF           ports.AgreementPDFRenderer
	AgreementStorage       ports.AgreementStorage

	// MigrationRunner runs this module's SQL migrations.
	// Provided by the platform: infra.Migrations.Run
	MigrationRunner func(ctx context.Context, moduleName string, fs interface{}) error
}

// Module is the Hrms module entry point.
type Module struct {
	// Services exposes this module's service interfaces to other modules.
	Services *Services
	Manifest templatemodule.Manifest

	handler *handlers.Handler
}

// Services bundles the public service interfaces.
type Services struct {
	Hrms ports.TenantService
}

// New wires the Hrms module. Returns error — never panics.
func New(ctx context.Context, opt Options) (*Module, error) {
	if opt.DB == nil {
		return nil, fmt.Errorf("hrms: DB pool is required")
	}
	log := logging.Ensure(opt.Log)

	if opt.MigrationRunner != nil {
		if err := opt.MigrationRunner(ctx, "hrms", migrations.FS); err != nil {
			log.Error().Err(err).Str("operation", "run managed migrations").Msg("hrms module startup failed")
			return nil, fmt.Errorf("hrms: migrations: %w", err)
		}
	} else if err := runMigrations(ctx, opt.DB); err != nil {
		log.Error().Err(err).Str("operation", "run embedded migrations").Msg("hrms module startup failed")
		return nil, fmt.Errorf("hrms: migrations: %w", err)
	}
	if err := ensureTenantBrandingSchema(ctx, opt.DB); err != nil {
		log.Error().Err(err).Str("operation", "ensure tenant branding schema").Msg("hrms module startup failed")
		return nil, fmt.Errorf("hrms: tenant branding schema: %w", err)
	}

	// Wire repo -> service -> handler
	store := postgres.New(opt.DB, log)
	serviceOptions := []services.TenantServiceOption{}
	objectStorage := storageadapter.NewObjectStorage(log)
	serviceOptions = append(serviceOptions, services.WithObjectStorage(objectStorage))
	serviceOptions = append(serviceOptions, services.WithEmailDeliverySender(email.NewMultiProviderSender(log)))
	serviceOptions = append(serviceOptions, services.WithCommunicationDeliverySender(communication.NewMultiProviderSender(log)))
	serviceOptions = append(serviceOptions, services.WithPushDeliverySender(push.NewMultiProviderSender(log)))
	if globalEmailProviderConfigured(opt.Cfg) {
		serviceOptions = append(serviceOptions, services.WithGlobalEmailProviderOnly(true))
	}
	if defaultEmail := defaultEmailProviderFromConfig(opt.Cfg); defaultEmail != nil {
		serviceOptions = append(serviceOptions, services.WithDefaultEmailProvider(defaultEmail))
	}
	if defaultCommunication := defaultCommunicationProviderFromConfig(opt.Cfg); defaultCommunication != nil {
		serviceOptions = append(serviceOptions, services.WithDefaultCommunicationProvider(defaultCommunication))
	}
	defaultStorage := defaultStorageProviderFromConfig(opt.Cfg)
	if defaultStorage != nil {
		serviceOptions = append(serviceOptions, services.WithDefaultStorageProvider(defaultStorage))
	}
	if defaultPush := defaultPushProviderFromConfig(opt.Cfg); defaultPush != nil {
		serviceOptions = append(serviceOptions, services.WithDefaultPushProvider(defaultPush))
	}
	tenantStorage := storageadapter.NewTenantFileStorage(store, objectStorage, defaultStorage, log)
	if opt.EmployeeIdentity != nil {
		serviceOptions = append(serviceOptions, services.WithEmployeeIdentityPort(opt.EmployeeIdentity))
	}
	if opt.LegacyPassword != nil {
		serviceOptions = append(serviceOptions, services.WithLegacyPasswordMigrationPort(opt.LegacyPassword))
	}
	if opt.PolicyStorage != nil {
		serviceOptions = append(serviceOptions, services.WithPolicyFileStorage(opt.PolicyStorage))
	} else {
		serviceOptions = append(serviceOptions, services.WithPolicyFileStorage(tenantStorage))
	}
	if opt.DocumentStorage != nil {
		serviceOptions = append(serviceOptions, services.WithEmployeeDocumentStorage(opt.DocumentStorage))
	} else {
		serviceOptions = append(serviceOptions, services.WithEmployeeDocumentStorage(tenantStorage))
	}
	serviceOptions = append(serviceOptions, services.WithHRCaseAttachmentStorage(tenantStorage))
	serviceOptions = append(serviceOptions, services.WithLearningCertificateStorage(tenantStorage))
	if opt.SalarySlipPDF != nil {
		serviceOptions = append(serviceOptions, services.WithSalarySlipPDFRenderer(opt.SalarySlipPDF))
	} else {
		serviceOptions = append(serviceOptions, services.WithSalarySlipPDFRenderer(pdf.NewSalarySlipRenderer()))
	}
	if opt.SalarySlipStorage != nil {
		serviceOptions = append(serviceOptions, services.WithSalarySlipStorage(opt.SalarySlipStorage))
	} else {
		serviceOptions = append(serviceOptions, services.WithSalarySlipStorage(tenantStorage))
	}
	if opt.EmployeeLetterPDF != nil {
		serviceOptions = append(serviceOptions, services.WithEmployeeLetterPDFRenderer(opt.EmployeeLetterPDF))
	} else {
		serviceOptions = append(serviceOptions, services.WithEmployeeLetterPDFRenderer(pdf.NewEmployeeLetterRenderer()))
	}
	if opt.EmployeeLetterStorage != nil {
		serviceOptions = append(serviceOptions, services.WithEmployeeLetterStorage(opt.EmployeeLetterStorage))
	} else {
		serviceOptions = append(serviceOptions, services.WithEmployeeLetterStorage(tenantStorage))
	}
	if opt.AgreementPDF != nil {
		serviceOptions = append(serviceOptions, services.WithAgreementPDFRenderer(opt.AgreementPDF))
	} else {
		serviceOptions = append(serviceOptions, services.WithAgreementPDFRenderer(pdf.NewAgreementRenderer()))
	}
	if opt.AgreementStorage != nil {
		serviceOptions = append(serviceOptions, services.WithAgreementStorage(opt.AgreementStorage))
	} else {
		serviceOptions = append(serviceOptions, services.WithAgreementStorage(tenantStorage))
	}
	svc := services.NewTenantService(store, log, serviceOptions...)
	h := handlers.New(svc, opt.TenantIDFromContext, opt.UserIDFromContext, opt.IsSuperAdmin, opt.RolesFromContext, opt.PermissionsFromContext, log)

	log.Info().Str("module", "hrms").Msg("Hrms module initialised")

	return &Module{
		Services: &Services{Hrms: svc},
		Manifest: permissions.Manifest(),
		handler:  h,
	}, nil
}

func defaultEmailProviderFromConfig(cfg Config) *domain.EmailProviderSettings {
	fromEmail := strings.TrimSpace(cfg.EmailFromEmail)
	if fromEmail == "" {
		return nil
	}
	provider := strings.TrimSpace(cfg.EmailProvider)
	if provider == "" {
		if strings.TrimSpace(cfg.SendGridAPIKey) != "" {
			provider = domain.EmailProviderSendGrid
		} else {
			provider = domain.EmailProviderSMTP
		}
	}
	port := cfg.SMTPPort
	var portPtr *int32
	if port > 0 {
		portPtr = &port
	}
	input := domain.EmailProviderSettingsInput{
		TenantID: uuid.New(), Provider: provider, IsEnabled: true, FromName: stringPtrOrNil(cfg.EmailFromName), FromEmail: fromEmail, ReplyToEmail: stringPtrOrNil(cfg.EmailReplyToEmail),
		SMTPHost: stringPtrOrNil(cfg.SMTPHost), SMTPPort: portPtr, SMTPUsername: stringPtrOrNil(cfg.SMTPUsername), SMTPPassword: stringPtrOrNil(cfg.SMTPPassword), SMTPEncryption: cfg.SMTPEncryption,
		SendGridAPIKey: stringPtrOrNil(cfg.SendGridAPIKey), SendGridSandboxMode: cfg.SendGridSandboxMode, WebhookSigningSecret: stringPtrOrNil(cfg.EmailWebhookSigningSecret),
	}
	item, err := domain.NewEmailProviderSettings(input)
	if err != nil {
		return nil
	}
	return item
}

func globalEmailProviderConfigured(cfg Config) bool {
	return strings.TrimSpace(cfg.EmailProvider) != "" ||
		strings.TrimSpace(cfg.EmailFromEmail) != "" ||
		strings.TrimSpace(cfg.SendGridAPIKey) != "" ||
		strings.TrimSpace(cfg.SMTPHost) != ""
}

func defaultCommunicationProviderFromConfig(cfg Config) *domain.CommunicationProviderSettings {
	if !cfg.SMSEnabled && !cfg.WhatsAppEnabled && strings.TrimSpace(cfg.SMSAuthKey) == "" && strings.TrimSpace(cfg.WhatsAppAuthKey) == "" {
		return nil
	}
	input := domain.CommunicationProviderSettingsInput{
		TenantID:             uuid.New(),
		SMSProvider:          cfg.SMSProvider,
		SMSEnabled:           cfg.SMSEnabled,
		SMSSenderID:          stringPtrOrNil(cfg.SMSSenderID),
		SMSAuthKey:           stringPtrOrNil(cfg.SMSAuthKey),
		SMSTemplateID:        stringPtrOrNil(cfg.SMSTemplateID),
		SMSRoute:             stringPtrOrNil(cfg.SMSRoute),
		SMSCountryCode:       stringPtrOrNil(cfg.SMSCountryCode),
		SMSBaseURL:           stringPtrOrNil(cfg.SMSBaseURL),
		WhatsAppProvider:     cfg.WhatsAppProvider,
		WhatsAppEnabled:      cfg.WhatsAppEnabled,
		WhatsAppAuthKey:      stringPtrOrNil(cfg.WhatsAppAuthKey),
		WhatsAppAppName:      stringPtrOrNil(cfg.WhatsAppAppName),
		WhatsAppSourceNumber: stringPtrOrNil(cfg.WhatsAppSourceNumber),
		WhatsAppTemplateID:   stringPtrOrNil(cfg.WhatsAppTemplateID),
		WhatsAppTemplateName: stringPtrOrNil(cfg.WhatsAppTemplateName),
		WhatsAppNamespace:    stringPtrOrNil(cfg.WhatsAppNamespace),
		WhatsAppBaseURL:      stringPtrOrNil(cfg.WhatsAppBaseURL),
		WebhookSigningSecret: stringPtrOrNil(cfg.CommunicationWebhookSecret),
	}
	item, err := domain.NewCommunicationProviderSettings(input)
	if err != nil {
		return nil
	}
	return item
}

func defaultStorageProviderFromConfig(cfg Config) *domain.StorageProviderSettings {
	if strings.TrimSpace(cfg.StorageBucket) == "" {
		return nil
	}
	input := domain.StorageProviderSettingsInput{
		TenantID:            uuid.New(),
		Provider:            cfg.StorageProvider,
		IsEnabled:           cfg.StorageEnabled,
		Bucket:              cfg.StorageBucket,
		Region:              stringPtrOrNil(cfg.StorageRegion),
		Endpoint:            stringPtrOrNil(cfg.StorageEndpoint),
		AccessKeyID:         stringPtrOrNil(cfg.StorageAccessKeyID),
		SecretAccessKey:     stringPtrOrNil(cfg.StorageSecretAccessKey),
		UseSSL:              cfg.StorageUseSSL,
		ForcePathStyle:      cfg.StorageForcePathStyle,
		ObjectPrefix:        stringPtrOrNil(cfg.StorageObjectPrefix),
		PublicBaseURL:       stringPtrOrNil(cfg.StoragePublicBaseURL),
		MaxFileSizeBytes:    cfg.StorageMaxFileSizeBytes,
		AllowedContentTypes: stringPtrOrNil(cfg.StorageAllowedContentTypes),
	}
	item, err := domain.NewStorageProviderSettings(input)
	if err != nil {
		return nil
	}
	return item
}

func defaultPushProviderFromConfig(cfg Config) *domain.PushProviderSettings {
	if !cfg.PushEnabled && strings.TrimSpace(cfg.PushProjectID) == "" && strings.TrimSpace(cfg.PushClientEmail) == "" {
		return nil
	}
	input := domain.PushProviderSettingsInput{
		TenantID: uuid.New(), Provider: cfg.PushProvider, IsEnabled: cfg.PushEnabled, ProjectID: stringPtrOrNil(cfg.PushProjectID), ClientEmail: stringPtrOrNil(cfg.PushClientEmail), PrivateKey: stringPtrOrNil(cfg.PushPrivateKey), PrivateKeyID: stringPtrOrNil(cfg.PushPrivateKeyID), AuthURI: stringPtrOrNil(cfg.PushAuthURI), TokenURI: stringPtrOrNil(cfg.PushTokenURI),
		AndroidEnabled: cfg.PushAndroidEnabled, IOSEnabled: cfg.PushIOSEnabled, WebEnabled: cfg.PushWebEnabled, DefaultClickAction: stringPtrOrNil(cfg.PushDefaultClickAction), DefaultImageURL: stringPtrOrNil(cfg.PushDefaultImageURL), TTLSeconds: cfg.PushTTLSeconds, CollapseKey: stringPtrOrNil(cfg.PushCollapseKey),
	}
	item, err := domain.NewPushProviderSettings(input)
	if err != nil {
		return nil
	}
	return item
}

func stringPtrOrNil(value string) *string {
	clean := strings.TrimSpace(value)
	if clean == "" {
		return nil
	}
	return &clean
}

// GetManifest returns the identity registration manifest for this module.
func (m *Module) GetManifest() templatemodule.Manifest {
	if m == nil {
		return permissions.Manifest()
	}
	return m.Manifest
}

// RegisterRoutes mounts Hrms HTTP routes on the root router.
func (m *Module) RegisterRoutes(r chi.Router, protected ...func(http.Handler) http.Handler) {
	httpx.RegisterRoutes(r, m.handler, protected...)
}

func runMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	entries, err := fs.ReadDir(migrations.FS, ".")
	if err != nil {
		return err
	}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			names = append(names, entry.Name())
		}
	}
	sort.Strings(names)
	for _, name := range names {
		content, err := migrations.FS.ReadFile(name)
		if err != nil {
			return err
		}
		if _, err := pool.Exec(ctx, string(content)); err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
	}
	return nil
}

func ensureTenantBrandingSchema(ctx context.Context, pool *pgxpool.Pool) error {
	content, err := migrations.FS.ReadFile("0005_upgrade_tenant_brandings.sql")
	if err != nil {
		return err
	}
	if _, err := pool.Exec(ctx, string(content)); err != nil {
		return fmt.Errorf("0005_upgrade_tenant_brandings.sql: %w", err)
	}
	return nil
}
