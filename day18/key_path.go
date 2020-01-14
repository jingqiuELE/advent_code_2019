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

type Position struct {
	Y int
	X int
}

type Vertex struct {
	Value      int
	Pos        Position
	Neighbours []*Vertex
}

type Record struct {
	steps int
	path  []int
}

func main() {
	var dataFile string

	flag.StringVarP(&dataFile, "data file name", "f", "", "")
	flag.Parse()

	maze, err := buildMaze(dataFile)
	if err != nil {
		log.Fatal("Failed to get maze from input file!", err)
	}

	vertexes := buildGraph(maze)

	var start []*Vertex
	for _, v := range vertexes {
		if v.Value == '@' {
			start = append(start, v)
		}
	}

	dataBase := make(map[string]map[string]Record)
	shortestSteps, shortestPath := calculatePath(vertexes, start, "", dataBase)
	fmt.Printf("shortestSteps=%v\n", shortestSteps)
	fmt.Printf("shortestPath=")
	for _, c := range shortestPath {
		fmt.Printf("%c", c)
	}
	fmt.Println("")
}

const DISTANCE = ('A' - 'a')

func calculatePath(vertexes []*Vertex, start []*Vertex,
	holdKeys string, dataBase map[string]map[string]Record) (int, []int) {
	var minSteps int
	var minPath []int

	var next []map[*Vertex]int
	for i, s := range start {
		fmt.Printf("start[%v]: %c\n", i, s.Value)
		nmap := BFS(s)
		next = append(next, nmap)
	}

	fmt.Printf("next:")
	for _, nmap := range next {
		for v, pathLen := range nmap {
			fmt.Printf("%c:%v ", v.Value, pathLen)
		}
	}
	fmt.Println("")

	minSteps = 0
	for i, nmap := range next {
		for v, pathLen := range nmap {
			var path []int

			fmt.Printf("move to %c, pathLen=%v\n", v.Value, pathLen)

			newStart := make([]*Vertex, len(start))
			copy(newStart, start)
			for j, _ := range newStart {
				if j == i {
					newStart[j] = v
				}
			}
			state := composeState(newStart)
			fmt.Println("state=", state)
			dmap := dataBase[state]
			if dmap == nil {
				dmap = make(map[string]Record)
				dataBase[state] = dmap
			}
			keys := composeKeys(holdKeys, v.Value)
			fmt.Println("keys=", keys)
			history, ok := dmap[keys]
			if ok {
				fmt.Println("hit dataBase:", history)
				pathLen += history.steps
				path = history.path
			} else {
				var gate, key Vertex
				changed := false
				gv := lookupGate(vertexes, v.Value+DISTANCE)
				if gv != nil {
					gate = *gv
					gv.Value = int('.')
					changed = true
				}
				key = *v
				v.Value = int('.')

				remainingSteps, remainingPath := calculatePath(vertexes, newStart, keys, dataBase)
				rc := Record{
					steps: remainingSteps,
					path:  remainingPath,
				}
				dmap[keys] = rc
				pathLen += remainingSteps
				if changed {
					*gv = gate
				}
				*v = key
				fmt.Printf("recursion at %c: pathLen=%v\n", v.Value, pathLen)
				fmt.Printf("steps so far:")
				fmt.Printf("%c", v.Value)
				for _, c := range remainingPath {
					fmt.Printf("%c", c)
				}
				fmt.Println("")
				path = remainingPath
			}

			if minSteps == 0 || pathLen < minSteps {
				minSteps = pathLen
				minPath = make([]int, len(path)+1)
				minPath[0] = v.Value
				copy(minPath[1:], path)
			}
		}
	}
	return minSteps, minPath
}

func composeState(start []*Vertex) string {
	var s string
	for _, v := range start {
		sy := strconv.FormatInt(int64(v.Pos.Y), 16)
		sx := strconv.FormatInt(int64(v.Pos.X), 16)
		s = s + sy + sx
	}
	return s
}

func lookupGate(vertexes []*Vertex, value int) *Vertex {
	for _, v := range vertexes {
		if v.Value == value {
			return v
		}
	}
	return nil
}

func composeKeys(current string, next int) string {
	set := false
	for i, c := range current {
		if int(c) > next {
			current = current[:i] + string(next) + current[i:]
			set = true
			break
		}
	}
	if set == false {
		current = current + string(next)
	}
	return current
}

func BFS(start *Vertex) map[*Vertex]int {
	type PV struct {
		*Vertex
		pathLen int
	}
	var queue []PV

	nmap := make(map[*Vertex]int)
	visited := make(map[*Vertex]bool)

	queue = append(queue, PV{start, 0})

	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		visited[v.Vertex] = true
		for _, n := range v.Neighbours {
			if visited[n] == true {
				continue
			}

			if n.Value >= 'A' && n.Value <= 'Z' {
				//If we found a gate, don't traverse behind the gate
				continue
			} else if n.Value >= 'a' && n.Value <= 'z' {
				nmap[n] = v.pathLen + 1
			} else {
				queue = append(queue, PV{n, v.pathLen + 1})
			}
		}
	}
	return nmap
}

func buildMaze(dataFile string) ([][]int, error) {
	var maze [][]int

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
		var l []int
		for _, c := range line {
			l = append(l, int(c))
		}
		maze = append(maze, l)
	}
	return maze, err
}

func buildGraph(maze [][]int) []*Vertex {
	var vmap map[Position]*Vertex
	var vertexes []*Vertex

	vmap = make(map[Position]*Vertex)

	for y, l := range maze {
		for x, c := range l {
			if c == '#' {
				continue
			}
			p := Position{
				Y: y,
				X: x,
			}
			v := NewVertex(p, c)
			vmap[p] = v
			vertexes = append(vertexes, v)
		}
	}

	//Do a second pass to add neighbours
	for _, v := range vertexes {
		v.AddNeighbours(vmap)
	}

	return vertexes
}

func NewVertex(p Position, c int) *Vertex {
	return &Vertex{
		Pos:   p,
		Value: c,
	}
}

func (v *Vertex) AddNeighbours(vmap map[Position]*Vertex) {
	n := vmap[east(v.Pos)]
	if n != nil {
		v.Neighbours = append(v.Neighbours, n)
	}

	n = vmap[west(v.Pos)]
	if n != nil {
		v.Neighbours = append(v.Neighbours, n)
	}

	n = vmap[north(v.Pos)]
	if n != nil {
		v.Neighbours = append(v.Neighbours, n)
	}

	n = vmap[south(v.Pos)]
	if n != nil {
		v.Neighbours = append(v.Neighbours, n)
	}
}

func (v *Vertex) Dump() {
	fmt.Printf("%c \n", v.Value)
}

func east(p Position) Position {
	return Position{
		Y: p.Y,
		X: p.X - 1,
	}
}

func west(p Position) Position {
	return Position{
		Y: p.Y,
		X: p.X + 1,
	}
}

func north(p Position) Position {
	return Position{
		Y: p.Y - 1,
		X: p.X,
	}
}

func south(p Position) Position {
	return Position{
		Y: p.Y + 1,
		X: p.X,
	}
}
