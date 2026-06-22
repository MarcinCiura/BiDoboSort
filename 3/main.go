package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/MarcinCiura/BiDoboSort/bidobo"
)

const (
	arraySize  = 1_000_000
	iterations = 10_000
)

func main() {
	// Initialize a local random generator
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		// Find the boundaries of the gap sequence array [ ... ]
		startIdx := strings.Index(line, "[")
		endIdx := strings.Index(line, "]")
		if startIdx == -1 || endIdx == -1 || startIdx >= endIdx {
			log.Printf("Skipping invalid line (missing brackets): %s\n", line)
			continue
		}

		// Extract the sequence from inside the brackets
		baseGapsStr := line[startIdx+1 : endIdx]
		parts := strings.Fields(baseGapsStr)
		var h []int
		for _, pStr := range parts {
			val, err := strconv.Atoi(pStr)
			if err != nil {
				log.Printf("Failed to parse gap %q: %v\n", pStr, err)
				continue
			}
			h = append(h, val)
		}

		// Extract the p ratio which is at the end of the line
		ratioStr := strings.TrimSpace(line[endIdx+1:])
		pRatio, err := strconv.ParseFloat(ratioStr, 64)
		if err != nil {
			log.Printf("Failed to parse ratio %q: %v\n", ratioStr, err)
			continue
		}

		// Reconstruct the full gap sequence
		// The condition next < arraySize ensures we don't generate gaps larger than the array
		for {
			last := h[len(h)-1]
			next := int(1.0 + pRatio*float64(last))
			if next >= arraySize {
				break
			}
			h = append(h, next)
		}

		// Slices to store execution times (in milliseconds)
		var totalTimes []float64
		var biDoboTimes []float64
		var insertionTimes []float64

		// Run the benchmark 1000 times
		for i := 0; i < iterations; i++ {
			// Generate an array of |arraySize| unique, shuffled elements
			Tint := rng.Perm(arraySize)
			T := []uint32{}
			for _, x := range Tint {
			       T = append(T, uint32(x))
			}

			// Benchmark BiDoboSort
			startBiDobo := time.Now()
			bidobo.BiDoboSort(T, h) // T is partially/fully sorted in-place
			elapsedBiDobo := float64(time.Since(startBiDobo).Microseconds()) / 1000.0

			// Benchmark InsertionSort
			startInsertion := time.Now()
			bidobo.InsertionSort(T) // Final pass to ensure it's completely sorted
			elapsedInsertion := float64(time.Since(startInsertion).Microseconds()) / 1000.0

			// Record the times
			biDoboTimes = append(biDoboTimes, elapsedBiDobo)
			insertionTimes = append(insertionTimes, elapsedInsertion)
			totalTimes = append(totalTimes, elapsedBiDobo+elapsedInsertion)
		}

		// Calculate statistics
		meanTotal, stdTotal := calculateStats(totalTimes)
		meanBi, stdBi := calculateStats(biDoboTimes)
		meanIns, stdIns := calculateStats(insertionTimes)

		// Print the results
		fmt.Printf("Base gaps: %v ratio: %f", parts, pRatio)
		fmt.Printf("  Total time:      mean %.6f ms, stddev %.6f ms", meanTotal, stdTotal)
		fmt.Printf("  BiDoboSort time: mean %.6f ms, stddev %.6f ms", meanBi, stdBi)
		fmt.Printf("  InsertSort time: mean %.6f ms, stddev %.6f ms\n", meanIns, stdIns)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading standard input: %v", err)
	}
}

// calculateStats computes the mean and standard deviation of a float64 slice.
func calculateStats(data []float64) (mean float64, stddev float64) {
	if len(data) == 0 {
		return 0, 0
	}

	var sum float64
	for _, v := range data {
		sum += v
	}
	mean = sum / float64(len(data))

	var variance float64
	for _, v := range data {
		variance += (v - mean) * (v - mean)
	}
	stddev = math.Sqrt(variance / float64(len(data)))

	return mean, stddev
}
