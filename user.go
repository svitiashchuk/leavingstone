package leavingstone

import "time"

type User struct {
	ID            int
	Name          string
	Email         string
	Token         string
	Password      string
	Leaves        []*Leave
	Started       time.Time
	ExtraVacation int
}

type UserService interface {
	AllUsers() ([]*User, error)
	LeavesUsed(u *User, leaveTypes []string, periodStart, periodEnd *time.Time) int
}
