package timefn

import (
	"fmt"
	"math"
	"slices"
	"sort"
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
// A step of 1 nanosecond would consider the following period to be only in the year 2020:
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

// Dates retrieves all the dates within the period, returning a slice of
// [time.Time]. Each date is represented by the start of the day, and they are
// returned in chronological order from the start to the end of the period. If
// the period is invalid, it returns nil.
func (p Period) Dates() []time.Time {
	return p.DatesStep(time.Nanosecond)
}

// DatesStep iterates over each date within the period, using a specified step
// interval. It generates a slice of [time.Time] representing each date from the
// start to the end of the period, not including the last date if it is equal to
// the end date minus the step interval. The step defines the minimum duration
// to advance when moving to the next date within the period. If any part of the
// period is invalid, it returns nil.
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

// SliceDates divides the [Period] into two periods based on a user-defined
// criterion. It iterates over each date within the period, invoking a callback
// function with the current date and its index. The callback's return value
// determines where the slicing occurs; if it returns true, slicing is performed
// at that date. The function returns two [Period]s representing the time spans
// before and after the slicing point, and a boolean indicating whether the
// slicing was successful. If no date satisfies the criterion, or if the
// [Period] is invalid, the original [Period] is returned as the first result,
// with an empty second [Period] and false for the boolean.
func (p Period) SliceDates(fn func(date time.Time, i int) bool) (Period, Period, bool) {
	return p.SliceDatesStep(time.Nanosecond, fn)
}

// SliceDatesStep divides the [Period] into two at a date determined by the
// provided callback function, which is called for each date in the period, with
// an additional step interval between dates. It returns two [Period]s: one
// before and one after the date where the callback returns true, along with a
// boolean indicating if such a date was found. If the period is invalid or no
// date satisfies the callback, it returns the original [Period], an empty
// [Period], and false.
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

// Cut removes specified periods from the receiver [Period] and returns a slice
// of the remaining [Period]s. This operation is non-destructive to the original
// [Period]. If no periods are specified for removal or if none of the specified
// periods intersect with the receiver, the result will contain the original
// [Period] unaltered. If an intersection occurs, the function returns a new set
// of [Period]s that represent the time spans before and after each
// intersection, effectively "cutting out" the intersecting ranges. The
// resulting slice is sorted by the start times of each [Period].
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

// CutInclusive trims the specified periods from the receiver [Period] and
// returns a slice of [Period]s that represent the remaining time ranges. It
// does so in an inclusive manner, where the end times of both the receiver and
// the specified periods are considered part of the cut. If a specified period
// to cut overlaps with or is within the bounds of the receiver period, it is
// trimmed accordingly, and the remaining non-overlapping parts are returned. If
// no overlap exists, the original [Period] is returned unchanged. The resulting
// slice of [Period]s is sorted by their start times.
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

// MergeStep merges the [Period] with a slice of other periods, ensuring that
// any overlapping periods are combined into continuous periods based on a
// specified minimum duration step. It returns a slice of merged periods, sorted
// by their start times. If no additional periods are provided, the result is a
// slice containing only the original period. The step parameter determines how
// much overlap is necessary for two periods to be considered as one continuous
// period.
func (p Period) MergeStep(step time.Duration, periods []Period) []Period {
	if len(periods) == 0 {
		return []Period{p}
	}

	periods = append([]Period{p}, periods...)

	sort.Slice(periods, func(i, j int) bool {
		return periods[i].Start.Before(periods[j].Start)
	})

	merged := []Period{periods[0]}

	for _, p := range periods[1:] {
		last := &merged[len(merged)-1]

		if last.OverlapsWithStep(step, p) {
			last.End = maxTime(last.End, p.End)
		} else if SameOrBefore(last.End, p.Start) {
			merged = append(merged, p)
		}
	}

	return merged
}

func maxTime(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

// Merge combines the receiver [Period] with a slice of other [Period]s and
// returns a new slice of merged [Period]s. Overlapping periods are consolidated
// into single periods, while non-overlapping periods remain separate. The merge
// process respects the chronological order of periods.
func (p Period) Merge(periods []Period) []Period {
	return p.MergeStep(0, periods)
}

func absoluteStep(step time.Duration) time.Duration {
	return time.Duration(math.Abs(float64(step)))
}
