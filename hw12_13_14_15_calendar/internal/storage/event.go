package storage

import "fmt"

type DateBusyError struct {
	Title    string
	DateFrom string
	DateTo   string
}

func (e DateBusyError) Error() string {
	return fmt.Sprintf("this date is busy with %s event: Date from - %s; Date to - %s", e.Title, e.DateFrom, e.DateTo)
}

type Event struct {
	ID            int
	Title         string
	DateTimeStart string
	DateTimeEnd   string
	Description   string
	UserID        int
	Delay         int
}
