package ptocker

import "time"

type Leave struct {
	ID       int
	Type     string
	Start    time.Time
	End      time.Time
	Approved bool
	UserID   int
}

type LeaveService interface {
	List(from, to time.Time, limit int) ([]*Leave, error)
}
