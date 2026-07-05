package services

import (
	"context"
	"html"
	"strings"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type renderedNotificationEmail struct {
	Subject  string
	TextBody string
	HTMLBody *string
}

func (s *TenantService) renderNotificationEmail(ctx context.Context, log *domain.NotificationLog, employee *domain.EmployeeListItem) renderedNotificationEmail {
	rendered := renderedNotificationEmail{}
	if log == nil {
		return rendered
	}
	var master *domain.NotificationMaster
	if log.NotificationMasterID != nil {
		item, err := s.notifications.GetNotificationMaster(ctx, log.TenantID, *log.NotificationMasterID)
		if err != nil {
			s.log.Warn().Err(err).Str("tenant_id", log.TenantID.String()).Str("notification_master_id", log.NotificationMasterID.String()).Msg("hrms: notification email template fallback used")
		} else {
			master = item
		}
	}
	branding := s.notificationEmailBranding(ctx, log)
	values := notificationEmailTemplateValues(log, master, employee, branding)

	subjectTemplate := "{{notification_title}}"
	textTemplate := "{{notification_message}}"
	htmlTemplate := ""
	if master != nil {
		if master.EmailSubjectTemplate != nil && strings.TrimSpace(*master.EmailSubjectTemplate) != "" {
			subjectTemplate = *master.EmailSubjectTemplate
		}
		if master.EmailTextTemplate != nil && strings.TrimSpace(*master.EmailTextTemplate) != "" {
			textTemplate = *master.EmailTextTemplate
		}
		if master.EmailHTMLTemplate != nil && strings.TrimSpace(*master.EmailHTMLTemplate) != "" {
			htmlTemplate = *master.EmailHTMLTemplate
		}
	}
	rendered.Subject = renderNotificationTemplate(subjectTemplate, values)
	rendered.TextBody = renderNotificationTemplate(textTemplate, values)
	if rendered.Subject == "" {
		rendered.Subject = values["notification_title"]
	}
	if rendered.TextBody == "" {
		rendered.TextBody = values["notification_message"]
	}

	bodyHTML := notificationTextToHTML(rendered.TextBody)
	if htmlTemplate != "" {
		bodyHTML = renderNotificationTemplate(htmlTemplate, values)
	}
	htmlBody := brandedNotificationHTML(branding, rendered.Subject, bodyHTML)
	rendered.HTMLBody = &htmlBody
	return rendered
}

type notificationEmailBrand struct {
	TenantName     string
	PrimaryColor   string
	SecondaryColor string
	LogoURL        string
}

func (s *TenantService) notificationEmailBranding(ctx context.Context, log *domain.NotificationLog) notificationEmailBrand {
	brand := notificationEmailBrand{TenantName: "Setika", PrimaryColor: "#588368", SecondaryColor: "#2f6f7d"}
	if log == nil || log.TenantID == uuid.Nil || s.branding == nil {
		return brand
	}
	branding, err := s.branding.GetTenantBranding(ctx, log.TenantID)
	if err != nil || branding == nil {
		if err != nil {
			s.log.Warn().Err(err).Str("tenant_id", log.TenantID.String()).Msg("hrms: tenant email branding fallback used")
		}
		return brand
	}
	if branding.DisplayName != nil && strings.TrimSpace(*branding.DisplayName) != "" {
		brand.TenantName = strings.TrimSpace(*branding.DisplayName)
	}
	if strings.TrimSpace(branding.PrimaryColor) != "" {
		brand.PrimaryColor = strings.TrimSpace(branding.PrimaryColor)
	}
	if strings.TrimSpace(branding.SecondaryColor) != "" {
		brand.SecondaryColor = strings.TrimSpace(branding.SecondaryColor)
	}
	if branding.LogoPath != nil && isEmailSafeURL(*branding.LogoPath) {
		brand.LogoURL = strings.TrimSpace(*branding.LogoPath)
	}
	return brand
}

func notificationEmailTemplateValues(log *domain.NotificationLog, master *domain.NotificationMaster, employee *domain.EmployeeListItem, brand notificationEmailBrand) map[string]string {
	title := ""
	if log.Subject != nil {
		title = strings.TrimSpace(*log.Subject)
	}
	message := ""
	if log.Body != nil {
		message = strings.TrimSpace(*log.Body)
	}
	code := ""
	if master != nil {
		code = master.Code
	}
	employeeName := "there"
	employeeEmail := ""
	employeeCode := ""
	if employee != nil {
		employeeName = strings.TrimSpace(strings.Join(employeeNameParts(employee), " "))
		if employeeName == "" {
			employeeName = "there"
		}
		if employee.Email != nil {
			employeeEmail = strings.TrimSpace(*employee.Email)
		}
		if employee.EmployeeCode != nil {
			employeeCode = strings.TrimSpace(*employee.EmployeeCode)
		}
	}
	return map[string]string{
		"tenant_name":           brand.TenantName,
		"notification_title":    title,
		"notification_message":  message,
		"notification_code":     code,
		"employee_name":         employeeName,
		"employee_email":        employeeEmail,
		"employee_code":         employeeCode,
		"brand_primary_color":   brand.PrimaryColor,
		"brand_secondary_color": brand.SecondaryColor,
	}
}

func employeeNameParts(employee *domain.EmployeeListItem) []string {
	if employee == nil {
		return nil
	}
	parts := []string{employee.Firstname}
	if employee.MiddleName != nil {
		parts = append(parts, *employee.MiddleName)
	}
	if employee.Lastname != nil {
		parts = append(parts, *employee.Lastname)
	}
	clean := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			clean = append(clean, trimmed)
		}
	}
	return clean
}

func renderNotificationTemplate(template string, values map[string]string) string {
	rendered := template
	for key, value := range values {
		rendered = strings.ReplaceAll(rendered, "{{"+key+"}}", value)
		rendered = strings.ReplaceAll(rendered, "{{ "+key+" }}", value)
	}
	return strings.TrimSpace(rendered)
}

func notificationTextToHTML(value string) string {
	paragraphs := strings.Split(value, "\n\n")
	var builder strings.Builder
	for _, paragraph := range paragraphs {
		trimmed := strings.TrimSpace(paragraph)
		if trimmed == "" {
			continue
		}
		escaped := html.EscapeString(trimmed)
		escaped = strings.ReplaceAll(escaped, "\n", "<br>")
		builder.WriteString(`<p style="margin:0 0 14px;color:#374151;font-size:14px;line-height:1.7;">`)
		builder.WriteString(escaped)
		builder.WriteString(`</p>`)
	}
	if builder.Len() == 0 {
		return `<p style="margin:0;color:#374151;font-size:14px;line-height:1.7;">Notification</p>`
	}
	return builder.String()
}

func brandedNotificationHTML(brand notificationEmailBrand, subject string, bodyHTML string) string {
	logo := `<div style="font-size:18px;font-weight:800;color:#111827;">` + html.EscapeString(brand.TenantName) + `</div>`
	if brand.LogoURL != "" {
		logo = `<img src="` + html.EscapeString(brand.LogoURL) + `" alt="` + html.EscapeString(brand.TenantName) + `" style="max-height:44px;max-width:180px;display:block;">`
	}
	return `<!doctype html><html><body style="margin:0;background:#f6f7f8;font-family:Inter,Arial,sans-serif;color:#111827;">
<table role="presentation" width="100%" cellspacing="0" cellpadding="0" style="background:#f6f7f8;padding:24px 12px;">
<tr><td align="center">
<table role="presentation" width="100%" cellspacing="0" cellpadding="0" style="max-width:640px;background:#ffffff;border:1px solid #e5e7eb;border-radius:12px;overflow:hidden;">
<tr><td style="padding:24px 28px;border-top:5px solid ` + html.EscapeString(brand.PrimaryColor) + `;">` + logo + `</td></tr>
<tr><td style="padding:0 28px 8px;"><h1 style="margin:0;color:#111827;font-size:22px;line-height:1.35;">` + html.EscapeString(subject) + `</h1></td></tr>
<tr><td style="padding:10px 28px 28px;">` + bodyHTML + `</td></tr>
<tr><td style="padding:18px 28px;background:#f9fafb;border-top:1px solid #edf1ef;color:#6b7280;font-size:12px;line-height:1.5;">This email was sent by ` + html.EscapeString(brand.TenantName) + ` through Setika HRMS.</td></tr>
</table>
</td></tr>
</table>
</body></html>`
}

func isEmailSafeURL(value string) bool {
	clean := strings.TrimSpace(strings.ToLower(value))
	return strings.HasPrefix(clean, "https://") || strings.HasPrefix(clean, "http://")
}
