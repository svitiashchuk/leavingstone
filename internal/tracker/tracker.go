package tracker

import (
	"leavingstone/internal/model"
	"math"
	"sort"
	"time"
)

// Day types
const (
	Vacation    = "vacation"
	SickDay     = "sick"
	DayOff      = "dayoff"
	Workday     = "workday"
	Weekend     = "weekend"
	BankHoliday = "bankholiday"
)

type Tracker struct {
	us model.UserService
	ls model.LeaveService
}

func NewTracker(us model.UserService, ls model.LeaveService) *Tracker {
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
	Day        time.Time
}

type Calendar map[string]LeaveDay

type Employee struct {
	Name     string
	Calendar *Calendar
}

type WorkforceStat struct {
	AbsentEmployees int
	WorkforcePower  int
}

type LeavesStat struct {
	Pending   int
	AllLeaves int
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

func (t *Tracker) WorkforceStat(period []time.Time, ee []*Employee) *WorkforceStat {
	absentDays := 0
	absentEmployees := make(map[int]interface{})

	firstDay := period[0]
	lastDay := period[len(period)-1]
	periodDays := len(period) * len(ee)

	for k, e := range ee {
		for _, ld := range *e.Calendar {
			if ld.isLeave() && ld.IsApproved && ld.Day.After(firstDay) && ld.Day.Before(lastDay) {
				absentDays++
				absentEmployees[k] = nil
			}
		}
	}

	return &WorkforceStat{
		AbsentEmployees: len(absentEmployees),
		WorkforcePower:  int(math.Round(100.0 * (float64(periodDays-absentDays) / float64(periodDays)))),
	}
}

func (t *Tracker) LeavesStat(period []time.Time, ee []*Employee) *LeavesStat {
	allLeaves := make(map[int]interface{})
	pendingLeaves := make(map[int]interface{})

	firstDay := period[0]
	lastDay := period[len(period)-1]

	for _, e := range ee {
		for _, ld := range *e.Calendar {
			if !ld.isLeave() || ld.Day.Before(firstDay) || ld.Day.After(lastDay) {
				continue
			}

			if !ld.IsApproved {
				pendingLeaves[ld.LeaveID] = nil
			}

			if ld.isLeave() {
				allLeaves[ld.LeaveID] = nil
			}
		}
	}

	return &LeavesStat{
		Pending:   len(pendingLeaves),
		AllLeaves: len(allLeaves),
	}
}

// leavesToCalendar converts from slice of Leaves that contains
// information about start and end of Leave to map where key-string
// is a date in Calendar map and value is LeaveDay which stores information
// about particular type of Leave and whether it was approved.
func leavesToCalendar(ll []*model.Leave) *Calendar {
	c := make(Calendar)

	for _, l := range ll {
		for d := l.Start; d.Before(l.End); d = d.Add(24 * time.Hour) {
			c[d.Format("2006-01-02")] = LeaveDay{l.ID, l.Type, l.Approved, d}
		}
	}

	return &c
}

func (c Calendar) Get(day time.Time) LeaveDay {
	if isWeekend(day) {
		return LeaveDay{0, Weekend, true, day}
	}

	if isBankHoliday(day) {
		return LeaveDay{0, BankHoliday, true, day}
	}

	if ld, ok := c[day.Format("2006-01-02")]; ok {
		return ld
	}

	return LeaveDay{0, Workday, true, day}
}

func isWeekend(day time.Time) bool {
	return day.Weekday() == 0 || day.Weekday() == 6
}

func isBankHoliday(day time.Time) bool {
	return false
}

func (ld *LeaveDay) isLeave() bool {
	return ld.Type == Vacation || ld.Type == SickDay || ld.Type == DayOff
}

func DayTypes() []string {
	return []string{Vacation, SickDay, DayOff, Workday, Weekend, BankHoliday}
}

func LeaveTypes() []string {
	return []string{Vacation, SickDay, DayOff}
}
