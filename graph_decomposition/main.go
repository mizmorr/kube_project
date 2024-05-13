package main

import (
	"decomposition/core"
	"fmt"
	"time"
)

func main() {
	t := time.Now()
	core.Get_git()
	fmt.Println(time.Since(t).Seconds())
}
