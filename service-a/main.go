package main

import (
    "encoding/json"
    "log"
    "net/http"
    "os"

    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/zipkin"
    "go.opentelemetry.io/otel/sdk/resource"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
    "go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

type CEPInput struct {
    CEP string `json:"cep"`
}

func initTracer() {
    endpoint := os.Getenv("OTEL_EXPORTER_ZIPKIN_ENDPOINT")
    exporter, err := zipkin.New(endpoint)
    if err != nil {
        log.Fatal(err)
    }
    tp := sdktrace.NewTracerProvider(
        sdktrace.WithBatcher(exporter),
        sdktrace.WithResource(resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceNameKey.String("service-a"))),
    )
    otel.SetTracerProvider(tp)
    tracer = otel.Tracer("service-a")
}

func handler(w http.ResponseWriter, r *http.Request) {
    ctx, span := tracer.Start(r.Context(), "Validate and Forward CEP")
    defer span.End()

    if r.Method != http.MethodPost {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var input CEPInput
    err := json.NewDecoder(r.Body).Decode(&input)
    if err != nil || len(input.CEP) != 8 {
        http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
        return
    }

    req, _ := http.NewRequestWithContext(ctx, "GET", "http://service-b:8080/weather/"+input.CEP, nil)
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        http.Error(w, "error contacting service-b", http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    w.WriteHeader(resp.StatusCode)
    w.Header().Set("Content-Type", "application/json")
    json.NewDecoder(resp.Body).Decode(w)
}

func main() {
    initTracer()
    http.HandleFunc("/cep", handler)
    log.Println("Service A running on port 8081")
    log.Fatal(http.ListenAndServe(":8081", nil))
}
