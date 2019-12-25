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

const (
	UP = iota
	LEFT
	RIGHT
	DOWN
)

const (
	TURN_LEFT  = 0
	TURN_RIGHT = 1
)

type Position struct {
	X int
	Y int
}

func main() {
	var dataFile string
	var painted map[Position]bool

	flag.StringVarP(&dataFile, "data file name", "f", "", "")
	flag.Parse()

	program, err := buildList(dataFile)
	if err != nil {
		log.Fatal("Failed to get program from input file!", err)
	}

	panels := make([][]int, 201)
	for i, _ := range panels {
		panels[i] = make([]int, 201)
	}
	painted = make(map[Position]bool)

	pos := Position{
		X: 100,
		Y: 100,
	}
	facing := UP
	panels[pos.X][pos.Y] = 1

	input := make(chan int, 1)
	output := make(chan int, 2)
	go runProgram(program, input, output)
	input <- 1

	done := false
	for done == false {
		select {
		case color, ok := <-output:
			if ok == false {
				done = true
				break
			}

			direction, ok := <-output
			if ok == false {
				done = true
				break
			}
			fmt.Printf("Color=%v, Dir=%v\n", color, direction)
			panels[pos.X][pos.Y] = color
			painted[pos] = true
			pos, facing = calculateNewPos(pos, facing, direction)
			fmt.Println("NewPos:", pos)
			input <- panels[pos.X][pos.Y]
		}
	}

	count := 0
	for k, _ := range painted {
		fmt.Println(k)
		count++
	}
	fmt.Println("count=", count)

	for _, v := range panels {
		for _, k := range v {
			if k == 1 {
				fmt.Printf("X")
			} else {
				fmt.Printf(" ")
			}
		}
		fmt.Println("")
	}
}

func calculateNewPos(pos Position, facing int, dir int) (Position, int) {
	npos := pos
	nfacing := facing

	switch facing {
	case UP:
		switch dir {
		case TURN_LEFT:
			nfacing = LEFT
			npos.X--
		case TURN_RIGHT:
			nfacing = RIGHT
			npos.X++
		}
	case LEFT:
		switch dir {
		case TURN_LEFT:
			nfacing = DOWN
			npos.Y++
		case TURN_RIGHT:
			nfacing = UP
			npos.Y--
		}
	case DOWN:
		switch dir {
		case TURN_LEFT:
			nfacing = RIGHT
			npos.X++
		case TURN_RIGHT:
			nfacing = LEFT
			npos.X--
		}
	case RIGHT:
		switch dir {
		case TURN_LEFT:
			nfacing = UP
			npos.Y--
		case TURN_RIGHT:
			nfacing = DOWN
			npos.Y++
		}
	}

	return npos, nfacing
}

func runProgram(program []int, input chan int, output chan int) {
	p := make([]int, len(program)+1024*4)
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
			log.Fatal("unrecognized opcode=", opcode)
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
