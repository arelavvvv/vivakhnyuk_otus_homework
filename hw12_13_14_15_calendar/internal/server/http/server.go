package internalhttp

import (
	"context"
	"fmt"
	"net/http"

	//nolint:depguard
	"github.com/arelavvvv/vivakhnyuk_otus_homework/hw12_13_14_15_calendar/internal/app"
)

type Server struct {
	server *http.Server
	app    app.App
}

type MyHandler struct{}

func NewServer(app app.App, host string, port int) *Server {
	handler := &MyHandler{}

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.testHandler)
	fmt.Println("server will start on " + fmt.Sprintf("%s:%d", host, port))

	logger := NewLogger(mux, app.Logger)

	server := &http.Server{ //#nosec G112
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: logger,
	}

	return &Server{
		server,
		app,
	}
}

func (s *Server) Start(ctx context.Context) error {
	err := s.server.ListenAndServe()
	if err != nil {
		return err
	}

	<-ctx.Done()
	return nil
}

func (s *Server) Stop(_ context.Context) error {
	err := s.server.Close()
	return err
}

func (s *MyHandler) testHandler(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("Hello world!"))
	if err != nil {
		return
	}
}
