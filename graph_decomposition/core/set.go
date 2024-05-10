package core

import (
	"decomposition/maps"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func Make_sort() {
	git := Get_git()
	set := Make_k_decomp(git)
	for key, value := range set.numbers {
		file_path := fmt.Sprintf("./git_sets/set_%d.txt", value)

		f, err := os.OpenFile(file_path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err.Error())
		}
		defer f.Close()
		f.WriteString(fmt.Sprintf("%d\n", key))
	}
}
func make_graph(g *Undirected_Graph, nodes []int64) *Undirected_Graph {
	k := NewUndirectedGraph()
	for _, n := range nodes {
		node := g.Get_node(n)
		k.Insert_node(node.adj_list, node.ID)
	}
	return k

}
func check_edges(g *Undirected_Graph, git_label map[int64]int64) map[int64]int64 {
	rating := make(map[int64]int64, 0)
	for _, node := range g.nodes {
		for _, edge := range node.adj_list {
			if _, ok := rating[git_label[edge]]; !ok {
				rating[git_label[edge]] = 1
			} else {
				rating[git_label[edge]] += 1
			}
		}
	}
	maxes := maps.Get_max_of_max(rating, 3)
	return maxes
}
func get_maxes(mp map[int64]int64) {
	// var filename string = "./git_sets/git_cohesion.txt"
	// m := make(map[int64]int64)
	// names := make(chan string)
	// readerr := make(chan error)
	// go GetLine(filename, names, readerr)
	// for name := range names {
	// 	str_slice := strings.Split(name, ",")
	// 	key, _ := strconv.ParseInt(str_slice[0], 10, 64)
	// 	val, _ := strconv.ParseInt(str_slice[1], 10, 64)
	// 	m[key] = val
	// }
	// return m
}

func make_git_label() {
	git_set := Make_k_decomp(Get_git())

	f, err := os.OpenFile("./git_sets/git_label.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer f.Close()
	for key, val := range git_set.numbers {
		f.WriteString(fmt.Sprintf("%d,%d\n", key, val))
	}
}
func get_label_from_txt() map[int64]int64 {
	var filename string = "./git_sets/git_label.txt"
	m := make(map[int64]int64)
	names := make(chan string)
	readerr := make(chan error)
	go GetLine(filename, names, readerr)
	for name := range names {
		str_slice := strings.Split(name, ",")
		key, _ := strconv.ParseInt(str_slice[0], 10, 64)
		val, _ := strconv.ParseInt(str_slice[1], 10, 64)
		m[key] = val
	}
	return m
}
func Get_cohesion() map[int64][]int64 {

	root_name := "./git_sets/"
	f, _ := os.ReadDir(root_name)
	result_path := root_name + "git_cohesion.txt"
	if _, err := os.Stat(result_path); err == nil {
		dir, _ := os.Getwd()
		// fmt.Println("exists", dir+"/git_sets/git_cohesion.txt")
		err := os.Remove(dir + "/git_sets/git_cohesion.txt")

		if err != nil {
			panic(err)
		}
	}
	git := Get_git()

	git_label := get_label_from_txt()

	file, err := os.OpenFile(result_path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err.Error())
	}
	var key, value int64 = 1, -1
	defer file.Close()
	new_map := []map[int64]int64{}
	for index := 1; index < len(f); index++ {
		filename := root_name + "set_" + fmt.Sprint(index) + ".txt"
		names := make(chan string)
		readerr := make(chan error)
		go GetLine(filename, names, readerr)
		nodes := []int64{}
		for k := range names {
			cur, err := strconv.Atoi(k)
			if err != nil {
				log.Fatal(err.Error())
			} else {
				nodes = append(nodes, int64(cur))
			}
		}
		new_graph := make_graph(git, nodes)
		maxes := check_edges(new_graph, git_label)
		delete(maxes, int64(index))
		if value < maxes[34] {
			value = maxes[34]
			key = int64(index)
		}
		new_map = append(new_map, maxes)
		file.WriteString(fmt.Sprintf("set %d, - %v\n", index, maxes))
	}

	for index, m := range new_map {
		if int64(index) != key-1 {
			delete(m, 34)
		}
		new_map[index] = maps.Get_max_of_max(m, 1)
	}
	process_map := make(map[int64][]int64, 0)
	result_map := make(map[int64][]int64, 0)
	for i, m := range new_map {
		for k := range m {
			var pool int64
			if i == 32 {
				delete(m, 33)
			}
			for i, j := range process_map {
				if ok := maps.Interpolation_search(j, 0, int64(len(j)-1), k); ok != -1 {
					pool = i
					break
				}

			}
			if pool == 0 {
				process_map[int64(len(process_map))] = append(process_map[int64(len(process_map))], k)
				result_map[int64(len(result_map))] = append(result_map[int64(len(result_map))], int64(i+1))
			} else {
				process_map[pool] = append(process_map[pool], k)
				result_map[pool] = append(result_map[pool], int64(i+1))
			}
		}
	}
	// fmt.Println(result_map)
	return result_map

}

func Make_markup() {
	names := map[int64]string{0: "zero", 1: "first", 2: "second", 3: "third", 4: "fourth", 5: "fifth", 6: "sixth", 7: "seventh", 8: "eighth", 9: "nineth"}
	new_map := make(map[int64]string, 0)
	m := Get_cohesion()
	for k, v := range m {
		for _, elem := range v {
			new_map[elem] = names[k]
		}
	}
	f, err := os.Create("markup.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	// fmt.Println(new_map)
	for k, v := range new_map {
		f.WriteString(fmt.Sprintf("%v:%s\n", k, v))
	}
}
