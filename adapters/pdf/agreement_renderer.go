package pdf

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/jung-kurt/gofpdf"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

type AgreementRenderer struct{}

func NewAgreementRenderer() *AgreementRenderer {
	return &AgreementRenderer{}
}

func (r *AgreementRenderer) RenderAgreementPDF(ctx context.Context, doc ports.AgreementDocument) ([]byte, error) {
	if doc.Agreement == nil {
		return nil, fmt.Errorf("agreement is required")
	}
	title := doc.Agreement.Title
	if doc.Agreement.Subject != nil && strings.TrimSpace(*doc.Agreement.Subject) != "" {
		title = strings.TrimSpace(*doc.Agreement.Subject)
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
	p.CellFormat(174, 7, "Agreement Type: "+strings.Title(strings.ReplaceAll(doc.Agreement.AgreementType, "_", " ")), "", 1, "L", false, 0, "")
	if doc.Agreement.EffectiveDate != nil {
		p.CellFormat(174, 7, "Effective Date: "+doc.Agreement.EffectiveDate.Format("02 Jan 2006"), "", 1, "L", false, 0, "")
	}
	p.Ln(5)
	p.SetFont("Helvetica", "", 11)
	for _, paragraph := range htmlParagraphs(doc.Agreement.RenderedHTML) {
		p.MultiCell(174, 6, paragraph, "", "L", false)
		p.Ln(2)
	}
	if doc.Agreement.SignatureCompletedAt != nil {
		p.Ln(8)
		p.SetFont("Helvetica", "B", 10)
		p.CellFormat(174, 7, "Digitally signed by "+employeeLetterValueOrDefault(doc.Agreement.SignerName, "signer"), "", 1, "L", false, 0, "")
		if doc.Agreement.SignatureHash != nil {
			p.SetFont("Helvetica", "", 8)
			p.MultiCell(174, 4, "Signature hash: "+*doc.Agreement.SignatureHash, "", "L", false)
		}
	}
	var buf bytes.Buffer
	if err := p.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
