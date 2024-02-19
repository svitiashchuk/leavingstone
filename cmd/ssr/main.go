package main

import (
	"fmt"
	"leavingstone/internal/auth"
	"leavingstone/internal/session"
	"leavingstone/internal/sqlite"
	"leavingstone/internal/tracker"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	errLogger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	appLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

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

	app := tracker.NewApp(session.NewKeeper(), auth, us, ls, t, ac, appLogger, errLogger)

	app.RegisterRoutes()

	addr := ":8080"
	fmt.Printf("Serving... at %s\n", addr)
	http.ListenAndServe(addr, nil)
}
