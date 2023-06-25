package ptocker

import "time"

type User struct {
	ID            int
	Name          string
	Email         string
	Token         string
	Leaves        []*Leave
	Started       time.Time
	ExtraVacation int
}

type UserService interface {
	AllUsers() ([]*User, error)
}
