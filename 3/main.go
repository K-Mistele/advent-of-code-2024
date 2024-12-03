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

func main() {
	fmt.Println("Day 3!")

	fileContents := getProblemString("input.txt")
	fmt.Println("File contents:", fileContents)

	expressions := findMulExpressions(fileContents)

	// Calculate the number of goroutines that will be needed and set up a buffered channel of the same size
	routines := len(expressions)
	productsChannel := make(chan int, routines)

	// Launch routines
	for _, expression := range expressions {
		go worker(productsChannel, expression)
	}

	// Do the summation by pulling from the channel
	var sumOfProducts int = 0
	for i := 0; i < routines; i++ {
		sum := <-productsChannel
		sumOfProducts += sum
	}

	fmt.Println("Sum of products:", sumOfProducts)

}
