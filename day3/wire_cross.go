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

type Wire struct {
	Lines []Line
}

const MaxUint64 = ^uint64(0)
const MaxInt64 = int64(MaxUint64 >> 1)

func main() {
	var dataFile string

	flag.StringVarP(&dataFile, "data file name", "f", "", "")
	flag.Parse()

	wires, err := buildWires(dataFile)
	if err != nil {
		log.Fatal("Failed to build wires from input file!", err)
	}
	if len(wires) != 2 {
		log.Fatal("Input should have 2 wires, got ", len(wires))
	}

	cross_points := calculateCross(wires[0], wires[1])
	var min_distance int64
	min_distance = MaxInt64
	for _, p := range cross_points {
		d := mahattan(p, Point{0, 0})
		if d > 0 && d < min_distance {
			min_distance = d
		}
	}
	fmt.Println("min_distance = ", min_distance)

	var min_steps int64
	min_steps = MaxInt64

	for _, p := range cross_points {
		fmt.Println("p=", p)
		steps_a := calculateSteps(wires[0], p)
		fmt.Println("steps_a:", steps_a)
		steps_b := calculateSteps(wires[1], p)
		fmt.Println("steps_b:", steps_b)
		sum := steps_a + steps_b
		if sum > 0 && sum < min_steps {
			min_steps = sum
		}
	}
	fmt.Println("min_steps:", min_steps)
}

func buildWires(dataFile string) ([]Wire, error) {
	var wires []Wire
	var sp Point
	var steps int64

	f, err := os.Open(dataFile)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	buff := bufio.NewReader(f)
	for {
		var wire Wire
		strWire, err := buff.ReadString('\n')
		if err != nil {
			break
		}
		strWire = strings.TrimSuffix(strWire, "\n")
		lists := strings.Split(strWire, ",")
		sp = Point{0, 0}
		steps = 0
		for _, l := range lists {
			num, err := strconv.ParseInt(l[1:], 10, 64)
			if err != nil {
				log.Fatal(err)
				os.Exit(2)
			}
			ep := NewPoint(sp, string(l[0]), num)
			line := NewLine(sp, ep, steps)
			wire.Lines = append(wire.Lines, line)
			sp = ep
			steps += line.Len
		}
		wires = append(wires, wire)
	}
	return wires, err
}

func calculateCross(wire_a Wire, wire_b Wire) []Point {
	var result []Point

	for _, l_a := range wire_a.Lines {
		for _, l_b := range wire_b.Lines {
			points := Cross(l_a, l_b)
			if points != nil {
				fmt.Println("l_a:", l_a)
				fmt.Println("l_b:", l_b)
				for _, p := range points {
					fmt.Println(p)
				}
			}
			result = append(result, points...)
		}
	}
	return result
}

func Cross(l Line, t Line) []Point {
	var result []Point

	switch l.Dir {
	case "horizontal":
		if t.Dir == "horizontal" {
			if l.A.Y == t.A.Y {
				/* Make sure l is the left line. */
				if l.A.X > t.A.X {
					temp := t
					t = l
					l = temp
				}

				if l.A.X <= t.A.X && l.B.X >= t.A.X &&
					t.A.X <= l.B.X && t.B.X >= l.B.X {
					result = append(result, t.A, l.B)
				}
				if l.A.X <= t.A.X && l.B.X >= t.B.X {
					result = append(result, t.A, t.B)
				}
			}
		} else {
			if l.A.Y >= t.A.Y && l.A.Y <= t.B.Y &&
				t.A.X >= l.A.X && t.A.X <= l.B.X {
				point := Point{
					X: t.A.X,
					Y: l.A.Y,
				}
				result = append(result, point)
			}
		}
	case "vertical":
		if t.Dir == "vertical" {
			if l.A.X == t.A.X {
				/* Make sure l is the left line. */
				if l.A.Y > t.A.Y {
					temp := t
					t = l
					l = temp
				}

				if l.A.Y <= t.A.Y && l.B.Y >= t.A.Y &&
					t.A.Y <= l.B.Y && t.B.Y >= l.B.Y {
					result = append(result, t.A, l.B)
				}
				if l.A.Y <= t.A.Y && l.B.Y >= t.B.Y {
					result = append(result, t.A, t.B)
				}
			}
		} else {
			if l.A.X >= t.A.X && l.A.X <= t.B.X &&
				t.A.Y >= l.A.Y && t.A.Y <= l.B.Y {
				point := Point{
					X: l.A.X,
					Y: t.A.Y,
				}
				result = append(result, point)
			}
		}
	}
	return result
}

func mahattan(start Point, end Point) int64 {
	var abs int64

	if start.X > end.X {
		abs += start.X - end.X
	} else {
		abs += end.X - start.X
	}

	if start.Y > end.Y {
		abs += start.Y - end.Y
	} else {
		abs += end.Y - start.Y
	}
	return abs
}

func calculateSteps(wire Wire, p Point) int64 {
	for _, l := range wire.Lines {
		if l.Contains(p) {
			return l.StepsTake(p)
		}
	}
	return 0
}
