package weather

import (
	"context"

	"github.com/hectormalot/omgo"
)

type weatherGeter struct {
	client *omgo.Client
}

func New() (*weatherGeter, error) {

	c, err := omgo.NewClient()
	if err != nil {
		return nil, err
	}

	return &weatherGeter{client: &c}, nil

}

func (w *weatherGeter) Info(ctx context.Context, timezone string, lat, lon float64) error { // TODO return value

	loc, err := omgo.NewLocation(lat, lon)
	if err != nil {
		return err
	}

	opts := omgo.Options{
		TemperatureUnit:   "celsius",
		WindspeedUnit:     "ms",
		PrecipitationUnit: "mm",
		Timezone:          timezone,
		PastDays:          0,
		HourlyMetrics:     []string{"temperature_2m,windspeed_10m,rain,weathercode"},
	}

	_, err = w.client.Forecast(ctx, loc, &opts)
	if err != nil {
		return err
	}

	return nil

}
