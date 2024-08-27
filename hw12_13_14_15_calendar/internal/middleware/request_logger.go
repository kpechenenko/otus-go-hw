package middleware

import (
	"github.com/kpechenenko/hw12_13_14_15_calendar/internal/logger"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type RequestLogger struct {
	handler http.Handler
	logger  *zap.SugaredLogger
}

type responseWriterWrapper struct {
	http.ResponseWriter
	code int
}

func (rw *responseWriterWrapper) WriteHeader(statusCode int) {
	rw.code = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (l *RequestLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	rw := &responseWriterWrapper{w, -1}
	l.handler.ServeHTTP(w, r)
	//IP клиента;
	//дата и время запроса;
	//метод, path и версия HTTP;
	//код ответа;
	//latency (время обработки запроса, посчитанное, например, с помощью middleware);
	//user agent, если есть.
	//todo check statuscode
	//zap.Int("statusCode", w.Header()),
	l.logger.Infow("handle request",
		zap.String("ip", r.RemoteAddr),
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
		zap.String("httpVersion", r.Proto),
		zap.Int64("latency", time.Since(start).Milliseconds()),
		zap.String("userAgent", r.UserAgent()),
		zap.Int("code", rw.code),
	)
}

func NewRequestLogger(handlerToWrap http.Handler) *RequestLogger {
	const loggerName = "requestLogger"
	return &RequestLogger{handlerToWrap, logger.GetNamed(loggerName)}
}
