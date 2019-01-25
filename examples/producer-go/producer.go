package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("A producer written in Go!")

	start := time.Now()
	c := time.Tick(5 * time.Second)
	for now := range c {
		fmt.Println("running since", now.Sub(start))
	}
}
