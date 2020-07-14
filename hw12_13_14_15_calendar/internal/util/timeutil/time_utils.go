package timeutil

import "time"

const (
	DayDuration  = 86400
	WeekDuration = 604800
)

// beginning of the day.
func Bod(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

// end of the day.
func Eod(t time.Time) time.Time {
	year, month, day := t.Date()
	start := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	return start.Add(time.Second * (DayDuration - 1))
}

// beginning of the week.
func Bow(t time.Time) time.Time {
	weekday := time.Duration(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	year, month, day := t.Date()
	today := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	return today.Add(-1 * (weekday - 1) * 24 * time.Hour)
}

// end of the week.
func Eow(t time.Time) time.Time {
	weekday := time.Duration(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	year, month, day := t.Date()
	today := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	return today.Add(-1 * (weekday - 1) * 24 * time.Hour).Add(time.Second * (WeekDuration - 1))
}

// beginning of the month.
func Bom(t time.Time) time.Time {
	year, month, _ := t.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
}

// end of the month.
func Eom(t time.Time) time.Time {
	year, month, _ := t.Date()
	bom := time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
	return bom.AddDate(0, 1, -1).Add(time.Second * (DayDuration - 1))
}

// checks if myTime is between start and end range.
func TimeInRange(myTime, start, end time.Time) bool {
	if (myTime.After(start) || myTime.Equal(start)) && (myTime.Before(end) || myTime.Equal(end)) {
		return true
	}
	return false
}
