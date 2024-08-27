package repository

import (
	"context"
	"fmt"
	"github.com/kpechenenko/hw12_13_14_15_calendar/internal/model"
	"sync"
)

const (
	startEventID   = 1
	stepForEventID = 1
)

// InMemoryEventRepository хранит события в оперативной памяти.
type InMemoryEventRepository struct {
	mu             sync.Mutex
	lastEventID    model.EventID
	eventIDToEvent map[model.EventID]model.Event
}

// NewInMemoryEventRepository конструктор с параметрами.
func NewInMemoryEventRepository() *InMemoryEventRepository {
	return &InMemoryEventRepository{
		eventIDToEvent: make(map[model.EventID]model.Event),
		lastEventID:    startEventID,
	}
}

func (r *InMemoryEventRepository) generateNextEventID() model.EventID {
	r.lastEventID += stepForEventID
	return r.lastEventID
}

func (r *InMemoryEventRepository) CreateEvent(ctx context.Context, params CreateEventParams) (id model.EventID, err error) {
	if len(params.Title) == 0 || params.Date.IsZero() || params.Duration.Abs() == 0 || params.OwnerUserID == 0 {
		err = fmt.Errorf("%w: event title, date, ownerUserId must be provided", ErrInvalidParam)
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	id = r.generateNextEventID()
	event := model.Event{
		EventID:     id,
		Title:       params.Title,
		Date:        params.Date,
		Duration:    params.Duration,
		Description: params.Description,
		OwnerUserID: params.OwnerUserID,
		NotifyTime:  params.NotifyTime,
	}
	r.eventIDToEvent[id] = event
	return
}

func (r *InMemoryEventRepository) UpdateEvent(ctx context.Context, params UpdateEventParams) error {
	if params.EventID == 0 || len(params.Title) == 0 || params.Date.IsZero() || params.Duration.Abs() == 0 {
		return fmt.Errorf("%w: params id, title, date must be provided", ErrInvalidParam)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	current, ok := r.eventIDToEvent[params.EventID]
	if !ok {
		return fmt.Errorf("%w: eventID #%d", ErrEventDoesNotExist, params.EventID)
	}
	current.Title = params.Title
	current.Date = params.Date
	current.Duration = params.Duration
	current.Description = params.Description
	current.NotifyTime = params.NotifyTime
	r.eventIDToEvent[params.EventID] = current
	return nil
}

func (r *InMemoryEventRepository) DeleteEvent(ctx context.Context, eventID model.EventID) (err error) {
	if eventID == 0 {
		err = fmt.Errorf("%w: eventID must be provided", ErrInvalidParam)
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.eventIDToEvent, eventID)
	return nil
}

func (r *InMemoryEventRepository) FindEvent(ctx context.Context, params FindEventParams) (events []model.Event, err error) {
	if params.BeginDate.IsZero() || params.EndDate.IsZero() {
		err = fmt.Errorf("%w: beginDate and endDate must be provided", ErrInvalidParam)
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	events = make([]model.Event, 0, 10)
	// Не захотел реализовать самодельный индекс по дате, чтобы работало быстрее, поэтому только полный перебор.
	for _, event := range r.eventIDToEvent {
		if params.BeginDate.Year() <= event.Date.Year() && event.Date.Year() <= params.EndDate.Year() &&
			params.BeginDate.YearDay() <= event.Date.YearDay() && event.Date.YearDay() <= params.EndDate.YearDay() {
			events = append(events, event)
		}
	}
	return
}
