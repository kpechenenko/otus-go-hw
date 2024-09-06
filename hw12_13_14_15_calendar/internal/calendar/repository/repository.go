package repository

import (
	"context"
	"time"

	"github.com/kpechenenko/hw12_13_14_15_calendar/calendar/internal/calendar/model"
)

const (
	DateFormat     = "2006-01-02"
	DateTimeFormat = time.RFC3339
)

// Repository хранилище событий.
type Repository interface {
	AddEvent(context.Context, AddEventParams) (model.EventID, error)
	UpdateEvent(context.Context, UpdateEventParams) error
	DeleteEvent(context.Context, model.EventID) error
	GetEvents(context.Context, GetEventParams) ([]model.Event, error)
}

type GetEventParams struct {
	BeginDate time.Time `json:"beginDate"` // Дата начала поиска, включительно. Время не учитывается.
	EndDate   time.Time `json:"endDate"`   // Дата окончания поиска, включительно. Время не учитывается.
}

// AddEventParams параметры для добавления события.
type AddEventParams struct {
	Title       string         `json:"title"`       // Заголовок - короткий текст.
	Date        time.Time      `json:"date"`        // Дата и время события.
	Duration    time.Duration  `json:"duration"`    // Длительность события.
	Description *string        `json:"description"` // Описание события.
	OwnerID     model.UserID   `json:"ownerId"`     // Код пользователя владельца события.
	NotifyFor   *time.Duration `json:"notifyFor"`   // За сколько времени высылать уведомление?
}

// UpdateEventParams параметры для обновления события.
type UpdateEventParams struct {
	ID          model.EventID  `json:"id"`          //  Уникальный идентификатор события.
	Title       string         `json:"title"`       // Заголовок - короткий текст.
	Date        time.Time      `json:"date"`        // Дата и время события.
	Duration    time.Duration  `json:"duration"`    // Длительность события.
	Description *string        `json:"description"` // Описание события.
	NotifyFor   *time.Duration `json:"notifyFor"`   // За сколько времени высылать уведомление?
}
