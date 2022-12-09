package timefn

import (
	"fmt"
	"math"
	"strings"
	"text/template"
	"time"

	"github.com/bounoable/timefn/internal/slice"
	"golang.org/x/exp/slices"
)

// DefaultPeriodFormat is the default format string used by [Period.Format].
var DefaultPeriodFormat = "{{ .Start }} -> {{ .End }}"

// Period is an open-closed time period, i.e. the period includes the start
// time, but excludes the end time. A period with a start time at
// 2020-01-01 00:00:00 and end time at 2020-01-02 00:00:00 represents the date
// 2020-01-01 from 00:00:00 to 23:59:59 (1 nanosecond before the next day).
type Period struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

func (p Period) String() string {
	return p.Format()
}

// Format formats the period to a string using the [DefaultPeriodFormat].
func (p Period) Format() string {
	return p.FormatAs(DefaultPeriodFormat)
}

// FormatAs formats the period to a string using the provided format string.
// The format string will be compiled as a template that has access to the
// period:
//
//	"{{ .Start }} -> {{ .End }}"
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

// IsZero returns whether the period's start and end date are both zero-value
// as determined by [time.Time.IsZero()].
func (p Period) IsZero() bool {
	return p.Start.IsZero() && p.End.IsZero()
}

// Validate validates the time period. Start and End must be non-zero, and End
// must be after Start.
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

// Add returns a new period with the start and end time of p shifted by d.
func (p Period) Add(d time.Duration) Period {
	return Period{
		Start: p.Start.Add(d),
		End:   p.End.Add(d),
	}
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

// Years returns the years of the period. Years returns [Period.YearsStep] with
// a step of 1 nanosecond, meaning that the period must be in a year for at
// least 1 nanosecond for the year to be included in the result.
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

// InYear returns whether the period includes a time that is in the given year.
// InYear is equivalent to [Period.InYearStep] with a step of 1 nanosecond.
func (p Period) InYear(year int) bool {
	return p.InYearStep(time.Nanosecond, year)
}

// InYearStep returns whether the period includes a time that is in the given year.
// The step parameter is passed to [Period.YearsStep] to extract the years of the period.
func (p Period) InYearStep(step time.Duration, year int) bool {
	return slices.Contains(p.YearsStep(step), year)
}

// Dates returns the dates in the period. Dates returns [Period.DatesStep] with
// a step of 1 nanosecond, meaning that the period must be in a day for at least
// 1 nanosecond for the day to be included in the result.
func (p Period) Dates() []time.Time {
	return p.DatesStep(time.Nanosecond)
}

// DatesStep returns the the dates in the period. The step defines the minimum
// duration the period must be in a day for the day to be included in the result.
// For example, if the step is 1 hour, the period must be in a day for at least
// 1 hour to be included in the result.
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
		current = StartOfDay(current.Add(24 * time.Hour))
		if current.After(end) {
			break
		}
	}

	return out
}

// SliceDates slices `p` around the first date that satisfies the given predicate.
// SliceDates is equivalent to [Period.SliceDatesStep] with a step of 1 nanosecond.
func (p Period) SliceDates(fn func(date time.Time, i int) bool) (Period, Period, bool) {
	return p.SliceDatesStep(time.Nanosecond, fn)
}

// SliceDatesStep slices `p` around the first date that satisfies the given
// predicate. The date that satisfies the predicate and all subsequent dates
// are not included in `before`. The remaining period is returned in `after`.
// If no dates were cut from `p`, or if `p` is not a valid period,
// Cut returns `p, Period{}, false`. The step parameter is passed to
// [Period.DatesStep] to extract the dates of the period.
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

// Cut cuts periods out of `p`. Cut treats the periods as closed-open intervals,
// meaning that the start of a period is inclusive but the end is not.
//
// # Example
//
//	p := Period{
//		Start: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), // Jan 1, 2020
//		End: time.Date(2020, 1, 7, 0, 0, 0, 0, time.UTC), // Jan 7, 2020
//	}
//	cut := Period{
//		Start: time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC), // Jan 1, 2020
//		End: time.Date(2020, 1, 5, 0, 0, 0, 0, time.UTC), // Jan 5, 2020
//	}
//
//	result := p.Cut(cut)
//	// result == []Period{
//	//	{
//	//		Start: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
//	//		End: time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
//	//	},
//	//	{
//	//		Start: time.Date(2020, 1, 5, 0, 0, 0, 0, time.UTC),
//	//		End: time.Date(2020, 1, 7, 0, 0, 0, 0, time.UTC),
//	//	},
//	}
func (p Period) Cut(cut ...Period) []Period {
	slices.SortFunc(cut, func(a, b Period) bool {
		return a.Start.Before(b.Start)
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

// CutInclusive is the same as Cut, but treats the periods as closed intervals,
// meaning that the start and end of the period are both included.
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
