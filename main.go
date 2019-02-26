package main

import (
	"fmt"

	machinery "github.com/evandroferreiras/machinery-tutorial/machinery"
)

func main() {
	fmt.Println("Starting the application...")
	err := machinery.NewBbroker: 'redis://localhost:6379'
	default_queue: machinery_tasks
	
	result_backend: 'redis://localhost:6379'
	results_expire_in: 3600000uilder().Do()
	if err != nil {
		panic(err)
	}
	fmt.Println("Terminating the application...")
}
