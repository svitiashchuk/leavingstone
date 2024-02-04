package leavingstone

import (
	"math"
	"time"
)

type Leave struct {
	ID       int
	Type     string
	Start    time.Time
	End      time.Time
	Approved bool
	UserID   int
	User     *User
}

func (l *Leave) Duration() time.Duration {
	return l.End.Sub(l.Start)
}

func (l *Leave) DurationDays() int {
	return int(math.Round(l.Duration().Hours() / 24))
}

type LeaveService interface {
	List(from, to time.Time, limit int) ([]*Leave, error)
	Create(userID int, from, to time.Time, leaveType string) error
	Approve(id int) error
	Reject(id int) error
	Delete(id int) error
}
