CREATE TABLE IF NOT EXISTS hrms.flex_pay_runs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    run_code TEXT NOT NULL,
    title TEXT NOT NULL,
    run_type TEXT NOT NULL DEFAULT 'mixed',
    status TEXT NOT NULL DEFAULT 'draft',
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    payout_date DATE,
    currency_code CHAR(3) NOT NULL DEFAULT 'INR',
    source_policy TEXT NOT NULL DEFAULT 'approved_sources',
    invoice_count INTEGER NOT NULL DEFAULT 0,
    item_count INTEGER NOT NULL DEFAULT 0,
    gross_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    tds_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    gst_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    net_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    generated_at TIMESTAMPTZ,
    submitted_at TIMESTAMPTZ,
    submitted_by UUID,
    approved_at TIMESTAMPTZ,
    approved_by UUID,
    rejected_at TIMESTAMPTZ,
    rejected_by UUID,
    paid_at TIMESTAMPTZ,
    paid_by UUID,
    payment_reference TEXT,
    export_batch_ref TEXT,
    notes TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    CONSTRAINT flex_pay_runs_code_check CHECK (btrim(run_code) <> ''),
    CONSTRAINT flex_pay_runs_title_check CHECK (btrim(title) <> ''),
    CONSTRAINT flex_pay_runs_type_check CHECK (run_type IN ('hourly', 'milestone', 'retainer', 'stipend', 'invoice', 'mixed')),
    CONSTRAINT flex_pay_runs_status_check CHECK (status IN ('draft', 'generated', 'submitted', 'approved', 'rejected', 'payment_pending', 'paid', 'cancelled')),
    CONSTRAINT flex_pay_runs_date_check CHECK (period_end >= period_start),
    CONSTRAINT flex_pay_runs_currency_check CHECK (currency_code ~ '^[A-Z]{3}$'),
    CONSTRAINT flex_pay_runs_amounts_check CHECK (gross_amount >= 0 AND tds_amount >= 0 AND gst_amount >= 0 AND net_amount >= 0),
    CONSTRAINT flex_pay_runs_metadata_object_check CHECK (jsonb_typeof(metadata) = 'object')
);

CREATE UNIQUE INDEX IF NOT EXISTS flex_pay_runs_tenant_code_active_uq
    ON hrms.flex_pay_runs (tenant_id, lower(run_code))
    WHERE inactive = FALSE;

CREATE INDEX IF NOT EXISTS flex_pay_runs_tenant_status_period_idx
    ON hrms.flex_pay_runs (tenant_id, status, period_start DESC, period_end DESC)
    WHERE inactive = FALSE;

CREATE TABLE IF NOT EXISTS hrms.contractor_invoices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    flex_pay_run_id UUID REFERENCES hrms.flex_pay_runs(id) ON DELETE SET NULL,
    worker_profile_id UUID NOT NULL REFERENCES hrms.worker_profiles(id) ON DELETE RESTRICT,
    engagement_id UUID REFERENCES hrms.engagements(id) ON DELETE SET NULL,
    invoice_number TEXT NOT NULL,
    invoice_date DATE NOT NULL,
    due_date DATE,
    status TEXT NOT NULL DEFAULT 'draft',
    vendor_name TEXT NOT NULL,
    vendor_gstin TEXT,
    place_of_supply TEXT,
    reverse_charge BOOLEAN NOT NULL DEFAULT FALSE,
    currency_code CHAR(3) NOT NULL DEFAULT 'INR',
    gross_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    tds_section TEXT,
    tds_rate NUMERIC(5,2) NOT NULL DEFAULT 0,
    tds_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    gst_rate NUMERIC(5,2) NOT NULL DEFAULT 0,
    gst_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    net_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    submitted_at TIMESTAMPTZ,
    submitted_by UUID,
    approved_at TIMESTAMPTZ,
    approved_by UUID,
    rejected_at TIMESTAMPTZ,
    rejected_by UUID,
    rejection_reason TEXT,
    paid_at TIMESTAMPTZ,
    paid_by UUID,
    payment_reference TEXT,
    attachment_path TEXT,
    notes TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    CONSTRAINT contractor_invoices_number_check CHECK (btrim(invoice_number) <> ''),
    CONSTRAINT contractor_invoices_vendor_check CHECK (btrim(vendor_name) <> ''),
    CONSTRAINT contractor_invoices_status_check CHECK (status IN ('draft', 'submitted', 'approved', 'rejected', 'payment_pending', 'paid', 'cancelled')),
    CONSTRAINT contractor_invoices_tds_section_check CHECK (tds_section IS NULL OR tds_section IN ('194C', '194J', 'none')),
    CONSTRAINT contractor_invoices_currency_check CHECK (currency_code ~ '^[A-Z]{3}$'),
    CONSTRAINT contractor_invoices_amounts_check CHECK (gross_amount >= 0 AND tds_rate >= 0 AND tds_amount >= 0 AND gst_rate >= 0 AND gst_amount >= 0 AND net_amount >= 0),
    CONSTRAINT contractor_invoices_metadata_object_check CHECK (jsonb_typeof(metadata) = 'object')
);

CREATE UNIQUE INDEX IF NOT EXISTS contractor_invoices_tenant_invoice_active_uq
    ON hrms.contractor_invoices (tenant_id, lower(invoice_number))
    WHERE inactive = FALSE;

CREATE INDEX IF NOT EXISTS contractor_invoices_tenant_run_status_idx
    ON hrms.contractor_invoices (tenant_id, flex_pay_run_id, status)
    WHERE inactive = FALSE;

CREATE INDEX IF NOT EXISTS contractor_invoices_tenant_worker_idx
    ON hrms.contractor_invoices (tenant_id, worker_profile_id, invoice_date DESC)
    WHERE inactive = FALSE;

CREATE TABLE IF NOT EXISTS hrms.flex_pay_run_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    flex_pay_run_id UUID NOT NULL REFERENCES hrms.flex_pay_runs(id) ON DELETE CASCADE,
    contractor_invoice_id UUID REFERENCES hrms.contractor_invoices(id) ON DELETE SET NULL,
    worker_profile_id UUID NOT NULL REFERENCES hrms.worker_profiles(id) ON DELETE RESTRICT,
    engagement_id UUID REFERENCES hrms.engagements(id) ON DELETE SET NULL,
    source_type TEXT NOT NULL,
    source_id UUID,
    description TEXT NOT NULL,
    quantity NUMERIC(12,2) NOT NULL DEFAULT 1,
    rate_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    gross_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    tds_section TEXT,
    tds_rate NUMERIC(5,2) NOT NULL DEFAULT 0,
    tds_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    gst_rate NUMERIC(5,2) NOT NULL DEFAULT 0,
    gst_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    net_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'draft',
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    CONSTRAINT flex_pay_run_items_source_check CHECK (source_type IN ('work_log', 'milestone', 'retainer', 'stipend', 'manual_invoice', 'adjustment')),
    CONSTRAINT flex_pay_run_items_status_check CHECK (status IN ('draft', 'submitted', 'approved', 'rejected', 'payment_pending', 'paid', 'cancelled')),
    CONSTRAINT flex_pay_run_items_description_check CHECK (btrim(description) <> ''),
    CONSTRAINT flex_pay_run_items_tds_section_check CHECK (tds_section IS NULL OR tds_section IN ('194C', '194J', 'none')),
    CONSTRAINT flex_pay_run_items_amounts_check CHECK (quantity >= 0 AND rate_amount >= 0 AND gross_amount >= 0 AND tds_rate >= 0 AND tds_amount >= 0 AND gst_rate >= 0 AND gst_amount >= 0 AND net_amount >= 0),
    CONSTRAINT flex_pay_run_items_metadata_object_check CHECK (jsonb_typeof(metadata) = 'object')
);

CREATE UNIQUE INDEX IF NOT EXISTS flex_pay_run_items_source_active_uq
    ON hrms.flex_pay_run_items (tenant_id, source_type, source_id)
    WHERE source_id IS NOT NULL AND inactive = FALSE;

CREATE INDEX IF NOT EXISTS flex_pay_run_items_tenant_run_idx
    ON hrms.flex_pay_run_items (tenant_id, flex_pay_run_id, status)
    WHERE inactive = FALSE;

CREATE TABLE IF NOT EXISTS hrms.flex_pay_run_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    flex_pay_run_id UUID REFERENCES hrms.flex_pay_runs(id) ON DELETE CASCADE,
    contractor_invoice_id UUID REFERENCES hrms.contractor_invoices(id) ON DELETE CASCADE,
    event_type TEXT NOT NULL,
    from_status TEXT,
    to_status TEXT,
    comment TEXT,
    actor_id UUID,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT flex_pay_run_events_type_check CHECK (btrim(event_type) <> ''),
    CONSTRAINT flex_pay_run_events_metadata_object_check CHECK (jsonb_typeof(metadata) = 'object')
);

CREATE INDEX IF NOT EXISTS flex_pay_run_events_tenant_run_created_idx
    ON hrms.flex_pay_run_events (tenant_id, flex_pay_run_id, created_at DESC);

CREATE INDEX IF NOT EXISTS flex_pay_run_events_tenant_invoice_created_idx
    ON hrms.flex_pay_run_events (tenant_id, contractor_invoice_id, created_at DESC);

DROP TRIGGER IF EXISTS trg_flex_pay_runs_updated_at ON hrms.flex_pay_runs;
CREATE TRIGGER trg_flex_pay_runs_updated_at BEFORE UPDATE ON hrms.flex_pay_runs FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();
DROP TRIGGER IF EXISTS trg_contractor_invoices_updated_at ON hrms.contractor_invoices;
CREATE TRIGGER trg_contractor_invoices_updated_at BEFORE UPDATE ON hrms.contractor_invoices FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();
DROP TRIGGER IF EXISTS trg_flex_pay_run_items_updated_at ON hrms.flex_pay_run_items;
CREATE TRIGGER trg_flex_pay_run_items_updated_at BEFORE UPDATE ON hrms.flex_pay_run_items FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();

ALTER TABLE hrms.flex_pay_runs ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.contractor_invoices ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.flex_pay_run_items ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.flex_pay_run_events ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS flex_pay_runs_tenant_isolation ON hrms.flex_pay_runs;
CREATE POLICY flex_pay_runs_tenant_isolation ON hrms.flex_pay_runs USING (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true)) WITH CHECK (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true));
DROP POLICY IF EXISTS contractor_invoices_tenant_isolation ON hrms.contractor_invoices;
CREATE POLICY contractor_invoices_tenant_isolation ON hrms.contractor_invoices USING (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true)) WITH CHECK (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true));
DROP POLICY IF EXISTS flex_pay_run_items_tenant_isolation ON hrms.flex_pay_run_items;
CREATE POLICY flex_pay_run_items_tenant_isolation ON hrms.flex_pay_run_items USING (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true)) WITH CHECK (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true));
DROP POLICY IF EXISTS flex_pay_run_events_tenant_isolation ON hrms.flex_pay_run_events;
CREATE POLICY flex_pay_run_events_tenant_isolation ON hrms.flex_pay_run_events USING (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true)) WITH CHECK (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true));
