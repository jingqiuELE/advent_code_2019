package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	flag "github.com/spf13/pflag"
)

type Node struct {
	Name     string
	Children []*Node
	Parent   *Node
	Orbits   uint32
	Path     uint32
	Visited  bool
}

func main() {
	var dataFile string

	flag.StringVarP(&dataFile, "data file name", "f", "", "")
	flag.Parse()

	mnodes := buildTree(dataFile)

	walkTree(mnodes["COM"])
	num_orbits := Orbits(mnodes["COM"])
	fmt.Println("num_orbits:", num_orbits)

	BFS(mnodes["YOU"])
	for k, node := range mnodes {
		fmt.Printf("%s: path:%v\n", k, node.Path)
	}
	fmt.Println("paths:", mnodes["SAN"].Path-2)
}

func NewNode(name string) *Node {
	return &Node{
		Name:   name,
		Orbits: 0,
	}
}

func (node *Node) AddChild(child *Node) {
	node.Children = append(node.Children, child)
}

func (node *Node) SetParent(parent *Node) {
	node.Parent = parent
}

func buildTree(dataFile string) map[string]*Node {
	var mnodes map[string]*Node
	mnodes = make(map[string]*Node)

	f, err := os.Open(dataFile)
	if err != nil {
		return nil
	}

	defer f.Close()

	buff := bufio.NewReader(f)
	for {
		var a, b *Node

		line, err := buff.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSuffix(line, "\n")
		lists := strings.Split(line, ")")
		if mnodes[lists[0]] == nil {
			a = NewNode(lists[0])
			mnodes[a.Name] = a
		} else {
			a = mnodes[lists[0]]
		}

		if mnodes[lists[1]] == nil {
			b = NewNode(lists[1])
			mnodes[b.Name] = b
		} else {
			b = mnodes[lists[1]]
		}
		a.AddChild(b)
		b.SetParent(a)
	}
	return mnodes
}

func walkTree(root *Node) {
	for _, child := range root.Children {
		child.Orbits = root.Orbits + 1
		walkTree(child)
	}
}

func Orbits(root *Node) uint32 {
	total_orbits := root.Orbits
	for _, child := range root.Children {
		total_orbits += Orbits(child)
	}
	return total_orbits
}

func BFS(start *Node) {
	var queue []*Node

	fmt.Println("visiting:", start.Name)
	start.Visited = true
	for _, node := range start.Children {
		if node.Visited == false {
			node.Path = start.Path + 1
			queue = append(queue, node)
		}
	}

	if start.Parent != nil && start.Parent.Visited == false {
		start.Parent.Path = start.Path + 1
		queue = append(queue, start.Parent)
	}

	for _, next := range queue {
		BFS(next)
	}
}
