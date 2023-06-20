package server

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"ptocker"
)

// Day types
const (
	Vacation    = "Vacation"
	SickDay     = "SickDay"
	DayOff      = "DayOff"
	Workday     = "Workday"
	Weekend     = "Weekend"
	BankHoliday = "BankHoliday"
)

type MonthPeriod struct {
	Month int
	Year  int
}

type Navigation struct {
	Prev MonthPeriod
	Now  MonthPeriod
	Next MonthPeriod
}

func Serve() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles(
			"templates/layout.html",
			"templates/list.html",
		))
		tmpl.ExecuteTemplate(w, "layout", nil)
	})

	http.HandleFunc("/tracker", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles(
			"templates/fragments/tracker.html",
		))

		y, err := strconv.Atoi(r.URL.Query().Get("year"))
		if err != nil {
			y = time.Now().Year()
		}

		m, err := strconv.Atoi(r.URL.Query().Get("month"))
		if err != nil {
			m = int(time.Now().Month())
		}

		// days := days(time.Now().Year())
		days := month(y, m)
		uut := usersForTemplate(days, users())

		next := time.Date(y, time.Month(m+1), 1, 0, 0, 0, 0, time.UTC)
		prev := time.Date(y, time.Month(m-1), 1, 0, 0, 0, 0, time.UTC)

		fmt.Println(prev)
		fmt.Println(next)

		nav := Navigation{
			Prev: MonthPeriod{Month: int(prev.Month()), Year: prev.Year()},
			Now:  MonthPeriod{Month: int(time.Month(m)), Year: y},
			Next: MonthPeriod{Month: int(next.Month()), Year: next.Year()},
		}

		data := map[string]interface{}{
			"Nav":   nav,
			"Users": uut,
			"Days":  days,
		}

		fmt.Println(tmpl.DefinedTemplates())
		tmpl.ExecuteTemplate(w, "tracker.html", data)
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles(
			"templates/layout.html",
			"templates/login.html",
		))

		if r.Method == "POST" {
			r.ParseForm()

			// TODO - retrieve user from DB and set token in cookie

			w.Header().Add("HX-Redirect", "/login")
			w.Header().Add("Set-Cookie", "auth_token=TOKEN; Path=/; HttpOnly")
			tmpl.ExecuteTemplate(w, "layout", nil)
		} else {
			fmt.Printf(r.Method)
			fmt.Println(tmpl.DefinedTemplates())
			tmpl.ExecuteTemplate(w, "layout", nil)
		}
	})

	http.ListenAndServe(":8080", nil)
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

type Leaves map[string]string

func users() []*ptocker.User {
	return []*ptocker.User{
		// {Name: "John", Leaves: Leaves{"01.01.2023": Vacation, "02.01.2023": SickDay}},
		// {Name: "Oliver", Leaves: Leaves{"03.02.2023": Vacation, "06.02.2023": Vacation, "07.02.2023": Vacation}},
	}
}

type Day struct {
	Date time.Time
	Type string
}

type UserForTemplate struct {
	Name     string
	UserDays []Day
}

func usersForTemplate(days []time.Time, uu []User) []UserForTemplate {
	uut := make([]UserForTemplate, len(uu))

	for i, u := range uu {

		ud := make([]Day, len(days))
		for j, d := range days {
			ud[j] = Day{Date: d, Type: u.DayType(d)}
		}

		uut[i] = UserForTemplate{
			Name:     u.Name,
			UserDays: ud,
		}
	}

	return uut
}

func (u *User) DayType(day time.Time) string {
	if leave, ok := u.Leaves[day.Format("02.01.2006")]; ok {
		return leave
	}

	if isWeekend(day) {
		return Weekend
	}

	if isBankHoliday(day) {
		return BankHoliday
	}

	return Workday
}

func isWeekend(day time.Time) bool {
	return day.Weekday() == 0 || day.Weekday() == 6
}

func isBankHoliday(day time.Time) bool {
	return false
}

func GetMonthNumber(month time.Month) int {
	return int(month)
}
