// Package hrms provides the Hrms Spur module.
//
// Install:
//   spur add module hrms
//
// Or manually:
//   go get y@latest
//
// Wire in app.go:
//   hrmsModule, err := hrms.New(ctx, hrms.Options{DB: dbPool, Log: log, Cfg: cfg.Hrms})
//   hrmsModule.RegisterRoutes(r)
package hrms

import (
	"context"
	"fmt"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ranakdinesh/spur-platform/logger"
	"y/adapters/httpx"
	"y/adapters/httpx/handlers"
	"y/adapters/postgres"
	"y/core/ports"
	"y/core/services"
	"y/sql/migrations"
)

// Config holds all Hrms configuration from environment variables.
type Config struct {
	// TODO: add your config fields
	// Example:
	// APIKey string `env:"HRMS_API_KEY"`
}

// Options is passed by app.go when constructing this module.
type Options struct {
	DB  *pgxpool.Pool
	Log *logger.Loggerx
	Cfg Config

	// MigrationRunner runs this module's SQL migrations.
	// Provided by the platform: infra.Migrations.Run
	MigrationRunner func(ctx context.Context, moduleName string, fs interface{}) error
}

// Module is the Hrms module entry point.
type Module struct {
	// Services exposes this module's service interfaces to other modules.
	Services *Services

	handler *handlers.Handler
}

// Services bundles the public service interfaces.
type Services struct {
	Hrms ports.HrmsService
}

// New wires the Hrms module. Returns error — never panics.
func New(ctx context.Context, opt Options) (*Module, error) {
	if opt.DB == nil {
		return nil, fmt.Errorf("hrms: DB pool is required")
	}

	// Run migrations
	if opt.MigrationRunner != nil {
		if err := opt.MigrationRunner(ctx, "hrms", migrations.FS); err != nil {
			return nil, fmt.Errorf("hrms: migrations: %w", err)
		}
	}

	// Wire repo → service → handler
	store := postgres.New(opt.DB)
	svc   := services.NewHrmsService(store, opt.Log)
	h     := handlers.New(svc)

	opt.Log.Info(ctx).Str("module", "hrms").Msg("Hrms module initialised")

	return &Module{
		Services: &Services{ Hrms: svc }, 
		handler:  h,
	}, nil
}

// RegisterRoutes mounts Hrms HTTP routes on the root router.
func (m *Module) RegisterRoutes(r chi.Router) {
	httpx.RegisterRoutes(r, m.handler)
}
