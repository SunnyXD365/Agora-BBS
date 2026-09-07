package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Temporal Worker initial running...")
	for {
		time.Sleep(10 * time.Second)
	}
}
