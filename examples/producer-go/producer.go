package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	for {
		fmt.Println("Running since", time.Now().Sub(start))
		time.Sleep(time.Second * 4)
	}
}
