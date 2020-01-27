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
	Label      string
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

	vertexes := buildGraph(maze)

	var start, end *Vertex
	for _, v := range vertexes {
		if v.Label == "AA" {
			start = v
			fmt.Println("Found start vertex")
		}

		if v.Label == "ZZ" {
			end = v
			fmt.Println("Found end vertex")
		}
	}

	nmap := BFS(start)
	fmt.Println(nmap[end])
}

func calculatePath(vertexes []*Vertex, start *Vertex) (int, []int) {
	var minSteps int
	var minPath []int

	return minSteps, minPath
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
			nmap[n] = v.pathLen + 1
			queue = append(queue, PV{n, v.pathLen + 1})
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
			if c == '.' {
				p := Position{
					Y: y,
					X: x,
				}
				label := scanLabel(maze, p)
				if label != "" {
					fmt.Printf("got label %v at %v\n", label, p)
				}
				v := NewVertex(p, label)
				vmap[p] = v
				vertexes = append(vertexes, v)
			}
		}
	}

	//Do a second pass to add neighbours
	for _, v := range vertexes {
		v.AddNeighbours(vmap)
	}

	//Do a third pass to connect the portals
	for _, v := range vertexes {
		if v.Label != "" {
			for _, t := range vertexes {
				if t.Label == v.Label && t != v {
					fmt.Printf("Appending to %v with %v\n", v, t)
					v.Neighbours = append(v.Neighbours, t)
				}
			}
		}
	}
	return vertexes
}

func scanLabel(maze [][]int, p Position) string {
	var result string

	n := east(p)
	c := maze[n.Y][n.X]
	if c >= 'A' && c <= 'Z' {
		n0 := east(n)
		c0 := maze[n0.Y][n0.X]
		result = string(c)
		result = result + string(c0)
		return result
	}

	n = south(p)
	c = maze[n.Y][n.X]
	if c >= 'A' && c <= 'Z' {
		n0 := south(n)
		c0 := maze[n0.Y][n0.X]
		result = string(c)
		result = result + string(c0)
		return result
	}

	n = west(p)
	c = maze[n.Y][n.X]
	if c >= 'A' && c <= 'Z' {
		n0 := west(n)
		c0 := maze[n0.Y][n0.X]
		result = string(c0)
		result = result + string(c)
		return result
	}

	n = north(p)
	c = maze[n.Y][n.X]
	if c >= 'A' && c <= 'Z' {
		n0 := north(n)
		c0 := maze[n0.Y][n0.X]
		result = string(c0)
		result = result + string(c)
		return result
	}

	return ""
}

func NewVertex(p Position, c string) *Vertex {
	return &Vertex{
		Pos:   p,
		Label: c,
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

type dirFunc func(Position) Position

func east(p Position) Position {
	return Position{
		Y: p.Y,
		X: p.X + 1,
	}
}

func west(p Position) Position {
	return Position{
		Y: p.Y,
		X: p.X - 1,
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
