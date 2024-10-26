package example

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// a program that produces a stream of timing logs
func example() {
	for {
		arr := []int{}
		arr_size := 10000
		logElapsedTime := elapsedTimeLogger(fmt.Sprintf("%d elements", arr_size))
		for i := 0; i < arr_size; i++ {
			arr = append(arr, rand.Intn(100))
		}
		logElapsedTime("initializing array")

		sum := merge_serial(arr)
		fmt.Println("sum:", sum)
		logElapsedTime("merge_serial")

		sum = merge_parallel(arr)
		fmt.Println("sum:", sum)
		logElapsedTime("merge_parallel")
	}
}

// some slow merge operation on two array elements
func slow_combine(a int, b int) int {
	waste_time(1)
	return a + b
}

// just a plain for loop
func merge_serial(arr []int) int {
	sum := 0
	for _, n := range arr {
		sum = slow_combine(sum, n)
	}
	return sum
}

// using recursion to solve subproblems in parallel
func merge_parallel(arr []int) int {
	if len(arr) < 100 {
		return merge_serial(arr)
	}
	middle_index := len(arr) / 2
	var wg sync.WaitGroup
	left := 0
	right := 0
	wg.Add(2)
	go func() {
		defer wg.Done()
		left = merge_parallel(arr[:middle_index])
	}()
	go func() {
		defer wg.Done()
		right = merge_parallel(arr[middle_index:])
	}()
	wg.Wait()
	return slow_combine(left, right)
}

func elapsedTimeLogger(namespace string) func(string) {
	start := time.Now().UnixMicro()
	return func(msg string) {
		now := time.Now().UnixMicro()
		fmt.Println(fmt.Sprintf("%s: %s took %d micros", namespace, msg, (now - start)))
		start = time.Now().UnixMicro()
	}
}

func waste_time(micros int64) {
	start := time.Now().UnixMicro()
	//waste some time
	for time.Now().UnixMicro()-start < micros {
	}
}
