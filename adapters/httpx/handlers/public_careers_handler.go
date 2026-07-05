package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

const careerPortalSource = "Career Portal"

type publicCareerPage struct {
	Tenant   publicCareerTenant      `json:"tenant"`
	Branding *domain.TenantBranding  `json:"branding,omitempty"`
	Content  publicCareerPageContent `json:"content"`
	Openings []*publicCareerOpening  `json:"openings"`
}

type publicCareerTenant struct {
	TenantID    uuid.UUID `json:"tenant_id"`
	Subdomain   string    `json:"subdomain"`
	DisplayName string    `json:"display_name"`
}

type publicCareerPageContent struct {
	Headline        string   `json:"headline"`
	WelcomeMessage  string   `json:"welcome_message"`
	About           string   `json:"about"`
	CoreValues      []string `json:"core_values"`
	Notices         []string `json:"notices"`
	SEOTitle        string   `json:"seo_title"`
	SEODescription  string   `json:"seo_description"`
	FeaturedJobIDs  []string `json:"featured_job_ids"`
	CandidateCTA    string   `json:"candidate_cta"`
	LoginButtonText string   `json:"login_button_text"`
}

type publicCareerOpening struct {
	*domain.JobPosting
	Slug      string `json:"slug"`
	PublicURL string `json:"public_url"`
}

type publicCareerApplicationRequest struct {
	JobPostingID      uuid.UUID `json:"job_posting_id"`
	Firstname         string    `json:"firstname"`
	Lastname          string    `json:"lastname"`
	Email             string    `json:"email"`
	Phone             string    `json:"phone"`
	DOB               string    `json:"dob"`
	ResumeURL         string    `json:"resume_url"`
	CoverLetter       string    `json:"cover_letter"`
	TotalExperience   *float64  `json:"total_experience,omitempty"`
	CurrentCompany    string    `json:"current_company"`
	ExpectedCTC       *float64  `json:"expected_ctc,omitempty"`
	NoticePeriod      *int32    `json:"notice_period,omitempty"`
	PreferredLocation string    `json:"preferred_location"`
}

type publicCareerApplicationResponse struct {
	Candidate   *domain.Candidate            `json:"candidate"`
	Application *domain.CandidateApplication `json:"application"`
}

func (h *Handler) GetPublicCareers(w http.ResponseWriter, r *http.Request) {
	branding, err := h.resolveCareerBranding(r)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "resolve public careers tenant", err, "tenant careers page not found")
		return
	}
	openings, err := h.svc.ListPublishedJobPostings(r.Context(), branding.TenantID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list public career openings", err, "failed to list career openings")
		return
	}
	content := h.publicCareerContent(r, branding)
	publicOpenings := h.publicCareerOpenings(r, openings)
	respondJSON(w, http.StatusOK, publicCareerPage{
		Tenant: publicCareerTenant{
			TenantID:    branding.TenantID,
			Subdomain:   branding.Subdomain,
			DisplayName: h.careerTenantName(branding),
		},
		Branding: branding,
		Content:  content,
		Openings: publicOpenings,
	})
}

func (h *Handler) ApplyPublicCareerPosting(w http.ResponseWriter, r *http.Request) {
	branding, err := h.resolveCareerBranding(r)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "resolve public careers tenant", err, "tenant careers page not found")
		return
	}
	var req publicCareerApplicationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode public career application", err, "invalid request body")
		return
	}
	if req.JobPostingID == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, "validate public career application", domain.ErrInvalidJobPostingID, "job posting is required")
		return
	}
	posting, err := h.svc.GetJobPosting(r.Context(), branding.TenantID, req.JobPostingID)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get public career posting", err, "job opening not found")
		return
	}
	if !posting.IsPublished || posting.JobStatus == nil || *posting.JobStatus != domain.JobPostingStatusOpen {
		h.respondError(w, r, http.StatusBadRequest, "validate public career posting", errors.New("job opening is not accepting applications"), "job opening is not accepting applications")
		return
	}
	firstname, lastname, email := strings.TrimSpace(req.Firstname), strings.TrimSpace(req.Lastname), strings.TrimSpace(req.Email)
	if firstname == "" || lastname == "" || email == "" {
		h.respondError(w, r, http.StatusBadRequest, "validate public career applicant", domain.ErrInvalidCandidateName, "first name, last name, and email are required")
		return
	}
	dob, ok := parsePublicDate(w, r, h, "validate public career applicant", req.DOB, "date of birth")
	if !ok {
		return
	}
	source := careerPortalSource
	candidate, err := h.svc.CreateCandidate(r.Context(), ports.CandidateCommand{
		TenantID:          branding.TenantID,
		Firstname:         publicCareerStringPtr(firstname),
		Lastname:          publicCareerStringPtr(lastname),
		Email:             publicCareerStringPtr(email),
		Phone:             publicCareerStringPtr(req.Phone),
		DOB:               dob,
		TotalExperience:   req.TotalExperience,
		CurrentCompany:    publicCareerStringPtr(req.CurrentCompany),
		ExpectedSalary:    req.ExpectedCTC,
		NoticePeriod:      req.NoticePeriod,
		PreferredLocation: publicCareerStringPtr(req.PreferredLocation),
		Source:            &source,
		ResumeURL:         publicCareerStringPtr(req.ResumeURL),
	})
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "create public career candidate", err, err.Error())
		return
	}
	status := domain.CandidateApplicationStatusNew
	appliedAt := time.Now().UTC()
	application, err := h.svc.CreateCandidateApplication(r.Context(), ports.CandidateApplicationCommand{
		TenantID:     branding.TenantID,
		CandidateID:  &candidate.ID,
		JobPostingID: &req.JobPostingID,
		ResumeURL:    publicCareerStringPtr(req.ResumeURL),
		CoverLetter:  publicCareerStringPtr(req.CoverLetter),
		ExpectedCTC:  req.ExpectedCTC,
		NoticePeriod: req.NoticePeriod,
		Source:       &source,
		SourceDetail: publicCareerStringPtr(r.Host),
		Status:       &status,
		AppliedAt:    &appliedAt,
	})
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "create public career application", err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, publicCareerApplicationResponse{Candidate: candidate, Application: application})
}

func (h *Handler) resolveCareerBranding(r *http.Request) (*domain.TenantBranding, error) {
	host := strings.TrimSpace(r.URL.Query().Get("host"))
	if host == "" {
		host = r.Host
	}
	return h.svc.ResolveTenantBrandingByHost(r.Context(), host)
}

func (h *Handler) publicCareerContent(r *http.Request, branding *domain.TenantBranding) publicCareerPageContent {
	tenantName := h.careerTenantName(branding)
	content := publicCareerPageContent{
		Headline:        "Careers at " + tenantName,
		WelcomeMessage:  "Welcome to " + tenantName + ". Explore our current openings and apply directly.",
		About:           "Build meaningful work with a team that values ownership, clarity, and service.",
		CoreValues:      []string{"Ownership", "Customer trust", "Learning", "Respect"},
		Notices:         []string{},
		SEOTitle:        tenantName + " Careers and Job Openings",
		SEODescription:  "Explore current job openings at " + tenantName + " and apply online.",
		CandidateCTA:    "View openings",
		LoginButtonText: "Login",
	}
	setting, err := h.svc.GetTenantSetting(r.Context(), branding.TenantID, "careers")
	if err != nil || setting == nil {
		return content
	}
	if value, ok := setting.Value["headline"].(string); ok && strings.TrimSpace(value) != "" {
		content.Headline = strings.TrimSpace(value)
	}
	if value, ok := setting.Value["welcome_message"].(string); ok && strings.TrimSpace(value) != "" {
		content.WelcomeMessage = strings.TrimSpace(value)
	}
	if value, ok := setting.Value["about"].(string); ok && strings.TrimSpace(value) != "" {
		content.About = strings.TrimSpace(value)
	}
	if values := publicCareerStringSlice(setting.Value["core_values"]); len(values) > 0 {
		content.CoreValues = values
	}
	if values := publicCareerStringSlice(setting.Value["notices"]); len(values) > 0 {
		content.Notices = values
	}
	if values := publicCareerStringSlice(setting.Value["featured_job_ids"]); len(values) > 0 {
		content.FeaturedJobIDs = values
	}
	if value, ok := setting.Value["seo_title"].(string); ok && strings.TrimSpace(value) != "" {
		content.SEOTitle = strings.TrimSpace(value)
	}
	if value, ok := setting.Value["seo_description"].(string); ok && strings.TrimSpace(value) != "" {
		content.SEODescription = strings.TrimSpace(value)
	}
	if value, ok := setting.Value["candidate_cta"].(string); ok && strings.TrimSpace(value) != "" {
		content.CandidateCTA = strings.TrimSpace(value)
	}
	if value, ok := setting.Value["login_button_text"].(string); ok && strings.TrimSpace(value) != "" {
		content.LoginButtonText = strings.TrimSpace(value)
	}
	return content
}

func (h *Handler) publicCareerOpenings(r *http.Request, openings []*domain.JobPosting) []*publicCareerOpening {
	baseURL := publicCareerBaseURL(r)
	result := make([]*publicCareerOpening, 0, len(openings))
	for _, opening := range openings {
		if opening == nil {
			continue
		}
		slug := publicCareerJobSlug(opening)
		result = append(result, &publicCareerOpening{
			JobPosting: opening,
			Slug:       slug,
			PublicURL:  baseURL + "/jobs/" + slug,
		})
	}
	return result
}

func (h *Handler) careerTenantName(branding *domain.TenantBranding) string {
	if branding != nil && branding.DisplayName != nil && strings.TrimSpace(*branding.DisplayName) != "" {
		return strings.TrimSpace(*branding.DisplayName)
	}
	if branding != nil && strings.TrimSpace(branding.Subdomain) != "" {
		return strings.TrimSpace(branding.Subdomain)
	}
	return "Setika"
}

func parsePublicDate(w http.ResponseWriter, r *http.Request, h *Handler, operation string, raw string, field string) (*time.Time, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, true
	}
	parsed, err := time.Parse("2006-01-02", raw)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid "+field)
		return nil, false
	}
	return &parsed, true
}

func publicCareerStringPtr(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return &value
}

func publicCareerStringSlice(raw any) []string {
	values, ok := raw.([]any)
	if !ok {
		return nil
	}
	clean := make([]string, 0, len(values))
	for _, item := range values {
		if value, ok := item.(string); ok && strings.TrimSpace(value) != "" {
			clean = append(clean, strings.TrimSpace(value))
		}
	}
	return clean
}

func publicCareerBaseURL(r *http.Request) string {
	proto := strings.TrimSpace(r.Header.Get("X-Forwarded-Proto"))
	if proto == "" {
		proto = "https"
		if r.TLS == nil && strings.HasPrefix(r.Host, "localhost") {
			proto = "http"
		}
	}
	host := strings.TrimSpace(r.URL.Query().Get("host"))
	if host == "" {
		host = r.Host
	}
	return proto + "://" + strings.TrimRight(host, "/")
}

func publicCareerJobSlug(opening *domain.JobPosting) string {
	title := "job"
	if opening.Title != nil && strings.TrimSpace(*opening.Title) != "" {
		title = strings.TrimSpace(*opening.Title)
	}
	base := slugifyPublicCareerText(title)
	if base == "" {
		base = "job"
	}
	return base + "-" + opening.ID.String()
}

func slugifyPublicCareerText(value string) string {
	var b strings.Builder
	lastDash := false
	for _, r := range strings.ToLower(value) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
			lastDash = false
			continue
		}
		if !lastDash {
			b.WriteByte('-')
			lastDash = true
		}
	}
	return strings.Trim(b.String(), "-")
}
