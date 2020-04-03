package timefn_test

import (
	"testing"
	"time"

	"github.com/bounoable/timefn"
	"github.com/stretchr/testify/assert"
)

func TestStartOfSecond(t *testing.T) {
	assert.Equal(t, time.Date(2020, 3, 1, 15, 15, 15, 0, time.UTC), timefn.StartOfSecond(time.Date(2020, 3, 1, 15, 15, 15, 15, time.UTC)))
}

func TestEndOfSecond(t *testing.T) {
	assert.Equal(t, time.Date(2020, 3, 1, 15, 15, 16, 0, time.UTC).Add(-time.Nanosecond), timefn.EndOfSecond(time.Date(2020, 3, 1, 15, 15, 15, 15, time.UTC)))
}

func TestStartOfMinute(t *testing.T) {
	assert.Equal(t, time.Date(2020, 3, 1, 15, 15, 0, 0, time.UTC), timefn.StartOfMinute(time.Date(2020, 3, 1, 15, 15, 15, 15, time.UTC)))
}

func TestEndOfMinute(t *testing.T) {
	assert.Equal(t, time.Date(2020, 3, 1, 15, 16, 0, 0, time.UTC).Add(-time.Nanosecond), timefn.EndOfMinute(time.Date(2020, 3, 1, 15, 15, 15, 15, time.UTC)))
}

func TestStartOfHour(t *testing.T) {
	assert.Equal(t, time.Date(2020, 3, 1, 15, 0, 0, 0, time.UTC), timefn.StartOfHour(time.Date(2020, 3, 1, 15, 15, 15, 15, time.UTC)))
}

func TestEndOfHour(t *testing.T) {
	assert.Equal(t, time.Date(2020, 3, 1, 16, 0, 0, 0, time.UTC).Add(-time.Nanosecond), timefn.EndOfHour(time.Date(2020, 3, 1, 15, 15, 15, 15, time.UTC)))
}

func TestStartOfDay(t *testing.T) {
	assert.Equal(t, time.Date(2020, 3, 1, 0, 0, 0, 0, time.UTC), timefn.StartOfDay(time.Date(2020, 3, 1, 15, 15, 15, 15, time.UTC)))
}

func TestEndOfDay(t *testing.T) {
	assert.Equal(t, time.Date(2020, 3, 2, 0, 0, 0, 0, time.UTC).Add(-time.Nanosecond), timefn.EndOfDay(time.Date(2020, 3, 1, 15, 15, 15, 15, time.UTC)))
}

func TestStartOfWeek(t *testing.T) {
	tests := []struct {
		Time     time.Time
		Expected time.Time
	}{
		{
			Time:     time.Date(2020, 3, 30, 15, 15, 15, 15, time.UTC),
			Expected: time.Date(2020, 3, 29, 0, 0, 0, 0, time.UTC),
		},
		{
			Time:     time.Date(2020, 3, 1, 15, 15, 15, 15, time.UTC),
			Expected: time.Date(2020, 3, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			Time:     time.Date(2020, 4, 8, 15, 15, 15, 15, time.UTC),
			Expected: time.Date(2020, 4, 5, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.Expected, timefn.StartOfWeek(test.Time))
	}
}

func TestEndOfWeek(t *testing.T) {
	tests := []struct {
		Time     time.Time
		Expected time.Time
	}{
		{
			Time:     time.Date(2020, 3, 30, 15, 15, 15, 15, time.UTC),
			Expected: time.Date(2020, 4, 5, 0, 0, 0, 0, time.UTC).Add(-time.Nanosecond),
		},
		{
			Time:     time.Date(2020, 3, 1, 15, 15, 15, 15, time.UTC),
			Expected: time.Date(2020, 3, 8, 0, 0, 0, 0, time.UTC).Add(-time.Nanosecond),
		},
		{
			Time:     time.Date(2020, 4, 8, 15, 15, 15, 15, time.UTC),
			Expected: time.Date(2020, 4, 12, 0, 0, 0, 0, time.UTC).Add(-time.Nanosecond),
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.Expected, timefn.EndOfWeek(test.Time))
	}
}

func TestStartOfISOWeek(t *testing.T) {
	tests := []struct {
		Time     time.Time
		Expected time.Time
	}{
		{
			Time:     time.Date(2020, 3, 30, 15, 15, 15, 15, time.UTC),
			Expected: time.Date(2020, 3, 30, 0, 0, 0, 0, time.UTC),
		},
		{
			Time:     time.Date(2020, 3, 1, 15, 15, 15, 15, time.UTC),
			Expected: time.Date(2020, 2, 24, 0, 0, 0, 0, time.UTC),
		},
		{
			Time:     time.Date(2020, 4, 8, 15, 15, 15, 15, time.UTC),
			Expected: time.Date(2020, 4, 6, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.Expected, timefn.StartOfISOWeek(test.Time))
	}
}

func TestEndOfISOWeek(t *testing.T) {
	tests := []struct {
		Time     time.Time
		Expected time.Time
	}{
		{
			Time:     time.Date(2020, 3, 30, 15, 15, 15, 15, time.UTC),
			Expected: time.Date(2020, 4, 6, 0, 0, 0, 0, time.UTC).Add(-time.Nanosecond),
		},
		{
			Time:     time.Date(2020, 3, 1, 15, 15, 15, 15, time.UTC),
			Expected: time.Date(2020, 3, 2, 0, 0, 0, 0, time.UTC).Add(-time.Nanosecond),
		},
		{
			Time:     time.Date(2020, 4, 8, 15, 15, 15, 15, time.UTC),
			Expected: time.Date(2020, 4, 13, 0, 0, 0, 0, time.UTC).Add(-time.Nanosecond),
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.Expected, timefn.EndOfISOWeek(test.Time))
	}
}

func TestStartOfMonth(t *testing.T) {
	assert.Equal(t, time.Date(2020, 3, 1, 0, 0, 0, 0, time.UTC), timefn.StartOfMonth(time.Date(2020, 3, 15, 15, 15, 15, 15, time.UTC)))
}

func TestEndOfMonth(t *testing.T) {
	assert.Equal(t, time.Date(2020, 4, 1, 0, 0, 0, 0, time.UTC).Add(-time.Nanosecond), timefn.EndOfMonth(time.Date(2020, 3, 15, 15, 15, 15, 15, time.UTC)))
}

func TestStartOfYear(t *testing.T) {
	assert.Equal(t, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), timefn.StartOfYear(time.Date(2020, 3, 15, 15, 15, 15, 15, time.UTC)))
}

func TestEndOfYear(t *testing.T) {
	assert.Equal(t, time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC).Add(-time.Nanosecond), timefn.EndOfYear(time.Date(2020, 3, 15, 15, 15, 15, 15, time.UTC)))
}

func TestBetween(t *testing.T) {
	tests := []struct {
		Time     time.Time
		Left     time.Time
		Right    time.Time
		Expected bool
	}{
		{
			Time:     time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			Left:     time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			Right:    time.Date(2020, 1, 10, 0, 0, 0, 0, time.UTC),
			Expected: false,
		},
		{
			Time:     time.Date(2020, 1, 10, 0, 0, 0, 0, time.UTC),
			Left:     time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			Right:    time.Date(2020, 1, 10, 0, 0, 0, 0, time.UTC),
			Expected: false,
		},
		{
			Time:     time.Date(2020, 1, 1, 0, 0, 0, 1, time.UTC),
			Left:     time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			Right:    time.Date(2020, 1, 10, 0, 0, 0, 0, time.UTC),
			Expected: true,
		},
		{
			Time:     time.Date(2020, 1, 10, 0, 0, 0, 0, time.UTC).Add(-time.Nanosecond),
			Left:     time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			Right:    time.Date(2020, 1, 10, 0, 0, 0, 0, time.UTC),
			Expected: true,
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.Expected, timefn.Between(test.Time, test.Left, test.Right))
	}
}

func TestBetweenInclusive(t *testing.T) {
	tests := []struct {
		Time     time.Time
		Left     time.Time
		Right    time.Time
		Expected bool
	}{
		{
			Time:     time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			Left:     time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			Right:    time.Date(2020, 1, 10, 0, 0, 0, 0, time.UTC),
			Expected: true,
		},
		{
			Time:     time.Date(2020, 1, 10, 0, 0, 0, 0, time.UTC),
			Left:     time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			Right:    time.Date(2020, 1, 10, 0, 0, 0, 0, time.UTC),
			Expected: true,
		},
		{
			Time:     time.Date(2020, 1, 1, 0, 0, 0, 1, time.UTC),
			Left:     time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			Right:    time.Date(2020, 1, 10, 0, 0, 0, 0, time.UTC),
			Expected: true,
		},
		{
			Time:     time.Date(2020, 1, 10, 0, 0, 0, 0, time.UTC).Add(-time.Nanosecond),
			Left:     time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			Right:    time.Date(2020, 1, 10, 0, 0, 0, 0, time.UTC),
			Expected: true,
		},
		{
			Time:     time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Add(-time.Nanosecond),
			Left:     time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			Right:    time.Date(2020, 1, 10, 0, 0, 0, 0, time.UTC),
			Expected: false,
		},
		{
			Time:     time.Date(2020, 1, 10, 0, 0, 0, 0, time.UTC).Add(time.Nanosecond),
			Left:     time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			Right:    time.Date(2020, 1, 10, 0, 0, 0, 0, time.UTC),
			Expected: false,
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.Expected, timefn.BetweenInclusive(test.Time, test.Left, test.Right))
	}
}

func TestSameOrBefore(t *testing.T) {
	tests := []struct {
		Time     time.Time
		Right    time.Time
		Expected bool
	}{
		{
			Time:     time.Date(2020, 1, 9, 0, 0, 0, 0, time.UTC),
			Right:    time.Date(2020, 1, 10, 0, 0, 0, 0, time.UTC),
			Expected: true,
		},
		{
			Time:     time.Date(2020, 1, 10, 0, 0, 0, 0, time.UTC),
			Right:    time.Date(2020, 1, 10, 0, 0, 0, 0, time.UTC),
			Expected: true,
		},
		{
			Time:     time.Date(2020, 1, 10, 0, 0, 0, 1, time.UTC),
			Right:    time.Date(2020, 1, 10, 0, 0, 0, 0, time.UTC),
			Expected: false,
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.Expected, timefn.SameOrBefore(test.Time, test.Right))
	}
}

func TestSameOrAfter(t *testing.T) {
	tests := []struct {
		Time     time.Time
		Left     time.Time
		Expected bool
	}{
		{
			Time:     time.Date(2020, 1, 11, 0, 0, 0, 0, time.UTC),
			Left:     time.Date(2020, 1, 10, 0, 0, 0, 0, time.UTC),
			Expected: true,
		},
		{
			Time:     time.Date(2020, 1, 10, 0, 0, 0, 0, time.UTC),
			Left:     time.Date(2020, 1, 10, 0, 0, 0, 0, time.UTC),
			Expected: true,
		},
		{
			Time:     time.Date(2020, 1, 10, 0, 0, 0, 0, time.UTC).Add(-time.Nanosecond),
			Left:     time.Date(2020, 1, 10, 0, 0, 0, 0, time.UTC),
			Expected: false,
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.Expected, timefn.SameOrAfter(test.Time, test.Left))
	}
}
