package main

import "fmt"

type Slab struct {
	Name string
	Num  uint64
}

type Recipe struct {
	Target      *Slab
	Ingredients []*Slab
}

type Vertex struct {
	Name      string
	Num       uint64
	UpDegrees uint64
	Weight    uint64
	Uplinks   []*Edge
	Downlinks []*Edge
	InHeap    bool
	Recipe    *Recipe
}

func NewVertex(name string) *Vertex {
	return &Vertex{
		Name: name,
	}
}

func (v *Vertex) AddUplink(edge *Edge) {
	v.Uplinks = append(v.Uplinks, edge)
}

func (v *Vertex) AddDownlink(edge *Edge) {
	v.Downlinks = append(v.Downlinks, edge)
}

func (v *Vertex) Dump() {
	fmt.Printf("dumpVertex: %v\n", v)
	for _, e := range v.Downlinks {
		fmt.Printf("%s-->%s\n", e.Head.Name, e.Tail.Name)
	}
}
