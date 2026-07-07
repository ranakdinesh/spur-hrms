-- Older seeded/demo branding rows sometimes stored a color value in
-- sidebar_background. The field represents a background mode, not a color.

UPDATE hrms.tenant_brandings
SET sidebar_background = 'solid'
WHERE sidebar_background = '#111827';
