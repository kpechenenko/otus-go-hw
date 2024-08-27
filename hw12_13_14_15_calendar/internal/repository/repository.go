package repository

import (
	"context"
	"github.com/kpechenenko/hw12_13_14_15_calendar/internal/model"
	"time"
)

const (
	DateFormat     = "2006-01-02"
	DateTimeFormat = time.RFC3339
)

// EventRepository хранилище для взаимодействия с событиями.
type EventRepository interface {
	// CreateEvent создать событие.
	CreateEvent(ctx context.Context, params CreateEventParams) (model.EventID, error)
	// UpdateEvent обновить событие.
	UpdateEvent(ctx context.Context, params UpdateEventParams) error
	// DeleteEvent удалить событие.
	DeleteEvent(ctx context.Context, eventID model.EventID) error
	// FindEvent получить список событий.
	FindEvent(ctx context.Context, params FindEventParams) ([]model.Event, error)
}

// CreateEventParams параметры для сохранения события.
type CreateEventParams struct {
	Title       string         // Заголовок - короткий текст.
	Date        time.Time      // Дата и время события.
	Duration    time.Duration  // Длительность события.
	Description *string        // Описание события.
	OwnerUserID model.UserID   // Код пользователя владельца события.
	NotifyTime  *time.Duration // За сколько времени высылать уведомление/
}

// UpdateEventParams параметры для обновления события.
type UpdateEventParams struct {
	EventID     model.EventID  //  Уникальный идентификатор события.
	Title       string         // Заголовок - короткий текст.
	Date        time.Time      // Дата и время события.
	Duration    time.Duration  // Длительность события.
	Description *string        // Описание события.
	NotifyTime  *time.Duration // За сколько времени высылать уведомление/
}

// FindEventParams параметры для получения событий.
type FindEventParams struct {
	BeginDate time.Time // Дата начала поиска, включительно. Время не учитывается.
	EndDate   time.Time // Дата окончания поиска, включительно. Время не учитывается.
}
