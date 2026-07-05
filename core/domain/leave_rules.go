package domain

import "time"

func LeaveDayUnits(dayType string) (float64, error) {
	dayType, err := ValidateLeaveDayType(dayType)
	if err != nil {
		return 0, err
	}
	if dayType == LeaveDayFullDay {
		return 1, nil
	}
	return 0.5, nil
}

func CalculateLeaveDays(startDate, endDate time.Time, startDayType, endDayType string) (float64, error) {
	startDate = dateOnly(startDate)
	endDate = dateOnly(endDate)
	if endDate.Before(startDate) {
		return 0, ErrInvalidDateRange
	}
	normalizedStartDayType, err := ValidateLeaveDayType(startDayType)
	if err != nil {
		return 0, err
	}
	normalizedEndDayType, err := ValidateLeaveDayType(endDayType)
	if err != nil {
		return 0, err
	}
	startUnits, err := LeaveDayUnits(startDayType)
	if err != nil {
		return 0, err
	}
	endUnits, err := LeaveDayUnits(endDayType)
	if err != nil {
		return 0, err
	}
	if startDate.Equal(endDate) {
		if normalizedStartDayType == LeaveDayFullDay || normalizedEndDayType == LeaveDayFullDay {
			return 1, nil
		}
		if normalizedStartDayType == normalizedEndDayType {
			return 0.5, nil
		}
		return 1, nil
	}
	middleDays := endDate.Sub(startDate).Hours()/24 - 1
	return startUnits + endUnits + middleDays, nil
}

func IsSandwich(startDate, endDate time.Time, holidays []time.Time, weekoffs []time.Weekday) (bool, float64) {
	startDate = dateOnly(startDate)
	endDate = dateOnly(endDate)
	if !startDate.Before(endDate) {
		return false, 0
	}
	holidaySet := make(map[time.Time]struct{}, len(holidays))
	for _, holiday := range holidays {
		holidaySet[dateOnly(holiday)] = struct{}{}
	}
	weekoffSet := make(map[time.Weekday]struct{}, len(weekoffs))
	for _, weekoff := range weekoffs {
		weekoffSet[weekoff] = struct{}{}
	}
	var gapDays float64
	for day := startDate.AddDate(0, 0, 1); day.Before(endDate); day = day.AddDate(0, 0, 1) {
		_, holiday := holidaySet[day]
		_, weekoff := weekoffSet[day.Weekday()]
		if !holiday && !weekoff {
			return false, 0
		}
		gapDays++
	}
	return gapDays > 0, gapDays
}

func DistributeInto12Months(totalDays float64) [12]int {
	wholeDays := int(totalDays)
	base := wholeDays / 12
	remainder := wholeDays % 12
	var months [12]int
	for i := range months {
		months[i] = base
	}
	for i := 0; i < remainder; i++ {
		months[i]++
	}
	return months
}
