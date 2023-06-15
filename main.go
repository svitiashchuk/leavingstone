package main

import (
	"fmt"
	"ptocker/server"
	"time"
)

const YearDays = 366
const VacationsPerYear = 15
const SickDaysPerYear = 10

type Vacation struct {
	started  time.Time
	ended    time.Time
	approved bool
}

type User struct {
	name  string
	email string
	// date when user started to work:
	started        time.Time
	extraVacations int
}

func main() {
	server.Serve()
}

func experiment() {
	// create a new user:
	u := User{
		name:           "John",
		email:          "john.doe@ptocker.com",
		started:        time.Date(2022, 8, 21, 0, 0, 0, 0, time.UTC),
		extraVacations: 0,
	}

	fmt.Printf("User: %v\n", u)
	fmt.Printf("Vacations total: %v\n", vacationsTotal(u))
	fmt.Printf("Vacations used: %v\n", usedVacations(u))
	fmt.Printf("Vacations left: %v\n", vacationsTotal(u)-usedVacations(u))
}

// todo calculate only for single year.
func usedVacations(u User) int {
	used := 0
	vv := []Vacation{
		{mustParseTime("2023-01-03"), mustParseTime("2023-01-05"), true}, // tu,wed,th = 3
		{mustParseTime("2023-02-14"), mustParseTime("2023-02-20"), true}, // tu->mo = 5
		// {mustParseTime("2023-04-19"), mustParseTime("2023-04-20"), true}, // we, th = 2
	}

	// TODO take into account holidays.
	// count workdays between two dates:
	for _, v := range vv {
		for d := v.started; d.Before(v.ended) || d.Equal(v.ended); d = d.Add(24 * time.Hour) {
			if d.Weekday() == time.Saturday || d.Weekday() == time.Sunday {
				continue
			}

			used++
		}
	}

	return used
}

func vacationsTotal(u User) int {
	vacationsForPeriod := 0

	if u.started.Year() == time.Now().Year() {
		vacationsForPeriod = daysSinceJoined(u.started) * VacationsPerYear / YearDays
	} else {
		vacationsForPeriod = daysSinceNewYear() * VacationsPerYear / YearDays
	}

	return vacationsForPeriod + u.extraVacations
}

func daysSinceNewYear() int {
	now := time.Now()
	thisYear := now.Year()
	newYear := time.Date(thisYear, 1, 1, 0, 0, 0, 0, time.UTC)

	daysSince := int(now.Sub(newYear).Hours() / 24)

	return daysSince
}

func daysSinceJoined(initial time.Time) int {
	now := time.Now()
	daysSince := int(now.Sub(initial).Hours() / 24)

	return daysSince
}

func mustParseTime(d string) time.Time {
	t, err := time.Parse("2006-01-02", d)
	if err != nil {
		panic(err)
	}

	return t
}
