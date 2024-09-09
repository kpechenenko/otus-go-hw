package service

import (
	"context"
	"time"

	"github.com/kpechenenko/hw12_13_14_15_calendar/calendar/internal/calendar/model"
	"github.com/kpechenenko/hw12_13_14_15_calendar/calendar/internal/calendar/repository"
)

type Service interface {
	CreateEvent(ctx context.Context, e repository.AddEventParams) (model.EventID, error)
	UpdateEvent(ctx context.Context, e repository.UpdateEventParams) error
	DeleteEvent(ctx context.Context, id model.EventID) error
	GetEventsForDay(ctx context.Context, day time.Time) ([]model.Event, error)
	GetEventsForWeek(ctx context.Context, beginDate time.Time) ([]model.Event, error)
	GetEventsForMonth(ctx context.Context, beginDate time.Time) ([]model.Event, error)
}
