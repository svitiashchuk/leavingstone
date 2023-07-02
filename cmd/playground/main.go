package main

import (
	"fmt"
	"ptocker/internal/pkg/tracker"
	"ptocker/sqlite"
)

func main() {
	us, _ := sqlite.NewUserService()
	t := tracker.NewTracker(us)

	ee := t.List()

	for _, e := range ee {
		fmt.Printf("%#\n", e)
	}
}
