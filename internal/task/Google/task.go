package task

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/sync/errgroup"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/tasks/v1"
)

type TaskResult struct {
	Summary string
	Date    string
}

func (s TaskResult) String() string {
	return fmt.Sprintf("%s - %s", s.Summary, s.Date)
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

func merge(in1, in2 <-chan TaskResult) <-chan TaskResult {
	out := make(chan TaskResult)
	go func() {
		defer close(out)
		for in1 != nil || in2 != nil {
			select {
			case val1, ok := <-in1:
				if ok {
					out <- val1
				} else {
					in1 = nil
				}

			case val2, ok := <-in2:
				if ok {
					out <- val2
				} else {
					in2 = nil
				}
			}
		}
	}()
	return out
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

	errgr, ctx := errgroup.WithContext(context.Background())

	resultcal := make(chan TaskResult)
	errgr.Go(func() error {
		defer close(resultcal)
		if err := g.infoCal(calendarid, resultcal, ctx); err != nil {
			return err
		}
		return nil
	})

	resulttask := make(chan TaskResult)
	errgr.Go(func() error {
		defer close(resulttask)
		if err := g.infoTasks(resulttask, ctx); err != nil {
			return err
		}
		return nil
	})

	resultch := merge(resultcal, resulttask)

	result := []TaskResult{}
	for item := range resultch {
		result = append(result, item)
	}

	if err := errgr.Wait(); err != nil {
		return nil, err
	} else {
		return result, nil
	}

}

func (g *GCalendar) infoCal(calendarid string, out chan<- TaskResult, ctx context.Context) error {

	curtime := time.Now()
	tstart := time.Date(curtime.Year(), curtime.Month(), curtime.Day(), 0, 0, 0, 0, curtime.Location()).Format(time.RFC3339)
	tend := time.Date(curtime.Year(), curtime.Month(), curtime.Day(), 23, 59, 59, 59, curtime.Location()).Format(time.RFC3339)

	events, err := g.calendarService.Events.List(calendarid).ShowDeleted(false).
		SingleEvents(true).TimeMin(tstart).TimeMax(tend).MaxResults(10).OrderBy("startTime").Do()
	if err != nil {
		return err
	}

	for _, item := range events.Items {
		var result TaskResult
		if item.Start.DateTime == "" {
			result.Date = "весь день"
		} else {
			date, _ := time.Parse(time.RFC3339, item.Start.DateTime)
			result.Date = date.Format("15:04:05")
		}
		result.Summary = item.Summary

		select {
		case <-ctx.Done():
			return ctx.Err()
		case out <- result:
		}

	}

	return nil

}

func (g *GCalendar) infoTasks(out chan<- TaskResult, ctx context.Context) error {

	tasklists, err := g.taskService.Tasklists.List().Do()
	if err != nil {
		return err
	}

	curtime := time.Now()
	startofday := time.Date(curtime.Year(), curtime.Month(), curtime.Day(), 0, 0, 0, 0, time.UTC)
	tend := startofday.AddDate(0, 0, 1).Format(time.RFC3339)
	tstart := startofday.Format(time.RFC3339)

	for _, tasklist := range tasklists.Items {

		events, err := g.taskService.Tasks.List(tasklist.Id).ShowCompleted(false).
			ShowDeleted(false).ShowHidden(false).DueMin(tstart).DueMax(tend).Do()
		if err != nil {
			return err
		}

		if len(events.Items) > 0 {
			for _, item := range events.Items {
				var result TaskResult
				result.Summary = item.Title
				result.Date = "сегодня"

				select {
				case <-ctx.Done():
					return ctx.Err()
				case out <- result:
				}

			}
		}

	}

	return nil

}
