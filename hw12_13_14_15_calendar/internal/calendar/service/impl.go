package service

import (
	"context"
	"time"

	"github.com/kpechenenko/hw12_13_14_15_calendar/calendar/internal/calendar/model"
	"github.com/kpechenenko/hw12_13_14_15_calendar/calendar/internal/calendar/repository"
)

type service struct {
	repo repository.Repository
}

func (s *service) CreateEvent(ctx context.Context, params repository.AddEventParams) (model.EventID, error) {
	return s.repo.AddEvent(ctx, params)
}

func (s *service) UpdateEvent(ctx context.Context, params repository.UpdateEventParams) error {
	return s.repo.UpdateEvent(ctx, params)
}

func (s *service) DeleteEvent(ctx context.Context, id model.EventID) error {
	return s.repo.DeleteEvent(ctx, id)
}

func (s *service) GetEventsForDay(ctx context.Context, day time.Time) ([]model.Event, error) {
	params := repository.GetEventParams{BeginDate: day, EndDate: day}
	return s.repo.GetEvents(ctx, params)
}

func (s *service) GetEventsForWeek(ctx context.Context, beginDate time.Time) ([]model.Event, error) {
	params := repository.GetEventParams{BeginDate: beginDate, EndDate: getLastDayOfWeek(beginDate)}
	return s.repo.GetEvents(ctx, params)
}

func (s *service) GetEventsForMonth(ctx context.Context, beginDate time.Time) ([]model.Event, error) {
	params := repository.GetEventParams{BeginDate: beginDate, EndDate: getLastDayOfMonth(beginDate)}
	return s.repo.GetEvents(ctx, params)
}

func New(repo repository.Repository) Service {
	return &service{repo: repo}
}
