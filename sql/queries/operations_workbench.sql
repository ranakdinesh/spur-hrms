-- name: ListOperationsWorkbenchCards :many
WITH cards AS (
    SELECT
        la.tenant_id,
        ('leave_approval:' || la.id::text) AS card_key,
        'approvals'::text AS lane,
        'leave'::text AS category,
        'leave'::text AS source_module,
        'leave_approval'::text AS source_type,
        la.id AS source_id,
        l.user_id AS employee_user_id,
        ('Approve leave for ' || trim(e.firstname || ' ' || COALESCE(e.lastname, '')))::text AS title,
        (COALESCE(lt.name, 'Leave') || ' from ' || l.start_date::text || ' to ' || l.end_date::text || ' (' || COALESCE(l.days, 0)::text || ' days)')::text AS summary,
        la.status::text AS status,
        CASE
            WHEN l.start_date <= CURRENT_DATE THEN 'high'
            WHEN l.start_date <= CURRENT_DATE + INTERVAL '2 days' THEN 'medium'
            ELSE 'low'
        END::text AS severity,
        CASE
            WHEN l.start_date <= CURRENT_DATE THEN 20
            WHEN l.start_date <= CURRENT_DATE + INTERVAL '2 days' THEN 40
            ELSE 70
        END::int AS priority,
        l.start_date::timestamptz AS due_at,
        la.created_at AS detected_at,
        'Review leave'::text AS action_label,
        'leave-approvals'::text AS route_section,
        l.id AS route_record_id,
        jsonb_build_object(
            'leave_id', l.id,
            'approver_id', la.approver_id,
            'employee_name', trim(e.firstname || ' ' || COALESCE(e.lastname, '')),
            'employee_code', e.employee_code,
            'department', d.name,
            'leave_type', COALESCE(lt.name, 'Leave'),
            'start_date', l.start_date,
            'end_date', l.end_date,
            'start_day_type', l.start_day_type,
            'end_day_type', l.end_day_type,
            'days', l.days,
            'reason', l.reason,
            'applied_date', l.applied_date,
            'source_status', l.status
        ) AS metadata
    FROM hrms.leave_approvals la
    JOIN hrms.leaves l ON l.tenant_id = la.tenant_id AND l.id = la.leave_id AND NOT l.inactive
    JOIN hrms.employees e ON e.tenant_id = l.tenant_id AND e.user_id = l.user_id AND NOT e.inactive
    LEFT JOIN hrms.leave_types lt ON lt.tenant_id = l.tenant_id AND lt.id = l.leave_type_id AND NOT lt.inactive
    LEFT JOIN hrms.departments d ON d.tenant_id = e.tenant_id AND d.id = e.department_id AND NOT d.inactive
    WHERE la.tenant_id = sqlc.arg('tenant_id')
      AND la.status = 'pending'
      AND NOT la.inactive

    UNION ALL

    SELECT
        ar.tenant_id,
        ('attendance_request:' || ar.id::text) AS card_key,
        CASE WHEN ar.payroll_blocking THEN 'payroll_blockers' ELSE 'exceptions' END::text AS lane,
        'attendance'::text AS category,
        'attendance'::text AS source_module,
        'attendance_request'::text AS source_type,
        ar.id AS source_id,
        ar.user_id AS employee_user_id,
        ('Resolve attendance request for ' || trim(e.firstname || ' ' || COALESCE(e.lastname, '')))::text AS title,
        (COALESCE(ar.request_type, ar.requested_type, 'attendance') || ' on ' || ar.date::text || COALESCE(': ' || ar.reason, ''))::text AS summary,
        ar.status::text AS status,
        CASE
            WHEN ar.payroll_blocking THEN 'critical'
            WHEN ar.escalation_due_at IS NOT NULL AND ar.escalation_due_at <= NOW() THEN 'high'
            ELSE 'medium'
        END::text AS severity,
        CASE
            WHEN ar.payroll_blocking THEN 10
            WHEN ar.escalation_due_at IS NOT NULL AND ar.escalation_due_at <= NOW() THEN 25
            ELSE 45
        END::int AS priority,
        COALESCE(ar.escalation_due_at, ar.date::timestamptz) AS due_at,
        ar.created_at AS detected_at,
        CASE WHEN ar.payroll_blocking THEN 'Clear blocker' ELSE 'Review request' END::text AS action_label,
        'attendance'::text AS route_section,
        ar.id AS route_record_id,
        jsonb_build_object(
            'employee_code', e.employee_code,
            'request_type', COALESCE(ar.request_type, ar.requested_type),
            'route_mode', ar.route_mode,
            'payroll_blocking', ar.payroll_blocking
        ) AS metadata
    FROM hrms.attendance_requests ar
    JOIN hrms.employees e ON e.tenant_id = ar.tenant_id AND e.user_id = ar.user_id AND NOT e.inactive
    WHERE ar.tenant_id = sqlc.arg('tenant_id')
      AND ar.status = 'pending'
      AND NOT ar.inactive

    UNION ALL

    SELECT
        er.tenant_id,
        ('employee_exit:' || er.id::text) AS card_key,
        'exit'::text AS lane,
        'exit'::text AS category,
        'employee_exits'::text AS source_module,
        'employee_exit'::text AS source_type,
        er.id AS source_id,
        er.employee_user_id,
        ('Manage exit for ' || trim(e.firstname || ' ' || COALESCE(e.lastname, '')))::text AS title,
        (er.exit_type || ' exit, last working day ' || COALESCE(er.last_working_date::text, 'not set'))::text AS summary,
        er.status::text AS status,
        CASE
            WHEN COALESCE(tasks.blocked_tasks, 0) > 0 THEN 'high'
            WHEN er.last_working_date IS NOT NULL AND er.last_working_date <= CURRENT_DATE + INTERVAL '7 days' THEN 'medium'
            ELSE 'low'
        END::text AS severity,
        CASE
            WHEN COALESCE(tasks.blocked_tasks, 0) > 0 THEN 25
            WHEN er.last_working_date IS NOT NULL AND er.last_working_date <= CURRENT_DATE + INTERVAL '7 days' THEN 45
            ELSE 75
        END::int AS priority,
        er.last_working_date::timestamptz AS due_at,
        er.created_at AS detected_at,
        CASE WHEN er.status = 'submitted' THEN 'Review exit' ELSE 'Track clearance' END::text AS action_label,
        'employee-exits'::text AS route_section,
        er.id AS route_record_id,
        jsonb_build_object(
            'employee_code', e.employee_code,
            'total_tasks', COALESCE(tasks.total_tasks, 0),
            'completed_tasks', COALESCE(tasks.completed_tasks, 0),
            'blocked_tasks', COALESCE(tasks.blocked_tasks, 0)
        ) AS metadata
    FROM hrms.employee_exit_requests er
    JOIN hrms.employees e ON e.tenant_id = er.tenant_id AND e.id = er.employee_id AND NOT e.inactive
    LEFT JOIN LATERAL (
        SELECT
            COUNT(*) AS total_tasks,
            COUNT(*) FILTER (WHERE status IN ('completed','waived')) AS completed_tasks,
            COUNT(*) FILTER (WHERE status = 'blocked') AS blocked_tasks
        FROM hrms.employee_exit_tasks t
        WHERE t.tenant_id = er.tenant_id AND t.exit_request_id = er.id AND NOT t.inactive
    ) tasks ON TRUE
    WHERE er.tenant_id = sqlc.arg('tenant_id')
      AND er.status IN ('submitted','approved')
      AND NOT er.inactive

    UNION ALL

    SELECT
        cot.tenant_id,
        ('candidate_onboarding_task:' || cot.id::text) AS card_key,
        'joining'::text AS lane,
        'onboarding'::text AS category,
        'onboarding'::text AS source_module,
        'candidate_onboarding_task'::text AS source_type,
        cot.id AS source_id,
        NULL::uuid AS employee_user_id,
        ('Complete joining task: ' || ot.title)::text AS title,
        (trim(COALESCE(c.firstname, '') || ' ' || COALESCE(c.lastname, '')) || ' onboarding is ' || co.onboarding_status)::text AS summary,
        cot.status::text AS status,
        CASE
            WHEN cot.due_at IS NOT NULL AND cot.due_at < NOW() THEN 'high'
            WHEN cot.due_at IS NOT NULL AND cot.due_at <= NOW() + INTERVAL '3 days' THEN 'medium'
            ELSE 'low'
        END::text AS severity,
        CASE
            WHEN cot.due_at IS NOT NULL AND cot.due_at < NOW() THEN 25
            WHEN cot.due_at IS NOT NULL AND cot.due_at <= NOW() + INTERVAL '3 days' THEN 45
            ELSE 75
        END::int AS priority,
        cot.due_at,
        cot.created_at AS detected_at,
        'Open onboarding'::text AS action_label,
        'candidate-onboarding'::text AS route_section,
        co.id AS route_record_id,
        jsonb_build_object(
            'candidate_id', c.id,
            'candidate_email', c.email,
            'task_required', ot.is_required,
            'progress_percentage', co.progress_percentage
        ) AS metadata
    FROM hrms.candidate_onboarding_tasks cot
    JOIN hrms.candidate_onboardings co ON co.tenant_id = cot.tenant_id AND co.id = cot.candidate_onboarding_id AND NOT co.inactive
    JOIN hrms.onboarding_tasks ot ON ot.tenant_id = cot.tenant_id AND ot.id = cot.onboarding_task_id AND NOT ot.inactive
    JOIN hrms.candidates c ON c.tenant_id = cot.tenant_id AND c.id = co.candidate_id AND NOT c.inactive
    WHERE cot.tenant_id = sqlc.arg('tenant_id')
      AND cot.status IN ('Pending','InProgress')
      AND co.onboarding_status IN ('Pending','InProgress')
      AND NOT cot.inactive

    UNION ALL

    SELECT
        el.tenant_id,
        ('employee_letter:' || el.id::text) AS card_key,
        'approvals'::text AS lane,
        'documents'::text AS category,
        'employee_letters'::text AS source_module,
        'employee_letter'::text AS source_type,
        el.id AS source_id,
        el.user_id AS employee_user_id,
        (CASE WHEN el.status = 'Generated' THEN 'Approve letter: ' ELSE 'Track signed letter: ' END || el.subject)::text AS title,
        (trim(e.firstname || ' ' || COALESCE(e.lastname, '')) || ' - ' || el.letter_type)::text AS summary,
        el.status::text AS status,
        CASE WHEN el.status = 'Generated' THEN 'medium' ELSE 'low' END::text AS severity,
        CASE WHEN el.status = 'Generated' THEN 45 ELSE 70 END::int AS priority,
        COALESCE(el.effective_date, el.issue_date)::timestamptz AS due_at,
        el.created_at AS detected_at,
        CASE WHEN el.status = 'Generated' THEN 'Review letter' ELSE 'Track signature' END::text AS action_label,
        'employee-letters'::text AS route_section,
        el.id AS route_record_id,
        jsonb_build_object('employee_code', e.employee_code, 'letter_type', el.letter_type, 'requires_signature', el.signature_token IS NOT NULL) AS metadata
    FROM hrms.employee_letters el
    JOIN hrms.employees e ON e.tenant_id = el.tenant_id AND e.id = el.employee_id AND NOT e.inactive
    WHERE el.tenant_id = sqlc.arg('tenant_id')
      AND el.status IN ('Generated','Sent')
      AND NOT el.inactive

    UNION ALL

    SELECT
        a.tenant_id,
        ('agreement:' || a.id::text) AS card_key,
        'approvals'::text AS lane,
        'documents'::text AS category,
        'agreements'::text AS source_module,
        'agreement'::text AS source_type,
        a.id AS source_id,
        NULL::uuid AS employee_user_id,
        (CASE WHEN a.status = 'Generated' THEN 'Send agreement: ' ELSE 'Track agreement signature: ' END || a.title)::text AS title,
        COALESCE(wp.display_name, a.agreement_type)::text AS summary,
        a.status::text AS status,
        CASE WHEN a.status = 'Generated' THEN 'medium' ELSE 'low' END::text AS severity,
        CASE WHEN a.status = 'Generated' THEN 45 ELSE 70 END::int AS priority,
        COALESCE(a.effective_date, a.issue_date)::timestamptz AS due_at,
        a.created_at AS detected_at,
        CASE WHEN a.status = 'Generated' THEN 'Review agreement' ELSE 'Track signature' END::text AS action_label,
        'agreements'::text AS route_section,
        a.id AS route_record_id,
        jsonb_build_object('agreement_type', a.agreement_type, 'worker_code', wp.worker_code, 'requires_signature', a.signature_token IS NOT NULL) AS metadata
    FROM hrms.agreements a
    LEFT JOIN hrms.worker_profiles wp ON wp.tenant_id = a.tenant_id AND wp.id = a.worker_profile_id AND NOT wp.inactive
    WHERE a.tenant_id = sqlc.arg('tenant_id')
      AND a.status IN ('Generated','Sent')
      AND NOT a.inactive

    UNION ALL

    SELECT
        ci.tenant_id,
        ('compliance_item:' || ci.id::text) AS card_key,
        CASE WHEN cr.blocks_payroll THEN 'payroll_blockers' ELSE 'compliance' END::text AS lane,
        'compliance'::text AS category,
        'compliance'::text AS source_module,
        'compliance_checklist_item'::text AS source_type,
        ci.id AS source_id,
        NULL::uuid AS employee_user_id,
        ('Resolve compliance: ' || cr.title)::text AS title,
        COALESCE(wp.display_name, e.title, cr.category)::text AS summary,
        ci.status::text AS status,
        CASE
            WHEN cr.blocks_payroll OR ci.status IN ('expired','non_compliant') THEN 'high'
            WHEN ci.due_date IS NOT NULL AND ci.due_date <= CURRENT_DATE + INTERVAL '7 days' THEN 'medium'
            ELSE cr.severity
        END::text AS severity,
        CASE
            WHEN cr.blocks_payroll OR ci.status IN ('expired','non_compliant') THEN 20
            WHEN ci.due_date IS NOT NULL AND ci.due_date <= CURRENT_DATE + INTERVAL '7 days' THEN 40
            ELSE 65
        END::int AS priority,
        ci.due_date::timestamptz AS due_at,
        ci.created_at AS detected_at,
        CASE WHEN cr.blocks_payroll THEN 'Clear compliance blocker' ELSE 'Review compliance' END::text AS action_label,
        'compliance'::text AS route_section,
        ci.id AS route_record_id,
        jsonb_build_object('rule_id', cr.id, 'rule_code', cr.code, 'blocks_payroll', cr.blocks_payroll, 'requires_evidence', cr.requires_evidence) AS metadata
    FROM hrms.compliance_checklist_items ci
    JOIN hrms.compliance_rules cr ON cr.tenant_id = ci.tenant_id AND cr.id = ci.rule_id AND NOT cr.inactive
    LEFT JOIN hrms.worker_profiles wp ON wp.tenant_id = ci.tenant_id AND wp.id = ci.worker_profile_id AND NOT wp.inactive
    LEFT JOIN hrms.engagements e ON e.tenant_id = ci.tenant_id AND e.id = ci.engagement_id AND NOT e.inactive
    WHERE ci.tenant_id = sqlc.arg('tenant_id')
      AND ci.status IN ('pending','in_review','non_compliant','expired')
      AND NOT ci.inactive

    UNION ALL

    SELECT
        hc.tenant_id,
        ('hr_case:' || hc.id::text) AS card_key,
        'employee_requests'::text AS lane,
        'helpdesk'::text AS category,
        'hr_cases'::text AS source_module,
        'hr_case'::text AS source_type,
        hc.id AS source_id,
        COALESCE(hc.subject_employee_user_id, hc.requester_user_id) AS employee_user_id,
        (hc.case_number || ': ' || hc.title)::text AS title,
        (COALESCE(hcc.name, hc.case_type, 'HR case') || COALESCE(' - ' || NULLIF(trim(COALESCE(se.firstname, '') || ' ' || COALESCE(se.lastname, '')), ''), ''))::text AS summary,
        hc.status::text AS status,
        CASE
            WHEN hc.confidentiality_level IN ('sensitive','grievance') THEN 'critical'
            WHEN hc.status = 'escalated' OR (hc.due_at IS NOT NULL AND hc.due_at < NOW()) THEN 'high'
            WHEN hc.priority IN ('urgent','high') THEN 'medium'
            ELSE 'low'
        END::text AS severity,
        CASE
            WHEN hc.confidentiality_level IN ('sensitive','grievance') THEN 10
            WHEN hc.status = 'escalated' OR (hc.due_at IS NOT NULL AND hc.due_at < NOW()) THEN 20
            WHEN hc.priority = 'urgent' THEN 25
            WHEN hc.priority = 'high' THEN 40
            ELSE 65
        END::int AS priority,
        hc.due_at,
        hc.created_at AS detected_at,
        CASE WHEN hc.status IN ('new','open') THEN 'Start case' ELSE 'Handle case' END::text AS action_label,
        'hr-helpdesk'::text AS route_section,
        hc.id AS route_record_id,
        jsonb_build_object(
            'case_number', hc.case_number,
            'case_type', hc.case_type,
            'confidentiality_level', hc.confidentiality_level,
            'owner_role', hc.owner_role,
            'category_code', hcc.code,
            'requester_employee_code', re.employee_code,
            'subject_employee_code', se.employee_code
        ) AS metadata
    FROM hrms.hr_cases hc
    LEFT JOIN hrms.hr_case_categories hcc ON hcc.tenant_id = hc.tenant_id AND hcc.id = hc.category_id AND NOT hcc.inactive
    LEFT JOIN hrms.employees re ON re.tenant_id = hc.tenant_id AND re.user_id = hc.requester_user_id AND NOT re.inactive
    LEFT JOIN hrms.employees se ON se.tenant_id = hc.tenant_id AND se.user_id = hc.subject_employee_user_id AND NOT se.inactive
    WHERE hc.tenant_id = sqlc.arg('tenant_id')
      AND hc.status NOT IN ('resolved','closed','cancelled')
      AND NOT hc.inactive

    UNION ALL

    SELECT
        i.tenant_id,
        ('insight:' || i.id::text) AS card_key,
        'ai_recommendations'::text AS lane,
        'ai'::text AS category,
        'insights'::text AS source_module,
        'insight'::text AS source_type,
        i.id AS source_id,
        i.employee_user_id,
        i.title::text AS title,
        i.summary::text AS summary,
        i.status::text AS status,
        i.severity::text AS severity,
        CASE i.severity WHEN 'critical' THEN 10 WHEN 'high' THEN 20 WHEN 'medium' THEN 45 ELSE 75 END::int AS priority,
        i.due_at,
        i.detected_at,
        'Review insight'::text AS action_label,
        'insights'::text AS route_section,
        i.id AS route_record_id,
        jsonb_build_object('insight_type', i.insight_type, 'confidence_score', i.confidence_score, 'source', i.source) AS metadata
    FROM hrms.insights i
    WHERE i.tenant_id = sqlc.arg('tenant_id')
      AND i.status IN ('open','reviewing')
      AND NOT i.inactive

    UNION ALL

    SELECT
        aa.tenant_id,
        ('ai_action:' || aa.id::text) AS card_key,
        'ai_recommendations'::text AS lane,
        'ai'::text AS category,
        'ai_actions'::text AS source_module,
        'ai_agent_action'::text AS source_type,
        aa.id AS source_id,
        aa.employee_user_id,
        aa.title::text AS title,
        aa.summary::text AS summary,
        aa.status::text AS status,
        aa.severity::text AS severity,
        CASE aa.severity WHEN 'critical' THEN 10 WHEN 'high' THEN 20 WHEN 'medium' THEN 45 ELSE 75 END::int AS priority,
        aa.created_at + INTERVAL '2 days' AS due_at,
        aa.created_at AS detected_at,
        'Review recommendation'::text AS action_label,
        'insights'::text AS route_section,
        aa.id AS route_record_id,
        jsonb_build_object('agent_key', aa.agent_key, 'agent_name', aa.agent_name, 'requires_human_review', aa.requires_human_review, 'confidence_score', aa.confidence_score) AS metadata
    FROM hrms.ai_agent_action_logs aa
    WHERE aa.tenant_id = sqlc.arg('tenant_id')
      AND aa.status IN ('proposed','reviewing','failed')
      AND NOT aa.inactive
)
SELECT
    tenant_id,
    card_key::text AS card_key,
    lane,
    category,
    source_module,
    source_type,
    source_id,
    COALESCE(employee_user_id, '00000000-0000-0000-0000-000000000000'::uuid) AS employee_user_id,
    COALESCE(NULLIF(trim(title), ''), 'HR work item')::text AS title,
    COALESCE(NULLIF(trim(summary), ''), source_module, 'HRMS')::text AS summary,
    status,
    severity,
    priority,
    due_at,
    detected_at,
    action_label,
    route_section,
    route_record_id,
    metadata
FROM cards
WHERE (sqlc.narg('lane')::text IS NULL OR lane = sqlc.narg('lane')::text)
  AND (sqlc.narg('category')::text IS NULL OR category = sqlc.narg('category')::text)
  AND (sqlc.narg('severity')::text IS NULL OR severity = sqlc.narg('severity')::text)
  AND (
      sqlc.narg('search')::text IS NULL
      OR COALESCE(title, '') ILIKE '%' || sqlc.narg('search')::text || '%'
      OR COALESCE(summary, '') ILIKE '%' || sqlc.narg('search')::text || '%'
      OR source_module ILIKE '%' || sqlc.narg('search')::text || '%'
  )
ORDER BY priority ASC, due_at ASC NULLS LAST, detected_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');
