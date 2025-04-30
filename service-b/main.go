package main

import (
    "encoding/json"
    "errors"
    "fmt"
    "log"
    "net/http"
    "os"
    "regexp"

    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/zipkin"
    "go.opentelemetry.io/otel/sdk/resource"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
    "go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

type WeatherResponse struct {
    City  string  `json:"city"`
    TempC float64 `json:"temp_C"`
    TempF float64 `json:"temp_F"`
    TempK float64 `json:"temp_K"`
}

type ViaCEPResponse struct {
    Localidade string `json:"localidade"`
    Erro       bool   `json:"erro"`
}

type WeatherAPIResponse struct {
    Current struct {
        TempC float64 `json:"temp_c"`
    } `json:"current"`
}

func initTracer() {
    endpoint := os.Getenv("OTEL_EXPORTER_ZIPKIN_ENDPOINT")
    exporter, err := zipkin.New(endpoint)
    if err != nil {
        log.Fatal(err)
    }
    tp := sdktrace.NewTracerProvider(
        sdktrace.WithBatcher(exporter),
        sdktrace.WithResource(resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceNameKey.String("service-b"))),
    )
    otel.SetTracerProvider(tp)
    tracer = otel.Tracer("service-b")
}

func isValidCEP(cep string) bool {
    match, _ := regexp.MatchString(`^\d{8}$`, cep)
    return match
}

func getCityByCEP(cep string) (string, error) {
    _, span := tracer.Start(r.Context(), "ViaCEP Lookup")
    defer span.End()

    resp, err := http.Get(fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep))
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    var data ViaCEPResponse
    if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
        return "", err
    }
    if data.Erro {
        return "", errors.New("can not find zipcode")
    }
    return data.Localidade, nil
}

func getTemperatureByCity(city string) (float64, error) {
    _, span := tracer.Start(r.Context(), "WeatherAPI Lookup")
    defer span.End()

    key := os.Getenv("WEATHER_API_KEY")
    url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s", key, city)

    resp, err := http.Get(url)
    if err != nil {
        return 0, err
    }
    defer resp.Body.Close()

    var data WeatherAPIResponse
    if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
        return 0, err
    }
    return data.Current.TempC, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
    ctx, span := tracer.Start(r.Context(), "Handle Weather by CEP")
    defer span.End()

    cep := r.URL.Path[len("/weather/"):]
    if !isValidCEP(cep) {
        http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
        return
    }

    city, err := getCityByCEP(cep)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    tempC, err := getTemperatureByCity(city)
    if err != nil {
        http.Error(w, "temperature error", http.StatusInternalServerError)
        return
    }

    resp := WeatherResponse{
        City:  city,
        TempC: tempC,
        TempF: tempC*1.8 + 32,
        TempK: tempC + 273,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}

func main() {
    initTracer()
    http.HandleFunc("/weather/", handler)
    log.Println("Service B running on port 8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
