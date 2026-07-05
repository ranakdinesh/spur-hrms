package domain

import "github.com/google/uuid"

type ReportDataset struct {
	TenantID    uuid.UUID       `json:"tenant_id"`
	ReportCode  string          `json:"report_code"`
	Title       string          `json:"title"`
	Description string          `json:"description,omitempty"`
	PeriodLabel string          `json:"period_label,omitempty"`
	Columns     []ReportColumn  `json:"columns"`
	Rows        [][]string      `json:"rows"`
	Summary     []ReportMetric  `json:"summary,omitempty"`
	Branding    *ReportBranding `json:"branding,omitempty"`
}

type ReportColumn struct {
	Key   string `json:"key"`
	Label string `json:"label"`
}

type ReportMetric struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type ReportBranding struct {
	DisplayName    string `json:"display_name"`
	PrimaryColor   string `json:"primary_color"`
	SecondaryColor string `json:"secondary_color"`
	LogoPath       string `json:"logo_path,omitempty"`
}
