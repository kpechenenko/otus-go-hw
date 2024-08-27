package model

import (
	"time"
)

// EventID идентификатор события.
type EventID int64

// Event событие в календаре.
type Event struct {
	EventID     EventID        // Уникальный идентификатор события.
	Title       string         // Заголовок.
	Date        time.Time      // Дата и время события.
	Duration    time.Duration  // Длительность события, наносекунды.
	Description *string        // Описание события.
	OwnerUserID UserID         // Код пользователя владельца события.
	NotifyTime  *time.Duration // За сколько времени высылать уведомление? Наносекунды. Временная метка: Date - NotifyTime
}
