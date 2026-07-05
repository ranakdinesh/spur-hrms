# spur-hrms — Hrms Module

Hrms support for Spur projects

---

## Install

```bash
spur add module hrms
```

Or manually:
```bash
go get github.com/ranakdinesh/spur-hrms@latest
```

---

## Wire into app.go

```go
import hrms "github.com/ranakdinesh/spur-hrms"

hrmsModule, err := hrms.New(ctx, hrms.Options{
    DB:              dbPool,
    Log:             log,
    Cfg:             cfg.Hrms,
    MigrationRunner: infra.Migrations.Run,
})
if err != nil {
    return nil, fmt.Errorf("hrms: %w", err)
}
if err := identityModule.Services.ModuleService.RegisterManifest(ctx, hrmsModule.Manifest); err != nil {
    return nil, fmt.Errorf("hrms manifest: %w", err)
}
hrmsModule.RegisterRoutes(r)
```

---

## Configuration

```bash
# deployments/.env
# TODO: add your env vars here
# HRMS_API_KEY=
```

---

## HTTP Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET    | /hrms       | List all (tenant-scoped) |
| POST   | /hrms       | Create new |
| GET    | /hrms/{id}  | Get by ID |
| DELETE | /hrms/{id}  | Delete |

---

## Using in Another Module

```go
// Inject via Options
type Options struct {
    HrmsSvc hrms.ports.HrmsService
}
```

---

## Development

```bash
# Generate sqlc code after editing sql/queries/
sqlc generate

# Run tests
go test ./...

# Build check
go build ./...
```

---

## Adding to the Registry

After pushing to GitHub, add an entry to
`github.com/ranakdinesh/spur-registry/modules.json`.
Copy the contents of `spur.json` in this repo.
