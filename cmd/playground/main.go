package main

import (
	"fmt"
	"ptocker/internal/pkg/tracker"
	"ptocker/sqlite"
)

func main() {
	us, _ := sqlite.NewUserService()
	ls, _ := sqlite.NewLeaveService()
	t := tracker.NewTracker(us, ls)

	ee := t.List()

	for _, e := range ee {
		fmt.Printf("%#v\n", e)
	}
}
