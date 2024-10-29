# loghist

logs -> stdin -> live updating histograms

Usage:

`go run example/example.go | go run main.go`

`example/example.go` is a stand in for any program you want to plot the timing logs of.

Logs are expected on stdout and should follow the format `[any string describing a task] took [number] [time unit]`

Ex: `update: physics took 10.387667 millis`

Current valid strings for time units:
`nanos`
`nanoseconds`
`ns`
`micros`
`microseconds`
`µs`
`millis`
`milliseconds`
`ms`
`seconds`
`s`
`sec`
`secs`

`main.go` will capture the timing log output of the program and serve a page charting it as a stacked histogram at `localhost:8080`

Any command line args to `main.go` will be interpreted as filters on the resulting histogram. If supplied, only task names matching at least one filter will be shown.

Ex:
`go run example/example.go | go run main.go fast slow`
