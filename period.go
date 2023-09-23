package timefn

import (
	"fmt"
	"math"
	"slices"
	"strings"
	"text/template"
	"time"

	"github.com/bounoable/timefn/internal/slice"
)

// DefaultPeriodFormat is a variable used as the default format for presenting a
// period of time. It uses two placeholders, .Start and .End, to represent the
// start and end of the period respectively. The format string is used in the
// Format method of the Period type to generate a string representation of the
// time period.
var DefaultPeriodFormat = "{{ .Start }} -> {{ .End }}"

// Period represents a duration of time between two points in time, defined by a
// start and end time. It provides methods for formatting the period into a
// string, validating the period, adding a duration to the period, and checking
// if a given time falls within the period.
//
// Period also includes methods for checking if it overlaps with another period,
// retrieving the years within the period, checking if the period falls within a
// specific year, getting all dates within the period, and dividing the period
// at a given date.
//
// Additionally, it provides functionalities to cut out certain periods from
// itself and return the remaining periods.
type Period struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// String returns a string representation of the Period. It leverages the Format
// method to generate this string, using the default period format defined
// within the package.
func (p Period) String() string {
	return p.Format()
}

// Format returns a string representation of the Period. The formatting is based
// on the DefaultPeriodFormat which represents the start and end time of the
// period. If there's an error during the formatting process, it returns a
// descriptive error message in a string format.
func (p Period) Format() string {
	return p.FormatAs(DefaultPeriodFormat)
}

// FormatAs formats the period using the given format string. The format string
// can contain placeholders for the start and end times of the period. If an
// empty string is passed as the format, the default format "{{ .Start }} -> {{
// .End }}" is used. If an error occurs during formatting, it returns a string
// representation of the error message.
func (p Period) FormatAs(format string) string {
	if format == "" {
		format = "{{ .Start }} -> {{ .End }}"
	}

	var buf strings.Builder
	tpl, err := template.New("").Parse(format)
	if err != nil {
		return fmt.Sprintf("<failed to format period: %s>", err)
	}

	if err := tpl.Execute(&buf, p); err != nil {
		return fmt.Sprintf("<failed to format period: %s>", err)
	}

	return buf.String()
}

// IsZero checks if the start and end times of the period are both zero,
// indicating that the period does not represent a valid time span. It returns
// true if both times are zero, and false otherwise.
func (p Period) IsZero() bool {
	return p.Start.IsZero() && p.End.IsZero()
}

// Validate checks the validity of the [Period]. It returns an error if the
// Start time is zero, the End time is zero, or if the End time is equal to or
// before the Start time. If none of these conditions are met, it returns nil
// indicating that the [Period] is valid.
func (p Period) Validate() error {
	if p.Start.IsZero() {
		return fmt.Errorf("start is zero")
	}

	if p.End.IsZero() {
		return fmt.Errorf("end is zero")
	}

	if p.End.Equal(p.Start) {
		return fmt.Errorf("end must be after start; is the same (%v)", p.End)
	}

	if p.End.Before(p.Start) {
		return fmt.Errorf("end (%v) is %v before start (%v)", p.End, p.Start.Sub(p.End), p.Start)
	}

	return nil
}

// Add extends the start and end times of the period by a specified duration. It
// returns a new Period with the updated start and end times.
func (p Period) Add(d time.Duration) Period {
	return Period{
		Start: p.Start.Add(d),
		End:   p.End.Add(d),
	}
}

// Contains checks whether a given time falls within the period. It returns true
// if the time is the same as or after the start of the period, and before the
// end of the period.
func (p Period) Contains(t time.Time) bool {
	return SameOrBefore(p.Start, t) && p.End.After(t)
}

// ContainsInclusive checks if a given time is within the period, including the
// start and end times. It returns true if the time falls within or exactly on
// the start and end times of the period, otherwise it returns false.
func (p Period) ContainsInclusive(t time.Time) bool {
	return SameOrBefore(p.Start, t) && SameOrAfter(p.End, t)
}

// OverlapsWith returns whether p and p2 overlap.
func (p Period) OverlapsWith(p2 Period) bool {
	return p.OverlapsWithStep(time.Nanosecond, p2)
}

// OverlapsWithStep returns whether p and p2 overlap. The step parameter defines
// the minimum duration the two periods must overlap for them to be considered
// overlapping. For example, if the step is 1 hour, the periods must overlap for
// at least 1 hour to be considered overlapping. A step of 0 would consider the
// following two periods to be overlapping:
//
//	"2020-01-01 00:00:00 -> 2020-01-02 00:00:00"
//	"2020-01-02 00:00:00 -> 2020-01-03 00:00:00"
//
// [OverlapsWith] is equivalent to OverlapsWithStep with a step of 1 nanosecond.
func (p Period) OverlapsWithStep(step time.Duration, p2 Period) bool {
	if p.IsZero() || p2.IsZero() {
		return false
	}

	step = absoluteStep(step)
	pEnd := p.End.Add(-step)
	p2End := p2.End.Add(-step)

	return BetweenInclusive(p.Start, p2.Start, p2End) ||
		BetweenInclusive(pEnd, p2.Start, p2End) ||
		BetweenInclusive(p2.Start, p.Start, pEnd) ||
		BetweenInclusive(p2End, p.Start, pEnd)
}

// Years returns a slice of integers representing the years that fall within the
// period. It calculates this based on the start and end dates of the period.
// The function includes a year in the result if any part of that year is within
// the period.
func (p Period) Years() []int {
	return p.YearsStep(time.Nanosecond)
}

// YearsStep returns the years of the period. The step defines the minimum
// duration the period must be in a year for it to be included in the result.
// For example, if the step is 1 hour, the period must be in a year for at least
// 1 hour to be included in the result.
//
// A step of 0 would consider the following period to be in the years [2020, 2021]:
//
//	"2020-12-31 00:00:00 -> 2021-01-01 00:00:00"
//
// A step of 1 nanosecond would consider the following period to only be in the year 2020:
//
//	"2020-12-31 00:00:00 -> 2021-01-01 00:00:00"
func (p Period) YearsStep(step time.Duration) []int {
	step = absoluteStep(step)
	min := p.Start.Year()
	max := p.End.Add(-step).Year()

	if min > max {
		min, max = max, min
	}

	if min == max {
		return []int{min}
	}

	out := make([]int, (max-min)+1)
	for i := range out {
		out[i] = min + i
	}

	return out
}

// InYear checks if the period falls within the specified year. It returns true
// if the period is at least for a nanosecond in the given year, otherwise it
// returns false. The year is determined by using the start and end times of the
// period.
func (p Period) InYear(year int) bool {
	return p.InYearStep(time.Nanosecond, year)
}

// InYearStep checks if the period occurs within a specified year, given a step
// duration. The step duration is the minimum time the period must be within the
// year for it to be considered as occurring within that year. The function
// returns true if the period occurs within the specified year, otherwise it
// returns false.
func (p Period) InYearStep(step time.Duration, year int) bool {
	return slices.Contains(p.YearsStep(step), year)
}

// Dates returns a slice of all dates within the period, from the start date to
// the end date. Each element in the returned slice represents a single day
// within the period. The start and end times of the period are included in this
// range. If the period is not valid, it returns nil.
func (p Period) Dates() []time.Time {
	return p.DatesStep(time.Nanosecond)
}

// DatesStep returns a slice of dates within the defined time period. The step
// duration argument defines the minimum duration that must pass within a day
// for that day to be considered part of the period. For instance, if step is 1
// hour, at least 1 hour must pass within a day for that day to be included in
// the period. The function validates the period before attempting to generate
// the dates. If the validation fails, it returns nil. Otherwise, it generates
// and returns an array of time.Time values representing each date in the period
// at which the specified step duration passes.
func (p Period) DatesStep(step time.Duration) []time.Time {
	if err := p.Validate(); err != nil {
		return nil
	}

	var out []time.Time
	step = absoluteStep(step)
	end := p.End.Add(-step)
	current := StartOfDay(p.Start)

	for {
		out = append(out, current)
		current = StartOfDay(current.AddDate(0, 0, 1))
		if current.After(end) {
			break
		}
	}

	return out
}

// SliceDates divides the period into two at the first date for which the
// provided function returns true. The function is applied to each date within
// the period, in chronological order. If such a date is found, it returns two
// new periods: one before and one after that date, and a boolean value
// indicating that a split was made. If no such date is found, it returns the
// original period, an empty period, and false.
func (p Period) SliceDates(fn func(date time.Time, i int) bool) (Period, Period, bool) {
	return p.SliceDatesStep(time.Nanosecond, fn)
}

// SliceDatesStep takes a step duration and a function as arguments. It attempts
// to divide the period into two at the first date where the provided function
// returns true. The function is run on each date in the period, starting from
// the start date, with each date and its index as arguments. If such a date is
// found, it splits the period into two at that date, returning these two
// periods and a boolean value of true. If not, it returns the original period,
// an empty period, and false. The step argument defines the minimum duration
// that must pass within a day for that day to be considered part of the period.
// For instance, if step is 1 hour, at least 1 hour must pass within a day for
// that day to be included in the period.
func (p Period) SliceDatesStep(step time.Duration, fn func(date time.Time, i int) bool) (before Period, after Period, found bool) {
	if err := p.Validate(); err != nil {
		return p, Period{}, false
	}

	dates := p.DatesStep(step)

	// This can never happen because the we know the period is valid.
	if len(dates) == 0 {
		return p, Period{}, false
	}

	var cutDates []time.Time
	for i, date := range dates {
		if fn(date, i) {
			if i < len(dates) {
				after.Start = date
				after.End = p.End
			}
			break
		}
		cutDates = append(cutDates, date)
	}

	if len(cutDates) == 0 {
		return p, Period{}, false
	}

	before.Start = cutDates[0]
	before.End = cutDates[len(cutDates)-1]
	found = true

	return
}

// Cut takes a list of periods and removes them from the original period,
// returning the remaining periods. If a period in the list overlaps with the
// original period, it will be subtracted from it, potentially splitting the
// original period into two or more smaller periods. If a period in the list
// does not overlap with the original period, it is ignored. The function
// handles multiple overlapping periods and adjusts the original period
// accordingly. Before performing these operations, it sorts the periods in
// ascending order based on their start times to ensure a consistent result.
func (p Period) Cut(cut ...Period) []Period {
	slices.SortFunc(cut, func(a, b Period) int {
		if a.Start.Before(b.Start) {
			return -1
		}
		if a.Start.After(b.Start) {
			return 1
		}
		return 0
	})

	remaining := []Period{p}

	for _, c := range cut {
		newRemaining := make([]Period, 0, len(remaining))

		for _, r := range remaining {
			if cutted, ok := r.cut(c); ok {
				newRemaining = append(newRemaining, cutted...)
				continue
			}
			newRemaining = append(newRemaining, r)
		}

		remaining = newRemaining
	}

	return remaining
}

func (p Period) cut(cut Period) ([]Period, bool) {
	cutStartZero := cut.Start.IsZero()
	cutEndZero := cut.End.IsZero()

	if !cutStartZero && SameOrBefore(cut.Start, p.Start) && SameOrAfter(cut.End, p.End) {
		return nil, true
	}

	if !cutEndZero && SameOrAfter(cut.End, p.End) && SameOrBefore(cut.Start, p.Start) {
		return nil, true
	}

	if (!cutStartZero && p.End.Before(cut.Start)) || (!cutEndZero && p.Start.After(cut.End)) {
		return []Period{p}, false
	}

	beforeStart := !cutStartZero && p.Start.Before(cut.Start)
	afterEnd := !cutEndZero && p.End.After(cut.End)

	if beforeStart && afterEnd {
		return []Period{
			{Start: p.Start, End: cut.Start},
			{Start: cut.End, End: p.End},
		}, true
	}

	if beforeStart {
		return []Period{{Start: p.Start, End: cut.Start}}, true
	}

	if afterEnd {
		return []Period{{Start: cut.End, End: p.End}}, true
	}

	return []Period{p}, false
}

// CutInclusive removes the given periods from the original [Period] and returns
// the remaining periods. In contrast to the Cut function, this function
// considers periods that either start at the end time or end at the start time
// of the original period as overlapping. These overlapping periods are then
// removed from the original period. The function ensures that the end times of
// the resulting periods are exclusive by subtracting a nanosecond.
func (p Period) CutInclusive(cut ...Period) []Period {
	periodEndZero := p.End.IsZero()

	if !periodEndZero {
		p.End = p.End.Add(time.Nanosecond)
	}

	cut = slice.Map(cut, func(p Period) Period {
		if !p.End.IsZero() {
			p.End = p.End.Add(time.Nanosecond)
		}
		return p
	})

	result := p.Cut(cut...)

	if !periodEndZero {
		result = slice.Map(result, func(p Period) Period {
			p.End = p.End.Add(-time.Nanosecond)
			return p
		})
	}

	return result
}

func absoluteStep(step time.Duration) time.Duration {
	return time.Duration(math.Abs(float64(step)))
}
