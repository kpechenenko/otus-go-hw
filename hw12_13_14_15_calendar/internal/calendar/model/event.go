package model

import (
	"time"

	"github.com/google/uuid"
)

type EventID uuid.UUID // Тип идентификатора события в системе.

func (e EventID) String() string {
	return uuid.UUID(e).String()
}

type UserID int64 // Тип идентификатора пользователя в системе.

// Event событие в календаре.
type Event struct {
	ID          EventID       `json:"id,omitempty"`          // Уникальный идентификатор события.
	Title       string        `json:"title,omitempty"`       // Заголовок - короткий текст.
	Date        time.Time     `json:"date,omitempty"`        // Дата и время события.
	Duration    time.Duration `json:"duration,omitempty"`    // Длительность события.
	Description *string       `json:"description,omitempty"` // Описание события - длинный текст, опционально.
	OwnerID     UserID        `json:"ownerId,omitempty"`     // ID пользователя, владельца события.
	// За сколько времени до наступления события высылать уведомление, опционально.
	NotifyFor *time.Duration `json:"notifyFor,omitempty"`
}
