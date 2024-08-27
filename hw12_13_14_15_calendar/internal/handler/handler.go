package handler

import (
	"github.com/kpechenenko/hw12_13_14_15_calendar/internal/logger"
	"github.com/kpechenenko/hw12_13_14_15_calendar/internal/service"
	"go.uber.org/zap"
)

type Handler struct {
	eventService service.EventService
	logger       *zap.SugaredLogger
}

func NewHandler(eventService service.EventService) *Handler {
	const loggerName = "handler"
	return &Handler{eventService: eventService, logger: logger.GetNamed(loggerName)}
}
