package main

import "fmt"

// Declare a type constraint as an interface
type Number interface {
	int64 | float64
}

func main() {
	// Initialize a map for the integer values
	ints := map[string]int64{
		"first":  34,
		"second": 12,
	}

	// Initialize a map for the float values
	floats := map[string]float64{
		"first":  35.98,
		"second": 26.99,
	}

	fmt.Printf("Non-Generic Sums: %v and %v\n",
		SumInts(ints),
		SumFloats(floats))

	fmt.Printf("Generic Sums: %v and %v\n",
		SumIntsOrFloats[int64](ints),
		SumIntsOrFloats[float64](floats))
}

func SumInts(m map[string]int64) int64 {
	var sum int64

	for _, v := range m {
		sum += v
	}

	return sum
}

func SumFloats(m map[string]float64) float64 {
	var sum float64

	for _, v := range m {
		sum += v
	}

	return sum
}

func SumIntsOrFloats[V Number](m map[string]V) V {
	var sum V

	for _, v := range m {
		sum += v
	}

	return sum
}
