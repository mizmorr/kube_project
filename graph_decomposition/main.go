package main

import (
	"decomposition/core"
	"fmt"
)

func main() {
	names := map[int64]string{0: "zero", 1: "first", 2: "second", 3: "third", 4: "fourth", 5: "fifth", 6: "sixth", 7: "seventh", 8: "eighth", 9: "nineth"}
	new_map := make(map[int64]string, 0)
	m := core.Get_cohesion()
	for k, v := range m {
		for _, elem := range v {
			new_map[elem] = names[k]
		}
	}
	fmt.Println(new_map)
}
