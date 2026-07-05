# AGENTS.md - spur-hrms

## Module Scope

`spur-hrms` is a Spur product module, not a standalone monolith. Use `BUSINESS_LOGIC_REFERENCE.md` as the source of HRMS business rules, but adapt all implementation to this module repository.

Do not create a monolithic `cmd/server`, `internal/app`, or nested `modules/*` application layout from the reference document. Keep this repository as one module with feature-focused files under the existing hexagonal structure.

For HRMS-077 onward, follow `MODERN_HRMS_IMPLEMENTATION_GUARDRAILS.md` before `MODERN_HRMS_REFERENCE.md`. The guardrails document maps the modern reference into this repo's migration numbering, permission naming, sqlc ownership, route conventions, and file boundaries.

## Task Evaluation and Clarification Gate

Before performing any HRMS task, evaluate the request for:

- Business goal and user role.
- Affected module areas and existing reusable workflows.
- Backend data model, migration, sqlc, permissions, services, handlers, and integration impact.
- Frontend navigation, modal forms, helper text, mobile behavior, and role-specific workflow impact.
- Tenant isolation, privacy/confidentiality, payroll/compliance risk, and data migration risk.
- Verification plan and whether live/dev server or database access is required.

If any requirement is ambiguous, acceptance criteria are missing, roles are unclear, data migration behavior is uncertain, UI flow conflicts with these rules, or there is doubt about production/dev impact, ask concise clarification questions first. Do not execute implementation until those questions are answered. Once answered, restate the resolved scope briefly, update or reference the relevant `TASKLIST.md` entry, and then execute.

## Product Research Requirement

Before implementing any HRMS task, research what the market expects and how comparable HRMS products handle the same workflow. Look for practical patterns around HR operations, employee self-service, payroll, attendance, compliance, reporting, mobile usage, and India-specific statutory needs when relevant.

Every completed task must include a short research outcome in the implementation summary, task update, or `TASKLIST.md` entry. Capture what was learned, what Setika will adopt, and what was intentionally avoided.

Use research for product judgment, not copy-paste parity. The final implementation must still follow Setika's tenant isolation, permission model, module architecture, and user experience rules.

## Required Architecture

Follow hexagonal architecture for every feature:

- Domain models and invariants live in `core/domain`.
- Core interfaces/ports live in `core/ports`.
- Business logic lives in `core/services`.
- PostgreSQL repositories live in `adapters/postgres`.
- Generated sqlc code lives only in `adapters/postgres/sqlc` and must not be hand-edited.
- SQL statements live only in `sql/queries/*.sql`.
- Schema changes live only in `sql/migrations/*.sql`.
- HTTP handlers live in `adapters/httpx/handlers`.
- HTTP route registration lives in `adapters/httpx/routes` or `adapters/httpx/routes.go` if the module keeps a single route file.
- Module construction and dependency wiring live in `hrms-module.go`.

Keep `core/ports` split by feature or aggregate. Do not add repository interfaces, service interfaces, commands, or DTOs to a single large `ports.go` catch-all file. Use focused files such as `employee_ports.go`, `leave_ports.go`, `attendance_ports.go`, `salary_ports.go`, and keep cross-cutting service contracts in a small dedicated file such as `service_ports.go`.

Repositories must call generated `sqlc.Queries` methods. Do not embed ad hoc SQL strings in repositories, services, handlers, or module wiring when sqlc can own the query. If a query is needed, add it to `sql/queries`, run `sqlc generate`, then map generated rows in `adapters/postgres`.

Keep `adapters/postgres/store.go` limited to adapter construction, shared dependencies, transaction helpers, and shared logging helpers. Put repository methods in focused files such as `employee_repository.go`, `leave_repository.go`, `attendance_repository.go`, and `salary_repository.go`.

Keep `core/services/service.go` limited to service construction, shared dependencies, and cross-cutting helpers such as `RunAsSystem`. Put business methods in focused service files named after the feature or aggregate, such as `branch_service.go`, `holiday_service.go`, `leave_service.go`, `attendance_service.go`, and `salary_service.go`. Do not add new feature CRUD/business logic to a large catch-all service file.

Keep sqlc-to-domain row mapping in dedicated mapper files in `adapters/postgres`. Do not mix database row mapping into service or handler files.

## Database Rules

All HRMS-owned tables must be in migrations under `sql/migrations`. Use the module schema consistently. Prefer `hrms.<table>` for HRMS-owned tables unless an existing migration intentionally defines otherwise.

Every tenant-scoped table must include:

- `tenant_id UUID NOT NULL`.
- `inactive BOOLEAN NOT NULL DEFAULT FALSE` for soft delete where the entity is user/business visible.
- `created_at` and `updated_at` where mutations are supported.
- RLS enabled when tenant data is isolated by PostgreSQL.
- Policies compatible with `app.tenant_id` and `app.is_super_admin` context.

Never recreate or directly mutate identity-owned auth schema tables. Interact with users, roles, permissions, and password migration through spur-identity ports/services. HRMS may reference identity IDs but must not bypass identity service boundaries.

If a flow touches RLS-protected tables and must run as system, it must go through a service-level system transaction. Do not bypass service boundaries from HTTP handlers.

## Logging Rules

Every adapter, repository, service, handler, scheduled job, and integration must have access to a logger instance through construction or dependency injection.

Log returned errors at the first boundary where operational context is known:

- Repositories log database failures with structured fields such as operation, tenant_id, user_id, entity_id, and query purpose.
- Services log validation and business-rule failures with operation and relevant domain identifiers.
- Handlers log request-level failures with route, method, tenant_id when available, and status.
- Integrations log external provider failures with provider, operation, status code when available, and correlation identifiers.

Warnings must be logged for recoverable but important conditions, including missing notification masters, missing reporting managers, skipped optional integrations, default working-hour fallback, and disabled provider configuration.

Logs must be visible on the console in development. If the host application does not pass a logger, create a safe module logger that writes to stdout/stderr in development rather than silently discarding errors.

## Business Guardrails

Tenant isolation is mandatory. Never trust tenant_id from request bodies for protected endpoints. Use authenticated tenant context from middleware/JWT.

Never hard-delete business records. Use `inactive = TRUE` except for replace-style child snapshots where the reference explicitly requires delete-and-reinsert inside a transaction, such as salary slip items.

Never return inactive records unless an explicit admin/audit API asks for them.

Leave rules:

- Apply leave increments `pending_days`; it must not deduct `used_days`.
- Approve leave moves days from `pending_days` to `used_days` and writes leave ledger rows.
- Reject leave decrements `pending_days` and does not touch `used_days`.
- Leave approval must run in one transaction.
- Leave overlap must handle full day, first half, and second half with a day-part bitmask.
- Leave dates must be within the same active financial year.
- Earn leave auto-conversion applies only when the leave type shortcode is `earnleave`.

Salary rules:

- Only one salary template can be active per tenant and financial year.
- Salary slip generation must reject existing slips unless an explicit regenerate flag is provided.
- Salary slips must snapshot salary items and leave balances at generation time.
- Add the LWP deduction line only when absent deduction is greater than zero.

Attendance rules:

- Attendance status priority is: before joining skip, approved leave, explicit attendance status, check-in present, holiday, weekoff, past absent, today empty.
- Working-hour fallback order is user, branch, tenant, default Mon-Fri 09:00-18:00 with Sat-Sun off.

Notifications and integrations:

- Email, SMS, WhatsApp, storage, PDF, and push delivery must be behind ports/interfaces.
- Do not hard-code SendGrid, MSG91, Gupshup, AWS S3, Firebase, or SMTP directly into business services.
- Missing notification master is a logged warning and a no-op, not a crash.

## Permissions and Manifest

Every HRMS feature must declare permissions in this module's manifest in the same change as the feature. Permission keys are module-local, for example `employees.list`, `leave.approve`, and `salary.generate`. Do not prefix keys with `hrms.` inside the key.

The host app must register the HRMS manifest with identity before bootstrap role assignment. Manifest registration must be idempotent and must upsert new permissions without duplicating existing ones.

Frontend role and permission screens must read permissions from the identity API. Do not add static HRMS permission lists to the frontend to hide missing manifest declarations.

## Frontend Completion Rule

A backend feature task is not `Complete` until its corresponding UI is implemented in `/Users/dinesh/workplace/setika-new/frontend` when the feature is user-facing.

Use `/Users/dinesh/workplace/setika-new/tailwind/template/src` as the visual reference for HRMS screens. Relevant references include:

- Employees: `employees.html`, `employees-grid.html`, `employee-details.html`, `employee-report.html`.
- Leave: `leaves.html`, `leaves-employee.html`, `leave-type.html`, `leave-settings.html`, `leave-report.html`.
- Attendance: `attendance-admin.html`, `attendance-employee.html`, `attendance-report.html`.
- Salary: `payroll.html`, `employee-salary.html`, `payslip.html`, `payslip-report.html`, `payroll-deduction.html`.
- Dashboard: `dashboard.html`, `employee-dashboard.html`.
- Tenant setup: `departments.html`, `designations.html`, `holidays.html`, `policy.html`, `bussiness-settings.html`.
- Onboarding: `candidates.html`, `candidates-grid.html`, `candidates-kanban.html`, `job-list.html`, `job-grid.html`.

Frontend auth token handling must stay centralized in shared auth/API helpers. Do not scatter token storage or direct Authorization header logic across pages.

## HRMS UI Simplicity Rules

HRMS screens must be clean, focused, and usable by non-technical HR users. Do not crowd a page with every field, rule, and explanation.

- All create, edit, setup, approval-comment, and configuration forms must open in modal popups unless the flow genuinely needs a dedicated multi-step wizard.
- After login, do not place marketing copy, obvious role explanations, or static instructional text directly in the main screen body. Authenticated screens must use working space for data, actions, charts, records, or controls. Put necessary explanation behind an information button, tooltip, popover, or contextual help drawer.
- Helper text, business-rule explanations, eligibility notes, and long guidance must be hidden behind an information button, tooltip, popover, or contextual help drawer. Do not place long helper paragraphs directly in page bodies or form rows.
- Main pages should emphasize summary cards, filters, tables/lists, status chips, and primary actions. Advanced settings should stay collapsed or move into a modal.
- Split complex workflows into tabs, segmented views, accordions, or task-specific modals instead of adding everything to one screen.
- For mobile or small screens, prefer compact action queues and bottom sheets/modals over desktop-sized forms and wide tables.
- If a UI starts to look cluttered, simplify the visible surface before adding more copy or controls.

## Task Tracking

Use `TASKLIST.md` as the implementation tracker. Every task must have one status: `Pending`, `In Progress`, `Blocked`, or `Complete`.

A task can be marked `Complete` only when:

- SQL migrations and sqlc queries are in the required folders.
- sqlc code has been regenerated when queries or schema changed.
- Repositories, ports, services, handlers, and routes follow the architecture above.
- Errors and warnings are logged at the correct boundary.
- Manifest permissions are declared when applicable.
- The matching frontend UI exists for user-facing work and references the Tailwind template where applicable.
- Tests or build checks have been run, or the task explicitly records why they could not be run.

Keep tasks small enough that a single feature slice can move from `Pending` to `Complete` without mixing unrelated domains.

## Engineering Standards

Apply DRY and single responsibility principles. Prefer small cohesive ports, services, repositories, handlers, and mappers over large catch-all files.

When a feature grows beyond a few related methods, split it by aggregate or workflow before adding more logic. New services, repositories, handlers, and mappers should be named for the business capability they implement, not for generic layers.

Validate at the domain/service boundary before repository calls. Keep handlers focused on transport concerns: request parsing, auth context extraction, service calls, response mapping, and request-level logging.

Do not hand-edit generated files. Do not hide failures behind generic `Failed to fetch` responses without backend logs that explain the real cause.
