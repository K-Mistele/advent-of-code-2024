package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func getProblemString(filename string) string {

	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Unable to open file %s: %s\n", filename, err)
		os.Exit(1)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			fmt.Println("Error closing file %s: %s\n", filename, err)
			os.Exit(1)
		}
	}()

	fileContents, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("Unable to read file %s: %s\n", filename, err)
		os.Exit(1)
	}
	contentsString := string(fileContents)
	return contentsString
}

func findMulExpressions(haystack string) []string {

	needle := `mul\(\d{1,3},\d{1,3}\)`
	re, err := regexp.Compile(needle)
	if err != nil {
		fmt.Println("Error: unable to compile regex.", err)
		os.Exit(1)
	}

	matches := re.FindAllString(haystack, -1)
	return matches
}

// handleMultiplicationExpression handles each of the mul(X,Y) expressions
func handleMultiplicationExpression(expression string) int {
	intermediate := strings.Replace(expression, "mul(", "", 1)
	delimitedNumbers := strings.Replace(intermediate, ")", "", 1)
	operands := strings.Split(delimitedNumbers, ",")

	left, err := strconv.Atoi(operands[0])
	if err != nil {
		fmt.Printf("Error converting operand %s to string: %s\n", operands[0], err)
		return -9999999999
	}
	right, err := strconv.Atoi(operands[1])
	if err != nil {
		fmt.Printf("Error converting operand %s to string: %s\n", operands[1], err)
		return -99999999999
	}

	return left * right
}

// wrapper for a goroutine that will execute these in parallel.
func worker(products chan<- int, expression string) {
	product := handleMultiplicationExpression(expression)
	products <- product
}

// reprocessProblemWithToggles re-does the regex parsing, but also includes do() and don't() toggles
func reprocessProblemWithToggles(input string) []string {

	pattern := `(?:mul\(\d{1,3},\d{1,3}\))|(?:do\(\))|(?:don't\(\))`
	re := regexp.MustCompile(pattern)
	matches := re.FindAllString(input, -1)

	// Go through and remove expressions that are between "don't()" and the next "do()"
	var shouldDo bool = true

	expressionsToDo := make([]string, 0)

	for _, match := range matches {

		// Toggle do / don't
		if match == "don't()" {
			shouldDo = false
		} else if match == "do()" {
			shouldDo = true

			// for regular matches, add them to the list only if we're in "do" mode.
		} else {
			if shouldDo {
				expressionsToDo = append(expressionsToDo, match)
			}
		}
	}

	return expressionsToDo
}

func handleWorkAsynchronously(workqueue []string) int {
	// Calculate the number of goroutines that will be needed and set up a buffered channel of the same size
	routines := len(workqueue)
	productsChannel := make(chan int, routines)

	// Launch routines
	for _, expression := range workqueue {
		go worker(productsChannel, expression)
	}

	// Do the summation by pulling from the channel
	var sumOfProducts int = 0
	for i := 0; i < routines; i++ {
		sum := <-productsChannel
		sumOfProducts += sum
	}
	return sumOfProducts
}

func main() {
	fmt.Println("Day 3!")

	fileContents := getProblemString("input.txt")
	fmt.Println("File contents:", fileContents)

	expressions := findMulExpressions(fileContents)

	sumOfProducts := handleWorkAsynchronously(expressions)

	fmt.Println("Sum of products:", sumOfProducts)

	revisedExpressions := reprocessProblemWithToggles(fileContents)
	sumOfProducts = handleWorkAsynchronously(revisedExpressions)
	fmt.Println("Revised sum:", sumOfProducts)
}
