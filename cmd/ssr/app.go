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

type Navigation struct {
	Prev MonthPeriod
	Now  MonthPeriod
	Next MonthPeriod
}

type App struct {
	sm   SessionManager
	auth *Authenticator
	us   *sqlite.UserService
	ls   *sqlite.LeaveService
	t    *tracker.Tracker
	ac   *tracker.Accountant
}

type CommonTemplateData struct {
	IsAuthenticated bool
}

type CommonFormTemplateData struct {
	Errors []string
}

type ProfileTemplateData struct {
	CommonTemplateData
	UpcomingLeaves []*leavingstone.Leave
	User           *leavingstone.User
	VacationsMax   int
	VacationsUsed  int
	VacationsLeft  int
	SickdaysMax    int
	SickdaysUsed   int
	SickdaysLeft   int
}

type TrackerTemplateData struct {
	CommonTemplateData
	Nav           Navigation
	Employees     []*tracker.Employee
	Days          []time.Time
	WorkforceStat *tracker.WorkforceStat
	LeavesStat    *tracker.LeavesStat
}

type CalendarTemplateData struct {
	CommonTemplateData
	Today         time.Time
	Weekdays      []string
	MonthWeekDays [][]time.Time
	SelectedYear  int
	SelectedMonth time.Month
	Nav           CalendarNav
}

type CalendarNav struct {
	Prev MonthPeriod
	Next MonthPeriod
}

func (app *App) registerRoutes() {
	http.HandleFunc("/login", app.handleLogin)
	http.HandleFunc("/", app.authenticate(app.requireAuth(app.handleIndex)))
	http.HandleFunc("/plan-leave", app.authenticate(app.requireAuth(app.handlePlanLeave)))
	http.HandleFunc("/profile", app.authenticate(app.requireAuth(app.handleProfile)))
	http.HandleFunc("/overview", app.authenticate(app.requireAuth(app.handleOverview)))

	// assets for frontend
	http.HandleFunc("/dist/", app.handleDist)

	http.HandleFunc("/leaves/approve", app.authenticate(app.requireAuth(app.handleLeaveApprove)))
	http.HandleFunc("/leaves/reject", app.authenticate(app.requireAuth(app.handleLeaveReject)))

	// fragments
	http.HandleFunc("/tracker", app.authenticate(app.requireAuth(app.handleTracker)))
	http.HandleFunc("/fragments/calendar", app.authenticate(app.requireAuth(app.handleCalendar)))
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
	tmpl := template.Must(
		template.
			New("profile").
			Funcs(templateFuncs()).
			ParseFiles(
				"frontend/src/templates/layout.html",
				"frontend/src/templates/profile.html",
			),
	)

	u, err := app.us.FindByID(app.userID(r))
	if err != nil {
		panic(err)
	}

	leaves, err := app.ls.Upcoming(app.userID(r))
	if err != nil {
		panic(err)
	}

	templateData := &ProfileTemplateData{
		CommonTemplateData: *app.commonTemplateData(r),
		UpcomingLeaves:     leaves,
		User:               u,
		VacationsMax:       app.ac.MaxVacationDays() + u.ExtraVacation,
		VacationsLeft:      app.ac.VacationsLeft(u),
		VacationsUsed:      app.ac.VacationsUsed(u),
		SickdaysMax:        app.ac.MaxSickDays(),
		SickdaysLeft:       app.ac.SickdaysLeft(u),
		SickdaysUsed:       app.ac.SickdaysUsed(u),
	}

	err = tmpl.ExecuteTemplate(w, "layout", templateData)
	if err != nil {
		panic(err)
	}
}

func (app *App) handlePlanLeave(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		startDate := r.Form.Get("start_date")
		endDate := r.Form.Get("end_date")
		leaveType := r.Form.Get("leave_type")

		// normalize dates so leave continues from 00:00 to 23:59
		start, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			panic(err)
		}
		start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())

		end, err := time.Parse("2006-01-02", endDate)
		if err != nil {
			panic(err)
		}

		end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 0, end.Location())

		err = app.ls.Create(app.userID(r), start, end, leaveType)
		if err != nil {
			panic(err)
		}

		// TODO: app.sm.SetFlash(w, "Leave request created successfully")
		// simple http redirect
		http.Redirect(w, r, "/overview", http.StatusFound)
	} else {
		// Render the form for planning leave
		tmpl := template.Must(template.ParseFiles(
			"frontend/src/templates/layout.html",
			"frontend/src/templates/plan_leave.html",
		))

		data := struct {
			CommonFormTemplateData
			CommonTemplateData
			LeaveTypes []string
		}{
			CommonFormTemplateData: CommonFormTemplateData{
				Errors: []string{},
			},
			CommonTemplateData: *app.commonTemplateData(r),
			LeaveTypes:         tracker.LeaveTypes(),
		}

		tmpl.ExecuteTemplate(w, "layout", data)
	}
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

func (app *App) handleCalendar(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(
		template.
			New("calendar").
			Funcs(templateFuncs()).
			ParseFiles(
				"frontend/src/templates/fragments/calendar.html",
			),
	)

	monthNum, err := strconv.Atoi(r.URL.Query().Get("month"))
	if err != nil {
		monthNum = int(time.Now().Month())
	}
	selectedMonth := time.Month(monthNum)

	selectedYear, err := strconv.Atoi(r.URL.Query().Get("year"))
	if err != nil {
		selectedYear = time.Now().Year()
	}

	data := &CalendarTemplateData{
		CommonTemplateData: *app.commonTemplateData(r),
		Weekdays:           []string{"Mo", "Tu", "We", "Th", "Fr", "Sa", "Su"},
		MonthWeekDays:      calendarMonth(selectedYear, int(selectedMonth)),
		SelectedYear:       selectedYear,
		SelectedMonth:      selectedMonth,
		Today:              time.Now(),
		Nav: CalendarNav{
			Prev: MonthPeriod{
				Month: time.Month((int(selectedMonth)+10)%12 + 1),
				Year:  selectedYear,
			},
			Next: MonthPeriod{
				Month: time.Month((int(selectedMonth) + 1) % 12),
				Year:  selectedYear,
			},
		},
	}

	err = tmpl.ExecuteTemplate(w, "calendar.html", data)
	if err != nil {
		panic(err)
	}
}

func (app *App) htmxRedirect(w http.ResponseWriter, r *http.Request, url string) {
	w.Header().Add("HX-Redirect", url)
}

func (app *App) commonTemplateData(r *http.Request) *CommonTemplateData {
	return &CommonTemplateData{
		IsAuthenticated: app.isAuthenticated(r),
	}
}
