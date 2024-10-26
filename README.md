# loghist
logs -> stdin -> live updating histograms

Usage:

`go run example/example.go | go run main.go`

`example/example.go` is a stand in for any program you want to plot the timing logs of.

`main.go` will capture the log output of the program and serve a page charting it at localhost:8080

You can add filters to the plotted histograms. 

Try:
`go run example/example.go | go run main.go fast slow`

Would recommend aliasing for quick invocation:

```
alias loghist="open http://localhost:8080 && go run [absolute path to main.go]"
go run example/example.go | loghist
```
