package main

import (
	"fmt"
	"leavingstone/internal/auth"
	"leavingstone/internal/session"
	"leavingstone/internal/sqlite"
	"leavingstone/internal/tracker"
	"net/http"
)

func main() {
	us, err := sqlite.NewUserService()
	if err != nil {
		panic(err)
	}

	ls, err := sqlite.NewLeaveService()
	if err != nil {
		panic(err)
	}

	ac := tracker.NewAccountant(us)
	auth := auth.New(us)
	t := tracker.NewTracker(us, ls)

	app := tracker.NewApp(session.NewKeeper(), auth, us, ls, t, ac)

	app.RegisterRoutes()

	addr := ":8080"
	fmt.Printf("Serving... at %s\n", addr)
	http.ListenAndServe(addr, nil)
}
