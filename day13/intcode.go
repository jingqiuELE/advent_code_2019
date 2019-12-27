package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

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
	Y int
	X int
}

const N = 50

func main() {
	var dataFile string

	flag.StringVarP(&dataFile, "data file name", "f", "", "")
	flag.Parse()

	program, err := buildList(dataFile)
	if err != nil {
		log.Fatal("Failed to get program from input file!", err)
	}

	input := make(chan int, 1)
	output := make(chan int, 3)
	games := make(chan map[Position]int)

	go runProgram(program, input, output)
	go playGames(input, games)

	panel := make(map[Position]int)
	var score int

	tick := time.Tick(10 * time.Millisecond)
	done := false
	max_y := 0
	max_x := 0
	for done == false {
		select {
		case x, ok := <-output:
			if ok == false {
				done = true
				break
			}

			y, ok := <-output
			if ok == false {
				done = true
				break
			}

			tile, ok := <-output
			if ok == false {
				done = true
				break
			}

			if x == -1 && y == 0 {
				score = tile
			} else {
				if y > max_y {
					max_y = y
				}

				if x > max_x {
					max_x = x
				}

				p := Position{
					Y: y,
					X: x,
				}
				panel[p] = tile
			}
		case <-tick:
			refreshScreen(panel, max_y, max_x)
			games <- panel
		}
	}
	fmt.Println("score=", score)
}

func refreshScreen(panel map[Position]int, max_y int, max_x int) {
	counter := 0
	var frame [][]int

	frame = make([][]int, max_y+1)
	for i, _ := range frame {
		frame[i] = make([]int, max_x+1)
	}

	for pos, tile := range panel {
		frame[pos.Y][pos.X] = tile
	}

	for _, line := range frame {
		for _, p := range line {
			switch p {
			case 0:
				fmt.Printf(" ")
			case 1:
				fmt.Printf("=")
			case 2:
				fmt.Printf("X")
				counter++
			case 3:
				fmt.Printf("_")
			case 4:
				fmt.Printf("O")
			}
		}
		fmt.Println("")
	}
}

func playGames(control chan int, pchannel chan map[Position]int) {
	var target Position
	var padel Position

	for {
		select {
		case panel := <-pchannel:
			for pos, tile := range panel {
				if tile == 4 {
					target = pos
				} else if tile == 3 {
					padel = pos
				}
			}
		}

		if padel.X < target.X {
			control <- 1
		} else if padel.X > target.X {
			control <- -1
		} else {
			control <- 0
		}
	}
}

func runControl(input chan int) {
	// disable input buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	b := make([]byte, 1)
	for {
		_, err := os.Stdin.Read(b)
		if err == io.EOF {
			break
		}
		fmt.Println("input=", b[0])
		switch b[0] {
		case 97:
			input <- -1
		case 111:
			input <- 1
		default:
			input <- 0
		}
	}
}

func runProgram(program []int, input chan int, output chan int) {
	p := make([]int, len(program)+1024*4)
	copy(p, program)

	p[0] = 2

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
			fmt.Println("machine got input", p[rpos])
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
