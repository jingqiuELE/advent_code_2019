package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	flag "github.com/spf13/pflag"
)

func main() {
	var dataFile string

	flag.StringVarP(&dataFile, "data file name", "f", "", "")
	flag.Parse()

	vertexes, uplinks, downlinks, err := buildGraph(dataFile)
	if err != nil {
		log.Fatal("Failed to build graph from input file!", err)
	}

	calcWeight(vertexes, uplinks)
	ore_needed := calcOre(vertexes, downlinks, 1)
	fmt.Println("result=", ore_needed)

	var target uint64
	var fuel uint64

	target = 1000000000000
	min_fuel := target / ore_needed
	max_fuel := min_fuel * 3
	for min_fuel < max_fuel-1 {
		fuel = (max_fuel + min_fuel) / 2
		ore_needed = calcOre(vertexes, downlinks, fuel)
		if ore_needed < target {
			min_fuel = fuel
		} else if ore_needed > target {
			max_fuel = fuel
		} else {
			break
		}
		fmt.Printf("min=%v, max=%v\n", min_fuel, max_fuel)
	}

	if min_fuel >= max_fuel-1 {
		fuel = min_fuel
	}
	fmt.Printf("fuel=%v\n", fuel)
	ore_needed = calcOre(vertexes, downlinks, fuel)
	fmt.Println("ore needed = ", ore_needed)
}

func buildGraph(dataFile string) ([]*Vertex, []*Edge, []*Edge, error) {
	var uplinks []*Edge
	var downlinks []*Edge
	var vertexes []*Vertex
	var vertex_map map[string]*Vertex

	f, err := os.Open(dataFile)
	if err != nil {
		return nil, nil, nil, err
	}

	defer f.Close()

	vertex_map = make(map[string]*Vertex)

	buff := bufio.NewReader(f)
	for {
		line, err := buff.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, nil, nil, err
		}
		line = strings.TrimSuffix(line, "\n")
		index := strings.Index(line, "=>")

		/* Process target */
		num, name := retrieveInfo(line[index+2:])
		target := Slab{
			Name: name,
			Num:  num,
		}
		tvertex := createVertex(vertex_map, name)

		/* Process list of sources */
		var ingredients []*Slab
		ilist := strings.Split(line[:index], ",")
		for _, l := range ilist {
			num, name := retrieveInfo(l)
			s := Slab{
				Name: name,
				Num:  num,
			}
			ingredients = append(ingredients, &s)
			svertex := createVertex(vertex_map, name)
			uplink := NewEdge(svertex, tvertex)
			uplinks = append(uplinks, uplink)

			downlink := NewEdge(tvertex, svertex)
			downlinks = append(downlinks, downlink)

			svertex.AddUplink(uplink)
			tvertex.UpDegrees++
			tvertex.AddDownlink(downlink)
		}

		/* Create recipe */
		recipe := Recipe{
			Target:      &target,
			Ingredients: ingredients,
		}
		tvertex.Recipe = &recipe
	}

	for _, v := range vertex_map {
		vertexes = append(vertexes, v)
	}
	return vertexes, uplinks, downlinks, nil
}

func createVertex(vertex_map map[string]*Vertex, name string) *Vertex {
	if vertex_map[name] == nil {
		vertex := NewVertex(name)
		vertex_map[name] = vertex
	}
	return vertex_map[name]
}

func retrieveInfo(l string) (uint64, string) {
	l = strings.TrimSpace(l)
	slab := strings.Split(l, " ")
	num, err := strconv.ParseUint(slab[0], 10, 64)
	if err != nil {
		log.Fatal(err)
		return 0, slab[1]
	}
	return num, slab[1]
}

func calcWeight(vertexes []*Vertex, uplinks []*Edge) {
	var root *Vertex
	var queue []*Vertex

	for _, v := range vertexes {
		if v.Name == "ORE" {
			root = v
			break
		}
	}

	if root.UpDegrees != 0 {
		log.Fatal("We assume root has 0 UpDegrees, not satisfied!")
		return
	}

	/* calculate longest path from root to every other node in the uplink DAG. */
	queue = append(queue, root)
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		for _, e := range v.Uplinks {
			if v.Weight+1 > e.Tail.Weight {
				e.Tail.Weight = v.Weight + 1
			}
			e.Tail.UpDegrees--
			if e.Tail.UpDegrees == 0 {
				queue = append(queue, e.Tail)
			}
		}
	}
}

func calcOre(vertexes []*Vertex, downlinks []*Edge, needed uint64) uint64 {
	var remaining VertexHeap
	var root *Vertex

	/* Find root with name "FUEL" */
	for _, v := range vertexes {
		if v.Name == "FUEL" {
			root = v
		}
		v.Num = 0
	}

	root.Num = needed
	heap.Push(&remaining, root)
	for len(remaining) > 0 {
		v := heap.Pop(&remaining).(*Vertex)
		if v.Recipe == nil {
			return v.Num
		}
		num_units := roundup(v.Num, v.Recipe.Target.Num)
		for _, edge := range v.Downlinks {
			if edge.Tail.InHeap == false {
				heap.Push(&remaining, edge.Tail)
			}
			i := findIngredientIndex(edge.Tail.Name, v.Recipe.Ingredients)
			edge.Tail.Num += num_units * v.Recipe.Ingredients[i].Num
		}
	}
	return 0
}

func roundup(needed uint64, unit uint64) uint64 {
	tmp := float64(needed) / float64(unit)
	num_unit := uint64(math.Ceil(tmp))
	return num_unit
}

func findIngredientIndex(name string, ingredients []*Slab) int {
	for index, ingredient := range ingredients {
		if ingredient.Name == name {
			return index
		}
	}
	log.Fatalf("Failed to find ingredient %s\n", name)
	return -1
}
