package main

import (
	"time"
)

type MonthPeriod struct {
	Month time.Month
	Year  int
}

func (mp MonthPeriod) MonthNum() int {
	return int(mp.Month)
}

func calendarMonth(year, month int) [][]time.Time {
	y := year
	m := time.Month(month)

	// TODO problem with September 2024, where first day is Sunday
	firstDayOfMonth := time.Date(y, m, 1, 0, 0, 0, 0, time.UTC)
	diffDaysToWeekStart := firstDayOfMonth.Weekday() - time.Monday
	if (diffDaysToWeekStart) < 0 {
		diffDaysToWeekStart = 7 + diffDaysToWeekStart
	}
	firstDayForCalendar := time.Date(y, m, 1-int(diffDaysToWeekStart), 0, 0, 0, 0, time.UTC)

	lastDayOfMonth := time.Date(y, m+1, 0, 0, 0, 0, 0, time.UTC)
	diffDaysToWeekEnd := 7 - lastDayOfMonth.Weekday()
	lastDayForCalendar := time.Date(y, m+1, int(diffDaysToWeekEnd), 0, 0, 0, 0, time.UTC)

	// fmt.Print("Mo\tTu\tWe\tTh\tFr\tSa\tSu\n")

	monthWeeks := [][]time.Time{}
	monthWeeks = append(monthWeeks, []time.Time{})

	currentWeek := 0
	d := firstDayForCalendar
	for i := 0; !d.After(lastDayForCalendar); i += 1 {
		monthWeeks[currentWeek] = append(monthWeeks[currentWeek], d)
		// fmt.Printf("%d\t", d.Day())
		if d.Weekday() == time.Sunday {
			currentWeek += 1
			monthWeeks = append(monthWeeks, []time.Time{})
			//fmt.Print("\n")
		}

		d = d.AddDate(0, 0, 1)
	}

	return monthWeeks
}

func month(year, month int) []time.Time {
	days := []time.Time{}
	initial := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)

	for i := 0; i < daysInMonth(year, month); i++ {
		days = append(days, initial.AddDate(0, 0, i))
	}

	return days
}

func daysInMonth(year, month int) int {
	d := time.Date(year, time.Month(month+1), 0, 0, 0, 0, 0, time.UTC)

	return d.Day()
}

func days(year int) []time.Time {
	days := []time.Time{}
	initial := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	for i := 0; i < daysInYear(year); i++ {
		days = append(days, initial.AddDate(0, 0, i))
	}

	return days
}

func daysInYear(year int) int {
	if isLeap(year) {
		return 366
	}

	return 365
}

func isLeap(year int) bool {
	return year%4 == 0 && year%100 != 0 || year%400 == 0
}

func period(start, end time.Time) []time.Time {
	days := []time.Time{}
	for d := start; d.Before(end); d = d.AddDate(0, 0, 1) {
		days = append(days, d)
	}

	return days
}
