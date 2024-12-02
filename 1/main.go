package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func getSortedArrays(filename string) (left []int, right []int) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error trying to open input.txt: %v", err)

	}
	// Defer closing the file
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Error trying to close file", err)
		}
	}(file)

	var separator string = "   "
	leftNumbers := make([]int, 0)
	rightNumbers := make([]int, 0)

	// Create a scanner
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, separator)
		left := parts[0]
		right := strings.TrimRight(parts[1], "\n")

		if len(left) == 0 || len(right) == 0 {
			fmt.Printf("Error: line '%s' does not split correctly\n", line)
			os.Exit(1)
		} else {
			fmt.Printf("Read '%s' as left and '%s' as right\n", left, right)
		}
		// read the line and split it.
		leftNum, err := strconv.Atoi(left)
		if err != nil {
			fmt.Printf("Error converting left to integer: %v\n", err)
			os.Exit(1)
		}
		rightNum, err := strconv.Atoi(right)
		if err != nil {
			fmt.Printf("Error converting right to integer: %v\n", err)
			os.Exit(1)
		}

		leftNumbers = append(leftNumbers, leftNum)
		rightNumbers = append(rightNumbers, rightNum)
	}

	// Once arrays are built, sort and return
	sort.Ints(leftNumbers)
	sort.Ints(rightNumbers)
	return leftNumbers, rightNumbers
}

// Go's native ABS function expects and returns float64, this is cheaper due to lack of casting
func easyAbs(number int) int {
	if number >= 0 {
		return number
	} else {
		return number * -1
	}
}

// count the occurrences of each number in a number slice in go, return a map mapping each number to occurrences
func countOccurrences(numbers []int) map[int]int {

	occurrences := make(map[int]int)
	for i := 0; i < len(numbers); i++ {

		if value, exists := occurrences[numbers[i]]; exists {
			occurrences[numbers[i]] = value + 1
		} else {
			occurrences[numbers[i]] = 1
		}
	}
	return occurrences
}

// Calculate the complex distance
func calculateComplexDistance(left []int, occurrences map[int]int) int {
	fmt.Printf("Occurrences: %v\n", occurrences)

	var totalDistance int = 0
	for _, number := range left {

		count := occurrences[number]
		totalDistance += number * count
	}

	return totalDistance
}

func main() {
	leftColSorted, rightColSorted := getSortedArrays("input.txt")

	var totalDifference int = 0
	for i := 0; i < len(leftColSorted); i++ {

		totalDifference += easyAbs(leftColSorted[i] - rightColSorted[i])
	}

	fmt.Println("Total distance:", totalDifference)

	occurrences := countOccurrences(rightColSorted)
	complexDistance := calculateComplexDistance(leftColSorted, occurrences)
	fmt.Println("Complex distance:", complexDistance)
}
