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

	nmap := BFS(start, vertexes)
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
	var nvertexes [][]*Vertex

	nmap := make(map[*Vertex]int)
	visited := make(map[*Vertex]bool)

	//Push the first layer of map into nvertexes
	nvertexes = append(nvertexes, vertexes)

	queue = append(queue, PV{start, 0})

	foundExit := false
	for !foundExit && len(queue) > 0 {
		var new []*Vertex

		v := queue[0]
		queue = queue[1:]
		visited[v.Vertex] = true

		if v.IsInnerPortal() {
			if len(nvertexes) <= v.Level+1 {
				fmt.Printf("generating graph at [%v] %v Level=%d, steps=%d\n", v.Pos, v.Label, v.Level, nmap[v.Vertex])
				new = generateGraph(nvertexes[v.Level], v.Level+1)
				nvertexes = append(nvertexes, new)
			} else {
				new = nvertexes[v.Level+1]
			}
			for _, n := range new {
				if n.Label == v.Label && n.Pos != v.Pos && !v.HasNeighbour(n) {
					fmt.Printf("connecting downward to %v from %v\n", n, v.Vertex)
					v.Neighbours = append(v.Neighbours, n)
					break
				}
			}
		} else if v.IsOuterPortal() {
			for _, n := range nvertexes[v.Level-1] {
				if n.Label == v.Label && n.Pos != v.Pos && !v.HasNeighbour(n) {
					fmt.Printf("Connecting upward to %v from %v\n", n, v.Vertex)
					v.Neighbours = append(v.Neighbours, n)
				}
			}
		}

		if v.Label != "" {
			fmt.Printf("Visiting [%v] %v at level %v, pathLen=%v\n",
				v.Pos, v.Label, v.Level, nmap[v.Vertex])
		}

		for _, n := range v.Neighbours {
			if visited[n] == true {
				continue
			}
			nmap[n] = v.pathLen + 1
			if n.Label == "ZZ" && n.Level == 0 {
				foundExit = true
				break
			} else {
				queue = append(queue, PV{n, v.pathLen + 1})
			}
		}
	}
	if !foundExit {
		fmt.Println("Exit not founded!")
	}
	return nmap
}

func generateGraph(vertexes []*Vertex, level int) []*Vertex {
	var new []*Vertex

	vmap := make(map[Position]*Vertex)

	for _, v := range vertexes {
		n := v.Dup()
		n.Level = level
		vmap[n.Pos] = n
		new = append(new, n)
	}

	for _, v := range new {
		v.AddNeighbours(vmap)
	}

	return new
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
				label := scanLabel(maze, p)
				on := onEdge(maze, p)
				if on {
					mark = OUTER_EDGE
				} else if label != "" {
					mark = INNER_EDGE
				} else {
					mark = MIDDLE
				}
				if label != "" {
					fmt.Printf("got label %v at %v, mark=%d\n", label, p, mark)
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

	if p.X == 2 || p.X == len(maze[0])-3 || p.Y == 2 || p.Y == len(maze)-3 {
		return true
	}

	return false
}

func NewVertex(p Position, c string, on int) *Vertex {
	return &Vertex{
		Pos:   p,
		Label: c,
		Mark:  on,
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

func (v *Vertex) HasNeighbour(n *Vertex) bool {
	for _, t := range v.Neighbours {
		if t == n {
			return true
		}
	}
	return false
}

func (v *Vertex) IsInnerPortal() bool {
	return v.Mark == INNER_EDGE
}

func (v *Vertex) IsOuterPortal() bool {
	if v.Label != "AA" && v.Label != "ZZ" && v.Label != "" && v.Mark == OUTER_EDGE && v.Level > 0 {
		return true
	}
	return false
}

func (v *Vertex) Dump() {
	if v.Label != "" {
		fmt.Printf("Pos[%v, %v], Label[%v] Level:%d\n", v.Pos.Y, v.Pos.X, v.Label, v.Level)
	}
}

func (v *Vertex) Dup() *Vertex {
	return &Vertex{
		Pos:   v.Pos,
		Label: v.Label,
		Level: v.Level,
		Mark:  v.Mark,
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
