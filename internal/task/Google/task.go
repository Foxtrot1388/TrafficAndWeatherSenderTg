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
	"google.golang.org/api/tasks/v1"
)

type TaskResult struct {
	Summary string
	Date    string
}

type GCalendar struct {
	calendarService *calendar.Service
	taskService     *tasks.Service
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

func New(ctx context.Context, tokenfile, clientsecret string) (*GCalendar, error) {

	b, err := os.ReadFile("../" + clientsecret)
	if err != nil {
		return nil, err
	}

	conf, err := google.ConfigFromJSON(b, tasks.TasksReadonlyScope, calendar.CalendarReadonlyScope)
	if err != nil {
		return nil, err
	}

	token, _ := tokenFromFile("../" + tokenfile)

	taskService, err := tasks.NewService(ctx, option.WithTokenSource(conf.TokenSource(ctx, token)))
	if err != nil {
		return nil, err
	}

	calendarService, err := calendar.NewService(ctx, option.WithTokenSource(conf.TokenSource(ctx, token)))
	if err != nil {
		return nil, err
	}

	return &GCalendar{calendarService: calendarService, taskService: taskService}, nil

}

func (g *GCalendar) Info(calendarid string) ([]TaskResult, error) {

	resultcal, err := g.infoCal(calendarid)
	if err != nil {
		return nil, err
	}

	resulttask, err := g.infoTasks()
	if err != nil {
		return nil, err
	}

	result := append(resultcal, resulttask...)
	return result, nil

}

func (g *GCalendar) infoCal(calendarid string) ([]TaskResult, error) {

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

func (g *GCalendar) infoTasks() ([]TaskResult, error) {

	tasklists, err := g.taskService.Tasklists.List().Do()
	if err != nil {
		return nil, err
	}

	curtime := time.Now()
	startofday := time.Date(curtime.Year(), curtime.Month(), curtime.Day(), 0, 0, 0, 0, time.UTC)
	tend := startofday.AddDate(0, 0, 1).Format(time.RFC3339)
	tstart := startofday.Format(time.RFC3339)

	resultall := make([]TaskResult, 0)
	for _, tasklist := range tasklists.Items {

		events, err := g.taskService.Tasks.List(tasklist.Id).ShowCompleted(false).
			ShowDeleted(false).ShowHidden(false).DueMin(tstart).DueMax(tend).Do()
		if err != nil {
			return nil, err
		}

		if len(events.Items) > 0 {
			result := make([]TaskResult, len(events.Items))
			for i, item := range events.Items {
				result[i].Summary = item.Title
				result[i].Date = "сегодня"
			}
			resultall = append(resultall, result...)
		}

	}

	return resultall, nil

}

func (s TaskResult) String() string {
	return fmt.Sprintf("%s - %s", s.Summary, s.Date)
}
