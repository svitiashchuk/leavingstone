package main

import (
	"fmt"
	"leavingstone/internal/pkg/tracker"
	"leavingstone/sqlite"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type MonthPeriod struct {
	Month time.Month
	Year  int
}

type Navigation struct {
	Prev MonthPeriod
	Now  MonthPeriod
	Next MonthPeriod
}

type App struct {
	sm   SessionManager
	auth *Authenticator
	us   *sqlite.UserService
	t    *tracker.Tracker
}

func (app *App) registerRoutes() {
	http.HandleFunc("/", app.requireAuth(app.handleIndex))
	http.HandleFunc("/login", app.handleLogin)
	http.HandleFunc("/profile", app.requireAuth(app.handleProfile))
	http.HandleFunc("/tracker", app.requireAuth(app.handleTracker))
	http.HandleFunc("/overview", app.requireAuth(app.handleOverview))

	// assets for frontend
	http.HandleFunc("/dist/", app.handleDist)

	http.HandleFunc("/leaves/approve", app.requireAuth(app.handleLeaveApprove))
	http.HandleFunc("/leaves/reject", app.requireAuth(app.handleLeaveReject))
}

func (s *App) handleIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		"frontend/src/templates/layout.html",
		"frontend/src/templates/list.html",
	))
	tmpl.ExecuteTemplate(w, "layout", nil)
}

func (s *App) handleLogin(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		"frontend/src/templates/layout.html",
		"frontend/src/templates/login.html",
	))

	if r.Method == "POST" {
		r.ParseForm()
		email := r.Form.Get("email")
		passPlain := r.Form.Get("password")

		u, err := s.us.Find(email)
		if err != nil {
			panic(err)
		}

		err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(passPlain))
		if err != nil {
			panic(err)
		}

		c := fmt.Sprintf("auth_token=%s; Path=/; HttpOnly", u.Token)
		s.htmxRedirect(w, r, "/profile")
		w.Header().Add("Set-Cookie", c)
	}

	tmpl.ExecuteTemplate(w, "layout", nil)
}

func (s *App) handleProfile(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		"frontend/src/templates/layout.html",
		"frontend/src/templates/profile.html",
	))

	tmpl.ExecuteTemplate(w, "layout", nil)
}

func (s *App) handleOverview(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		"frontend/src/templates/layout.html",
		"frontend/src/templates/overview.html",
	))

	tmpl.ExecuteTemplate(w, "layout", nil)
}

func (s *App) handleTracker(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		"frontend/src/templates/fragments/tracker.html",
	))

	y, err := strconv.Atoi(r.URL.Query().Get("year"))
	if err != nil {
		y = time.Now().Year()
	}

	m, err := strconv.Atoi(r.URL.Query().Get("month"))
	if err != nil {
		m = int(time.Now().Month())
	}

	days := month(y, m)
	ee := s.t.List()

	next := time.Date(y, time.Month(m+1), 1, 0, 0, 0, 0, time.UTC)
	prev := time.Date(y, time.Month(m-1), 1, 0, 0, 0, 0, time.UTC)

	nav := Navigation{
		Prev: MonthPeriod{Month: prev.Month(), Year: prev.Year()},
		Now:  MonthPeriod{Month: time.Month(m), Year: y},
		Next: MonthPeriod{Month: next.Month(), Year: next.Year()},
	}

	workforceStat := s.t.WorkforceStat(days, ee)
	leavesStat := s.t.LeavesStat(days, ee)

	data := map[string]interface{}{
		"Nav":           nav,
		"Users":         ee,
		"Days":          days,
		"WorkforceStat": workforceStat,
		"LeavesStat":    leavesStat,
	}

	tmpl.ExecuteTemplate(w, "tracker.html", data)
}

func (s *App) handleDist(w http.ResponseWriter, r *http.Request) {
	// TODO use embed:
	b, err := os.ReadFile("frontend/dist/output.css")
	if err != nil {
		panic(err)
	}

	w.Header().Add("Content-type", "text/css")
	w.Write(b)
}

func (s *App) handleLeaveApprove(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	id, err := strconv.Atoi(r.Form.Get("id"))
	if err != nil {
		panic(err)
	}

	s.t.ApproveLeave(id)

	// send hx-trigger header to reload full tracker
	w.Header().Add("HX-Trigger", "reloadTracker")
}

func (s *App) handleLeaveReject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	id, err := strconv.Atoi(r.Form.Get("id"))
	if err != nil {
		panic(err)
	}

	s.t.RejectLeave(id)
	// send hx-trigger header to reload full tracker
	w.Header().Add("HX-Trigger", "reloadTracker")
}

func (s *App) htmxRedirect(w http.ResponseWriter, r *http.Request, url string) {
	w.Header().Add("HX-Redirect", url)
}

func (mp MonthPeriod) MonthNum() int {
	return int(mp.Month)
}

func month(year, month int) []time.Time {
	days := []time.Time{}
	initial := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)

	for i := 0; i < daysInMonth(year, month); i++ {
		days = append(days, initial.AddDate(0, 0, i))
	}

	return days
}

func daysInMonth(year, month int) int {
	d := time.Date(year, time.Month(month+1), 0, 0, 0, 0, 0, time.UTC)

	return d.Day()
}

func period(start, end time.Time) []time.Time {
	days := []time.Time{}
	for d := start; d.Before(end); d = d.AddDate(0, 0, 1) {
		days = append(days, d)
	}

	return days
}

func days(year int) []time.Time {
	days := []time.Time{}
	initial := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	for i := 0; i < daysInYear(year); i++ {
		days = append(days, initial.AddDate(0, 0, i))
	}

	return days
}

func daysInYear(year int) int {
	if isLeap(year) {
		return 366
	}

	return 365
}

func isLeap(year int) bool {
	return year%4 == 0 && year%100 != 0 || year%400 == 0
}