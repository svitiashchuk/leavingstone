package main

import (
	"fmt"
	"leavingstone"
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

type CommonTemplateData struct {
	IsAuthenticated bool
}

type ProfileTemplateData struct {
	CommonTemplateData
	User *leavingstone.User
}

type TrackerTemplateData struct {
	CommonTemplateData
	Nav           Navigation
	Employees     []*tracker.Employee
	Days          []time.Time
	WorkforceStat *tracker.WorkforceStat
	LeavesStat    *tracker.LeavesStat
}

func (app *App) registerRoutes() {
	http.HandleFunc("/", app.authenticate(app.requireAuth(app.handleIndex)))
	http.HandleFunc("/login", app.handleLogin)
	http.HandleFunc("/profile", app.authenticate(app.requireAuth((app.handleProfile))))
	http.HandleFunc("/tracker", app.authenticate(app.requireAuth((app.handleTracker))))
	http.HandleFunc("/overview", app.authenticate(app.requireAuth(app.handleOverview)))

	// assets for frontend
	http.HandleFunc("/dist/", app.handleDist)

	http.HandleFunc("/leaves/approve", app.authenticate(app.requireAuth(app.handleLeaveApprove)))
	http.HandleFunc("/leaves/reject", app.authenticate(app.requireAuth(app.handleLeaveReject)))
}

func (app *App) handleIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		"frontend/src/templates/layout.html",
		"frontend/src/templates/list.html",
	))

	tmpl.ExecuteTemplate(w, "layout", app.commonTemplateData(r))
}

func (app *App) handleLogin(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		"frontend/src/templates/layout.html",
		"frontend/src/templates/login.html",
	))

	if r.Method == "POST" {
		r.ParseForm()
		email := r.Form.Get("email")
		passPlain := r.Form.Get("password")

		u, err := app.us.Find(email)
		if err != nil {
			panic(err)
		}

		err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(passPlain))
		if err != nil {
			panic(err)
		}

		c := fmt.Sprintf("auth_token=%s; Path=/; HttpOnly", u.Token)
		app.htmxRedirect(w, r, "/profile")
		w.Header().Add("Set-Cookie", c)
	}

	tmpl.ExecuteTemplate(w, "layout", app.commonTemplateData(r))
}

func (app *App) handleProfile(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		"frontend/src/templates/layout.html",
		"frontend/src/templates/profile.html",
	))

	u, err := app.us.FindByID(app.userID(r))
	if err != nil {
		panic(err)
	}

	templateData := &ProfileTemplateData{
		CommonTemplateData: *app.commonTemplateData(r),
		User:               u,
	}

	tmpl.ExecuteTemplate(w, "layout", templateData)
}

func (app *App) handleOverview(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		"frontend/src/templates/layout.html",
		"frontend/src/templates/overview.html",
	))

	tmpl.ExecuteTemplate(w, "layout", app.commonTemplateData(r))
}

func (app *App) handleTracker(w http.ResponseWriter, r *http.Request) {
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
	ee := app.t.List()

	next := time.Date(y, time.Month(m+1), 1, 0, 0, 0, 0, time.UTC)
	prev := time.Date(y, time.Month(m-1), 1, 0, 0, 0, 0, time.UTC)

	nav := Navigation{
		Prev: MonthPeriod{Month: prev.Month(), Year: prev.Year()},
		Now:  MonthPeriod{Month: time.Month(m), Year: y},
		Next: MonthPeriod{Month: next.Month(), Year: next.Year()},
	}

	data := &TrackerTemplateData{
		CommonTemplateData: *app.commonTemplateData(r),
		Nav:                nav,
		Employees:          ee,
		Days:               days,
		WorkforceStat:      app.t.WorkforceStat(days, ee),
		LeavesStat:         app.t.LeavesStat(days, ee),
	}

	tmpl.ExecuteTemplate(w, "tracker.html", data)
}

func (app *App) handleDist(w http.ResponseWriter, r *http.Request) {
	// TODO use embed:
	b, err := os.ReadFile("frontend/dist/output.css")
	if err != nil {
		panic(err)
	}

	w.Header().Add("Content-type", "text/css")
	w.Write(b)
}

func (app *App) handleLeaveApprove(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	id, err := strconv.Atoi(r.Form.Get("id"))
	if err != nil {
		panic(err)
	}

	app.t.ApproveLeave(id)

	// send hx-trigger header to reload full tracker
	w.Header().Add("HX-Trigger", "reloadTracker")
}

func (app *App) handleLeaveReject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	id, err := strconv.Atoi(r.Form.Get("id"))
	if err != nil {
		panic(err)
	}

	app.t.RejectLeave(id)
	// send hx-trigger header to reload full tracker
	w.Header().Add("HX-Trigger", "reloadTracker")
}

func (app *App) htmxRedirect(w http.ResponseWriter, r *http.Request, url string) {
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

func (app *App) commonTemplateData(r *http.Request) *CommonTemplateData {
	return &CommonTemplateData{
		IsAuthenticated: app.isAuthenticated(r),
	}
}
