package pdf

import (
	"bytes"
	"context"
	"fmt"
	"html"
	"regexp"
	"strings"

	"github.com/jung-kurt/gofpdf"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

type EmployeeLetterRenderer struct{}

func NewEmployeeLetterRenderer() *EmployeeLetterRenderer {
	return &EmployeeLetterRenderer{}
}

func (r *EmployeeLetterRenderer) RenderEmployeeLetterPDF(ctx context.Context, doc ports.EmployeeLetterDocument) ([]byte, error) {
	if doc.Letter == nil {
		return nil, fmt.Errorf("employee letter is required")
	}
	title := "Employee Letter"
	if doc.Letter.Subject != nil && strings.TrimSpace(*doc.Letter.Subject) != "" {
		title = strings.TrimSpace(*doc.Letter.Subject)
	}
	p := gofpdf.New("P", "mm", "A4", "")
	p.SetMargins(18, 18, 18)
	p.SetAutoPageBreak(true, 18)
	p.AddPage()
	p.SetDrawColor(229, 231, 235)
	p.SetFillColor(88, 131, 104)
	p.Rect(18, 18, 174, 12, "F")
	p.SetTextColor(255, 255, 255)
	p.SetFont("Helvetica", "B", 14)
	p.CellFormat(174, 12, title, "", 1, "C", false, 0, "")
	p.SetTextColor(17, 24, 39)
	p.Ln(8)
	p.SetFont("Helvetica", "", 10)
	p.CellFormat(174, 7, "Employee: "+employeeLetterName(doc), "", 1, "L", false, 0, "")
	if doc.Letter.IssueDate != nil {
		p.CellFormat(174, 7, "Issue Date: "+doc.Letter.IssueDate.Format("02 Jan 2006"), "", 1, "L", false, 0, "")
	}
	p.Ln(5)
	p.SetFont("Helvetica", "", 11)
	for _, paragraph := range htmlParagraphs(doc.Letter.RenderedHTML) {
		p.MultiCell(174, 6, paragraph, "", "L", false)
		p.Ln(2)
	}
	if doc.Letter.SignatureCompletedAt != nil {
		p.Ln(8)
		p.SetFont("Helvetica", "B", 10)
		p.CellFormat(174, 7, "Digitally signed by "+employeeLetterValueOrDefault(doc.Letter.SignerName, "employee"), "", 1, "L", false, 0, "")
		if doc.Letter.SignatureHash != nil {
			p.SetFont("Helvetica", "", 8)
			p.MultiCell(174, 4, "Signature hash: "+*doc.Letter.SignatureHash, "", "L", false)
		}
	}
	var buf bytes.Buffer
	if err := p.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func employeeLetterName(doc ports.EmployeeLetterDocument) string {
	if doc.Employee != nil {
		name := strings.TrimSpace(strings.Join([]string{doc.Employee.Firstname, ptrString(doc.Employee.MiddleName), ptrString(doc.Employee.Lastname)}, " "))
		if doc.Employee.EmployeeCode != nil && *doc.Employee.EmployeeCode != "" {
			return fmt.Sprintf("%s (%s)", name, *doc.Employee.EmployeeCode)
		}
		if name != "" {
			return name
		}
	}
	if doc.Letter.EmployeeFirstname != nil || doc.Letter.EmployeeLastname != nil {
		return strings.TrimSpace(strings.Join([]string{ptrString(doc.Letter.EmployeeFirstname), ptrString(doc.Letter.EmployeeLastname)}, " "))
	}
	return doc.Letter.UserID.String()
}

var htmlTagPattern = regexp.MustCompile(`<[^>]+>`)

func htmlParagraphs(value *string) []string {
	if value == nil || strings.TrimSpace(*value) == "" {
		return []string{" "}
	}
	clean := strings.NewReplacer("</p>", "\n\n", "<br>", "\n", "<br/>", "\n", "<br />", "\n", "</div>", "\n\n").Replace(*value)
	clean = htmlTagPattern.ReplaceAllString(clean, "")
	clean = html.UnescapeString(clean)
	parts := strings.Split(clean, "\n")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			out = append(out, trimmed)
		}
	}
	if len(out) == 0 {
		return []string{" "}
	}
	return out
}

func employeeLetterValueOrDefault(value *string, fallback string) string {
	if value == nil || strings.TrimSpace(*value) == "" {
		return fallback
	}
	return strings.TrimSpace(*value)
}
