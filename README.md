# loghist
logs -> stdin -> live updating histograms

Usage:

`go run main.go | go run aggregate.go | python hist.py`

main.go is a stand in for any program you want to plot the timing logs of.

You can add filters to the plotted histograms. 

Try:
`go run main.go | go run aggregate.go | python hist.py fast slow`
