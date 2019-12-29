package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	flag "github.com/spf13/pflag"
)

func main() {
	var dataFile string
	//var basePattern = []int{0, 1, 0, -1}

	flag.StringVarP(&dataFile, "data file name", "f", "", "")
	flag.Parse()

	input, err := buildList(dataFile)
	if err != nil {
		log.Fatal("Failed to get input file!", err)
	}

	//result := FFT(input, basePattern, 100)
	//fmt.Println("result=", result)

	message := FFTExtract(input, 10000, 100)
	fmt.Println("message=", message)

}

/* Extract 8 number message from the bottom half of FFT result.
 * The bottom half of FFT pattern matrix is special with all 0 followed by all 1.
 * As such, we can calculate messages residing in this bit position.
 */
func FFTExtract(input []int, repeatTimes int, iteration int) []int {
	var skipBits int
	var realInput []int

	for i := 0; i < 7; i++ {
		skipBits *= 10
		skipBits += input[i]
	}
	fmt.Println("skipBits=", skipBits)
	l := len(input)*repeatTimes - skipBits
	fmt.Println("len of l=", l)

	times := l / len(input)
	start := len(input) - l%len(input)

	fmt.Println("start=", start)
	realInput = append(realInput, input[start:]...)

	for i := 0; i < times; i++ {
		realInput = append(realInput, input...)
	}

	fmt.Println("input=", input)
	fmt.Println("len of realInput=", len(realInput))

	for i := 0; i < iteration; i++ {
		output := make([]int, l)
		sum := 0
		for _, d := range realInput {
			sum += d
		}
		output[0] = sum % 10
		for j := 1; j < l; j++ {
			sum = sum - realInput[j-1]
			output[j] = sum % 10
		}
		realInput = output
	}

	return realInput[:8]
}

func FFT(input []int, basePattern []int, iteration int) []int {
	for i := 0; i < iteration; i++ {
		//fmt.Printf("[%v]:%v\n", i, input)
		output := make([]int, len(input))
		for j := 0; j < len(input); j++ {
			pattern := calculatePattern(basePattern, j+1, len(input))
			sum := sumMultiply(input, pattern)
			output[j] = sum
		}
		input = output
	}
	return input
}

func sumMultiply(input []int, pattern []int) int {
	if len(input) != len(pattern) {
		log.Fatal("Len of input and pattern not matched!")
	}

	sum := 0
	//fmt.Println("input:", input)
	//fmt.Println("pattern:", pattern)
	for i, v := range input {
		sum += v * pattern[i]
	}
	digit := sum % 10
	return int(math.Abs(float64(digit)))
}

func calculatePattern(basePattern []int, repeat int, plen int) []int {
	baseLen := len(basePattern)

	pattern := make([]int, plen+1)
	for i, _ := range pattern {
		j := (i / repeat) % baseLen
		pattern[i] = basePattern[j]
	}
	return pattern[1:]
}

func buildList(dataFile string) ([]int, error) {
	var input []int

	f, err := os.Open(dataFile)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	buff := bufio.NewReader(f)
	line, err := buff.ReadString('\n')
	if err != nil {
		return input, err
	}
	line = strings.TrimSuffix(line, "\n")
	for _, l := range line {
		num, err := strconv.Atoi(string(l))
		if err != nil {
			log.Fatal(err)
			os.Exit(2)
		}
		input = append(input, num)
	}
	return input, err
}
