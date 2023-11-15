package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	data := map[string][]int{}
	go func() {
		for {
			line, _ := reader.ReadString('\n')
			trimmed := strings.Trim(line, "\n")
			words := strings.Split(trimmed, " ")
			if words[len(words)-1] == "micros" || words[len(words)-1] == "millis" {
				key := strings.Join(words[:len(words)-3], " ")
				val, err := strconv.Atoi(words[len(words)-2])
				if err != nil {
					continue
				}
				data[key] = append(data[key], val)
			}
		}
	}()

	for {
		time.Sleep(time.Millisecond * 100)
		if len(data) > 0 {
			json, _ := json.Marshal(data)
			fmt.Println(string(json))
			data = map[string][]int{}
		}
	}
}
