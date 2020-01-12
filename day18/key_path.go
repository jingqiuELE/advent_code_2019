package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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

func main() {
	var dataFile string

	flag.StringVarP(&dataFile, "data file name", "f", "", "")
	flag.Parse()

	maze, err := buildMaze(dataFile)
	if err != nil {
		log.Fatal("Failed to get maze from input file!", err)
	}

	vertexes, kmap := buildGraph(maze)
	for key, v := range kmap {
		fmt.Printf("key:%c, v:%v\n", key, v)
	}
	dataBase := make(map[int]map[string]int)
	shortest_path := calculatePath(vertexes, kmap[int('@')], kmap, "", dataBase)
	fmt.Println("shortest path=", shortest_path)
}

const DISTANCE = ('A' - 'a')

func calculatePath(vertexes []*Vertex, start *Vertex, kmap map[int]*Vertex,
	holdKeys string, dataBase map[int]map[string]int) int {
	var minPath int

	next := BFS(start)
	fmt.Printf("next:")
	for v, pathLen := range next {
		fmt.Printf("%c:%v ", v.Value, pathLen)
	}
	fmt.Println("")

	minPath = 0
	for v, pathLen := range next {
		fmt.Printf("move to %c, pathLen=%v\n", v.Value, pathLen)
		dmap := dataBase[v.Value]
		if dmap == nil {
			dmap = make(map[string]int)
			dataBase[v.Value] = dmap
		}
		keys := composeKeys(holdKeys, v.Value)
		fmt.Println("keys=", keys)
		history, ok := dmap[keys]
		if ok {
			fmt.Println("hit dataBase:", history)
			pathLen += history
		} else {
			var gate, key Vertex
			changed := false
			gv := kmap[v.Value+DISTANCE]
			if gv != nil {
				gate = *gv
				gv.Value = int('.')
				changed = true
			}
			key = *v
			v.Value = int('.')
			remainingSteps := calculatePath(vertexes, v, kmap, keys, dataBase)
			dmap[keys] = remainingSteps
			pathLen += remainingSteps
			if changed {
				//Restore the kmap
				*gv = gate
			}
			*v = key
			fmt.Printf("after recurtion at %c: pathLen=%v\n", v.Value, pathLen)
		}

		if minPath == 0 || pathLen < minPath {
			minPath = pathLen
		}
	}
	return minPath
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

	visited := make(map[*Vertex]bool)
	next := make(map[*Vertex]int)

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
				next[n] = v.pathLen + 1
			} else {
				queue = append(queue, PV{n, v.pathLen + 1})
			}
		}
	}
	return next
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

func buildGraph(maze [][]int) ([]*Vertex, map[int]*Vertex) {
	var kmap map[int]*Vertex
	var vmap map[Position]*Vertex
	var vertexes []*Vertex

	vmap = make(map[Position]*Vertex)
	kmap = make(map[int]*Vertex)

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
			if c != int('.') {
				kmap[c] = v
			}
			vertexes = append(vertexes, v)
		}
	}

	//Do a second pass to add neighbours
	for _, v := range vertexes {
		v.AddNeighbours(vmap)
	}

	return vertexes, kmap
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
