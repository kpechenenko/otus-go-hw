package service

import (
	"github.com/kpechenenko/hw12_13_14_15_calendar/internal/repository"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetLastDayOfWeek(t *testing.T) {
	sunday, err := time.Parse(repository.DateFormat, "2024-06-09")
	assert.NoError(t, err)
	monday, err := time.Parse(repository.DateFormat, "2024-06-03")
	assert.NoError(t, err)
	wednesday, err := time.Parse(repository.DateFormat, "2024-06-05")
	assert.NoError(t, err)
	tests := []struct {
		name          string
		day           time.Time
		wantedLastDay time.Time
	}{
		{
			"check last for sunday",
			sunday,
			sunday,
		},
		{
			"check last for monday",
			monday,
			sunday,
		},
		{
			"check last for wednesday",
			wednesday,
			sunday,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.wantedLastDay, getLastDayOfWeek(tc.day))
		})
	}
}

func TestGetLastDayOfMonth(t *testing.T) {
	middleOfJune, err := time.Parse(repository.DateFormat, "2024-06-15")
	assert.NoError(t, err)
	endOfJune, err := time.Parse(repository.DateFormat, "2024-06-30")
	assert.NoError(t, err)
	beginOfFebruary, err := time.Parse(repository.DateFormat, "2024-02-01")
	assert.NoError(t, err)
	endOfFebruary, err := time.Parse(repository.DateFormat, "2024-02-29")
	assert.NoError(t, err)
	endOfDecember, err := time.Parse(repository.DateFormat, "2023-12-31")
	assert.NoError(t, err)
	tests := []struct {
		name          string
		day           time.Time
		wantedLastDay time.Time
	}{
		{
			"check middleOfJune",
			middleOfJune,
			endOfJune,
		},
		{
			"check beginOfFebruary",
			beginOfFebruary,
			endOfFebruary,
		},
		{
			"check endOfDecember",
			endOfDecember,
			endOfDecember,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.wantedLastDay, getLastDayOfMonth(tc.day))
		})
	}
}
