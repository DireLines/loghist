package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var globalData = make(map[string][]float64)
var globalDataLock sync.RWMutex
var filters []string

var timeUnits map[string]float64 = map[string]float64{
	"nanos":        1.0 / 1000.0,
	"nanoseconds":  1.0 / 1000.0,
	"ns":           1.0 / 1000.0,
	"micros":       1.0,
	"microseconds": 1.0,
	"µs":           1.0,
	"millis":       1000,
	"milliseconds": 1000,
	"ms":           1000,
	"seconds":      1000000,
	"s":            1000000,
	"sec":          1000000,
	"secs":         1000000,
}

func main() {
	filters = os.Args[1:]
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	// Start HTTP server in a goroutine
	go startHTTPServer(port)

	reader := bufio.NewReader(os.Stdin)
	batch := map[string][]float64{}
	var batchLock sync.Mutex
	go func() {
		for {
			line, _ := reader.ReadString('\n')
			trimmed := strings.Trim(line, "\n")
			words := strings.Split(trimmed, " ")
			timeUnit := words[len(words)-1]
			timeUnitMultiplier, isValidUnit := timeUnits[timeUnit]
			if isValidUnit {
				key := strings.Join(words[:len(words)-3], " ")
				val, err := strconv.ParseFloat(words[len(words)-2], 64)
				if err != nil {
					continue
				}
				batchLock.Lock()
				batch[key] = append(batch[key], val*timeUnitMultiplier)
				batchLock.Unlock()
			}
		}
	}()

	for {
		time.Sleep(time.Millisecond * 20)
		if len(batch) > 0 {
			batchLock.Lock()
			if len(filters) > 0 {
				batch = filterDict(filtersContainKey, batch)
			}
			globalDataLock.Lock()
			for k, v := range batch {
				globalData[k] = append(globalData[k], v...)
			}
			globalDataLock.Unlock()
			batch = map[string][]float64{}
			batchLock.Unlock()
		}
	}
}

func filterDict(predicate func(string, []float64) bool, dictObj map[string][]float64) map[string][]float64 {
	result := make(map[string][]float64)
	for key, value := range dictObj {
		if predicate(key, value) {
			result[key] = value
		}
	}
	return result
}

func filtersContainKey(key string, value []float64) bool {
	for _, filter := range filters {
		if strings.Contains(key, filter) {
			return true
		}
	}
	return false
}

func startHTTPServer(port string) {
	http.HandleFunc("/", serveHTML)
	http.HandleFunc("/data", serveData)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error starting HTTP server: %v\n", err)
		os.Exit(1)
	}
}

func serveHTML(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, htmlPage)
}

func serveData(w http.ResponseWriter, r *http.Request) {
	globalDataLock.RLock()
	defer globalDataLock.RUnlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(globalData)
}

const htmlPage = `
<!DOCTYPE html>
<html>
<head>
    <title>Live Histogram</title>
    <script src="https://cdn.plot.ly/plotly-latest.min.js"></script>
    <style>
        body {
            background-color: #333;
            color: #ccc;
            font-family: Arial, sans-serif;
        }
        #plot {
            width: 100%;
            height: 100vh;
        }
    </style>
</head>
<body>
    <div id="plot"></div>
    <script>
        function fetchDataAndUpdatePlot() {
            fetch('/data')
                .then(response => response.json())
                .then(data => {
                    var traces = [];
                    var keys = Object.keys(data);
                    for (var i = 0; i < keys.length; i++) {
                        var key = keys[i];
                        var values = data[key];
                        var trace = {
                            x: values,
                            type: 'histogram',
                            name: key,
                            opacity: 0.75
                        };
                        traces.push(trace);
                    }
                    var layout = {
                        title: {text: 'distribution of timing logs by task'},
                        barmode: 'stack',
                        bargap: 0.1,
                        bargroupgap: 0.1,
                        paper_bgcolor: '#333',
                        plot_bgcolor: '#333',
                        font: { color: '#ccc' },
                        xaxis: { title: { text: 'execution time (micros)' } },
                        yaxis: { title: { text: 'count of occurrences' } },
						showlegend: true,
                    };
                    Plotly.newPlot('plot', traces, layout);
                })
                .catch(error => console.error('Error fetching data:', error));
        }
        setInterval(fetchDataAndUpdatePlot, 100);
        fetchDataAndUpdatePlot();
    </script>
</body>
</html>
`
