package service

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/igromanas/go-final/pkg/model"
)

const Layout = "20060102"

type Rule struct {
	mode   string
	first  []int
	second []time.Month // TODO
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
				rule.second = append(rule.second, time.Month(num))
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
	if date.Before(now) {
		date = now
	}

	dw := int(date.Weekday())
	for _, w := range weekdays {
		if w > dw {
			return date.AddDate(0, 0, w-dw)
		}
	}

	if weekdays[0] == dw {
		return date.AddDate(0, 0, 7)
	}

	return date.AddDate(0, 0, weekdays[0]+7-dw)
}

func createDate(year int, month time.Month, day int) time.Time {
	test := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Now().UTC().Location())
	for test.Day() < day {
		test = test.AddDate(0, 0, 1)
	}
	return test
}

func addSomeDays(now time.Time, date time.Time, days []int, months []time.Month) time.Time {
	if date.Before(now) {
		date = now
	}

	var candidates []time.Time
	if months == nil {
		for _, day := range days {
			switch day {
			case -1:
				d := createDate(date.Year(), date.Month()+1, 0) // last day of current month
				if date.After(d) || date.Day() == d.Day() {
					candidates = append(candidates, createDate(date.Year(), date.Month()+2, 0)) // last day of next month
				} else {
					candidates = append(candidates, d)
				}
			case -2:
				d := createDate(date.Year(), date.Month()+1, -1) // penultimate day of current month
				if date.After(d) || date.Day() == d.Day() {
					candidates = append(candidates, createDate(date.Year(), date.Month()+2, -1)) // penultimate day of next month
				} else {
					candidates = append(candidates, d)
				}
			default:
				d := createDate(date.Year(), date.Month(), day)
				if date.After(d) || date.Day() == d.Day() {
					candidates = append(candidates, createDate(date.Year(), date.Month()+1, day))
				} else {
					candidates = append(candidates, d)
				}
			}
		}
	} else {
		for _, m := range months {
			switch {
			case date.Month() < m:
				for _, day := range days {
					switch day {
					case -1:
						candidates = append(candidates, createDate(date.Year(), m+1, 0)) // last day of m month
					case -2:
						candidates = append(candidates, createDate(date.Year(), m+1, -1)) // penultimate day of m month
					default:
						candidates = append(candidates, createDate(date.Year(), m, day))
					}
				}
			case date.Month() > m:
				nextYear := date.Year() + 1
				for day := range days {
					switch day {
					case -1:
						candidates = append(candidates, createDate(nextYear, m+1, 0)) // last day of m month next year
					case -2:
						candidates = append(candidates, createDate(nextYear, m+1, -1)) // penultimate day of m month next year
					default:
						candidates = append(candidates, createDate(nextYear, m, day))
					}
				}
			case date.Month() == m:
				for _, day := range days {
					switch day {
					case -1:
						d := createDate(date.Year(), date.Month()+1, 0) // last day of current month
						if date.Before(d) {
							candidates = append(candidates, d)
						}
					case -2:
						d := createDate(date.Year(), date.Month()+1, -1) // penultimate day of current month
						if date.Before(d) {
							candidates = append(candidates, d)
						}
					default:
						d := createDate(date.Year(), date.Month(), day)
						if date.Before(d) {
							candidates = append(candidates, d)
						}
					}
				}
			}
		}
	}

	sort.Slice(candidates, func(i, j int) bool { return candidates[i].Before(candidates[j]) })
	return candidates[0]
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
	case "w":
		return addWeekdays(now, dstart, rule.first), nil
	case "m":
		return addSomeDays(now, dstart, rule.first, rule.second), nil
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

	return nil
}
