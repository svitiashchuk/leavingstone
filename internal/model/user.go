package model

import "time"

type User struct {
	ID            int
	Name          string
	Email         string
	Token         string
	Password      string
	Leaves        []*Leave
	Started       time.Time
	TeamID        int
	ExtraVacation int
}

type MemberInfo struct {
	ID            int
	Name          string
	Email         string
	Started       time.Time
	ExtraVacation int
	VacationsUsed int
	SickdaysUsed  int
	TodayStatus   string
}

type UserService interface {
	AllUsers() ([]*User, error)
	LeavesUsed(u *User, leaveTypes []string, periodStart, periodEnd *time.Time) int
}
