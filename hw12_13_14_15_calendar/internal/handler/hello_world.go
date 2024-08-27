package handler

import "net/http"

func (h *Handler) HelloWorld(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("hello world"))
	if err != nil {
		h.logger.Errorf("write \"hello world\": %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
