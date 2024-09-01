package handler

import "github.com/kpechenenko/hw12_13_14_15_calendar/calendar/internal/service"

const PingPath = "ping"

type Handler struct {
	srv service.Service
}

func New(srv service.Service) *Handler {
	return &Handler{srv: srv}
}
