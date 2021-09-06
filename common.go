package main

func IBound(v, min, max int) int {
	if min > max {
		return 0
	}
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func IntAbs(a int) int {
	if a >= 0 {
		return a
	}
	return a * -1
}
