//go:generate go run asm.go -out bidobo.s -stubs stub.go
package bidobo

type Direction int

const (
	UPWARD   Direction = +1
	DOWNWARD Direction = -1
)

func BiDoboSort[E uint32 | uint64](T []E, h []int) {
	i := 0
	for h[i+1] < len(T) {
		i++
	}
	dir := UPWARD
	for ; i > 0; i-- {
		switch dir {
		case UPWARD:
			blockSortUpward(T, len(T)-2*h[i], h[i])
		case DOWNWARD:
			blockSortDownward(T, len(T)-2*h[i], h[i])
		}
		dir = -dir
	}
}

func blockSortUpward[E uint32 | uint64](T []E, end, gap int) {
	i := 0
	switch data := any(T).(type) {
	case []uint32:
		switch {
		case gap >= 8:
			i = blockSortUpwardBy8ElementsOf4Bytes(data, end, gap)
		case gap >= 4:
			i = blockSortUpwardBy4ElementsOf4Bytes(data, end, gap)
		}
	case []uint64:
		if gap >= 4 {
			i = blockSortUpwardBy4ElementsOf8Bytes(data, end, gap)
		}
	}
	for ; i < end; i++ {
		sortTwoElements(&T[i+gap], &T[i+2*gap])
		sortTwoElements(&T[i], &T[i+gap])
	}
}

func blockSortDownward[E uint32 | uint64](T []E, i, gap int) {
	switch data := any(T).(type) {
	case []uint32:
		switch {
		case gap >= 8:
			i = blockSortDownwardBy8ElementsOf4Bytes(data, i, gap)
		case gap >= 4:
			i = blockSortDownwardBy4ElementsOf4Bytes(data, i, gap)
		default:
			i--
		}
	case []uint64:
		switch {
		case gap >= 4:
			i = blockSortDownwardBy4ElementsOf8Bytes(data, i, gap)
		default:
			i--
		}
	}
	for ; i >= 0; i-- {
		sortTwoElements(&T[i], &T[i+gap])
		sortTwoElements(&T[i+gap], &T[i+2*gap])
	}
}

func sortTwoElements[E uint32 | uint64](a, b *E) {
	if *a > *b {
		*a, *b = *b, *a
	}
}

func InsertionSort[E uint32 | uint64](T []E) {
	for i := 1; i < len(T); i++ {
		ti := T[i]
		j := i - 1
		for ; j >= 0 && ti < T[j]; j-- {
			T[j+1] = T[j]
		}
		T[j+1] = ti
	}
}
