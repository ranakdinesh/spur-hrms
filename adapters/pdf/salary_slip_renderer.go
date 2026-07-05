package pdf

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jung-kurt/gofpdf"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

type SalarySlipRenderer struct{}

func NewSalarySlipRenderer() *SalarySlipRenderer {
	return &SalarySlipRenderer{}
}

func (r *SalarySlipRenderer) RenderSalarySlipPDF(ctx context.Context, doc ports.SalarySlipDocument) ([]byte, error) {
	if doc.Slip == nil {
		return nil, fmt.Errorf("salary slip is required")
	}
	title := "Salary Slip"
	primary := 17
	accentR, accentG, accentB := 88, 131, 104
	if doc.Format != nil {
		if strings.TrimSpace(doc.Format.Title) != "" {
			title = doc.Format.Title
		}
		accentR, accentG, accentB = hexColor(doc.Format.AccentColor, accentR, accentG, accentB)
		primary, _, _ = hexColor(doc.Format.PrimaryColor, primary, 24, 39)
	}

	p := gofpdf.New("P", "mm", "A4", "")
	p.SetMargins(14, 14, 14)
	p.SetAutoPageBreak(true, 14)
	p.AddPage()
	p.SetDrawColor(229, 231, 235)
	p.SetFillColor(accentR, accentG, accentB)
	p.Rect(14, 14, 182, 11, "F")
	p.SetTextColor(255, 255, 255)
	p.SetFont("Helvetica", "B", 15)
	p.CellFormat(182, 11, title, "", 1, "C", false, 0, "")
	p.SetTextColor(primary, primary, primary)
	if doc.Format != nil && doc.Format.Subtitle != nil {
		p.SetFont("Helvetica", "", 9)
		p.CellFormat(182, 7, *doc.Format.Subtitle, "", 1, "C", false, 0, "")
	}
	p.Ln(4)

	p.SetFont("Helvetica", "B", 11)
	p.CellFormat(91, 7, "Employee", "1", 0, "L", false, 0, "")
	p.CellFormat(91, 7, "Period", "1", 1, "L", false, 0, "")
	p.SetFont("Helvetica", "", 10)
	p.CellFormat(91, 8, employeeLabel(doc), "1", 0, "L", false, 0, "")
	p.CellFormat(91, 8, fmt.Sprintf("%02d/%04d", doc.Slip.Month, doc.Slip.Year), "1", 1, "L", false, 0, "")
	p.Ln(4)

	p.SetFont("Helvetica", "B", 11)
	p.CellFormat(45.5, 7, "Gross", "1", 0, "L", false, 0, "")
	p.CellFormat(45.5, 7, "Earnings", "1", 0, "L", false, 0, "")
	p.CellFormat(45.5, 7, "Deductions", "1", 0, "L", false, 0, "")
	p.CellFormat(45.5, 7, "Net", "1", 1, "L", false, 0, "")
	p.SetFont("Helvetica", "", 10)
	p.CellFormat(45.5, 8, money(doc.Slip.GrossSalary), "1", 0, "L", false, 0, "")
	p.CellFormat(45.5, 8, money(doc.Slip.TotalEarnings), "1", 0, "L", false, 0, "")
	p.CellFormat(45.5, 8, money(doc.Slip.TotalDeductions), "1", 0, "L", false, 0, "")
	p.CellFormat(45.5, 8, money(doc.Slip.NetSalary), "1", 1, "L", false, 0, "")
	p.Ln(4)

	p.SetFont("Helvetica", "B", 11)
	p.CellFormat(182, 8, "Salary Components", "1", 1, "L", false, 0, "")
	p.SetFont("Helvetica", "B", 9)
	p.CellFormat(92, 7, "Component", "1", 0, "L", false, 0, "")
	p.CellFormat(45, 7, "Type", "1", 0, "L", false, 0, "")
	p.CellFormat(45, 7, "Amount", "1", 1, "R", false, 0, "")
	p.SetFont("Helvetica", "", 9)
	for _, item := range doc.Slip.Items {
		p.CellFormat(92, 7, item.Name, "1", 0, "L", false, 0, "")
		p.CellFormat(45, 7, strings.Title(strings.ReplaceAll(item.ItemType, "_", " ")), "1", 0, "L", false, 0, "")
		p.CellFormat(45, 7, money(item.Amount), "1", 1, "R", false, 0, "")
	}

	if doc.Format == nil || doc.Format.ShowLeaveBalance {
		p.Ln(4)
		p.SetFont("Helvetica", "B", 11)
		p.CellFormat(182, 8, "Leave Balance", "1", 1, "L", false, 0, "")
		p.SetFont("Helvetica", "B", 9)
		p.CellFormat(92, 7, "Leave", "1", 0, "L", false, 0, "")
		p.CellFormat(30, 7, "Total", "1", 0, "R", false, 0, "")
		p.CellFormat(30, 7, "Used", "1", 0, "R", false, 0, "")
		p.CellFormat(30, 7, "Balance", "1", 1, "R", false, 0, "")
		p.SetFont("Helvetica", "", 9)
		for _, leave := range doc.Slip.Leaves {
			name := leave.LeaveTypeID.String()
			if leave.LeaveTypeName != nil && *leave.LeaveTypeName != "" {
				name = *leave.LeaveTypeName
			}
			p.CellFormat(92, 7, name, "1", 0, "L", false, 0, "")
			p.CellFormat(30, 7, fmt.Sprintf("%.1f", leave.TotalDays), "1", 0, "R", false, 0, "")
			p.CellFormat(30, 7, fmt.Sprintf("%.1f", leave.UsedDays), "1", 0, "R", false, 0, "")
			p.CellFormat(30, 7, fmt.Sprintf("%.1f", leave.BalanceDays), "1", 1, "R", false, 0, "")
		}
	}
	p.Ln(5)
	p.SetFont("Helvetica", "", 9)
	p.CellFormat(182, 6, fmt.Sprintf("Days: present %d, absent %d, total %d. LWP deduction %s.", doc.Slip.PresentDays, doc.Slip.AbsentDays, doc.Slip.TotalDays, money(doc.Slip.AbsentDeduction)), "", 1, "L", false, 0, "")
	if doc.Format != nil && doc.Format.FooterText != nil {
		p.MultiCell(182, 5, *doc.Format.FooterText, "", "C", false)
	}
	var buf bytes.Buffer
	if err := p.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

type LocalSalarySlipStorage struct {
	root string
}

func NewLocalSalarySlipStorage(root string) *LocalSalarySlipStorage {
	if strings.TrimSpace(root) == "" {
		root = filepath.Join(os.TempDir(), "setika-hrms", "salary-slips")
	}
	return &LocalSalarySlipStorage{root: root}
}

func (s *LocalSalarySlipStorage) StoreSalarySlipPDF(ctx context.Context, input ports.StoreSalarySlipPDFInput) (string, error) {
	dir := filepath.Join(s.root, input.TenantID.String(), input.UserID.String(), fmt.Sprintf("%04d", input.Year))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	name := safeFileName(input.FileName)
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, input.Content, 0o644); err != nil {
		return "", err
	}
	return path, nil
}

func employeeLabel(doc ports.SalarySlipDocument) string {
	if doc.Employee == nil {
		return doc.Slip.UserID.String()
	}
	name := strings.TrimSpace(strings.Join([]string{doc.Employee.Firstname, ptrString(doc.Employee.MiddleName), ptrString(doc.Employee.Lastname)}, " "))
	if doc.Employee.EmployeeCode != nil && *doc.Employee.EmployeeCode != "" {
		return fmt.Sprintf("%s (%s)", name, *doc.Employee.EmployeeCode)
	}
	return name
}

func ptrString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func money(value float64) string {
	return fmt.Sprintf("%.2f", value)
}

func hexColor(value string, fallbackR int, fallbackG int, fallbackB int) (int, int, int) {
	value = strings.TrimPrefix(strings.TrimSpace(value), "#")
	if len(value) != 6 {
		return fallbackR, fallbackG, fallbackB
	}
	var r, g, b int
	if _, err := fmt.Sscanf(value, "%02x%02x%02x", &r, &g, &b); err != nil {
		return fallbackR, fallbackG, fallbackB
	}
	return r, g, b
}

func safeFileName(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		value = "salary-slip.pdf"
	}
	replacer := strings.NewReplacer("/", "-", "\\", "-", ":", "-", " ", "-")
	return replacer.Replace(value)
}
