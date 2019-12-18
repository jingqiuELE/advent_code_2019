package main

import "fmt"

type Edge struct {
	Head *Vertex
	Tail *Vertex
}

func NewEdge(head *Vertex, tail *Vertex) *Edge {
	e := Edge{
		Head: head,
		Tail: tail,
	}
	return &e
}

func (e *Edge) Swap() {
	temp := e.Head
	e.Head = e.Tail
	e.Tail = temp
}

func (e *Edge) Dump() {
	fmt.Printf("%s-->%s\n", e.Head.Name, e.Tail.Name)
}
