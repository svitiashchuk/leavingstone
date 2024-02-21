package main

import (
	"database/sql"
	"fmt"
	"leavingstone/internal/auth"
	"leavingstone/internal/session"
	"leavingstone/internal/sqlite"
	"leavingstone/internal/tracker"
	"log/slog"
	"net/http"
	"os"
)

const DSN = "file:database.db?cache=shared&mode=rwc"

func main() {
	errLogger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	appLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	db, err := sql.Open("sqlite3", DSN)
	if err != nil {
		errLogger.Error(err.Error())
		os.Exit(1)
	}

	us := sqlite.NewUserService(db)
	ls := sqlite.NewLeaveService(db)
	teamService := sqlite.NewTeamService(db)
	ac := tracker.NewAccountant(us)
	auth := auth.New(us)
	t := tracker.NewTracker(us, ls)

	app := tracker.NewApp(session.NewKeeper(), auth, us, ls, teamService, t, ac, appLogger, errLogger)

	app.RegisterRoutes()

	addr := ":8080"
	fmt.Printf("Serving... at %s\n", addr)
	http.ListenAndServe(addr, nil)
}
