package http

import (
	"net/http"
	"os"
	"ptocker/internal/pkg/tracker"
	"ptocker/sqlite"
	"strconv"
	"text/template"
	"time"
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

type Server struct {
	t *tracker.Tracker
}

func NewServer() *Server {
	us, err := sqlite.NewUserService()
	if err != nil {
		panic(err)
	}

	ls, err := sqlite.NewLeaveService()
	if err != nil {
		panic(err)
	}

	return &Server{
		t: tracker.NewTracker(us, ls),
	}
}

func (s *Server) Serve() {
	s.registerRoutes()

	http.ListenAndServe(":8887", nil)
}

func (s *Server) registerRoutes() {
	http.HandleFunc("/", s.handleIndex)
	http.HandleFunc("/login", s.handleLogin)
	http.HandleFunc("/tracker", s.handleTracker)
	http.HandleFunc("/overview", s.handleOverview)

	// assets for frontend
	http.HandleFunc("/dist/", s.handleDist)

	http.HandleFunc("/leaves/approve", s.handleLeaveApprove)
	http.HandleFunc("/leaves/reject", s.handleLeaveReject)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		"frontend/src/templates/layout.html",
		"frontend/src/templates/list.html",
	))
	tmpl.ExecuteTemplate(w, "layout", nil)
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		"frontend/src/templates/layout.html",
		"frontend/src/templates/login.html",
	))

	if r.Method == "POST" {
		r.ParseForm()

		// TODO - retrieve user from DB and set token in cookie

		w.Header().Add("HX-Redirect", "/login")
		w.Header().Add("Set-Cookie", "auth_token=TOKEN; Path=/; HttpOnly")
		tmpl.ExecuteTemplate(w, "layout", nil)
	} else {

		tmpl.ExecuteTemplate(w, "layout", nil)
	}
}

func (s *Server) handleOverview(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		"frontend/src/templates/layout.html",
		"frontend/src/templates/overview.html",
	))

	tmpl.ExecuteTemplate(w, "layout", nil)
}

func (s *Server) handleTracker(w http.ResponseWriter, r *http.Request) {
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

	data := map[string]interface{}{
		"Nav":           nav,
		"Users":         ee,
		"Days":          days,
		"WorkforceStat": workforceStat,
	}

	tmpl.ExecuteTemplate(w, "tracker.html", data)
}

func (s *Server) handleDist(w http.ResponseWriter, r *http.Request) {
	// TODO use embed:
	b, err := os.ReadFile("frontend/dist/output.css")
	if err != nil {
		panic(err)
	}

	w.Header().Add("Content-type", "text/css")
	w.Write(b)
}

func (s *Server) handleLeaveApprove(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	id, err := strconv.Atoi(r.Form.Get("id"))
	if err != nil {
		panic(err)
	}

	s.t.ApproveLeave(id)
	s.handleTracker(w, r)
}

func (s *Server) handleLeaveReject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	id, err := strconv.Atoi(r.Form.Get("id"))
	if err != nil {
		panic(err)
	}

	s.t.RejectLeave(id)
	s.handleTracker(w, r)
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
