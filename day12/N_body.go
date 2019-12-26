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

const N = 4

type State struct {
	xState AxisState
	yState AxisState
	zState AxisState
}

type AxisState struct {
	pos [N]int
	v   [N]int
}

func main() {
	var dataFile string

	flag.StringVarP(&dataFile, "data file name", "f", "", "")
	flag.Parse()

	xState, yState, zState, err := buildMap(dataFile)
	if err != nil {
		log.Fatal("Failed to get data from input file!", err)
	}

	xchan := make(chan int)
	go simulateAxis(xState, xchan)

	ychan := make(chan int)
	go simulateAxis(yState, ychan)

	zchan := make(chan int)
	go simulateAxis(zState, zchan)

	x_repeat := <-xchan
	y_repeat := <-ychan
	z_repeat := <-zchan

	fmt.Printf("x, y, z repeats at (%v,%v,%v)\n", x_repeat, y_repeat, z_repeat)

	result := LCM(x_repeat, y_repeat, z_repeat)
	fmt.Println("result=", result)
}

func GCD(a, b int) int {
	for b != 0 {
		temp := b
		b = a % b
		a = temp
	}
	return a
}

func LCM(a, b int, integers ...int) int {
	result := a * b / GCD(a, b)
	for i := 0; i < len(integers); i++ {
		result = LCM(result, integers[i])
	}
	return result
}

func simulateAxis(astat AxisState, out chan int) {
	stats := make(map[AxisState]bool)
	for i := 0; ; i++ {
		astat = updateState(astat)
		if stats[astat] == true {
			out <- i
			break
		}
		stats[astat] = true
	}
}

func updateState(astat AxisState) AxisState {
	astat = updateVelocity(astat)
	astat = updatePosition(astat)
	return astat
}

func updatePosition(astat AxisState) AxisState {
	for i := 0; i < N; i++ {
		astat.pos[i] += astat.v[i]
	}
	return astat
}

func updateVelocity(astat AxisState) AxisState {
	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			if j != i {
				a := pointEffect(astat.pos[i], astat.pos[j])
				astat.v[i] += a
			}
		}
	}
	return astat
}

func pointEffect(target int, source int) int {
	if target > source {
		return -1
	} else if target < source {
		return 1
	} else {
		return 0
	}
}

func buildMap(dataFile string) (AxisState, AxisState, AxisState, error) {
	var xState AxisState
	var yState AxisState
	var zState AxisState

	f, err := os.Open(dataFile)
	if err != nil {
		return xState, yState, zState, err
	}

	defer f.Close()

	buff := bufio.NewReader(f)
	for i := 0; i < N; i++ {
		line, err := buff.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSuffix(line, ">\n")
		line = strings.TrimPrefix(line, "<")
		axies := strings.Split(line, ",")
		for _, b := range axies {
			b := strings.TrimPrefix(b, " ")
			vars := strings.Split(b, "=")
			num, err := strconv.Atoi(vars[1])
			if err != nil {
				log.Fatal(err)
			}
			switch vars[0] {
			case "x":
				xState.pos[i] = num
			case "y":
				yState.pos[i] = num
			case "z":
				zState.pos[i] = num
			}
		}
	}
	return xState, yState, zState, err
}
