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

var globalData = make(map[string][]int)
var globalDataLock sync.RWMutex
var filters []string

func main() {
	filters = os.Args[1:]
	// Start HTTP server in a goroutine
	go startHTTPServer()

	reader := bufio.NewReader(os.Stdin)
	batch := map[string][]int{}
	var batchLock sync.Mutex
	go func() {
		for {
			line, _ := reader.ReadString('\n')
			trimmed := strings.Trim(line, "\n")
			words := strings.Split(trimmed, " ")
			if words[len(words)-1] == "micros" || words[len(words)-1] == "millis" || words[len(words)-1] == "ms" || words[len(words)-1] == "µs" {
				key := strings.Join(words[:len(words)-3], " ")
				val, err := strconv.Atoi(words[len(words)-2])
				if err != nil {
					continue
				}
				batchLock.Lock()
				batch[key] = append(batch[key], val)
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
			batch = map[string][]int{}
			batchLock.Unlock()
		}
	}
}

func filterDict(predicate func(string, []int) bool, dictObj map[string][]int) map[string][]int {
	result := make(map[string][]int)
	for key, value := range dictObj {
		if predicate(key, value) {
			result[key] = value
		}
	}
	return result
}

func filtersContainKey(key string, value []int) bool {
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
                        barmode: 'stack',
                        bargap: 0.1,
                        bargroupgap: 0.1,
                        paper_bgcolor: '#333',
                        plot_bgcolor: '#333',
                        font: { color: '#ccc' }
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
