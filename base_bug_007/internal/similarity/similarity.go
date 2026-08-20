package similarity

import "strings"

func Score(a, b string) float64 {
	x, y := map[string]bool{}, map[string]bool{}
	for _, v := range strings.Fields(a) {
		x[v] = true
	}
	for _, v := range strings.Fields(b) {
		y[v] = true
	}
	if len(x) == 0 && len(y) == 0 {
		return 1
	}
	common := 0
	for v := range y {
		if x[v] {
			common++
		}
	}
	return float64(common) / float64(len(x)+len(y)-common)
}
func Similar(a, b string) bool { return Score(a, b) >= .65 }
