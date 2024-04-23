package analysis

import (
	"fmt"
	"math"
)

type Element struct {
	ID      int64
	Count   int64
	Percent float64
}

func Research(core_set map[int64]int64, sum int64) (string, map[int64]*Element) {
	expected_val, varience := (float64(sum) / float64(len(core_set))), 0.0
	for _, val := range core_set {
		varience += math.Pow(float64(val)-expected_val, 2)
	}
	varience /= float64(len(core_set))
	standard_deviation := math.Sqrt(varience)
	distribution := make(map[int64]*Element)
	for _, val := range core_set {
		if _, ok := distribution[val]; !ok {
			distribution[val] = &Element{ID: val, Count: 1}
		} else {
			distribution[val].Count++
		}
	}
	var modal int64
	for _, val := range distribution {
		if val.Count > modal {
			modal = val.Count
		}
	}
	for _, val := range distribution {
		val.Percent = float64(val.Count) / float64(len(core_set))
	}
	for _, val := range distribution {
		fmt.Println(val)
	}
	list := []string{"математическое ожидание значения номера k-ядра: ", "дисперсия значений номеров k-ядер: ", "среднеквадратичное отклонение значений номеров k-ядер: ", "мода значений номеров k-ядер: "}
	result := fmt.Sprintln(list[0], expected_val, "\n", list[1], varience, "\n", list[2], standard_deviation, "\n", list[3], modal)
	return result, distribution
}
