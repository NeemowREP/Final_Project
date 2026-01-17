package service

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func AfterNow(date, now time.Time) bool {
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return date.After(now)
}

func leapYear(year int) bool {
	return (year%4 == 0 && year%100 != 0) || (year%400 == 0)
}

func weekdayNum(t time.Time) int {
	d := int(t.Weekday())
	if d == 0 {
		d = 7
	}
	return d
}

func contains(arr []int, val int) bool {
	for _, x := range arr {
		if x == val {
			return true
		}
	}
	return false
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("repeat пустой")
	}

	date, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", fmt.Errorf("неверный формат dstart: %w", err)
	}
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	parts := strings.Split(repeat, " ")
	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", fmt.Errorf("не указан интервал дней")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days < 1 || days > 400 {
			return "", fmt.Errorf("недопустимый интервал дней: %s", parts[1])
		}
		date = date.AddDate(0, 0, days)
		for !AfterNow(date, now) {
			date = date.AddDate(0, 0, days)
		}

	case "y":
		date = date.AddDate(1, 0, 0)
		for !AfterNow(date, now) {
			year := date.Year() + 1
			month := date.Month()
			day := date.Day()
			if month == time.February && day == 29 {
				for !leapYear(year) {
					year++
				}
				month = time.March
				day = 1
			}
			date = time.Date(year, month, day, 0, 0, 0, 0, date.Location())
		}

	case "w":
		if len(parts) != 2 {
			return "", fmt.Errorf("неподдерживаемый формат w")
		}
		daysOfWeek := []int{}
		for _, s := range strings.Split(parts[1], ",") {
			n, err := strconv.Atoi(s)
			if err != nil || n < 1 || n > 7 {
				return "", fmt.Errorf("недопустимый день недели: %s", s)
			}
			daysOfWeek = append(daysOfWeek, n)
		}
		date = date.AddDate(0, 0, 1)
		for !AfterNow(date, now) || !contains(daysOfWeek, weekdayNum(date)) {
			date = date.AddDate(0, 0, 1)
		}

	case "m":
		if len(parts) < 2 {
			return "", fmt.Errorf("неподдерживаемый формат m")
		}
		daysOfMonth := []int{}
		for _, s := range strings.Split(parts[1], ",") {
			n, err := strconv.Atoi(s)
			if err != nil || n == 0 || n < -2 || n > 31 {
				return "", fmt.Errorf("недопустимый день месяца: %s", s)
			}
			daysOfMonth = append(daysOfMonth, n)
		}
		months := []int{}
		if len(parts) > 2 {
			for _, s := range strings.Split(parts[2], ",") {
				n, err := strconv.Atoi(s)
				if err != nil || n < 1 || n > 12 {
					return "", fmt.Errorf("недопустимый месяц: %s", s)
				}
				months = append(months, n)
			}
		} else {
			for i := 1; i <= 12; i++ {
				months = append(months, i)
			}
		}
		matchDate := func(d time.Time) bool {
			if !contains(months, int(d.Month())) {
				return false
			}
			day := d.Day()
			last := time.Date(d.Year(), d.Month()+1, 0, 0, 0, 0, 0, d.Location()).Day()
			for _, x := range daysOfMonth {
				if x > 0 && x == day {
					return true
				}
				if x == -1 && day == last {
					return true
				}
				if x == -2 && day == last-1 {
					return true
				}
			}
			return false
		}
		date = date.AddDate(0, 0, 1)
		for !AfterNow(date, now) || !matchDate(date) {
			date = date.AddDate(0, 0, 1)
		}

	default:
		return "", fmt.Errorf("неподдерживаемый формат: %s", repeat)
	}

	return date.Format("20060102"), nil
}
