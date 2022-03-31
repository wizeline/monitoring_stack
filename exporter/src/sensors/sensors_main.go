    package main

    import (
        "context"
        "log"
        "net/http"
        "time"
        "github.com/caarlos0/env"
        "github.com/prometheus/client_golang/prometheus"
        "github.com/prometheus/client_golang/prometheus/promauto"
        "github.com/prometheus/client_golang/prometheus/promhttp"
        "encoding/json"
        "fmt"
        "io/ioutil"
    )

    type Sensors struct {
        Temperature float32 `json:"temperature"`
        Humidity int `json:humidity`
//         Light bool `json:light`
        Gasanalog int `json:gasanalog`
        Gasdigital int `json:gasdigital`
    }

    // Config stores the parameters used to fetch the data
    type Config struct {
        pollingInterval time.Duration
        requestTimeout  time.Duration
        Port 			string `env:"SENSORS_PORT" envDefault:":2111"`
        APIUrl          string `env:"SENSORS_API_URL"`
        Duration        int    `env:"SENSORS_DURATION" envDefault:"60"`
    }

    func loadMetrics(ctx context.Context, url string, location int) <-chan error {
        var sensorsOutput Sensors

        errC := make(chan error)
        go func() {
            c := time.Tick(cfg.pollingInterval)
            for {
                select {
                case <-ctx.Done():
                    return // returning not to leak the goroutine
                case <-c:
                    spaceClient := http.Client{
                        Timeout: time.Second * 2, // Timeout after 2 seconds
                    }

                    req, err := http.NewRequest(http.MethodGet, url, nil)
                    if err != nil {


                        log.Fatal(err)
                    }
                    req.Header.Set("User-Agent", "live-sensors-test")

                    res, getErr := spaceClient.Do(req)
                    if getErr != nil {
                        log.Fatal(getErr)
                    }

                    // {"temperature":19.1,"humidity":59,"light":false,"gasanalog":93,"gasdigital":0}
                    if res.Body != nil {
                        defer res.Body.Close()
                    }

                    body, readErr := ioutil.ReadAll(res.Body)
                    if readErr != nil {
                        log.Fatal(readErr)
                    }
                    /*
                    body := []byte(`
                        {"temperature":19.1,"humidity":59,"light":false,"gasanalog":93,"gasdigital":0}
                        `)
                    */
                    jsonErr := json.Unmarshal(body, &sensorsOutput)
                    if jsonErr != nil {
                        log.Fatal(jsonErr)
                    }
                    fmt.Println("Struct is:", sensorsOutput)

                    temp.WithLabelValues(sensor_location).Set(float64(sensorsOutput.Temperature))
                    humidity.WithLabelValues(sensor_location).Set(float64(sensorsOutput.Humidity))
//                     light.WithLabelValues(sensor_location).Set(int64(sensorsOutput.Light))
                    gasa.WithLabelValues(sensor_location).Set(float64(sensorsOutput.Gasanalog))
                    gasd.WithLabelValues(sensor_location).Set(float64(sensorsOutput.Gasdigital))
                    //hardcoded to 7200 to set the same timezone as openweathermap
                    timezone.WithLabelValues("sensors_city_timezone").Set(7200)
                    log.Println("scraping OK for ", sensor_location)
                }
            }
        }()
        return errC
    }

    var (
        sensor_location = "store_sensor_location_willy"
        sensor_location_number = 999999
        cfg = Config{
            pollingInterval: 5 * time.Second,
            requestTimeout:  1 * time.Second,
        }
        temp = promauto.NewGaugeVec(prometheus.GaugeOpts{
            Namespace: "sensors",
            Name:      "temperature_celsius",
            Help:      "Temperature in °C",
        }, []string{"location"})

        humidity = promauto.NewGaugeVec(prometheus.GaugeOpts{
            Namespace: "sensors",
            Name:      "humidity_percent",
            Help:      "Humidity in Percent",
        }, []string{"location"})

        gasa = promauto.NewGaugeVec(prometheus.GaugeOpts{
            Namespace: "sensors",
            Name:      "gas_analog_degrees",
            Help:      "Gas analog in degrees",
        }, []string{"location"})

        gasd = promauto.NewGaugeVec(prometheus.GaugeOpts{
            Namespace: "sensors",
            Name:      "gas_digital_degrees",
            Help:      "Gas digital in degrees",
        }, []string{"location"})

        timezone = promauto.NewGaugeVec(prometheus.GaugeOpts{
            Namespace: "sensors",
            Name:      "city_timezone",
            Help:      "Timezone of the city / location",
        }, []string{"location"})

/*
        light = promauto.NewGaugeVec(prometheus.GaugeOpts{
            Namespace: "sensors",
            Name:      "light_bool",
            Help:      "Light in boolean",
        }, []string{"location"})
 */

    )

    func main() {

        if err := env.Parse(&cfg); err != nil {
            fmt.Printf("%+v\n", err)
        }
        prometheus.Register(temp)
        prometheus.Register(humidity)
        prometheus.Register(gasa)
        prometheus.Register(gasd)
        prometheus.Register(timezone)

        errC := loadMetrics(context.TODO(), cfg.APIUrl, sensor_location_number)
        go func() {
            for err := range errC {
                log.Println(err)
                fmt.Printf("%+v\n", err)
            }
        }()
        http.Handle("/metrics", promhttp.Handler())

        if err := http.ListenAndServe(cfg.Port, nil); err != nil {
            log.Printf("Error occur when start server %v", err)
            fmt.Printf("Error occur when start server %v", err)
        }
    }
