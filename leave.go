package ptocker

import "time"

type Leave struct {
	ID       int
	Type     string
	Start    time.Time
	End      time.Time
	Approved bool
}
