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
	NORTH = iota + 1
	SOUTH
	WEST
	EAST
)

const (
	NORMAL = iota
	WALL
	OXGEN_ROOM
	VISITED
)

type Vertex struct {
	Pos        Position
	Direction  int
	MinPath    int
	Visited    bool
	Property   int
	Neighbours []Position
}

type Position struct {
	X int
	Y int
}

func NewVertex(pos Position) *Vertex {
	v := Vertex{
		Pos: pos,
	}
	v.Neighbours = make([]Position, 4)
	for i, _ := range v.Neighbours {
		v.Neighbours[i] = pos
		switch i + 1 {
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

	vmap := make(map[Position]*Vertex)
	pos := Position{
		X: 0,
		Y: 0,
	}
	start := NewVertex(pos)
	vmap[pos] = start
	start.Visited = true
	BFS(start, vmap, control, output)

	var oxgen *Vertex
	for _, v := range vmap {
		if v.Property == OXGEN_ROOM {
			fmt.Println("MinPath to target:", v.MinPath)
			oxgen = v
		}
		v.Visited = false
	}

	printMaze(vmap)

	time := spreadCount(oxgen, vmap)
	fmt.Println("Time needed to fill all rooms with oxgen:", time)
}

func printMaze(vmap map[Position]*Vertex) {
	var minX, maxX, minY, maxY int

	for pos, _ := range vmap {
		if pos.X > maxX {
			maxX = pos.X
		}
		if pos.X < minX {
			minX = pos.X
		}
		if pos.Y > maxY {
			maxY = pos.Y
		}
		if pos.Y < minY {
			minY = pos.Y
		}
	}

	var frame [][]int
	frame = make([][]int, maxY-minY+1)
	for i, _ := range frame {
		frame[i] = make([]int, maxX-minX+1)
	}

	for _, v := range vmap {
		var x, y int
		x = v.Pos.X - minX
		y = v.Pos.Y - minY
		if v.Visited == false {
			frame[y][x] = v.Property
		} else {
			frame[y][x] = VISITED
		}
	}

	for _, line := range frame {
		for _, c := range line {
			switch c {
			case WALL:
				fmt.Printf("#")
			case NORMAL:
				fmt.Printf(" ")
			case OXGEN_ROOM:
				fmt.Printf("O")
			case VISITED:
				fmt.Printf("+")
			}
		}
		fmt.Println("")
	}
}

func spreadCount(start *Vertex, vmap map[Position]*Vertex) int {
	var queue []*Vertex
	var count int

	start.Visited = true
	queue = append(queue, start)

	count = -1
	for len(queue) > 0 {
		var nextQueue []*Vertex
		for _, v := range queue {
			for _, p := range v.Neighbours {
				n := vmap[p]
				if n == nil {
					log.Fatal("Couldn't find vertex for pos ", p)
				}
				if n.Visited == false && n.Property != WALL {
					n.Visited = true
					nextQueue = append(nextQueue, n)
				}
			}
		}
		queue = nextQueue
		count++
	}
	return count
}

func BFS(current *Vertex, vmap map[Position]*Vertex, control chan int, output chan int) {
	var queue []*Vertex
	var n *Vertex

	fmt.Println("visiting", current.Pos)
	for i, p := range current.Neighbours {
		if vmap[p] == nil {
			n = NewVertex(p)
			vmap[p] = n
		} else {
			n = vmap[p]
		}
		if n.Visited == false {
			n.Visited = true
			n.MinPath = current.MinPath + 1
			n.Direction = i + 1
			control <- i + 1
			feedback := <-output
			switch feedback {
			case 0:
				fmt.Printf("hit wall:%v,", n.Pos)
				n.Property = WALL
			case 1:
				fmt.Printf("add %v,", n.Pos)
				n.Property = NORMAL
				queue = append(queue, n)
				control <- oppositeDirection(i + 1)
				<-output
			case 2:
				fmt.Printf("add %v,", n.Pos)
				n.Property = OXGEN_ROOM
				queue = append(queue, n)
				control <- oppositeDirection(i + 1)
				<-output
			}
		}
	}

	fmt.Println("")
	for _, q := range queue {
		control <- q.Direction
		<-output
		BFS(q, vmap, control, output)
		control <- oppositeDirection(q.Direction)
		<-output
	}
}

func oppositeDirection(dir int) int {
	switch dir {
	case NORTH:
		return SOUTH
	case SOUTH:
		return NORTH
	case WEST:
		return EAST
	case EAST:
		return WEST
	default:
		log.Fatal("unable to handle", dir)
	}
	return -1
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
