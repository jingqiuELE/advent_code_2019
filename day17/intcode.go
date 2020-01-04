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

const (
	SCAFF = iota
	CROSS
)

const (
	FORWARD = iota
	TURN_RIGHT
	TURN_LEFT
)

const (
	NORTH = iota
	SOUTH
	EAST
	WEST
)

type Vertex struct {
	Pos        Position
	Visited    bool
	Property   int
	Neighbours []Position
}

func NewVertex(pos Position) *Vertex {
	v := Vertex{
		Pos: pos,
	}
	v.Neighbours = make([]Position, 4)
	for i, _ := range v.Neighbours {
		v.Neighbours[i] = pos
		switch i {
		case NORTH:
			v.Neighbours[i].Y--
		case SOUTH:
			v.Neighbours[i].Y++
		case WEST:
			v.Neighbours[i].X--
		case EAST:
			v.Neighbours[i].X++
		}
	}
	return &v
}

func main() {
	var dataFile string

	flag.StringVarP(&dataFile, "data file name", "f", "", "")
	flag.Parse()

	program, err := buildList(dataFile)
	if err != nil {
		log.Fatal("Failed to get program from input file!", err)
	}

	control := make(chan int, 1)
	output := make(chan int, 1)

	go runProgram(program, control, output)

	var frame [][]int
	var line []int
	done := false
	for done == false {
		select {
		case ascii, ok := <-output:
			if ok == false {
				done = true
				break
			}
			switch ascii {
			case 10:
				if len(line) > 0 {
					frame = append(frame, line)
				}
				line = []int{}
			default:
				line = append(line, ascii)
			}
		}
	}

	sum := 0
	for y, line := range frame {
		if y == len(frame)-1 || y == 0 {
			continue
		}
		for x, _ := range line {
			if x == len(line)-1 || x == 0 {
				continue
			}
			if frame[y][x] == '#' {
				if frame[y+1][x] == '#' && frame[y-1][x] == '#' &&
					frame[y][x+1] == '#' && frame[y][x-1] == '#' {
					frame[y][x] = 'O'
					alignment := y * x
					sum += alignment
				}
			}
		}
	}
	fmt.Println("sum of alignment=", sum)
	dumpVideo(frame)

	var start *Vertex
	var dir int
	vmap := make(map[Position]*Vertex)
	for y, line := range frame {
		for x, _ := range line {
			var v *Vertex
			if frame[y][x] != '.' {
				p := Position{
					Y: y,
					X: x,
				}
				v = NewVertex(p)
				vmap[p] = v
			}

			switch frame[y][x] {
			case '^':
				start = v
				dir = NORTH
			case '#':
				v.Property = SCAFF
			case 'O':
				v.Property = CROSS
			}
		}
	}

	steps := DFS(vmap, start, dir)
	fmt.Println(steps)
	str := stringSteps(steps)
	fmt.Println("result:", str)
	mainRoutine := "A,B,A,B,A,C,B,C,A,C"
	funcA := "R,4,L,10,L,10"
	funcB := "L,8,R,12,R,10,R,4"
	funcC := "L,8,L,8,R,10,R,4"
	wantFeed := "n"

	control_1 := make(chan int, 1)
	output_1 := make(chan int, 1)
	program[0] = 2

	go runProgram(program, control_1, output_1)
	var command string
	for {
		c := <-output_1
		if c < 255 {
			fmt.Printf("%c", c)
		} else {
			fmt.Printf("result: %v\n", c)
		}
		if c == '\n' {
			//We got a command
			switch command {
			case "Main:":
				sendControl(control_1, mainRoutine)
			case "Function A:":
				sendControl(control_1, funcA)
			case "Function B:":
				sendControl(control_1, funcB)
			case "Function C:":
				sendControl(control_1, funcC)
			case "Continuous video feed?":
				sendControl(control_1, wantFeed)
			}
			command = ""
		} else {
			command = command + string(c)
		}
	}
}

func sendControl(control chan int, data string) {
	fmt.Println("sending:", data)
	for _, c := range data {
		control <- int(c)
	}
	control <- '\n'
}

func stringSteps(steps []int) string {
	var result, pending_action string
	var counter int

	counter = 0

	for _, v := range steps {
		switch v {
		case TURN_RIGHT:
			result = result + pending_action + strconv.Itoa(counter)
			result += ", "
			pending_action = "R,"
			counter = 0
		case TURN_LEFT:
			result = result + pending_action + strconv.Itoa(counter)
			result += ", "
			pending_action = "L,"
			counter = 0
		case FORWARD:
			counter++
		}
	}
	result = result + pending_action + strconv.Itoa(counter)
	return result
}

func DFS(vmap map[Position]*Vertex, start *Vertex, dir int) []int {
	var steps []int

	start.Visited = true
	p := start.Neighbours[dir]
	if vmap[p] != nil {
		/* Favour straght forward direction first. */
		n := vmap[p]
		if n.Visited == false || n.Property == CROSS {
			steps = append(steps, FORWARD)
			remaining := DFS(vmap, vmap[p], dir)
			steps = append(steps, remaining...)
		}
	} else {
		for i, p := range start.Neighbours {
			if vmap[p] != nil {
				if vmap[p].Visited == false || vmap[p].Property == CROSS {
					action := getAction(dir, i)
					steps = append(steps, action, FORWARD)
					remaining := DFS(vmap, vmap[p], i)
					steps = append(steps, remaining...)
				}
			}
		}
	}
	return steps
}

func getAction(dir int, ndir int) int {
	fmt.Printf("dir=%v, ndir=%v\n", dir, ndir)
	switch dir {
	case NORTH:
		switch ndir {
		case WEST:
			return TURN_LEFT
		case EAST:
			return TURN_RIGHT
		default:
			log.Fatal("illegal case!")
			return FORWARD
		}
	case SOUTH:
		switch ndir {
		case WEST:
			return TURN_RIGHT
		case EAST:
			return TURN_LEFT
		default:
			log.Fatal("illegal case!")
		}
	case WEST:
		switch ndir {
		case NORTH:
			return TURN_RIGHT
		case SOUTH:
			return TURN_LEFT
		default:
			log.Fatal("illegal case!")
		}
	case EAST:
		switch ndir {
		case NORTH:
			return TURN_LEFT
		case SOUTH:
			return TURN_RIGHT
		default:
			log.Fatal("illegal case!")
		}
	}
	log.Fatal("Should never get here!")
	return -1
}

func dumpVideo(frame [][]int) {
	for _, line := range frame {
		for _, c := range line {
			fmt.Printf("%c", rune(c))
		}
		fmt.Println("")
	}
}

func runProgram(program []int, input chan int, output chan int) {
	p := make([]int, len(program)+1024*8)
	copy(p, program)
	fmt.Println("p[0]=", p[0])

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
			fmt.Println("HALT")
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
