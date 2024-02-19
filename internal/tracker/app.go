package tracker

import (
	"fmt"
	"leavingstone/internal/auth"
	"leavingstone/internal/middleware"
	"leavingstone/internal/model"
	"leavingstone/internal/session"
	"leavingstone/internal/sqlite"
	"leavingstone/internal/templ"
	"log/slog"
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
	sm          *session.Keeper
	auth        *auth.Authenticator
	us          *sqlite.UserService
	ls          *sqlite.LeaveService
	teamService *sqlite.TeamService
	t           *Tracker
	ac          *Accountant
	errorLogger *slog.Logger
	appLogger   *slog.Logger
}

func NewApp(
	sm *session.Keeper,
	auth *auth.Authenticator,
	us *sqlite.UserService,
	ls *sqlite.LeaveService,
	t *Tracker,
	ac *Accountant,
	appLogger *slog.Logger,
	errorLogger *slog.Logger,
) *App {
	return &App{
		sm:          sm,
		auth:        auth,
		us:          us,
		ls:          ls,
		t:           t,
		ac:          ac,
		appLogger:   appLogger,
		errorLogger: errorLogger,
	}
}

func (app *App) userID(r *http.Request) int {
	return r.Context().Value(auth.UserIDContextKey).(int)
}

type CommonTemplateData struct {
	IsAuthenticated bool
	Alert           string
}

type CommonFormTemplateData struct {
	Errors []string
}

type ProfileTemplateData struct {
	CommonTemplateData
	UpcomingLeaves []*model.Leave
	User           *model.User
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
	Employees     []*Employee
	Days          []time.Time
	WorkforceStat *WorkforceStat
	LeavesStat    *LeavesStat
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

type OverviewTemplateData struct {
	CommonTemplateData
	UpcomingLeaves []*model.Leave
}

type CalendarNav struct {
	Prev MonthPeriod
	Next MonthPeriod
}

func (app *App) RegisterRoutes() {

	// assets for frontend
	http.HandleFunc("/dist/", app.handleDist)

	// main routes
	mainMiddleware := middleware.
		NewChain().
		Use(app.sm.Provide).
		Use(app.auth.Authenticate)

	http.HandleFunc("/login", app.handleLogin)
	http.HandleFunc("/logout", mainMiddleware.Then(app.handleLogout))
	http.HandleFunc("/", mainMiddleware.Then(app.handleIndex))
	http.HandleFunc("/profile", mainMiddleware.Then(app.handleProfile))
	http.HandleFunc("/overview", mainMiddleware.Then(app.handleOverview))
	http.HandleFunc("/teams/create", mainMiddleware.Then(app.CreateTeam))

	http.HandleFunc("/leaves/plan", mainMiddleware.Then(app.handlePlanLeave))
	http.HandleFunc("/leaves/approve", mainMiddleware.Then(app.handleLeaveApprove))
	http.HandleFunc("/leaves/reject", mainMiddleware.Then(app.handleLeaveReject))

	// fragments
	http.HandleFunc("/tracker", mainMiddleware.Then(app.handleTracker))
	http.HandleFunc("/fragments/calendar", mainMiddleware.Then(app.handleCalendar))
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
			app.internalError(w, err)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(passPlain))
		if err != nil {
			app.internalError(w, err)
			return
		}

		c := fmt.Sprintf("auth_token=%s; Path=/; HttpOnly", u.Token)
		w.Header().Add("Set-Cookie", c)

		http.Redirect(w, r, "/overview", http.StatusFound)
	}

	tmpl.ExecuteTemplate(w, "layout", app.commonTemplateData(r))
}

func (app *App) handleLogout(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Set-Cookie", "auth_token=; Path=/; HttpOnly")
	http.Redirect(w, r, "/login", http.StatusFound)
}

func (app *App) handleProfile(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(
		template.
			New("profile").
			Funcs(templ.Funcs()).
			ParseFiles(
				"frontend/src/templates/layout.html",
				"frontend/src/templates/profile.html",
			),
	)

	u, err := app.us.FindByID(app.userID(r))
	if err != nil {
		app.internalError(w, err)
		return
	}

	leaves, err := app.ls.Upcoming(app.userID(r))
	if err != nil {
		app.internalError(w, err)
		return
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

	if err := tmpl.ExecuteTemplate(w, "layout", templateData); err != nil {
		app.internalError(w, err)
		return
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
			app.clientError(w, err)
			return
		}
		start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())

		end, err := time.Parse("2006-01-02", endDate)
		if err != nil {
			app.clientError(w, err)
			return
		}

		end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 0, end.Location())

		err = app.ls.Create(app.userID(r), start, end, leaveType)
		if err != nil {
			app.internalError(w, err)
			return
		}

		app.sm.Get(r.Context().Value(session.SessionContextKey).(string)).Flash("Leave planned!")
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
			LeaveTypes:         LeaveTypes(),
		}

		tmpl.ExecuteTemplate(w, "layout", data)
	}
}

func (app *App) handleOverview(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(
		template.
			New("overview").
			Funcs(templ.Funcs()).
			ParseFiles(
				"frontend/src/templates/layout.html",
				"frontend/src/templates/overview.html",
			),
	)

	ll, err := app.ls.AllUpcoming()
	if err != nil {
		app.internalError(w, err)
		return
	}

	data := &OverviewTemplateData{
		CommonTemplateData: *app.commonTemplateData(r),
		UpcomingLeaves:     ll,
	}

	err = tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		app.internalError(w, err)
		return
	}
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
		app.internalError(w, err)
		return
	}

	w.Header().Add("Content-type", "text/css")
	w.Write(b)
}

func (app *App) handleLeaveApprove(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.internalError(w, err)
		return
	}

	id, err := strconv.Atoi(r.Form.Get("id"))
	if err != nil {
		app.internalError(w, err)
		return
	}

	app.t.ApproveLeave(id)

	// send hx-trigger header to reload full tracker
	w.Header().Add("HX-Trigger", "reloadTracker")
}

func (app *App) handleLeaveReject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, err)
		return
	}

	id, err := strconv.Atoi(r.Form.Get("id"))
	if err != nil {
		app.clientError(w, err)
		return
	}

	app.t.RejectLeave(id)
	// send hx-trigger header to reload full tracker
	w.Header().Add("HX-Trigger", "reloadTracker")
}

func (app *App) handleCalendar(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(
		template.
			New("calendar").
			Funcs(templ.Funcs()).
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
		app.internalError(w, err)
		return
	}
}

func (app *App) commonTemplateData(r *http.Request) *CommonTemplateData {
	var alert string

	ctxVal := r.Context().Value(session.SessionContextKey)
	if ctxVal != nil {
		session := app.sm.Get(ctxVal.(string))
		if session != nil {
			alert = session.GetFlash()
		}
	}

	return &CommonTemplateData{
		IsAuthenticated: true,
		Alert:           alert,
	}
}
