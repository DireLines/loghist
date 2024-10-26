package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	data := map[string][]int{}
	var dataLock sync.Mutex
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
				dataLock.Lock()
				data[key] = append(data[key], val)
				dataLock.Unlock()
			}
		}
	}()

	for {
		time.Sleep(time.Millisecond * 100)
		if len(data) > 0 {
			dataLock.Lock()
			json, _ := json.Marshal(data)
			fmt.Println(string(json))
			data = map[string][]int{}
			dataLock.Unlock()
		}
	}
}
