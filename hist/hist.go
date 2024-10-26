package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

var data = make(map[string][]float64)
var dataMutex sync.RWMutex
var filters []string

func main() {
	filters = os.Args[1:]

	// Start HTTP server in a goroutine
	go startHTTPServer()

	// Read from stdin
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
				continue
			}
		}

		var batch map[string][]float64
		err = json.Unmarshal(line, &batch)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing JSON: %v\n", err)
			continue
		}

		if len(filters) > 0 {
			batch = filterDict(filtersContainKey, batch)
		}

		dataMutex.Lock()
		for k, v := range batch {
			data[k] = append(data[k], v...)
		}
		dataMutex.Unlock()
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

func startHTTPServer() {
	http.HandleFunc("/", serveHTML)
	http.HandleFunc("/data", serveData)
	err := http.ListenAndServe(":8080", nil)
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
	dataMutex.RLock()
	defer dataMutex.RUnlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

const htmlPage = `
<!DOCTYPE html>
<html>
<head>
    <title>Live Histogram</title>
    <script src="https://cdn.plot.ly/plotly-latest.min.js"></script>
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
                        barmode: 'stack',
                        bargap: 0.1,
                        bargroupgap: 0.1
                    };
                    Plotly.newPlot('plot', traces, layout);
                })
                .catch(error => console.error('Error fetching data:', error));
        }
        setInterval(fetchDataAndUpdatePlot, 1000);
        fetchDataAndUpdatePlot();
    </script>
</body>
</html>
`
