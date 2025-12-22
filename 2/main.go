package main

import (
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/MarcinCiura/BiDoboSort/bidobo"
)

func main() {
	N := 100_000
	for h1 := 2; h1 <= 8; h1++ {
		for h2 := 2; h2 <= 4*h1; h2++ {
			for h3 := 2; h3 <= 4*h2; h3++ {
				for h4 := 2; h4 <= 4*h3; h4++ {
					for p := 1.25; p < 1.4; p += 0.001 {
						h := []int{1, h1, h2, h3, h4}
						for h[len(h)-1] < N {
							h = append(h, int(p*float64(h[len(h)-1])))
						}
						T := []uint32{}
						for i := range N {
							T = append(T, uint32(i))
						}
						rand.Shuffle(N, func(a, b int) { T[a], T[b] = T[b], T[a] })

						start := time.Now()
						bidobo.BiDoboSort(T, h)
						first := time.Since(start)
						bidobo.InsertionSort(T)
						second := time.Since(start)
						fmt.Printf("%v %v %v %v %f\n", second, first, second-first, h[:6], p)
					}
				}
			}
		}
	}
}
