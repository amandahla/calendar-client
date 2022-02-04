package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"githum.com/amandahla/calendar-client/calendar"
)

func main() {
	radicaleServerURL := os.Getenv("RADICALE_SERVERURL")
	radicalePath := os.Getenv("RADICALE_PATH")

	my_calendar := calendar.Radicale{ServerURL: radicaleServerURL, Path: radicalePath}

	log.Println("Add Big Event")
	start := time.Now().UTC().AddDate(0, 0, 1).UTC()
	end := start.AddDate(0, 0, 1).UTC()
	err := my_calendar.AddEvent(calendar.RadicaleEvent{Start: start, End: end, Summary: "Big Event"})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Big Event added")

	log.Println("List events from now until next 10 days")
	start = time.Now().UTC()
	end = start.AddDate(0, 0, 10).UTC()
	number_events_next_ten_days, err := my_calendar.CountEvents(start, end)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(number_events_next_ten_days)
}
