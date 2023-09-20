package task

import (
	"context"
	"fmt"
	"time"

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

func New(ctx context.Context, credentialsfile string) (*GCalendar, error) {

	calendarService, err := calendar.NewService(ctx, option.WithCredentialsFile("../"+credentialsfile))
	if err != nil {
		return nil, err
	}

	return &GCalendar{calendarService: calendarService}, nil

}

func (g *GCalendar) Info(calendarid string) ([]TaskResult, error) {

	curtime := time.Now()
	tstart := time.Date(curtime.Year(), curtime.Month(), curtime.Day(), 0, 0, 0, 0, curtime.Location()).Format(time.RFC3339)
	tend := time.Date(curtime.Year(), curtime.Month(), curtime.Day(), 23, 59, 59, 59, curtime.Location()).Format(time.RFC3339)
	events, err := g.calendarService.Events.List(calendarid).ShowDeleted(false).
		SingleEvents(true).TimeMin(tstart).TimeMax(tend).MaxResults(10).OrderBy("startTime").Do()
	if err != nil {
		return nil, err
	}

	result := make([]TaskResult, len(events.Items))
	for i, item := range events.Items {
		var date time.Time
		if item.Start.DateTime == "" {
			result[i].Date = "весь день"
		} else {
			date, _ = time.Parse(time.RFC3339, item.Start.DateTime)
			result[i].Date = date.Format("15:04:05")
		}
		result[i].Summary = item.Summary
	}

	return result, nil

}

func (s TaskResult) String() string {
	return fmt.Sprintf("%s - %s", s.Summary, s.Date)
}
