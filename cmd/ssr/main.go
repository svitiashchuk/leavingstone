package main

import (
	"fmt"
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

	auth := &Authenticator{
		us: us,
	}

	t := tracker.NewTracker(us, ls)

	app := &App{
		sm:   NewSessionKeeper(),
		auth: auth,
		us:   us,
		ls:   ls,
		t:    t,
		ac:   ac,
	}

	app.registerRoutes()

	addr := ":8080"
	fmt.Printf("Serving... at %s\n", addr)
	http.ListenAndServe(addr, nil)
}
