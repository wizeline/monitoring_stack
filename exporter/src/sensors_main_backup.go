    package main

    import (
        "encoding/json"
        "fmt"
        "io/ioutil"
        "log"
        "net/http"
        "time"
    )

    type Sensors struct {
        Temperature float32 `json:"temperature"`
        Humidity int `json:humidity`
        Light bool `json:light`
        Gasanalog int `json:gasanalog`
        Gasdigital int `json:gasdigital`
    }

    func main() {
        var sensorsOutput Sensors
     	//url := string env.Parse(`env:"SENSORS_API_URL" envDefault:"http://walii.dynu.net/sensors`)
     	//url := "http://walii.dynu.net/sensors"
        url := "http://api.open-notify.org/astros.json"

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

        bodyOriginal, readErr := ioutil.ReadAll(res.Body)
        if readErr != nil {
            log.Fatal(readErr)
        }
        body := []byte(`
            {"temperature":19.1,"humidity":59,"light":false,"gasanalog":93,"gasdigital":0}
            `)

        jsonErr := json.Unmarshal(body, &sensorsOutput)
        if jsonErr != nil {
            log.Fatal(jsonErr)
        }
        fmt.Println("Struct is:", sensorsOutput)
        fmt.Println(sensorsOutput.Temperature)
        fmt.Println(sensorsOutput.Humidity)
        fmt.Println(sensorsOutput.Light)
        fmt.Println(sensorsOutput.Gasanalog)
        fmt.Println(sensorsOutput.Gasdigital)
        fmt.Println(bodyOriginal)
    }