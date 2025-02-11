package memorystorage

import (
	"errors"
	"sync"

	//nolint:depguard
	"github.com/arelavvvv/vivakhnyuk_otus_homework/hw12_13_14_15_calendar/internal/storage"
	//nolint:depguard
	"github.com/golang-module/carbon"
)

type Storage struct {
	events []storage.Event
	mu     sync.RWMutex
}

func New() *Storage {
	var events []storage.Event
	return &Storage{
		events,
		sync.RWMutex{},
	}
}

func (s *Storage) IntervalCheck(id int, dateFrom string, dateTo string) error {
	for _, event := range s.events {
		if (dateFrom >= event.DateTimeStart && dateFrom <= event.DateTimeEnd ||
			dateTo >= event.DateTimeStart && dateTo <= event.DateTimeEnd) && !(id > 0 && event.ID == id) {
			return storage.DateBusyError{
				Title:    event.Title,
				DateFrom: event.DateTimeStart,
				DateTo:   event.DateTimeEnd,
			}
		}
	}
	return nil
}

func (s *Storage) NewEvent(
	title string,
	dateTimeStart string,
	dateTimeEnd string,
	description string,
	userID int,
	delay int,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := s.IntervalCheck(0, dateTimeStart, dateTimeEnd)
	if err != nil {
		return err
	}

	event := storage.Event{
		ID:            s.GenerateID(),
		Title:         title,
		DateTimeStart: dateTimeStart,
		DateTimeEnd:   dateTimeEnd,
		Description:   description,
		UserID:        userID,
		Delay:         delay,
	}
	s.events = append(s.events, event)
	// TODO errors/validation
	return nil
}

func (s *Storage) GenerateID() int {
	if len(s.events) == 0 {
		return 1
	}

	lastEvent := s.events[len(s.events)-1]
	return lastEvent.ID + 1
}

func (s *Storage) EventDelete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, event := range s.events {
		if event.ID == id {
			s.events = append(s.events[:i], s.events[i+1:]...)
			return nil
		}
	}
	return errors.New("no such event")
}

func (s *Storage) UpdEvent(
	id int,
	title string,
	dateTimeStart string,
	dateTimeEnd string,
	description string,
	userID int,
	delay int,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := s.IntervalCheck(id, dateTimeStart, dateTimeEnd)
	if err != nil {
		return err
	}

	updatedEvent := storage.Event{
		ID:            s.GenerateID(),
		Title:         title,
		DateTimeStart: dateTimeStart,
		DateTimeEnd:   dateTimeEnd,
		Description:   description,
		UserID:        userID,
		Delay:         delay,
	}

	for i, event := range s.events {
		if event.ID == id {
			s.events[i] = updatedEvent
			return nil
		}
	}
	return errors.New("no such event")
}

func (s *Storage) GetEvents() ([]storage.Event, error) {
	return s.events, nil
}

func (s *Storage) GetEventsDay(day string) ([]storage.Event, error) {
	return s.GetEventsInterval(day, "day")
}

func (s *Storage) GetEventsWeek(day string) ([]storage.Event, error) {
	return s.GetEventsInterval(day, "week")
}

func (s *Storage) GetEventsMonth(day string) ([]storage.Event, error) {
	return s.GetEventsInterval(day, "month")
}

func (s *Storage) GetEventsInterval(day string, interval string) ([]storage.Event, error) {
	pday := carbon.Parse(day)
	startDate := pday.Format("Y-m-d")

	var endDate string
	switch interval {
	case "day":
		endDate = pday.AddDay().Format("Y-m-d")
	case "week":
		endDate = pday.AddWeek().Format("Y-m-d")
	case "month":
		endDate = pday.AddMonth().Format("Y-m-d")
	default:
		return nil, errors.New("wrong interval name ")
	}

	var events []storage.Event
	for _, event := range s.events {
		if event.DateTimeStart >= startDate && event.DateTimeStart < endDate {
			events = append(events, event)
		}
	}

	return events, nil
}
