package timefn

import "time"

// StartOfSecond returns a new time.Time value representing the start of the
// second for the given time. The returned time will have its nanosecond field
// set to 0, effectively rounding down to the nearest second.
func StartOfSecond(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, t.Location())
}

// EndOfSecond returns the last nanosecond of the second specified by the
// provided [time.Time]. The returned [time.Time] has the same year, month, day,
// hour, minute and second as the input, but its nanoseconds field is set to one
// less than a full second.
func EndOfSecond(t time.Time) time.Time {
	return StartOfSecond(t).Add(time.Second).Add(-time.Nanosecond)
}

// StartOfMinute returns a new instance of [time.Time] representing the start of
// the minute of the provided time, with the second and nanosecond fields set to
// zero. The returned time maintains the same year, month, day, hour and minute
// as the input but resets the second and nanosecond to their minimum possible
// values.
func StartOfMinute(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location())
}

// EndOfMinute returns the exact time at the end of the minute for a given time.
// This is effectively one nanosecond before the start of the next minute.
func EndOfMinute(t time.Time) time.Time {
	return StartOfMinute(t).Add(time.Minute).Add(-time.Nanosecond)
}

// StartOfHour returns a new instance of [time.Time] representing the start of
// the hour of the provided time value. All minutes, seconds, and nanoseconds
// values are set to zero while year, month, day, and hour are preserved from
// the original time value. The location (time zone) of the returned time is
// also preserved.
func StartOfHour(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location())
}

// EndOfHour returns the time instance representing the end of the hour for the
// provided time [t]. The end of the hour is defined as one nanosecond before
// the start of the next hour.
func EndOfHour(t time.Time) time.Time {
	return StartOfHour(t).Add(time.Hour).Add(-time.Nanosecond)
}

// StartOfDay returns a new instance of [time.Time] representing the start of
// the day of the given time, with the hour, minute, second, and nanosecond
// fields set to zero while maintaining the same year, month, day and location
// as the original.
func StartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// EndOfDay returns the end of the day for a given time, represented as a
// time.Time value. The end of the day is defined as the last possible moment
// before the start of the next day. This is equivalent to one nanosecond before
// midnight of the next day in the same location as the input time.
func EndOfDay(t time.Time) time.Time {
	return StartOfDay(t).AddDate(0, 0, 1).Add(-time.Nanosecond)
}

// StartOfWeek returns the start of the week for a given time. The week starts
// on Sunday as per Go's time package definition. The returned time has the same
// location and date but the hour, minute, second, and nanosecond are set to
// their zero values.
func StartOfWeek(t time.Time) time.Time {
	return StartOfDay(t.AddDate(0, 0, -int(t.Weekday())))
}

// EndOfWeek returns the end of the week for a given time. The end of the week
// is defined as 23:59:59 on the last day of the week, which depends on the
// Weekday of the input time. The returned time is in the same location as the
// input time.
func EndOfWeek(t time.Time) time.Time {
	return EndOfDay(t.AddDate(0, 0, 6-int(t.Weekday())))
}

// StartOfISOWeek returns a new time.Time representing the start of the ISO 8601
// week for the given time. The start of a week is considered to be Monday. The
// returned time has the same location and year, month, and day fields as t but
// the hour, minute, second, and nanosecond fields are all set to zero.
func StartOfISOWeek(t time.Time) time.Time {
	return StartOfDay(t.AddDate(0, 0, -int((t.Weekday()+6)%7)))
}

// EndOfISOWeek returns the last instant within the same ISO week as the
// provided [time.Time]. An ISO week starts on Monday and ends on Sunday. The
// returned time will be at the end of the day, just before midnight, in the
// location of the provided time.
func EndOfISOWeek(t time.Time) time.Time {
	return EndOfDay(t.AddDate(0, 0, 6-int((t.Weekday()+6)%7)))
}

// StartOfMonth returns a new instance of [time.Time] set to the first day of
// the provided time's month, with the hour, minute, second, and nanosecond
// fields set to zero. The location is preserved.
func StartOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

// EndOfMonth takes a [time.Time] value and returns a new [time.Time] value
// representing the exact end of the same month. The returned time is the last
// nanosecond of the last day of the month in the same location as the input
// time.
func EndOfMonth(t time.Time) time.Time {
	return StartOfMonth(t).AddDate(0, 1, 0).Add(-time.Nanosecond)
}

// StartOfYear returns the time representing the start of the year for the given
// time [t]. The returned time will have a date component equal to January 1st
// of the year of [t], and a time component set to midnight in [t]'s location.
func StartOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
}

// EndOfYear returns the latest possible time within the same year as the given
// time [t]. The returned time is one nanosecond before the start of the next
// year.
func EndOfYear(t time.Time) time.Time {
	return StartOfYear(t).AddDate(1, 0, 0).Add(-time.Nanosecond)
}

// Between checks if a given time [t] falls after time [l] and before time [r].
// Returns true if [t] is between [l] and [r], otherwise returns false.
func Between(t, l, r time.Time) bool {
	return t.After(l) && t.Before(r)
}

// BetweenInclusive determines if a given time ([time.Time]) is within or equal
// to the range of two other specified times, inclusive. The function returns
// true if the time is equal to or falls between the other two times, otherwise
// it returns false.
func BetweenInclusive(t, l, r time.Time) bool {
	return SameOrBefore(t, r) && SameOrAfter(t, l)
}

// SameOrBefore checks whether a given time is the same as or precedes another
// specified time. It takes two arguments, both of type [time.Time], and returns
// a boolean. If the first argument is either equal to or comes before the
// second argument in chronological order, the function returns true. Otherwise,
// it returns false.
func SameOrBefore(t, r time.Time) bool {
	return t.Equal(r) || t.Before(r)
}

// SameOrAfter determines if the given time [t] is the same as or after another
// time [l]. It returns true if [t] is either equal to [l] or occurs after [l],
// and false otherwise.
func SameOrAfter(t, l time.Time) bool {
	return t.Equal(l) || t.After(l)
}

// AtTime sets the time of the provided [time.Time] to the specified hours,
// minutes, seconds and nanoseconds while keeping the same date and location.
// The returned [time.Time] will have the same year, month, day and location as
// the original, but the hour, minute, second and nanosecond values will be
// replaced with the ones provided as arguments.
func AtTime(t time.Time, h, m, s, ns int) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), h, m, s, ns, t.Location())
}
