package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Get readings from the file
func getReadings(fileName string) [][]int {

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Failed to close file:", err)
			os.Exit(1)
		}
	}(file)

	var reports [][]int = make([][]int, 0)

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		report := scanner.Text()
		report = strings.TrimRight(report, "\n")

		// Split the line on spaces
		readings := strings.Split(report, " ")
		numericalReadings := make([]int, 0)
		for _, reading := range readings {
			numericalReading, err := strconv.Atoi(reading)
			if err != nil {
				fmt.Printf("Error converting %s to number:%s", numericalReading, err)
			}
			numericalReadings = append(numericalReadings, numericalReading)
		}

		// Once parsed, add the readings to the report
		reports = append(reports, numericalReadings)
	}
	// Return the slice of slices
	return reports

}

// Go's native ABS function expects and returns float64, this is cheaper due to lack of casting
func easyAbs(number int) int {
	if number >= 0 {
		return number
	} else {
		return number * -1
	}
}

func worker(reading []int) bool {

	increasing := reading[1] > reading[0]

	for i := 1; i < len(reading); i++ {
		readingDistance := easyAbs(reading[i] - reading[i-1])
		//fmt.Printf("\tReading distance between %d and %d: %d\n", reading[i], reading[i-1], readingDistance)
		// if we are increasing and one reading is LESS THAN the one before it
		if increasing && reading[i] < reading[i-1] {
			return false
		}
		// if we are DECREASING and one reading is MORE than the one before it.
		if !increasing && reading[i] > reading[i-1] {
			return false

		}
		// if the reading distance is < 1 or greater than 4
		if readingDistance < 1 || readingDistance > 3 {
			return false
		}
	}

	return true
}

// Worker to determine if each one is incrementing or decrementing
func asyncWorker(done chan<- bool, reading []int) {
	done <- worker(reading)
}

func main() {
	fmt.Println("Welcome to day 2!")
	readings := getReadings("input.txt")
	//fmt.Println("Readings:", readings)

	start := time.Now()
	var safeReadings int = 0

	numWorkers := len(readings)
	// Create a buffered channel with as many spots as numbers of workers to prevent blocking
	// otherwise we can get blocked by processing, can slow us down a lot
	readingsChan := make(chan bool, numWorkers)
	//var wg sync.WaitGroup

	// Dispatch goroutines
	for _, reading := range readings {
		go asyncWorker(readingsChan, reading)
	}

	// read until all workers have returned to the channel
	for i := 0; i < numWorkers; i++ {
		safe := <-readingsChan
		if safe {
			safeReadings++
		}
	}

	end := time.Now()
	fmt.Println("Safe readings:", safeReadings)
	fmt.Println("Finished in:", end.Sub(start))
	close(readingsChan)
}
