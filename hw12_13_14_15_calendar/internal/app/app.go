package app

import (
	"context"

	//nolint:depguard
	"github.com/arelavvvv/vivakhnyuk_otus_homework/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	Logger  Logger
	Storage Storage
}

type Logger interface {
	Info(msg string)
	Error(msg string)
	Warn(msg string)
	Debug(msg string)
}

type Storage interface {
	NewEvent(
		Title string,
		DateTimeStart string,
		DateTimeEnd string,
		Description string,
		UserID int,
		Delay int,
	) error
	UpdEvent(
		ID int,
		Title string,
		DateTimeStart string,
		DateTimeEnd string,
		Description string,
		UserID int,
		Delay int,
	) error
	EventDelete(ID int) error
	GetEvents() ([]storage.Event, error)
	GetEventsDay(day string) ([]storage.Event, error)
	GetEventsWeek(day string) ([]storage.Event, error)
	GetEventsMonth(day string) ([]storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{
		Logger:  logger,
		Storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error { //nolint:revive
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
