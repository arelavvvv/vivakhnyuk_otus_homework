package internalhttp

import (
	"net/http"
	"time"

	//nolint:depguard
	"github.com/arelavvvv/vivakhnyuk_otus_homework/hw12_13_14_15_calendar/internal/app"
	//nolint:depguard
	"github.com/golang-module/carbon"
)

type Logger struct {
	handler http.Handler
	logger  app.Logger
}

func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	rCtx := r.Clone(r.Context())
	l.handler.ServeHTTP(w, rCtx)

	log := r.RemoteAddr + " " +
		carbon.Now().Format("Y-m-d H:i:s") + " " +
		r.Method + " " +
		r.Proto + " " +
		r.RequestURI + " " +
		r.UserAgent() +
		"; request_time: " + time.Since(start).String()
	l.logger.Debug(log)
}

func NewLogger(handlerToWrap http.Handler, logger app.Logger) *Logger {
	return &Logger{handlerToWrap, logger}
}
