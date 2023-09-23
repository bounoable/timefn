package timefn_test

import (
	"slices"
	"testing"
	"time"

	"github.com/bounoable/timefn"
)

func TestPeriod_OverlapsWithStep(t *testing.T) {
	jan1 := time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	jan3 := time.Date(2023, time.January, 3, 0, 0, 0, 0, time.UTC)
	jan7 := time.Date(2023, time.January, 7, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name string
		a, b timefn.Period
		step time.Duration
		want bool
	}{
		{
			name: "adjacent periods (a.End==b.Start),step=0",
			a: timefn.Period{
				Start: jan1,
				End:   jan3,
			},
			b: timefn.Period{
				Start: jan3,
				End:   jan7,
			},
			step: 0,
			want: true,
		},
		{
			name: "adjacent periods (a.End==b.Start),step=1ns",
			a: timefn.Period{
				Start: jan1,
				End:   jan3,
			},
			b: timefn.Period{
				Start: jan3,
				End:   jan7,
			},
			step: time.Nanosecond,
			want: false,
		},
		{
			name: "adjacent periods (a.End==b.Start),step=2ns",
			a: timefn.Period{
				Start: jan1,
				End:   jan3,
			},
			b: timefn.Period{
				Start: jan3,
				End:   jan7,
			},
			step: 2 * time.Nanosecond,
			want: false,
		},
		{
			name: "overlapping by 1ns,step=1ns",
			a: timefn.Period{
				Start: jan1,
				End:   jan3,
			},
			b: timefn.Period{
				Start: jan3.Add(-time.Nanosecond),
				End:   jan7,
			},
			step: time.Nanosecond,
			want: true,
		},
		{
			name: "overlapping by 1ns,step=2ns",
			a: timefn.Period{
				Start: jan1,
				End:   jan3,
			},
			b: timefn.Period{
				Start: jan3.Add(-time.Nanosecond),
				End:   jan7,
			},
			step: 2 * time.Nanosecond,
			want: false,
		},
		{
			name: "overlapping by 2ns,step=2ns",
			a: timefn.Period{
				Start: jan1,
				End:   jan3,
			},
			b: timefn.Period{
				Start: jan3.Add(-2 * time.Nanosecond),
				End:   jan7,
			},
			step: 2 * time.Nanosecond,
			want: true,
		},
		{
			name: "overlapping by 14m,step=15m",
			a: timefn.Period{
				Start: jan1,
				End:   jan3,
			},
			b: timefn.Period{
				Start: jan3.Add(-14 * time.Minute),
				End:   jan7,
			},
			step: 15 * time.Minute,
			want: false,
		},
		{
			name: "overlapping by 15m,step=15m",
			a: timefn.Period{
				Start: jan1,
				End:   jan3,
			},
			b: timefn.Period{
				Start: jan3.Add(-15 * time.Minute),
				End:   jan7,
			},
			step: 15 * time.Minute,
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			overlap := tt.a.OverlapsWithStep(tt.step, tt.b)

			if overlap && !tt.want {
				t.Errorf("%s should not overlap with %s with a step of %s", tt.a, tt.b, tt.step)
			} else if !overlap && tt.want {
				t.Errorf("%s should overlap with %s with a step of %s", tt.a, tt.b, tt.step)
			}
		})
	}
}

func TestPeriod_CutInclusive(t *testing.T) {
	jan1 := time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	jan2 := time.Date(2023, time.January, 2, 0, 0, 0, 0, time.UTC)
	jan2End := timefn.EndOfDay(jan2)
	jan3 := time.Date(2023, time.January, 3, 0, 0, 0, 0, time.UTC)
	jan3End := timefn.EndOfDay(jan3)
	jan4 := time.Date(2023, time.January, 4, 0, 0, 0, 0, time.UTC)
	jan4End := timefn.EndOfDay(jan4)
	jan5 := time.Date(2023, time.January, 5, 0, 0, 0, 0, time.UTC)
	jan6 := time.Date(2023, time.January, 6, 0, 0, 0, 0, time.UTC)
	jan6End := timefn.EndOfDay(jan6)
	jan7 := time.Date(2023, time.January, 7, 0, 0, 0, 0, time.UTC)
	jan7End := timefn.EndOfDay(jan7)

	tests := []struct {
		name    string
		period  timefn.Period
		cut     timefn.Period
		want    []timefn.Period
		wantCut bool
	}{
		{
			name:    "empty period + empty cut",
			period:  timefn.Period{},
			cut:     timefn.Period{},
			want:    []timefn.Period{{}},
			wantCut: false,
		},
		{
			name: "partial overlap with start",
			period: timefn.Period{
				Start: jan1,
				End:   jan7End,
			},
			cut: timefn.Period{
				End: jan4End,
			},
			want: []timefn.Period{{
				Start: jan5,
				End:   jan7End,
			}},
			wantCut: true,
		},
		{
			name: "partial overlap with end",
			period: timefn.Period{
				Start: jan1,
				End:   jan7End,
			},
			cut: timefn.Period{
				Start: jan4,
			},
			want: []timefn.Period{{
				Start: jan1,
				End:   jan3End,
			}},
			wantCut: true,
		},
		{
			name: "total overlap",
			period: timefn.Period{
				Start: jan1,
				End:   jan7End,
			},
			cut: timefn.Period{
				Start: jan3,
				End:   jan6End,
			},
			want: []timefn.Period{
				{
					Start: jan1,
					End:   jan2End,
				},
				{
					Start: jan7,
					End:   jan7End,
				},
			},
			wantCut: true,
		},
		{
			name: "exact cut",
			period: timefn.Period{
				Start: jan1,
				End:   jan7End,
			},
			cut: timefn.Period{
				Start: jan1,
				End:   jan7End,
			},
			want:    nil,
			wantCut: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.period.CutInclusive(tt.cut)

			if !slices.Equal(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeriod_Cut(t *testing.T) {
	jan1 := time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	jan3 := time.Date(2023, time.January, 3, 0, 0, 0, 0, time.UTC)
	jan4 := time.Date(2023, time.January, 4, 0, 0, 0, 0, time.UTC)
	jan6 := time.Date(2023, time.January, 6, 0, 0, 0, 0, time.UTC)
	jan7 := time.Date(2023, time.January, 7, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		period  timefn.Period
		cut     timefn.Period
		want    []timefn.Period
		wantCut bool
	}{
		{
			name:    "empty period + empty cut",
			period:  timefn.Period{},
			cut:     timefn.Period{},
			want:    []timefn.Period{{}},
			wantCut: false,
		},
		{
			name: "partial overlap with start",
			period: timefn.Period{
				Start: jan1,
				End:   jan7,
			},
			cut: timefn.Period{
				End: jan4,
			},
			want: []timefn.Period{{
				Start: jan4,
				End:   jan7,
			}},
			wantCut: true,
		},
		{
			name: "partial overlap with end",
			period: timefn.Period{
				Start: jan1,
				End:   jan7,
			},
			cut: timefn.Period{
				Start: jan4,
			},
			want: []timefn.Period{{
				Start: jan1,
				End:   jan4,
			}},
			wantCut: true,
		},
		{
			name: "total overlap",
			period: timefn.Period{
				Start: jan1,
				End:   jan7,
			},
			cut: timefn.Period{
				Start: jan3,
				End:   jan6,
			},
			want: []timefn.Period{
				{
					Start: jan1,
					End:   jan3,
				},
				{
					Start: jan6,
					End:   jan7,
				},
			},
			wantCut: true,
		},
		{
			name: "exact cut",
			period: timefn.Period{
				Start: jan1,
				End:   jan7,
			},
			cut: timefn.Period{
				Start: jan1,
				End:   jan7,
			},
			want:    nil,
			wantCut: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.period.Cut(tt.cut)

			if !slices.Equal(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeriod_YearsStep(t *testing.T) {
	tests := []struct {
		period timefn.Period
		step   time.Duration
		want   []int
	}{
		{
			period: timefn.Period{
				Start: time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			step: 0,
			want: []int{2020, 2021, 2022, 2023},
		},
		{
			period: timefn.Period{
				Start: time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			step: time.Nanosecond,
			want: []int{2020, 2021, 2022},
		},
		{
			period: timefn.Period{
				Start: time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2023, time.January, 6, 0, 0, 0, 0, time.UTC),
			},
			step: time.Nanosecond,
			want: []int{2020, 2021, 2022, 2023},
		},
		{
			period: timefn.Period{
				Start: time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			step: time.Minute,
			want: []int{2020, 2021, 2022},
		},
		{
			period: timefn.Period{
				Start: time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2023, time.January, 1, 1, 0, 0, 0, time.UTC),
			},
			step: time.Hour,
			want: []int{2020, 2021, 2022, 2023},
		},
	}

	for _, tt := range tests {
		years := tt.period.YearsStep(tt.step)

		if !slices.Equal(years, tt.want) {
			t.Errorf("expected years of %s with step %s to be %v; got %v", tt.period, tt.step, tt.want, years)
		}
	}
}

func TestPeriod_DatesStep(t *testing.T) {
	tests := []struct {
		period timefn.Period
		step   time.Duration
		want   []time.Time
	}{
		{
			period: timefn.Period{
				Start: time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2020, time.January, 5, 0, 0, 0, 0, time.UTC),
			},
			step: 0,
			want: []time.Time{
				time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2020, time.January, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2020, time.January, 3, 0, 0, 0, 0, time.UTC),
				time.Date(2020, time.January, 4, 0, 0, 0, 0, time.UTC),
				time.Date(2020, time.January, 5, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			period: timefn.Period{
				Start: time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2020, time.January, 5, 0, 0, 0, 0, time.UTC),
			},
			step: time.Nanosecond,
			want: []time.Time{
				time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2020, time.January, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2020, time.January, 3, 0, 0, 0, 0, time.UTC),
				time.Date(2020, time.January, 4, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			period: timefn.Period{
				Start: time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2020, time.January, 5, 1, 0, 0, 0, time.UTC),
			},
			step: time.Hour,
			want: []time.Time{
				time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2020, time.January, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2020, time.January, 3, 0, 0, 0, 0, time.UTC),
				time.Date(2020, time.January, 4, 0, 0, 0, 0, time.UTC),
				time.Date(2020, time.January, 5, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, tt := range tests {
		dates := tt.period.DatesStep(tt.step)

		if !slices.Equal(dates, tt.want) {
			t.Errorf("expected dates of period to be %v; got %v", tt.want, dates)
		}
	}
}
