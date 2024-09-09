package repository

import (
	"github.com/google/uuid"
	"github.com/kpechenenko/hw12_13_14_15_calendar/calendar/internal/calendar/model"
)

func GenerateEventID() model.EventID {
	return model.EventID(uuid.New())
}
