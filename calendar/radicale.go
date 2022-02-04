package calendar

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dolanor/caldav-go/caldav"
	"github.com/dolanor/caldav-go/caldav/entities"
	"github.com/dolanor/caldav-go/icalendar/components"
)

/*Radicale defines server information
ServerURL examples:
http://localhost/radicale/myuser/ (if its behind nginx)
http://localhost:5232/myuser/

Path examples:
/0ci89cde8-fa17-2396-efd8-b55d389cd4yy/
*/
type Radicale struct {
	ServerURL string
	Path      string
	client    *caldav.Client
}

type RadicaleEvent struct {
	Start   time.Time
	End     time.Time
	Summary string
}

func (r *Radicale) setClient() error {
	if r.client == nil {
		server, err := caldav.NewServer(r.ServerURL)
		if err != nil {
			return err
		}
		r.client = caldav.NewClient(server, http.DefaultClient)
	}

	err := r.client.ValidateServer(r.Path)
	if err != nil {
		return err
	}

	return nil
}

func (r *Radicale) CountEvents(start, end time.Time) (int, error) {
	err := r.setClient()
	if err != nil {
		return 0, err
	}

	query, err := entities.NewEventRangeQuery(start, end)
	if err != nil {
		return 0, err
	}

	events, err := r.client.QueryEvents(r.Path, query)
	if err != nil {
		return 0, err
	}

	printEvents(events)

	return len(events), err
}

func (r *Radicale) AddEvent(event interface{}) error {
	err := r.setClient()
	if err != nil {
		return err
	}

	radicaleEvent, ok := event.(RadicaleEvent)
	if !ok {
		return fmt.Errorf("invalid radicale event")
	}
	uuid := fmt.Sprintf("test-single-event-%d", radicaleEvent.End.Unix())
	putEvent := components.NewEventWithEnd(uuid, radicaleEvent.Start, radicaleEvent.End)
	putEvent.Summary = radicaleEvent.Summary

	err = r.client.PutEvents(r.Path, putEvent)
	if err != nil {
		return err
	}
	return nil
}

func printEvents(events []*components.Event) {
	for _, event := range events {
		fmt.Println(event.Summary)
	}
}
