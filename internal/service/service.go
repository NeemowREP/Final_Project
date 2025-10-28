package service

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func weekdayNum(t time.Time) int {
	d := int(t.Weekday())
	if d == 0 {
		d = 7
	}
	return d
}

func leapYear(year int) bool {
	return (year%4 == 0 && year%100 != 0) || (year%400 == 0)
}

func matchMonthDay(date time.Time, dayMask [32]bool, monthMask [13]bool) bool {
	day := date.Day()
	month := int(date.Month())
	lastDay := time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, date.Location()).Day()

	matchDay := dayMask[day] || (day == lastDay && dayMask[31]) || (day == lastDay-1 && dayMask[30])
	matchMonth := monthMask[month]

	return matchDay && matchMonth
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
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	parts := strings.Split(repeat, " ")

	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", fmt.Errorf("недопустимый формат d")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days < 1 || days > 400 {
			return "", fmt.Errorf("недопустимое значение дней: %s", parts[1])
		}
		for !date.After(now) {
			date = date.AddDate(0, 0, days)
		}

	case "y":
		for !date.After(now) {
			year := date.Year() + 1
			month := date.Month()
			day := date.Day()
			if month == time.February && day == 29 {
				for !leapYear(year) {
					year++
				}
			}
			date = time.Date(year, month, day, 0, 0, 0, 0, date.Location())
		}

	case "w":
		if len(parts) != 2 {
			return "", fmt.Errorf("неподдерживаемый формат w")
		}
		var weekdays [8]bool
		for _, s := range strings.Split(parts[1], ",") {
			n, err := strconv.Atoi(s)
			if err != nil || n < 1 || n > 7 {
				return "", fmt.Errorf("недопустимый день недели: %s", s)
			}
			weekdays[n] = true
		}
		for !date.After(now) || !weekdays[weekdayNum(date)] {
			date = date.AddDate(0, 0, 1)
		}

	case "m":
		if len(parts) < 2 {
			return "", fmt.Errorf("неподдерживаемый формат m")
		}
		var dayMask [32]bool
		var monthMask [13]bool

		for _, s := range strings.Split(parts[1], ",") {
			n, err := strconv.Atoi(s)
			if err != nil || n < -2 || n > 31 || n == 0 {
				return "", fmt.Errorf("недопустимый день месяца: %s", s)
			}
			if n > 0 {
				dayMask[n] = true
			} else {
				dayMask[31+n+1] = true
			}
		}

		if len(parts) > 2 {
			for _, s := range strings.Split(parts[2], ",") {
				n, err := strconv.Atoi(s)
				if err != nil || n < 1 || n > 12 {
					return "", fmt.Errorf("недопустимый месяц: %s", s)
				}
				monthMask[n] = true
			}
		} else {
			for i := 1; i <= 12; i++ {
				monthMask[i] = true
			}
		}

		for !date.After(now) || !matchMonthDay(date, dayMask, monthMask) {
			date = date.AddDate(0, 0, 1)
		}

	default:
		return "", fmt.Errorf("неподдерживаемый формат: %s", repeat)
	}

	return date.Format("20060102"), nil
}
