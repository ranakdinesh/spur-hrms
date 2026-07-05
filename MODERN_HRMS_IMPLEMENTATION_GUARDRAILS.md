# Modern HRMS Implementation Guardrails

This document is the repo-specific plan for HRMS-077 onward. It adapts
`../MODERN_HRMS_REFERENCE.md` to the current `spur-hrms` module architecture.

`../MODERN_HRMS_TRANSFORMATION.md` is referenced by `TASKLIST.md`, but it is not
present in this workspace as of HRMS-076. Until it exists, use this document,
`../MODERN_HRMS_REFERENCE.md`, `BUSINESS_LOGIC_REFERENCE.md`, and `AGENTS.md` as
the implementation sources of truth.

## Research Outcome

Modern HCM products and market guidance point toward one consolidated workforce
system instead of disconnected point tools: employee and contingent workforce
coverage, skills-based planning, internal mobility, flexible work/payment models,
privacy-aware analytics, and AI/automation with human approval and audit trails.
For Setika, that means modern features must extend the existing tenant-scoped
HRMS module and identity model rather than creating separate apps, separate
worker identities, or AI-managed source-of-truth mutations.

Useful market and compliance references reviewed:

- Workday/SAP-style HCM direction: workforce skills, talent intelligence, and
  external workforce management are expected enterprise capabilities.
- HR automation trend: consolidation and single workforce data are preferred over
  fragmented tools and manual re-entry.
- India-specific needs: contractor/consultant classifications must consider
  CLRA, PF/ESIC/PT applicability, TDS 194C/194J handling, and DPDPA consent,
  notice, retention, and erasure expectations.

## Non-Negotiable Scope Rules

- Keep modern HRMS inside this `spur-hrms` module. Do not create nested
  `modules/workforce`, `modules/projects`, `modules/skills`, or a new monolith.
- Keep the current hexagonal layout:
  - Domain: `core/domain`
  - Ports: `core/ports`
  - Services: `core/services`
  - Repositories and mappers: `adapters/postgres`
  - Generated sqlc: `adapters/postgres/sqlc`
  - SQL queries: `sql/queries`
  - Migrations: `sql/migrations`
  - HTTP handlers: `adapters/httpx/handlers`
  - Route registration: `adapters/httpx/routes.go`
  - Module wiring and options: `hrms-module.go`
- Do not redefine identity-owned users, roles, permissions, or passwords. HRMS
  can reference identity user IDs through ports and tenant context only.
- Do not turn non-employee workers into ordinary employees by default. HRMS-078
  must define a worker profile extension that can link to an existing employee
  record when the worker is also an employee.
- Do not let AI or automation approve payment, terminate engagements, change pay,
  or mutate source-of-truth records without a human action.
- Do not store raw PII, salary values, bank data, statutory IDs, exact GPS
  coordinates, or free-text sensitive remarks in AI signal/action tables.

## Migration Numbering Plan

The current latest migration is `0054_insights_foundation.sql`. Modern HRMS
migrations start at `0055`. Never reuse or renumber an existing migration.

| Task | Primary Migration |
|---|---|
| HRMS-077 | `0055_worker_type_taxonomy.sql` |
| HRMS-078 | `0056_worker_profile_extensions.sql` |
| HRMS-079 | `0057_engagements.sql` |
| HRMS-080 | `0058_work_logs.sql` |
| HRMS-081 | `0059_projects_milestones.sql` |
| HRMS-082 | `0060_agreements_lifecycle.sql` |
| HRMS-083 | `0061_esign_provider_settings.sql` |
| HRMS-084 | `0062_flexible_worker_payroll.sql` |
| HRMS-085 | `0063_compliance_rules_checklists.sql` |
| HRMS-086 | `0064_skills_catalog_worker_skills.sql` |
| HRMS-087 | `0065_project_skill_requirements.sql` |
| HRMS-088 | `0066_talent_marketplace.sql` |
| HRMS-089 | `0067_okr_foundation.sql` |
| HRMS-090 | `0068_performance_checkins_feedback.sql` |
| HRMS-091 | `0069_pulse_wellbeing.sql` |
| HRMS-092 | `0070_ai_action_layer.sql` |
| HRMS-093 | `0071_bounded_ai_agent_state.sql` |
| HRMS-094 | `0072_modern_people_analytics.sql` |
| HRMS-095 | `0073_ecosystem_privacy_hardening.sql` |

If a task needs more than one migration, keep the base number for the main table
set and use the next available number for follow-up indexes or seed changes. Do
not combine unrelated task schemas into one migration just to save numbers.

## Database Rules For Modern Tables

- Use the `hrms` schema for all HRMS-owned modern tables.
- Prefer table names that match existing repo style:
  - `hrms.worker_types`
  - `hrms.worker_profiles`
  - `hrms.engagements`
  - `hrms.work_logs`
  - `hrms.projects`
  - `hrms.project_milestones`
  - `hrms.agreement_templates`
  - `hrms.agreements`
  - `hrms.skills`
  - `hrms.worker_skills`
  - `hrms.project_skill_requirements`
  - `hrms.marketplace_opportunities`
  - `hrms.okr_cycles`
  - `hrms.objectives`
  - `hrms.key_results`
  - `hrms.pulse_surveys`
  - `hrms.pulse_responses`
  - `hrms.modern_ai_signals`
- Every tenant-scoped table must include `tenant_id`, `inactive`, `created_at`,
  and `updated_at` unless it is an immutable append-only audit/event table.
- Enable RLS on tenant tables and add policies compatible with `app.tenant_id`
  and `app.is_super_admin`.
- Use soft delete for business records. Do not hard-delete except replace-style
  child snapshots inside an explicit transaction.
- Use database constraints for invariant protection:
  - Unique active worker type code per tenant.
  - No duplicate active worker skill per tenant/worker/skill.
  - No duplicate work log for the same tenant/worker/engagement/date.
  - No duplicate pay run item for the same pay run and worker.
- Add indexes during the owning task, not as a later cleanup, when filters or
  joins are part of the user workflow.
- Use PostgreSQL enum types only when the allowed values are stable and owned by
  the module. Use lookup tables when tenants must customize values.

## SQLC Ownership

Every new aggregate must own a query file in `sql/queries`:

| Aggregate | Query File |
|---|---|
| Worker types and worker profiles | `modern_workforce.sql` |
| Engagements and work logs | `engagements.sql` |
| Projects and milestones | `projects.sql` |
| Agreements | `agreements.sql` |
| E-sign provider state | `esign.sql` |
| Flexible worker payroll | `flex_payroll.sql` |
| Compliance rules/checklists | `compliance.sql` |
| Skills and worker skills | `skills.sql` |
| Marketplace | `talent_marketplace.sql` |
| OKR and performance | `performance.sql` |
| Pulse and wellbeing | `wellbeing.sql` |
| AI action layer | `ai_actions.sql` |
| Modern analytics snapshots | `modern_analytics.sql` |
| Privacy/integration hardening | `privacy_integrations.sql` |

Repositories must call generated sqlc methods only. If a query is missing, add
or update the query file, run `sqlc generate`, then map rows in
`adapters/postgres/*_mappers.go`.

## Permission Naming

Permission keys remain module-local. Do not prefix keys with `hrms.` in
`pkg/permissions`.

Use consistent action suffixes:

- Read: `.list`, `.view`
- Mutation: `.create`, `.update`, `.delete`
- Workflow action: `.submit`, `.approve`, `.reject`, `.cancel`, `.complete`,
  `.terminate`, `.publish`, `.send`, `.sign`, `.export`
- Configuration: `.manage`

Suggested modern permission groups:

- `worker_types.*`
- `workers.*`
- `engagements.*`
- `work_logs.*`
- `projects.*`
- `milestones.*`
- `agreements.*`
- `esign_settings.manage`
- `flex_pay_runs.*`
- `contractor_invoices.*`
- `compliance_rules.manage`
- `compliance_checklists.*`
- `skills.*`
- `worker_skills.*`
- `skill_gaps.view`
- `marketplace_opportunities.*`
- `okr.*`
- `performance_checkins.*`
- `feedback.*`
- `pulse_surveys.*`
- `wellbeing.view`
- `ai_actions.view`
- `ai_actions.review`
- `modern_analytics.view`
- `privacy_requests.*`
- `integrations.manage`

Every backend feature task must add permissions and role-template grants in
`pkg/permissions` in the same change as handlers and routes.

## Backend Ownership Boundaries

Use focused files. Do not add modern feature code to catch-all files.

| Task | Domain | Ports | Service | Repository | Handler |
|---|---|---|---|---|---|
| HRMS-077 | `worker_type.go` | `worker_type_ports.go` | `worker_type_service.go` | `worker_type_repository.go` | `worker_type_handler.go` |
| HRMS-078 | `worker_profile.go` | `worker_profile_ports.go` | `worker_profile_service.go` | `worker_profile_repository.go` | `worker_profile_handler.go` |
| HRMS-079 | `engagement.go` | `engagement_ports.go` | `engagement_service.go` | `engagement_repository.go` | `engagement_handler.go` |
| HRMS-080 | `work_log.go` | `work_log_ports.go` | `work_log_service.go` | `work_log_repository.go` | `work_log_handler.go` |
| HRMS-081 | `project.go`, `milestone.go` | `project_ports.go` | `project_service.go`, `milestone_service.go` | `project_repository.go` | `project_handler.go` |
| HRMS-082 | `agreement.go` | `agreement_ports.go` | `agreement_service.go` | `agreement_repository.go` | `agreement_handler.go` |
| HRMS-083 | `esign.go` | `esign_ports.go` | `esign_service.go` | `esign_repository.go` | `esign_handler.go` |
| HRMS-084 | `flex_payroll.go` | `flex_payroll_ports.go` | `flex_payroll_service.go` | `flex_payroll_repository.go` | `flex_payroll_handler.go` |
| HRMS-085 | `compliance.go` | `compliance_ports.go` | `compliance_service.go` | `compliance_repository.go` | `compliance_handler.go` |
| HRMS-086 | `skill.go` | `skill_ports.go` | `skill_service.go` | `skill_repository.go` | `skill_handler.go` |
| HRMS-087 | `skill_gap.go` | `skill_gap_ports.go` | `skill_gap_service.go` | use `skill_repository.go` | `skill_gap_handler.go` |
| HRMS-088 | `talent_marketplace.go` | `talent_marketplace_ports.go` | `talent_marketplace_service.go` | `talent_marketplace_repository.go` | `talent_marketplace_handler.go` |
| HRMS-089 | `okr.go` | `okr_ports.go` | `okr_service.go` | `okr_repository.go` | `okr_handler.go` |
| HRMS-090 | `performance.go` | `performance_ports.go` | `performance_service.go` | `performance_repository.go` | `performance_handler.go` |
| HRMS-091 | `wellbeing.go` | `wellbeing_ports.go` | `wellbeing_service.go` | `wellbeing_repository.go` | `wellbeing_handler.go` |
| HRMS-092/093 | `ai_action.go` | `ai_action_ports.go` | `ai_action_service.go` | `ai_action_repository.go` | `ai_action_handler.go` |
| HRMS-094 | `modern_analytics.go` | `modern_analytics_ports.go` | `modern_analytics_service.go` | `modern_analytics_repository.go` | `modern_analytics_handler.go` |
| HRMS-095 | `privacy.go` | `privacy_ports.go` | `privacy_service.go` | `privacy_repository.go` | `privacy_handler.go` |

`store.go`, `service.go`, and `routes.go` can be touched for construction and
route registration only. They must not accumulate feature logic.

## Route Conventions

The reference document uses `/api/v1`. In this repo, use existing HRMS route
prefixes:

- Tenant routes: `/hrms/<resource>`
- Super-admin tenant routes: `/hrms/tenants/{tenantID}/<resource>`
- Public webhook routes only when required, for example
  `/hrms/webhooks/esign/{provider}` with HMAC verification and replay checks.

Every protected endpoint must derive tenant ID from middleware/JWT context or
super-admin route parameter. Never trust tenant ID from JSON bodies.

## Frontend Rules

- Do not add static permission lists to the frontend.
- Add navigation placeholders only when the backend manifest permission exists
  and the page can read permissions from identity.
- Keep setup/create/edit workflows in modals. Use dense lists, filters, cards,
  status chips, segmented controls, and tabs for complex modern modules.
- Heavy setup remains web-first. Mobile gets approval queues, alerts, people
  lookup, and lightweight self-service actions only.
- Reference Tailwind screens from `/Users/dinesh/workplace/setika-new/tailwind`
  for visual density, but do not copy long helper text into page bodies.

## Integration And AI Rules

- Redis/event bus, Python AI sidecar, e-sign, WhatsApp, Slack, email, storage,
  and push must be behind ports/interfaces.
- If Redis or the Python sidecar is unavailable, services must degrade to
  deterministic rules, logged warnings, and human-review queues.
- AI agents may write to AI/insight/action tables and notification/checklist
  queues. They must not directly mutate employee, engagement, agreement, salary,
  invoice, milestone, or payroll source-of-truth rows.
- Every AI recommendation must include reason codes, confidence, source object
  references, generated timestamp, and human review status.
- Every externally triggered webhook must verify signature, timestamp/replay,
  provider identity, tenant or reference mapping, and idempotency key.

## India Compliance And Privacy Rules

- Worker classification must drive statutory defaults, not free text:
  PF, ESIC, PT, LWF, CLRA, TDS section, attendance mode, and pay mode.
- TDS calculation for non-employees must be separate from salary TDS under
  section 192. HRMS-084 must handle 194C/194J paths explicitly.
- Contractor and agency workflows must capture evidence and waivers through
  compliance checklist records, not comments fields.
- DPDPA-related features must support notice, consent/audit references,
  retention policy, correction/erasure workflow state, and grievance/escalation
  metadata. Do not physically erase records needed for payroll, statutory, or
  audit retention without a policy decision.

## Testing And Completion Gates

Each modern feature slice must record:

- `sqlc generate` when migrations or query files change.
- `GOCACHE=/Users/dinesh/workplace/setika-new/spur-hrms/.gocache go test ./...`
  for backend changes.
- `npm run build` for frontend changes.
- Focused unit tests for state transitions and calculations:
  engagement status, work-log approval, milestone acceptance, TDS thresholds,
  compliance rule detection, skill matching order, OKR rollups, pulse privacy,
  and AI action review/override.
- `TASKLIST.md` research outcome explaining what market/compliance pattern was
  adopted and what was intentionally kept out of scope.
