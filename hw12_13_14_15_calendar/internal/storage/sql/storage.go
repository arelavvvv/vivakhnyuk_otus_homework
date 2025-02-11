package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"

	//nolint:depguard
	"github.com/arelavvvv/vivakhnyuk_otus_homework/hw12_13_14_15_calendar/internal/storage"
	//nolint:depguard
	"github.com/golang-module/carbon"
	//nolint:depguard
	_ "github.com/jackc/pgx/stdlib"
)

type Storage struct {
	username string
	password string
	host     string
	port     int
	dbname   string
	db       *sql.DB
}

func (s *Storage) IntervalCheck(id int, dateFrom string, dateTo string) error {
	query := fmt.Sprintf("where start_date between "+
		"%s and %s or end_date between %s and %s",
		dateFrom, dateTo, dateFrom, dateTo)

	if id > 0 {
		query += fmt.Sprintf(" and id != %d", id)
	}

	events, err := s.ReqEvents(query)
	if err != nil {
		return err
	}

	if len(events) > 0 {
		return storage.DateBusyError{
			Title:    events[0].Title,
			DateFrom: events[0].DateTimeStart,
			DateTo:   events[0].DateTimeEnd,
		}
	}

	return nil
}

func New(username string, password string, host string, port int, dbname string) *Storage {
	return &Storage{
		username,
		password,
		host,
		port,
		dbname,
		nil,
	}
}

func (s *Storage) Connect(
	ctx context.Context,
) error {
	dsn := "user=" + s.username +
		" dbname=" + s.dbname +
		" sslmode=disable password=" + s.password +
		" host=" + s.host +
		" port=" + strconv.Itoa(s.port)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return err
	}

	err = db.PingContext(ctx)
	if err != nil {
		return err
	}

	s.db = db
	return nil
}

func (s *Storage) Close() error {
	if s.db != nil {
		err := s.db.Close()
		return err
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
	err := s.IntervalCheck(0, dateTimeStart, dateTimeEnd)
	if err != nil {
		return err
	}
	ctx := context.Background()
	err = s.Connect(ctx)
	if err != nil {
		return err
	}
	defer func(s *Storage) {
		err := s.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(s)

	query := "insert into " +
		"events(title, start_date, end_date, description, user_id, delay)" +
		" values($1, $2, $3, $4, $5, $6)"
	_, err = s.db.Exec(query, title, dateTimeStart, dateTimeEnd, description, userID, delay)

	return err
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
	err := s.IntervalCheck(id, dateTimeStart, dateTimeEnd)
	if err != nil {
		return err
	}

	ctx := context.Background()
	err = s.Connect(ctx)
	if err != nil {
		return err
	}
	defer func(s *Storage) {
		err := s.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(s)

	query := "update events set " +
		"title=$1, start_date=$2, end_date=$3, description=$4, user_id=$5, delay=$6 where id=$7"
	_, err = s.db.Exec(query, title, dateTimeStart, dateTimeEnd, description, userID, delay, id)
	return err
}

func (s *Storage) EventDelete(id int) error {
	ctx := context.Background()
	err := s.Connect(ctx)
	if err != nil {
		return err
	}
	defer func(s *Storage) {
		err := s.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(s)

	query := "delete from events where id=$1"
	_, err = s.db.Exec(query, id)
	return err
}

func (s *Storage) GetEvents() ([]storage.Event, error) {
	events, err := s.ReqEvents("")
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (s *Storage) ReqEvents(query string) ([]storage.Event, error) {
	ctx := context.Background()
	err := s.Connect(ctx)
	if err != nil {
		return nil, err
	}
	defer func(s *Storage) {
		err := s.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(s)

	evquery := "select * from events " + query //#nosec G202
	rows, err := s.db.Query(evquery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []storage.Event
	for rows.Next() {
		var id, userID, delay int
		var title, dateTimeStart, dateTimeEnd, description string
		if err := rows.Scan(&id, &title, &dateTimeStart, &dateTimeEnd, &description, &userID, &delay); err != nil {
			fmt.Println(err)
		}
		events = append(events, storage.Event{
			ID:            id,
			Description:   description,
			Delay:         delay,
			DateTimeStart: dateTimeStart,
			DateTimeEnd:   dateTimeEnd,
			UserID:        userID,
			Title:         title,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
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
	dateStart := pday.Format("Y-m-d")

	var dateEnd string
	switch interval {
	case "day":
		dateEnd = pday.AddDay().Format("Y-m-d")
	case "week":
		dateEnd = pday.AddWeek().Format("Y-m-d")
	case "month":
		dateEnd = pday.AddMonth().Format("Y-m-d")
	default:
		return nil, errors.New("wrong interval")
	}

	query := fmt.Sprintf(" and where start_date >= %s and start_date < %s", dateStart, dateEnd)
	return s.ReqEvents(query)
}
