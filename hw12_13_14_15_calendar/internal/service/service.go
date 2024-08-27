package service

import (
	"context"
	"github.com/kpechenenko/hw12_13_14_15_calendar/internal/logger"
	"github.com/kpechenenko/hw12_13_14_15_calendar/internal/model"
	"github.com/kpechenenko/hw12_13_14_15_calendar/internal/repository"
	"go.uber.org/zap"
	"time"
)

// EventService бизнес логика для взаимодействия с событиями календаря.
type EventService interface {
	CreateEvent(ctx context.Context, params repository.CreateEventParams) (model.EventID, error)
	UpdateEvent(ctx context.Context, params repository.UpdateEventParams) error
	DeleteEvent(ctx context.Context, eventID model.EventID) error
	GetEventsForDay(ctx context.Context, date time.Time) ([]model.Event, error)
	GetEventsForWeek(ctx context.Context, beginDate time.Time) ([]model.Event, error)
	GetEventsForMonth(ctx context.Context, beginDate time.Time) ([]model.Event, error)
}

// EventServiceImpl бизнес логика для взаимодействия с событиями календаря.
type EventServiceImpl struct {
	repo   repository.EventRepository
	logger *zap.SugaredLogger
}

func NewEventService(repo repository.EventRepository) *EventServiceImpl {
	const loggerName = "eventService"
	return &EventServiceImpl{
		repo:   repo,
		logger: logger.GetNamed(loggerName),
	}
}

func (s *EventServiceImpl) CreateEvent(ctx context.Context, params repository.CreateEventParams) (model.EventID, error) {
	return s.repo.CreateEvent(ctx, params)
}

func (s *EventServiceImpl) UpdateEvent(ctx context.Context, params repository.UpdateEventParams) error {
	return s.repo.UpdateEvent(ctx, params)
}

func (s *EventServiceImpl) DeleteEvent(ctx context.Context, eventID model.EventID) error {
	return s.repo.DeleteEvent(ctx, eventID)
}

func (s *EventServiceImpl) GetEventsForDay(ctx context.Context, date time.Time) ([]model.Event, error) {
	params := repository.FindEventParams{
		BeginDate: date,
		EndDate:   date,
	}
	return s.repo.FindEvent(ctx, params)
}

func (s *EventServiceImpl) GetEventsForWeek(ctx context.Context, beginDate time.Time) ([]model.Event, error) {
	params := repository.FindEventParams{
		BeginDate: beginDate,
		EndDate:   getLastDayOfWeek(beginDate),
	}
	return s.repo.FindEvent(ctx, params)
}

func (s *EventServiceImpl) GetEventsForMonth(ctx context.Context, beginDate time.Time) ([]model.Event, error) {
	params := repository.FindEventParams{
		BeginDate: beginDate,
		EndDate:   getLastDayOfMonth(beginDate),
	}
	return s.repo.FindEvent(ctx, params)
}
