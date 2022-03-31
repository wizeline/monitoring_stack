package main

import (
	"context"
	"log"
	"net/http"
	"time"

/*
// Once I learn go, I need to change this again, to cycle through locations -
// perhaps with import "strconv"
// s := strconvItoa(97)  // convert int to string
// 3128760:4259418:5879400:4907985:4043416:1269843
// Barcelona, Indianapolis, IN : Anchorage, AK : Rockton, IL : Guam : Hyderabad, India
*/

	owm "github.com/briandowns/openweathermap"
	"github.com/caarlos0/env"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Config stores the parameters used to fetch the data
type Config struct {
	pollingInterval time.Duration
	requestTimeout  time.Duration
	Port 			string `env:"OWM_PORT" envDefault:":2112"`
	APIKey          	string `env:"OWM_API_KEY"`
	Location 		int    `env:"OWM_LOCATION" envDefault:"4259418"`
	Duration        	int    `env:"OWM_DURATION" envDefault:"60"`
}

func loadMetrics(ctx context.Context, location int) <-chan error {
	errC := make(chan error)
	go func() {
		c := time.Tick(cfg.pollingInterval)
		for {
			select {
			case <-ctx.Done():
				return // returning not to leak the goroutine
			case <-c:
				client := &http.Client{
					Timeout: cfg.requestTimeout,
				}

				w, err := owm.NewCurrent("C", "EN", cfg.APIKey, owm.WithHttpClient(client)) // (internal - OpenWeatherMap reference for fahrenheit) with English output
				if err != nil {
					errC <- err
					continue
				}

				err = w.CurrentByID(location)
				if err != nil {
					errC <- err
					continue
				}

				temp.WithLabelValues(w.Name).Set(w.Main.Temp)

				pressure.WithLabelValues(w.Name).Set(w.Main.Pressure)

				humidity.WithLabelValues(w.Name).Set(float64(w.Main.Humidity))

				wind.WithLabelValues(w.Name).Set(w.Wind.Speed)

				wind_dir.WithLabelValues(w.Name).Set(w.Wind.Deg)

				clouds.WithLabelValues(w.Name).Set(float64(w.Clouds.All))

				rain.WithLabelValues(w.Name).Set(w.Rain.ThreeH)

				timezone.WithLabelValues(w.Name).Set(float64(w.Timezone))

				sunrise.WithLabelValues(w.Name).Set(float64(w.Sys.Sunrise))

				sunset.WithLabelValues(w.Name).Set(float64(w.Sys.Sunset))

				var scraped_weather = w.Weather[0].Description
				if scraped_weather ==  last_weather {
					weather.WithLabelValues(w.Name, scraped_weather).Set(1)
				} else {
					weather.WithLabelValues(w.Name, scraped_weather).Set(1)
					weather.WithLabelValues(w.Name, last_weather).Set(0)
					last_weather = scraped_weather
				}
				log.Println(w.Weather[0].Description)
				log.Println("scraping OK for ", w.Name)
			}
		}
	}()
	return errC
}

var (
	cfg = Config{
		pollingInterval: 5 * time.Second,
		requestTimeout:  1 * time.Second,
	}
	temp = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "openweathermap",
		Name:      "temperature_celsius",
		Help:      "Temperature in °C",
	}, []string{"location"})

	pressure = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "openweathermap",
		Name:      "pressure_hpa",
		Help:      "Atmospheric pressure in hPa",
	}, []string{"location"})

	humidity = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "openweathermap",
		Name:      "humidity_percent",
		Help:      "Humidity in Percent",
	}, []string{"location"})

	wind = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "openweathermap",
		Name:      "wind_mph",
		Help:      "Wind speed in mph",
	}, []string{"location"})

	wind_dir = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "openweathermap",
		Name:      "wind_direction",
		Help:      "Wind direction in degrees",
	}, []string{"location"})

	clouds = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "openweathermap",
		Name:      "cloudiness_percent",
		Help:      "Cloudiness in Percent",
	}, []string{"location"})

	rain = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "openweathermap",
		Name:      "rain",
		Help:      "Rain contents 3h",
	}, []string{"location"})

	sunrise = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "openweathermap",
		Name:      "unix_sunrise",
		Help:      "Sunrise in Unix time",
	}, []string{"location"})

	sunset = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "openweathermap",
		Name:      "unix_sunset",
		Help:      "Sunset in Unix time",
	}, []string{"location"})

	timezone = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "openweathermap",
		Name:      "city_timezone",
		Help:      "Timezone of the city / location",
	}, []string{"location"})

	city_name = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "openweathermap",
		Name:      "city_name",
		Help:      "Name of the city / location",
	}, []string{"location"})

	weather = prometheus.NewGaugeVec(prometheus.GaugeOpts{
    	Namespace: 	"openweathermap",
        Name: 		"weather",
        Help: 		"The weather label.",
    }, []string{"location", "weather"})

    last_weather = ""
)

func main() {

	env.Parse(&cfg)
	prometheus.Register(temp)
	prometheus.Register(pressure)
	prometheus.Register(humidity)
	prometheus.Register(wind)
	prometheus.Register(clouds)
	prometheus.Register(wind_dir)
	prometheus.Register(sunrise)
	prometheus.Register(sunset)
	prometheus.Register(timezone)
	prometheus.Register(city_name)
	prometheus.Register(weather)

	errC := loadMetrics(context.TODO(), cfg.Location)
	go func() {
		for err := range errC {
			log.Println(err)
		}
	}()
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(cfg.Port, nil)
}
