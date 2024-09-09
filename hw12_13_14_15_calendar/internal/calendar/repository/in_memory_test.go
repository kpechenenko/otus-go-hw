package repository

import (
	"context"
	"testing"
	"time"

	"github.com/kpechenenko/hw12_13_14_15_calendar/calendar/internal/calendar/model"
	"github.com/stretchr/testify/assert"
)

func TestInMemory_AddEvent(t *testing.T) {
	dur := time.Duration(300)
	now := time.Now()
	saveParamsList := []AddEventParams{
		{Title: "event1", Date: now, Duration: dur, OwnerID: 1},
		{Title: "event2", Date: now, Duration: dur, OwnerID: 2, NotifyFor: &dur},
	}
	wantedMemory := make(map[model.EventID]model.Event)
	srcRepo := NewInMemory()
	repo, ok := srcRepo.(*inMemory)
	assert.True(t, ok)

	t.Run("insert multiple events", func(t *testing.T) {
		for _, p := range saveParamsList {
			ctx := context.Background()
			eventID, err := repo.AddEvent(ctx, p)
			assert.NoError(t, err)
			wantedMemory[eventID] = model.Event{
				ID:          eventID,
				Title:       p.Title,
				Date:        p.Date,
				Duration:    p.Duration,
				Description: p.Description,
				OwnerID:     p.OwnerID,
				NotifyFor:   p.NotifyFor,
			}
		}
		assert.EqualValues(t, wantedMemory, repo.events)
	})

	t.Run("insert wrong event", func(t *testing.T) {
		p := AddEventParams{}
		ctx := context.Background()
		_, err := repo.AddEvent(ctx, p)
		assert.ErrorContains(t, err, ErrInvalidParams.Error())
	})
}

func TestInMemory_UpdateEvent(t *testing.T) {
	dur := time.Duration(300)
	now := time.Now()
	id1 := GenerateEventID()
	event := model.Event{ID: id1, Title: "event1", Date: now, Duration: dur, OwnerID: 1}
	initMemory := map[model.EventID]model.Event{
		id1: event,
	}

	srcRepo := NewInMemory()
	repo, ok := srcRepo.(*inMemory)
	assert.True(t, ok)
	repo.events = initMemory

	newTitle := "newEvent1"
	newNow := time.Now()
	descr := "newEvent1Description"
	newDur := time.Duration(500)
	wantedMemory := map[model.EventID]model.Event{
		id1: {ID: id1, Title: newTitle, Date: newNow, Duration: newDur, OwnerID: 1, Description: &descr},
	}

	p := UpdateEventParams{ID: id1, Title: newTitle, Date: newNow, Duration: newDur, Description: &descr}

	t.Run("update existing event", func(t *testing.T) {
		ctx := context.Background()
		err := repo.UpdateEvent(ctx, p)
		assert.NoError(t, err)
		assert.EqualValues(t, wantedMemory, repo.events)
	})

	t.Run("update non-existent event", func(t *testing.T) {
		id2 := GenerateEventID()
		p.ID = id2
		ctx := context.Background()
		err := repo.UpdateEvent(ctx, p)
		assert.Nil(t, err)
		assert.EqualValues(t, wantedMemory, repo.events)
	})
}

func TestInMemory_DeleteEvent(t *testing.T) {
	id1 := GenerateEventID()
	event := model.Event{ID: id1, Title: "event1", Date: time.Now(), Duration: time.Duration(300), OwnerID: 1}
	initMemory := map[model.EventID]model.Event{id1: event}

	srcRepo := NewInMemory()
	repo, ok := srcRepo.(*inMemory)
	assert.True(t, ok)
	repo.events = initMemory

	wantedMemory := map[model.EventID]model.Event{}

	t.Run("delete non-existing event", func(t *testing.T) {
		id2 := GenerateEventID()
		ctx := context.Background()
		err := repo.DeleteEvent(ctx, id2)
		assert.NoError(t, err)
		assert.Equal(t, initMemory, repo.events)
	})

	t.Run("delete existing event", func(t *testing.T) {
		ctx := context.Background()
		err := repo.DeleteEvent(ctx, id1)
		assert.NoError(t, err)
		assert.EqualValues(t, wantedMemory, repo.events)
	})
}

func TestInMemory_FindEvent(t *testing.T) {
	id1 := GenerateEventID()
	d1, err := time.Parse(DateFormat, "2024-06-05")
	assert.NoError(t, err)
	event1 := model.Event{ID: id1, Title: "event1", Date: d1, Duration: time.Duration(300), OwnerID: 1}

	id2 := GenerateEventID()
	d2, err := time.Parse(DateFormat, "2024-06-10")
	assert.NoError(t, err)
	event2 := model.Event{ID: id2, Title: "event2", Date: d2, Duration: time.Duration(600), OwnerID: 3}

	id3 := GenerateEventID()
	d3, err := time.Parse(DateFormat, "2024-06-12")
	assert.NoError(t, err)
	event3 := model.Event{ID: id3, Title: "event3", Date: d3, Duration: time.Duration(1000), OwnerID: 10}

	id4 := GenerateEventID()
	d4, err := time.Parse(DateFormat, "2024-06-20")
	assert.NoError(t, err)
	event4 := model.Event{ID: id4, Title: "event4", Date: d4, Duration: time.Duration(10000000), OwnerID: 5}

	initMemory := map[model.EventID]model.Event{
		id1: event1,
		id2: event2,
		id3: event3,
		id4: event4,
	}

	srcRepo := NewInMemory()
	repo, ok := srcRepo.(*inMemory)
	assert.True(t, ok)
	repo.events = initMemory

	t.Run("find for empty period", func(t *testing.T) {
		bd, err := time.Parse(DateFormat, "2024-06-01")
		assert.NoError(t, err)
		ed, err := time.Parse(DateFormat, "2024-06-03")
		assert.NoError(t, err)
		p := GetEventParams{BeginDate: bd, EndDate: ed}
		ctx := context.Background()
		events, err := repo.GetEvents(ctx, p)
		assert.NoError(t, err)
		assert.NotNil(t, events)
		assert.Empty(t, events)
	})

	t.Run("find for empty day", func(t *testing.T) {
		bd, err := time.Parse(DateFormat, "2024-06-01")
		assert.NoError(t, err)
		p := GetEventParams{BeginDate: bd, EndDate: bd}
		ctx := context.Background()
		events, err := repo.GetEvents(ctx, p)
		assert.NoError(t, err)
		assert.NotNil(t, events)
		assert.Empty(t, events)
	})

	t.Run("find for day", func(t *testing.T) {
		p := GetEventParams{BeginDate: d1, EndDate: d1}
		ctx := context.Background()
		events, err := repo.GetEvents(ctx, p)
		wanted := []model.Event{event1}
		assert.NoError(t, err)
		assert.NotNil(t, events)
		assert.ElementsMatch(t, events, wanted)
	})

	t.Run("find for small period", func(t *testing.T) {
		p := GetEventParams{BeginDate: d2, EndDate: d3}
		ctx := context.Background()
		events, err := repo.GetEvents(ctx, p)
		wanted := []model.Event{event2, event3}
		assert.NoError(t, err)
		assert.NotNil(t, events)
		assert.ElementsMatch(t, events, wanted)
	})

	t.Run("find for period ", func(t *testing.T) {
		p := GetEventParams{BeginDate: d1, EndDate: d4}
		ctx := context.Background()
		events, err := repo.GetEvents(ctx, p)
		wanted := []model.Event{event1, event2, event3, event4}
		assert.NoError(t, err)
		assert.NotNil(t, events)
		assert.ElementsMatch(t, events, wanted)
	})
}
