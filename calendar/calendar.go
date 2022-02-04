package calendar

import "time"

type Calendar interface {
	CountEvents(start, end time.Time) (int, error)
	AddEvent(event interface{}) error
}
