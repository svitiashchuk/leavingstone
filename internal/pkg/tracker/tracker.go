package tracker

import (
	"ptocker"
	"sort"
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
	ls ptocker.LeaveService
}

func NewTracker(us ptocker.UserService, ls ptocker.LeaveService) *Tracker {
	return &Tracker{us, ls}
}

type Leave struct {
	ID       int
	Type     string
	Start    time.Time
	End      time.Time
	Approved bool
	Employee *EmployeeSummary
}

type EmployeeSummary struct {
	ID   int
	Name string
}

// todo: consider adding raw date as time.Time to allow different formatting if needed
type LeaveDay struct {
	LeaveID    int
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

	sort.Slice(ee, func(i, j int) bool { return ee[i].Name > ee[j].Name })

	return ee
}

func (t *Tracker) UpcomingLeaves() []*Leave {
	from := time.Now()
	to := from.Add(time.Duration(30 * 24 * time.Hour))
	ll, err := t.ls.List(from, to, 50)
	if err != nil {
		panic(err)
	}

	ull := []*Leave{}
	for _, l := range ll {
		ull = append(ull, &Leave{
			l.ID,
			l.Type,
			l.Start,
			l.End,
			l.Approved,
			nil,
		})
	}

	return ull
}

func (t *Tracker) ApproveLeave(id int) {
	t.ls.Approve(id)
}

func (t *Tracker) RejectLeave(id int) {
	t.ls.Delete(id)
}

// TODO calculate workforce power based on period
func (t *Tracker) WorkforcePower(ee []*Employee) int {
	return 87
}

// leavesToCalendar converts from slice of Leaves that contains
// information about start and end of Leave to map where key-string
// is a date in Calendar map and value is LeaveDay which stores information
// about particular type of Leave and whether it was approved.
func leavesToCalendar(ll []*ptocker.Leave) *Calendar {
	c := make(Calendar)

	for _, l := range ll {
		for d := l.Start; d.Before(l.End); d = d.Add(24 * time.Hour) {
			c[d.Format("2006-01-02")] = LeaveDay{l.ID, l.Type, l.Approved}
		}
	}

	return &c
}

func (c Calendar) Get(day time.Time) LeaveDay {
	if isWeekend(day) {
		return LeaveDay{0, Weekend, true}
	}

	if isBankHoliday(day) {
		return LeaveDay{0, BankHoliday, true}
	}

	if ld, ok := c[day.Format("2006-01-02")]; ok {
		return ld
	}

	return LeaveDay{0, Workday, true}
}

func isWeekend(day time.Time) bool {
	return day.Weekday() == 0 || day.Weekday() == 6
}

func isBankHoliday(day time.Time) bool {
	return false
}
