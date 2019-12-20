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

	program, err := buildList(dataFile)
	if err != nil {
		log.Fatal("Failed to get program from input file!", err)
	}

	runProgram(program, 5)
}

func runProgram(p []int64, input int64) []int64 {
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
			p[rpos] = input
			i += 2
		case LOAD:
			param0 := loadParam(p[i+1], pmode[0], p)
			fmt.Printf("output[%d]: %v\n", i, param0)
			i += 2
		case JUMP_IF_TRUE:
			param0 := loadParam(p[i+1], pmode[0], p)
			param1 := loadParam(p[i+2], pmode[1], p)
			if param0 != int64(0) {
				i = int(param1)
			} else {
				i += 3
			}
		case JUMP_IF_FALSE:
			param0 := loadParam(p[i+1], pmode[0], p)
			param1 := loadParam(p[i+2], pmode[1], p)
			if param0 == int64(0) {
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
			fmt.Println("HALT at ", i)
			return p
		default:
			log.Fatal("unrecognized opcode=", opcode)
			return p
		}
	}
	fmt.Println("Program finished without HALT! i=", i)
	return p
}

func parseCode(code int64) (opcode int64, pmode []int64) {
	var ps int64

	opcode = code % 100

	switch opcode {
	case ADD:
		pmode = make([]int64, 3)
	case MULTIPLY:
		pmode = make([]int64, 3)
	case LOAD:
		pmode = make([]int64, 1)
	case STORE:
		pmode = make([]int64, 1)
	case JUMP_IF_TRUE:
		pmode = make([]int64, 2)
	case JUMP_IF_FALSE:
		pmode = make([]int64, 2)
	case LESS_THAN:
		pmode = make([]int64, 3)
	case EQUALS:
		pmode = make([]int64, 3)
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

func loadParam(data int64, mode int64, p []int64) int64 {
	switch mode {
	case IMMEDIATE:
		return data
	case POSITION:
		return p[data]
	}
	log.Fatal("Shouldn't get here")
	return -1
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
