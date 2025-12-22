package main

import (
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/MarcinCiura/BiDoboSort/bidobo"
)

func main() {
	N := 100_000
	for range 100_000 {
		T := []uint32{}
		for i := range N {
			T = append(T, uint32(i))
		}
		rand.Shuffle(N, func(a, b int) { T[a], T[b] = T[b], T[a] })
		h := []int{1, 4, 4, 8, 16}
		p := 1.35
		for h[len(h)-1] < N {
			h = append(h, int(p*float64(h[len(h)-1])))
		}
		start := time.Now()
		bidobo.BiDoboSort(T, h)
		first := time.Since(start)
		bidobo.InsertionSort(T)
		second := time.Since(start)
		fmt.Printf("%v %v %v %v %f\n", second, first, second-first, h[:6], p)
	}
}
