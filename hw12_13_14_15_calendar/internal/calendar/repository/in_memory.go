package repository

import (
	"context"
	"fmt"
	"sync"

	"github.com/kpechenenko/hw12_13_14_15_calendar/calendar/internal/calendar/model"
)

type inMemory struct {
	mu     sync.RWMutex
	events map[model.EventID]model.Event
}

func (s *inMemory) AddEvent(_ context.Context, params AddEventParams) (id model.EventID, err error) {
	if len(params.Title) == 0 || params.Date.IsZero() || params.Duration.Abs() == 0 || params.OwnerID == 0 {
		err = fmt.Errorf("%w: event title, date, ownerUserId must be provided", ErrInvalidParams)
		return
	}
	id = GenerateEventID()
	event := model.Event{
		ID:          id,
		Title:       params.Title,
		Date:        params.Date,
		Duration:    params.Duration,
		Description: params.Description,
		OwnerID:     params.OwnerID,
		NotifyFor:   params.NotifyFor,
	}
	s.mu.Lock()
	s.events[id] = event
	s.mu.Unlock()
	return
}

func (s *inMemory) UpdateEvent(_ context.Context, params UpdateEventParams) (err error) {
	if len(params.Title) == 0 || params.Date.IsZero() || params.Duration.Abs() == 0 {
		err = fmt.Errorf("%w: params id, title, date must be provided", ErrInvalidParams)
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	var event model.Event
	var ok bool
	if event, ok = s.events[params.ID]; !ok {
		return
	}
	event.Title = params.Title
	event.Date = params.Date
	event.Duration = params.Duration
	event.Description = params.Description
	event.NotifyFor = params.NotifyFor
	s.events[event.ID] = event
	return
}

func (s *inMemory) DeleteEvent(_ context.Context, id model.EventID) (err error) {
	s.mu.Lock()
	delete(s.events, id)
	s.mu.Unlock()
	return
}

func (s *inMemory) GetEvents(_ context.Context, params GetEventParams) (events []model.Event, err error) {
	if params.BeginDate.IsZero() || params.EndDate.IsZero() {
		err = fmt.Errorf("%w: beginDate and endDate must be provided", ErrInvalidParams)
		return
	}
	events = make([]model.Event, 0)
	s.mu.RLock()
	// Не захотел реализовать самодельный индекс по дате, чтобы работало быстрее, поэтому только полный перебор.
	for _, event := range s.events {
		if params.BeginDate.Year() <= event.Date.Year() && event.Date.Year() <= params.EndDate.Year() &&
			params.BeginDate.YearDay() <= event.Date.YearDay() && event.Date.YearDay() <= params.EndDate.YearDay() {
			events = append(events, event)
		}
	}
	s.mu.RUnlock()
	return
}

func NewInMemory() Repository {
	return &inMemory{events: make(map[model.EventID]model.Event)}
}
