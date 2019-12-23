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

const (
	ADD           = 1
	MULTIPLY      = 2
	STORE         = 3
	LOAD          = 4
	JUMP_IF_TRUE  = 5
	JUMP_IF_FALSE = 6
	LESS_THAN     = 7
	EQUALS        = 8
	HALT          = 99
)

const (
	POSITION  = 0
	IMMEDIATE = 1
)

func main() {
	var dataFile string

	flag.StringVarP(&dataFile, "data file name", "f", "", "")
	flag.Parse()

	candidates := []int{5, 6, 7, 8, 9}
	phases := generatePhases(candidates)

	program, err := buildList(dataFile)
	if err != nil {
		log.Fatal("Failed to get program from input file!", err)
	}

	fmt.Println("num of phases=", len(phases))
	fmt.Println("phases=", phases)

	var max int
	max = -1
	for _, phase := range phases {
		input := make(chan int, 1)
		output0 := amplifier(program, phase[0], input)
		output1 := amplifier(program, phase[1], output0)
		output2 := amplifier(program, phase[2], output1)
		output3 := amplifier(program, phase[3], output2)
		output4 := amplifier(program, phase[4], output3)

		input <- 0
		var result int
		for result = range output4 {
			fmt.Println("result=", result)
			input <- result
		}

		fmt.Println("tp0: result=", result)
		if result > max {
			max = result
		}
	}
	fmt.Println("max=", max)
}

func generatePhases(candidates []int) [][]int {
	var result [][]int

	if len(candidates) == 1 {
		result = append(result, []int{candidates[0]})
	} else {
		for i, k := range candidates {
			nextCandidates := make([]int, len(candidates)-1)
			copy(nextCandidates, candidates[0:i])
			copy(nextCandidates[i:], candidates[i+1:])
			next := generatePhases(nextCandidates)
			for _, n := range next {
				n = append(n, k)
				result = append(result, n)
			}
		}
	}
	return result
}

func amplifier(program []int, phase int, input chan int) chan int {
	var p []int
	p = make([]int, len(program))
	copy(p, program)
	output := make(chan int)
	go runProgram(p, input, output)
	input <- phase
	return output
}

func runProgram(p []int, input chan int, output chan int) {
	i := 0
	for i < len(p) {
		opcode, pmode := parseCode(p[i])
		switch opcode {
		case ADD:
			param0 := loadParam(p[i+1], pmode[0], p)
			param1 := loadParam(p[i+2], pmode[1], p)
			rpos := p[i+3]
			p[rpos] = param0 + param1
			i += 4
		case MULTIPLY:
			param0 := loadParam(p[i+1], pmode[0], p)
			param1 := loadParam(p[i+2], pmode[1], p)
			rpos := p[i+3]
			p[rpos] = param0 * param1
			i += 4
		case STORE:
			rpos := p[i+1]
			p[rpos] = <-input
			i += 2
		case LOAD:
			param0 := loadParam(p[i+1], pmode[0], p)
			fmt.Printf("output: %v\n", param0)
			output <- param0
			i += 2
		case JUMP_IF_TRUE:
			param0 := loadParam(p[i+1], pmode[0], p)
			param1 := loadParam(p[i+2], pmode[1], p)
			if param0 != int(0) {
				i = int(param1)
			} else {
				i += 3
			}
		case JUMP_IF_FALSE:
			param0 := loadParam(p[i+1], pmode[0], p)
			param1 := loadParam(p[i+2], pmode[1], p)
			if param0 == int(0) {
				i = int(param1)
			} else {
				i += 3
			}
		case LESS_THAN:
			param0 := loadParam(p[i+1], pmode[0], p)
			param1 := loadParam(p[i+2], pmode[1], p)
			rpos := p[i+3]
			if param0 < param1 {
				p[rpos] = 1
			} else {
				p[rpos] = 0
			}
			i += 4
		case EQUALS:
			param0 := loadParam(p[i+1], pmode[0], p)
			param1 := loadParam(p[i+2], pmode[1], p)
			rpos := p[i+3]
			if param0 == param1 {
				p[rpos] = 1
			} else {
				p[rpos] = 0
			}
			i += 4
		case HALT:
			fmt.Println("HALT")
			close(output)
			return
		default:
			log.Fatal("unrecognized opcode=", opcode)
		}
	}
}

func parseCode(code int) (opcode int, pmode []int) {
	var ps int

	opcode = code % 100

	switch opcode {
	case ADD:
		pmode = make([]int, 3)
	case MULTIPLY:
		pmode = make([]int, 3)
	case LOAD:
		pmode = make([]int, 1)
	case STORE:
		pmode = make([]int, 1)
	case JUMP_IF_TRUE:
		pmode = make([]int, 2)
	case JUMP_IF_FALSE:
		pmode = make([]int, 2)
	case LESS_THAN:
		pmode = make([]int, 3)
	case EQUALS:
		pmode = make([]int, 3)
	case HALT:
	default:
		log.Fatal("unrecognized opcode:", opcode)
	}

	ps = (code - opcode) / 100
	for i := 0; ps > 0; i++ {
		digit := ps % 10
		pmode[i] = digit
		ps = (ps - digit) / 10
	}
	return opcode, pmode
}

func loadParam(data int, mode int, p []int) int {
	switch mode {
	case IMMEDIATE:
		return data
	case POSITION:
		return p[data]
	}
	log.Fatal("Shouldn't get here")
	return -1
}

func buildList(dataFile string) ([]int, error) {
	var program []int

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
			program = append(program, int(num))
		}
	}
	return program, err
}
