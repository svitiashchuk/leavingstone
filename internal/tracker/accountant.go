package tracker

import (
	"leavingstone/internal/model"
	"time"
)

type Accountant struct {
	us     model.UserService
	policy model.CompanyPolicy
}

func NewAccountant(us model.UserService) *Accountant {
	return &Accountant{
		us: us,
		policy: model.CompanyPolicy{
			MaxVacationDays: 20,
			MaxSickDays:     10,
		},
	}
}

func (a *Accountant) MaxVacationDays() int {
	return a.policy.MaxVacationDays
}

func (a *Accountant) VacationsLeft(u *model.User) int {
	return a.policy.MaxVacationDays + u.ExtraVacation - a.VacationsUsed(u)
}

func (a *Accountant) VacationsUsed(u *model.User) int {
	yearStart := time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.UTC)
	yearEnd := time.Date(time.Now().Year(), 12, 31, 23, 59, 59, 0, time.UTC)

	return a.us.LeavesUsed(u, []string{"vacation", "dayoff"}, &yearStart, &yearEnd)
}

func (a *Accountant) SickdaysLeft(u *model.User) int {
	return a.policy.MaxSickDays - a.SickdaysUsed(u)
}

func (a *Accountant) SickdaysUsed(u *model.User) int {
	yearStart := time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.UTC)
	yearEnd := time.Date(time.Now().Year(), 12, 31, 23, 59, 59, 0, time.UTC)

	return a.us.LeavesUsed(u, []string{"sick"}, &yearStart, &yearEnd)
}

func (a *Accountant) MaxSickDays() int {
	return a.policy.MaxSickDays
}
