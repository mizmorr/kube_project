package core

import (
	"bufio"
	"decomposition/analysis"
	"decomposition/maps"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	mutex = sync.RWMutex{}
)

type Graph interface {
	add_Edge()
	k_core_label()
	append_node()
	insert_node()
	remove_node(id int64)
	get_new_id() int64
}
type Node struct {
	ID          int64
	adj_list    map[int64]int64
	edge_number int64
}
type Undirected_Graph struct {
	nodes map[int64]*Node
}
type Core_set struct {
	numbers map[int64]int64
	value   int64
}

func Get_set(set *Core_set) map[int64]int64 {
	return set.numbers
}
func (g *Undirected_Graph) add_Edge(source, dest int64) {
	if ok := search_map(g.nodes[source].adj_list, dest); !ok {
		mutex.Lock()
		g.nodes[source].adj_list[int64(len(g.nodes[source].adj_list))] = dest
		g.nodes[dest].adj_list[int64(len(g.nodes[dest].adj_list))] = source
		g.nodes[source].edge_number++
		g.nodes[dest].edge_number++
		mutex.Unlock()
	}
}
func pop_front(s []*Node) (int64, []*Node, error) {
	if len(s) == 0 {
		return -1, nil, fmt.Errorf("bad")
	}
	return s[0].ID, s[1:], nil
}
func pop(s []int64) (int64, []int64, error) {
	if len(s) == 0 {
		return -1, nil, fmt.Errorf("bad")
	}
	return s[0], s[1:], nil
}
func (g *Undirected_Graph) get_new_id() int64 {
	if len(g.nodes) == 0 {
		return 0
	}
	return g.nodes[int64(len(g.nodes)-1)].ID + 1
}
func (set *Core_set) get_max() int64 {
	var wg sync.WaitGroup
	var max int64
	wg.Add(3)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for i := int64(0); i < int64(len(set.numbers)/3); i++ {
			if set.numbers[i] > max {
				max = set.numbers[i]
			}
		}
	}(&wg)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for i := int64(len(set.numbers) / 3); i < int64(2*len(set.numbers)/3); i++ {
			if set.numbers[i] > max {
				max = set.numbers[i]
			}
		}
	}(&wg)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for i := int64(2 * len(set.numbers) / 3); i < int64(len(set.numbers)); i++ {
			if set.numbers[i] > max {
				max = set.numbers[i]
			}
		}
	}(&wg)
	wg.Wait()
	return max
}

func (g *Undirected_Graph) k_core_label() *Core_set {
	core_set := Core_set{numbers: make(map[int64]int64, len(g.nodes))}
	var (
		k, v int64
		err  error
	)
	for len(g.nodes) > 0 {
		k++
		unprocessed := maps.Keys_ordered(g.nodes)
		for len(unprocessed) > 0 {
			v, unprocessed, err = pop(unprocessed)
			if err == nil && g.nodes[v] != nil {
				if g.nodes[v].edge_number < k {
					for _, node := range g.nodes[v].adj_list {
						if g.nodes[node] != nil {
							g.nodes[node].edge_number--
							unprocessed = append(unprocessed, g.nodes[node].ID)
						}
					}
					core_set.numbers[v] = k - 1
					core_set.value += k - 1
					g.remove_node(v)
				}
			}
		}

	}
	return &core_set
}

func Samples(num, opt int) (string, string) {
	g := NewUndirectedGraph()
	paths := []string{"samples/graph1.txt", "samples/last_fm.txt", "samples/git.txt"}
	separator := []string{" ", ",", ","}
	filename := flag.String("filename", paths[num], "The file to parse")
	flag.Parse()
	var result []int64
	if *filename == "" {
		log.Fatal("Provide a file to parse")
	}

	names := make(chan string)
	readerr := make(chan error)
	go GetLine(*filename, names, readerr)

	for name := range names {
		k := strings.Split(name, separator[num])
		for _, num := range k {
			postj, err := strconv.Atoi(num)
			if err != nil {
				panic(err)
			}
			result = append(result, int64(postj))
		}

	}
	if err := <-readerr; err != nil {
		log.Fatal(err)
	}
	adjes := map[int64]map[int64]int64{}
	for _, node := range result {
		if _, value := adjes[node]; !value {
			adjes[node] = make(map[int64]int64, 0)
		}
	}
	for i := int64(0); i < int64(len(result)-1); i += 2 {
		adjes[result[i]][int64(len(adjes[result[i]]))] = result[i+1]
		adjes[result[i+1]][int64(len(adjes[result[i+1]]))] = result[i]
	}
	for i := int64(0); i < int64(len(adjes)); i++ {
		g.append_node(adjes[i])
	}
	for _, node := range g.nodes {
		node.edge_number = int64(len(node.adj_list))
	}
	t := time.Now()
	set := g.k_core_label()
	res_time := fmt.Sprint(time.Since(t).Seconds())
	switch opt {
	case 1:
		return res_time, fmt.Sprint(set.get_max())
	case 0:
		res := ""
		for k, v := range set.numbers {
			res += fmt.Sprint(k) + ":" + fmt.Sprint(v) + "\n"
		}
		return res_time, res
	}
	return "", ""

}
func (g *Undirected_Graph) remove_node(id int64) {

	for _, from := range g.nodes[id].adj_list {
		if g.nodes[from] != nil {
			adj := maps.Values2(g.nodes[from].adj_list)
			id_adj := interpolation_search(adj, 0, int64(len(adj)-1), id)
			delete(g.nodes[from].adj_list, id_adj)
		}
	}
	delete(g.nodes, id)

}

func (g *Undirected_Graph) append_node(ids map[int64]int64) {
	g.nodes[g.get_new_id()] = &Node{g.get_new_id(), ids, int64(len(ids))}
}

func (g *Undirected_Graph) Get_node(index int64) *Node {
	return g.nodes[index]
}
func (g *Undirected_Graph) Insert_node(adj map[int64]int64, ID int64) {
	g.nodes[ID] = &Node{ID, adj, int64(len(adj))}
}
func NewUndirectedGraph() *Undirected_Graph {
	return &Undirected_Graph{nodes: make(map[int64]*Node, 0)}
}

func (g *Undirected_Graph) print() {
	for _, node := range g.nodes {
		fmt.Printf("id - %v, adj - %v, edge num - %v\n", node.ID, node.adj_list, node.edge_number)
		time.Sleep(10000)
	}
}

func interpolation_search(arr []int64, low, high, search int64) int64 {

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
			return interpolation_search(arr, pos+1, high, search)
		case arr[pos] > search:
			return interpolation_search(arr, low, pos-1, search)
		}
	}
	return -1
}

func GetLine(filename string, names chan string, readerr chan error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		names <- scanner.Text()
	}
	close(names) // close causes range on channel to break out of loop
	readerr <- scanner.Err()
}
func Random_test(nodes_num int64, prob float64, opt int) (string, string, map[int64]int64, int64) {
	g := random_graph(nodes_num, prob)
	t := time.Now()
	set := g.k_core_label()
	analysis.Research(set.numbers, set.value)

	res_time := fmt.Sprint(time.Since(t).Seconds())
	result := ""
	switch opt {
	case 1:
		return res_time, fmt.Sprint(set.get_max()), set.numbers, set.value
	case 0:
		// for k, v := range set.numbers {
		// 	result += fmt.Sprint(k) + ":" + fmt.Sprint(v) + "\n"
		// }
		return res_time, result, set.numbers, set.value
	}
	return "", "", map[int64]int64{}, int64(0)
}
func random_graph(nodes_num int64, prob float64) *Undirected_Graph {
	g := NewUndirectedGraph()
	for i := int64(0); i < nodes_num; i++ {
		g.append_node(map[int64]int64{})
	}
	// var wg sync.WaitGroup
	// wg.Add(2)
	// go func(wg *sync.WaitGroup) {
	// 	defer wg.Done()

	// }(&wg)
	// go func(wg *sync.WaitGroup) {
	// 	defer wg.Done()
	// 	for i := nodes_num / 2; i < nodes_num-1; i++ {
	// 		for j := i + 1; j < nodes_num; j++ {
	// 			if rand.Float64() < prob {
	// 				g.add_Edge(i, j)
	// 			}
	// 		}
	// 	}
	// }(&wg)
	// wg.Wait()
	for i := int64(0); i < nodes_num/2; i++ {
		for j := i + 1; j < nodes_num; j++ {
			if rand.Float64() < prob {
				g.add_Edge(i, j)
			}
		}
	}
	return g
}
func search_map(m map[int64]int64, search int64) bool {
	is_searched := false
	var wg sync.WaitGroup
	wg.Add(2)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for i := int64(0); i < int64(len(m)/2); i++ {
			if m[i] == search {
				is_searched = true
			}
		}
	}(&wg)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for i := int64(len(m) / 2); i < int64(len(m)); i++ {
			if m[i] == search {
				is_searched = true
			}
		}
	}(&wg)
	wg.Wait()
	return is_searched

}
func TestRand() {
	g := random_graph(10, 0.4)
	g.print()
}

func Get_git() *Undirected_Graph {
	g := NewUndirectedGraph()
	var result []int64
	var filename string = "samples/git.txt"
	names := make(chan string)
	readerr := make(chan error)
	go GetLine(filename, names, readerr)

	for name := range names {
		k := strings.Split(name, ",")
		for _, num := range k {
			postj, err := strconv.Atoi(num)
			if err != nil {
				panic(err)
			}
			result = append(result, int64(postj))
		}

	}
	if err := <-readerr; err != nil {
		log.Fatal(err)
	}
	adjes := map[int64]map[int64]int64{}
	for _, node := range result {
		if _, value := adjes[node]; !value {
			adjes[node] = make(map[int64]int64, 0)
		}
	}
	for i := int64(0); i < int64(len(result)-1); i += 2 {
		adjes[result[i]][int64(len(adjes[result[i]]))] = result[i+1]
		adjes[result[i+1]][int64(len(adjes[result[i+1]]))] = result[i]
	}
	for i := int64(0); i < int64(len(adjes)); i++ {
		g.append_node(adjes[i])
	}
	for _, node := range g.nodes {
		node.edge_number = int64(len(node.adj_list))
	}
	return g
}
func Make_k_decomp(g *Undirected_Graph) *Core_set {
	set := g.k_core_label()

	return set
}
