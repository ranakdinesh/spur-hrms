# HRMS Navigation And Command Center UX Audit

Date: 2026-06-27
Target: `http://localhost:3000/dashboard`
Evidence:
- `01-user-reported-command-center.png` - live screenshot supplied by user.
- Source review: `frontend/src/app/dashboard/page.tsx`, `HRCommandCenterSection.tsx`, `OperationsWorkbenchSection.tsx`.

## Scope

This audit reviews the new Setika HRMS dashboard shell, side menu, command center, and operations work entry points for HR/Tenant Admin users. The goal is to explain why the new design feels confusing and what structure should replace it.

Browser automation note: the in-app browser runtime and local Playwright capture were unavailable in this session. The audit is therefore based on the user-provided live screenshot plus source inspection of the rendered navigation and command-center components.

## User Goal

An HR user should log in and immediately know:

1. What needs action today.
2. Where routine employee requests and approvals live.
3. Where to manage people records.
4. Where to configure company settings.
5. How to start a new request or operation.

The current UI does not support that mental model reliably.

## Findings

### 1. The Left Menu Is Organized By Implementation, Not User Work

Evidence:
- `My Work` contains Attendance, Shift Scheduling, Benefits, Projects, Work Logs, Learning, Talent Marketplace, OKRs, Performance, Wellbeing, Leaves, Payslips, Agreements, Policies, HR Helpdesk, Employee Relations, Inbox, and My Documents.
- `My Team` contains Operations Workbench, Workflow Inbox, HR Helpdesk, Employee Relations, Employees, Leave Approvals, Leave Reports, Projects, Skill Gaps, Learning, Talent Marketplace, OKRs, Performance, Wellbeing, Work Logs, Agreements, Reports, Insights, People Analytics, Attendance, Shift Scheduling, and Benefits.
- `People` repeats Operations Workbench, Workflow Inbox, Employee Relations, Employees, Projects, Work Logs, Compliance, Skills, Learning, Analytics, Agreements, Document Sign, Asset & Access, Exits, and Onboarding.

Impact:
- Users see the same destination in multiple places.
- "My Team" and "People" overlap heavily, so neither name explains what belongs there.
- The user cannot predict where a module is before searching visually.

Recommendation:
Replace the current top-level menu with role-task groups:

- `Inbox`
  - My Tasks
  - Approvals
  - Requests
  - Notifications
- `People`
  - Employees
  - Workforce Hub
  - Documents
  - Assets & Access
  - Exits
- `Time`
  - Attendance
  - Shifts
  - Leave
  - Holidays
- `Payroll`
  - Payroll Runs
  - Salary
  - Payslips
  - Compensation
  - Benefits
- `Hiring`
  - Requisitions
  - Jobs
  - Candidates
  - Interviews
  - Offers
  - Onboarding
- `Performance`
  - OKRs
  - Reviews
  - Skills
  - Learning
  - Succession
- `Reports`
  - Reports
  - People Analytics
  - Insights
- `Setup`
  - Company
  - Policies
  - Leave Rules
  - Payroll Rules
  - Notifications
  - Integrations
  - Roles

### 2. The Dashboard Has Two Competing Task Concepts

Evidence:
- The default screen is `HR Command Center`.
- The menu also exposes `Operations Workbench`.
- The menu also exposes `Workflow Inbox`.
- The command center itself has "Priority queue", "Quick actions", "In-context AI", "New", and "Refresh".

Impact:
- A user cannot tell when to use Command Center versus Operations Workbench versus Workflow Inbox.
- The product is trying to implement an inbox model, but it still exposes three different task surfaces.

Recommendation:
Make the operating model explicit:

- Default landing page should be `Inbox`.
- `Inbox` should have Outlook-style folders:
  - `Assigned to me`
  - `Approvals`
  - `Watching`
  - `Created by me`
  - `Completed`
  - `Snoozed`
  - `Delegated`
- Command Center should become an HR leadership overview, not the default task UI.
- Operations Workbench should be hidden from everyday users or renamed `All Work` under Inbox for HR Ops only.

### 3. The Current "New" Button Is Too Abstract

Evidence:
- Header has a global `New` button.
- Command Center also has a `New` button.
- The user's question earlier was exactly whether this should behave like a new email button.

Impact:
- "New" does not tell HR users what they are creating.
- It competes with workflow-specific create buttons.

Recommendation:
Keep one global create button, but label and structure it as:

`+ New Request`

Inside the modal:
- Employee Request
- HR Case
- Leave / Attendance Request
- Document Sign Request
- Asset / Access Request
- Hiring Request
- Payroll Change
- Company Setup Change
- Tenant Operation

Each option should show only:
- Name
- One-line purpose
- Required approval path
- Owner team

Do not show technical metadata such as source module, workflow template keys, or permission names in the primary UI.

### 4. Active Menu Expansion Creates Hidden Navigation

Evidence:
- In `dashboard/page.tsx`, a parent menu expands only when it is active.
- Clicking a parent with children navigates to the first visible child instead of simply expanding/collapsing.

Impact:
- Users lose control: clicking a menu category changes the page.
- Only one active branch is visible, so users must hunt through top-level categories.
- This is especially painful in a long menu with repeated destinations.

Recommendation:
- Parent click should expand/collapse only.
- Child click should navigate.
- Allow 1-2 branches to stay open, or use a pinned "Favorites" group.
- Add a clear selected state for the actual child, not only the parent.

### 5. Search Looks Useful But Is Not Functional Enough

Evidence:
- Header search input says `Search in HRMS`.
- Source shows the input has no visible result panel or command behavior in the inspected area.

Impact:
- In a large HRMS, search becomes the safety net for poor navigation. If it does not search modules, employees, tasks, and actions, users feel trapped.

Recommendation:
Implement command search as a first-class navigation tool:

- Search modules: "attendance", "employees", "payroll".
- Search actions: "approve leave", "create employee", "mark attendance".
- Search people: employee name/code.
- Search open work: task title, request number.
- Keyboard shortcut should open a modal palette, not a passive input.

### 6. Visual Hierarchy Is Too Heavy For A Work App

Evidence from screenshot:
- Large sidebar cards consume major width.
- Main content has very large cards even when all metrics are zero.
- There is too much empty whitespace in the default state.
- The user profile and header action area are visually heavy.

Impact:
- The screen feels oversized and slow to scan.
- HR users doing repeated daily work need density and predictability, not a marketing-dashboard feel.

Recommendation:
- Make the shell denser.
- Reduce sidebar width from 292px to about 248px.
- Use smaller menu row height.
- Use one-line menu rows with icons and labels; move item counts to subtle badges only where useful.
- Replace metric cards with a compact status strip when counts are zero.
- Show "No priority work" with suggested next actions, not a large empty panel.

### 7. Accessibility Risks

Likely risks from screenshot and source:
- Many icon-only header buttons rely on title/aria labels, but visual meaning is weak.
- The `i` info buttons are text-only circles and may be unclear to screen reader users without stronger labels.
- Parent menu buttons both navigate and expand, which is poor keyboard behavior.
- Long side navigation creates high tab-stop volume.
- Low-contrast gray-green labels may be difficult for some users.

Recommendations:
- Parent menu buttons should have `aria-controls` and only expand/collapse.
- Add visible labels or tooltips for all icon-only header actions.
- Provide skip links to main content.
- Keep focus visible inside side nav and modals.
- Test with keyboard only at 100%, 150%, and 200% zoom.

## Recommended Target IA

For HR/Tenant Admin:

1. `Inbox`
   Daily actionable work: approvals, requests, delegated tasks, completed work.
2. `People`
   Employee records and workforce lifecycle.
3. `Time`
   Attendance, shifts, leave, holidays.
4. `Payroll`
   Payroll runs, salary, compensation, benefits.
5. `Hiring`
   Hiring pipeline and onboarding.
6. `Performance`
   OKRs, reviews, skills, learning, succession.
7. `Reports`
   Reporting, analytics, AI insights.
8. `Setup`
   Company configuration, policies, roles, integrations.

For Employees:

1. `Home`
2. `My Inbox`
3. `Attendance`
4. `Leave`
5. `Payslips`
6. `Documents`
7. `Learning`
8. `Helpdesk`

For Super Admin:

1. `Platform`
2. `Tenants`
3. `Tenant Operations`
4. `Billing`
5. `Users & Roles`
6. `System Settings`

Do not show HR operational modules inside Super Admin unless a tenant has been selected.

## Implementation Plan To Fix

### P0 - Make Navigation Usable

- Replace overlapping `My Work`, `My Team`, and `People` group logic with role-specific IA.
- Remove duplicate destinations from multiple groups unless there is a clear task reason.
- Parent menu click expands only; child click navigates.
- Add a `Favorites` or `Pinned` group for most-used items.

### P0 - Consolidate Task Surfaces

- Make `Workflow Inbox` the default task interface.
- Move Command Center to `Reports > Command Center` or `Inbox > Overview`.
- Move Operations Workbench to `Inbox > All Work` for HR Ops/Admin only.

### P1 - Fix Global New

- Rename `New` to `+ New Request`.
- Use a clean modal request catalog grouped by user intent.
- Hide technical metadata by default.
- Route creation to owning workflows.

### P1 - Improve Search

- Turn header search into a command palette.
- Search modules, actions, employees, and work items.
- Show keyboard-first results with section labels.

### P1 - Reduce Visual Noise

- Reduce sidebar width and menu row height.
- Make zero-state Command Center compact.
- Use fewer oversized cards.
- Keep helper explanations behind info buttons, but label those buttons clearly.

## Step Review

1. Login/default landing: partially healthy. The app opens to the correct product shell, but the default HR screen is not self-explanatory enough for daily work.
2. Command Center: weak. It summarizes work but competes with Inbox and Workbench; empty states waste space.
3. Side menu: unhealthy. Too many repeated modules, unclear group names, and active-only expansion make it hard to navigate.
4. Workflow Inbox concept: promising. This should become the main operating surface.
5. Global New: promising but unclear. It needs request-oriented language and a simpler catalog.
6. Mobile behavior: not fully verified. Source suggests the same large navigation model collapses behind a hamburger, which likely makes the IA problem worse on mobile.

## Bottom Line

The redesign failed because it added a modern-looking shell without simplifying the product model. The right direction is not more cards or more menu groups. Setika needs one clear HR operating model:

`Inbox for work, modules for records, setup for configuration, reports for insight.`

Once that model is applied, the same existing modules can become much more usable without removing functionality.
