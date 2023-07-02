package tracker

import (
	"ptocker"
	"time"
)

// Day types
const (
	Vacation    = "vacation"
	SickDay     = "sickday"
	DayOff      = "dayoff"
	Workday     = "workday"
	Weekend     = "weekend"
	BankHoliday = "bankholiday"
)

type Tracker struct {
	us ptocker.UserService
}

func NewTracker(us ptocker.UserService) *Tracker {
	return &Tracker{us}
}

// todo: consider adding raw date as time.Time to allow different formatting if needed
type LeaveDay struct {
	Type       string
	IsApproved bool
}

type Calendar map[string]LeaveDay

type Employee struct {
	Name     string
	Calendar *Calendar
}

func (t *Tracker) List() []*Employee {
	uu, err := t.us.AllUsers()
	if err != nil {
		panic(err)
	}

	ee := make([]*Employee, len(uu))
	for i, u := range uu {
		ee[i] = &Employee{u.Name, leavesToCalendar(u.Leaves)}
	}

	return ee
}

// leavesToCalendar converts from slice of Leaves that contains
// information about start and end of Leave to map where key-string
// is a date in Calendar map and value is LeaveDay which stores information
// about particular type of Leave and whether it was approved.
func leavesToCalendar(ll []*ptocker.Leave) *Calendar {
	c := make(Calendar)

	for _, l := range ll {
		for d := l.Start; d.Before(l.End); d = d.Add(24 * time.Hour) {
			c[d.Format("2006-01-02")] = LeaveDay{l.Type, l.Approved}
		}
	}

	return &c
}

func (c Calendar) Get(day time.Time) LeaveDay {
	if isWeekend(day) {
		return LeaveDay{Weekend, true}
	}

	if isBankHoliday(day) {
		return LeaveDay{BankHoliday, true}
	}

	if ld, ok := c[day.Format("2006-01-02")]; ok {
		return ld
	}

	return LeaveDay{Workday, true}
}

func isWeekend(day time.Time) bool {
	return day.Weekday() == 0 || day.Weekday() == 6
}

func isBankHoliday(day time.Time) bool {
	return false
}
