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
	RELATIVE_BASE = 9
	HALT          = 99
)

const (
	POSITION  = 0
	IMMEDIATE = 1
	RELATIVE  = 2
)

type Position struct {
	X int
	Y int
}

var program []int

func main() {
	var dataFile string
	var err error

	flag.StringVarP(&dataFile, "data file name", "f", "", "")
	flag.Parse()

	program, err = buildList(dataFile)
	if err != nil {
		log.Fatal("Failed to get program from input file!", err)
	}

	result := make([][]int, 50)
	for i, _ := range result {
		result[i] = make([]int, 50)
	}

	for y := 0; y < 50; y++ {
		for x := 0; x < 50; x++ {
			pulled := scan(y, x)
			fmt.Printf("result=%v\n", pulled)
			result[y][x] = pulled
		}
	}
	count := 0
	for _, line := range result {
		for _, c := range line {
			if c == 1 {
				count++
				fmt.Printf("#")
			} else {
				fmt.Printf(".")
			}
		}
		fmt.Println("")
	}
	fmt.Println("count=", count)
	last_start := 0
	last_width := 0
	for y := 50; ; y++ {
		start := -1
		width := 0
		for x := last_start; ; x++ {
			pulled := scan(y, x)
			if pulled == 1 {
				start = x
				break
			}
		}

		for x := start + last_width; ; x++ {
			pulled := scan(y, x)
			if pulled == 0 {
				width = x - start
				fmt.Println("width=", width)
				ld := scan(y+99, x-100)
				rd := scan(y+99, x-1)
				fmt.Printf("ld=%v, rd=%v\n", ld, rd)
				if ld == 1 && rd == 1 {
					fmt.Printf("y=%v, x=%v\n", y, x-100)
					scan(946, 716)
					return
				}
				break
			}
		}
		last_start = start
		last_width = width - 5
		if last_width < 0 {
			last_width = 0
		}
	}
}

func scan(y int, x int) int {
	control := make(chan int, 2)
	output := make(chan int, 1)
	go runProgram(program, control, output)
	fmt.Printf("feeding (%v, %v), ", y, x)
	control <- x
	control <- y
	pulled := <-output
	fmt.Printf("%v\n", pulled)
	return pulled
}

func runProgram(program []int, input chan int, output chan int) {
	p := make([]int, len(program)+1024*8)
	copy(p, program)

	relative_base := 0
	i := 0
	for i < len(p) {
		opcode, pmode := parseCode(p[i])
		switch opcode {
		case ADD:
			param0 := loadParam(p[i+1], pmode[0], relative_base, p)
			param1 := loadParam(p[i+2], pmode[1], relative_base, p)
			rpos := loadPos(p[i+3], pmode[2], relative_base)
			p[rpos] = param0 + param1
			i += 4
		case MULTIPLY:
			param0 := loadParam(p[i+1], pmode[0], relative_base, p)
			param1 := loadParam(p[i+2], pmode[1], relative_base, p)
			rpos := loadPos(p[i+3], pmode[2], relative_base)
			p[rpos] = param0 * param1
			i += 4
		case STORE:
			rpos := loadPos(p[i+1], pmode[0], relative_base)
			p[rpos] = <-input
			i += 2
		case LOAD:
			param0 := loadParam(p[i+1], pmode[0], relative_base, p)
			output <- param0
			i += 2
		case JUMP_IF_TRUE:
			param0 := loadParam(p[i+1], pmode[0], relative_base, p)
			param1 := loadParam(p[i+2], pmode[1], relative_base, p)
			if param0 != int(0) {
				i = int(param1)
			} else {
				i += 3
			}
		case JUMP_IF_FALSE:
			param0 := loadParam(p[i+1], pmode[0], relative_base, p)
			param1 := loadParam(p[i+2], pmode[1], relative_base, p)
			if param0 == int(0) {
				i = int(param1)
			} else {
				i += 3
			}
		case LESS_THAN:
			param0 := loadParam(p[i+1], pmode[0], relative_base, p)
			param1 := loadParam(p[i+2], pmode[1], relative_base, p)
			rpos := loadPos(p[i+3], pmode[2], relative_base)
			if param0 < param1 {
				p[rpos] = 1
			} else {
				p[rpos] = 0
			}
			i += 4
		case EQUALS:
			param0 := loadParam(p[i+1], pmode[0], relative_base, p)
			param1 := loadParam(p[i+2], pmode[1], relative_base, p)
			rpos := loadPos(p[i+3], pmode[2], relative_base)
			if param0 == param1 {
				p[rpos] = 1
			} else {
				p[rpos] = 0
			}
			i += 4
		case RELATIVE_BASE:
			param0 := loadParam(p[i+1], pmode[0], relative_base, p)
			relative_base += param0
			i += 2
		case HALT:
			close(output)
			return
		default:
			log.Fatalf("unrecognized opcode=%v at %v\n", opcode, i)
		}
	}
}

func parseCode(code int) (opcode int, pmode []int) {
	var ps int

	opcode = code % 100

	pmode = make([]int, 3)

	ps = (code - opcode) / 100
	for i := 0; ps > 0; i++ {
		digit := ps % 10
		pmode[i] = digit
		ps = (ps - digit) / 10
	}
	return opcode, pmode
}

func loadParam(data int, mode int, base int, p []int) int {
	switch mode {
	case IMMEDIATE:
		return data
	case POSITION:
		return p[data]
	case RELATIVE:
		return p[base+data]
	}
	log.Fatal("Shouldn't get here")
	return -1
}

func loadPos(pos int, mode int, base int) int {
	switch mode {
	case POSITION:
		return pos
	case RELATIVE:
		return base + pos
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
