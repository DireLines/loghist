# loghist
logs -> stdin -> live updating histograms

Usage:

`go run example/example.go | go run main.go`

`example/example.go` is a stand in for any program you want to plot the timing logs of.

Logs are expected on stdout and should follow the format `[any string describing a task] took [integer] [time unit]`

Current valid strings for time units:  `millis`, `ms`, `micros`, and `µs`

`main.go` will capture the timing log output of the program and serve a page charting it at localhost:8080

You can add filters to the plotted histograms. 

Try:
`go run example/example.go | go run main.go fast slow`
