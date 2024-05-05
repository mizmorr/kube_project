package maps

import (
	"sort"
)

// Keys returns the keys of the map m.
// The keys will be in an indeterminate order.
func Keys[M ~map[K]V, K comparable, V any](m M) []K {
	r := make([]K, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	return r
}

func Keys_ordered[M ~map[K]V, K int64, V any](m M) []K {
	r := make([]K, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	sort.Slice(r, func(i, j int) bool {
		return r[i] < r[j]
	})
	return r
}

// Values returns the values of the map m.
// The values will be in an indeterminate order.
func Values[M ~map[K]V, K comparable, V any](m M) []V {
	r := make([]V, 0, len(m))
	for _, v := range m {
		r = append(r, v)
	}
	return r
}

func Values2[M ~map[int64]V, V any](m M) []V {
	r := make([]V, len(m))
	for i := int64(0); i < int64(len(m)); i++ {
		r[i] = m[i]
	}
	return r
}

func Pop[T any](stack *[]T) {
	if len(*stack) == 0 {
		return
	}
	*stack = (*stack)[:len(*stack)-1]
}
func Pop_front[T any](stack *[]T) T {
	if len(*stack) == 0 {
		var NULL T
		return NULL
	} else {
		poped := (*stack)[0]
		*stack = (*stack)[1:]
		return poped
	}

}

func get_map_maxes(m map[int64]int64, index, current int, maxes map[int64]int64) map[int64]int64 {
	if len(m) == 0 || index == current {
		return maxes
	}
	var key_max int64
	for key := range m {
		key_max = key
		break
	}
	for key := range m {
		if m[key] > m[key_max] {
			key_max = key
		}
	}
	maxes[key_max] = m[key_max]
	if index != 1 {
		delete(m, key_max)
	}
	current++
	return get_map_maxes(m, index, current, maxes)

}
func Get_max_of_max(m map[int64]int64, number_of_maxes int) map[int64]int64 {
	return get_map_maxes(m, number_of_maxes, 0, map[int64]int64{})
}

func Interpolation_search(arr []int64, low, high, search int64) int64 {

	if low <= high && search >= arr[low] && search <= arr[high] {

		if arr[high]-arr[low] == 0 {
			switch {
			case arr[len(arr)-1] == search:
				return int64(len(arr) - 1)
			default:
				return -1
			}
		}
		pos := low + (((high - low) / (arr[high] - arr[low])) * (search - arr[low]))
		switch {
		case arr[pos] == search:
			return pos
		case arr[pos] < search:
			return Interpolation_search(arr, pos+1, high, search)
		case arr[pos] > search:
			return Interpolation_search(arr, low, pos-1, search)
		}
	}
	return -1
}
