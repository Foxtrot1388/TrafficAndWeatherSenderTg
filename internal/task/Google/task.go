package task

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type TaskResult struct {
	Summary string
	Date    string
}

type GCalendar struct {
	calendarService *calendar.Service
}

func New(ctx context.Context) (*GCalendar, error) {

	b, err := os.ReadFile("../credentials.json")
	if err != nil {
		return nil, err
	}

	conf, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		return nil, err
	}

	token, _ := tokenFromFile("../token.json")
	calendarService, err := calendar.NewService(ctx, option.WithTokenSource(conf.TokenSource(ctx, token)))
	if err != nil {
		return nil, err
	}

	return &GCalendar{calendarService: calendarService}, nil

}

func (g *GCalendar) Info() ([]TaskResult, error) {

	curtime := time.Now()
	tstart := time.Date(curtime.Year(), curtime.Month(), curtime.Day(), 0, 0, 0, 0, curtime.Location()).Format(time.RFC3339)
	tend := time.Date(curtime.Year(), curtime.Month(), curtime.Day(), 23, 59, 59, 59, curtime.Location()).Format(time.RFC3339)
	events, err := g.calendarService.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).TimeMin(tstart).TimeMax(tend).MaxResults(10).OrderBy("startTime").Do()
	if err != nil {
		return nil, err
	}

	result := make([]TaskResult, len(events.Items))
	for i, item := range events.Items {
		date := item.Start.DateTime
		if date == "" {
			date = item.Start.Date
		}
		result[i].Summary = item.Summary
		result[i].Date = date
	}

	return result, nil

}

func (s TaskResult) String() string {
	return fmt.Sprintf("%s - %s", s.Summary, s.Date)
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}
