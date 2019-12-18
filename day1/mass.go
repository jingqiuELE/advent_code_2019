package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	flag "github.com/spf13/pflag"
)

func main() {
	var dataFile string
	var total_fuel int64

	flag.StringVarP(&dataFile, "data file name", "f", "", "")
	flag.Parse()

	masses, err := buildList(dataFile)
	if err != nil {
		log.Fatal("Failed to get masses from input file!", err)
	}
	total_fuel = 0
	for _, mass := range masses {
		total_fuel += calculate_fuel(mass)
	}
	fmt.Println("fuel=", total_fuel)
}

func calculate_fuel(mass int64) int64 {
	var fuel int64
	fuel = (mass/3 - 2)
	if fuel < 0 {
		return 0
	} else {
		return fuel + calculate_fuel(fuel)
	}
}

func buildList(dataFile string) ([]int64, error) {
	var masses []int64

	f, err := os.Open(dataFile)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	buff := bufio.NewReader(f)
	for {
		line, err := buff.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSuffix(line, "\n")
		mass, err := strconv.ParseInt(line, 10, 64)
		if err != nil {
			log.Fatal(err)
			os.Exit(2)
		}
		masses = append(masses, mass)
	}
	return masses, err
}
