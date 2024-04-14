package enums

import (
	"time"
)

type DateRange int

const (
	Today DateRange = iota + 1
	Yesterday
	Last7Days
	LastWeek
	Last30Days
	Custom
)

var DateRangeMap = map[string]DateRange{
	"Yesterday":  Yesterday,
	"Today":      Today,
	"Last7Days":  Last7Days,
	"LastWeek":   LastWeek,
	"Last30Days": Last30Days,
	"Custom":     Custom,
}

func (d DateRange) GetValues() []DateRange {
	var values []DateRange
	for _, v := range DateRangeMap {
		values = append(values, v)
	}

	return values
}

func (d DateRange) GetValue(key string) DateRange {
	return DateRangeMap[key]
}

func (d DateRange) GetRanges(to, from time.Time) (time.Time, time.Time) {
	var start time.Time
	var end time.Time
	now := time.Now()
	switch d {
	case Yesterday:
		start = now.Add(-24 * time.Hour)
		end = now.Add(-24 * time.Hour)
	case Today:
		start = now
		end = now
	case Last7Days:
		start = now.AddDate(0, 0, -7)
		end = now
	case LastWeek:
		now = time.Date(2024, 4, 7, 0, 0, 0, 0, now.Location())
		weekDay := now.Weekday()
		if weekDay == time.Saturday {
			end = now.AddDate(0, 0, -7)
			start = end.AddDate(0, 0, -6)
		} else {
			end = now.AddDate(0, 0, int(time.Saturday)-now.Day())
			start = end.AddDate(0, 0, -6)
		}
	case Last30Days:
		start = now.AddDate(0, 0, -30)
		end = now
	case Custom:
		start = to
		end = from
	default:
		return start, end
	}

	return start, end
}
