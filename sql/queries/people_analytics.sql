-- name: GetPeopleAnalyticsWorkspace :one
WITH workforce AS (
    SELECT
        COUNT(*)::int AS total_workers,
        COUNT(*) FILTER (WHERE wp.profile_status = 'active')::int AS active_workers,
        COUNT(*) FILTER (WHERE wp.profile_status IN ('paused', 'ended'))::int AS inactive_or_paused,
        COALESCE(jsonb_agg(jsonb_build_object('label', COALESCE(wt.name, 'Unclassified'), 'count', grouped.total_count) ORDER BY grouped.total_count DESC), '[]'::jsonb) AS type_mix
    FROM hrms.worker_profiles wp
    LEFT JOIN hrms.worker_types wt ON wt.id = wp.worker_type_id
    LEFT JOIN LATERAL (
        SELECT COUNT(*)::int AS total_count
        FROM hrms.worker_profiles inner_wp
        WHERE inner_wp.tenant_id = wp.tenant_id
          AND inner_wp.worker_type_id = wp.worker_type_id
          AND NOT inner_wp.inactive
    ) grouped ON TRUE
    WHERE wp.tenant_id = $1 AND NOT wp.inactive
),
renewals AS (
    SELECT
        COUNT(*) FILTER (WHERE status = 'active' AND renewal_due_date BETWEEN CURRENT_DATE AND CURRENT_DATE + INTERVAL '30 days')::int AS due_30,
        COUNT(*) FILTER (WHERE status = 'active' AND renewal_due_date < CURRENT_DATE)::int AS overdue,
        COUNT(*) FILTER (WHERE status = 'active')::int AS active_engagements
    FROM hrms.engagements
    WHERE tenant_id = $1 AND NOT inactive
),
projects AS (
    SELECT
        COUNT(*) FILTER (WHERE status = 'active')::int AS active_projects,
        COUNT(*) FILTER (WHERE status = 'on_hold')::int AS on_hold_projects,
        COUNT(*) FILTER (WHERE status = 'active' AND due_date < CURRENT_DATE)::int AS overdue_projects,
        COUNT(*) FILTER (WHERE priority IN ('high', 'critical') AND status = 'active')::int AS high_priority_projects
    FROM hrms.projects
    WHERE tenant_id = $1 AND NOT inactive
),
skills AS (
    SELECT
        COUNT(*)::int AS skill_records,
        COUNT(*) FILTER (WHERE verification_status = 'hr_verified')::int AS verified_skills,
        COUNT(*) FILTER (WHERE verification_status IN ('self_declared', 'expired'))::int AS verification_backlog
    FROM hrms.worker_skills
    WHERE tenant_id = $1 AND NOT inactive
),
requirements AS (
    SELECT COUNT(*)::int AS open_skill_requirements
    FROM hrms.project_skill_requirements
    WHERE tenant_id = $1 AND NOT inactive
),
wellbeing AS (
    SELECT
        COUNT(*)::int AS score_count,
        ROUND(COALESCE(AVG(wellbeing_score), 0), 2)::float8 AS wellbeing_index,
        COUNT(*) FILTER (WHERE risk_level IN ('high', 'critical'))::int AS high_risk_count
    FROM hrms.wellbeing_scores
    WHERE tenant_id = $1 AND NOT inactive
      AND score_date >= CURRENT_DATE - INTERVAL '90 days'
),
costs AS (
    SELECT
        ROUND(COALESCE(SUM(net_salary), 0), 2)::float8 AS payroll_cost_12m,
        ROUND(COALESCE(AVG(net_salary), 0), 2)::float8 AS average_net_salary,
        COUNT(*)::int AS salary_slip_count
    FROM hrms.salary_slips
    WHERE tenant_id = $1 AND NOT inactive
      AND make_date(year, month, 1) >= date_trunc('month', CURRENT_DATE) - INTERVAL '12 months'
),
risks AS (
    SELECT jsonb_build_array(
        jsonb_build_object('area', 'Contract renewals', 'severity', CASE WHEN renewals.overdue > 0 THEN 'high' WHEN renewals.due_30 > 0 THEN 'medium' ELSE 'low' END, 'count', renewals.overdue + renewals.due_30),
        jsonb_build_object('area', 'Project delivery', 'severity', CASE WHEN projects.overdue_projects > 0 THEN 'high' WHEN projects.on_hold_projects > 0 THEN 'medium' ELSE 'low' END, 'count', projects.overdue_projects + projects.on_hold_projects),
        jsonb_build_object('area', 'Skills coverage', 'severity', CASE WHEN requirements.open_skill_requirements > skills.verified_skills THEN 'medium' ELSE 'low' END, 'count', GREATEST(requirements.open_skill_requirements - skills.verified_skills, 0)),
        jsonb_build_object('area', 'Wellbeing', 'severity', CASE WHEN wellbeing.high_risk_count > 0 THEN 'high' ELSE 'low' END, 'count', wellbeing.high_risk_count)
    ) AS heatmap
    FROM renewals, projects, skills, requirements, wellbeing
)
SELECT jsonb_build_object(
    'generated_at', NOW(),
    'workforce', jsonb_build_object(
        'total_workers', workforce.total_workers,
        'active_workers', workforce.active_workers,
        'inactive_or_paused', workforce.inactive_or_paused,
        'type_mix', workforce.type_mix
    ),
    'engagement_health', jsonb_build_object(
        'active_engagements', renewals.active_engagements,
        'renewals_due_30', renewals.due_30,
        'renewals_overdue', renewals.overdue
    ),
    'project_health', to_jsonb(projects),
    'skills_intelligence', jsonb_build_object(
        'skill_records', skills.skill_records,
        'verified_skills', skills.verified_skills,
        'verification_backlog', skills.verification_backlog,
        'open_skill_requirements', requirements.open_skill_requirements
    ),
    'wellbeing', to_jsonb(wellbeing),
    'cost_intelligence', to_jsonb(costs),
    'risk_heatmap', risks.heatmap
) AS workspace
FROM workforce, renewals, projects, skills, requirements, wellbeing, costs, risks;
