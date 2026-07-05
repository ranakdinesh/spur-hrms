ALTER TABLE hrms.employees
    ADD COLUMN IF NOT EXISTS middle_name VARCHAR(100);

ALTER TABLE hrms.document_types
    ADD COLUMN IF NOT EXISTS description TEXT,
    ADD COLUMN IF NOT EXISTS is_required BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS instructions TEXT,
    ADD COLUMN IF NOT EXISTS allowed_content_types TEXT NOT NULL DEFAULT 'application/pdf,image/jpeg,image/png,image/webp,application/msword,application/vnd.openxmlformats-officedocument.wordprocessingml.document',
    ADD COLUMN IF NOT EXISTS max_file_size_bytes BIGINT NOT NULL DEFAULT 10485760,
    ADD COLUMN IF NOT EXISTS display_order INT NOT NULL DEFAULT 0;

ALTER TABLE hrms.employee_documents
    ADD COLUMN IF NOT EXISTS status VARCHAR(30) NOT NULL DEFAULT 'pending_review',
    ADD COLUMN IF NOT EXISTS review_remarks TEXT,
    ADD COLUMN IF NOT EXISTS reviewed_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS reviewed_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS original_file_name VARCHAR(255),
    ADD COLUMN IF NOT EXISTS content_type VARCHAR(255),
    ADD COLUMN IF NOT EXISTS file_size_bytes BIGINT,
    ADD COLUMN IF NOT EXISTS encrypted BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS encryption_algorithm VARCHAR(50) NOT NULL DEFAULT 'AES-256-GCM';

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'employee_documents_status_check'
          AND conrelid = 'hrms.employee_documents'::regclass
    ) THEN
        ALTER TABLE hrms.employee_documents
            ADD CONSTRAINT employee_documents_status_check
            CHECK (status IN ('pending_review', 'approved', 'rejected', 'resubmission_requested'));
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS employee_documents_tenant_status_idx ON hrms.employee_documents(tenant_id, status) WHERE NOT inactive;
