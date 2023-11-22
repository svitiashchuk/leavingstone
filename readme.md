### Running
```shell
docker run -it --rm \
  -w "$PWD" \
  -e "air_wd=$PWD" \
  -v "$PWD":"$PWD" \
  -p 8080:8887 \
  cosmtrek/air -c .air.toml
```


- user: name, email, start-date
- request: day-off
- request: sick-day
- request: vacation

vaction-balance: start-date + months*coef_vacation_off - user_vacation_during_period
sick-date: start-date + months*coef_sick_days - user_sick_days_during_period

user:
vacation: start-date + 12*coef_vacation_off
vacation left until new period starts: vacation - user_vacation_during_period
vacation left (available): start-date + cur_month*coef_vacation_off - user_vacation_during_period


When approved - nothing changes, all calculations are done on the fly
Balance in db updated only in New Year (5 extra transferred) manually

For new people who works less than 12 months:
beginning for calculation day-by-day is the day they started to work

For new people who works more than 12 months:
beginning for calculation day-by-day is 1 Jan


Today is 23 Apr 2023
User started to work on 21 Apr 2021
Calculation: from 1 Jan 2023 to today (daily) + Extra days transferred from 2022 (max 5 days)

User started to work on 21 Oct 2022
Calculation: from 1 Jan 2023 to today (daily) + Extra days transferred from 2022 (max 5 days)

User started to work on 21 Feb 2023
Calculation: from 1 Feb 2023 to today (daily)


So, if user joined company this year - they don't have extra days and calulation is done from the day they started to work
If user joined company last year - they have extra days and calulation is done from 1 Jan of this year

```go
	thisYear := now.Year()
	lastYearlyPeriod := time.Date(thisYear, initial.Month(), initial.Day(), 0, 0, 0, 0, time.UTC)

	if time.Now().Before(lastYearlyPeriod) {
		lastYearlyPeriod = time.Date(thisYear-1, initial.Month(), initial.Day(), 0, 0, 0, 0, time.UTC)
	}
```

---

How to run:

