package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

const (
	DataFile                   = "loremipsum.txt"
	amountOfParellellRoutinies = 20
)

// Return the word frequencies of the text argument.
//
// Split load optimally across processor cores.
func WordCount(text string) map[string]int {
	freq := make(map[string]int)
	ch := make(chan map[string]int)
	words := strings.Fields(text)
	stepSize := len(words) / amountOfParellellRoutinies

	//telling workers to work
	for i := 0; i < amountOfParellellRoutinies; i++ {
		go countForEachSlice(words[i*stepSize:(i+1)*stepSize], ch)
	}
	// if stepSize has dropped floating value, we need to calculate those parts aswell.
	go countForEachSlice(words[amountOfParellellRoutinies*stepSize:], ch)

	//getting the "profits"
	for i := 0; i < amountOfParellellRoutinies+1; i++ {
		profit := <-ch

		for key, val := range profit {
			freq[key] += val
		}
	}
	return freq
}

func countForEachSlice(words []string, ch chan<- map[string]int) {
	sendOver := make(map[string]int)
	for _, word := range words {
		sendOver[strings.TrimRight(strings.ToLower(word), ".,")] += 1
	}
	ch <- sendOver
}

// Benchmark how long it takes to count word frequencies in text numRuns times.
//
// Return the total time elapsed.
func benchmark(text string, numRuns int) int64 {
	start := time.Now()
	for i := 0; i < numRuns; i++ {
		WordCount(text)
	}
	runtimeMillis := time.Since(start).Nanoseconds() / 1e6

	return runtimeMillis
}

// Print the results of a benchmark
func printResults(runtimeMillis int64, numRuns int) {
	fmt.Printf("amount of runs: %d\n", numRuns)
	fmt.Printf("total time: %d ms\n", runtimeMillis)
	average := float64(runtimeMillis) / float64(numRuns)
	fmt.Printf("average time/run: %.2f ms\n", average)
}

func main() {
	// read in DataFile as a string called data
	dataInBytes, err := os.ReadFile(DataFile)
	if err != nil {
		log.Fatal(err)
	}
	data := string(dataInBytes)

	numRuns := 1000
	runtimeMillis := benchmark(string(data), numRuns)
	printResults(runtimeMillis, numRuns)
}
