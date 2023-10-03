package main

import (
	"fmt"
	"leavingstone/internal/pkg/tracker"
	"leavingstone/sqlite"
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
