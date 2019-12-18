package main

import "fmt"

const Min = 134792
const Max = 675810

func main() {
	var result1 []uint
	var result2 []uint

	for i := Min; i <= Max; i++ {
		digits := calculateDigits(uint(i))
		if filterRisingOrder(digits) && filterRepeat(digits) {
			result1 = append(result1, uint(i))
		}

		if filterRisingOrder(digits) && filterRepeatTwo(digits) {
			result2 = append(result2, uint(i))
		}
	}
	for _, r := range result2 {
		fmt.Printf("%v ", r)
	}
	fmt.Println(len(result2))
}

/* This function returns a slice of digits with the least
   significant number in pos 0*/
func calculateDigits(data uint) []uint {
	var digits []uint
	var code uint
	code = 10

	for data >= code {
		digit := data % code
		digits = append(digits, digit)
		data = (data - digit) / code
	}
	digits = append(digits, data)
	return digits
}

func filterRisingOrder(digits []uint) bool {
	last := digits[0]
	for _, digit := range digits[1:] {
		if digit > last {
			return false
		}
		last = digit
	}
	return true
}

func filterRepeat(digits []uint) bool {
	last := digits[0]
	for _, digit := range digits[1:] {
		if digit == last {
			return true
		}
		last = digit
	}
	return false
}

func filterRepeatTwo(digits []uint) bool {
	groupSize := 1
	last := digits[0]
	for _, digit := range digits[1:] {
		if digit == last {
			groupSize += 1
		} else {
			if groupSize == 2 {
				return true
			}
			groupSize = 1
		}
		last = digit
	}

	if groupSize == 2 {
		return true
	}
	return false
}
