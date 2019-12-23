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

func main() {
	var dataFile string

	flag.StringVarP(&dataFile, "data file name", "f", "", "")
	flag.Parse()

	width := 25
	height := 6
	frames, err := buildFrames(width, height, dataFile)
	if err != nil {
		log.Fatal("Failed to build frames from input file!", err)
	}
	counter := make([][3]int, len(frames))
	for i, f := range frames {
		for _, bit := range f {
			counter[i][bit] += 1
		}
	}

	min_zero_count := width * height
	min_zero_frame := 0
	for i, c := range counter {
		if c[0] < min_zero_count {
			min_zero_count = c[0]
			min_zero_frame = i
		}
	}
	fmt.Println("result=", counter[min_zero_frame][1]*counter[min_zero_frame][2])

	final_frame := make([]int, width*height)
	for i := 0; i < len(final_frame); i++ {
		final_frame[i] = lookDown(frames, i)
	}

	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			v := final_frame[i*width+j]
			if v == 0 {
				fmt.Printf(" ")
			} else if v == 1 {
				fmt.Printf("M")
			}
		}
		fmt.Printf("\n")
	}
}

func lookDown(frames [][]int, pos int) int {
	for _, f := range frames {
		switch f[pos] {
		case 0:
			return 0
		case 1:
			return 1
		}
	}
	log.Fatal("transparent at pos", pos)
	return 2
}

func buildFrames(width int, height int, dataFile string) ([][]int, error) {
	var pixels []int
	var frames [][]int

	pixels_in_frame := width * height

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
		for _, l := range line {
			num, err := strconv.ParseInt(string(l), 10, 64)
			if err != nil {
				log.Fatal(err)
				os.Exit(2)
			}
			pixels = append(pixels, int(num))
		}
	}
	num_frames := len(pixels) / pixels_in_frame
	frames = make([][]int, num_frames)
	for i := 0; i < num_frames; i++ {
		frames[i] = pixels[i*pixels_in_frame : (i+1)*pixels_in_frame]
	}
	return frames, err
}
