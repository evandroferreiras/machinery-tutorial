package main

import (
	"fmt"

	machinery "github.com/evandroferreiras/machinery-tutorial/machinery"
)

func main() {
	fmt.Println("Starting the application...")
	err := machinery.NewBuilder().Do()
	if err != nil {
		panic(err)
	}
	fmt.Println("Terminating the application...")
}
