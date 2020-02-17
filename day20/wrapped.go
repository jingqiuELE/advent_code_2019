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

const (
	INNER_EDGE = iota
	OUTER_EDGE
	MIDDLE
)

type Vertex struct {
	Label      string
	Pos        Position
	Level      int
	Mark       int
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

func BFS(start *Vertex, vertexes []*Vertex) map[*Vertex]int {
	type PV struct {
		*Vertex
		pathLen int
	}
	var queue []PV

	nmap := make(map[*Vertex]int)
	visited := make(map[*Vertex]bool)

	queue = append(queue, PV{start, 0})

	foundExit := false
	for !foundExit && len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		visited[v.Vertex] = true

		if v.IsInnerPortal() {
			generateGraph(v)
		}

		for _, n := range v.Neighbours {
			if visited[n] == true {
				continue
			}
			nmap[n] = v.pathLen + 1
			if n.Label == "ZZ" && n.Level == 0 {
				fountExit = true
				break
			} else {
				queue = append(queue, PV{n, v.pathLen + 1})
			}
		}
	}
	return nmap
}

func generateGraph(portal *Vertex, vertexes []*Vertex) {
	var nvertexes []*Vertex

	if portal.IsInnerPortal() == false {
		log.Fatal("I can only generate graph from inner portal!")
		return
	}

	for _, n := range vertexes {
		if n.Label == portal.Label && n != portal {
			start = n
			break
		}
	}

	visited := make(map[*Vertex]bool)
	vmap := make(map[Position]*Vertex)

	if start == nil {
		log.Fatal("Couldn't find start point!")
		return
	}
	var queue []*Vertex

	queue = append(queue, start)

	for len(queue) != 0 {
		n := queue[0]
		queue = queue[1:]
		visited[n] = true
		t := n.Dup()
		t.Level++
		vmap[t.Position] = t
		nvertexes = append(nvertexes, t)

		for _, nn := range n.Neighbours {
			if visited[nn] == false {
				queue = append(queue, nn)
			}
		}
	}

	for _, v := range vertexes {
		dv = vmap[v.Position]
		for _, nv := range v.Neighbours {
			dn := vmap[nv.Position]
			if dn == nil {
				log.Fatal("Couldn't find duplicated vertex for postion %v", nv.Position)
				return
			}
			dv.AddNeighbours(dn)
		}
	}

	portal.AddNeighbours[vmap[start.Position]]

	for _, v := range nvertexes {
		if v.IsOuterPortal {
			for _, n := range vertexes {
				if n.Label == v.Label && n.Position != v.Position {
					v.AddNeighbours(n)
				}
			}
		}
	}
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
	var mark int

	vmap = make(map[Position]*Vertex)

	for y, l := range maze {
		for x, c := range l {
			if c == '.' {
				p := Position{
					Y: y,
					X: x,
				}
				label, direction := scanLabel(maze, p)
				if label != "" {
					fmt.Printf("got label %v at %v\n", label, p)
				}
				on := onEdge(maze, p, direction)
				if on {
					mark = OUTER_EDGE
				} else if label != "" {
					mark = INNER_EDGE
				} else {
					mark = MIDDLE
				}
				v := NewVertex(p, label, mark)
				vmap[p] = v
				vertexes = append(vertexes, v)
			}
		}
	}

	//Do a second pass to add neighbours
	for _, v := range vertexes {
		v.AddNeighbours(vmap)
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

func onEdge(maze [][]int, p Position) bool {

	if p.X == 2 || p.X == len(maze[0]-2) || p.Y == 2 || p.Y == len(maze)-2 {
		return true
	}

	return false
}

func NewVertex(p Position, c string, on int) *Vertex {
	return &Vertex{
		Pos:    p,
		Label:  c,
		OnEdge: on,
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

func (v *Vertex) IsInnerPortal() bool {
	if v.Label != "" && v.OnEdge == INNER {
		return true
	}
	return false
}

func (v *Vertex) IsOuterrPortal() bool {
	if v.Label != "" && v.OnEdge == OUTER && v.Level > 0 {
		return true
	}
	return false
}

func (v *Vertex) Dup() *Vertex {
	return &Vertex{
		Pos:    v.Pos,
		Label:  v.Label,
		Level:  v.Level,
		OnEdge: v.OnEdge,
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
