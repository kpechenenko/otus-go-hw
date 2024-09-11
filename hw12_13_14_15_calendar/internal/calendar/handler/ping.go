package handler

import "net/http"

func (h *Handler) Ping(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("pong"))
}

func (h *Handler) PingWithParams(w http.ResponseWriter, _ *http.Request, _ map[string]string) {
	_, _ = w.Write([]byte("pong"))
}
