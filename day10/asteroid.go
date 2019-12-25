package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strings"

	flag "github.com/spf13/pflag"
)

type Position struct {
	Y        int
	X        int
	Distance float64
}

type TanList struct {
	Atan  float64
	Plist []Position
}

func main() {
	var dataFile string

	flag.StringVarP(&dataFile, "data file name", "f", "", "")
	flag.Parse()

	asteriods, err := buildMap(dataFile)
	if err != nil {
		log.Fatal("Failed to get data from input file!", err)
	}

	detectedMap := detectAsteriods(asteriods)
	max_record := 0
	var max_pos Position
	for y, d := range detectedMap {
		for x, a := range d {
			if a > max_record {
				max_record = a
				max_pos = Position{
					Y: y,
					X: x,
				}
			}
		}
	}
	fmt.Printf("max=%v, pos=%v\n", max_record, max_pos)

	amap := make(map[float64]*TanList)
	for y := 0; y < len(asteriods); y++ {
		for x := 0; x < len(asteriods[0]); x++ {
			if (y != max_pos.Y || x != max_pos.X) && asteriods[y][x] == 1 {
				atan, distance := calculateRelativePos(max_pos, y, x)
				fmt.Printf("atan(%v, %v) = %v\n", y, x, atan)
				p := Position{
					Y:        y,
					X:        x,
					Distance: distance,
				}
				if amap[atan] == nil {
					tanList := TanList{
						Atan: atan,
					}
					amap[atan] = &tanList
				}
				t := amap[atan]
				t.Plist = append(t.Plist, p)
			}
		}
	}

	var tlist []*TanList
	for _, t := range amap {
		tlist = append(tlist, t)
	}

	sort.Sort(TanSorter(tlist))
	counter := 0
	for _, t := range tlist {
		sort.Sort(PositionSorter(t.Plist))
		fmt.Println(t)
		counter += len(t.Plist)
	}

	for i := 0; i < counter; {
		for _, t := range tlist {
			if len(t.Plist) > 0 {
				fmt.Printf("[%v]:destroying: (%v,%v)\n", i, t.Plist[0].Y, t.Plist[0].X)
				t.Plist = t.Plist[1:]
				i++
			}
		}
	}
}

type TanSorter []*TanList

func (s TanSorter) Len() int {
	return len(s)
}

func (s TanSorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]

}

func (s TanSorter) Less(i, j int) bool {
	return s[i].Atan < s[j].Atan
}

type PositionSorter []Position

func (s PositionSorter) Len() int {
	return len(s)
}

func (s PositionSorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]

}

func (s PositionSorter) Less(i, j int) bool {
	return s[i].Distance < s[j].Distance
}

func calculateRelativePos(base Position, y int, x int) (float64, float64) {
	dy := -(y - base.Y)
	dx := x - base.X
	atan := math.Atan2(float64(dx), float64(dy))
	fmt.Printf("atan(%v, %v)=%v ", dx, dy, atan)
	if atan < 0 {
		atan = 2*math.Pi + atan
	}
	distance := math.Sqrt(math.Pow(float64(dy), 2) + math.Pow(float64(dx), 2))
	return atan, distance
}

func buildMap(dataFile string) ([][]int, error) {
	var asteriods [][]int

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
		l := make([]int, len(line))
		for i, c := range line {
			switch c {
			case '.':
				l[i] = 0
			case '#':
				l[i] = 1
			}
		}
		asteriods = append(asteriods, l)
	}
	return asteriods, err
}

func detectAsteriods(asteriods [][]int) [][]int {
	var result [][]int
	result = make([][]int, len(asteriods))
	for i := 0; i < len(result); i++ {
		result[i] = make([]int, len(asteriods[0]))
	}

	for y := 0; y < len(asteriods); y++ {
		for x := 0; x < len(asteriods[0]); x++ {
			if asteriods[y][x] == 1 {
				result[y][x] = lineOfSight(asteriods, y, x)
			}
		}
	}
	return result
}

func lineOfSight(asteriods [][]int, base_y int, base_x int) int {
	counter := 0
	for y := 0; y < len(asteriods); y++ {
		for x := 0; x < len(asteriods[0]); x++ {
			if asteriods[y][x] == 1 {
				if y != base_y || x != base_x {
					if blocked(asteriods, base_y, base_x, y, x) == false {
						counter++
					}
				}
			}
		}
	}
	return counter
}

func blocked(asteriods [][]int, base_y int, base_x int, y int, x int) bool {
	dy := y - base_y
	dx := x - base_x
	gcd := GCD(dy, dx)
	if gcd < 0 {
		gcd *= -1
	}
	if gcd == 1 {
		return false
	}

	step_y := dy / gcd
	step_x := dx / gcd
	for i := 1; i < gcd; i++ {
		inter_y := y - step_y*i
		inter_x := x - step_x*i
		if asteroidExist(asteriods, inter_y, inter_x) {
			return true
		}
	}
	return false
}

func asteroidExist(asteriods [][]int, y int, x int) bool {
	len_y := len(asteriods)
	len_x := len(asteriods[0])
	if y < 0 || y >= len_y || x < 0 || x >= len_x {
		return false
	}
	if asteriods[y][x] == 1 {
		return true
	} else {
		return false
	}
}

func GCD(a int, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}
