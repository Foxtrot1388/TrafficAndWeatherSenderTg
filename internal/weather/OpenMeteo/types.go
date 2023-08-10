package weather

import (
	"fmt"
	"strings"
	"time"
)

var wmoCodes = map[float64]string{
	0:  "Чистое небо",
	1:  "Преимущественно ясно",
	2:  "Переменная облачность",
	3:  "Пасмурно",
	45: "Туман",
	48: "изморозь",
	51: "Морось легкая",
	53: "Морось умеренная",
	55: "Морось плотная",
	56: "Ледяная морось легкая",
	57: "Ледяная морось плотная",
	61: "Дождь слабый",
	63: "Дождь умеренный",
	65: "Дождь сильный",
	66: "Ледяной дождь легкая интенсивность",
	67: "Ледяной дождь сильная интенсивность",
	71: "Снегопад слабая интенсивность",
	73: "Снегопад умеренная интенсивность",
	75: "Снегопад сильная интенсивность",
	77: "Снежные зерна",
	80: "Ливневые дожди слабые",
	81: "Ливневые дожди умеренные",
	82: "Ливневые дожди сильные",
	85: "Слабый снег",
	86: "Сильный снег",
	// Thunderstorm forecast with hail is only available in Central Europe
	95: "Гроза",
	96: "Гроза со слабым градом",
	99: "Гроза с сильным градом",
}

type WeatherInfo struct {
	CurrentWeather struct {
		Temperature float64
		WeatherInfo string
		WindSpeed   float64
	}
	Hourly            []hourlyinfo
	TemperatureUnit   string
	WindspeedUnit     string
	PrecipitationUnit string
}

type hourlyinfo struct {
	Time        time.Time
	WeatherInfo string
	Temperature float64
	WindSpeed   float64
	Rain        float64
}

func (w WeatherInfo) String() string {

	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("Сейчас %s температура %g%s ветер %g%s",
		w.CurrentWeather.WeatherInfo,
		w.CurrentWeather.Temperature,
		w.TemperatureUnit,
		w.CurrentWeather.WindSpeed,
		w.WindspeedUnit))

	for i := 0; i < len(w.Hourly); i++ {
		sb.WriteString(fmt.Sprintf("\r\nНа %s %s температура %g%s ветер %g%s осадки %g%s",
			w.Hourly[i].Time.Format("15:04"),
			w.Hourly[i].WeatherInfo,
			w.Hourly[i].Temperature,
			w.TemperatureUnit,
			w.Hourly[i].WindSpeed,
			w.WindspeedUnit,
			w.Hourly[i].Rain,
			w.PrecipitationUnit))
	}

	return sb.String()

}
