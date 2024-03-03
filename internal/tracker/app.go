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
	templator   *templ.Templator
	errorLogger *slog.Logger
	appLogger   *slog.Logger
}

func NewApp(
	sm *session.Keeper,
	auth *auth.Authenticator,
	us *sqlite.UserService,
	ls *sqlite.LeaveService,
	teamService *sqlite.TeamService,
	t *Tracker,
	ac *Accountant,
	templator *templ.Templator,
	appLogger *slog.Logger,
	errorLogger *slog.Logger,
) *App {
	return &App{
		sm:          sm,
		auth:        auth,
		us:          us,
		ls:          ls,
		teamService: teamService,
		t:           t,
		ac:          ac,
		templator:   templator,
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
	Theme           string
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
	http.HandleFunc("/teams/details", mainMiddleware.Then(app.TeamDetails))
	http.HandleFunc("/teams/add-member", mainMiddleware.Then(app.handleAddMember))
	http.HandleFunc("/teams/delete-member", mainMiddleware.Then(app.handleDeleteMember))

	http.HandleFunc("/leaves/plan", mainMiddleware.Then(app.handlePlanLeave))
	http.HandleFunc("/leaves/approve", mainMiddleware.Then(app.handleLeaveApprove))
	http.HandleFunc("/leaves/reject", mainMiddleware.Then(app.handleLeaveReject))

	http.HandleFunc("/settings", mainMiddleware.Then(app.handleSettings))

	// fragments
	http.HandleFunc("/tracker", mainMiddleware.Then(app.handleTracker))
	http.HandleFunc("/fragments/calendar", mainMiddleware.Then(app.handleCalendar))
	http.HandleFunc("/fragments/teams/delete-member-dialog", mainMiddleware.Then(app.handleDeleteMemberDialog))
	http.HandleFunc("/fragments/teams/search-members", mainMiddleware.Then(app.handleSearchMembers))

	http.HandleFunc("/fragments/leaves/decision-dialog", mainMiddleware.Then(app.handleLeaveDecisionDialog))

	// settings
	http.HandleFunc("/user-settings/theme", mainMiddleware.Then(app.handleThemeChange))
}

func (app *App) handleIndex(w http.ResponseWriter, r *http.Request) {
	if err := app.templator.Page("list", nil).ExecuteTemplate(w, "layout", app.commonTemplateData(r)); err != nil {
		app.internalError(w, err)
		return
	}
}

func (app *App) handleLogin(w http.ResponseWriter, r *http.Request) {
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
		return
	}

	tmpl := app.templator.Page("login", nil)
	if err := tmpl.ExecuteTemplate(w, "layout", app.commonTemplateData(r)); err != nil {
		app.internalError(w, err)
		return
	}
}

func (app *App) handleLogout(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Set-Cookie", "auth_token=; Path=/; HttpOnly")
	http.Redirect(w, r, "/login", http.StatusFound)
}

func (app *App) handleProfile(w http.ResponseWriter, r *http.Request) {
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

	tmpl := app.templator.Page("profile", templ.Funcs())
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
		return
	} else {
		// Render the form for planning leave
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

		tmpl := app.templator.Page("plan_leave", templ.Funcs())
		if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
			app.internalError(w, err)
			return
		}
	}
}

func (app *App) handleOverview(w http.ResponseWriter, r *http.Request) {
	ll, err := app.ls.AllUpcoming()
	if err != nil {
		app.internalError(w, err)
		return
	}

	data := &OverviewTemplateData{
		CommonTemplateData: *app.commonTemplateData(r),
		UpcomingLeaves:     ll,
	}

	tmpl := app.templator.Page("overview", templ.Funcs())
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		app.internalError(w, err)
		return
	}
}

func (app *App) handleTracker(w http.ResponseWriter, r *http.Request) {
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

	tmpl := app.templator.Fragment("tracker", nil)
	if err := tmpl.ExecuteTemplate(w, "tracker.html", data); err != nil {
		app.internalError(w, err)
		return
	}
}

func (app *App) handleDist(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]
	// TODO use embed:
	b, err := os.ReadFile("frontend/" + path)
	if err != nil {
		app.internalError(w, err)
		return
	}

	if path[len(path)-4:] == ".svg" {
		w.Header().Add("Content-type", "image/svg+xml")
	} else {
		w.Header().Add("Content-type", "text/css")
	}
	w.Write(b)
}

func (app *App) handleLeaveDecisionDialog(w http.ResponseWriter, r *http.Request) {
	leaveID, err := strconv.Atoi(r.URL.Query().Get("leave_id"))
	if err != nil {
		app.internalError(w, err)
		return
	}

	tmpl := app.templator.Fragment("leave_decision_dialog", nil)
	if err := tmpl.ExecuteTemplate(w, "leave_decision_dialog.html", struct{ LeaveID int }{LeaveID: leaveID}); err != nil {
		app.internalError(w, err)
		return
	}
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

	tmpl := app.templator.Fragment("calendar", templ.Funcs())
	if err := tmpl.ExecuteTemplate(w, "calendar.html", data); err != nil {
		app.internalError(w, err)
		return
	}
}

func (app *App) commonTemplateData(r *http.Request) *CommonTemplateData {
	var alert string
	theme := "dark"

	ctxVal := r.Context().Value(session.SessionContextKey)
	if ctxVal != nil {
		session := app.sm.Get(ctxVal.(string))
		if session != nil {
			alert = session.GetFlash()
		}

		theme = session.Get("theme")
	}

	return &CommonTemplateData{
		IsAuthenticated: true,
		Alert:           alert,
		Theme:           theme,
	}
}
