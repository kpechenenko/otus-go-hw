package repository

import (
	"context"
	"github.com/kpechenenko/hw12_13_14_15_calendar/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestInMemoryEventRepository_SaveEvent(t *testing.T) {
	dur := time.Duration(300)
	now := time.Now()
	saveParamsList := []CreateEventParams{
		{Title: "event1", Date: now, Duration: dur, OwnerUserID: 1},
		{Title: "event2", Date: now, Duration: dur, OwnerUserID: 2, NotifyTime: &dur},
	}
	wantedMemory := make(map[model.EventID]model.Event)
	repository := NewInMemoryEventRepository()
	t.Run("insert multiple events", func(t *testing.T) {
		for _, p := range saveParamsList {
			ctx := context.Background()
			eventID, err := repository.CreateEvent(ctx, p)
			assert.NoError(t, err)
			wantedMemory[eventID] = model.Event{EventID: eventID, Title: p.Title, Date: p.Date, Duration: p.Duration, Description: p.Description,
				OwnerUserID: p.OwnerUserID, NotifyTime: p.NotifyTime}
		}
		assert.Equal(t, wantedMemory, repository.eventIDToEvent)
	})

	t.Run("insert wrong event", func(t *testing.T) {
		p := CreateEventParams{}
		ctx := context.Background()
		_, err := repository.CreateEvent(ctx, p)
		assert.ErrorContains(t, err, ErrInvalidParam.Error())
	})
}

func TestInMemoryEventRepository_UpdateEvent(t *testing.T) {
	dur := time.Duration(300)
	now := time.Now()
	id1 := model.EventID(1)
	initMemory := map[model.EventID]model.Event{
		id1: {EventID: id1, Title: "event1", Date: now, Duration: dur, OwnerUserID: 1},
	}

	repository := NewInMemoryEventRepository()
	repository.eventIDToEvent = initMemory

	newTitle := "newEvent1"
	newNow := time.Now()
	descr := "newEvent1Description"
	newDur := time.Duration(500)
	wantedMemory := map[model.EventID]model.Event{
		id1: {EventID: id1, Title: newTitle, Date: newNow, Duration: newDur, OwnerUserID: 1, Description: &descr},
	}

	p := UpdateEventParams{EventID: id1, Title: newTitle, Date: newNow, Duration: newDur, Description: &descr}

	t.Run("update existing event", func(t *testing.T) {
		ctx := context.Background()
		err := repository.UpdateEvent(ctx, p)
		assert.NoError(t, err)
		assert.Equal(t, wantedMemory, repository.eventIDToEvent)
	})

	t.Run("update non-existent event", func(t *testing.T) {
		id2 := model.EventID(2)
		p.EventID = id2
		ctx := context.Background()
		err := repository.UpdateEvent(ctx, p)
		assert.ErrorContains(t, err, ErrEventDoesNotExist.Error())
		assert.Equal(t, wantedMemory, repository.eventIDToEvent)
	})
}

func TestInMemoryEventRepository_DeleteEvent(t *testing.T) {
	id1 := model.EventID(1)

	initMemory := map[model.EventID]model.Event{
		id1: {EventID: id1, Title: "event1", Date: time.Now(), Duration: time.Duration(300), OwnerUserID: 1},
	}

	repository := NewInMemoryEventRepository()
	repository.eventIDToEvent = initMemory

	wantedMemory := map[model.EventID]model.Event{}

	t.Run("delete non-existing event", func(t *testing.T) {
		id2 := model.EventID(2)
		ctx := context.Background()
		err := repository.DeleteEvent(ctx, id2)
		assert.NoError(t, err)
		assert.Equal(t, initMemory, repository.eventIDToEvent)
	})

	t.Run("delete existing event", func(t *testing.T) {
		ctx := context.Background()
		err := repository.DeleteEvent(ctx, id1)
		assert.NoError(t, err)
		assert.Equal(t, wantedMemory, repository.eventIDToEvent)
	})

	t.Run("delete with wrong id", func(t *testing.T) {
		id3 := model.EventID(0)
		ctx := context.Background()
		err := repository.DeleteEvent(ctx, id3)
		assert.ErrorContains(t, err, ErrInvalidParam.Error())
		assert.Equal(t, wantedMemory, repository.eventIDToEvent)
	})
}

func TestInMemoryEventRepository_FindEvent(t *testing.T) {
	id1 := model.EventID(1)
	d1, err := time.Parse(DateFormat, "2024-06-05")
	assert.NoError(t, err)
	event1 := model.Event{EventID: id1, Title: "event1", Date: d1, Duration: time.Duration(300), OwnerUserID: 1}

	id2 := model.EventID(2)
	d2, err := time.Parse(DateFormat, "2024-06-10")
	assert.NoError(t, err)
	event2 := model.Event{EventID: id2, Title: "event2", Date: d2, Duration: time.Duration(600), OwnerUserID: 3}

	id3 := model.EventID(3)
	d3, err := time.Parse(DateFormat, "2024-06-12")
	assert.NoError(t, err)
	event3 := model.Event{EventID: id3, Title: "event3", Date: d3, Duration: time.Duration(1000), OwnerUserID: 10}

	id4 := model.EventID(4)
	d4, err := time.Parse(DateFormat, "2024-06-20")
	assert.NoError(t, err)
	event4 := model.Event{EventID: id4, Title: "event4", Date: d4, Duration: time.Duration(10000000), OwnerUserID: 5}

	initMemory := map[model.EventID]model.Event{
		id1: event1,
		id2: event2,
		id3: event3,
		id4: event4,
	}

	repository := NewInMemoryEventRepository()
	repository.eventIDToEvent = initMemory

	t.Run("find for empty period", func(t *testing.T) {
		bd, err := time.Parse(DateFormat, "2024-06-01")
		assert.NoError(t, err)
		ed, err := time.Parse(DateFormat, "2024-06-03")
		assert.NoError(t, err)
		p := FindEventParams{BeginDate: bd, EndDate: ed}
		ctx := context.Background()
		events, err := repository.FindEvent(ctx, p)
		assert.NoError(t, err)
		assert.NotNil(t, events)
		assert.Empty(t, events)
	})

	t.Run("find for empty day", func(t *testing.T) {
		bd, err := time.Parse(DateFormat, "2024-06-01")
		assert.NoError(t, err)
		p := FindEventParams{BeginDate: bd, EndDate: bd}
		ctx := context.Background()
		events, err := repository.FindEvent(ctx, p)
		assert.NoError(t, err)
		assert.NotNil(t, events)
		assert.Empty(t, events)
	})

	t.Run("find for day", func(t *testing.T) {
		p := FindEventParams{BeginDate: d1, EndDate: d1}
		ctx := context.Background()
		events, err := repository.FindEvent(ctx, p)
		wanted := []model.Event{event1}
		assert.NoError(t, err)
		assert.NotNil(t, events)
		assert.ElementsMatch(t, events, wanted)
	})

	t.Run("find for small period", func(t *testing.T) {
		p := FindEventParams{BeginDate: d2, EndDate: d3}
		ctx := context.Background()
		events, err := repository.FindEvent(ctx, p)
		wanted := []model.Event{event2, event3}
		assert.NoError(t, err)
		assert.NotNil(t, events)
		assert.ElementsMatch(t, events, wanted)
	})

	t.Run("find for period ", func(t *testing.T) {
		p := FindEventParams{BeginDate: d1, EndDate: d4}
		ctx := context.Background()
		events, err := repository.FindEvent(ctx, p)
		wanted := []model.Event{event1, event2, event3, event4}
		assert.NoError(t, err)
		assert.NotNil(t, events)
		assert.ElementsMatch(t, events, wanted)
	})
}
