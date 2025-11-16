package service

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/igromanas/go-final/pkg/model"
)

const Layout = "20060102"

type Rule struct {
	mode   string
	first  []int
	second []int // TODO
}

func parseRule(input string, rule *Rule) error {
	parts := strings.Split(input, " ")
	rule.mode = parts[0]

	switch rule.mode {
	case "y":
		return nil
	case "d":
		if len(parts) < 2 {
			return fmt.Errorf("not enough parameters for mode '%s'", rule.mode)
		}
		num, err := strconv.Atoi(parts[1])
		if err != nil {
			return fmt.Errorf("invalid parameter error: %w", err)
		}
		if num < 1 || num > 400 {
			return fmt.Errorf("number out of range: %d", num)
		}
		rule.first = append(rule.first, num)
	case "w":
		if len(parts) < 2 {
			return fmt.Errorf("not enough parameters for mode '%s'", rule.mode)
		}
		nums := strings.Split(parts[1], ",")
		for _, n := range nums {
			num, err := strconv.Atoi(n)
			if err != nil {
				return fmt.Errorf("invalid parameter error: %w", err)
			}
			rule.first = append(rule.first, num)
		}
	case "m":
		if len(parts) < 2 {
			return fmt.Errorf("not enough parameters for mode '%s'", rule.mode)
		}
		nums := strings.Split(parts[1], ",")
		for _, n := range nums {
			num, err := strconv.Atoi(n)
			if err != nil {
				return fmt.Errorf("invalid first parameter error: %w", err)
			}
			rule.first = append(rule.first, num)
		}
		if len(parts) == 3 {
			nums := strings.Split(parts[2], ",")
			for _, n := range nums {
				num, err := strconv.Atoi(n)
				if err != nil {
					return fmt.Errorf("invalid second parameter error: %w", err)
				}
				rule.second = append(rule.second, num)
			}
		}
	}

	return nil
}

func addYears(now time.Time, date time.Time, years int) time.Time {
	for {
		date = date.AddDate(years, 0, 0)
		if now.Before(date) {
			break
		}
	}
	return date
}

func addDays(now time.Time, date time.Time, days int) time.Time {
	for {
		date = date.AddDate(0, 0, days)
		if now.Before(date) {
			break
		}
	}
	return date
}

func addWeekdays(now time.Time, date time.Time, weekdays []int) time.Time {
	if now.Before(date) {
		date = now
	}
	dw := int(date.Weekday())
	for _, w := range weekdays {
		if w > dw {
			return date.AddDate(0, 0, w-dw)
		}
	}
	return date.AddDate(0, 0, weekdays[0]-dw)
}

func NextDate(now time.Time, dstart time.Time, repeat string) (time.Time, error) {
	rule := Rule{}
	err := parseRule(repeat, &rule)
	if err != nil {
		return time.Time{}, err
	}
	switch rule.mode {
	case "d":
		return addDays(now, dstart, rule.first[0]), nil

	case "y":
		return addYears(now, dstart, 1), nil
	// == TODO ==
	case "w":
		return addWeekdays(now, dstart, rule.first), nil
	// case "m":
	default:
		return time.Time{}, fmt.Errorf("unexpected rule mode")
	}

}

func ModifyDate(task *model.Task) error {
	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format(Layout)
		return nil
	}

	dstart, err := time.Parse(Layout, task.Date)
	if err != nil {
		return fmt.Errorf("date parse error: %v", err)
	}

	if now.Format(Layout) > task.Date {
		if task.Repeat == "" {
			task.Date = now.Format(Layout)
		} else {
			nd, err := NextDate(now, dstart, task.Repeat)
			if err != nil {
				return fmt.Errorf("next date error: %v", err)
			}
			task.Date = nd.Format(Layout)
		}
	}
	// var next string
	// if task.Repeat != "" {
	// 	nd, err := NextDate(now, dstart, task.Repeat)
	// 	if err != nil {
	// 		return fmt.Errorf("next date error: {%v}", err)
	// 	}
	// 	next = nd.Format(LAYOUT)
	// }

	// // если сегодня (now) больше task.Date (t)
	// if now.After(dstart) {
	// 	if len(task.Repeat) == 0 {
	// 		// если правила повторения нет, то берём сегодняшнее число
	// 		task.Date = now.Format(LAYOUT)
	// 	} else {
	// 		// в противном случае, берём вычисленную ранее следующую дату
	// 		task.Date = next
	// 	}
	// }

	return nil
}
