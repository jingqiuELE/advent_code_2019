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

	flag.StringVarP(&dataFile, "data file name", "f", "", "")
	flag.Parse()

	program, err := buildList(dataFile)
	if err != nil {
		log.Fatal("Failed to get program from input file!", err)
	}

	var noun int64
	var verb int64
	for noun = 0; noun <= 99; noun++ {
		for verb = 0; verb <= 99; verb++ {
			result := runProgram(program, noun, verb)
			if result[0] == 19690720 {
				fmt.Printf("Found! noun=%v, verb=%v, result=%v\n", noun, verb, 100*noun+verb)
				return
			}
		}
	}
}

func runProgram(program []int64, noun int64, verb int64) []int64 {
	var p []int64
	p = make([]int64, len(program))
	copy(p, program)
	p[1] = noun
	p[2] = verb
	i := 0
	for p[i] != 99 && i < len(p) {
		if p[i] == 1 {
			pos1 := p[i+1]
			pos2 := p[i+2]
			rpos := p[i+3]
			p[rpos] = p[pos1] + p[pos2]
			i += 4
		} else if p[i] == 2 {
			pos1 := p[i+1]
			pos2 := p[i+2]
			rpos := p[i+3]
			p[rpos] = p[pos1] * p[pos2]
			i += 4
		} else {
			fmt.Println("Unrecognized opcode!")
			i++
		}
	}
	return p
}

func buildList(dataFile string) ([]int64, error) {
	var program []int64

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
		lists := strings.Split(line, ",")
		for _, l := range lists {
			num, err := strconv.ParseInt(l, 10, 64)
			if err != nil {
				log.Fatal(err)
				os.Exit(2)
			}
			program = append(program, num)
		}
	}
	return program, err
}
