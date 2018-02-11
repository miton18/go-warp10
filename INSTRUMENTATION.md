# Instrumentation

Create a counter on a passive exporter

```go
    pe, err := w.NewPassiveExporter("127.0.0.1:9100", "/metrics")
    if err != nil {
        fmt.Println(err)
    }

    fmt.Println("new counter")
    c := w.NewMetricCounter("http.calls", w.Labels{
        "code": "500",
    }, "the number of HTTP calls on the API")

    // Add the counter to the exporter
    pe.Register(c)

    go func() {
        t := time.NewTicker(time.Second)
        for {
            select {
            case <-t.C:
                c.Add(1)
            }
        }
    }()
```
