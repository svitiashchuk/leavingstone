package leavingstone

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
	Create(userID int, from, to time.Time, leaveType string) error
	Approve(id int) error
	Reject(id int) error
	Delete(id int) error
}
