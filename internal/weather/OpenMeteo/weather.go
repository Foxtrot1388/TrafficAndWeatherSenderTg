package weather

import (
	"context"
	"time"

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

func (w *weatherGeter) Info(ctx context.Context, timezone string, lat float64, lon float64, timeAt int, timeTo int) (*WeatherInfo, error) { // TODO return value

	loc, err := omgo.NewLocation(lat, lon)
	if err != nil {
		return nil, err
	}

	// TODO forecast_days

	opts := omgo.Options{
		TemperatureUnit:   "celsius",
		WindspeedUnit:     "ms",
		PrecipitationUnit: "mm",
		Timezone:          timezone,
		PastDays:          0,
		HourlyMetrics:     []string{"temperature_2m,windspeed_10m,rain,weathercode"},
	}

	resultforecst, err := w.client.Forecast(ctx, loc, &opts)
	if err != nil {
		return nil, err
	}

	var result WeatherInfo
	result.TemperatureUnit = "C"
	result.PrecipitationUnit = "мм"
	result.WindspeedUnit = "м/с"
	result.CurrentWeather.Temperature = resultforecst.CurrentWeather.Temperature
	result.CurrentWeather.WindSpeed = resultforecst.CurrentWeather.WindSpeed
	result.CurrentWeather.WeatherInfo = wmoCodes[resultforecst.CurrentWeather.WeatherCode]

	curday := time.Now().Day()
	index := filter(resultforecst.HourlyTimes, func(t time.Time) bool {
		return t.Day() == curday && t.Hour() >= timeAt && t.Hour() <= timeTo
	})
	result.Hourly = make([]hourlyinfo, len(index))
	for p, i := range index {
		var value hourlyinfo
		value.Temperature = resultforecst.HourlyMetrics["temperature_2m"][i]
		value.WeatherInfo = wmoCodes[resultforecst.HourlyMetrics["weathercode"][i]]
		value.Rain = resultforecst.HourlyMetrics["rain"][i]
		value.WindSpeed = resultforecst.HourlyMetrics["windspeed_10m"][i]
		value.Time = resultforecst.HourlyTimes[i]
		result.Hourly[p] = value
	}

	return &result, nil

}

func filter(vs []time.Time, f func(time.Time) bool) []int {
	vsf := make([]int, 0)
	for i, v := range vs {
		if f(v) {
			vsf = append(vsf, i)
		}
	}
	return vsf
}
