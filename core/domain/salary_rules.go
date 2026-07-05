package domain

type SalaryItem struct {
	ItemType  string  `json:"item_type"`
	Code      string  `json:"code"`
	Name      string  `json:"name"`
	Amount    float64 `json:"amount"`
	SortOrder int     `json:"sort_order"`
}

type SalaryResult struct {
	TotalEarnings   float64      `json:"total_earnings"`
	TotalDeductions float64      `json:"total_deductions"`
	AbsentDeduction float64      `json:"absent_deduction"`
	NetSalary       float64      `json:"net_salary"`
	Items           []SalaryItem `json:"items,omitempty"`
}

func CalculateUserSalary(grossSalary float64, structure []SalaryItem, presentDays, absentDays, daysInMonth int, isSpecial bool) SalaryResult {
	var earnings float64
	var deductions float64
	for _, item := range structure {
		if item.ItemType == SalaryItemEarning {
			earnings += item.Amount
			continue
		}
		if item.ItemType == SalaryItemDeduction {
			deductions += item.Amount
		}
	}
	var absentDeduction float64
	if absentDays > 0 && daysInMonth > 0 {
		absentDeduction = (grossSalary / float64(daysInMonth)) * float64(absentDays)
		deductions += absentDeduction
	}
	if isSpecial && presentDays > 0 && daysInMonth > 0 {
		earnings = (earnings / float64(daysInMonth)) * float64(presentDays)
	}
	result := SalaryResult{TotalEarnings: roundMoney(earnings), TotalDeductions: roundMoney(deductions), AbsentDeduction: roundMoney(absentDeduction), NetSalary: roundMoney(earnings - deductions)}
	result.Items = append(result.Items, structure...)
	if result.AbsentDeduction > 0 {
		result.Items = append(result.Items, SalaryItem{ItemType: SalaryItemDeduction, Code: SalaryCodeLWP, Name: "Leave Without Pay", Amount: result.AbsentDeduction, SortOrder: 999})
	}
	return result
}
