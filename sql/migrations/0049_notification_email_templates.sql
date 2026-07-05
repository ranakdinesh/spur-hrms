ALTER TABLE hrms.notification_masters
    ADD COLUMN IF NOT EXISTS email_subject_template TEXT,
    ADD COLUMN IF NOT EXISTS email_text_template TEXT,
    ADD COLUMN IF NOT EXISTS email_html_template TEXT;
