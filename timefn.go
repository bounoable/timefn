package timefn

import "time"

// StartOfSecond returns the start of the second.
func StartOfSecond(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, t.Location())
}

// EndOfSecond returns the end of the second.
func EndOfSecond(t time.Time) time.Time {
	return StartOfSecond(t).Add(time.Second).Add(-time.Nanosecond)
}

// StartOfMinute returns the start of the minute.
func StartOfMinute(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location())
}

// EndOfMinute returns the end of the minute.
func EndOfMinute(t time.Time) time.Time {
	return StartOfMinute(t).Add(time.Minute).Add(-time.Nanosecond)
}

// StartOfHour returns the start of the hour.
func StartOfHour(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location())
}

// EndOfHour returns the end of the hour.
func EndOfHour(t time.Time) time.Time {
	return StartOfHour(t).Add(time.Hour).Add(-time.Nanosecond)
}

// StartOfDay returns the start of the day.
func StartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// EndOfDay returns the end of the day.
func EndOfDay(t time.Time) time.Time {
	return StartOfDay(t).AddDate(0, 0, 1).Add(-time.Nanosecond)
}

// StartOfWeek returns the start (sunday) of the week.
func StartOfWeek(t time.Time) time.Time {
	return StartOfDay(t.AddDate(0, 0, -int(t.Weekday())))
}

// EndOfWeek returns the end (saturday) of the week.
func EndOfWeek(t time.Time) time.Time {
	return EndOfDay(t.AddDate(0, 0, 6-int(t.Weekday())))
}

// StartOfISOWeek returns the ISO start (monday) of the week.
func StartOfISOWeek(t time.Time) time.Time {
	return StartOfDay(t.AddDate(0, 0, -int((t.Weekday()+6)%7)))
}

// EndOfISOWeek returns the ISO end (sunday) of the week.
func EndOfISOWeek(t time.Time) time.Time {
	return EndOfDay(t.AddDate(0, 0, 6-int((t.Weekday()+6)%7)))
}

// StartOfMonth returns the start of the month.
func StartOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

// EndOfMonth returns the start of the month.
func EndOfMonth(t time.Time) time.Time {
	return StartOfMonth(t).AddDate(0, 1, 0).Add(-time.Nanosecond)
}

// StartOfYear returns the start of the year.
func StartOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
}

// EndOfYear returns the end of the year.
func EndOfYear(t time.Time) time.Time {
	return StartOfYear(t).AddDate(1, 0, 0).Add(-time.Nanosecond)
}

// Between determines if t is between l and r (exclusive).
func Between(t, l, r time.Time) bool {
	return t.After(l) && t.Before(r)
}

// BetweenInclusive determines if t is between l and r (inclusive).
func BetweenInclusive(t, l, r time.Time) bool {
	return SameOrBefore(t, r) && SameOrAfter(t, l)
}

// SameOrBefore determines if t is the same as or before r.
func SameOrBefore(t, r time.Time) bool {
	return t.Equal(r) || t.Before(r)
}

// SameOrAfter determines if t is the same as or after l.
func SameOrAfter(t, l time.Time) bool {
	return t.Equal(l) || t.After(l)
}

// AtTime returns t with its time set to h hours, m minutes, s seconds and ns nanoseconds.
func AtTime(t time.Time, h, m, s, ns int) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), h, m, s, ns, t.Location())
}
