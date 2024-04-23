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
func get_maxes(m map[int64]int64) {

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
func Get_cohesion() {

	git := Get_git()
	root_name := "./git_sets/"
	f, _ := os.ReadDir(root_name)
	git_label := get_label_from_txt()
	result_path := root_name + "git_cohesion.txt"
	file, err := os.OpenFile(result_path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer file.Close()
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
		file.WriteString(fmt.Sprintf("set %d - %#v\n", index, maxes))
	}

}
