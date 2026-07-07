-- Align default tenant branding with Setika's current public site typography and
-- green/orange product palette. This updates only rows that still have the old
-- untouched defaults so tenant-customized branding is preserved.

ALTER TABLE hrms.tenant_brandings ALTER COLUMN secondary_color SET DEFAULT '#e87839';
ALTER TABLE hrms.tenant_brandings ALTER COLUMN tertiary_color SET DEFAULT '#f2b36d';
ALTER TABLE hrms.tenant_brandings ALTER COLUMN topbar_color SET DEFAULT '#fffaf4';
ALTER TABLE hrms.tenant_brandings ALTER COLUMN sidebar_color SET DEFAULT '#426b53';
ALTER TABLE hrms.tenant_brandings ALTER COLUMN font_family SET DEFAULT '"Plus Jakarta Sans", "Segoe UI", sans-serif';

UPDATE hrms.tenant_brandings
SET secondary_color = '#e87839'
WHERE secondary_color = '#2f6f7d';

UPDATE hrms.tenant_brandings
SET tertiary_color = '#f2b36d'
WHERE tertiary_color = '#e87839';

UPDATE hrms.tenant_brandings
SET topbar_color = '#fffaf4'
WHERE topbar_color = '#ffffff';

UPDATE hrms.tenant_brandings
SET sidebar_color = '#426b53'
WHERE sidebar_color = '#111827';

UPDATE hrms.tenant_brandings
SET font_family = '"Plus Jakarta Sans", "Segoe UI", sans-serif'
WHERE font_family = 'Inter, sans-serif';
